/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;
    var render = SalesTransactions.render;

    SalesTransactions.columns = {
        VendTransaction: [
            ...col.id('transactionId'),
            ...col.col('machineId', 'Machine'),
            ...col.date('timestamp', 'Timestamp'),
            ...col.col('slotId', 'Slot'),
            ...col.col('productName', 'Product'),
            ...col.money('price', 'Price'),
            ...col.enum('paymentMethod', 'Payment', null, render.paymentMethod),
            ...col.status('status', 'Status', null, render.transactionStatus),
            ...col.col('cardType', 'Card Type')
        ],
        VendSettlement: [
            ...col.id('settlementId'),
            ...col.col('machineId', 'Machine'),
            ...col.date('settlementDate', 'Settlement Date'),
            ...col.number('transactionCount', 'Transactions'),
            ...col.money('totalAmount', 'Total Amount'),
            ...col.col('status', 'Status')
        ]
    };
})();
