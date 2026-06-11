/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */

// Register custom view type for multi-series inventory chart
(function() {
    'use strict';
    if (window.Layer8DViewFactory) {
        Layer8DViewFactory.register('inventory-chart', function(options) {
            return {
                init: function() { VendInventoryChart.init(options.containerId); },
                refresh: function() { VendInventoryChart.init(options.containerId); },
                destroy: function() {
                    var c = document.getElementById(options.containerId);
                    if (c) c.innerHTML = '';
                    VendInventoryChart._reset();
                }
            };
        });
    }
})();

window.VendInventoryChart = {
    MAX_MACHINES: 10,
    COLORS: ['#0ea5e9','#22c55e','#f59e0b','#ef4444','#8b5cf6','#ec4899','#14b8a6','#f97316','#6366f1','#84cc16'],
    _containerId: null,
    _machineList: [],
    _series: [],  // { machineId, machineName, snapshots: [], dayData: {}, color }

    _reset: function() {
        this._series = [];
        this._machineList = [];
    },

    init: function(containerId) {
        this._containerId = containerId;
        this._series = [];
        var container = document.getElementById(containerId);
        if (!container) return;

        container.innerHTML = '<div style="padding: 16px; color: var(--layer8d-text-muted);">Loading machines...</div>';

        // Fetch machine list for the dropdowns
        var config = Layer8DConfig.getConfig();
        var prefix = (config && config.app && config.app.apiPrefix) || '/vend';
        var query = encodeURIComponent(JSON.stringify({ text: 'select machineId,name from VendFleetMachine' }));

        fetch(prefix + '/10/Machine?body=' + query, {
            method: 'GET',
            headers: typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
        })
        .then(function(r) { return r.json(); })
        .then(function(data) {
            var machines = data.list || [];
            VendInventoryChart._machineList = machines.map(function(m) {
                return { id: m.machineId, name: m.name || m.machineId };
            });
            VendInventoryChart._buildUI(container);
            // Auto-load first machine
            if (VendInventoryChart._machineList.length > 0) {
                VendInventoryChart._addMachine(VendInventoryChart._machineList[0].id);
            }
        })
        .catch(function() {
            container.innerHTML = '<div style="padding: 16px; color: var(--layer8d-error);">Failed to load machine list</div>';
        });
    },

    _metrics: [
        { key: 'fillPct', label: 'Fill %', unit: '%', max: 100 },
        { key: 'dailyRevenue', label: 'Daily Revenue', unit: '$', max: 0 },
        { key: 'revenue', label: 'Revenue (since restock)', unit: '$', max: 0 },
        { key: 'totalStock', label: 'Total Stock', unit: '', max: 0 },
        { key: 'emptySlots', label: 'Empty Slots', unit: '', max: 0 }
    ],
    _activeMetric: 'fillPct',

    _buildUI: function(container) {
        var html = '<div style="background: var(--layer8d-bg-white); border-radius: 12px; padding: 20px; box-shadow: 0 1px 3px rgba(0,0,0,0.08);">';

        // Title row with metric selector
        html += '<div style="display: flex; align-items: center; gap: 12px; margin-bottom: 16px;">';
        html += '<h3 style="margin: 0; color: var(--layer8d-text-dark);">30 Day History —</h3>';
        html += '<select id="inv-chart-metric-select" style="padding: 6px 10px; border: 1px solid var(--layer8d-border); border-radius: 6px; font-size: 14px; font-weight: 600; background: var(--layer8d-bg-input); color: var(--layer8d-text-dark);">';
        this._metrics.forEach(function(m) {
            html += '<option value="' + m.key + '">' + m.label + '</option>';
        });
        html += '</select>';
        html += '</div>';

        // Controls row
        html += '<div style="display: flex; align-items: center; gap: 12px; margin-bottom: 16px; flex-wrap: wrap;">';
        html += '<select id="inv-chart-machine-select" style="padding: 6px 10px; border: 1px solid var(--layer8d-border); border-radius: 6px; font-size: 13px; background: var(--layer8d-bg-input); color: var(--layer8d-text-dark);">';
        this._machineList.forEach(function(m) {
            html += '<option value="' + m.id + '">' + m.name + '</option>';
        });
        html += '</select>';
        html += '<button id="inv-chart-add-btn" style="padding: 6px 14px; border: none; border-radius: 6px; background: var(--layer8d-primary); color: white; font-size: 13px; cursor: pointer;">Add to Graph</button>';
        html += '<span id="inv-chart-count" style="font-size: 12px; color: var(--layer8d-text-muted);">0/' + this.MAX_MACHINES + ' machines</span>';
        html += '</div>';

        // Chart area
        html += '<div id="inv-chart-svg" style="width: 100%; height: 350px;"></div>';

        // Legend
        html += '<div id="inv-chart-legend" style="margin-top: 12px; display: flex; flex-wrap: wrap; gap: 8px;"></div>';
        html += '</div>';

        container.innerHTML = html;

        // Attach add button handler
        var btn = document.getElementById('inv-chart-add-btn');
        if (btn) {
            btn.addEventListener('click', function() {
                var sel = document.getElementById('inv-chart-machine-select');
                if (sel) VendInventoryChart._addMachine(sel.value);
            });
        }

        // Attach metric selector change
        var metricSel = document.getElementById('inv-chart-metric-select');
        if (metricSel) {
            metricSel.addEventListener('change', function() {
                VendInventoryChart._activeMetric = metricSel.value;
                // Re-process all series with new metric and re-render
                VendInventoryChart._reloadAllSeries();
            });
        }
    },

    _addMachine: function(machineId) {
        // Check if already added
        for (var i = 0; i < this._series.length; i++) {
            if (this._series[i].machineId === machineId) return;
        }
        if (this._series.length >= this.MAX_MACHINES) return;

        var machineName = machineId;
        for (var j = 0; j < this._machineList.length; j++) {
            if (this._machineList[j].id === machineId) {
                machineName = this._machineList[j].name;
                break;
            }
        }

        var color = this.COLORS[this._series.length % this.COLORS.length];
        var newSeries = { machineId: machineId, machineName: machineName, snapshots: [], dayData: {}, color: color };
        this._series.push(newSeries);
        this._updateCount();

        // Fetch this machine's snapshots
        var config = Layer8DConfig.getConfig();
        var prefix = (config && config.app && config.app.apiPrefix) || '/vend';
        var query = encodeURIComponent(JSON.stringify({
            text: 'select * from VendInventorySnapshot where machineId=' + machineId + ' sort-by timestamp limit 500'
        }));

        var self = this;
        fetch(prefix + '/10/InvSnap?body=' + query, {
            method: 'GET',
            headers: typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
        })
        .then(function(r) { return r.json(); })
        .then(function(data) {
            newSeries.snapshots = data.list || [];
            self._processSeriesDayData(newSeries);
            self._renderChart();
        })
        .catch(function() {
            self._series.pop();
            self._updateCount();
        });
    },

    _updateCount: function() {
        var el = document.getElementById('inv-chart-count');
        if (el) el.textContent = this._series.length + '/' + this.MAX_MACHINES + ' machines';
    },

    _processSeriesDayData: function(series) {
        var metric = this._activeMetric;
        series.dayData = {};
        series.snapshots.forEach(function(s) {
            var ts = s.timestamp;
            if (typeof ts === 'string') ts = parseInt(ts);
            if (!ts) return;
            var day = Math.floor(ts / 86400) * 86400;
            if (!series.dayData[day]) {
                series.dayData[day] = { total: 0, count: 0 };
            }
            var val = s[metric] || 0;
            if (typeof val === 'string') val = parseFloat(val) || 0;
            // Revenue fields are in cents — convert to dollars for display
            if (metric === 'revenue' || metric === 'dailyRevenue') val = val / 100;
            series.dayData[day].total += val;
            series.dayData[day].count++;
        });
    },

    _reloadAllSeries: function() {
        var self = this;
        this._series.forEach(function(s) {
            self._processSeriesDayData(s);
        });
        this._renderChart();
    },

    _renderChart: function() {
        var svgContainer = document.getElementById('inv-chart-svg');
        var legendContainer = document.getElementById('inv-chart-legend');
        if (!svgContainer) return;

        var w = svgContainer.offsetWidth || 800;
        var h = 350;
        var pad = { top: 20, right: 20, bottom: 40, left: 50 };
        var plotW = w - pad.left - pad.right;
        var plotH = h - pad.top - pad.bottom;

        // Find active metric config
        var metricCfg = this._metrics.find(function(m) { return m.key === this._activeMetric; }.bind(this)) || this._metrics[0];

        // Collect all days and find max value
        var allDays = {};
        var maxVal = 0;
        this._series.forEach(function(s) {
            for (var day in s.dayData) {
                allDays[day] = true;
                var avg = s.dayData[day].total / s.dayData[day].count;
                if (avg > maxVal) maxVal = avg;
            }
        });
        if (metricCfg.max > 0) maxVal = metricCfg.max;
        else maxVal = Math.ceil(maxVal * 1.1) || 1;
        var days = Object.keys(allDays).map(Number).sort(function(a, b) { return a - b; });

        if (days.length === 0) {
            svgContainer.innerHTML = '<div style="padding: 40px; text-align: center; color: var(--layer8d-text-muted);">No data available</div>';
            legendContainer.innerHTML = '';
            return;
        }

        var minDay = days[0];
        var maxDay = days[days.length - 1];
        var dayRange = maxDay - minDay || 86400;

        // Build SVG
        var svg = '<svg width="' + w + '" height="' + h + '" xmlns="http://www.w3.org/2000/svg">';

        // Grid lines (5 ticks)
        var tickStep = maxVal / 4;
        for (var ti = 0; ti <= 4; ti++) {
            var tickVal = Math.round(ti * tickStep);
            var y = pad.top + plotH - (tickVal / maxVal) * plotH;
            svg += '<line x1="' + pad.left + '" y1="' + y + '" x2="' + (w - pad.right) + '" y2="' + y + '" stroke="var(--layer8d-border)" stroke-dasharray="3,3" opacity="0.5"/>';
            var tickLabel = metricCfg.unit === '$' ? '$' + tickVal : tickVal + (metricCfg.unit || '');
            svg += '<text x="' + (pad.left - 8) + '" y="' + (y + 4) + '" text-anchor="end" fill="var(--layer8d-text-muted)" font-size="11">' + tickLabel + '</text>';
        }

        // X-axis labels (~7 labels)
        var labelStep = Math.max(1, Math.floor(days.length / 7));
        days.forEach(function(day, i) {
            if (i % labelStep !== 0 && i !== days.length - 1) return;
            var x = pad.left + ((day - minDay) / dayRange) * plotW;
            var d = new Date(day * 1000);
            var label = (d.getMonth() + 1) + '/' + d.getDate();
            svg += '<text x="' + x + '" y="' + (h - pad.bottom + 16) + '" text-anchor="middle" fill="var(--layer8d-text-muted)" font-size="11">' + label + '</text>';
        });

        // Draw lines per series
        var legendHtml = '';
        this._series.forEach(function(series) {
            var points = [];
            days.forEach(function(day) {
                var data = series.dayData[day];
                if (data && data.count > 0) {
                    var avg = data.total / data.count;
                    var x = pad.left + ((day - minDay) / dayRange) * plotW;
                    var y = pad.top + plotH - (avg / maxVal) * plotH;
                    points.push(x + ',' + y);
                }
            });

            if (points.length > 1) {
                svg += '<polyline points="' + points.join(' ') + '" fill="none" stroke="' + series.color + '" stroke-width="2.5" stroke-linejoin="round"/>';
            } else if (points.length === 1) {
                var p = points[0].split(',');
                svg += '<circle cx="' + p[0] + '" cy="' + p[1] + '" r="5" fill="' + series.color + '"/>';
            }

            legendHtml += '<span style="display: inline-flex; align-items: center; gap: 4px; font-size: 12px; color: var(--layer8d-text-dark); padding: 3px 8px; border: 1px solid var(--layer8d-border); border-radius: 4px;">';
            legendHtml += '<span style="width: 14px; height: 3px; background: ' + series.color + '; display: inline-block; border-radius: 2px;"></span>';
            legendHtml += series.machineName;
            legendHtml += '<span style="cursor: pointer; margin-left: 4px; color: var(--layer8d-text-muted);" onclick="VendInventoryChart._removeMachine(\'' + series.machineId + '\')">&times;</span>';
            legendHtml += '</span>';
        });

        svg += '</svg>';
        svgContainer.innerHTML = svg;
        legendContainer.innerHTML = legendHtml;
    },

    _removeMachine: function(machineId) {
        this._series = this._series.filter(function(s) { return s.machineId !== machineId; });
        // Reassign colors
        for (var i = 0; i < this._series.length; i++) {
            this._series[i].color = this.COLORS[i % this.COLORS.length];
        }
        this._updateCount();
        this._renderChart();
    }
};
