/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.ReportsData = window.ReportsData || {};

    var factory = window.Layer8EnumFactory;

    var REPORT_FORMAT = factory.simple([
        'Unspecified', 'PDF', 'CSV', 'Excel', 'JSON', 'HTML'
    ]);

    var REPORT_FREQUENCY = factory.simple([
        'Unspecified', 'Daily', 'Weekly', 'Monthly',
        'Quarterly', 'Yearly', 'Once'
    ]);

    ReportsData.enums = {
        REPORT_FORMAT: REPORT_FORMAT,
        REPORT_FREQUENCY: REPORT_FREQUENCY
    };

    ReportsData.render = {};
})();
