import './style.css';
import './app.css';

function appBinding(name) {
    const fn = window.go?.main?.App?.[name];
    if (typeof fn !== 'function') throw new Error(`Wails binding ${name} is unavailable`);
    return fn;
}

async function refresh(els) {
    const [state, logs] = await Promise.all([
        appBinding('GetUpdateState')(),
        appBinding('GetLogs')(),
    ]);
    els.state.textContent = state;
    const lines = Array.isArray(logs) ? logs : [];
    els.log.textContent = ['Updater log:', ...lines].join('\n');
    els.log.scrollTop = els.log.scrollHeight;
}

async function init() {
    const els = {
        version: document.getElementById('version'),
        info: document.getElementById('info'),
        state: document.getElementById('state'),
        log: document.getElementById('log'),
    };
    try {
        const [version, info] = await Promise.all([
            appBinding('GetVersion')(),
            appBinding('GetInfo')(),
        ]);
        els.version.textContent = version;
        els.info.textContent = `OS: ${info.os} | Arch: ${info.arch} | Signing key: ${info.keyId}`;
        await refresh(els);
        setInterval(() => { refresh(els).catch((e) => console.error(e)); }, 1500);
    } catch (err) {
        els.log.textContent = `Error: ${err}`;
        console.error(err);
    }
}

const appEl = document.getElementById('app');
if (appEl) appEl.style.display = 'none';
init();
