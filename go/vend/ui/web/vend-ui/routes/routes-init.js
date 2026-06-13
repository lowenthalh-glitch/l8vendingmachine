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

        // Inject "Generate Routes" button — retry until DOM is ready
        function tryInjectButton(attempt) {
            var container = document.querySelector('.l8-service-view[data-service="routes"]');
            if (!container) {
                if (attempt < 10) setTimeout(function() { tryInjectButton(attempt + 1); }, 300);
                return;
            }
            // Avoid duplicate buttons
            if (container.querySelector('.vend-generate-routes-btn')) return;

            var btn = document.createElement('button');
            btn.className = 'layer8d-btn layer8d-btn-primary layer8d-btn-small vend-generate-routes-btn';
            btn.textContent = 'Generate Routes';
            btn.style.marginBottom = '12px';
            btn.addEventListener('click', function() {
                btn.disabled = true;
                btn.textContent = 'Generating...';

                var prefix = Layer8DConfig.getApiPrefix();
                var body = JSON.stringify({ plannedDate: Math.floor(Date.now() / 1000) + 86400 });
                fetch(prefix + '/10/OptRoute', {
                    method: 'POST',
                    headers: { 'Authorization': 'Bearer ' + sessionStorage.bearerToken, 'Content-Type': 'application/json' },
                    body: body
                })
                .then(function(r) { return r.text(); })
                .then(function(text) {
                    var data = {};
                    try { data = JSON.parse(text); } catch(e) { /* response may not be JSON */ }
                    if (data.error) {
                        Layer8DNotification.error('Route generation: ' + data.error);
                        return;
                    }
                    var count = data.generatedCount || 0;
                    var listA = data.listACount || 0;
                    var listB = data.listBAdded || 0;
                    if (count > 0) {
                        Layer8DNotification.success('Generated ' + count + ' routes (' + listA + ' urgent, ' + listB + ' opportunistic)');
                    } else {
                        Layer8DNotification.info('No machines need restocking');
                    }
                    // Refresh routes table
                    if (window.Routes && Routes._state && Routes._state.serviceTables) {
                        var table = Routes._state.serviceTables['routes'];
                        if (table && table.fetchData) table.fetchData(1);
                    }
                })
                .catch(function(err) {
                    Layer8DNotification.error('Route generation failed: ' + (err.message || err));
                })
                .finally(function() {
                    btn.disabled = false;
                    btn.textContent = 'Generate Routes';
                });
            });

            container.insertBefore(btn, container.firstChild);
        }
        tryInjectButton(0);
    };
})();
