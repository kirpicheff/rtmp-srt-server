// Конфигурация API
const API_BASE = '/api';
const API_CREDENTIALS = btoa('admin:secret'); // В продакшене лучше получать через форму входа

// Глобальные переменные
let currentTab = 'status';
let refreshInterval;

// Инициализация приложения
document.addEventListener('DOMContentLoaded', function () {
    setupNavigation();
    loadCurrentTab();
    startAutoRefresh();
});

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
            loadOutputs();
            break;
        case 'settings':
            loadSettings();
            break;
    }
}

// Автообновление данных
function startAutoRefresh() {
    refreshInterval = setInterval(() => {
        if (currentTab === 'status') {
            loadStatus();
        } else if (currentTab === 'outputs') {
            loadOutputs();
        }
    }, 5000); // Обновляем статус каждые 5 секунд
}

// Обновление текущей вкладки
function refreshCurrentTab() {
    loadCurrentTab();
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
        showError('Ошибка загрузки статуса: ' + error.message);
    }
}

// Отображение статуса
function displayStatus(statuses) {
    const container = document.getElementById('status-content');
    const currentTime = new Date().toLocaleTimeString();

    if (!statuses || statuses.length === 0) {
        container.innerHTML = '<div class="alert alert-error">Нет активных входов</div>';
        return;
    }

    const html = `
        <div style="margin-bottom: 20px; color: #6c757d; font-size: 0.9em;">
            Последнее обновление: ${currentTime}
        </div>
        ${statuses.map(input => `
            <div class="status-card">
                <h3>
                    <span class="status-indicator ${input.active ? 'status-active' : 'status-inactive'}"></span>
                    ${input.name}
                </h3>
                <p><strong>Путь:</strong> ${input.url_path}</p>
                <p><strong>Статус:</strong> ${input.active ? 'Активен' : 'Неактивен'}</p>
                <p><strong>Подключения:</strong> ${input.connections}</p>
                <p><strong>Ошибки:</strong> ${input.error_count}</p>
                
                ${input.outputs && input.outputs.length > 0 ? `
                    <h4 style="margin-top: 15px; margin-bottom: 10px;">Выходы:</h4>
                    ${input.outputs.map(output => `
                        <div class="output-item">
                            <div class="output-info">
                                <div class="output-url">${output.url}</div>
                                <div class="output-stats">
                                    Статус: <span style="color: ${output.active ? '#27ae60' : '#e74c3c'}; font-weight: bold;">${output.active ? 'Активен' : 'Неактивен'}</span> | 
                                    Битрейт: <strong>${(output.bitrate_kbps || 0).toFixed(1)} kbps</strong> | 
                                    Аптайм: <strong>${output.uptime || '00:00:00'}</strong> | 
                                    Ошибки: <strong>${output.error_count || 0}</strong>
                                </div>
                            </div>
                            <div class="output-actions">
                                <button class="btn btn-warning btn-small" onclick="reconnectOutput('${input.name}', '${output.url}')" title="Реконнект">
                                    🔄
                                </button>
                            </div>
                        </div>
                    `).join('')}
                ` : '<p>Нет выходов</p>'}
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
        showError('Ошибка загрузки входов: ' + error.message);
    }
}

// Отображение входов
function displayInputs(inputs) {
    const container = document.getElementById('inputs-content');

    const html = `
        <div class="form-group">
            <h3>Добавить новый вход</h3>
            <form id="add-input-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>Имя входа:</label>
                        <input type="text" class="form-control" id="input-name" placeholder="например: obs" required>
                    </div>
                    <div class="form-group">
                        <label>Путь RTMP:</label>
                        <input type="text" class="form-control" id="input-path" placeholder="/live/stream" required>
                    </div>
                </div>
                <div class="form-group">
                    <label>Выходы (по одному на строку):</label>
                    <textarea class="form-control" id="input-outputs" rows="3" placeholder="srt://example.com:9000&#10;rtmp://example.com/live/stream"></textarea>
                </div>
                <button type="submit" class="btn btn-success">Добавить вход</button>
            </form>
        </div>
        
        <div class="form-group">
            <h3>Существующие входы</h3>
            ${inputs.length === 0 ? '<p>Нет входов</p>' : inputs.map(input => `
                <div class="output-item">
                    <div class="output-info">
                        <div class="output-url">${input.Name || input.name}</div>
                        <div class="output-stats">
                            Путь: ${input.URLPath || input.url_path} | 
                            Выходов: ${(input.Outputs || input.outputs || []).length}
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

        showSuccess('Вход успешно добавлен');
        loadInputs();
    } catch (error) {
        showError('Ошибка добавления входа: ' + error.message);
    }
}

// Удаление входа
async function removeInput(name) {
    if (!confirm(`Удалить вход "${name}"?`)) return;

    try {
        await apiRequest(`/inputs/remove?name=${encodeURIComponent(name)}`);
        showSuccess('Вход успешно удален');
        loadInputs();
    } catch (error) {
        showError('Ошибка удаления входа: ' + error.message);
    }
}

// Загрузка выходов
async function loadOutputs() {
    try {
        const inputs = await apiRequest('/inputs');
        const statuses = await apiRequest('/status/all');
        displayOutputs(inputs, statuses);
    } catch (error) {
        showError('Ошибка загрузки выходов: ' + error.message);
        // Показываем индикатор загрузки при ошибке
        document.getElementById('outputs-content').innerHTML = '<div class="loading">Ошибка загрузки данных...</div>';
    }
}

// Отображение выходов
function displayOutputs(inputs, statuses) {
    const container = document.getElementById('outputs-content');

    if (inputs.length === 0) {
        container.innerHTML = '<div class="alert alert-error">Нет входов для управления выходами</div>';
        return;
    }

    const currentTime = new Date().toLocaleTimeString();

    const html = `
        <div class="form-group">
            <h3>Добавить выход</h3>
            <form id="add-output-form">
                <div class="form-row">
                    <div class="form-group">
                        <label>Вход:</label>
                        <select class="form-control" id="output-input" required>
                            ${inputs.map(input => `<option value="${input.Name || input.name}">${input.Name || input.name}</option>`).join('')}
                        </select>
                    </div>
                    <div class="form-group">
                        <label>URL выхода:</label>
                        <input type="text" class="form-control" id="output-url" placeholder="srt://example.com:9000" required>
                    </div>
                </div>
                <button type="submit" class="btn btn-success">Добавить выход</button>
            </form>
        </div>
        
        <div class="form-group">
            <h3>Управление выходами <small style="color: #6c757d;">(обновлено: ${currentTime})</small></h3>
            ${statuses.length === 0 ? '<p>Нет активных входов</p>' : statuses.map(input => `
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
                                    Статус: <span style="color: ${output.active ? '#27ae60' : '#e74c3c'}; font-weight: bold;">${output.active ? 'Активен' : 'Неактивен'}</span> | 
                                    Битрейт: <strong>${(output.bitrate_kbps || 0).toFixed(1)} kbps</strong> | 
                                    Аптайм: <strong>${output.uptime || '00:00:00'}</strong> | 
                                    Ошибки: <strong>${output.error_count || 0}</strong>
                                </div>
                            </div>
                            <div class="output-actions">
                                <button class="btn btn-warning btn-small" onclick="reconnectOutput('${input.name}', '${output.url}')" title="Реконнект">
                                    🔄
                                </button>
                                <button class="btn btn-danger btn-small" onclick="removeOutput('${input.name}', '${output.url}')" title="Удалить">
                                    🗑️
                                </button>
                            </div>
                        </div>
                    `).join('') : '<p>Нет выходов</p>'}
                </div>
            `).join('')}
        </div>
    `;

    container.innerHTML = html;

    // Добавляем обработчик формы
    const form = document.getElementById('add-output-form');
    if (form) {
        form.addEventListener('submit', addOutput);
    }
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

        showSuccess('Выход успешно добавлен');
        loadOutputs();
    } catch (error) {
        showError('Ошибка добавления выхода: ' + error.message);
    }
}

// Удаление выхода
async function removeOutput(inputName, url) {
    if (!confirm(`Удалить выход "${url}"?`)) return;

    try {
        await apiRequest('/outputs/remove', {
            method: 'POST',
            body: JSON.stringify({
                name: inputName,
                url: url
            })
        });

        showSuccess('Выход успешно удален');
        loadOutputs();
    } catch (error) {
        showError('Ошибка удаления выхода: ' + error.message);
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

        showSuccess('Реконнект инициирован');
        setTimeout(() => loadOutputs(), 1000);
    } catch (error) {
        showError('Ошибка реконнекта: ' + error.message);
    }
}

// Загрузка настроек
async function loadSettings() {
    try {
        const settings = await apiRequest('/settings');
        displaySettings(settings);
    } catch (error) {
        showError('Ошибка загрузки настроек: ' + error.message);
    }
}

// Отображение настроек
function displaySettings(settings) {
    const container = document.getElementById('settings-content');

    const html = `
        <form id="settings-form">
            <div class="form-group">
                <h3>Настройки SRT</h3>
                <div class="form-row">
                    <div class="form-group">
                        <label>Latency (мс):</label>
                        <input type="number" class="form-control" id="srt-latency" value="${settings.SRTSettings?.Latency || settings.srt_settings?.latency || 120}" min="0">
                    </div>
                    <div class="form-group">
                        <label>Connect Timeout (мс):</label>
                        <input type="number" class="form-control" id="srt-timeout" value="${settings.SRTSettings?.ConnectTimeout || settings.srt_settings?.connect_timeout || 5000}" min="0">
                    </div>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label>Passphrase:</label>
                        <input type="text" class="form-control" id="srt-passphrase" value="${settings.SRTSettings?.Passphrase || settings.srt_settings?.passphrase || ''}">
                    </div>
                    <div class="form-group">
                        <label>Stream ID:</label>
                        <input type="text" class="form-control" id="srt-streamid" value="${settings.SRTSettings?.StreamID || settings.srt_settings?.streamid || ''}">
                    </div>
                </div>
            </div>
            
            <div class="form-group">
                <h3>Настройки логирования</h3>
                <div class="form-row">
                    <div class="form-group">
                        <label>
                            <input type="checkbox" id="log-to-file" ${settings.LogToFile || settings.log_to_file ? 'checked' : ''}>
                            Логировать в файл
                        </label>
                    </div>
                    <div class="form-group">
                        <label>Файл лога:</label>
                        <input type="text" class="form-control" id="log-file" value="${settings.LogFile || settings.log_file || 'server.log'}">
                    </div>
                </div>
            </div>
            
            <div class="form-group">
                <h3>Интервал переподключения</h3>
                <div class="form-group">
                    <label>Интервал (секунды):</label>
                    <input type="number" class="form-control" id="reconnect-interval" value="${settings.ReconnectInterval || settings.reconnect_interval || 5}" min="1">
                </div>
            </div>
            
            <div class="form-group">
                <button type="submit" class="btn btn-success">Сохранить настройки</button>
                <button type="button" class="btn btn-warning" onclick="reloadSettings()">Перезагрузить из файла</button>
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

        showSuccess('Настройки успешно сохранены');
    } catch (error) {
        showError('Ошибка сохранения настроек: ' + error.message);
    }
}

// Перезагрузка настроек
async function reloadSettings() {
    try {
        await apiRequest('/settings/reload', { method: 'POST' });
        showSuccess('Настройки перезагружены из файла');
        loadSettings();
    } catch (error) {
        showError('Ошибка перезагрузки настроек: ' + error.message);
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