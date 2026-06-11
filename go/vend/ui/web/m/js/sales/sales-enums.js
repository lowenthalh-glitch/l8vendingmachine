/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    const factory = window.Layer8EnumFactory;
    const { createStatusRenderer, renderEnum } = Layer8MRenderers;

    window.MobileSalesTransactions = window.MobileSalesTransactions || {};

    const PAYMENT_METHOD = factory.simple([
        'Unspecified', 'Cash', 'Credit Card', 'NFC Contactless',
        'QR Code', 'Mobile Wallet', 'Prepaid', 'Free'
    ]);

    const TRANSACTION_STATUS = factory.create([
        ['Unspecified', null, ''],
        ['Completed', 'completed', 'active'],
        ['Failed', 'failed', 'terminated'],
        ['Refunded', 'refunded', 'pending'],
        ['Pending', 'pending', 'pending']
    ]);

    MobileSalesTransactions.enums = {
        PAYMENT_METHOD: PAYMENT_METHOD,
        TRANSACTION_STATUS: TRANSACTION_STATUS
    };

    MobileSalesTransactions.render = {
        paymentMethod: (v) => renderEnum(v, PAYMENT_METHOD.enum),
        transactionStatus: createStatusRenderer(TRANSACTION_STATUS.enum, TRANSACTION_STATUS.classes)
    };

    MobileSalesTransactions.primaryKeys = {
        VendTransaction: 'transactionId',
        VendSettlement: 'settlementId'
    };
})();
