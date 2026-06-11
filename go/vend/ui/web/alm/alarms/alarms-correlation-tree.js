/*
Layer 8 Alarms - Correlation Tree View
Injects a "Correlation" tab into the Alarm detail popup showing
linked root-cause and symptom alarms using Layer8DTreeGrid.
Loads AFTER alm-init.js so Alm._showDetailsModal exists.
*/

(function() {
    'use strict';

    var esc = Layer8DUtils.escapeHtml;
    var ALARM_ENDPOINT = '/10/Alarm';
    var MAX_DEPTH = 10;

    // ========================================================================
    // 1. Wrap Alm._showDetailsModal to inject Correlation tab for Alarm model
    // ========================================================================

    var origShowDetails = Alm._showDetailsModal;

    Alm._showDetailsModal = function(service, item, itemId) {
        origShowDetails.call(Alm, service, item, itemId);

        if (service.model !== 'Alarm') return;

        setTimeout(function() {
            var body = Layer8DPopup.getBody();
            if (!body) return;
            injectCorrelationTab(body, item);
        }, 100);
    };

    // ========================================================================
    // 2. Tab injection
    // ========================================================================

    function injectCorrelationTab(body, alarmData) {
        var tabsBar = body.querySelector('.probler-popup-tabs');
        var tabContent = body.querySelector('.probler-popup-tab-content');
        if (!tabsBar || !tabContent) return;

        var tabIndex = 'tab-correlation';

        var tabBtn = document.createElement('div');
        tabBtn.className = 'probler-popup-tab';
        tabBtn.dataset.tab = tabIndex;
        tabBtn.textContent = 'Correlation';
        tabsBar.appendChild(tabBtn);

        var pane = document.createElement('div');
        pane.className = 'probler-popup-tab-pane';
        pane.dataset.pane = tabIndex;
        pane.innerHTML = '<div class="alm-corr-loading">Loading correlation data...</div>';
        tabContent.appendChild(pane);

        fetchAndRender(alarmData, pane);
    }

    // ========================================================================
    // 3. Data fetching
    // ========================================================================

    function buildQuery(model, where) {
        return 'select * from ' + model + ' where ' + where;
    }

    async function queryAlarms(where) {
        var url = Layer8DConfig.resolveEndpoint(ALARM_ENDPOINT)
            + '?body=' + encodeURIComponent(JSON.stringify({ text: buildQuery('Alarm', where) }));

        var resp = await fetch(url, {
            method: 'GET',
            headers: Object.assign(
                { 'Content-Type': 'application/json' },
                typeof getAuthHeaders === 'function' ? getAuthHeaders() : {}
            )
        });

        if (!resp.ok) return [];
        var data = await resp.json();
        return data.list || [];
    }

    // Recursively collect all correlated alarms into a flat array
    async function fetchAllCorrelated(alarm, visited, depth) {
        if (depth > MAX_DEPTH) return [];
        if (visited[alarm.alarmId]) return [];
        visited[alarm.alarmId] = true;

        var result = [alarm];

        var symptoms = await queryAlarms('RootCauseAlarmId=' + alarm.alarmId);
        for (var i = 0; i < symptoms.length; i++) {
            var childResults = await fetchAllCorrelated(symptoms[i], visited, depth + 1);
            result = result.concat(childResults);
        }

        return result;
    }

    // ========================================================================
    // 4. Render using Layer8DTreeGrid
    // ========================================================================

    async function fetchAndRender(alarm, container) {
        try {
            if (!alarm.isRootCause && !alarm.rootCauseAlarmId) {
                container.innerHTML = renderEmpty();
                return;
            }

            var flatAlarms = await fetchAllCorrelated(alarm, {}, 0);
            if (flatAlarms.length === 0) {
                container.innerHTML = renderEmpty();
                return;
            }

            // Build container HTML
            var html = '';

            // Parent link if this alarm has a root cause
            if (alarm.rootCauseAlarmId) {
                var parents = await queryAlarms('AlarmId=' + alarm.rootCauseAlarmId);
                if (parents.length > 0) {
                    html += renderParentLink(parents[0]);
                }
            }

            var gridId = 'alm-corr-grid-' + alarm.alarmId;
            html += '<div id="' + gridId + '"></div>';

            container.innerHTML = html;

            // Attach parent link click handler
            var parentLink = container.querySelector('.alm-corr-parent-link');
            if (parentLink) {
                parentLink.addEventListener('click', function() {
                    openAlarmDetail(parentLink.dataset.alarmId);
                });
            }

            // Create tree grid
            var currentAlarmId = alarm.alarmId;
            var treeGrid = new Layer8DTreeGrid({
                containerId: gridId,
                primaryKey: 'alarmId',
                viewConfig: {
                    parentIdField: 'rootCauseAlarmId',
                    idField: 'alarmId',
                    expandedByDefault: true
                },
                columns: [
                    {
                        key: 'name',
                        label: 'Alarm',
                        render: function(item) {
                            return renderAlarmCell(item, currentAlarmId);
                        }
                    },
                    {
                        key: 'severity',
                        label: 'Severity',
                        render: function(item) {
                            if (AlmAlarms.render && AlmAlarms.render.severity) {
                                return AlmAlarms.render.severity(item.severity);
                            }
                            return esc(String(item.severity || ''));
                        }
                    },
                    {
                        key: 'state',
                        label: 'State',
                        render: function(item) {
                            if (AlmAlarms.render && AlmAlarms.render.state) {
                                return AlmAlarms.render.state(item.state);
                            }
                            return esc(String(item.state || ''));
                        }
                    },
                    { key: 'nodeName', label: 'Node' }
                ],
                onItemClick: function(item) {
                    if (item.alarmId !== currentAlarmId) {
                        openAlarmDetail(item.alarmId);
                    }
                }
            });

            treeGrid.init();
            treeGrid.setData(flatAlarms);

            // Highlight current alarm row
            highlightCurrentRow(container, currentAlarmId);
        } catch (err) {
            console.error('Correlation tree error:', err);
            container.innerHTML = '<div class="alm-corr-error">Failed to load correlation data</div>';
        }
    }

    // ========================================================================
    // 5. Cell renderers and helpers
    // ========================================================================

    function renderAlarmCell(item, currentAlarmId) {
        var isRoot = item.isRootCause && !item.rootCauseAlarmId;
        var badgeClass = isRoot ? 'alm-corr-badge-root' : 'alm-corr-badge-symptom';
        var badgeLabel = isRoot ? 'ROOT' : 'SYMPTOM';
        var isCurrent = item.alarmId === currentAlarmId;
        var nameClass = isCurrent ? ' alm-corr-name-current' : '';

        return '<span class="alm-corr-badge alm-corr-badge-inline ' + badgeClass + '">'
            + badgeLabel + '</span> '
            + '<span class="alm-corr-name' + nameClass + '">'
            + esc(item.name || item.alarmId) + '</span>';
    }

    function highlightCurrentRow(container, alarmId) {
        var row = container.querySelector('.layer8d-tree-grid-row[data-id="' + alarmId + '"]');
        if (row) {
            row.classList.add('alm-corr-current-row');
        }
    }

    function renderEmpty() {
        return '<div class="alm-corr-empty">'
            + '<div class="alm-corr-empty-icon">&#x1F517;</div>'
            + '<div>No correlation data</div>'
            + '<div style="font-size:12px;margin-top:4px;">This alarm is not part of a correlation group</div>'
            + '</div>';
    }

    function renderParentLink(parentAlarm) {
        var sevHtml = '';
        if (AlmAlarms.render && AlmAlarms.render.severity) {
            sevHtml = '<span class="alm-corr-severity">' + AlmAlarms.render.severity(parentAlarm.severity) + '</span>';
        }
        var stateHtml = '';
        if (AlmAlarms.render && AlmAlarms.render.state) {
            stateHtml = '<span class="alm-corr-severity">' + AlmAlarms.render.state(parentAlarm.state) + '</span>';
        }

        return '<div class="alm-corr-parent-link" data-alarm-id="' + esc(parentAlarm.alarmId) + '">'
            + '<span class="alm-corr-parent-arrow">&#x2191;</span>'
            + '<span class="alm-corr-parent-label">Root Cause:</span>'
            + sevHtml
            + '<span class="alm-corr-name">' + esc(parentAlarm.name || parentAlarm.alarmId) + '</span>'
            + stateHtml
            + '<span class="alm-corr-meta">' + esc(parentAlarm.nodeName || parentAlarm.nodeId || '') + '</span>'
            + '</div>';
    }

    // ========================================================================
    // 6. Navigation — open stacked detail popups
    // ========================================================================

    async function openAlarmDetail(alarmId) {
        try {
            var alarms = await queryAlarms('AlarmId=' + alarmId);
            if (alarms.length === 0) {
                Layer8DNotification.warning('Alarm not found');
                return;
            }
            var alarm = alarms[0];
            var service = { key: 'alarms', label: 'Active Alarms', endpoint: '/10/Alarm', model: 'Alarm' };
            Alm._showDetailsModal(service, alarm, alarm.alarmId);
        } catch (err) {
            console.error('Failed to open alarm detail:', err);
            Layer8DNotification.error('Failed to load alarm');
        }
    }

})();
