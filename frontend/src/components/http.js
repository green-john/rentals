import Axios from 'axios';

const baseUrl = process.env.VUE_APP_API_ADDR;
export let $http = Axios.create({baseURL: baseUrl});
