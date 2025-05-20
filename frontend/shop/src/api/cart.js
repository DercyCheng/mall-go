import request from './request';

// Get cart list
export function getCartList() {
    return request({
        url: '/cart/list',
        method: 'get'
    });
}

// Get cart list with promotion info
export function getCartListWithPromotion() {
    return request({
        url: '/cart/list/promotion',
        method: 'get'
    });
}

// Add product to cart
export function addToCart(data) {
    return request({
        url: '/cart/add',
        method: 'post',
        data: data
    });
}

// Update cart item quantity
export function updateCartQuantity(id, quantity) {
    return request({
        url: '/cart/update/quantity',
        method: 'get',
        params: {
            id,
            quantity
        }
    });
}

// Delete cart item
export function deleteCartItem(ids) {
    return request({
        url: '/cart/delete',
        method: 'post',
        data: {
            ids
        }
    });
}

// Clear cart
export function clearCart() {
    return request({
        url: '/cart/clear',
        method: 'post'
    });
}
