import { authHeader } from "../_helpers";
import { userService } from "../_services";

export const apiService = {
  getAll,
  createApp,
  inspectApp,
  deployApp,
  stopApp,
  startApp,
  logs,
  deleteApp
};

function getAll() {
  const requestOptions = {
    method: "GET",
    headers: authHeader()
  };
  return fetch(`/api/app`, requestOptions)
    .then(response => {
      if (!response.ok) {
        if (response.status === 401) {
          // auto logout if 401 response returned from api
          logout();
          location.reload(true);
        }
        const error = (data && data.message) || response.statusText;
        return Promise.reject(error);
      }
      return response.text();
    })
    .then(text => {
      const apps = [];
      const responseObjects = text.split("\n");
      responseObjects
        .filter(value => {
          return value != "";
        })
        .forEach(value => {
          apps.push(JSON.parse(value)["message"]);
        });
      return apps;
    });
}

function createApp(name) {
  const requestOptions = {
    method: "POST",
    headers: authHeader()
  };
  return fetch(`/api/app/${name}`, requestOptions)
    .then(response => {
      if (!response.ok) {
        if (response.status === 401) {
          // auto logout if 401 response returned from api
          logout();
          location.reload(true);
        }
        const error = (data && data.message) || response.statusText;
        return Promise.reject(error);
      }
      return response.body;
    })
    .then(body => {
      return body.getReader();
    });
}

function inspectApp(name) {
  const requestOptions = {
    method: "GET",
    headers: authHeader()
  };
  return fetch(
    `/api/app/${name}`,
    requestOptions
  ).then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.text();
  });
}

function deployApp(name) {
  const requestOptions = {
    method: "POST",
    headers: authHeader()
  };
  return fetch(
    `/api/app/${name}/deploy`,
    requestOptions
  ).then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.text();
  });
}

function stopApp(name) {
  const requestOptions = {
    method: "POST",
    headers: authHeader()
  };
  return fetch(
    `/api/app/${name}/stop`,
    requestOptions
  )
  .then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.text();
  });
}

function startApp(name) {
  const requestOptions = {
    method: "POST",
    headers: authHeader()
  };
  return fetch(
    `/api/app/${name}/start`,
    requestOptions
  )
  .then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.text();
  });
}

function logs(name) {
  const requestOptions = {
    method: "GET",
    headers: authHeader()
  };
  return fetch(`/api/app/${name}/logs`, requestOptions)
  .then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.body;
  })
  .then(body => {
    return body.getReader();
  });
}

function deleteApp(name) {
  const requestOptions = {
    method: "DELETE",
    headers: authHeader()
  };
  return fetch(
    `/api/app/${name}`,
    requestOptions
  )
  .then(response => {
    if (!response.ok) {
      if (response.status === 401) {
        // auto logout if 401 response returned from api
        logout();
        location.reload(true);
      }
      const error = (data && data.message) || response.statusText;
      return Promise.reject(error);
    }
    return response.text();
  });
}