/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';

    // Listen for navigation to routes service and inject Generate button
    document.addEventListener('click', function(e) {
        var serviceCard = e.target.closest('[data-service="routes"]');
        if (!serviceCard) return;

        setTimeout(function() {
            var dataContainer = document.querySelector('.mobile-data-container');
            if (!dataContainer) return;
            if (dataContainer.querySelector('.vend-generate-routes-btn')) return;

            var btn = document.createElement('button');
            btn.className = 'layer8d-btn layer8d-btn-primary layer8d-btn-small vend-generate-routes-btn';
            btn.textContent = 'Generate Routes';
            btn.style.cssText = 'margin:12px 16px; width:calc(100% - 32px);';
            btn.addEventListener('click', function() {
                btn.disabled = true;
                btn.textContent = 'Generating...';

                var prefix = Layer8MConfig.getConfig().app.apiPrefix || '/vend';
                Layer8MAuth.post(prefix + '/10/OptRoute', {
                    plannedDate: Math.floor(Date.now() / 1000) + 86400
                })
                .then(function(data) {
                    var count = data.generatedCount || 0;
                    Layer8MUtils.showSuccess('Generated ' + count + ' routes');
                })
                .catch(function(err) {
                    Layer8MUtils.showError('Route generation failed');
                })
                .finally(function() {
                    btn.disabled = false;
                    btn.textContent = 'Generate Routes';
                });
            });

            dataContainer.insertBefore(btn, dataContainer.firstChild);
        }, 500);
    });
})();
