import { $http } from "./http";

export default {
    createClientAccount(username, password) {
        return $http.post('/newClient', {username, password})
    }
}