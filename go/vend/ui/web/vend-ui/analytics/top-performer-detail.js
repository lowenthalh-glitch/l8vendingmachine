/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */

// Custom detail popup for Top Performers — adds a daily revenue chart tab.
(function() {
    'use strict';

    var origInit = null;

    function hookTopPerformerDetail() {
        if (!window.Analytics || !Analytics._showDetailsModal) return;

        var origShowDetails = Analytics._showDetailsModal;
        Analytics._showDetailsModal = function(service, item, id) {
            if (service.model === 'VendTopPerformer' && item) {
                showTopPerformerPopup(service, item);
            } else {
                origShowDetails.call(Analytics, service, item, id);
            }
        };
    }

    function showTopPerformerPopup(service, item) {
        var machineName = item.machineName || item.machineId || 'Machine';
        var machineId = item.machineId || item.performerId;

        // Build tabbed content
        var content = '<div class="probler-popup-tabs">' +
            '<div class="probler-popup-tab active" data-tab="chart">Revenue Chart</div>' +
            '<div class="probler-popup-tab" data-tab="details">Details</div>' +
            '</div>' +
            '<div class="probler-popup-tab-content">' +
            '<div class="probler-popup-tab-pane active" data-pane="chart">' +
            '<div id="top-perf-chart" style="padding: 16px;"><div style="color: var(--layer8d-text-muted);">Loading chart...</div></div>' +
            '</div>' +
            '<div class="probler-popup-tab-pane" data-pane="details">' +
            '<div id="top-perf-details" style="padding: 16px;"></div>' +
            '</div>' +
            '</div>';

        Layer8DPopup.show({
            title: machineName,
            content: content,
            size: 'large',
            onShow: function() {
                // Render chart in the chart tab
                loadMachineRevenueChart(machineId, machineName);

                // Render form in the details tab
                var detailsContainer = document.getElementById('top-perf-details');
                if (detailsContainer && window.Layer8DForms) {
                    var formDef = AnalyticsData.forms.VendTopPerformer;
                    if (formDef) {
                        detailsContainer.innerHTML = Layer8DForms.generateFormHtml(formDef, item);
                    }
                }
            }
        });
    }

    function loadMachineRevenueChart(machineId, machineName) {
        var container = document.getElementById('top-perf-chart');
        if (!container) return;

        var config = Layer8DConfig.getConfig();
        var prefix = (config && config.app && config.app.apiPrefix) || '/vend';

        var query = encodeURIComponent(JSON.stringify({
            text: 'select * from VendInventorySnapshot where machineId=' + machineId + ' sort-by timestamp limit 500'
        }));

        fetch(prefix + '/10/InvSnap?body=' + query, {
            method: 'GET',
            headers: typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
        })
        .then(function(r) { return r.json(); })
        .then(function(data) {
            var snapshots = data.list || [];
            if (snapshots.length === 0) {
                container.innerHTML = '<div style="padding: 20px; color: var(--layer8d-text-muted);">No historical data available for this machine.</div>';
                return;
            }
            renderRevenueChart(container, snapshots, machineName);
        })
        .catch(function() {
            container.innerHTML = '<div style="padding: 20px; color: var(--layer8d-error);">Failed to load chart data.</div>';
        });
    }

    function renderRevenueChart(container, snapshots, machineName) {
        // Group by day, sum dailyRevenue
        var days = {};
        snapshots.forEach(function(s) {
            var ts = s.timestamp;
            if (typeof ts === 'string') ts = parseInt(ts);
            if (!ts) return;
            var day = Math.floor(ts / 86400) * 86400;
            if (!days[day]) days[day] = { revenue: 0, fill: 0, count: 0 };
            var rev = s.dailyRevenue || s.revenue || 0;
            if (typeof rev === 'string') rev = parseInt(rev);
            days[day].revenue = Math.max(days[day].revenue, rev); // daily max (cumulative)
            days[day].fill += (s.fillPct || 0);
            days[day].count++;
        });

        var dayKeys = Object.keys(days).map(Number).sort(function(a, b) { return a - b; });
        if (dayKeys.length === 0) {
            container.innerHTML = '<div style="padding: 20px; color: var(--layer8d-text-muted);">No data points.</div>';
            return;
        }

        var w = container.offsetWidth || 600;
        var h = 280;
        var pad = { top: 20, right: 20, bottom: 40, left: 60 };
        var plotW = w - pad.left - pad.right;
        var plotH = h - pad.top - pad.bottom;

        var maxRev = 0;
        dayKeys.forEach(function(d) {
            var val = days[d].revenue / 100; // cents to dollars
            if (val > maxRev) maxRev = val;
        });
        if (maxRev === 0) maxRev = 1;

        var minDay = dayKeys[0];
        var maxDay = dayKeys[dayKeys.length - 1];
        var dayRange = maxDay - minDay || 86400;

        var svg = '<svg width="' + w + '" height="' + h + '" xmlns="http://www.w3.org/2000/svg">';

        // Grid
        var tickStep = maxRev / 4;
        for (var ti = 0; ti <= 4; ti++) {
            var tickVal = Math.round(ti * tickStep);
            var y = pad.top + plotH - (tickVal / maxRev) * plotH;
            svg += '<line x1="' + pad.left + '" y1="' + y + '" x2="' + (w - pad.right) + '" y2="' + y + '" stroke="var(--layer8d-border)" stroke-dasharray="3,3" opacity="0.5"/>';
            svg += '<text x="' + (pad.left - 8) + '" y="' + (y + 4) + '" text-anchor="end" fill="var(--layer8d-text-muted)" font-size="11">$' + tickVal + '</text>';
        }

        // X-axis labels
        var labelStep = Math.max(1, Math.floor(dayKeys.length / 7));
        dayKeys.forEach(function(day, i) {
            if (i % labelStep !== 0 && i !== dayKeys.length - 1) return;
            var x = pad.left + ((day - minDay) / dayRange) * plotW;
            var d = new Date(day * 1000);
            svg += '<text x="' + x + '" y="' + (h - pad.bottom + 16) + '" text-anchor="middle" fill="var(--layer8d-text-muted)" font-size="11">' + (d.getMonth() + 1) + '/' + d.getDate() + '</text>';
        });

        // Revenue line
        var points = [];
        dayKeys.forEach(function(day) {
            var val = days[day].revenue / 100;
            var x = pad.left + ((day - minDay) / dayRange) * plotW;
            var y = pad.top + plotH - (val / maxRev) * plotH;
            points.push(x + ',' + y);
        });

        if (points.length > 1) {
            svg += '<polyline points="' + points.join(' ') + '" fill="none" stroke="var(--layer8d-success)" stroke-width="2.5" stroke-linejoin="round"/>';
        }

        // Data points
        points.forEach(function(p) {
            var coords = p.split(',');
            svg += '<circle cx="' + coords[0] + '" cy="' + coords[1] + '" r="3" fill="var(--layer8d-success)"/>';
        });

        svg += '</svg>';

        container.innerHTML = '<h4 style="margin: 0 0 12px 16px; color: var(--layer8d-text-dark);">Daily Revenue — ' + machineName + '</h4>' + svg;
    }

    // Hook after analytics initializes
    var waitCount = 0;
    function tryHook() {
        if (window.Analytics && Analytics._showDetailsModal) {
            hookTopPerformerDetail();
        } else if (waitCount < 20) {
            waitCount++;
            setTimeout(tryHook, 500);
        }
    }
    tryHook();
})();
