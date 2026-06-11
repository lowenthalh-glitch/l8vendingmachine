/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('sales', {
        title: 'Sales & Transactions',
        subtitle: 'Transactions, Settlements',
        icon: '💰',
        initFn: 'initializeSales',
        modules: [{
            key: 'transactions', label: 'Transactions', icon: '💰', isDefault: true,
            services: [
                { key: 'transactions', label: 'Transactions', icon: '💰', isDefault: true },
                { key: 'settlements', label: 'Settlements', icon: '📄' }
            ]
        }]
    });
})();
