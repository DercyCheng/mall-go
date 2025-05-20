import { defineStore } from 'pinia';
import { login, getUserInfo } from '../api/user';

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        userInfo: {},
        roles: []
    }),

    getters: {
        isLoggedIn: (state) => !!state.token,
        username: (state) => state.userInfo.username || ''
    },

    actions: {
        // Login and save token
        async login(userInfo) {
            try {
                const response = await login(userInfo);
                const token = response.data.token;
                this.token = token;
                localStorage.setItem('token', token);
                return Promise.resolve(response);
            } catch (error) {
                return Promise.reject(error);
            }
        },

        // Logout
        logout() {
            this.token = '';
            this.userInfo = {};
            this.roles = [];
            localStorage.removeItem('token');
        },

        // Get user info
        async getUserInfo() {
            try {
                const response = await getUserInfo();
                this.userInfo = response.data;
                this.roles = response.data.roles || [];
                return Promise.resolve(response);
            } catch (error) {
                return Promise.reject(error);
            }
        }
    }
});
