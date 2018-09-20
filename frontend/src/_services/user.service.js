import { authHeader } from '../_helpers';

export const userService = {
    login,
    logout,
    handleResponse
};

function login(username, password) {
    var formData = new FormData()

    formData.set("username", username)
    formData.set("password", password)

    const requestOptions = {
        method: 'POST',
        body: formData
    };

    return fetch(`/login`, requestOptions)
        .then(handleResponse)
        .then( user => {
            // login successful if there's a jwt token in the response
            if ("token" in user) {
                // store user details and jwt token in local storage to keep user logged in between page refreshes
                var data = {
                    "token": user.token,
                    "username": username
                }
                localStorage.setItem('user', JSON.stringify(data));
                return data
            }
            return user;
        });
}

function logout() {
    // remove user from local storage to log user out
    localStorage.removeItem('user');
}

function handleResponse(response) {
    return response.text().then(text => {
        const data = text && JSON.parse(text);
        if (!response.ok) {
            if (response.status === 401) {
                // auto logout if 401 response returned from api
                logout();
                location.reload(true);
            }
            const error = (data && data.message) || response.statusText;
            return Promise.reject(error);
        }
        return data;
    });
}