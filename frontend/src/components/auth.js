import {$http} from "./http";

export default {
    login(username, password) {
        if (!window.localStorage) {
            return;
        }

        if (localStorage.token) {
            return Promise.resolve({data: {token: this.getToken()}});
        }

        return $http.post('/login', {
            "username": username,
            "password": password
        });
    },

    getUserInfo() {
        if (!this.isLoggedIn()) {
            return Promise.reject({});
        }

        return $http.get('/profile', {
            headers: {Authorization: this.getToken()}
        });
    },

    logout() {
        delete localStorage.token;
    },

    getToken() {
        return localStorage.token;
    },

    isLoggedIn() {
        return !!localStorage.token;
    }
};
