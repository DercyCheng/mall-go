import { config } from '../../config/index';
import request from '../../utils/request';

/**
 * 获取地址列表
 */
export function fetchAddressList() {
    return request.get('/api/v1/addresses');
}

/**
 * 创建地址
 * @param {object} addressData 地址数据
 */
export function createAddress(addressData) {
    return request.post('/api/v1/addresses', addressData);
}

/**
 * 更新地址
 * @param {number} addressId 地址ID
 * @param {object} addressData 地址数据
 */
export function updateAddress(addressId, addressData) {
    return request.put(`/api/v1/addresses/${addressId}`, addressData);
}

/**
 * 删除地址
 * @param {number} addressId 地址ID
 */
export function deleteAddress(addressId) {
    return request.delete(`/api/v1/addresses/${addressId}`);
}

/**
 * 设置默认地址
 * @param {number} addressId 地址ID
 */
export function setDefaultAddress(addressId) {
    return request.put(`/api/v1/addresses/${addressId}/default`);
}

/**
 * 地址验证
 * @param {object} address 地址对象
 */
export function validateAddress(address) {
    const errors = [];

    if (!address.receiver_name || address.receiver_name.trim() === '') {
        errors.push('请输入收货人姓名');
    }

    if (!address.phone || !/^1[3-9]\d{9}$/.test(address.phone)) {
        errors.push('请输入正确的手机号码');
    }

    if (!address.province || address.province.trim() === '') {
        errors.push('请选择省份');
    }

    if (!address.city || address.city.trim() === '') {
        errors.push('请选择城市');
    }

    if (!address.district || address.district.trim() === '') {
        errors.push('请选择区县');
    }

    if (!address.detail || address.detail.trim() === '') {
        errors.push('请输入详细地址');
    }

    return {
        isValid: errors.length === 0,
        errors
    };
}

/**
 * 格式化地址显示
 * @param {object} address 地址对象
 */
export function formatAddress(address) {
    if (!address) return '';

    return `${address.province}${address.city}${address.district}${address.detail}`;
}

/**
 * 获取定位信息
 */
export function getCurrentLocation() {
    return new Promise((resolve, reject) => {
        wx.getLocation({
            type: 'gcj02',
            success: (res) => {
                // 可以调用地图API进行逆地理编码
                resolve({
                    latitude: res.latitude,
                    longitude: res.longitude
                });
            },
            fail: (err) => {
                console.error('获取定位失败:', err);
                reject(err);
            }
        });
    });
}
