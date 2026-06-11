/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = MobileSalesTransactions.render;

    MobileSalesTransactions.columns = {
        VendTransaction: [
            ...col.id('transactionId'),
            ...col.col('machineId', 'Machine'),
            ...col.date('timestamp', 'Timestamp'),
            ...col.col('slotId', 'Slot'),
            { key: 'productName', label: 'Product', primary: true, sortKey: 'productName', filterKey: 'productName' },
            { key: 'price', label: 'Price', secondary: true, sortKey: 'price', filterKey: 'price',
              render: (item) => Layer8MRenderers.renderMoney(item.price) },
            ...col.enum('paymentMethod', 'Payment', null, render.paymentMethod),
            ...col.status('status', 'Status', null, render.transactionStatus),
            ...col.col('cardType', 'Card Type')
        ],
        VendSettlement: [
            { key: 'settlementId', label: 'Settlement ID', primary: true, sortKey: 'settlementId', filterKey: 'settlementId' },
            ...col.col('machineId', 'Machine'),
            ...col.date('settlementDate', 'Settlement Date'),
            ...col.number('transactionCount', 'Transactions'),
            ...col.money('totalAmount', 'Total Amount'),
            ...col.col('status', 'Status')
        ]
    };
})();
