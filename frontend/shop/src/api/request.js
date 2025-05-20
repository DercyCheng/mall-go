import axios from 'axios';

// Create an axios instance
const service = axios.create({
    baseURL: '/api', // API base URL
    timeout: 15000 // Request timeout
});

// Request interceptor
service.interceptors.request.use(
    config => {
        // Add token to request headers if exists
        const token = localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    error => {
        console.error('Request error:', error);
        return Promise.reject(error);
    }
);

// Response interceptor
service.interceptors.response.use(
    response => {
        const res = response.data;

        // Check if the status code is OK (you can adjust this based on your backend API response structure)
        if (res.code === 200 || res.code === undefined) {
            return res;
        } else {
            // Handle other status codes
            console.error('API error:', res.message || 'Error');
            return Promise.reject(new Error(res.message || 'Error'));
        }
    },
    error => {
        console.error('Response error:', error);
        return Promise.reject(error);
    }
);

export default service;
