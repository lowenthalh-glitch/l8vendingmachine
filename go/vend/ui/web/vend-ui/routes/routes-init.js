/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8DModuleFactory.create({
        namespace: 'Routes',
        defaultModule: 'routes',
        defaultService: 'routes',
        sectionSelector: 'routes',
        initializerName: 'initializeRoutes',
        requiredNamespaces: ['RoutePlanning']
    });

    var origInit = window.initializeRoutes;
    window.initializeRoutes = function() {
        if (origInit) origInit();
        injectGeneratePanel(0);
    };

    function injectGeneratePanel(attempt) {
        var container = document.querySelector('.l8-service-view[data-service="routes"]')
            || document.querySelector('#routes-routes-table-container');
        if (!container) {
            if (attempt < 20) setTimeout(function() { injectGeneratePanel(attempt + 1); }, 300);
            return;
        }
        if (container.id && container.id.indexOf('table-container') !== -1) {
            container = container.parentElement;
        }
        if (!container || container.querySelector('.vend-route-gen-panel')) return;

        var panel = document.createElement('div');
        panel.className = 'vend-route-gen-panel';
        panel.style.cssText = 'display:flex;align-items:center;gap:12px;padding:8px 0;margin-bottom:8px;flex-wrap:wrap;';

        // Default to next weekday (skip Sat/Sun)
        var next = new Date(Date.now() + 86400000);
        while (next.getDay() === 0 || next.getDay() === 6) next.setDate(next.getDate() + 1);
        var dateStr = next.toISOString().split('T')[0];

        panel.innerHTML =
            '<label style="font-size:12px;color:var(--layer8d-text-medium);">From: ' +
            '<input type="date" id="gen-route-date" value="' + dateStr + '" style="font-size:12px;padding:3px 6px;border:1px solid var(--layer8d-border);border-radius:4px;">' +
            '</label>' +
            '<label style="font-size:12px;color:var(--layer8d-text-medium);">To: ' +
            '<input type="date" id="gen-route-date-end" value="" style="font-size:12px;padding:3px 6px;border:1px solid var(--layer8d-border);border-radius:4px;" placeholder="single day">' +
            '</label>' +
            '<label style="font-size:12px;color:var(--layer8d-text-medium);">Start Time: ' +
            '<input type="time" id="gen-route-time" value="06:00" style="font-size:12px;padding:3px 6px;border:1px solid var(--layer8d-border);border-radius:4px;">' +
            '</label>' +
            '<button id="gen-route-btn" class="layer8d-btn layer8d-btn-primary layer8d-btn-small">Generate Routes</button>' +
            '<span id="gen-route-status" style="font-size:12px;color:var(--layer8d-text-muted);"></span>';

        container.insertBefore(panel, container.firstChild);

        document.getElementById('gen-route-btn').addEventListener('click', generateRoutes);
    }

    function generateRoutes() {
        var btn = document.getElementById('gen-route-btn');
        var status = document.getElementById('gen-route-status');
        var dateInput = document.getElementById('gen-route-date');
        var timeInput = document.getElementById('gen-route-time');

        btn.disabled = true;
        btn.textContent = 'Generating...';
        status.textContent = '';

        var dateEndInput = document.getElementById('gen-route-date-end');

        var dateVal = dateInput.value || new Date(Date.now() + 86400000).toISOString().split('T')[0];
        var timeVal = timeInput.value || '06:00';
        var dateTime = new Date(dateVal + 'T' + timeVal + ':00');
        var plannedDate = Math.floor(dateTime.getTime() / 1000);

        var reqBody = {
            plannedDate: plannedDate,
            startTime: plannedDate
        };
        if (dateEndInput && dateEndInput.value) {
            var endDateTime = new Date(dateEndInput.value + 'T23:59:59');
            reqBody.plannedDateEnd = Math.floor(endDateTime.getTime() / 1000);
        }

        var prefix = Layer8DConfig.getApiPrefix();
        var body = JSON.stringify(reqBody);

        fetch(prefix + '/10/OptRoute', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + sessionStorage.bearerToken,
                'Content-Type': 'application/json'
            },
            body: body
        })
        .then(function(r) { return r.text(); })
        .then(function(text) {
            var data = {};
            try { data = JSON.parse(text); } catch(e) {}
            if (data.error) {
                Layer8DNotification.error('Route generation: ' + data.error);
                status.textContent = 'Error';
                return;
            }
            var count = data.generatedCount || 0;
            var listA = data.listACount || 0;
            var listB = data.listBAdded || 0;
            var days = data.plannedDays || 1;
            if (count > 0) {
                var msg = 'Generated ' + count + ' routes (' + listA + ' urgent, ' + listB + ' opportunistic)';
                if (days > 1) msg += ' across ' + days + ' days';
                Layer8DNotification.success(msg);
                status.textContent = count + ' routes' + (days > 1 ? ' / ' + days + ' days' : '');
            } else {
                Layer8DNotification.info('No machines need restocking');
                status.textContent = 'No machines need restocking';
            }
            if (window.Routes && Routes._state && Routes._state.serviceTables) {
                var table = Routes._state.serviceTables['routes'];
                if (table && table.fetchData) table.fetchData(1);
            }
        })
        .catch(function(err) {
            Layer8DNotification.error('Route generation failed: ' + (err.message || err));
            status.textContent = 'Failed';
        })
        .finally(function() {
            btn.disabled = false;
            btn.textContent = 'Generate Routes';
        });
    }
})();
