/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';
    window.LAYER8M_NAV_CONFIG = window.LAYER8M_NAV_CONFIG || {};
    var base = window.LAYER8M_NAV_CONFIG_BASE || {};
    var vend = window.LAYER8M_NAV_CONFIG_VEND || {};

    LAYER8M_NAV_CONFIG.modules = base.modules || [];

    // Merge vend module configs
    var keys = Object.keys(vend);
    for (var i = 0; i < keys.length; i++) {
        LAYER8M_NAV_CONFIG[keys[i]] = vend[keys[i]];
    }

    LAYER8M_NAV_CONFIG.icons = {};
    LAYER8M_NAV_CONFIG.getIcon = function(key) {
        return LAYER8M_NAV_CONFIG.icons[key] || '';
    };
})();
