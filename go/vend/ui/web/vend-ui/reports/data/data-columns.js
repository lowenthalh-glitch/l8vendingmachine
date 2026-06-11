/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;

    ReportsData.columns = {
        VendReport: [
            ...col.id('reportId'),
            ...col.col('code', 'Code'),
            ...col.col('name', 'Name'),
            ...col.col('reportType', 'Type'),
            ...col.col('category', 'Category'),
            ...col.boolean('isPublic', 'Public'),
            ...col.number('executionCount', 'Executions'),
            ...col.date('lastExecuted', 'Last Executed')
        ]
    };

    ReportsData.primaryKeys = {
        VendReport: 'reportId'
    };
})();
