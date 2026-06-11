/*
© 2025 Sharon Aicler (saichler@gmail.com)
Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
*/

// Get authentication headers with bearer token
function getAuthHeaders() {
    const bearerToken = sessionStorage.getItem('bearerToken');
    return {
        'Authorization': bearerToken ? 'Bearer ' + bearerToken : '',
        'Content-Type': 'application/json'
    };
}

// Logout function
function logout() {
    sessionStorage.removeItem('bearerToken');
    localStorage.removeItem('bearerToken');
    localStorage.removeItem('rememberedUser');
    window.location.href = 'l8ui/login/index.html';
}

// Show error popup before logging out
function showErrorAndLogout(message, detail) {
    sessionStorage.removeItem('bearerToken');
    localStorage.removeItem('bearerToken');
    localStorage.removeItem('rememberedUser');
    if (typeof Layer8DPopup !== 'undefined') {
        Layer8DPopup.show({
            title: 'Session Error',
            content: '<p>' + message + '</p>' + (detail ? '<p style="color:#718096;font-size:13px;">' + detail + '</p>' : ''),
            size: 'small',
            showFooter: true,
            saveButtonText: 'Log In',
            onSave: function() { window.location.href = 'l8ui/login/index.html'; }
        });
    } else {
        alert(message + (detail ? '\n\n' + detail : ''));
        window.location.href = 'l8ui/login/index.html';
    }
}

// App initialization (matches l8househm pattern)
(function() {
    'use strict';

    document.addEventListener('DOMContentLoaded', async function() {
        // Expose bearer token for iframes (targets UI reads from localStorage and window.parent)
        var bearerToken = sessionStorage.getItem('bearerToken');
        if (bearerToken) {
            localStorage.setItem('bearerToken', bearerToken);
            window.bearerToken = bearerToken;
        }

        // Load app config
        if (window.Layer8DConfig && Layer8DConfig.load) {
            try {
                await Layer8DConfig.load();
            } catch (e) {
                console.error('Failed to load app config', e);
            }
        }

        // Wire sidebar links to load sections
        const sidebarItems = document.querySelectorAll('.sidebar-item[data-section]');
        sidebarItems.forEach(function(item) {
            item.addEventListener('click', function(e) {
                e.preventDefault();
                const section = item.getAttribute('data-section');
                // Update active state
                sidebarItems.forEach(function(s) { s.classList.remove('active'); });
                item.classList.add('active');
                // Load the section
                if (typeof loadSection === 'function') {
                    loadSection(section);
                }
            });
        });

        // Listen for postMessage events from iframes (targets UI popup bridge)
        window.addEventListener('message', function(event) {
            if (!event.data || !event.data.type) return;

            switch (event.data.type) {
                case 'probler-popup-show':
                    if (typeof Layer8DPopup !== 'undefined') {
                        Layer8DPopup.show(event.data.config);
                    }
                    break;
                case 'probler-popup-close':
                    if (typeof Layer8DPopup !== 'undefined') {
                        Layer8DPopup.close();
                    }
                    break;
                case 'probler-popup-update':
                    if (typeof Layer8DPopup !== 'undefined') {
                        Layer8DPopup.updateContent(event.data.content);
                    }
                    break;
            }
        });

        // Load section from hash or default to dashboard
        if (typeof loadSection === 'function') {
            var hashParts = typeof getHashParts === 'function' ? getHashParts() : { section: '', service: '' };
            var initSection = hashParts.section && sections[hashParts.section] ? hashParts.section : 'dashboard';
            loadSection(initSection);

            // Activate saved service tab after section loads
            if (hashParts.service) {
                setTimeout(function() {
                    var navItem = document.querySelector('.l8-subnav-item[data-service="' + hashParts.service + '"]');
                    if (navItem) navItem.click();
                }, 500);
            }
        }
    });
})();
