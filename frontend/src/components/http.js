import Axios from 'axios';

const baseUrl = "https://trenlas.herokuapp.com";
export let $http = Axios.create({baseURL: baseUrl});
