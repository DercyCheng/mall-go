import { config } from '../../config/index';
import request from '../../utils/request';

/** 获取购物车mock数据 */
function mockFetchCartGroupData(params) {
  const { delay } = require('../_utils/delay');
  const { genCartGroupData } = require('../../model/cart');

  return delay().then(() => genCartGroupData(params));
}

/** 获取购物车数据 */
export function fetchCartGroupData(params) {
  if (config.useMock) {
    return mockFetchCartGroupData(params);
  }

  // 真实API调用
  return request.get('/api/v1/cart')
    .then(cartItems => {
      // 转换数据格式以匹配现有UI组件
      const groupedData = {
        storeGoodsList: [{
          storeName: '商城',
          storeId: 'store_1',
          goodsList: cartItems.map(item => ({
            id: item.id,
            goods: {
              goodsId: item.product_id,
              title: item.product?.name || '',
              price: item.product_sku?.price || item.product?.price || 0,
              originPrice: item.product?.price || 0,
              primaryImage: item.product?.main_image || '',
              skuId: item.sku_id || 0,
              specInfo: item.product_sku?.sku_name || ''
            },
            quantity: item.quantity,
            isChecked: true
          }))
        }],
        isAllSelected: true,
        totalAmount: 0,
        totalDiscountAmount: 0
      };

      // 计算总金额
      groupedData.totalAmount = groupedData.storeGoodsList[0].goodsList.reduce((sum, item) => {
        return sum + (item.goods.price * item.quantity);
      }, 0);

      return groupedData;
    })
    .catch(error => {
      console.error('获取购物车数据失败:', error);
      // 失败时返回mock数据
      return mockFetchCartGroupData(params);
    });
}

/** 添加到购物车 */
export function addToCart(productId, skuId, quantity = 1) {
  if (config.useMock) {
    return Promise.resolve();
  }

  return request.post('/api/v1/cart', {
    product_id: productId,
    sku_id: skuId,
    quantity: quantity
  });
}

/** 更新购物车 */
export function updateCartItem(cartId, quantity) {
  if (config.useMock) {
    return Promise.resolve();
  }

  return request.put(`/api/v1/cart/${cartId}`, {
    quantity: quantity
  });
}

/** 删除购物车项 */
export function removeCartItem(cartId) {
  if (config.useMock) {
    return Promise.resolve();
  }

  return request.delete(`/api/v1/cart/${cartId}`);
}
