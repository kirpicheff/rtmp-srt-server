/* ============================================================
 * RTMP/SRT/WHIP Server Manager
 * Vanilla modules (IIFE). No bundler, no framework.
 * ============================================================ */

/* -------------------- I18N -------------------- */
const I18N = (() => {
    const translations = {
        ru: {
            'nav.status': 'Дашборд',
            'nav.inputs': 'Входы',
            'nav.settings': 'Настройки',
            'nav.help': 'Помощь',

            'common.online': 'Online',
            'common.offline': 'Offline',
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
            'common.lastUpdate': 'обновлено',
            'common.updated': 'обновлено',
            'common.palette': 'Поиск',
            'common.copy': 'Копировать',
            'common.copied': 'Скопировано',
            'common.confirmRemove': 'Удалить?',
            'common.secondsAgo': 'с назад',
            'common.justNow': 'только что',
            'common.minutesAgo': 'мин назад',

            'topbar.summary': '{active}/{total} вх. · {outActive}/{outTotal} вых. · {conn} подкл.',

            'status.title': 'Панель мониторинга',
            'status.noInputs': 'Нет входов. Создайте первый вход во вкладке «Входы».',
            'status.noOutputs': 'Нет выходов',
            'status.loadError': 'Ошибка загрузки',
            'status.totalBitrate': 'Сумм. битрейт',
            'status.lastError': 'Последняя ошибка',

            'inputs.title': 'Управление входами',
            'inputs.subtitle': 'Добавляйте и удаляйте входы, просматривайте их выходы',
            'inputs.addNew': 'Новый вход',
            'inputs.inputName': 'Имя входа',
            'inputs.rtmpPath': 'Путь',
            'inputs.protocol': 'Протокол',
            'inputs.outputsPlaceholder': 'srt://example.com:9000\nrtmp://example.com/live/stream',
            'inputs.outputsLabel': 'Выходы (по одному на строку):',
            'inputs.existing': 'Все входы',
            'inputs.noInputs': 'Входов пока нет',
            'inputs.selectHint': 'Выберите вход слева или создайте новый',
            'inputs.addSuccess': 'Вход добавлен',
            'inputs.removeSuccess': 'Вход удалён',
            'inputs.removeConfirm': 'Удалить вход',
            'inputs.addError': 'Ошибка добавления входа',
            'inputs.removeError': 'Ошибка удаления входа',
            'inputs.loadError': 'Ошибка загрузки входов',
            'inputs.outputs': 'Выходы входа',
            'inputs.confirmRemoveBody': 'Вход будет остановлен, все активные выходы отключены.',

            'outputs.addSuccess': 'Выход добавлен',
            'outputs.removeSuccess': 'Выход удалён',
            'outputs.reconnectInitiated': 'Реконнект инициирован',
            'outputs.removeConfirm': 'Удалить выход',
            'outputs.addError': 'Ошибка добавления выхода',
            'outputs.removeError': 'Ошибка удаления выхода',
            'outputs.reconnectError': 'Ошибка реконнекта',
            'outputs.urlPlaceholder': 'srt:// · rtmp:// · file://',
            'outputs.editTitle': 'Редактировать выход',
            'outputs.urlLabel': 'URL выхода:',
            'outputs.editSuccess': 'Выход изменён',
            'outputs.editError': 'Ошибка изменения выхода',

            'settings.title': 'Глобальные настройки',
            'settings.subtitle': 'Параметры SRT, логов, интерфейса и реконнекта',
            'settings.unsaved': 'не сохранено',
            'settings.section.srt': 'SRT',
            'settings.section.logging': 'Логи и интерфейс',
            'settings.section.reconnect': 'Реконнект',
            'settings.srt': 'Настройки SRT',
            'settings.latency': 'Latency (мс):',
            'settings.timeout': 'Connect Timeout (мс):',
            'settings.passphrase': 'Passphrase:',
            'settings.passphraseHelp': 'Шифрование потока SRT. Оставьте пустым для отключения.',
            'settings.streamId': 'Stream ID:',
            'settings.streamIdHelp': 'Идентификатор потока, передаётся приёмнику.',
            'settings.logging': 'Настройки логов',
            'settings.logToFile': 'Логировать в файл',
            'settings.logFile': 'Файл лога:',
            'settings.minimizeToTray': 'Сворачивать в трей при закрытии',
            'settings.reconnectInterval': 'Интервал переподключения',
            'settings.reconnectIntervalLabel': 'Интервал (секунды):',
            'settings.reconnectIntervalHelp': 'Через сколько секунд повторять подключение упавших выходов.',
            'settings.saveSuccess': 'Настройки сохранены',
            'settings.reloadSuccess': 'Настройки перезагружены',
            'settings.saveError': 'Ошибка сохранения',
            'settings.reloadError': 'Ошибка перезагрузки',
            'settings.loadError': 'Ошибка загрузки',
            'settings.reload': 'Перезагрузить из файла',
            'settings.revert': 'Отменить изменения',
            'settings.dirtyHint': 'Есть несохранённые изменения',

            'help.title': 'Справка и руководство',
            'help.subtitle': 'Быстрый старт, протоколы, форматы выходов',
            'help.quickstart.title': 'Быстрый старт',
            'help.quickstart.body': 'Сервер работает как медиа-шлюз: принимает поток по одному протоколу (RTMP, SRT, WHIP) и одновременно транслирует его на несколько выходов (другие серверы, файлы). Не требует перекодирования — потоки передаются без изменений (copy).',
            'help.inputs.title': 'Как добавить вход',
            'help.inputs.body': 'Перейдите во вкладку «Входы». Заполните поля:',
            'help.inputs.field1': 'Имя — произвольный идентификатор входа (например: <code>obs</code>, <code>camera1</code>).',
            'help.inputs.field2': 'Путь — RTMP-путь, по которому encoder будет публиковать поток. Пример: <code>/live/stream</code>.',
            'help.inputs.field3': 'Выходы — список URL для ретрансляции (по одному на строку). Можно оставить пустым и добавить позже прямо с Дашборда.',
            'help.inputs.note': 'Вход можно создать до того, как encoder подключится. Статус «Неактивен» — норма при отсутствии публикующего потока.',
            'help.protocols.title': 'Поддерживаемые протоколы входа',
            'help.protocols.rtmp': '<b>RTMP</b> — стандартный протокол. Укажите путь в настройках OBS или другого encoder-а. Порт по умолчанию: <code>1935</code>.',
            'help.protocols.srt': '<b>SRT</b> — надёжный UDP-протокол для нестабильных сетей. Укажите адрес сервера и порт в encoder-е. Порт по умолчанию: <code>4000</code> (настраивается в конфиге).',
            'help.protocols.whip': '<b>WHIP</b> — WebRTC-протокол. Используется из браузерных encoder-ов. Endpoint: <code>/whip/&lt;path&gt;</code>.',
            'help.outputs.title': 'Как добавить выход',
            'help.outputs.body': 'Выход можно добавить двумя способами:',
            'help.outputs.way1': '<b>На Дашборде</b> — прямо в карточке потока введите URL в поле и нажмите «Добавить». Выход активируется немедленно.',
            'help.outputs.way2': '<b>При создании входа</b> — укажите список URL в поле «Выходы» при добавлении нового входа.',
            'help.outputs.formats': 'Примеры поддерживаемых форматов выходов:',
            'help.outputs.ex.rtmp': 'Ретрансляция на RTMP-сервер (YouTube, Twitch и др.)',
            'help.outputs.ex.srt': 'Ретрансляция по SRT',
            'help.outputs.ex.flv': 'Запись в файл FLV (нативно, ffmpeg не нужен)',
            'help.outputs.ex.mp4': 'Запись в файл MP4 (требует ffmpeg в PATH или bin/)',
            'help.outputs.mp4note': 'Для записи MP4 необходим ffmpeg. FLV пишется нативно без внешних зависимостей.',
            'help.ffmpeg.title': 'ffmpeg для записи файлов',
            'help.ffmpeg.body': 'Запись в MP4 использует ffmpeg как внешний процесс (stdin pipe). Бинарный файл ffmpeg должен находиться в папке <code>bin/</code> рядом с сервером или быть доступен в системном PATH. FLV запись не требует ffmpeg.',
            'help.ffmpeg.windows': 'Windows: скачайте готовую сборку с <b>gyan.dev</b> (раздел <i>release builds → essentials</i>). Разархивируйте, скопируйте <code>ffmpeg.exe</code> в папку <code>bin/</code>.',
            'help.ffmpeg.linux': 'Linux: установите через пакетный менеджер (<code>apt install ffmpeg</code>) или скачайте статический бинарник.',
            'help.ffmpeg.custom': 'Для работы WHIP требуется ffmpeg с поддержкой <code>libfdk_aac</code>.',
            'help.reconnect.title': 'Реконнект выходов',
            'help.reconnect.body': 'Если выход потерял соединение, он автоматически переподключается через интервал, заданный в Настройках. Принудительный реконнект доступен кнопкой на Дашборде (иконка «стрелки»).',
            'help.config.title': 'Файл конфигурации',
            'help.config.body': 'Настройки хранятся в файле <code>config.yaml</code> рядом с бинарным файлом сервера. Изменения через веб-интерфейс сохраняются в этот файл автоматически. Ручное редактирование также поддерживается — нажмите «Перезагрузить из файла» в разделе Настройки после правки вручную.',

            'palette.placeholder': 'Поиск входов, выходов, действий…',
            'palette.section.actions': 'Действия',
            'palette.section.inputs': 'Входы',
            'palette.section.outputs': 'Выходы',
            'palette.noResults': 'Ничего не найдено',
            'palette.action.addInput': 'Добавить вход',
            'palette.action.settings': 'Открыть настройки',
            'palette.action.help': 'Открыть помощь',
            'palette.action.refresh': 'Обновить данные',
            'palette.action.reloadConfig': 'Перезагрузить конфиг',
        },
        en: {
            'nav.status': 'Dashboard',
            'nav.inputs': 'Inputs',
            'nav.settings': 'Settings',
            'nav.help': 'Help',

            'common.online': 'Online',
            'common.offline': 'Offline',
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
            'common.lastUpdate': 'updated',
            'common.updated': 'updated',
            'common.palette': 'Search',
            'common.copy': 'Copy',
            'common.copied': 'Copied',
            'common.confirmRemove': 'Remove?',
            'common.secondsAgo': 's ago',
            'common.justNow': 'just now',
            'common.minutesAgo': 'm ago',

            'topbar.summary': '{active}/{total} in · {outActive}/{outTotal} out · {conn} conn.',

            'status.title': 'Monitoring',
            'status.noInputs': 'No inputs. Create your first input in the «Inputs» tab.',
            'status.noOutputs': 'No outputs',
            'status.loadError': 'Load error',
            'status.totalBitrate': 'Total bitrate',
            'status.lastError': 'Last error',

            'inputs.title': 'Input Management',
            'inputs.subtitle': 'Add, remove inputs and inspect their outputs',
            'inputs.addNew': 'New input',
            'inputs.inputName': 'Input name',
            'inputs.rtmpPath': 'Path',
            'inputs.protocol': 'Protocol',
            'inputs.outputsPlaceholder': 'srt://example.com:9000\nrtmp://example.com/live/stream',
            'inputs.outputsLabel': 'Outputs (one per line):',
            'inputs.existing': 'All inputs',
            'inputs.noInputs': 'No inputs yet',
            'inputs.selectHint': 'Select an input on the left or create a new one',
            'inputs.addSuccess': 'Input added',
            'inputs.removeSuccess': 'Input removed',
            'inputs.removeConfirm': 'Remove input',
            'inputs.addError': 'Error adding input',
            'inputs.removeError': 'Error removing input',
            'inputs.loadError': 'Error loading inputs',
            'inputs.outputs': 'Input outputs',
            'inputs.confirmRemoveBody': 'The input will be stopped and all active outputs disconnected.',

            'outputs.addSuccess': 'Output added',
            'outputs.removeSuccess': 'Output removed',
            'outputs.reconnectInitiated': 'Reconnect initiated',
            'outputs.removeConfirm': 'Remove output',
            'outputs.addError': 'Error adding output',
            'outputs.removeError': 'Error removing output',
            'outputs.reconnectError': 'Reconnect error',
            'outputs.urlPlaceholder': 'srt:// · rtmp:// · file://',
            'outputs.editTitle': 'Edit Output',
            'outputs.urlLabel': 'Output URL:',
            'outputs.editSuccess': 'Output modified',
            'outputs.editError': 'Error modifying output',

            'settings.title': 'Global settings',
            'settings.subtitle': 'SRT, logs, UI and reconnect parameters',
            'settings.unsaved': 'unsaved',
            'settings.section.srt': 'SRT',
            'settings.section.logging': 'Logs & UI',
            'settings.section.reconnect': 'Reconnect',
            'settings.srt': 'SRT settings',
            'settings.latency': 'Latency (ms):',
            'settings.timeout': 'Connect timeout (ms):',
            'settings.passphrase': 'Passphrase:',
            'settings.passphraseHelp': 'SRT stream encryption. Leave empty to disable.',
            'settings.streamId': 'Stream ID:',
            'settings.streamIdHelp': 'Stream identifier, passed to the receiver.',
            'settings.logging': 'Log settings',
            'settings.logToFile': 'Log to file',
            'settings.logFile': 'Log file:',
            'settings.minimizeToTray': 'Minimize to tray on close',
            'settings.reconnectInterval': 'Reconnect interval',
            'settings.reconnectIntervalLabel': 'Interval (seconds):',
            'settings.reconnectIntervalHelp': 'How often to retry failed outputs.',
            'settings.saveSuccess': 'Settings saved',
            'settings.reloadSuccess': 'Settings reloaded',
            'settings.saveError': 'Save error',
            'settings.reloadError': 'Reload error',
            'settings.loadError': 'Load error',
            'settings.reload': 'Reload from file',
            'settings.revert': 'Revert',
            'settings.dirtyHint': 'You have unsaved changes',

            'help.title': 'Help & Guide',
            'help.subtitle': 'Quick start, protocols, output formats',
            'help.quickstart.title': 'Quick Start',
            'help.quickstart.body': 'The server works as a media gateway: it receives a stream via one protocol (RTMP, SRT, WHIP) and simultaneously relays it to multiple outputs (other servers, files). No transcoding — streams are passed through without modification (copy).',
            'help.inputs.title': 'How to Add an Input',
            'help.inputs.body': 'Go to the «Inputs» tab. Fill in:',
            'help.inputs.field1': 'Name — any identifier (e.g. <code>obs</code>, <code>camera1</code>).',
            'help.inputs.field2': 'Path — the RTMP path your encoder will publish to. Example: <code>/live/stream</code>.',
            'help.inputs.field3': 'Outputs — list of destination URLs (one per line). Can be left empty and added later from the Dashboard.',
            'help.inputs.note': 'An input can be created before the encoder connects. Status «Inactive» is normal when no stream is being published.',
            'help.protocols.title': 'Supported Input Protocols',
            'help.protocols.rtmp': '<b>RTMP</b> — standard protocol. Default port: <code>1935</code>.',
            'help.protocols.srt': '<b>SRT</b> — reliable UDP for unstable networks. Default port: <code>4000</code> (configurable).',
            'help.protocols.whip': '<b>WHIP</b> — WebRTC-based protocol. Endpoint: <code>/whip/&lt;path&gt;</code>.',
            'help.outputs.title': 'How to Add an Output',
            'help.outputs.body': 'Outputs can be added in two ways:',
            'help.outputs.way1': '<b>On the Dashboard</b> — enter the URL directly in the stream card and click «Add».',
            'help.outputs.way2': '<b>When creating an input</b> — specify a list of URLs in the «Outputs» field.',
            'help.outputs.formats': 'Examples of supported output formats:',
            'help.outputs.ex.rtmp': 'Relay to RTMP server (YouTube, Twitch, etc.)',
            'help.outputs.ex.srt': 'Relay via SRT',
            'help.outputs.ex.flv': 'Record to FLV file (native, no ffmpeg required)',
            'help.outputs.ex.mp4': 'Record to MP4 file (requires ffmpeg in PATH or bin/)',
            'help.outputs.mp4note': 'MP4 recording requires ffmpeg. FLV is written natively.',
            'help.ffmpeg.title': 'ffmpeg for File Recording',
            'help.ffmpeg.body': 'MP4 recording uses ffmpeg as an external process. The ffmpeg binary must be located in the <code>bin/</code> folder or in the system PATH.',
            'help.ffmpeg.windows': 'Windows: download a pre-built binary from <b>gyan.dev</b> (release builds → essentials).',
            'help.ffmpeg.linux': 'Linux: install via package manager or download a static binary.',
            'help.ffmpeg.custom': 'WHIP recording requires ffmpeg built with <code>libfdk_aac</code> support.',
            'help.reconnect.title': 'Output Reconnection',
            'help.reconnect.body': 'If an output loses its connection, it auto-reconnects after the interval set in Settings. A manual reconnect is available via the button on the Dashboard.',
            'help.config.title': 'Configuration File',
            'help.config.body': 'Settings are stored in <code>config.yaml</code> next to the server binary. Changes via the web interface are saved automatically. Click «Reload from file» after manual edits.',

            'palette.placeholder': 'Search inputs, outputs, actions…',
            'palette.section.actions': 'Actions',
            'palette.section.inputs': 'Inputs',
            'palette.section.outputs': 'Outputs',
            'palette.noResults': 'No results',
            'palette.action.addInput': 'Add input',
            'palette.action.settings': 'Open settings',
            'palette.action.help': 'Open help',
            'palette.action.refresh': 'Refresh data',
            'palette.action.reloadConfig': 'Reload config',
        },
    };

    let current = localStorage.getItem('language') || 'ru';

    const t = (key, vars) => {
        let s = translations[current][key] || key;
        if (vars) {
            for (const k in vars) {
                s = s.replace(new RegExp('\\{' + k + '\\}', 'g'), vars[k]);
            }
        }
        return s;
    };

    const apply = () => {
        document.querySelectorAll('[data-i18n]').forEach(el => {
            const key = el.getAttribute('data-i18n');
            el.textContent = t(key);
        });
        document.querySelectorAll('[data-i18n-title]').forEach(el => {
            el.title = t(el.getAttribute('data-i18n-title'));
        });
        document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
            el.placeholder = t(el.getAttribute('data-i18n-placeholder'));
        });
    };

    const set = (lang) => {
        current = lang;
        localStorage.setItem('language', lang);
        apply();
    };

    const get = () => current;

    return { t, apply, set, get };
})();

/* -------------------- Api -------------------- */
const Api = (() => {
    const BASE = '/api';

    const request = async (endpoint, options = {}) => {
        const url = `${BASE}${endpoint}`;
        const response = await fetch(url, {
            headers: { 'Content-Type': 'application/json' },
            ...options,
            headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
        });
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        const ct = response.headers.get('content-type') || '';
        if (ct.includes('application/json')) {
            return await response.json();
        }
        return { success: true };
    };

    return {
        statusAll: () => request('/status/all'),
        inputs: () => request('/inputs'),
        addInput: (data) => request('/inputs/add', { method: 'POST', body: JSON.stringify(data) }),
        removeInput: (name) => request(`/inputs/remove?name=${encodeURIComponent(name)}`),
        addOutput: (data) => request('/outputs/add', { method: 'POST', body: JSON.stringify(data) }),
        removeOutput: (data) => request('/outputs/remove', { method: 'POST', body: JSON.stringify(data) }),
        reconnectOutput: (data) => request('/outputs/reconnect', { method: 'POST', body: JSON.stringify(data) }),
        settings: () => request('/settings'),
        saveSettings: (data) => request('/settings', { method: 'PUT', body: JSON.stringify(data) }),
        reloadSettings: () => request('/settings/reload', { method: 'POST' }),
    };
})();

/* -------------------- Formatters -------------------- */
const Formatters = (() => {
    const escapeHtml = (s) => String(s)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');

    const bitrate = (kbps) => {
        if (kbps == null || kbps <= 0) return '—';
        if (kbps >= 1000) return (kbps / 1000).toFixed(1) + ' Mbps';
        return kbps.toFixed(0) + ' kbps';
    };

    const uptime = (s) => {
        if (s == null) return '—';
        if (typeof s === 'string') {
            const trimmed = s.trim();
            if (!trimmed || trimmed === '0' || trimmed === '00:00:00' || trimmed === '0s') return '—';
            return trimmed;
        }
        if (typeof s !== 'number' || s < 0) return '—';
        s = Math.floor(s);
        const h = Math.floor(s / 3600);
        const m = Math.floor((s % 3600) / 60);
        const sec = s % 60;
        const pad = (n) => String(n).padStart(2, '0');
        if (h > 0) return `${h}:${pad(m)}:${pad(sec)}`;
        return `${pad(m)}:${pad(sec)}`;
    };

    const relativeTime = (timestamp) => {
        if (!timestamp) return '—';
        const diff = Math.max(0, Math.floor((Date.now() - timestamp) / 1000));
        if (diff < 5) return I18N.t('common.justNow');
        if (diff < 60) return diff + I18N.t('common.secondsAgo');
        const min = Math.floor(diff / 60);
        return min + I18N.t('common.minutesAgo');
    };

    const detectInputProto = (input) => {
        const lower = (input.name || '').toLowerCase();
        const path = input.url_path || '';
        if (lower.includes('whip') || path.includes('/whip')) return 'WHIP';
        if (lower.includes('srt') || path.includes('srt:')) return 'SRT';
        return 'RTMP';
    };

    const detectOutputType = (url) => {
        const u = (url || '').toLowerCase();
        if (u.includes('youtube.com') || u.includes('rtmp.youtube')) return { name: 'YT', cls: 'yt' };
        if (u.includes('twitch.tv')) return { name: 'TW', cls: 'tw' };
        if (u.startsWith('srt://')) return { name: 'SRT', cls: 'srt' };
        if (u.startsWith('file://')) return { name: 'FILE', cls: 'file' };
        if (u.startsWith('rtmp://')) return { name: 'RTMP', cls: 'rtmp' };
        return { name: 'URL', cls: 'rtmp' };
    };

    return { escapeHtml, bitrate, uptime, relativeTime, detectInputProto, detectOutputType };
})();

/* -------------------- Toast -------------------- */
const Toast = (() => {
    const container = () => document.getElementById('toast-container');

    const show = (message, type = 'success') => {
        const el = document.createElement('div');
        el.className = `toast ${type}`;
        el.textContent = message;
        container().appendChild(el);
        setTimeout(() => {
            el.style.opacity = '0';
            el.style.transition = 'opacity 0.2s';
            setTimeout(() => el.remove(), 200);
        }, 4000);
    };

    return {
        success: (m) => show(m, 'success'),
        error: (m) => show(m, 'error'),
    };
})();

/* -------------------- Modal -------------------- */
const Modal = (() => {
    let backdrop = null;
    let onCloseCb = null;
    const escHandler = (e) => {
        if (e.key === 'Escape') close();
    };

    const open = (html) => {
        close();
        backdrop = document.createElement('div');
        backdrop.className = 'modal-backdrop';
        backdrop.innerHTML = `<div class="modal">${html}</div>`;
        document.body.appendChild(backdrop);
        backdrop.addEventListener('click', (e) => {
            if (e.target === backdrop) close();
        });
        document.addEventListener('keydown', escHandler);
    };

    const close = () => {
        if (backdrop) {
            backdrop.remove();
            backdrop = null;
        }
        document.removeEventListener('keydown', escHandler);
        if (onCloseCb) {
            const cb = onCloseCb;
            onCloseCb = null;
            cb();
        }
    };

    const confirm = ({ title, body, confirmText, cancelText, danger }) => {
        return new Promise((resolve) => {
            open(`
                <div class="modal-header"><h3>${Formatters.escapeHtml(title)}</h3></div>
                <div class="modal-body">${body || ''}</div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" data-modal-cancel>${Formatters.escapeHtml(cancelText || I18N.t('common.cancel'))}</button>
                    <button class="btn ${danger ? 'btn-primary' : 'btn-primary'}" data-modal-confirm style="${danger ? 'background: var(--err);' : ''}">${Formatters.escapeHtml(confirmText || I18N.t('common.confirmRemove'))}</button>
                </div>
            `);
            const modal = backdrop.querySelector('.modal');
            modal.querySelector('[data-modal-cancel]').addEventListener('click', () => { close(); resolve(false); });
            modal.querySelector('[data-modal-confirm]').addEventListener('click', () => { close(); resolve(true); });
        });
    };

    const custom = ({ header, body, footer, onClose }) => {
        onCloseCb = onClose || null;
        open(`
            ${header ? `<div class="modal-header"><h3>${header}</h3></div>` : ''}
            <div class="modal-body" style="${footer ? '' : 'padding-bottom: var(--space-5);'}">${body}</div>
            ${footer ? `<div class="modal-footer">${footer}</div>` : ''}
        `);
    };

    return { open, close, confirm, custom };
})();

/* -------------------- TopBar -------------------- */
const TopBar = (() => {
    let lastUpdate = null;
    let serverOnline = true;
    let summaryData = { total: 0, active: 0, outTotal: 0, outActive: 0, conn: 0 };
    let relTimer = null;

    const led = () => document.getElementById('server-led');
    const text = () => document.getElementById('server-status-text');
    const rel = () => document.getElementById('last-update-rel');
    const summary = () => document.getElementById('topbar-summary');

    const setOnline = (online) => {
        serverOnline = online;
        led().classList.toggle('offline', !online);
        text().textContent = I18N.t(online ? 'common.online' : 'common.offline');
    };

    const setLastUpdate = (ts) => {
        lastUpdate = ts;
        updateRel();
    };

    const updateRel = () => {
        rel().textContent = lastUpdate ? Formatters.relativeTime(lastUpdate) : '—';
    };

    const setSummary = (data) => {
        summaryData = data;
        summary().textContent = I18N.t('topbar.summary', {
            active: data.active,
            total: data.total,
            outActive: data.outActive,
            outTotal: data.outTotal,
            conn: data.conn,
        });
    };

    const setDirty = (dirty) => {
        document.getElementById('topbar-dirty').classList.toggle('show', dirty);
    };

    const startTicker = () => {
        if (relTimer) return;
        relTimer = setInterval(updateRel, 1000);
    };

    return { setOnline, setLastUpdate, setSummary, setDirty, startTicker };
})();

/* -------------------- StreamCard (Dashboard render) -------------------- */
const StreamCard = (() => {
    const render = (input) => {
        const proto = Formatters.detectInputProto(input);
        const protoLower = proto.toLowerCase();
        const outputs = input.outputs || [];

        const activeOutputs = outputs.filter(o => o.active).length;
        const totalBitrate = outputs.reduce((acc, o) => acc + (o.active ? (o.bitrate_kbps || 0) : 0), 0);

        const metaItems = [
            { label: I18N.t('common.connections'), value: input.connections || 0, mono: true },
            { label: I18N.t('common.uptime'), value: input.uptime ? Formatters.uptime(input.uptime) : '—', mono: true, muted: !input.uptime },
            { label: I18N.t('status.totalBitrate'), value: Formatters.bitrate(totalBitrate), mono: true },
        ];

        if (input.last_error) {
            metaItems.push({ label: I18N.t('status.lastError'), value: Formatters.escapeHtml(input.last_error), mono: false, danger: true });
        }

        return `
        <article class="stream-card" data-stream="${Formatters.escapeHtml(input.name)}">
            <div class="stream-card-meta">
                <div class="stack">
                    <div class="stream-card-header" style="justify-content: space-between; flex-wrap: nowrap; align-items: center; width: 100%;">
                        <div class="row" style="gap: var(--space-2); min-width: 0;">
                            <span class="led ${input.active ? '' : 'off'}" data-cell="led"></span>
                            <span class="name" data-cell="name" style="text-overflow: ellipsis; overflow: hidden; white-space: nowrap;">${Formatters.escapeHtml(input.name)}</span>
                            <span class="proto-badge ${protoLower}" data-cell="proto">${proto}</span>
                        </div>
                        <button class="btn-icon danger" data-action="remove-input" data-name="${Formatters.escapeHtml(input.name)}" title="${I18N.t('common.remove')}" style="flex-shrink: 0; margin-left: auto; padding: 4px;">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                        </button>
                    </div>
                    <div class="stream-card-header">
                        <span class="path" data-cell="path">${Formatters.escapeHtml(input.url_path || '')}</span>
                    </div>
                    <div class="stream-card-header">
                        <span class="status-line ${input.active ? 'active' : 'inactive'}" data-cell="status">
                            ${input.active ? I18N.t('common.active') : I18N.t('common.inactive')}
                        </span>
                    </div>
                    <dl class="meta-grid">
                        ${metaItems.map(item => `
                            <dt>${item.label}</dt>
                            <dd class="${item.muted ? 'muted' : ''} ${item.danger ? 'danger' : ''}" data-cell-meta="${item.label}">${item.value}</dd>
                        `).join('')}
                    </dl>
                </div>
            </div>
            <div class="stream-card-body">
                <div class="outputs-section-header">
                    <span>${I18N.t('common.outputs')}</span>
                    <span class="count" data-cell="out-count">${activeOutputs}/${outputs.length}</span>
                </div>
                ${renderOutputsTable(input.name, outputs)}
                <form class="inline-add-form" data-action="add-output" data-name="${Formatters.escapeHtml(input.name)}">
                    <input type="text" class="form-input mono" placeholder="${Formatters.escapeHtml(I18N.t('outputs.urlPlaceholder'))}" required>
                    <button type="submit" class="btn btn-primary">${I18N.t('common.add')}</button>
                </form>
            </div>
        </article>
        `;
    };

    const renderOutputsTable = (inputName, outputs) => {
        if (!outputs || outputs.length === 0) {
            return `<div class="empty-state">${I18N.t('status.noOutputs')}</div>`;
        }
        return `
        <div class="outputs-table" data-outputs-of="${Formatters.escapeHtml(inputName)}">
            ${outputs.map(o => renderOutputRow(inputName, o)).join('')}
        </div>
        `;
    };

    const renderOutputRow = (inputName, o) => {
        const type = Formatters.detectOutputType(o.url);
        const bitrateVal = o.active && o.bitrate_kbps > 0 ? Formatters.bitrate(o.bitrate_kbps) : '—';
        const bitrateCls = o.active && o.bitrate_kbps > 0 ? '' : 'muted';
        const uptimeVal = o.uptime ? Formatters.uptime(o.uptime) : '—';
        const uptimeCls = o.active ? '' : 'muted';
        const errVal = o.error_count || 0;
        const errCls = errVal > 0 ? 'has-errors' : '';

        return `
        <div class="outputs-row ${o.active ? '' : 'inactive'}" data-output-url="${Formatters.escapeHtml(o.url)}">
            <span class="led cell ${o.active ? '' : 'off'}" data-cell="led"></span>
            <span class="cell" data-cell="type"><span class="proto-badge ${type.cls}">${type.name}</span></span>
            <span class="url cell" data-cell="url" title="${Formatters.escapeHtml(o.url)}" data-action="copy-url">${Formatters.escapeHtml(o.url)}</span>
            <span class="bitrate cell ${bitrateCls}" data-cell="bitrate">${bitrateVal}</span>
            <span class="uptime cell ${uptimeCls}" data-cell="uptime">${uptimeVal}</span>
            <span class="errors cell ${errCls}" data-cell="errors">${errVal}</span>
            <span class="actions cell">
                <button class="btn-icon" data-action="reconnect-output" data-name="${Formatters.escapeHtml(inputName)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('common.reconnect')}">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
                </button>
                <button class="btn-icon" data-action="edit-output" data-name="${Formatters.escapeHtml(inputName)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('outputs.editTitle')}">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                </button>
                <button class="btn-icon danger" data-action="remove-output" data-name="${Formatters.escapeHtml(inputName)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('common.remove')}">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                </button>
            </span>
        </div>
        `;
    };

    const renderOutputRowKey = (inputName, o) => inputName + '|' + o.url;

    return { render, renderOutputRow, renderOutputsTable, renderOutputRowKey };
})();

/* -------------------- Dashboard -------------------- */
const Dashboard = (() => {
    const container = () => document.getElementById('status-content');
    let lastStatuses = [];

    const emptyState = () => `
        <div class="card">
            <div class="card-body" style="text-align: center; padding: var(--space-6);">
                <div style="font-size: 0.9rem; color: var(--text-muted);">${I18N.t('status.noInputs')}</div>
            </div>
        </div>
    `;

    const loadingSkeleton = () => `
        <div class="stack-lg">
            <div class="skeleton skeleton-card"></div>
            <div class="skeleton skeleton-card"></div>
        </div>
    `;

    const renderAll = (statuses) => {
        if (!statuses || statuses.length === 0) {
            container().innerHTML = emptyState();
            lastStatuses = [];
            return;
        }
        const html = `<div class="stack-lg">${statuses.map(s => StreamCard.render(s)).join('')}</div>`;
        container().innerHTML = html;
        lastStatuses = statuses;
    };

    const diff = (newStatuses) => {
        if (!newStatuses || newStatuses.length === 0) {
            container().innerHTML = emptyState();
            lastStatuses = [];
            return;
        }

        if (lastStatuses.length === 0) {
            renderAll(newStatuses);
            return;
        }

        const oldMap = new Map(lastStatuses.map(s => [s.name, s]));
        const newMap = new Map(newStatuses.map(s => [s.name, s]));

        // Removed streams
        for (const [name] of oldMap) {
            if (!newMap.has(name)) {
                const el = container().querySelector(`[data-stream="${cssEscape(name)}"]`);
                if (el) el.remove();
            }
        }

        // Added streams
        for (const [name, s] of newMap) {
            if (!oldMap.has(name)) {
                const wrapper = document.createElement('div');
                wrapper.innerHTML = StreamCard.render(s);
                const card = wrapper.firstElementChild;
                const stack = container().querySelector('.stack-lg');
                if (stack) {
                    // insert in correct order
                    const idx = newStatuses.findIndex(x => x.name === name);
                    const next = newStatuses.slice(idx + 1).find(x => newMap.has(x.name) && !oldMap.has(x.name));
                    if (next) {
                        const ref = container().querySelector(`[data-stream="${cssEscape(next.name)}"]`);
                        if (ref) stack.insertBefore(card, ref);
                        else stack.appendChild(card);
                    } else {
                        stack.appendChild(card);
                    }
                }
            }
        }

        // Update existing
        for (const s of newStatuses) {
            const old = oldMap.get(s.name);
            if (!old) continue;
            updateStream(s, old);
        }

        lastStatuses = newStatuses;
    };

    const updateStream = (s, old) => {
        const card = container().querySelector(`[data-stream="${cssEscape(s.name)}"]`);
        if (!card) return;

        // LED
        const led = card.querySelector('[data-cell="led"]');
        if (led) led.classList.toggle('off', !s.active);

        // Status line
        const status = card.querySelector('[data-cell="status"]');
        if (status) {
            status.classList.toggle('active', s.active);
            status.classList.toggle('inactive', !s.active);
            status.textContent = I18N.t(s.active ? 'common.active' : 'common.inactive');
        }

        // Total bitrate
        const activeOutputs = (s.outputs || []).filter(o => o.active);
        const total = activeOutputs.reduce((acc, o) => acc + (o.bitrate_kbps || 0), 0);
        const oldActive = (old.outputs || []).filter(o => o.active);
        const oldTotal = oldActive.reduce((acc, o) => acc + (o.bitrate_kbps || 0), 0);
        if (total !== oldTotal) {
            const dd = [...card.querySelectorAll('[data-cell-meta]')].find(el => el.getAttribute('data-cell-meta') === I18N.t('status.totalBitrate'));
            if (dd) dd.textContent = Formatters.bitrate(total);
        }

        // Connections
        if ((s.connections || 0) !== (old.connections || 0)) {
            const dd = [...card.querySelectorAll('[data-cell-meta]')].find(el => el.getAttribute('data-cell-meta') === I18N.t('common.connections'));
            if (dd) dd.textContent = s.connections || 0;
        }

        // Uptime
        if ((s.uptime || 0) !== (old.uptime || 0)) {
            const dd = [...card.querySelectorAll('[data-cell-meta]')].find(el => el.getAttribute('data-cell-meta') === I18N.t('common.uptime'));
            if (dd) {
                dd.classList.toggle('muted', !s.uptime);
                dd.textContent = s.uptime ? Formatters.uptime(s.uptime) : '—';
            }
        }

        // Output count
        const outCount = card.querySelector('[data-cell="out-count"]');
        if (outCount) {
            const act = (s.outputs || []).filter(o => o.active).length;
            const tot = (s.outputs || []).length;
            outCount.textContent = `${act}/${tot}`;
        }

        // Diff outputs
        diffOutputs(s, old);
    };

    const diffOutputs = (s, old) => {
        const table = container().querySelector(`[data-outputs-of="${cssEscape(s.name)}"]`);
        const oldOuts = old.outputs || [];
        const newOuts = s.outputs || [];
        const oldMap = new Map(oldOuts.map(o => [o.url, o]));
        const newMap = new Map(newOuts.map(o => [o.url, o]));

        if (table) {
            // Removed outputs
            for (const [url] of oldMap) {
                if (!newMap.has(url)) {
                    const row = table.querySelector(`[data-output-url="${cssEscape(url)}"]`);
                    if (row) row.remove();
                }
            }

            // If empty after removals, replace with empty state
            if (newOuts.length === 0) {
                const emptyHtml = `<div class="empty-state">${I18N.t('status.noOutputs')}</div>`;
                table.outerHTML = emptyHtml;
                return;
            }

            // Added outputs
            for (const [url, o] of newMap) {
                if (!oldMap.has(url)) {
                    const wrapper = document.createElement('div');
                    wrapper.innerHTML = StreamCard.renderOutputRow(s.name, o);
                    const row = wrapper.firstElementChild;
                    // insert in correct order
                    const idx = newOuts.findIndex(x => x.url === url);
                    const nextAfter = newOuts.slice(idx + 1).find(x => newMap.has(x.url) && !oldMap.has(x.url));
                    if (nextAfter) {
                        const ref = table.querySelector(`[data-output-url="${cssEscape(nextAfter.url)}"]`);
                        if (ref) table.insertBefore(row, ref);
                        else table.appendChild(row);
                    } else {
                        table.appendChild(row);
                    }
                }
            }

            // Update existing output rows
            for (const o of newOuts) {
                const oldO = oldMap.get(o.url);
                if (!oldO) continue;
                updateOutputRow(table, s.name, o, oldO);
            }
        } else {
            // No table — need to create
            if (newOuts.length > 0) {
                const card = container().querySelector(`[data-stream="${cssEscape(s.name)}"]`);
                if (!card) return;
                const addForm = card.querySelector('form.inline-add-form');
                // Remove any existing empty-state sibling
                card.querySelectorAll('.empty-state').forEach(el => el.remove());
                const wrapper = document.createElement('div');
                wrapper.innerHTML = StreamCard.renderOutputsTable(s.name, newOuts);
                const newTable = wrapper.firstElementChild;
                if (addForm) {
                    card.querySelector('.stream-card-body').insertBefore(newTable, addForm);
                } else {
                    card.querySelector('.stream-card-body').appendChild(newTable);
                }
            }
        }
    };

    const updateOutputRow = (table, inputName, o, oldO) => {
        const row = table.querySelector(`[data-output-url="${cssEscape(o.url)}"]`);
        if (!row) return;

        // LED
        const led = row.querySelector('[data-cell="led"]');
        if (led) led.classList.toggle('off', !o.active);

        // Row active class
        row.classList.toggle('inactive', !o.active);

        // Bitrate
        if ((o.bitrate_kbps || 0) !== (oldO.bitrate_kbps || 0) || o.active !== oldO.active) {
            const cell = row.querySelector('[data-cell="bitrate"]');
            if (cell) {
                const v = o.active && o.bitrate_kbps > 0 ? Formatters.bitrate(o.bitrate_kbps) : '—';
                cell.textContent = v;
                cell.classList.toggle('muted', !(o.active && o.bitrate_kbps > 0));
            }
        }

        // Uptime
        if ((o.uptime || 0) !== (oldO.uptime || 0)) {
            const cell = row.querySelector('[data-cell="uptime"]');
            if (cell) {
                cell.textContent = o.uptime ? Formatters.uptime(o.uptime) : '—';
                cell.classList.toggle('muted', !o.active);
            }
        }

        // Errors
        if ((o.error_count || 0) !== (oldO.error_count || 0)) {
            const cell = row.querySelector('[data-cell="errors"]');
            if (cell) {
                const v = o.error_count || 0;
                cell.textContent = v;
                cell.classList.toggle('has-errors', v > 0);
            }
        }
    };

    const cssEscape = (s) => Formatters.escapeHtml(s);

    const load = async () => {
        if (lastStatuses.length === 0) {
            container().innerHTML = loadingSkeleton();
        }
        try {
            const statuses = await Api.statusAll();
            TopBar.setOnline(true);
            TopBar.setLastUpdate(Date.now());

            // Update summary
            const total = statuses.length;
            const active = statuses.filter(s => s.active).length;
            let outTotal = 0, outActive = 0, conn = 0;
            statuses.forEach(s => {
                conn += s.connections || 0;
                if (s.outputs) {
                    s.outputs.forEach(o => {
                        outTotal++;
                        if (o.active) outActive++;
                    });
                }
            });
            TopBar.setSummary({ total, active, outTotal, outActive, conn });

            diff(statuses);
        } catch (e) {
            TopBar.setOnline(false);
            container().innerHTML = `<div class="card"><div class="card-body" style="color: var(--err);">${I18N.t('status.loadError')}: ${Formatters.escapeHtml(e.message)}</div></div>`;
        }
    };

    // Action delegation
    const bindActions = () => {
        container().addEventListener('click', async (e) => {
            const target = e.target.closest('[data-action]');
            if (!target) return;
            const action = target.getAttribute('data-action');

            if (action === 'remove-input') {
                const name = target.getAttribute('data-name');
                const ok = await Modal.confirm({
                    title: I18N.t('inputs.removeConfirm') + ' "' + name + '"?',
                    body: I18N.t('inputs.confirmRemoveBody'),
                    confirmText: I18N.t('common.remove'),
                    danger: true,
                });
                if (!ok) return;
                try {
                    await Api.removeInput(name);
                    Toast.success(I18N.t('inputs.removeSuccess'));
                    load();
                } catch (err) {
                    Toast.error(I18N.t('inputs.removeError') + ': ' + err.message);
                }
            }

            if (action === 'reconnect-output') {
                const name = target.getAttribute('data-name');
                const url = target.getAttribute('data-url');
                try {
                    await Api.reconnectOutput({ name, url });
                    Toast.success(I18N.t('outputs.reconnectInitiated'));
                } catch (err) {
                    Toast.error(I18N.t('outputs.reconnectError') + ': ' + err.message);
                }
            }

            if (action === 'remove-output') {
                const name = target.getAttribute('data-name');
                const url = target.getAttribute('data-url');
                const ok = await Modal.confirm({
                    title: I18N.t('outputs.removeConfirm') + ' "' + url + '"?',
                    confirmText: I18N.t('common.remove'),
                    danger: true,
                });
                if (!ok) return;
                try {
                    await Api.removeOutput({ name, url });
                    Toast.success(I18N.t('outputs.removeSuccess'));
                    load();
                } catch (err) {
                    Toast.error(I18N.t('outputs.removeError') + ': ' + err.message);
                }
            }

            if (action === 'edit-output') {
                const name = target.getAttribute('data-name');
                const url = target.getAttribute('data-url');
                showEditOutputModal(name, url, async (newUrl) => {
                    if (newUrl === url) return;
                    if (!newUrl.startsWith('rtmp://') && !newUrl.startsWith('srt://') && !newUrl.startsWith('file://')) {
                        throw new Error('URL must start with srt://, rtmp:// or file://');
                    }
                    try {
                        await Api.removeOutput({ name, url });
                        await Api.addOutput({ name, url: newUrl });
                        Toast.success(I18N.t('outputs.editSuccess'));
                        load();
                    } catch (err) {
                        throw new Error(I18N.t('outputs.editError') + ': ' + err.message);
                    }
                });
            }

            if (action === 'copy-url') {
                const url = target.getAttribute('title') || target.textContent;
                try {
                    await navigator.clipboard.writeText(url);
                    Toast.success(I18N.t('common.copied'));
                } catch {
                    /* clipboard may not be available in WebView */
                }
            }
        });

        container().addEventListener('submit', async (e) => {
            const form = e.target.closest('form[data-action="add-output"]');
            if (!form) return;
            e.preventDefault();
            const name = form.getAttribute('data-name');
            const input = form.querySelector('input');
            const url = input.value.trim();
            if (!url) return;
            const btn = form.querySelector('button[type="submit"]');
            btn.disabled = true;
            try {
                await Api.addOutput({ name, url });
                Toast.success(I18N.t('outputs.addSuccess'));
                input.value = '';
                load();
            } catch (err) {
                Toast.error(I18N.t('outputs.addError') + ': ' + err.message);
            } finally {
                btn.disabled = false;
            }
        });
    };

    return { load, bindActions, renderAll };
})();

/* -------------------- InputsPage -------------------- */
const InputsPage = (() => {
    let inputs = [];
    let selected = null;

    const container = () => document.getElementById('inputs-content');

    const renderForm = () => `
        <div class="card">
            <div class="card-header"><h3>${I18N.t('inputs.addNew')}</h3></div>
            <div class="card-body">
                <form id="add-input-form">
                    <div class="form-row">
                        <div class="form-group">
                            <label>${I18N.t('inputs.inputName')}</label>
                            <input type="text" class="form-input mono" id="input-name" placeholder="obs" required>
                        </div>
                        <div class="form-group">
                            <label>${I18N.t('inputs.rtmpPath')}</label>
                            <input type="text" class="form-input mono" id="input-path" placeholder="/live/stream" required>
                        </div>
                    </div>
                    <div class="form-group">
                        <label>${I18N.t('inputs.protocol')}</label>
                        <div class="segment" id="input-proto">
                            <button type="button" class="active" data-proto="RTMP">RTMP</button>
                            <button type="button" data-proto="SRT">SRT</button>
                            <button type="button" data-proto="WHIP">WHIP</button>
                        </div>
                    </div>
                    <div class="form-group">
                        <label>${I18N.t('inputs.outputsLabel')}</label>
                        <textarea class="form-textarea mono" id="input-outputs" placeholder="${Formatters.escapeHtml(I18N.t('inputs.outputsPlaceholder'))}"></textarea>
                    </div>
                    <div class="row" style="justify-content: flex-end;">
                        <button type="submit" class="btn btn-primary">${I18N.t('common.add')}</button>
                    </div>
                </form>
            </div>
        </div>
    `;

    const renderList = () => {
        if (!inputs || inputs.length === 0) {
            return `<div class="empty-state">${I18N.t('inputs.noInputs')}</div>`;
        }
        return inputs.map(inp => {
            const name = inp.Name || inp.name;
            const path = inp.URLPath || inp.url_path;
            const outCount = (inp.Outputs || inp.outputs || []).length;
            const proto = Formatters.detectInputProto({ name, url_path: path });
            return `
                <div class="master-list-item ${selected === name ? 'active' : ''}" data-select-input="${Formatters.escapeHtml(name)}">
                    <span class="led"></span>
                    <div style="flex: 1; min-width: 0;">
                        <div class="name">${Formatters.escapeHtml(name)}</div>
                        <div class="path">${Formatters.escapeHtml(path)}</div>
                    </div>
                    <span class="proto-badge ${proto.toLowerCase()}">${proto}</span>
                    <span class="count" style="font-family: var(--font-mono); font-size: 0.72rem; color: var(--text-muted);">${outCount}</span>
                </div>
            `;
        }).join('');
    };

    const renderDetail = () => {
        if (!selected) {
            return `<div class="card"><div class="card-body" style="text-align: center; color: var(--text-subtle); padding: var(--space-6);">${I18N.t('inputs.selectHint')}</div></div>`;
        }
        const inp = inputs.find(i => (i.Name || i.name) === selected);
        if (!inp) return '';
        const name = inp.Name || inp.name;
        const path = inp.URLPath || inp.url_path;
        const outputs = inp.Outputs || inp.outputs || [];
        const proto = Formatters.detectInputProto({ name, url_path: path });

        return `
        <div class="card">
            <div class="card-header">
                <div class="row">
                    <span class="led"></span>
                    <h3>${Formatters.escapeHtml(name)}</h3>
                    <span class="proto-badge ${proto.toLowerCase()}">${proto}</span>
                </div>
                <button class="btn btn-danger" data-action="remove-input" data-name="${Formatters.escapeHtml(name)}" style="padding: 5px 10px; font-size: 0.75rem;">${I18N.t('common.remove')}</button>
            </div>
            <div class="card-body">
                <dl class="meta-grid">
                    <dt>${I18N.t('common.path')}</dt>
                    <dd>${Formatters.escapeHtml(path)}</dd>
                    <dt>${I18N.t('common.outputs')}</dt>
                    <dd>${outputs.length}</dd>
                </dl>
                <div class="outputs-section-header" style="margin-top: var(--space-4); margin-bottom: var(--space-2);">
                    <span>${I18N.t('inputs.outputs')}</span>
                </div>
                ${outputs.length > 0 ? `
                    <div class="outputs-table">
                        ${outputs.map(o => `
                            <div class="outputs-row ${o.active ? '' : 'inactive'}">
                                <span class="led cell ${o.active ? '' : 'off'}"></span>
                                <span class="cell" data-cell="type"><span class="proto-badge ${Formatters.detectOutputType(o.url).cls}">${Formatters.detectOutputType(o.url).name}</span></span>
                                <span class="url cell" title="${Formatters.escapeHtml(o.url)}">${Formatters.escapeHtml(o.url)}</span>
                                <span class="bitrate cell ${o.active && o.bitrate_kbps > 0 ? '' : 'muted'}">${o.active && o.bitrate_kbps > 0 ? Formatters.bitrate(o.bitrate_kbps) : '—'}</span>
                                <span class="uptime cell ${o.active ? '' : 'muted'}">${o.uptime ? Formatters.uptime(o.uptime) : '—'}</span>
                                <span class="errors cell ${o.error_count > 0 ? 'has-errors' : ''}">${o.error_count || 0}</span>
                                <span class="actions cell">
                                    <button class="btn-icon" data-action="reconnect-output" data-name="${Formatters.escapeHtml(name)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('common.reconnect')}">
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>
                                    </button>
                                    <button class="btn-icon" data-action="edit-output" data-name="${Formatters.escapeHtml(name)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('outputs.editTitle')}">
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                                    </button>
                                    <button class="btn-icon danger" data-action="remove-output" data-name="${Formatters.escapeHtml(name)}" data-url="${Formatters.escapeHtml(o.url)}" title="${I18N.t('common.remove')}">
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                                    </button>
                                </span>
                            </div>
                        `).join('')}
                    </div>
                ` : `<div class="empty-state">${I18N.t('status.noOutputs')}</div>`}
            </div>
        </div>
        `;
    };

    const render = () => {
        container().innerHTML = `
            <div class="stack-lg">
                ${renderForm()}
                <div class="master-detail">
                    <div class="master-list">
                        <div style="padding: var(--space-3) var(--space-4); border-bottom: 1px solid var(--border); font-size: 0.75rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em;">${I18N.t('inputs.existing')}</div>
                        ${renderList()}
                    </div>
                    ${renderDetail()}
                </div>
            </div>
        `;
        bind();
    };

    const bind = () => {
        // Form submit
        const form = document.getElementById('add-input-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const name = document.getElementById('input-name').value.trim();
                const path = document.getElementById('input-path').value.trim();
                const proto = document.querySelector('#input-proto .active')?.getAttribute('data-proto') || 'RTMP';
                const outputsText = document.getElementById('input-outputs').value;
                const outputs = outputsText.split('\n').map(s => s.trim()).filter(Boolean);
                // For SRT/WHIP, prepend protocol prefix to path
                let urlPath = path;
                if (proto === 'WHIP' && !urlPath.startsWith('/whip/')) {
                    urlPath = '/whip' + (urlPath.startsWith('/') ? urlPath : '/' + urlPath);
                }
                try {
                    await Api.addInput({ name, url_path: urlPath, outputs });
                    Toast.success(I18N.t('inputs.addSuccess'));
                    form.reset();
                    selected = name;
                    load();
                } catch (err) {
                    Toast.error(I18N.t('inputs.addError') + ': ' + err.message);
                }
            });
        }

        // Protocol segment
        const seg = document.getElementById('input-proto');
        if (seg) {
            seg.querySelectorAll('button').forEach(btn => {
                btn.addEventListener('click', () => {
                    seg.querySelectorAll('button').forEach(b => b.classList.remove('active'));
                    btn.classList.add('active');
                });
            });
        }

        // List select
        container().querySelectorAll('[data-select-input]').forEach(el => {
            el.addEventListener('click', () => {
                selected = el.getAttribute('data-select-input');
                render();
            });
        });

        // Action delegation
        container().addEventListener('click', async (e) => {
            const t = e.target.closest('[data-action]');
            if (!t) return;
            const action = t.getAttribute('data-action');
            if (action === 'remove-input') {
                const name = t.getAttribute('data-name');
                const ok = await Modal.confirm({
                    title: I18N.t('inputs.removeConfirm') + ' "' + name + '"?',
                    body: I18N.t('inputs.confirmRemoveBody'),
                    confirmText: I18N.t('common.remove'),
                    danger: true,
                });
                if (!ok) return;
                try {
                    await Api.removeInput(name);
                    Toast.success(I18N.t('inputs.removeSuccess'));
                    if (selected === name) selected = null;
                    load();
                } catch (err) {
                    Toast.error(I18N.t('inputs.removeError') + ': ' + err.message);
                }
            }
            if (action === 'remove-output' || action === 'reconnect-output') {
                const name = t.getAttribute('data-name');
                const url = t.getAttribute('data-url');
                if (action === 'reconnect-output') {
                    try {
                        await Api.reconnectOutput({ name, url });
                        Toast.success(I18N.t('outputs.reconnectInitiated'));
                    } catch (err) {
                        Toast.error(I18N.t('outputs.reconnectError') + ': ' + err.message);
                    }
                } else {
                    const ok = await Modal.confirm({
                        title: I18N.t('outputs.removeConfirm') + ' "' + url + '"?',
                        confirmText: I18N.t('common.remove'),
                        danger: true,
                    });
                    if (!ok) return;
                    try {
                        await Api.removeOutput({ name, url });
                        Toast.success(I18N.t('outputs.removeSuccess'));
                        load();
                    } catch (err) {
                        Toast.error(I18N.t('outputs.removeError') + ': ' + err.message);
                    }
                }
            }

            if (action === 'edit-output') {
                const name = t.getAttribute('data-name');
                const url = t.getAttribute('data-url');
                showEditOutputModal(name, url, async (newUrl) => {
                    if (newUrl === url) return;
                    if (!newUrl.startsWith('rtmp://') && !newUrl.startsWith('srt://') && !newUrl.startsWith('file://')) {
                        throw new Error('URL must start with srt://, rtmp:// or file://');
                    }
                    try {
                        await Api.removeOutput({ name, url });
                        await Api.addOutput({ name, url: newUrl });
                        Toast.success(I18N.t('outputs.editSuccess'));
                        load();
                    } catch (err) {
                        throw new Error(I18N.t('outputs.editError') + ': ' + err.message);
                    }
                });
            }
        });
    };

    const load = async () => {
        container().innerHTML = `<div class="loading">${I18N.t('common.loading')}</div>`;
        try {
            // Запрашиваем конфигурацию входов и их статусы параллельно
            const [inputsData, statusData] = await Promise.all([
                Api.inputs(),
                Api.statusAll()
            ]);

            // Объединяем информацию о статусах выходов в объект входа
            inputs = inputsData.map(inp => {
                const name = inp.Name || inp.name;
                const status = statusData.find(s => s.name === name);
                const outputs = status ? (status.outputs || []) : [];
                return {
                    ...inp,
                    Outputs: outputs,
                    outputs: outputs
                };
            });

            if (!selected && inputs.length > 0) {
                selected = (inputs[0].Name || inputs[0].name);
            }
            render();
        } catch (err) {
            container().innerHTML = `<div class="card"><div class="card-body" style="color: var(--err);">${I18N.t('inputs.loadError')}: ${Formatters.escapeHtml(err.message)}</div></div>`;
        }
    };

    return { load };
})();

/* -------------------- SettingsPage -------------------- */
const SettingsPage = (() => {
    let settings = null;
    let original = null;
    let section = 'srt';
    let dirty = false;
    const container = () => document.getElementById('settings-content');

    const sections = [
        { id: 'srt', label: 'settings.section.srt' },
        { id: 'logging', label: 'settings.section.logging' },
        { id: 'reconnect', label: 'settings.section.reconnect' },
    ];

    const setDirty = (val) => {
        dirty = val;
        TopBar.setDirty(val);
        updateActionsState();
    };

    const updateActionsState = () => {
        const revert = document.getElementById('settings-revert');
        if (revert) revert.disabled = !dirty;
    };

    const renderSrt = () => `
        <div class="card">
            <div class="card-header"><h3>${I18N.t('settings.srt')}</h3></div>
            <div class="card-body">
                <div class="form-row">
                    <div class="form-group">
                        <label>${I18N.t('settings.latency')}</label>
                        <input type="number" class="form-input mono" id="srt-latency" value="${settings.srt_settings?.latency ?? 120}" min="0" data-dirty>
                    </div>
                    <div class="form-group">
                        <label>${I18N.t('settings.timeout')}</label>
                        <input type="number" class="form-input mono" id="srt-timeout" value="${settings.srt_settings?.connect_timeout ?? 5000}" min="0" data-dirty>
                    </div>
                </div>
                <div class="form-group">
                    <label>${I18N.t('settings.passphrase')}</label>
                    <input type="text" class="form-input mono" id="srt-passphrase" value="${Formatters.escapeHtml(settings.srt_settings?.passphrase || '')}" data-dirty>
                    <div class="help">${I18N.t('settings.passphraseHelp')}</div>
                </div>
                <div class="form-group">
                    <label>${I18N.t('settings.streamId')}</label>
                    <input type="text" class="form-input mono" id="srt-streamid" value="${Formatters.escapeHtml(settings.srt_settings?.streamid || '')}" data-dirty>
                    <div class="help">${I18N.t('settings.streamIdHelp')}</div>
                </div>
            </div>
        </div>
    `;

    const renderLogging = () => `
        <div class="card">
            <div class="card-header"><h3>${I18N.t('settings.logging')}</h3></div>
            <div class="card-body">
                <div class="form-group">
                    <label class="form-check">
                        <input type="checkbox" id="log-to-file" ${settings.log_to_file ? 'checked' : ''} data-dirty>
                        ${I18N.t('settings.logToFile')}
                    </label>
                </div>
                <div class="form-group">
                    <label>${I18N.t('settings.logFile')}</label>
                    <input type="text" class="form-input mono" id="log-file" value="${Formatters.escapeHtml(settings.log_file || 'server.log')}" data-dirty>
                </div>
            </div>
        </div>
        <div class="card" style="margin-top: var(--space-4);">
            <div class="card-header"><h3>GUI</h3></div>
            <div class="card-body">
                <div class="form-group">
                    <label class="form-check">
                        <input type="checkbox" id="minimize-to-tray" ${settings.minimize_to_tray ? 'checked' : ''} data-dirty>
                        ${I18N.t('settings.minimizeToTray')}
                    </label>
                </div>
            </div>
        </div>
    `;

    const renderReconnect = () => `
        <div class="card">
            <div class="card-header"><h3>${I18N.t('settings.reconnectInterval')}</h3></div>
            <div class="card-body">
                <div class="form-group">
                    <label>${I18N.t('settings.reconnectIntervalLabel')}</label>
                    <input type="number" class="form-input mono" id="reconnect-interval" value="${settings.reconnect_interval ?? 5}" min="1" data-dirty>
                    <div class="help">${I18N.t('settings.reconnectIntervalHelp')}</div>
                </div>
            </div>
        </div>
    `;

    const renderSection = () => {
        if (section === 'srt') return renderSrt();
        if (section === 'logging') return renderLogging();
        if (section === 'reconnect') return renderReconnect();
        return '';
    };

    const render = () => {
        container().innerHTML = `
            <div class="settings-layout">
                <div class="settings-nav">
                    ${sections.map(s => `
                        <button data-section="${s.id}" class="${section === s.id ? 'active' : ''}">${I18N.t(s.label)}</button>
                    `).join('')}
                </div>
                <div>
                    ${renderSection()}
                    <div class="settings-actions">
                        <button class="btn btn-secondary" id="settings-reload">${I18N.t('settings.reload')}</button>
                        <button class="btn btn-secondary" id="settings-revert" disabled>${I18N.t('settings.revert')}</button>
                        <button class="btn btn-primary" id="settings-save">${I18N.t('common.save')}</button>
                    </div>
                </div>
            </div>
        `;
        bind();
    };

    const collect = () => ({
        srt_settings: {
            latency: document.getElementById('srt-latency') ? parseInt(document.getElementById('srt-latency').value) : (settings.srt_settings?.latency ?? 120),
            connect_timeout: document.getElementById('srt-timeout') ? parseInt(document.getElementById('srt-timeout').value) : (settings.srt_settings?.connect_timeout ?? 5000),
            passphrase: document.getElementById('srt-passphrase') ? document.getElementById('srt-passphrase').value : (settings.srt_settings?.passphrase || ''),
            streamid: document.getElementById('srt-streamid') ? document.getElementById('srt-streamid').value : (settings.srt_settings?.streamid || ''),
            encryption: 'none',
        },
        log_to_file: document.getElementById('log-to-file') ? document.getElementById('log-to-file').checked : (settings.log_to_file ?? false),
        log_file: document.getElementById('log-file') ? document.getElementById('log-file').value : (settings.log_file || 'server.log'),
        reconnect_interval: document.getElementById('reconnect-interval') ? parseInt(document.getElementById('reconnect-interval').value) : (settings.reconnect_interval ?? 5),
        minimize_to_tray: document.getElementById('minimize-to-tray') ? document.getElementById('minimize-to-tray').checked : (settings.minimize_to_tray ?? false),
    });

    const checkDirty = () => {
        const current = collect();
        setDirty(JSON.stringify(current) !== JSON.stringify(original));
    };

    const showDirtyDialog = () => new Promise((resolve) => {
        Modal.open(`
            <div class="modal">
                <div class="modal-header"><h3>${I18N.t('settings.dirtyHint')}</h3></div>
                <div class="modal-body">${I18N.t('settings.dirtyHint')}</div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" data-dirty-cancel>${I18N.t('common.cancel')}</button>
                    <button class="btn btn-secondary" data-dirty-discard>${I18N.t('settings.revert')}</button>
                    <button class="btn btn-primary" data-dirty-save>${I18N.t('common.save')}</button>
                </div>
            </div>
        `);
        const root = document.querySelector('.modal');
        root.querySelector('[data-dirty-cancel]').onclick = () => { Modal.close(); resolve('cancel'); };
        root.querySelector('[data-dirty-discard]').onclick = () => { Modal.close(); resolve('discard'); };
        root.querySelector('[data-dirty-save]').onclick = () => { Modal.close(); resolve('save'); };
    });

    const bind = () => {
        container().querySelectorAll('.settings-nav button').forEach(btn => {
            btn.addEventListener('click', async () => {
                if (dirty) {
                    const choice = await showDirtyDialog();
                    if (choice === 'cancel') return;
                    if (choice === 'save') {
                        await save();
                    }
                }
                section = btn.getAttribute('data-section');
                render();
            });
        });

        container().querySelectorAll('[data-dirty]').forEach(el => {
            el.addEventListener('input', checkDirty);
            el.addEventListener('change', checkDirty);
        });

        document.getElementById('settings-save')?.addEventListener('click', save);
        document.getElementById('settings-revert')?.addEventListener('click', () => {
            settings = JSON.parse(JSON.stringify(original));
            render();
            setDirty(false);
        });
        document.getElementById('settings-reload')?.addEventListener('click', async () => {
            try {
                await Api.reloadSettings();
                Toast.success(I18N.t('settings.reloadSuccess'));
                load();
            } catch (err) {
                Toast.error(I18N.t('settings.reloadError') + ': ' + err.message);
            }
        });
    };

    const save = async () => {
        const data = collect();
        try {
            await Api.saveSettings(data);
            Toast.success(I18N.t('settings.saveSuccess'));
            settings = data;
            original = JSON.parse(JSON.stringify(data));
            setDirty(false);
        } catch (err) {
            Toast.error(I18N.t('settings.saveError') + ': ' + err.message);
        }
    };

    const load = async () => {
        container().innerHTML = `<div class="loading">${I18N.t('common.loading')}</div>`;
        try {
            const data = await Api.settings();
            // Normalize keys
            settings = {
                srt_settings: {
                    latency: data.srt_settings?.latency ?? data.SRTSettings?.Latency ?? 120,
                    connect_timeout: data.srt_settings?.connect_timeout ?? data.SRTSettings?.ConnectTimeout ?? 5000,
                    passphrase: data.srt_settings?.passphrase ?? data.SRTSettings?.Passphrase ?? '',
                    streamid: data.srt_settings?.streamid ?? data.SRTSettings?.StreamID ?? '',
                },
                log_to_file: data.log_to_file ?? data.LogToFile ?? false,
                log_file: data.log_file ?? data.LogFile ?? 'server.log',
                reconnect_interval: data.reconnect_interval ?? data.ReconnectInterval ?? 5,
                minimize_to_tray: data.minimize_to_tray ?? data.MinimizeToTray ?? false,
            };
            original = JSON.parse(JSON.stringify(settings));
            setDirty(false);
            render();
        } catch (err) {
            container().innerHTML = `<div class="card"><div class="card-body" style="color: var(--err);">${I18N.t('settings.loadError')}: ${Formatters.escapeHtml(err.message)}</div></div>`;
        }
    };

    const isDirty = () => dirty;

    return { load, isDirty, save };
})();

/* -------------------- HelpPage -------------------- */
const HelpPage = (() => {
    const container = () => document.getElementById('help-content');

    const render = () => {
        const codeBlock = (lines) => '<pre><code>' + lines.join('\n') + '</code></pre>';
        container().innerHTML = `
            <div class="help-card">
                <h3>▶ ${I18N.t('help.quickstart.title')}</h3>
                <p>${I18N.t('help.quickstart.body')}</p>
            </div>
            <div class="help-card">
                <h3>＋ ${I18N.t('help.inputs.title')}</h3>
                <p>${I18N.t('help.inputs.body')}</p>
                <ul>
                    <li>${I18N.t('help.inputs.field1')}</li>
                    <li>${I18N.t('help.inputs.field2')}</li>
                    <li>${I18N.t('help.inputs.field3')}</li>
                </ul>
                <p><span class="highlight">${I18N.t('help.inputs.note')}</span></p>
            </div>
            <div class="help-card">
                <h3>📡 ${I18N.t('help.protocols.title')}</h3>
                <ul>
                    <li>${I18N.t('help.protocols.rtmp')}</li>
                    <li>${I18N.t('help.protocols.srt')}</li>
                    <li>${I18N.t('help.protocols.whip')}</li>
                </ul>
            </div>
            <div class="help-card">
                <h3>🚀 ${I18N.t('help.outputs.title')}</h3>
                <p>${I18N.t('help.outputs.body')}</p>
                <ul>
                    <li>${I18N.t('help.outputs.way1')}</li>
                    <li>${I18N.t('help.outputs.way2')}</li>
                </ul>
                <p>${I18N.t('help.outputs.formats')}</p>
                ${codeBlock([
                    'rtmp://a.rtmp.youtube.com/live2/xxxx-xxxx    # ' + I18N.t('help.outputs.ex.rtmp'),
                    'srt://relay.example.com:4000                 # ' + I18N.t('help.outputs.ex.srt'),
                    'file:///recordings/stream.flv                # ' + I18N.t('help.outputs.ex.flv'),
                    'file:///recordings/stream.mp4                # ' + I18N.t('help.outputs.ex.mp4'),
                ])}
                <p>${I18N.t('help.outputs.mp4note')}</p>
            </div>
            <div class="help-card">
                <h3>🎬 ${I18N.t('help.ffmpeg.title')}</h3>
                <p>${I18N.t('help.ffmpeg.body')}</p>
                <ul>
                    <li>${I18N.t('help.ffmpeg.windows')}</li>
                    <li>${I18N.t('help.ffmpeg.linux')}</li>
                    <li>${I18N.t('help.ffmpeg.custom')}</li>
                </ul>
                ${codeBlock(['# Windows', 'bin\\ffmpeg.exe', '', '# Linux / macOS', 'bin/ffmpeg'])}
            </div>
            <div class="help-card">
                <h3>🔁 ${I18N.t('help.reconnect.title')}</h3>
                <p>${I18N.t('help.reconnect.body')}</p>
            </div>
            <div class="help-card">
                <h3>⚙ ${I18N.t('help.config.title')}</h3>
                <p>${I18N.t('help.config.body')}</p>
            </div>
        `;
    };

    return { render };
})();

/* -------------------- Palette (Command) -------------------- */
const Palette = (() => {
    let inputs = [];
    let statuses = [];
    let activeIndex = 0;
    let filteredItems = [];
    let modal = null;
    let input = null;
    let list = null;

    const actions = [
        { id: 'addInput', title: () => I18N.t('palette.action.addInput'), run: () => { Router.go('inputs'); setTimeout(() => document.getElementById('input-name')?.focus(), 100); } },
        { id: 'settings', title: () => I18N.t('palette.action.settings'), run: () => Router.go('settings') },
        { id: 'help', title: () => I18N.t('palette.action.help'), run: () => Router.go('help') },
        { id: 'refresh', title: () => I18N.t('palette.action.refresh'), run: () => refresh() },
        { id: 'reload', title: () => I18N.t('palette.action.reloadConfig'), run: async () => { try { await Api.reloadSettings(); Toast.success(I18N.t('settings.reloadSuccess')); } catch (e) { Toast.error(e.message); } } },
    ];

    const collectItems = (query) => {
        const items = [];
        const q = query.toLowerCase().trim();

        const matches = (s) => !q || (s || '').toLowerCase().includes(q);

        actions.forEach(a => {
            const title = a.title();
            if (matches(title)) {
                items.push({ kind: I18N.t('palette.section.actions'), title, sub: '', run: a.run });
            }
        });

        inputs.forEach(inp => {
            const name = inp.Name || inp.name;
            const path = inp.URLPath || inp.url_path;
            if (matches(name) || matches(path)) {
                items.push({
                    kind: I18N.t('palette.section.inputs'),
                    title: name,
                    sub: path,
                    run: () => Router.go('inputs'),
                });
            }
            const outs = inp.Outputs || inp.outputs || [];
            outs.forEach(o => {
                if (matches(o.url)) {
                    items.push({
                        kind: I18N.t('palette.section.outputs'),
                        title: o.url,
                        sub: name,
                        run: () => Router.go('inputs'),
                    });
                }
            });
        });

        return items;
    };

    const renderItems = () => {
        if (filteredItems.length === 0) {
            list.innerHTML = `<div class="palette-empty">${I18N.t('palette.noResults')}</div>`;
            return;
        }
        list.innerHTML = filteredItems.map((it, i) => `
            <div class="palette-item ${i === activeIndex ? 'active' : ''}" data-idx="${i}">
                <span class="kind">${Formatters.escapeHtml(it.kind)}</span>
                <span class="title">${Formatters.escapeHtml(it.title)}</span>
                <span class="sub">${Formatters.escapeHtml(it.sub || '')}</span>
            </div>
        `).join('');
    };

    const refresh = async () => {
        try {
            statuses = await Api.statusAll();
            inputs = statuses.map(s => ({
                Name: s.name,
                URLPath: s.url_path,
                Outputs: s.outputs || [],
            }));
        } catch {
            // fall back to direct fetch
            try { inputs = await Api.inputs(); } catch {}
        }
    };

    const open = async () => {
        await refresh();
        activeIndex = 0;
        filteredItems = collectItems('');

        Modal.custom({
            body: `
                <input type="text" class="palette-search" placeholder="${Formatters.escapeHtml(I18N.t('palette.placeholder'))}" id="palette-input">
                <div class="palette-list" id="palette-list"></div>
            `,
            onClose: () => {},
        });
        modal = document.querySelector('.modal');
        input = document.getElementById('palette-input');
        list = document.getElementById('palette-list');
        modal.classList.add('palette');
        renderItems();
        setTimeout(() => input.focus(), 0);

        input.addEventListener('input', () => {
            activeIndex = 0;
            filteredItems = collectItems(input.value);
            renderItems();
        });

        input.addEventListener('keydown', (e) => {
            if (e.key === 'ArrowDown') {
                e.preventDefault();
                activeIndex = Math.min(activeIndex + 1, filteredItems.length - 1);
                renderItems();
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                activeIndex = Math.max(activeIndex - 1, 0);
                renderItems();
            } else if (e.key === 'Enter') {
                e.preventDefault();
                const item = filteredItems[activeIndex];
                if (item) {
                    Modal.close();
                    item.run();
                }
            }
        });

        list.addEventListener('click', (e) => {
            const el = e.target.closest('.palette-item');
            if (!el) return;
            const idx = parseInt(el.getAttribute('data-idx'));
            const item = filteredItems[idx];
            if (item) {
                Modal.close();
                item.run();
            }
        });
    };

    const close = () => Modal.close();

    return { open, close, refresh };
})();

/* -------------------- Router -------------------- */
const Router = (() => {
    let currentTab = 'status';
    let refreshTimer = null;

    const go = async (name) => {
        if (currentTab === 'settings' && name !== 'settings' && SettingsPage.isDirty()) {
            const ok = await Modal.confirm({
                title: I18N.t('settings.dirtyHint'),
                body: I18N.t('settings.dirtyHint'),
                confirmText: I18N.t('common.save'),
                cancelText: I18N.t('common.cancel'),
            });
            if (!ok) return;
            await SettingsPage.save();
        }
        document.querySelectorAll('.nav-tab').forEach(t => t.classList.toggle('active', t.getAttribute('data-tab') === name));
        document.querySelectorAll('.tab-content').forEach(t => t.classList.toggle('active', t.id === name));
        currentTab = name;
        load();
    };

    const load = () => {
        switch (currentTab) {
            case 'status': Dashboard.load(); break;
            case 'inputs': InputsPage.load(); break;
            case 'settings': SettingsPage.load(); break;
            case 'help': HelpPage.render(); break;
        }
    };

    const refresh = () => {
        if (currentTab === 'status') Dashboard.load();
    };

    const startAutoRefresh = () => {
        if (refreshTimer) clearInterval(refreshTimer);
        refreshTimer = setInterval(() => {
            if (currentTab === 'status') Dashboard.load();
        }, 5000);
    };

    return { go, load, refresh, startAutoRefresh, get current() { return currentTab; } };
})();

/* -------------------- Shortcuts -------------------- */
const Shortcuts = (() => {
    const bind = () => {
        document.addEventListener('keydown', (e) => {
            // Don't intercept while typing in inputs
            const tag = (e.target.tagName || '').toLowerCase();
            const isInput = tag === 'input' || tag === 'textarea' || tag === 'select';
            const cmdK = (e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k';

            if (cmdK) {
                e.preventDefault();
                Palette.open();
                return;
            }

            if (isInput) return;

            if (e.key === '?') {
                e.preventDefault();
                Router.go('help');
            } else if (e.key === '/') {
                e.preventDefault();
                Palette.open();
            } else if (e.key === 'g') {
                // gg = dashboard
                const handler = (e2) => {
                    if (e2.key === 'g') Router.go('status');
                    document.removeEventListener('keydown', handler);
                };
                document.addEventListener('keydown', handler, { once: true });
            } else if (e.key === 'Escape') {
                Modal.close();
            }
        });
    };

    return { bind };
})();

/* -------------------- Edit Output Modal Helper -------------------- */
const showEditOutputModal = (name, url, onSave) => {
    Modal.custom({
        header: I18N.t('outputs.editTitle'),
        body: `
            <form id="edit-output-form" onsubmit="event.preventDefault();">
                <div class="form-group" style="display: flex; flex-direction: column; gap: var(--space-2);">
                    <label style="font-weight: 500; color: var(--text-muted);">${I18N.t('outputs.urlLabel')}</label>
                    <input type="text" class="form-input mono" id="edit-output-url" value="${Formatters.escapeHtml(url)}" style="width: 100%;" required>
                </div>
            </form>
        `,
        footer: `
            <button class="btn btn-secondary" onclick="Modal.close()">${I18N.t('common.cancel')}</button>
            <button class="btn btn-primary" id="edit-output-save-btn">${I18N.t('common.save')}</button>
        `
    });

    const saveBtn = document.getElementById('edit-output-save-btn');
    saveBtn.addEventListener('click', async () => {
        const newUrl = document.getElementById('edit-output-url').value.trim();
        if (!newUrl) return;
        saveBtn.disabled = true;
        try {
            await onSave(newUrl);
            Modal.close();
        } catch (err) {
            Toast.error(err.message);
            saveBtn.disabled = false;
        }
    });
};

/* -------------------- Boot -------------------- */
const Boot = (() => {
    const setupLanguage = () => {
        const sel = document.getElementById('language-select');
        sel.value = I18N.get();
        sel.addEventListener('change', () => {
            I18N.set(sel.value);
            // Re-render current page
            Router.load();
        });
    };

    const setupNav = () => {
        document.querySelectorAll('.nav-tab').forEach(tab => {
            tab.addEventListener('click', () => Router.go(tab.getAttribute('data-tab')));
        });
    };

    const setupTopBar = () => {
        document.getElementById('topbar-refresh').addEventListener('click', () => {
            const icon = document.getElementById('refresh-icon');
            icon.style.transform = 'rotate(360deg)';
            setTimeout(() => { icon.style.transition = 'transform 0.4s'; icon.style.transform = ''; }, 400);
            Router.refresh();
        });
        document.getElementById('topbar-palette').addEventListener('click', () => Palette.open());
        document.getElementById('topbar-add-input').addEventListener('click', () => {
            Router.go('inputs');
            setTimeout(() => document.getElementById('input-name')?.focus(), 100);
        });
    };

    const run = () => {
        I18N.apply();
        setupLanguage();
        setupNav();
        setupTopBar();
        Dashboard.bindActions();
        Shortcuts.bind();
        TopBar.startTicker();
        Router.startAutoRefresh();
        Router.load();
    };

    return { run };
})();

document.addEventListener('DOMContentLoaded', Boot.run);
