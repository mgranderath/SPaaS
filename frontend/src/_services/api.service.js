import { authHeader, url } from "../_helpers";
import { userService } from "../_services";
import axios from "axios";

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
  return axios.get(url(`/api/app`), { headers: authHeader() })
    .then(response => response.data);
}

function createApp(name) {
  const requestOptions = {
    method: "POST",
    headers: authHeader()
  };
  return fetch(url(`/api/app/${name}`), requestOptions)
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
      url(`/api/app/${name}`),
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
      url(`/api/app/${name}/deploy`),
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
      url(`/api/app/${name}/stop`),
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
      url(`/api/app/${name}/start`),
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
  return fetch(url(`/api/app/${name}/logs`), requestOptions)
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
      url(`/api/app/${name}`),
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