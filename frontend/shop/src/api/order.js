import request from './request';

// Generate order confirmation
export function generateConfirmOrder() {
    return request({
        url: '/order/confirmOrder',
        method: 'post'
    });
}

// Submit order
export function submitOrder(data) {
    return request({
        url: '/order/generateOrder',
        method: 'post',
        data: data
    });
}

// Get order detail
export function getOrderDetail(id) {
    return request({
        url: `/order/detail/${id}`,
        method: 'get'
    });
}

// Get order list
export function getOrderList(params) {
    return request({
        url: '/order/list',
        method: 'get',
        params: params
    });
}

// Cancel order
export function cancelOrder(orderId) {
    return request({
        url: '/order/cancelOrder',
        method: 'post',
        params: {
            orderId
        }
    });
}

// Pay order success callback
export function paySuccess(orderId) {
    return request({
        url: '/order/paySuccess',
        method: 'post',
        params: {
            orderId
        }
    });
}

// Confirm receipt
export function confirmReceipt(orderId) {
    return request({
        url: '/order/confirmReceipt',
        method: 'post',
        params: {
            orderId
        }
    });
}

// Delete order (alias for cancelOrder)
export function deleteOrder(orderId) {
    return cancelOrder(orderId);
}

// Confirm receive order (alias for confirmReceipt)
export function confirmReceiveOrder(orderId) {
    return confirmReceipt(orderId);
}