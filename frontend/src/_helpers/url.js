import config from "config";

export function url(url) {
    if (config.apiUrl != null) {
        return `${config.apiUrl}${url}`;
    } else {
        return url;
    }
}