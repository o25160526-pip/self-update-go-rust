import './style.css';
import './app.css';

import {GetVersion, GetInfo, GetUpdateState} from '../wailsjs/go/main/App';

async function init() {
    const versionEl = document.getElementById('version');
    const infoEl = document.getElementById('info');
    const stateEl = document.getElementById('state');
    const logEl = document.getElementById('log');

    try {
        const version = await GetVersion();
        versionEl.textContent = version;

        const info = await GetInfo();
        infoEl.textContent = `OS: ${info.os} | Arch: ${info.arch}`;

        const state = await GetUpdateState();
        stateEl.textContent = state;

        logEl.textContent = `Updater log:\n[init] version=${version}, os=${info.os}, arch=${info.arch}, state=${state}`;
    } catch (err) {
        logEl.textContent = `Error: ${err}`;
        console.error(err);
    }
}

// Remove old #app content and run init
const appEl = document.getElementById('app');
if (appEl) appEl.style.display = 'none';

init();
