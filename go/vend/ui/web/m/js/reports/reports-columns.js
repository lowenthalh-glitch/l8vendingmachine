/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var col = window.Layer8ColumnFactory;

    MobileReportsData.columns = {
        VendReport: [
            ...col.id('reportId'),
            ...col.col('code', 'Code'),
            { key: 'name', label: 'Name', primary: true, sortKey: 'name', filterKey: 'name' },
            ...col.col('reportType', 'Type'),
            { key: 'category', label: 'Category', secondary: true, sortKey: 'category', filterKey: 'category' },
            ...col.boolean('isPublic', 'Public'),
            ...col.number('executionCount', 'Executions'),
            ...col.date('lastExecuted', 'Last Executed')
        ]
    };
})();
