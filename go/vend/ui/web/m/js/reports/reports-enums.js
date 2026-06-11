/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.MobileReportsData = window.MobileReportsData || {};

    var factory = window.Layer8EnumFactory;

    var REPORT_FORMAT = factory.simple([
        'Unspecified', 'PDF', 'CSV', 'Excel', 'JSON', 'HTML'
    ]);

    var REPORT_FREQUENCY = factory.simple([
        'Unspecified', 'Daily', 'Weekly', 'Monthly',
        'Quarterly', 'Yearly', 'Once'
    ]);

    MobileReportsData.enums = {
        REPORT_FORMAT: REPORT_FORMAT,
        REPORT_FREQUENCY: REPORT_FREQUENCY
    };

    MobileReportsData.render = {};

    MobileReportsData.primaryKeys = {
        VendReport: 'reportId'
    };
})();
