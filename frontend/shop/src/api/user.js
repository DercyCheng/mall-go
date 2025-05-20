import request from './request';

// User login
export function login(data) {
    return request({
        url: '/sso/login',
        method: 'post',
        data: data
    });
}

// User register
export function register(data) {
    return request({
        url: '/sso/register',
        method: 'post',
        data: data
    });
}

// Get user info
export function getUserInfo() {
    return request({
        url: '/member/info',
        method: 'get'
    });
}

// Update user info
export function updateUserInfo(data) {
    return request({
        url: '/member/update',
        method: 'post',
        data: data
    });
}

// Get user address list
export function getAddressList() {
    return request({
        url: '/member/address/list',
        method: 'get'
    });
}

// Add new address
export function addAddress(data) {
    return request({
        url: '/member/address/add',
        method: 'post',
        data: data
    });
}

// Update address
export function updateAddress(id, data) {
    return request({
        url: `/member/address/update/${id}`,
        method: 'post',
        data: data
    });
}

// Delete address
export function deleteAddress(id) {
    return request({
        url: `/member/address/delete/${id}`,
        method: 'post'
    });
}

// Get user coupon list
export function getCouponList() {
    return request({
        url: '/member/coupon/list',
        method: 'get'
    });
}

// Get available coupons for cart
export function getAvailableCoupons(cartId) {
    return request({
        url: `/member/coupon/list/cart/${cartId}`,
        method: 'get'
    });
}

// Claim coupon
export function claimCoupon(couponId) {
    return request({
        url: `/member/coupon/add/${couponId}`,
        method: 'post'
    });
}

// Get user favorites list
export function getFavoritesList(params) {
    return request({
        url: '/member/favorites/list',
        method: 'get',
        params: params
    });
}

// Add product to favorites
export function addToFavorites(productId) {
    return request({
        url: '/member/favorites/add',
        method: 'post',
        data: { productId }
    });
}

// Remove product from favorites
export function removeFromFavorites(productId) {
    return request({
        url: `/member/favorites/delete/${productId}`,
        method: 'post'
    });
}

// Check if product is in favorites
export function checkFavoriteStatus(productId) {
    return request({
        url: `/member/favorites/check/${productId}`,
        method: 'get'
    });
}

// Update user password
export function updatePassword(data) {
    return request({
        url: '/member/updatePassword',
        method: 'post',
        data: data
    });
}