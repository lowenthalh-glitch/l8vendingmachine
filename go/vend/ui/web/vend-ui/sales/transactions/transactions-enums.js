/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/
(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum, renderDate, renderBoolean } = Layer8DRenderers;

    window.SalesTransactions = window.SalesTransactions || {};

    const PAYMENT_METHOD = factory.simple([
        'Unspecified', 'Cash', 'Credit Card', 'NFC Contactless',
        'QR Code', 'Mobile Wallet', 'Prepaid', 'Free'
    ]);

    const TRANSACTION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Completed', 'completed', 'layer8d-status-active'],
        ['Failed', 'failed', 'layer8d-status-terminated'],
        ['Refunded', 'refunded', 'layer8d-status-pending'],
        ['Pending', 'pending', 'layer8d-status-pending']
    ]);

    SalesTransactions.enums = {
        PAYMENT_METHOD: PAYMENT_METHOD,
        TRANSACTION_STATUS: TRANSACTION_STATUS
    };

    SalesTransactions.render = {
        paymentMethod: (v) => renderEnum(v, PAYMENT_METHOD.enum),
        transactionStatus: createStatusRenderer(TRANSACTION_STATUS.enum, TRANSACTION_STATUS.classes)
    };

    SalesTransactions.primaryKeys = {
        VendTransaction: 'transactionId',
        VendSettlement: 'settlementId'
    };
})();
