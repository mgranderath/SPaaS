var vueInstance = null;

export const alertService = {
  error,
  init,
  success
};

function init(newInstance) {
  vueInstance = newInstance;
}

function error(title, content) {
  vueInstance.$snotify.error(content, title);
}

function success(title, content) {
  vueInstance.$snotify.success(content, title)
}