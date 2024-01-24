import {apps} from './xxoo/apps';

export function GetApps(corp) {
    if (corp === 'xxoo') {
        return apps;
    }
    return [];
}

export function GetApp(corp, id) {
    let apps = GetApps(corp);
    for (let i = 0; i < apps.length; i++) {
        if (apps[i].id === id) {
            return apps[i];
        }
    }

    return null;
}