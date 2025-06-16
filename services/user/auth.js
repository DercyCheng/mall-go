import { config } from '../../config/index';
import request from '../../utils/request';

/**
 * 微信登录
 * @param {string} code 微信授权码
 * @param {object} userInfo 用户信息
 */
export function wechatLogin(code, userInfo = {}) {
    return request.post('/api/v1/wechat/login', {
        code: code,
        nick_name: userInfo.nickName || '',
        avatar: userInfo.avatarUrl || '',
        gender: userInfo.gender || 0
    }).then(response => {
        // 存储token和用户信息
        if (response.token) {
            wx.setStorageSync('token', response.token);
            wx.setStorageSync('userInfo', response.user);
        }
        return response;
    });
}

/**
 * 获取用户信息
 */
export function getUserProfile() {
    return request.get('/api/v1/user/profile');
}

/**
 * 更新用户信息
 * @param {object} userInfo 用户信息
 */
export function updateUserProfile(userInfo) {
    return request.put('/api/v1/user/profile', userInfo);
}

/**
 * 检查登录状态
 */
export function checkLoginStatus() {
    const token = wx.getStorageSync('token');
    const userInfo = wx.getStorageSync('userInfo');

    return {
        isLoggedIn: !!(token && userInfo),
        token,
        userInfo
    };
}

/**
 * 退出登录
 */
export function logout() {
    wx.removeStorageSync('token');
    wx.removeStorageSync('userInfo');

    // 跳转到登录页面
    wx.switchTab({
        url: '/pages/user/index'
    });
}

/**
 * 微信授权登录流程
 */
export function wechatAuth() {
    return new Promise((resolve, reject) => {
        // 1. 获取登录凭证
        wx.login({
            success: (loginRes) => {
                if (loginRes.code) {
                    // 2. 获取用户信息
                    wx.getUserProfile({
                        desc: '用于完善会员资料',
                        success: (userRes) => {
                            // 3. 发送到后端登录
                            wechatLogin(loginRes.code, userRes.userInfo)
                                .then(resolve)
                                .catch(reject);
                        },
                        fail: (err) => {
                            console.error('获取用户信息失败:', err);
                            // 即使获取用户信息失败，也尝试登录
                            wechatLogin(loginRes.code)
                                .then(resolve)
                                .catch(reject);
                        }
                    });
                } else {
                    reject(new Error('获取登录凭证失败'));
                }
            },
            fail: (err) => {
                console.error('微信登录失败:', err);
                reject(new Error('微信登录失败'));
            }
        });
    });
}
