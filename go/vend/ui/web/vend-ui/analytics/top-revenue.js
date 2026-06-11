/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
function renderTopRevenueWidget() {
    // Find or create the container inside the content area
    var container = document.getElementById('analytics-top-revenue');
    if (!container) {
        // Insert before the module tabs or module content
        var contentArea = document.getElementById('content-area');
        if (!contentArea) return;
        var moduleTabs = contentArea.querySelector('.l8-module-tabs');
        if (!moduleTabs) return;
        container = document.createElement('div');
        container.id = 'analytics-top-revenue';
        container.style.cssText = 'padding: 16px 20px;';
        moduleTabs.parentNode.insertBefore(container, moduleTabs);
    }

    var config = Layer8DConfig.getConfig();
    var prefix = (config && config.app && config.app.apiPrefix) || '/vend';

    // Fetch historical snapshots (last 30 days of revenue data)
    var query = encodeURIComponent(JSON.stringify({ text: 'select * from VendInventorySnapshot limit 5000' }));
    fetch(prefix + '/10/InvSnap?body=' + query, {
        method: 'GET',
        headers: typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
    })
    .then(function(r) { return r.json(); })
    .then(function(data) {
        var snapshots = data.list || [];
        if (snapshots.length === 0) {
            container.innerHTML = '<div style="padding: 16px; color: var(--layer8d-text-muted);">No historical data available</div>';
            return;
        }

        // Sum revenue per machine across all snapshots
        var machineRevenue = {};
        var machineNames = {};
        snapshots.forEach(function(s) {
            var id = s.machineId;
            if (!id) return;
            if (!machineRevenue[id]) machineRevenue[id] = 0;
            machineRevenue[id] += s.revenue || 0;
            if (s.machineName) machineNames[id] = s.machineName;
        });

        var ranked = [];
        for (var id in machineRevenue) {
            if (machineRevenue[id] > 0) {
                ranked.push({ name: machineNames[id] || id, revenue: machineRevenue[id] });
            }
        }

        ranked.sort(function(a, b) { return b.revenue - a.revenue; });

        // Top 10%
        var topCount = Math.max(1, Math.ceil(ranked.length * 0.1));
        var top = ranked.slice(0, topCount);
        var maxRevenue = top.length > 0 ? top[0].revenue : 1;

        var html = '<div style="background: var(--layer8d-bg-white); border-radius: 12px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.08);">';
        html += '<h3 style="margin: 0 0 16px 0; color: var(--layer8d-text-dark);">Top Revenue Machines (30-day, top 10%)</h3>';

        top.forEach(function(m, i) {
            var barWidth = Math.round(m.revenue / maxRevenue * 100);
            html += '<div style="display: flex; align-items: center; gap: 12px; padding: 10px 0; border-bottom: 1px solid var(--layer8d-border);">';
            html += '<span style="min-width: 24px; font-size: 14px; font-weight: 700; color: var(--layer8d-text-muted);">#' + (i + 1) + '</span>';
            html += '<span style="min-width: 160px; font-size: 13px; color: var(--layer8d-text-dark); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">' + m.name + '</span>';
            html += '<div style="flex: 1; height: 10px; background: var(--layer8d-bg-light); border-radius: 5px; overflow: hidden;">';
            html += '<div style="height: 100%; width: ' + barWidth + '%; background: var(--layer8d-success); border-radius: 5px;"></div>';
            html += '</div>';
            html += '<span style="min-width: 80px; text-align: right; font-size: 14px; font-weight: 700; color: var(--layer8d-success);">' + VendInventoryUtils.formatMoney(m.revenue) + '</span>';
            html += '</div>';
        });

        html += '</div>';
        container.innerHTML = html;
    })
    .catch(function(err) {
        console.error('Top Revenue widget error:', err);
    });
}
