const http = require('k6/http')
const { sleep }= require('k6')

export let options = {
    stages: [
        { duration: '40s', target: 20 },
        { duration: '40s', target: 50 },
        { duration: '40s', target: 100 },
        { duration: '40s', target: 200 },
        { duration: '40s', target: 300 },
        { duration: '40s', target: 500 },
        { duration: '40s', target: 500 },
        { duration: '30s', target: 500 },
        { duration: '30s', target: 600 },
        { duration: '30s', target: 800 },
        { duration: '30s', target: 800 },
        { duration: '30s', target: 800 },
        { duration: '30s', target: 800 },
        { duration: '30s', target: 1000 },
        { duration: '30s', target: 1000 },
        { duration: '30s', target: 1000 },
        { duration: '30s', target: 1200 },
        { duration: '30s', target: 1200 },
        { duration: '30s', target: 1200 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1700 },
        { duration: '30s', target: 1700 },
        { duration: '30s', target: 1700 },
        { duration: '30s', target: 1700 },
        { duration: '30s', target: 1700 },
        { duration: '30s', target: 1500 },
        { duration: '30s', target: 1000 },
        { duration: '30s', target: 1000 },
        { duration: '30s', target: 1000 },
        { duration: '40s', target: 500 },
        { duration: '40s', target: 500 },
        { duration: '40s', target: 500 },
        { duration: '40s', target: 500 },
        { duration: '30s', target: 300},
        { duration: '30s', target: 100 },
    ],
    thresholds: {
        http_req_duration: ['p(99)<60'], // 99% of requests must complete below 1.5s
    },
};

const BASE_URL = '127.0.0.1';

export default () => {
    const data = "{\"ip\":\"217.121.148.207\"}";
    http.post(`${BASE_URL}/api/v1/ip/details'`, data, {
        headers: { 'Content-Type': 'application/json' },
    });

    sleep(1);
};
