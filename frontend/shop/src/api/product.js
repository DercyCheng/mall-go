import request from './request';

// Get product list with pagination and filters
export function getProductList(params) {
    return request({
        url: '/product/list',
        method: 'get',
        params: params
    });
}

// Get product detail
export function getProductDetail(id) {
    return request({
        url: `/product/detail/${id}`,
        method: 'get'
    });
}

// Get product categories
export function getCategories() {
    return request({
        url: '/product/category/list',
        method: 'get'
    });
}

// Get category detail
export function getCategoryDetail(id) {
    return request({
        url: `/product/category/${id}`,
        method: 'get'
    });
}

// Get brand list
export function getBrandList() {
    return request({
        url: '/product/brand/list',
        method: 'get'
    });
}

// Get brand detail
export function getBrandDetail(id) {
    return request({
        url: `/product/brand/${id}`,
        method: 'get'
    });
}

// Search products
export function searchProducts(keyword, params) {
    return request({
        url: '/product/search',
        method: 'get',
        params: {
            keyword,
            ...params
        }
    });
}

// Get recommended products
export function getRecommendProducts() {
    return request({
        url: '/product/recommend',
        method: 'get'
    });
}

// Get hot products
export function getHotProducts() {
    return request({
        url: '/product/hot',
        method: 'get'
    });
}

// Get new products
export function getNewProducts() {
    return request({
        url: '/product/new',
        method: 'get'
    });
}
