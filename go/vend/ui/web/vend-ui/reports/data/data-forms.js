/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    var f = window.Layer8FormFactory;

    ReportsData.forms = {
        VendReport: f.form('Report', [
            f.section('Report Information', [
                ...f.text('code', 'Code'),
                ...f.text('name', 'Name', true),
                ...f.textarea('description', 'Description'),
                ...f.text('reportType', 'Report Type'),
                ...f.text('category', 'Category'),
                ...f.checkbox('isPublic', 'Public'),
                ...f.textarea('query', 'Query')
            ])
        ])
    };
})();
