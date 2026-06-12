/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';

    // Austin TX center
    var CENTER = [30.2672, -97.7431];
    var ZOOM = 12;
    var map = null;
    var leafletLoaded = false;

    function loadLeaflet(callback) {
        if (leafletLoaded) { callback(); return; }
        // CSS
        var link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
        document.head.appendChild(link);
        // JS
        var script = document.createElement('script');
        script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js';
        script.onload = function() { leafletLoaded = true; callback(); };
        document.head.appendChild(script);
    }

    function getHeaders() {
        return {
            'Authorization': 'Bearer ' + sessionStorage.bearerToken,
            'Content-Type': 'application/json'
        };
    }

    function query(endpoint, model, callback) {
        var q = 'select * from ' + model;
        var body = encodeURIComponent(JSON.stringify({ text: q }));
        var prefix = Layer8DConfig.getApiPrefix();
        fetch(prefix + endpoint + '?body=' + body, { method: 'GET', headers: getHeaders() })
            .then(function(r) { return r.json(); })
            .then(function(data) { callback(data.list || []); })
            .catch(function() { callback([]); });
    }

    function colorIcon(color) {
        return L.divIcon({
            className: '',
            html: '<div style="width:14px;height:14px;background:' + color +
                  ';border:2px solid #fff;border-radius:50%;box-shadow:0 1px 4px rgba(0,0,0,0.3);"></div>',
            iconSize: [14, 14],
            iconAnchor: [7, 7],
            popupAnchor: [0, -10]
        });
    }

    function addFacilities() {
        query('/10/Facility', 'VendStockingFacility', function(items) {
            items.forEach(function(f) {
                if (!f.coordinates) return;
                var lat = f.coordinates.latitude;
                var lng = f.coordinates.longitude;
                if (!lat && !lng) return;
                var stockCount = (f.stock || []).length;
                L.marker([lat, lng], { icon: colorIcon('#e74c3c') })
                    .addTo(map)
                    .bindPopup('<b>' + esc(f.name) + '</b> (' + esc(f.code) + ')<br>' +
                        esc(addrLine(f.address)) + '<br>' +
                        'Storage: ' + (f.totalStorageSqFt || 0) + ' sq ft<br>' +
                        'Docks: ' + (f.loadingDocks || 0) + ' | Products: ' + stockCount);
            });
        });
    }

    function addLocations() {
        query('/10/Location', 'VendLocation', function(items) {
            items.forEach(function(loc) {
                if (!loc.coordinates) return;
                var lat = loc.coordinates.latitude;
                var lng = loc.coordinates.longitude;
                if (!lat && !lng) return;
                L.marker([lat, lng], { icon: colorIcon('#3498db') })
                    .addTo(map)
                    .bindPopup('<b>' + esc(loc.name) + '</b><br>' +
                        'Type: ' + esc(loc.locationType || '') + '<br>' +
                        'Contact: ' + esc(loc.contactName || ''));
            });
        });
    }

    function addTrucks() {
        query('/10/Truck', 'VendDeliveryTruck', function(items) {
            items.forEach(function(t) {
                var lat = t.lastLatitude;
                var lng = t.lastLongitude;
                if (!lat && !lng) return;
                if (typeof lat === 'string') lat = parseFloat(lat);
                if (typeof lng === 'string') lng = parseFloat(lng);
                var statusMap = {0:'Unknown',1:'Active',2:'Maintenance',3:'En Route',4:'Decommissioned'};
                L.marker([lat, lng], { icon: colorIcon('#2ecc71') })
                    .addTo(map)
                    .bindPopup('<b>' + esc(t.name) + '</b><br>' +
                        esc(t.make + ' ' + t.model) + '<br>' +
                        'Plate: ' + esc(t.plateNumber) + '<br>' +
                        'Status: ' + (statusMap[t.status] || t.status) + '<br>' +
                        'MPG: ' + (t.milesPerGallon ? t.milesPerGallon.toFixed(1) : '-'));
            });
        });
    }

    function addDrivers() {
        query('/10/Driver', 'VendDriver', function(items) {
            items.forEach(function(d) {
                var lat = d.currentLatitude;
                var lng = d.currentLongitude;
                if (typeof lat === 'string') lat = parseFloat(lat);
                if (typeof lng === 'string') lng = parseFloat(lng);
                if (!lat && !lng) return;
                var ago = '';
                if (d.lastLocationUpdate) {
                    var ts = typeof d.lastLocationUpdate === 'string' ? parseInt(d.lastLocationUpdate) : d.lastLocationUpdate;
                    var mins = Math.round((Date.now() / 1000 - ts) / 60);
                    ago = mins < 60 ? mins + 'm ago' : Math.round(mins / 60) + 'h ago';
                }
                L.marker([lat, lng], { icon: colorIcon('#f39c12') })
                    .addTo(map)
                    .bindPopup('<b>' + esc(d.firstName + ' ' + d.lastName) + '</b><br>' +
                        esc(addrLine(d.homeAddress)) + '<br>' +
                        'Phone: ' + esc(d.phone || '') + '<br>' +
                        'License: ' + licClass(d.licenseClass) +
                        (ago ? '<br>Updated: ' + ago : ''));
            });
        });
    }

    function addMachines() {
        query('/0/VCache', 'VendMachine', function(items) {
            items.forEach(function(entry) {
                if (!entry.machines) return;
                var machines = entry.machines;
                for (var key in machines) {
                    if (!machines.hasOwnProperty(key)) continue;
                    var m = machines[key];
                    var lat = m.locationLat;
                    var lng = m.locationLng;
                    if (typeof lat === 'string') lat = parseFloat(lat);
                    if (typeof lng === 'string') lng = parseFloat(lng);
                    if (!lat && !lng) continue;
                    var statusLabel = m.status || 'unknown';
                    L.marker([lat, lng], { icon: colorIcon('#9b59b6') })
                        .addTo(map)
                        .bindPopup('<b>' + esc(m.name || m.machineId) + '</b><br>' +
                            'Type: ' + esc(m.type || '') + '<br>' +
                            'Model: ' + esc(m.model || '') + '<br>' +
                            'Status: ' + statusLabel + '<br>' +
                            esc(m.locationAddress || '') + ', ' + esc(m.locationCity || '') + ' ' + esc(m.locationState || ''));
                }
            });
        });
    }

    function addrLine(addr) {
        if (!addr) return '';
        var parts = [addr.line1, addr.city, addr.stateProvince, addr.postalCode].filter(Boolean);
        return parts.join(', ');
    }

    function licClass(v) {
        return {0:'—',1:'Class C',2:'Class B',3:'Class A'}[v] || String(v);
    }

    function esc(s) {
        if (!s) return '';
        var d = document.createElement('div');
        d.textContent = s;
        return d.innerHTML;
    }

    window.initializeMap = function() {
        var container = document.getElementById('vend-map');
        if (!container) return;

        loadLeaflet(function() {
            if (map) { map.remove(); map = null; }

            map = L.map('vend-map').setView(CENTER, ZOOM);
            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
                attribution: '&copy; OpenStreetMap contributors',
                maxZoom: 19
            }).addTo(map);

            addFacilities();
            addLocations();
            addMachines();
            addTrucks();
            addDrivers();
        });
    };
})();
