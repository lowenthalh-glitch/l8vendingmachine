/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
function initializeDashboard() {
    var statsContainer = document.getElementById('dashboard-stats');
    if (!statsContainer) return;

    var config = Layer8DConfig.getConfig();
    var prefix = (config && config.app && config.app.apiPrefix) || '/vend';

    // Fetch fleet machines
    var machineQuery = encodeURIComponent(JSON.stringify({ text: 'select * from VendFleetMachine' }));
    fetch(prefix + '/10/Machine?body=' + machineQuery, {
        method: 'GET',
        headers: typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
    })
    .then(function(r) { return r.json(); })
    .then(function(data) {
        var machines = data.list || [];
        var totalMachines = (data.metadata && data.metadata.keyCount && data.metadata.keyCount.counts)
            ? (data.metadata.keyCount.counts.Total || machines.length)
            : machines.length;
        var online = 0, warning = 0, offline = 0;
        var totalSlots = 0, emptySlots = 0, lowStockSlots = 0;

        machines.forEach(function(m) {
            if (m.status === 'online') online++;
            else if (m.status === 'warning') warning++;
            else if (m.status === 'offline') offline++;
            totalSlots += m.totalSlots || 0;
            emptySlots += m.emptySlots || 0;
            lowStockSlots += m.lowStockSlots || 0;
        });

        statsContainer.innerHTML =
            kpiCard('🏭', 'Total Machines', totalMachines, '') +
            kpiCard('✅', 'Online', online, 'layer8d-status-active') +
            kpiCard('⚠️', 'Warning', warning, 'layer8d-status-pending') +
            kpiCard('❌', 'Offline', offline, 'layer8d-status-inactive') +
            kpiCard('📦', 'Total Slots', totalSlots, '') +
            kpiCard('🔴', 'Empty Slots', emptySlots, emptySlots > 0 ? 'layer8d-status-terminated' : '') +
            kpiCard('🟡', 'Low Stock Slots', lowStockSlots, lowStockSlots > 0 ? 'layer8d-status-pending' : '');

        renderRestockPriority(machines);
    })
    .catch(function() {
        statsContainer.innerHTML = '<p style="color: var(--layer8d-text-muted);">Unable to load dashboard data.</p>';
    });

}

function kpiCard(icon, label, value, badgeClass) {
    var badge = badgeClass ? ' layer8d-status-badge ' + badgeClass : '';
    return '<div style="background: var(--layer8d-bg-white); border-radius: 12px; padding: 20px; text-align: center; box-shadow: 0 1px 3px rgba(0,0,0,0.08);">' +
        '<div style="font-size: 28px; margin-bottom: 8px;">' + icon + '</div>' +
        '<div style="font-size: 32px; font-weight: 700; color: var(--layer8d-text-dark);">' +
        (badgeClass ? '<span class="' + badge + '">' + value + '</span>' : value) +
        '</div>' +
        '<div style="font-size: 13px; color: var(--layer8d-text-muted); margin-top: 4px;">' + label + '</div>' +
        '</div>';
}

function renderRestockPriority(machines) {
    var container = document.getElementById('dashboard-restock');
    if (!container || !machines || machines.length === 0) return;

    // Calculate revenue for all machines, find top 10% threshold
    var allWithRevenue = [];
    machines.forEach(function(m) {
        var revenue = VendInventoryUtils.calcRevenue(m.inventory);
        var pct = VendInventoryUtils.calcFillPct(m.inventory);
        allWithRevenue.push({ name: m.name || m.machineId, pct: pct, revenue: revenue });
    });
    allWithRevenue.sort(function(a, b) { return b.revenue - a.revenue; });

    // Top 10% revenue threshold
    var top10pctCount = Math.max(1, Math.ceil(allWithRevenue.length * 0.1));
    var revenueThreshold = allWithRevenue.length >= top10pctCount
        ? allWithRevenue[top10pctCount - 1].revenue : 0;

    // Filter: only top 10% revenue machines that need restock (fill < 60%)
    var needsRestock = allWithRevenue.filter(function(m) {
        return m.revenue >= revenueThreshold && m.pct >= 0 && m.pct < 60;
    });

    if (needsRestock.length === 0) {
        container.innerHTML = '<div style="padding: 16px; background: var(--layer8d-bg-light); border-radius: 8px; color: var(--layer8d-text-muted);">All top revenue machines adequately stocked</div>';
        return;
    }

    // Already sorted by revenue desc
    var top = needsRestock.slice(0, 10);

    var html = '<div style="background: var(--layer8d-bg-white); border-radius: 12px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.08);">';
    html += '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">';
    html += '<h3 style="margin: 0; color: var(--layer8d-text-dark);">Restock Priority</h3>';
    html += '<span style="font-size: 13px; color: var(--layer8d-text-muted);">' + needsRestock.length + ' machines need restock (sorted by revenue)</span>';
    html += '</div>';

    top.forEach(function(m) {
        html += '<div style="display: flex; align-items: center; gap: 12px; padding: 8px 0; border-bottom: 1px solid var(--layer8d-border);">';
        html += '<span style="min-width: 160px; font-size: 13px; color: var(--layer8d-text-dark); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">' + m.name + '</span>';
        html += '<div style="flex: 1;">' + VendInventoryUtils.fillBar(m.pct) + '</div>';
        html += '<span style="min-width: 70px; text-align: right; font-size: 12px; font-weight: 600; color: var(--layer8d-text-dark);">' + VendInventoryUtils.formatMoney(m.revenue) + '</span>';
        html += '</div>';
    });

    if (needsRestock.length > 10) {
        html += '<div style="text-align: center; padding-top: 12px;">';
        html += '<a href="#" onclick="loadSection(\'fleet\'); return false;" style="color: var(--layer8d-primary); font-size: 13px; text-decoration: none;">View all in Fleet &rarr;</a>';
        html += '</div>';
    }
    html += '</div>';
    container.innerHTML = html;
}

