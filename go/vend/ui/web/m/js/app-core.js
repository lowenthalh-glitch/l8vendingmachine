/* © 2025 Sharon Aicler (saichler@gmail.com) Layer 8 Ecosystem - Apache 2.0 */
(function() {
    'use strict';

    const SECTIONS = {
        'dashboard': 'sections/dashboard.html',
        'system': 'sections/system.html'
    };

    let currentSection = 'dashboard';
    let sectionCache = {};

    window.showErrorAndLogout = function(message, detail) {
        if (typeof Layer8MAuth !== 'undefined') {
            Layer8MAuth.showErrorAndLogout(message, detail);
        } else {
            alert(message + (detail ? '\n\n' + detail : ''));
            window.location.href = '/l8ui/login/';
        }
    };

    window.MobileApp = {
        async init() {
            if (!Layer8MAuth.requireAuth()) return;
            await Layer8MConfig.load();
            await Layer8DConfig.load();
            this.updateUserInfo();

            // L8VendingMachine does NOT use ModConfig -- skip Layer8DModuleFilter

            if (typeof Layer8MNav !== 'undefined') {
                Layer8MNav.showHome();
            }
        },

        updateUserInfo() {
            var username = Layer8MAuth.getUsername();
            var el = document.getElementById('user-display-name');
            if (el) el.textContent = username || 'User';
        },

        logout() {
            Layer8MAuth.logout();
        }
    };

    document.addEventListener('DOMContentLoaded', function() {
        MobileApp.init();
    });
})();
