import './style.css';
import './app.css';

function appBinding(name) {
    const fn = window.go?.main?.App?.[name];
    if (typeof fn !== 'function') throw new Error(`Wails binding ${name} is unavailable`);
    return fn;
}

async function init() {
    const versionEl = document.getElementById('version');
    const infoEl = document.getElementById('info');
    const stateEl = document.getElementById('state');
    const logEl = document.getElementById('log');
    try {
        const version = await appBinding('GetVersion')();
        const info = await appBinding('GetInfo')();
        const state = await appBinding('GetUpdateState')();
        versionEl.textContent = version;
        infoEl.textContent = `OS: ${info.os} | Arch: ${info.arch}`;
        stateEl.textContent = state;
        logEl.textContent = `Updater log:\n[init] version=${version}, os=${info.os}, arch=${info.arch}, state=${state}`;
    } catch (err) {
        logEl.textContent = `Error: ${err}`;
        console.error(err);
    }
}

const appEl = document.getElementById('app');
if (appEl) appEl.style.display = 'none';
init();
