import { $http } from "./http";
import $auth from "./auth";

export default {
    createClientAccount(username, password) {
        return $http.post('/newClient', {username, password})
    },

    getAllUsers() {
        return $http.get('/users', {
            headers: {Authorization: $auth.getToken()}
        }).then(res => {
            return res.data || [];
        }).catch(err => {
            if (err.response && err.response.status !== 403) {
                throw err
            } else {
                /// This error is expected for clients/realtors as they can't
                // see all users
                return [];
            }
        })
    }
}