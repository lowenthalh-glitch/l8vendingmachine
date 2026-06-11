/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

function showVendMachineDetail(item) {
    var esc = Layer8DUtils.escapeHtml;
    var machines = Array.isArray(item.machines) ? item.machines : Object.values(item.machines || {});

    var title = '<span style="font-weight:600">Management System</span> &mdash; ' + esc(item.machineId);

    var content = '<div class="probler-popup-tabs">' +
        '<div class="probler-popup-tab active" data-tab="machines">Vending Machines (' + machines.length + ')</div>' +
    '</div>' +
    '<div class="probler-popup-tab-content">' +
        '<div class="probler-popup-tab-pane active" data-pane="machines">' +
            buildMachinesTable(machines, esc) +
        '</div>' +
    '</div>';

    Layer8DPopup.show({
        titleHtml: title,
        content: content,
        size: 'xlarge',
        showFooter: false
    });
}

function buildMachinesTable(machines, esc) {
    if (!machines || machines.length === 0) {
        return '<p style="padding: 16px; color: var(--layer8d-text-muted);">No machines collected yet.</p>';
    }

    var statusBadge = function(status) {
        var cls = '';
        if (status === 'online') cls = 'layer8d-status-active';
        else if (status === 'offline') cls = 'layer8d-status-inactive';
        else if (status === 'warning') cls = 'layer8d-status-pending';
        return '<span class="layer8d-status-badge ' + cls + '">' + esc(status || 'Unknown') + '</span>';
    };

    var typeLabel = function(t) {
        var map = { 'vending_snack': 'Snack', 'vending_drink': 'Drink', 'vending_combo': 'Combo',
            'coffee': 'Coffee', 'ev_charger': 'EV Charger', 'laundry': 'Laundry', 'car_wash': 'Car Wash' };
        return map[t] || t || '';
    };

    var html = '<table class="layer8d-table" style="width:100%">' +
        '<thead><tr>' +
        '<th>Machine ID</th><th>Name</th><th>Type</th><th>Model</th>' +
        '<th>Status</th><th>City</th><th>State</th><th>Daily TXN</th><th>Device</th>' +
        '</tr></thead><tbody>';

    for (var i = 0; i < machines.length; i++) {
        var m = machines[i];
        html += '<tr>' +
            '<td>' + esc(m.machineId || '') + '</td>' +
            '<td>' + esc(m.name || '') + '</td>' +
            '<td>' + typeLabel(m.type) + '</td>' +
            '<td>' + esc(m.model || '') + '</td>' +
            '<td>' + statusBadge(m.status) + '</td>' +
            '<td>' + esc(m.locationCity || '') + '</td>' +
            '<td>' + esc(m.locationState || '') + '</td>' +
            '<td>' + (m.dailyTransactions || 0) + '</td>' +
            '<td>' + esc(m.deviceId || '') + '</td>' +
            '</tr>';
    }

    html += '</tbody></table>';
    return html;
}
