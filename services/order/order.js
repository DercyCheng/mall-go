import { config } from '../../config/index';
import request from '../../utils/request';

/**
 * 获取订单列表
 * @param {number} page 页码
 * @param {number} limit 每页数量
 * @param {number} status 订单状态
 */
export function fetchOrderList(page = 1, limit = 10, status = 0) {
    const params = { page, limit };
    if (status > 0) {
        params.status = status;
    }

    return request.get('/api/v1/orders', params);
}

/**
 * 获取订单详情
 * @param {number} orderId 订单ID
 */
export function fetchOrderDetail(orderId) {
    return request.get(`/api/v1/orders/${orderId}`);
}

/**
 * 创建订单
 * @param {object} orderData 订单数据
 */
export function createOrder(orderData) {
    return request.post('/api/v1/orders', orderData);
}

/**
 * 支付订单
 * @param {number} orderId 订单ID
 * @param {number} paymentType 支付方式
 */
export function payOrder(orderId, paymentType = 1) {
    return request.put(`/api/v1/orders/${orderId}/pay`, {
        payment_type: paymentType
    });
}

/**
 * 取消订单
 * @param {number} orderId 订单ID
 */
export function cancelOrder(orderId) {
    return request.put(`/api/v1/orders/${orderId}/cancel`);
}

/**
 * 微信支付
 * @param {object} paymentParams 支付参数
 */
export function wechatPay(paymentParams) {
    return new Promise((resolve, reject) => {
        wx.requestPayment({
            ...paymentParams,
            success: (res) => {
                wx.showToast({
                    title: '支付成功',
                    icon: 'success'
                });
                resolve(res);
            },
            fail: (err) => {
                console.error('支付失败:', err);
                if (err.errMsg.includes('cancel')) {
                    wx.showToast({
                        title: '支付已取消',
                        icon: 'none'
                    });
                } else {
                    wx.showToast({
                        title: '支付失败',
                        icon: 'none'
                    });
                }
                reject(err);
            }
        });
    });
}

/**
 * 订单状态映射
 */
export const ORDER_STATUS = {
    1: { text: '待付款', color: '#FA550A' },
    2: { text: '待发货', color: '#FA550A' },
    3: { text: '待收货', color: '#FA550A' },
    4: { text: '已完成', color: '#00A870' },
    5: { text: '已取消', color: '#BBBBBB' }
};

/**
 * 获取订单状态文本
 * @param {number} status 状态码
 */
export function getOrderStatusText(status) {
    return ORDER_STATUS[status]?.text || '未知状态';
}

/**
 * 获取订单状态颜色
 * @param {number} status 状态码
 */
export function getOrderStatusColor(status) {
    return ORDER_STATUS[status]?.color || '#000000';
}
