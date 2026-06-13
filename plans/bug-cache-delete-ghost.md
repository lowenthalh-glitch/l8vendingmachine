# Bug Report: Cache Delete Causes Ghost Missing Entry

## Summary

After deleting one entity from a CRUD service, a subsequent `select *` query returns N-2 items instead of N-1. One unrelated entity silently disappears from the result set. The entity still exists in the cache and can be found via a filtered query, but the unfiltered list skips it.

## Severity

**High** — data silently goes missing from list views after any delete. No error, no warning.

## Reproduction

1. Create 5 entities (e.g., 5 VendRoute via the route optimizer)
2. Delete 1 entity via the UI
3. Query `select * from VendRoute` — returns 3 instead of 4
4. Metadata says `Total: 5` (stale count) but only 3 items in the list
5. Query `select * from VendRoute where driverId=drv-002` — returns the "missing" entity (it exists)
6. Restart the service — the list returns the correct 4 items

**Reproduced in:**
- `l8vendingmachine` — VendRoute (delete 1 of 5, see 3)
- `l8physio` — therapist list (delete 1 of 3, see 1)

## Root Cause

**File:** `l8utils/go/utils/cache/internalQuery.go`, line 52

```go
func (this *internalQuery) prepare(cache map[string]interface{}, addedOrder []string, stamp int64, ...) {
    this.stamp = stamp
    data := make([]string, 0)

    if addedOrder != nil {
        data = addedOrder    // <-- BUG: shares the slice, does not copy
    } else {
        // ...
    }
    // ...
    this.data = data
}
```

When a `select *` query has no criteria, `prepare()` sets `data = addedOrder` — a **direct reference** to `internalCache.addedOrder`, not a copy. This means `internalQuery.data` and `internalCache.addedOrder` point to the **same underlying array**.

When `internalCache.delete()` runs (in `internalCache.go`, line 136):

```go
func (this *internalCache) delete(pk, uk string) (interface{}, bool) {
    // ...
    if idx, exists := this.key2order[pk]; exists {
        this.addedOrder[idx] = ""    // tombstone
        delete(this.key2order, pk)
        this.deleteCount++
    }
    this.stamp = time.Now().Unix()
    this.cleanupOrder()
    // ...
}
```

The tombstone (`""`) is written into `addedOrder[idx]`. Since the query's `data` slice shares the same backing array, the tombstone is now visible in the query's data too.

On the next `fetch()`, the stamp mismatch triggers `prepare()` again, which reassigns `data = addedOrder` (now the compacted slice from `cleanupOrder`). However, `cleanupOrder` only fires when tombstones exceed 25% of the slice or 100 count. For small lists (3-5 items), the threshold isn't met, so the tombstone persists.

In `fetch()` (line 209-214):

```go
for i := start; i < len(dq.data); i++ {
    key := dq.data[i]
    value, ok := this.cache[key]
    if ok {
        result = append(result, value)
    }
}
```

The empty string tombstone causes `this.cache[""]` → `!ok` → silently skipped. This accounts for the deleted entry. But the second missing entry is caused by the `sort.Slice` in `prepare()` (line 64) sorting the shared slice — the empty tombstone sorts to the beginning, shifting all indices, while `key2order` still has the old indices. This index mismatch causes one real entry to be unreachable.

## Fix

**File:** `l8utils/go/utils/cache/internalQuery.go`, line 51-52

Replace:
```go
if addedOrder != nil {
    data = addedOrder
}
```

With:
```go
if addedOrder != nil {
    data = make([]string, 0, len(addedOrder))
    for _, k := range addedOrder {
        if k != "" {
            data = append(data, k)
        }
    }
}
```

This copies the slice (breaking the shared reference) and filters out tombstones, ensuring the query data is always clean regardless of the `cleanupOrder` threshold.

## Why Not a Workaround

Per the `report-infra-bugs.md` rule, this is a framework bug in `l8utils` that affects all CRUD services across all L8 projects. Workarounds in consuming projects (e.g., forcing a page reload after delete) mask the real problem and don't fix it for other projects.
