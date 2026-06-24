/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';

    var CENTER = [30.2672, -97.7431];
    var ZOOM = 12;
    var map = null;
    var leafletLoaded = false;
    var layers = {};
    var rawData = { facilities: [], locations: [], machines: [], trucks: [], drivers: [], routes: [] };
    var machineCoords = {};
    var facCoords = {};

    function loadLeaflet(cb) {
        if (leafletLoaded) { cb(); return; }
        var link = document.createElement('link');
        link.rel = 'stylesheet'; link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
        document.head.appendChild(link);
        var s = document.createElement('script');
        s.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js';
        s.onload = function() { leafletLoaded = true; cb(); };
        document.head.appendChild(s);
    }

    function getHeaders() {
        return { 'Authorization': 'Bearer ' + sessionStorage.bearerToken, 'Content-Type': 'application/json' };
    }

    function query(endpoint, model, cb) {
        var q = 'select * from ' + model;
        var body = encodeURIComponent(JSON.stringify({ text: q }));
        var prefix = Layer8DConfig.getApiPrefix();
        fetch(prefix + endpoint + '?body=' + body, { method: 'GET', headers: getHeaders() })
            .then(function(r) { return r.json(); })
            .then(function(data) { cb(data.list || []); })
            .catch(function() { cb([]); });
    }

    function colorIcon(color, size) {
        var sz = size || 14;
        return L.divIcon({
            className: '',
            html: '<div style="width:'+sz+'px;height:'+sz+'px;background:'+color+
                  ';border:2px solid #fff;border-radius:50%;box-shadow:0 1px 4px rgba(0,0,0,0.3);"></div>',
            iconSize: [sz, sz], iconAnchor: [sz/2, sz/2], popupAnchor: [0, -sz/2]
        });
    }

    function getOSRMUrl() {
        try {
            var cfg = Layer8DConfig.getConfig();
            return (cfg && cfg.app && cfg.app.osrmUrl) || 'http://localhost:5000';
        } catch(e) { return 'http://localhost:5000'; }
    }

    // Fetch actual road geometry from OSRM for a list of waypoints.
    // Calls back with an array of [lat,lng] points for the full road path, or null on failure.
    function fetchOSRMRoute(waypoints, callback) {
        var coords = waypoints.map(function(p) { return p[1] + ',' + p[0]; }).join(';');
        var url = getOSRMUrl() + '/route/v1/driving/' + coords + '?overview=full&geometries=geojson';
        fetch(url)
            .then(function(r) { return r.json(); })
            .then(function(data) {
                if (data.code === 'Ok' && data.routes && data.routes[0] && data.routes[0].geometry) {
                    var geojsonCoords = data.routes[0].geometry.coordinates;
                    // GeoJSON is [lng,lat], Leaflet needs [lat,lng]
                    var path = geojsonCoords.map(function(c) { return [c[1], c[0]]; });
                    callback(path);
                } else {
                    callback(null);
                }
            })
            .catch(function() { callback(null); });
    }

    // Build a set of machine/facility/driver IDs on a single route's stops
    function buildHighlightSet(routes, routeId) {
        var set = { machines: {}, facilities: {}, driverId: '', truckId: '' };
        for (var i = 0; i < routes.length; i++) {
            if (routes[i].routeId !== routeId) continue;
            set.driverId = routes[i].driverId || '';
            set.truckId = routes[i].vehicleId || '';
            (routes[i].stops || []).forEach(function(s) {
                if (s.machineId) set.machines[s.machineId] = true;
                if (s.facilityId) set.facilities[s.facilityId] = true;
            });
            break;
        }
        return set;
    }

    function esc(s) { if (!s) return ''; var d = document.createElement('div'); d.textContent = s; return d.innerHTML; }
    function addrLine(a) { if (!a) return ''; return [a.line1,a.city,a.stateProvince,a.postalCode].filter(Boolean).join(', '); }
    function licClass(v) { return {0:'—',1:'Class C',2:'Class B',3:'Class A'}[v] || String(v); }
    function num(v) { return typeof v === 'string' ? parseFloat(v) : (v || 0); }
    function matchSearch(name, term) {
        if (!name || !term) return true;
        return name.toLowerCase().indexOf(term.toLowerCase()) !== -1;
    }

    // --- Render functions ---

    function renderFacilities(items, search, highlight) {
        var lg = L.layerGroup();
        items.forEach(function(f) {
            if (!f.coordinates) return;
            var lat = f.coordinates.latitude, lng = f.coordinates.longitude;
            if (!lat && !lng) return;
            facCoords[f.facilityId] = [lat, lng];
            if (search && !matchSearch(f.name, search)) return;
            var onRoute = !highlight || (highlight.facilities[f.facilityId]);
            var size = onRoute ? 18 : 12;
            var opacity = onRoute ? 1.0 : 0.3;
            var marker = L.marker([lat, lng], { icon: colorIcon('#e74c3c', size), opacity: opacity }).addTo(lg)
                .bindPopup('<b>'+esc(f.name)+'</b> ('+esc(f.code)+')<br>'+esc(addrLine(f.address))+
                    '<br>Storage: '+(f.totalStorageSqFt||0)+' sq ft | Docks: '+(f.loadingDocks||0));
        });
        return lg;
    }

    function renderLocations(items, search) {
        var lg = L.layerGroup();
        items.forEach(function(loc) {
            if (!loc.coordinates) return;
            var lat = loc.coordinates.latitude, lng = loc.coordinates.longitude;
            if (!lat && !lng) return;
            if (search && !matchSearch(loc.name, search)) return;
            L.marker([lat, lng], { icon: colorIcon('#3498db') }).addTo(lg)
                .bindPopup('<b>'+esc(loc.name)+'</b><br>Type: '+esc(loc.locationType||''));
        });
        return lg;
    }

    function renderMachines(items, statusFilter, search, highlight) {
        var lg = L.layerGroup();
        items.forEach(function(entry) {
            if (!entry.machines) return;
            for (var key in entry.machines) {
                if (!entry.machines.hasOwnProperty(key)) continue;
                var m = entry.machines[key];
                var mid = m.machineId || key;
                var lat = num(m.locationLat), lng = num(m.locationLng);
                if (!lat && !lng) continue;
                machineCoords[mid] = [lat, lng];
                if (statusFilter === 'needs-restock') {
                    var hasEmpty = false;
                    var slots = m.slots || [];
                    for (var si = 0; si < slots.length; si++) {
                        if (slots[si].currentStock === 0 || slots[si].status === 'empty' || slots[si].status === 'critical') {
                            hasEmpty = true; break;
                        }
                    }
                    if (!hasEmpty && m.status !== 'warning') continue;
                } else if (statusFilter && m.status !== statusFilter) continue;
                if (search && !matchSearch(m.name || mid, search)) continue;
                var onRoute = !highlight || highlight.machines[mid];
                var baseColor = (m.status === 'offline') ? '#95a5a6' : '#9b59b6';
                var size = onRoute ? 16 : 8;
                var opacity = onRoute ? 1.0 : 0.25;
                L.marker([lat, lng], { icon: colorIcon(baseColor, size), opacity: opacity }).addTo(lg)
                    .bindPopup('<b>'+esc(m.name||mid)+'</b><br>Type: '+esc(m.type||'')+
                        '<br>Status: '+esc(m.status||'')+'<br>'+esc(m.locationCity||''));
            }
        });
        return lg;
    }

    function renderTrucks(items, statusFilter, search, highlight) {
        var lg = L.layerGroup();
        var statusMap = {0:'Unknown',1:'Active',2:'Maintenance',3:'En Route',4:'Decommissioned'};
        items.forEach(function(t) {
            var lat = num(t.lastLatitude), lng = num(t.lastLongitude);
            if (!lat && !lng) return;
            if (statusFilter && String(t.status) !== statusFilter) return;
            if (search && !matchSearch(t.name, search)) return;
            var onRoute = !highlight || highlight.truckId === t.truckId;
            var size = onRoute ? 18 : 12;
            var opacity = onRoute ? 1.0 : 0.25;
            L.marker([lat, lng], { icon: colorIcon('#2ecc71', size), opacity: opacity }).addTo(lg)
                .bindPopup('<b>'+esc(t.name)+'</b><br>'+esc(t.make+' '+t.model)+
                    '<br>Plate: '+esc(t.plateNumber)+'<br>Status: '+(statusMap[t.status]||t.status)+
                    '<br>MPG: '+(t.milesPerGallon ? t.milesPerGallon.toFixed(1) : '-'));
        });
        return lg;
    }

    function renderDrivers(items, search, highlight) {
        var lg = L.layerGroup();
        items.forEach(function(d) {
            var lat = num(d.currentLatitude), lng = num(d.currentLongitude);
            if (!lat && !lng) return;
            var name = (d.firstName||'')+' '+(d.lastName||'');
            if (search && !matchSearch(name, search)) return;
            var ago = '';
            if (d.lastLocationUpdate) {
                var ts = num(d.lastLocationUpdate);
                var mins = Math.round((Date.now()/1000 - ts)/60);
                ago = mins < 60 ? mins+'m ago' : Math.round(mins/60)+'h ago';
            }
            var onRoute = !highlight || highlight.driverId === d.driverId;
            var size = onRoute ? 18 : 12;
            var opacity = onRoute ? 1.0 : 0.25;
            L.marker([lat, lng], { icon: colorIcon('#f39c12', size), opacity: opacity }).addTo(lg)
                .bindPopup('<b>'+esc(name)+'</b><br>'+esc(addrLine(d.homeAddress))+
                    '<br>License: '+licClass(d.licenseClass)+(ago ? '<br>Updated: '+ago : ''));
        });
        return lg;
    }

    var routeColors = ['#e74c3c','#3498db','#2ecc71','#f39c12','#9b59b6','#1abc9c','#e67e22','#34495e'];

    function renderRoutes(routes, selectedRouteIds, search) {
        var lg = L.layerGroup();
        routes.forEach(function(route, ri) {
            if (!route.stops || route.stops.length === 0) return;
            if (selectedRouteIds.length > 0 && selectedRouteIds.indexOf(route.routeId) === -1) return;
            if (search && !matchSearch(route.name, search)) return;
            var color = routeColors[ri % routeColors.length];
            var points = [];
            route.stops.forEach(function(stop) {
                var coords = null;
                if (stop.stopType === 'end') {
                    // End-of-day stop — use last known point + offset or skip if no coords
                    // The stop has no machineId/facilityId, coords come from the route data
                    return; // Will be drawn as the polyline endpoint
                } else if (stop.stopType === 'reload' && stop.facilityId) {
                    coords = facCoords[stop.facilityId];
                } else if (stop.machineId) {
                    coords = machineCoords[stop.machineId];
                }
                if (!coords) return;
                points.push(coords);
                var mc = stop.serviceUrgency==='high' ? '#e74c3c' : stop.serviceUrgency==='reload' ? '#3498db' :
                    stop.serviceUrgency==='break' ? '#95a5a6' : '#f1c40f';
                var stopLabel = stop.stopType==='reload' ? 'R' : stop.stopType==='break' ? 'B' : String(stop.stopOrder);
                var numberedIcon = L.divIcon({
                    className: '',
                    html: '<div style="width:22px;height:22px;background:'+color+';border:2px solid #fff;border-radius:50%;' +
                          'box-shadow:0 1px 4px rgba(0,0,0,0.4);display:flex;align-items:center;justify-content:center;' +
                          'font-size:10px;font-weight:bold;color:#fff;">' + stopLabel + '</div>',
                    iconSize: [22, 22], iconAnchor: [11, 11], popupAnchor: [0, -14]
                });
                var popupLabel = stop.stopType==='reload' ? ' (Reload)' : stop.stopType==='break' ? ' (Break)' :
                    ': '+esc(stop.machineName || stop.machineId);
                L.marker(coords, { icon: numberedIcon }).addTo(lg)
                    .bindPopup('<b>'+esc(route.name)+'</b><br>Stop #'+stop.stopOrder+ popupLabel +
                        '<br>'+esc(stop.locationAddress||'')+' '+esc(stop.locationCity||'')+
                        '<br>Urgency: '+esc(stop.serviceUrgency||''));
            });
            if (points.length >= 2) {
                var popupText = '<b>'+esc(route.name)+'</b><br>Distance: '+
                    (route.totalDistance ? route.totalDistance.toFixed(1)+' mi' : '-')+
                    '<br>Duration: '+(route.totalDuration||0)+' min'+
                    '<br>Fuel: $'+(route.estimatedFuelCost ? route.estimatedFuelCost.toFixed(2) : '-');
                // Try OSRM for actual road geometry
                fetchOSRMRoute(points, function(roadPath) {
                    if (roadPath) {
                        L.polyline(roadPath, { color: color, weight: 4, opacity: 0.8 }).addTo(lg)
                            .bindPopup(popupText);
                    } else {
                        // Fallback to straight lines
                        L.polyline(points, { color: color, weight: 3, opacity: 0.7, dashArray: '8,4' }).addTo(lg)
                            .bindPopup(popupText);
                    }
                });
            }
        });
        return lg;
    }

    // --- Refresh map with current filters ---

    function refreshMap() {
        if (!map) return;
        var toggles = getToggles();
        var mf = val('map-machine-status');
        var tf = val('map-truck-status');
        var rf = getSelectedRouteIds();
        var search = val('map-search').trim();

        for (var k in layers) { if (layers[k]) map.removeLayer(layers[k]); }
        layers = {};

        // Always build machineCoords/facCoords even if layer hidden (routes need them)
        machineCoords = {};
        facCoords = {};
        rawData.facilities.forEach(function(f) {
            if (f.coordinates) facCoords[f.facilityId] = [f.coordinates.latitude, f.coordinates.longitude];
        });
        rawData.machines.forEach(function(entry) {
            if (!entry.machines) return;
            for (var key in entry.machines) {
                var m = entry.machines[key];
                var lat = num(m.locationLat), lng = num(m.locationLng);
                if (lat || lng) machineCoords[m.machineId || key] = [lat, lng];
            }
        });

        // Build highlight set: when exactly 1 route selected, highlight its stops
        var highlight = null;
        if (rf.length === 1) {
            highlight = buildHighlightSet(rawData.routes, rf[0]);
        }

        if (toggles.facilities) { layers.facilities = renderFacilities(rawData.facilities, search, highlight); layers.facilities.addTo(map); }
        if (toggles.locations) { layers.locations = renderLocations(rawData.locations, search); layers.locations.addTo(map); }
        if (toggles.machines) { layers.machines = renderMachines(rawData.machines, mf, search, highlight); layers.machines.addTo(map); }
        if (toggles.trucks) { layers.trucks = renderTrucks(rawData.trucks, tf, search, highlight); layers.trucks.addTo(map); }
        if (toggles.drivers) { layers.drivers = renderDrivers(rawData.drivers, search, highlight); layers.drivers.addTo(map); }
        if (toggles.routes) { layers.routes = renderRoutes(rawData.routes, rf, search); layers.routes.addTo(map); }
    }

    function getToggles() {
        var r = {};
        document.querySelectorAll('#map-toggles input[data-layer]').forEach(function(cb) { r[cb.dataset.layer] = cb.checked; });
        return r;
    }

    function val(id) { var el = document.getElementById(id); return el ? el.value : ''; }

    function getSelectedRouteIds() {
        var panel = document.getElementById('map-route-panel');
        if (!panel) return [];
        var ids = [];
        panel.querySelectorAll('input[type="checkbox"]:checked').forEach(function(cb) {
            ids.push(cb.value);
        });
        return ids;
    }

    function populateRouteDropdown() {
        var panel = document.getElementById('map-route-panel');
        var btn = document.getElementById('map-route-btn');
        if (!panel || !btn) return;
        panel.innerHTML = '';

        if (rawData.routes.length === 0) {
            panel.innerHTML = '<div style="padding:6px 12px;font-size:11px;color:var(--layer8d-text-muted);">No routes</div>';
            return;
        }

        // Select All option
        var allLabel = document.createElement('label');
        allLabel.style.cssText = 'display:flex;align-items:center;gap:6px;padding:4px 12px;font-size:11px;cursor:pointer;border-bottom:1px solid var(--layer8d-border);margin-bottom:2px;';
        var allCb = document.createElement('input');
        allCb.type = 'checkbox';
        allCb.checked = true;
        allCb.addEventListener('change', function() {
            panel.querySelectorAll('input.route-cb').forEach(function(cb) { cb.checked = allCb.checked; });
            updateRouteButtonLabel();
            refreshMap();
        });
        allLabel.appendChild(allCb);
        allLabel.appendChild(document.createTextNode('Select All'));
        panel.appendChild(allLabel);

        // Build driver/truck lookup for labels
        var driverMap = {};
        rawData.drivers.forEach(function(d) {
            driverMap[d.driverId] = (d.firstName || '') + ' ' + (d.lastName || '');
        });
        var truckMap = {};
        rawData.trucks.forEach(function(t) {
            truckMap[t.truckId] = t.name || t.plateNumber || t.truckId;
        });

        rawData.routes.forEach(function(r, ri) {
            var driverName = driverMap[r.driverId] || r.driverId || '';
            var truckName = truckMap[r.vehicleId] || r.vehicleId || '';
            var stops = (r.stops || []).length;
            var routeLabel = (r.name || r.routeId) + ' — ' + driverName.trim() + ' / ' + truckName + ' (' + stops + ' stops)';

            var label = document.createElement('label');
            label.style.cssText = 'display:flex;align-items:center;gap:6px;padding:4px 12px;font-size:11px;cursor:pointer;line-height:1.4;';
            label.addEventListener('mouseenter', function() { label.style.background = 'var(--layer8d-bg-light)'; });
            label.addEventListener('mouseleave', function() { label.style.background = ''; });
            var cb = document.createElement('input');
            cb.type = 'checkbox';
            cb.className = 'route-cb';
            cb.value = r.routeId;
            cb.checked = true;
            cb.addEventListener('change', function() { updateRouteButtonLabel(); refreshMap(); });
            var dot = document.createElement('span');
            var color = routeColors[ri % routeColors.length];
            dot.style.cssText = 'display:inline-block;width:8px;height:8px;border-radius:50%;background:' + color + ';flex-shrink:0;';
            label.appendChild(cb);
            label.appendChild(dot);
            label.appendChild(document.createTextNode(routeLabel));
            panel.appendChild(label);
        });

        // Toggle dropdown on button click
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            panel.style.display = panel.style.display === 'none' ? 'block' : 'none';
        });
        // Close on outside click
        document.addEventListener('click', function(e) {
            if (!document.getElementById('map-route-dropdown').contains(e.target)) {
                panel.style.display = 'none';
            }
        });
        // Stop panel clicks from closing
        panel.addEventListener('click', function(e) { e.stopPropagation(); });

        updateRouteButtonLabel();
    }

    function updateRouteButtonLabel() {
        var btn = document.getElementById('map-route-btn');
        if (!btn) return;
        var ids = getSelectedRouteIds();
        var total = rawData.routes.length;
        if (ids.length === 0) btn.textContent = 'Routes (none) ▾';
        else if (ids.length === total) btn.textContent = 'Routes (all) ▾';
        else btn.textContent = 'Routes (' + ids.length + '/' + total + ') ▾';
    }

    function loadAllData(cb) {
        var pending = 6;
        function done() { pending--; if (pending === 0) cb(); }
        query('/10/Facility', 'VendStockingFacility', function(d) { rawData.facilities = d; done(); });
        query('/10/Location', 'VendLocation', function(d) { rawData.locations = d; done(); });
        query('/0/VCache', 'VendMachine', function(d) { rawData.machines = d; done(); });
        query('/10/Truck', 'VendDeliveryTruck', function(d) { rawData.trucks = d; done(); });
        query('/10/Driver', 'VendDriver', function(d) { rawData.drivers = d; done(); });
        query('/10/Route', 'VendRoute', function(d) { rawData.routes = d; done(); });
    }

    function attachFilterListeners() {
        document.querySelectorAll('#map-toggles input[data-layer]').forEach(function(cb) {
            cb.addEventListener('change', refreshMap);
        });
        ['map-machine-status','map-truck-status'].forEach(function(id) {
            var el = document.getElementById(id);
            if (el) el.addEventListener('change', refreshMap);
        });
        var searchInput = document.getElementById('map-search');
        if (searchInput) {
            var debounce = null;
            searchInput.addEventListener('input', function() {
                clearTimeout(debounce);
                debounce = setTimeout(refreshMap, 300);
            });
        }
    }

    window.initializeMap = function() {
        var container = document.getElementById('vend-map');
        if (!container) return;
        loadLeaflet(function() {
            if (map) { map.remove(); map = null; }
            layers = {}; machineCoords = {}; facCoords = {};

            map = L.map('vend-map').setView(CENTER, ZOOM);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '&copy; OpenStreetMap contributors', maxZoom: 19
            }).addTo(map);

            loadAllData(function() {
                populateRouteDropdown();
                refreshMap();
            });
            attachFilterListeners();
        });
    };
})();
