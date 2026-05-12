const axios = require('axios');
const api = axios.create({ baseURL: 'http://localhost:8080/api/v1/' });
console.log(api.getUri({ url: 'menus' }));
console.log(api.getUri({ url: '/menus' }));
