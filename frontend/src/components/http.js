import Axios from 'axios';

// const baseUrl = "localhost:8083";
const baseUrl = process.env.VUE_APP_DEBUG;
console.log(baseUrl);
console.log(process.env);
console.log(process.env.NODE_ENV);
export let $http = Axios.create({baseURL: baseUrl});
