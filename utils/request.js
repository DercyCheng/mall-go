import { config } from '../config/index';

// 统一请求封装
class Request {
    constructor() {
        this.baseURL = config.baseURL;
        this.timeout = 10000;
    }

    // 通用请求方法
    request(options) {
        return new Promise((resolve, reject) => {
            // 获取token
            const token = wx.getStorageSync('token');

            // 设置请求头
            const header = {
                'Content-Type': 'application/json',
                ...options.header
            };

            if (token) {
                header.Authorization = `Bearer ${token}`;
            }

            wx.request({
                url: this.baseURL + options.url,
                method: options.method || 'GET',
                data: options.data || {},
                header,
                timeout: this.timeout,
                success: (res) => {
                    // 统一处理响应
                    if (res.statusCode === 200) {
                        if (res.data.code === 200) {
                            // 兼容新的响应格式
                            if (res.data.data !== undefined) {
                                resolve(res.data.data);
                            } else {
                                resolve(res.data);
                            }
                        } else if (res.data.code === 401) {
                            // token失效，跳转登录
                            wx.removeStorageSync('token');
                            wx.removeStorageSync('userInfo');
                            wx.showToast({
                                title: '请重新登录',
                                icon: 'none'
                            });
                            // 跳转到登录页面
                            wx.switchTab({
                                url: '/pages/user/index'
                            });
                            reject(new Error('未授权'));
                        } else {
                            wx.showToast({
                                title: res.data.message || '请求失败',
                                icon: 'none'
                            });
                            reject(new Error(res.data.message || '请求失败'));
                        }
                    } else {
                        wx.showToast({
                            title: '网络错误',
                            icon: 'none'
                        });
                        reject(new Error('网络错误'));
                    }
                },
                fail: (err) => {
                    console.error('Request failed:', err);
                    wx.showToast({
                        title: '网络连接失败',
                        icon: 'none'
                    });
                    reject(new Error('网络连接失败'));
                }
            });
        });
    }

    // GET请求
    get(url, data = {}, options = {}) {
        return this.request({
            url,
            method: 'GET',
            data,
            ...options
        });
    }

    // POST请求
    post(url, data = {}, options = {}) {
        return this.request({
            url,
            method: 'POST',
            data,
            ...options
        });
    }

    // PUT请求
    put(url, data = {}, options = {}) {
        return this.request({
            url,
            method: 'PUT',
            data,
            ...options
        });
    }

    // DELETE请求
    delete(url, data = {}, options = {}) {
        return this.request({
            url,
            method: 'DELETE',
            data,
            ...options
        });
    }
}

// 创建实例
const request = new Request();

export default request;

// 导出常用方法
export const { get, post, put, delete: del } = request;
