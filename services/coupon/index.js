import request from '../../utils/request';

/**
 * 获取可用优惠券列表
 */
export function fetchCouponList(status = 'default') {
  return request.get('/api/v1/coupons')
    .then(coupons => {
      return coupons.map(coupon => ({
        id: coupon.id,
        name: coupon.name,
        type: coupon.type === 1 ? 'price' : 'discount',
        value: coupon.amount * 100, // 转换为分
        base: coupon.min_amount * 100,
        startTime: coupon.start_time,
        endTime: coupon.end_time,
        isAvailable: coupon.total > coupon.used,
        status: 'available'
      }));
    })
    .catch(error => {
      console.error('获取优惠券列表失败:', error);
      return [];
    });
}

/**
 * 领取优惠券
 * @param {number} couponId 优惠券ID
 */
export function claimCoupon(couponId) {
  return request.post(`/api/v1/coupons/${couponId}/claim`);
}

/**
 * 获取我的优惠券
 * @param {string} status 状态筛选
 */
export function fetchUserCouponList(status = 'default') {
  const statusMap = {
    'default': 0,
    'available': 1,
    'used': 2,
    'expired': 3
  };

  const params = status !== 'default' ? { status: statusMap[status] } : {};
  return request.get('/api/v1/coupons/my', params)
    .then(userCoupons => {
      return userCoupons.map(userCoupon => ({
        id: userCoupon.id,
        couponId: userCoupon.coupon_id,
        name: userCoupon.coupon?.name || '',
        type: userCoupon.coupon?.type === 1 ? 'price' : 'discount',
        value: (userCoupon.coupon?.amount || 0) * 100,
        base: (userCoupon.coupon?.min_amount || 0) * 100,
        status: userCoupon.status === 1 ? 'available' : userCoupon.status === 2 ? 'used' : 'expired',
        usedTime: userCoupon.used_time,
        createTime: userCoupon.created_at
      }));
    })
    .catch(error => {
      console.error('获取我的优惠券失败:', error);
      return [];
    });
}

/** 获取优惠券 详情 */
function mockFetchCouponDetail(id, status) {
  const { delay } = require('../_utils/delay');
  const { getCoupon } = require('../../model/coupon');
  const { genAddressList } = require('../../model/address');

  return delay().then(() => {
    const result = {
      detail: getCoupon(id, status),
      storeInfoList: genAddressList(),
    };

    result.detail.useNotes = `1个订单限用1张，除运费券外，不能与其它类型的优惠券叠加使用（运费券除外）\n2.仅适用于各区域正常售卖商品，不支持团购、抢购、预售类商品`;
    result.detail.storeAdapt = `商城通用`;

    if (result.detail.type === 'price') {
      result.detail.desc = `减免 ${result.detail.value / 100} 元`;

      if (result.detail.base) {
        result.detail.desc += `，满${result.detail.base / 100}元可用`;
      }

      result.detail.desc += '。';
    } else if (result.detail.type === 'discount') {
      result.detail.desc = `${result.detail.value}折`;

      if (result.detail.base) {
        result.detail.desc += `，满${result.detail.base / 100}元可用`;
      }

      result.detail.desc += '。';
    }

    return result;
  });
}

/** 获取优惠券 详情 */
export function fetchCouponDetail(id, status = 'default') {
  if (config.useMock) {
    return mockFetchCouponDetail(id, status);
  }
  return new Promise((resolve) => {
    resolve('real api');
  });
}
