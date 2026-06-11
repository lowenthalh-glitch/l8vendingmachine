/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    var f = window.Layer8FormFactory;
    var enums = SalesTransactions.enums;

    SalesTransactions.forms = {
        VendTransaction: f.form('Transaction', [
            f.section('Transaction Details', [
                ...f.text('transactionId', 'Transaction ID', false, { readOnly: true }),
                ...f.text('machineId', 'Machine', false, { readOnly: true }),
                ...f.text('slotId', 'Slot', false, { readOnly: true }),
                ...f.text('productName', 'Product', false, { readOnly: true }),
                ...f.money('price', 'Price', { readOnly: true }),
                ...f.select('paymentMethod', 'Payment Method', enums.PAYMENT_METHOD.enum, false, { readOnly: true }),
                ...f.select('status', 'Status', enums.TRANSACTION_STATUS.enum, false, { readOnly: true }),
                ...f.text('cardType', 'Card Type', false, { readOnly: true }),
                ...f.text('cardLastFour', 'Card Last Four', false, { readOnly: true }),
                ...f.text('authorizationCode', 'Auth Code', false, { readOnly: true })
            ])
        ]),
        VendSettlement: f.form('Settlement', [
            f.section('Settlement Info', [
                ...f.reference('machineId', 'Machine', 'VendMachine'),
                ...f.date('settlementDate', 'Settlement Date'),
                ...f.number('transactionCount', 'Transaction Count'),
                ...f.money('totalAmount', 'Total Amount'),
                ...f.text('processorReference', 'Processor Reference'),
                ...f.text('status', 'Status')
            ])
        ])
    };
})();
