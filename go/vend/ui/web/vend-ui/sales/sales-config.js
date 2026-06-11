/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    Layer8ModuleConfigFactory.create({
        namespace: 'Sales',
        modules: {
            'transactions': {
                label: 'Transactions', icon: '💳',
                services: [
                    { key: 'transactions', label: 'Transactions', icon: '💳', endpoint: '/10/Txn', model: 'VendTransaction', supportedViews: ['table', 'chart', 'timeline'] },
                    { key: 'settlements', label: 'Settlements', icon: '🏦', endpoint: '/10/Settlemnt', model: 'VendSettlement' }
                ]
            }
        },
        submodules: ['SalesTransactions']
    });
})();
