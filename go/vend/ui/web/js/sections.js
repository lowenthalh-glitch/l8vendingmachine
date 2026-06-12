// Section Navigation and Loading Module

// Section mapping to HTML files
const sections = {
    dashboard: 'sections/dashboard.html',
    map: 'sections/map.html',
    nayax: 'sections/nayax.html',
    fleet: 'sections/fleet.html',
    inventory: 'sections/inventory.html',
    sales: 'sections/sales.html',
    maintenance: 'sections/maintenance.html',
    alarms: 'sections/alarms.html',
    routes: 'sections/routes.html',
    analytics: 'sections/analytics.html',
    warehouse: 'sections/warehouse.html',
    compliance: 'sections/compliance.html',
    reports: 'sections/reports.html',
    system: 'sections/system.html'
};

// Section initialization functions
const sectionInitializers = {
    dashboard: () => {
        if (typeof initializeDashboard === 'function') {
            initializeDashboard();
        }
    },
    map: () => {
        if (typeof initializeMap === 'function') {
            initializeMap();
        }
    },
    nayax: () => {
        if (typeof initializeNayax === 'function') {
            initializeNayax();
        }
    },
    fleet: () => {
        if (typeof initializeFleet === 'function') {
            initializeFleet();
        }
    },
    inventory: () => {
        if (typeof initializeInventory === 'function') {
            initializeInventory();
        }
    },
    sales: () => {
        if (typeof initializeSales === 'function') {
            initializeSales();
        }
    },
    maintenance: () => {
        if (typeof initializeMaintenance === 'function') {
            initializeMaintenance();
        }
    },
    alarms: () => {
        if (typeof initializeAlm === 'function') {
            initializeAlm();
        }
    },
    routes: () => {
        if (typeof initializeRoutes === 'function') {
            initializeRoutes();
        }
    },
    analytics: () => {
        if (typeof initializeAnalytics === 'function') {
            initializeAnalytics();
        }
    },
    warehouse: () => {
        if (typeof initializeWarehouse === 'function') {
            initializeWarehouse();
        }
    },
    compliance: () => {
        if (typeof initializeCompliance === 'function') {
            initializeCompliance();
        }
    },
    reports: () => {
        if (typeof initializeReports === 'function') {
            initializeReports();
        }
    },
    system: () => {
        if (typeof initializeL8Sys === 'function') {
            initializeL8Sys();
        }
    }
};

// Track active service in hash (section/service)
function updateHash(sectionName, serviceKey) {
    if (serviceKey) {
        window.location.hash = sectionName + '/' + serviceKey;
    } else {
        window.location.hash = sectionName;
    }
}

function getHashParts() {
    var hash = window.location.hash.replace('#', '');
    var parts = hash.split('/');
    return { section: parts[0] || '', service: parts[1] || '' };
}

// Listen for sub-nav clicks to update hash with service
document.addEventListener('click', function(e) {
    var navItem = e.target.closest('.l8-subnav-item');
    if (navItem && navItem.dataset.service) {
        var hashParts = getHashParts();
        if (hashParts.section) {
            updateHash(hashParts.section, navItem.dataset.service);
        }
    }
});

// Load section content dynamically
function loadSection(sectionName) {
    updateHash(sectionName, '');
    const contentArea = document.getElementById('content-area');
    const sectionFile = sections[sectionName];

    if (!sectionFile) {
        contentArea.innerHTML = '<div class="section-container"><h2 class="section-title">Error</h2><div class="section-content">Section not found.</div></div>';
        return;
    }

    contentArea.style.opacity = '0';
    contentArea.style.transform = 'translateY(20px)';

    fetch(sectionFile + '?t=' + new Date().getTime())
        .then(response => {
            if (!response.ok) {
                throw new Error('Section not found');
            }
            return response.text();
        })
        .then(html => {
            setTimeout(() => {
                contentArea.innerHTML = html;

                const placeholder = contentArea.querySelector('[id$="-section-placeholder"]');
                if (placeholder && window.Layer8SectionGenerator) {
                    const generatedHtml = Layer8SectionGenerator.generate(sectionName);
                    const temp = document.createElement('div');
                    temp.innerHTML = generatedHtml;
                    placeholder.replaceWith(...temp.children);
                }

                setTimeout(() => {
                    contentArea.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
                    contentArea.style.opacity = '1';
                    contentArea.style.transform = 'translateY(0)';
                }, 50);

                const sectionContainer = contentArea.querySelector('.section-container');
                if (sectionContainer) {
                    sectionContainer.style.animation = 'fade-in-up 0.6s ease-out';
                }

                if (sectionInitializers[sectionName]) {
                    sectionInitializers[sectionName]();
                }

                if (window.Layer8DModuleFilter) {
                    Layer8DModuleFilter.applyToSection(sectionName);
                }

                if (window.Layer8DPermissionFilter) {
                    Layer8DPermissionFilter.applyToSection(sectionName);
                }
            }, 200);
        })
        .catch(error => {
            contentArea.innerHTML = '<div class="section-container"><h2 class="section-title">Error</h2><div class="section-content">Failed to load section content.</div></div>';
            contentArea.style.opacity = '1';
            contentArea.style.transform = 'translateY(0)';
        });
}
