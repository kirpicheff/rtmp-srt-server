// Система локализации
const translations = {
    ru: {
        // Заголовок
        'header.subtitle': 'Управление трансляциями и настройками сервера',

        // Навигация
        'nav.status': '📊 Статус',
        'nav.inputs': '📥 Входы',
        'nav.outputs': '📤 Выходы',
        'nav.settings': '⚙️ Настройки',

        // Общие
        'common.loading': 'Загрузка...',
        'common.refresh': 'Обновить',
        'common.add': 'Добавить',
        'common.remove': 'Удалить',
        'common.reconnect': 'Реконнект',
        'common.save': 'Сохранить',
        'common.cancel': 'Отмена',
        'common.active': 'Активен',
        'common.inactive': 'Неактивен',
        'common.bitrate': 'Битрейт',
        'common.uptime': 'Аптайм',
        'common.errors': 'Ошибки',
        'common.status': 'Статус',
        'common.name': 'Имя',
        'common.url': 'URL',
        'common.path': 'Путь',
        'common.outputs': 'Выходы',
        'common.connections': 'Подключения',
        'common.lastUpdate': 'Последнее обновление',
        'common.updated': 'обновлено',

        // Статус
        'status.title': 'Общий статус сервера',
        'status.noInputs': 'Нет активных входов',
        'status.noOutputs': 'Нет выходов',
        'status.loadError': 'Ошибка загрузки статуса',

        // Входы
        'inputs.title': 'Управление входами',
        'inputs.addNew': 'Добавить новый вход',
        'inputs.inputName': 'Имя входа',
        'inputs.rtmpPath': 'Путь RTMP',
        'inputs.outputsPlaceholder': 'srt://example.com:9000\nrtmp://example.com/live/stream',
        'inputs.outputsLabel': 'Выходы (по одному на строку):',
        'inputs.existing': 'Существующие входы',
        'inputs.noInputs': 'Нет входов',
        'inputs.addSuccess': 'Вход успешно добавлен',
        'inputs.removeSuccess': 'Вход успешно удален',
        'inputs.removeConfirm': 'Удалить вход',
        'inputs.addError': 'Ошибка добавления входа',
        'inputs.removeError': 'Ошибка удаления входа',
        'inputs.loadError': 'Ошибка загрузки входов',

        // Выходы
        'outputs.title': 'Управление выходами',
        'outputs.addOutput': 'Добавить выход',
        'outputs.input': 'Вход',
        'outputs.outputUrl': 'URL выхода',
        'outputs.management': 'Управление выходами',
        'outputs.noInputsForManagement': 'Нет входов для управления выходами',
        'outputs.noActiveInputs': 'Нет активных входов',
        'outputs.addSuccess': 'Выход успешно добавлен',
        'outputs.removeSuccess': 'Выход успешно удален',
        'outputs.reconnectInitiated': 'Реконнект инициирован',
        'outputs.removeConfirm': 'Удалить выход',
        'outputs.addError': 'Ошибка добавления выхода',
        'outputs.removeError': 'Ошибка удаления выхода',
        'outputs.reconnectError': 'Ошибка реконнекта',
        'outputs.loadError': 'Ошибка загрузки выходов',

        // Настройки
        'settings.title': 'Глобальные настройки',
        'settings.srt': 'Настройки SRT',
        'settings.latency': 'Latency (мс):',
        'settings.timeout': 'Connect Timeout (мс):',
        'settings.passphrase': 'Passphrase:',
        'settings.streamId': 'Stream ID:',
        'settings.logging': 'Настройки логирования',
        'settings.logToFile': 'Логировать в файл',
        'settings.logFile': 'Файл лога:',
        'settings.reconnectInterval': 'Интервал переподключения',
        'settings.reconnectIntervalLabel': 'Интервал (секунды):',
        'settings.saveSuccess': 'Настройки успешно сохранены',
        'settings.reloadSuccess': 'Настройки перезагружены',
        'settings.saveError': 'Ошибка сохранения настроек',
        'settings.reloadError': 'Ошибка перезагрузки настроек',
        'settings.loadError': 'Ошибка загрузки настроек',
        'settings.reload': 'Перезагрузить из файла'
    },
    en: {
        // Header
        'header.subtitle': 'Stream and server settings management',

        // Navigation
        'nav.status': '📊 Status',
        'nav.inputs': '📥 Inputs',
        'nav.outputs': '📤 Outputs',
        'nav.settings': '⚙️ Settings',

        // Common
        'common.loading': 'Loading...',
        'common.refresh': 'Refresh',
        'common.add': 'Add',
        'common.remove': 'Remove',
        'common.reconnect': 'Reconnect',
        'common.save': 'Save',
        'common.cancel': 'Cancel',
        'common.active': 'Active',
        'common.inactive': 'Inactive',
        'common.bitrate': 'Bitrate',
        'common.uptime': 'Uptime',
        'common.errors': 'Errors',
        'common.status': 'Status',
        'common.name': 'Name',
        'common.url': 'URL',
        'common.path': 'Path',
        'common.outputs': 'Outputs',
        'common.connections': 'Connections',
        'common.lastUpdate': 'Last update',
        'common.updated': 'updated',

        // Status
        'status.title': 'Server Status Overview',
        'status.noInputs': 'No active inputs',
        'status.noOutputs': 'No outputs',
        'status.loadError': 'Error loading status',

        // Inputs
        'inputs.title': 'Input Management',
        'inputs.addNew': 'Add New Input',
        'inputs.inputName': 'Input Name',
        'inputs.rtmpPath': 'RTMP Path',
        'inputs.outputsPlaceholder': 'srt://example.com:9000\nrtmp://example.com/live/stream',
        'inputs.outputsLabel': 'Outputs (one per line):',
        'inputs.existing': 'Existing Inputs',
        'inputs.noInputs': 'No inputs',
        'inputs.addSuccess': 'Input successfully added',
        'inputs.removeSuccess': 'Input successfully removed',
        'inputs.removeConfirm': 'Remove input',
        'inputs.addError': 'Error adding input',
        'inputs.removeError': 'Error removing input',
        'inputs.loadError': 'Error loading inputs',

        // Outputs
        'outputs.title': 'Output Management',
        'outputs.addOutput': 'Add Output',
        'outputs.input': 'Input',
        'outputs.outputUrl': 'Output URL',
        'outputs.management': 'Output Management',
        'outputs.noInputsForManagement': 'No inputs available for output management',
        'outputs.noActiveInputs': 'No active inputs',
        'outputs.addSuccess': 'Output successfully added',
        'outputs.removeSuccess': 'Output successfully removed',
        'outputs.reconnectInitiated': 'Reconnect initiated',
        'outputs.removeConfirm': 'Remove output',
        'outputs.addError': 'Error adding output',
        'outputs.removeError': 'Error removing output',
        'outputs.reconnectError': 'Reconnect error',
        'outputs.loadError': 'Error loading outputs',

        // Settings
        'settings.title': 'Global Settings',
        'settings.srt': 'SRT Settings',
        'settings.latency': 'Latency (ms):',
        'settings.timeout': 'Connect Timeout (ms):',
        'settings.passphrase': 'Passphrase:',
        'settings.streamId': 'Stream ID:',
        'settings.logging': 'Logging Settings',
        'settings.logToFile': 'Log to file',
        'settings.logFile': 'Log file:',
        'settings.reconnectInterval': 'Reconnect Interval',
        'settings.reconnectIntervalLabel': 'Interval (seconds):',
        'settings.saveSuccess': 'Settings successfully saved',
        'settings.reloadSuccess': 'Settings reloaded',
        'settings.saveError': 'Error saving settings',
        'settings.reloadError': 'Error reloading settings',
        'settings.loadError': 'Error loading settings',
        'settings.reload': 'Reload from file'
    }
};

let currentLanguage = localStorage.getItem('language') || 'ru';

// Функция перевода
function t(key) {
    return translations[currentLanguage][key] || key;
}

// Функция применения переводов
function applyTranslations() {
    // Переводим элементы с data-i18n
    document.querySelectorAll('[data-i18n]').forEach(element => {
        const key = element.getAttribute('data-i18n');
        element.textContent = t(key);
    });

    // Переводим атрибуты title с data-i18n-title
    document.querySelectorAll('[data-i18n-title]').forEach(element => {
        const key = element.getAttribute('data-i18n-title');
        element.title = t(key);
    });

    // Переводим placeholder с data-i18n-placeholder
    document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
        const key = element.getAttribute('data-i18n-placeholder');
        element.placeholder = t(key);
    });
}

// Функция смены языка
function changeLanguage(lang) {
    currentLanguage = lang;
    localStorage.setItem('language', lang);
    applyTranslations();

    // Перезагружаем текущую вкладку чтобы обновить динамический контент
    loadCurrentTab();
}

// Конфигурация API
const API_BASE = '/api';
const API_CREDENTIALS = btoa('admin:secret'); // В продакшене лучше получать через форму входа

// Глобальные переменные
let currentTab = 'status';
let refreshInterval;

// Инициализация приложения
document.addEventListener('DOMContentLoaded', function () {
    setupNavigation();
    setupLanguageSwitcher();
    applyTranslations();
    loadCurrentTab();
    startAutoRefresh();
});

// Настройка переключателя языка
function setupLanguageSwitcher() {
    const languageSelect = document.getElementById('language-select');
    languageSelect.value = currentLanguage;

    languageSelect.addEventListener('change', function () {
        changeLanguage(this.value);
    });
}

// Настройка навигации
function setupNavigation() {
    const tabs = document.querySelectorAll('.nav-tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function () {
            const tabName = this.getAttribute('data-tab');
            switchTab(tabName);
        });
    });
}

// Переключение вкладок
function switchTab(tabName) {
    // Обновляем активную вкладку
    document.querySelectorAll('.nav-tab').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

    // Скрываем все контенты
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });

    // Показываем нужный контент
    document.getElementById(tabName).classList.add('active');
    currentTab = tabName;

    // Загружаем данные для вкладки
    loadCurrentTab();
}

// Загрузка данных для текущей вкладки
function loadCurrentTab() {
    switch (currentTab) {
        case 'status':
            loadStatus();
            break;
        case 'inputs':
            loadInputs();
            break;
        case 'outputs':
            loadOutputsInitial(); // Используем специальную функцию для первоначальной загрузки
            break;
        case 'settings':
            loadSettings();
            break;
    }
}

// Автообновление данных
function startAutoRefresh() {
    refreshInterval = setInterval(async () => {
        if (currentTab === 'status') {
            loadStatus();
        } else if (currentTab === 'outputs') {
            // Обновляем только статусы выходов, не трогаем форму
            loadOutputs(); // Используем упрощенную функцию
        }
    }, 5000); // Обновляем статус каждые 5 секунд
}

// Обновление текущей вкладки
function refreshCurrentTab() {
    if (currentTab === 'outputs') {
        // Для выходов обновляем только статусы
        loadOutputsInitial();
    } else {
        loadCurrentTab();
    }
}

// API функции
async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`;
    const defaultOptions = {
        headers: {
            'Authorization': `Basic ${API_CREDENTIALS}`,
            'Content-Type': 'application/json'
        }
    };

    try {
        const response = await fetch(url, { ...defaultOptions, ...options });

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        if (response.headers.get('content-type')?.includes('application/json')) {
            return await response.json();
        }

        return { success: true };
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// Загрузка статуса
async function loadStatus() {
    try {
        const statuses = await apiRequest('/status/all');
        displayStatus(statuses);
    } catch (error) {
        showError(t('status.loadError') + ': ' + error.message);
    }
}

// Отображение статуса
function displayStatus(statuses) {
    const container = document.getElementById('status-content');
    const currentTime = new Date().toLocaleTimeString();

    if (!statuses || statuses.length === 0) {
        container.innerHTML = `<div class="alert alert-error">${t('status.noInputs')}</div>`;
        return;
    }

    const html = `
        <div style="margin-bottom: 20px; color: #6c757d; font-size: 0.9em;">
            ${t('common.lastUpdate')}: ${currentTime}
        </div>
        ${statuses.map(input => `
            <div class="status-card">
                <h3>
                    <span class="status-indicator ${input.active ? 'status-active' : 'status-inactive'}"></span>
                    ${input.name}
                </h3>
                <p><strong>${t('common.path')}:</strong> ${input.url_path}</p>
                <p><strong>${t('common.status')}:</strong> ${input.active ? t('common.active') : t('common.inactive')}</p>
                <p><strong>${t('common.connections')}:</strong> ${input.connections}</p>
                <p><strong>${t('common.errors')}:</strong> ${input.error_count}</p>
                
                ${input.outputs && input.outputs.length > 0 ? `
                    <h4 style="margin-top: 15px; margin-bottom: 10px;">${t('common.outputs')}:</h4>
                    ${input.outputs.map(output => `
                        <div class="output-item">
                            <div class="output-info">
                                <div class="output-url">${output.url}</div>
                                <div class="output-stats">
                                    ${t('common.status')}: <span style="color: ${output.active ? '#27ae60' : '#e74c3c'}; font-weight: bold;">${output.active ? t('common.active') : t('common.inactive')}</span> | 
                                    ${t('common.bitrate')}: <strong>${(output.bitrate_kbps || 0).toFixed(1)} kbps</strong> | 
                                    ${t('common.uptime')}: <strong>${output.uptime || '00:00:00'}</strong> | 
                                    ${t('common.errors')}: <strong>${output.error_count || 0}</strong>
                                </div>
                            </div>
                            <div class="output-actions">
                                <button class="btn btn-warning btn-small" onclick="reconnectOutput('${input.name}', '${output.url}')" title="${t('common.reconnect')}">
                                    🔄
                                </button>
                            </div>
                        </div>
                    `).join('')}
                ` : `<p>${t('status.noOutputs')}</p>`}
            </div>
        `).join('')}
    `;

    container.innerHTML = html;
}

// Загрузка входов
async function loadInputs() {
    try {
        const inputs = await apiRequest('/inputs');
        displayInputs(inputs);
    } catch (error) {
        showError(t('inputs.loadError') + ': ' + error.message);
    }
}

// Отображение входов
function displayInputs(inputs) {
    const container = document.getElementById('inputs-content');

    const html = `
        <div class="form-group">
            <h3>${t('inputs.addNew')}</h3>
            <form id="add-input-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>${t('inputs.inputName')}:</label>
                        <input type="text" class="form-control" id="input-name" placeholder="например: obs" required>
                    </div>
                    <div class="form-group">
                        <label>${t('inputs.rtmpPath')}:</label>
                        <input type="text" class="form-control" id="input-path" placeholder="/live/stream" required>
                    </div>
                </div>
                <div class="form-group">
                    <label>${t('inputs.outputsLabel')}</label>
                    <textarea class="form-control" id="input-outputs" rows="3" placeholder="${t('inputs.outputsPlaceholder')}"></textarea>
                </div>
                <button type="submit" class="btn btn-success">${t('common.add')} ${t('nav.inputs').replace('📥 ', '').toLowerCase()}</button>
            </form>
        </div>
        
        <div class="form-group">
            <h3>${t('inputs.existing')}</h3>
            ${inputs.length === 0 ? `<p>${t('inputs.noInputs')}</p>` : inputs.map(input => `
                <div class="output-item">
                    <div class="output-info">
                        <div class="output-url">${input.Name || input.name}</div>
                        <div class="output-stats">
                            ${t('common.path')}: ${input.URLPath || input.url_path} | 
                            ${t('common.outputs')}: ${(input.Outputs || input.outputs || []).length}
                        </div>
                    </div>
                    <div class="output-actions">
                        <button class="btn btn-danger btn-small" onclick="removeInput('${input.Name || input.name}')">
                            🗑️
                        </button>
                    </div>
                </div>
            `).join('')}
        </div>
    `;

    container.innerHTML = html;

    // Добавляем обработчик формы
    document.getElementById('add-input-form').addEventListener('submit', addInput);
}

// Добавление входа
async function addInput(event) {
    event.preventDefault();

    const name = document.getElementById('input-name').value;
    const path = document.getElementById('input-path').value;
    const outputsText = document.getElementById('input-outputs').value;
    const outputs = outputsText.split('\n').filter(url => url.trim() !== '');

    try {
        await apiRequest('/inputs/add', {
            method: 'POST',
            body: JSON.stringify({
                name: name,
                url_path: path,
                outputs: outputs
            })
        });

        showSuccess(t('inputs.addSuccess'));
        loadInputs();
    } catch (error) {
        showError(t('inputs.addError') + ': ' + error.message);
    }
}

// Удаление входа
async function removeInput(name) {
    if (!confirm(`${t('inputs.removeConfirm')} "${name}"?`)) return;

    try {
        await apiRequest(`/inputs/remove?name=${encodeURIComponent(name)}`);
        showSuccess(t('inputs.removeSuccess'));
        loadInputs();
    } catch (error) {
        showError(t('inputs.removeError') + ': ' + error.message);
    }
}

// Первоначальная загрузка выходов (при переключении на вкладку)
async function loadOutputsInitial() {
    try {
        const inputs = await apiRequest('/inputs');
        const statuses = await apiRequest('/status/all');

        // Создаем форму только один раз или если она не существует
        if (!document.getElementById('add-output-form')) {
            createOutputForm(inputs);
        } else {
            // Обновляем только список входов в форме
            updateOutputFormInputs(inputs);
        }

        // Обновляем статусы выходов
        displayOutputsStatus(statuses);
    } catch (error) {
        showError(t('outputs.loadError') + ': ' + error.message);
        document.getElementById('outputs-status-container').innerHTML = `<div class="loading">${t('outputs.loadError')}...</div>`;
    }
}

// Загрузка выходов (старая функция - теперь используется только для обновления статусов)
async function loadOutputs() {
    try {
        const statuses = await apiRequest('/status/all');
        displayOutputsStatus(statuses);
    } catch (error) {
        showError(t('outputs.loadError') + ': ' + error.message);
        document.getElementById('outputs-status-container').innerHTML = `<div class="loading">${t('outputs.loadError')}...</div>`;
    }
}

// Создание формы добавления выхода (один раз)
function createOutputForm(inputs) {
    const container = document.getElementById('add-output-form-container');

    if (inputs.length === 0) {
        container.innerHTML = `<div class="alert alert-error">${t('outputs.noInputsForManagement')}</div>`;
        return;
    }

    const html = `
        <div class="form-group">
            <h3>${t('outputs.addOutput')}</h3>
            <form id="add-output-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>${t('outputs.input')}:</label>
                        <select class="form-control" id="output-input" required>
                            ${inputs.map(input => `<option value="${input.Name || input.name}">${input.Name || input.name}</option>`).join('')}
                        </select>
                    </div>
                    <div class="form-group">
                        <label>${t('outputs.outputUrl')}:</label>
                        <input type="text" class="form-control" id="output-url" placeholder="srt://example.com:9000" required>
                    </div>
                </div>
                <button type="submit" class="btn btn-success">${t('common.add')} ${t('nav.outputs').replace('📤 ', '').toLowerCase()}</button>
            </form>
        </div>
    `;

    container.innerHTML = html;

    // Добавляем обработчик формы
    const form = document.getElementById('add-output-form');
    if (form) {
        form.addEventListener('submit', addOutput);
    }
}

// Обновление списка входов в форме (без сброса значений)
function updateOutputFormInputs(inputs) {
    const select = document.getElementById('output-input');
    if (!select) return;

    const currentValue = select.value;

    // Обновляем опции
    select.innerHTML = inputs.map(input =>
        `<option value="${input.Name || input.name}" ${(input.Name || input.name) === currentValue ? 'selected' : ''}>${input.Name || input.name}</option>`
    ).join('');
}

// Отображение статусов выходов (обновляемая часть)
function displayOutputsStatus(statuses) {
    const container = document.getElementById('outputs-status-container');
    const currentTime = new Date().toLocaleTimeString();

    const html = `
        <div class="form-group">
            <h3>${t('outputs.management')} <small style="color: #6c757d;">(${t('common.updated')}: ${currentTime})</small></h3>
            ${statuses.length === 0 ? `<p>${t('outputs.noActiveInputs')}</p>` : statuses.map(input => `
                <div class="status-card">
                    <h3>
                        <span class="status-indicator ${input.active ? 'status-active' : 'status-inactive'}"></span>
                        ${input.name}
                    </h3>
                    ${input.outputs && input.outputs.length > 0 ? input.outputs.map(output => `
                        <div class="output-item">
                            <div class="output-info">
                                <div class="output-url">${output.url}</div>
                                <div class="output-stats">
                                    ${t('common.status')}: <span style="color: ${output.active ? '#27ae60' : '#e74c3c'}; font-weight: bold;">${output.active ? t('common.active') : t('common.inactive')}</span> | 
                                    ${t('common.bitrate')}: <strong>${(output.bitrate_kbps || 0).toFixed(1)} kbps</strong> | 
                                    ${t('common.uptime')}: <strong>${output.uptime || '00:00:00'}</strong> | 
                                    ${t('common.errors')}: <strong>${output.error_count || 0}</strong>
                                </div>
                            </div>
                            <div class="output-actions">
                                <button class="btn btn-warning btn-small" onclick="reconnectOutput('${input.name}', '${output.url}')" title="${t('common.reconnect')}">
                                    🔄
                                </button>
                                <button class="btn btn-danger btn-small" onclick="removeOutput('${input.name}', '${output.url}')" title="${t('common.remove')}">
                                    🗑️
                                </button>
                            </div>
                        </div>
                    `).join('') : `<p>${t('status.noOutputs')}</p>`}
                </div>
            `).join('')}
        </div>
    `;

    container.innerHTML = html;
}

// Добавление выхода
async function addOutput(event) {
    event.preventDefault();

    const inputName = document.getElementById('output-input').value;
    const url = document.getElementById('output-url').value;

    try {
        await apiRequest('/outputs/add', {
            method: 'POST',
            body: JSON.stringify({
                name: inputName,
                url: url
            })
        });

        showSuccess(t('outputs.addSuccess'));

        // Очищаем только поле URL
        document.getElementById('output-url').value = '';

        // Обновляем только статусы выходов, не трогаем форму
        loadOutputs();
    } catch (error) {
        showError(t('outputs.addError') + ': ' + error.message);
    }
}

// Удаление выхода
async function removeOutput(inputName, url) {
    if (!confirm(`${t('outputs.removeConfirm')} "${url}"?`)) return;

    try {
        await apiRequest('/outputs/remove', {
            method: 'POST',
            body: JSON.stringify({
                name: inputName,
                url: url
            })
        });

        showSuccess(t('outputs.removeSuccess'));

        // Обновляем только статусы выходов, не трогаем форму
        loadOutputs();
    } catch (error) {
        showError(t('outputs.removeError') + ': ' + error.message);
    }
}

// Реконнект выхода
async function reconnectOutput(inputName, url) {
    try {
        await apiRequest('/outputs/reconnect', {
            method: 'POST',
            body: JSON.stringify({
                name: inputName,
                url: url
            })
        });

        showSuccess(t('outputs.reconnectInitiated'));

        // Обновляем только статусы выходов через короткое время
        setTimeout(() => loadOutputs(), 1000);
    } catch (error) {
        showError(t('outputs.reconnectError') + ': ' + error.message);
    }
}

// Загрузка настроек
async function loadSettings() {
    try {
        const settings = await apiRequest('/settings');
        displaySettings(settings);
    } catch (error) {
        showError(t('settings.loadError') + ': ' + error.message);
    }
}

// Отображение настроек
function displaySettings(settings) {
    const container = document.getElementById('settings-content');

    const html = `
        <form id="settings-form">
            <div class="form-group">
                <h3>${t('settings.srt')}</h3>
                <div class="form-row">
                    <div class="form-group">
                        <label>${t('settings.latency')}</label>
                        <input type="number" class="form-control" id="srt-latency" value="${settings.SRTSettings?.Latency || settings.srt_settings?.latency || 120}" min="0">
                    </div>
                    <div class="form-group">
                        <label>${t('settings.timeout')}</label>
                        <input type="number" class="form-control" id="srt-timeout" value="${settings.SRTSettings?.ConnectTimeout || settings.srt_settings?.connect_timeout || 5000}" min="0">
                    </div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label>${t('settings.passphrase')}</label>
                        <input type="text" class="form-control" id="srt-passphrase" value="${settings.SRTSettings?.Passphrase || settings.srt_settings?.passphrase || ''}">
                    </div>
                    <div class="form-group">
                        <label>${t('settings.streamId')}</label>
                        <input type="text" class="form-control" id="srt-streamid" value="${settings.SRTSettings?.StreamID || settings.srt_settings?.streamid || ''}">
                    </div>
                </div>
            </div>
            
            <div class="form-group">
                <h3>${t('settings.logging')}</h3>
                <div class="form-row">
                    <div class="form-group">
                        <label>
                            <input type="checkbox" id="log-to-file" ${settings.LogToFile || settings.log_to_file ? 'checked' : ''}>
                            ${t('settings.logToFile')}
                        </label>
                    </div>
                    <div class="form-group">
                        <label>${t('settings.logFile')}</label>
                        <input type="text" class="form-control" id="log-file" value="${settings.LogFile || settings.log_file || 'server.log'}">
                    </div>
                </div>
            </div>
            
            <div class="form-group">
                <h3>${t('settings.reconnectInterval')}</h3>
                <div class="form-group">
                    <label>${t('settings.reconnectIntervalLabel')}</label>
                    <input type="number" class="form-control" id="reconnect-interval" value="${settings.ReconnectInterval || settings.reconnect_interval || 5}" min="1">
                </div>
            </div>
            
            <div class="form-group">
                <button type="submit" class="btn btn-success">${t('common.save')}</button>
                <button type="button" class="btn btn-warning" onclick="reloadSettings()">${t('settings.reload')}</button>
            </div>
        </form>
    `;

    container.innerHTML = html;

    // Добавляем обработчик формы
    document.getElementById('settings-form').addEventListener('submit', saveSettings);
}

// Сохранение настроек
async function saveSettings(event) {
    event.preventDefault();

    const settings = {
        srt_settings: {
            latency: parseInt(document.getElementById('srt-latency').value),
            connect_timeout: parseInt(document.getElementById('srt-timeout').value),
            passphrase: document.getElementById('srt-passphrase').value,
            streamid: document.getElementById('srt-streamid').value,
            encryption: 'none'
        },
        log_to_file: document.getElementById('log-to-file').checked,
        log_file: document.getElementById('log-file').value,
        reconnect_interval: parseInt(document.getElementById('reconnect-interval').value)
    };

    try {
        await apiRequest('/settings', {
            method: 'PUT',
            body: JSON.stringify(settings)
        });

        showSuccess(t('settings.saveSuccess'));
    } catch (error) {
        showError(t('settings.saveError') + ': ' + error.message);
    }
}

// Перезагрузка настроек
async function reloadSettings() {
    try {
        await apiRequest('/settings/reload', { method: 'POST' });
        showSuccess(t('settings.reloadSuccess'));
        loadSettings();
    } catch (error) {
        showError(t('settings.reloadError') + ': ' + error.message);
    }
}

// Показ уведомлений
function showSuccess(message) {
    showAlert(message, 'success');
}

function showError(message) {
    showAlert(message, 'error');
}

function showAlert(message, type) {
    const alert = document.createElement('div');
    alert.className = `alert alert-${type}`;
    alert.textContent = message;

    const container = document.querySelector('.content');
    container.insertBefore(alert, container.firstChild);

    setTimeout(() => {
        alert.remove();
    }, 5000);
} 