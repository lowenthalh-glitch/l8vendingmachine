/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    Layer8SectionConfigs.register('compliance', {
        title: 'Compliance & Inspections',
        subtitle: 'Inspections, Findings, Certifications',
        icon: '✅',
        initFn: 'initializeCompliance',
        modules: [{
            key: 'inspections', label: 'Inspections', icon: '✅', isDefault: true,
            services: [
                { key: 'inspections', label: 'Inspections', icon: '🔍', isDefault: true },
                { key: 'findings', label: 'Findings', icon: '📋' },
                { key: 'certifications', label: 'Certifications', icon: '📜' }
            ]
        }]
    });
})();
