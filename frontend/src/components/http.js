import Axios from 'axios';

const baseUrl = "http://localhost:8083";
export let $http = Axios.create({baseURL: baseUrl});
