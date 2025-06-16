import { config } from '../../config/index';
import request from '../../utils/request';

/** 获取商品列表 */
export function fetchGoodsList(pageIndex = 1, pageSize = 20, categoryId = 0, keyword = '') {
  const params = {
    page: pageIndex,
    limit: pageSize
  };

  if (categoryId) {
    params.category_id = categoryId;
  }

  if (keyword) {
    params.keyword = keyword;
  }

  return request.get('/api/v1/products', params)
    .then(response => {
      // 适配新的响应格式
      const list = response.data || response.list || [];
      return list.map(item => ({
        spuId: item.id,
        thumb: item.main_image,
        title: item.name,
        price: item.price,
        originPrice: item.price,
        tags: item.sub_title ? [item.sub_title] : [],
        etitle: item.sub_title || ''
      }));
    })
    .catch(error => {
      console.error('获取商品列表失败:', error);
      return [];
    });
}

/** 获取商品详情 */
export function fetchGoodsDetail(goodsId) {
  return request.get(`/api/v1/products/${goodsId}`)
    .then(product => {
      return {
        goods: {
          goodsId: product.id,
          title: product.name,
          price: product.price,
          originPrice: product.price,
          primaryImage: product.main_image,
          images: product.sub_images ? product.sub_images.split(',') : [product.main_image],
          detail: product.detail,
          stock: product.stock,
          category: product.category?.name || '',
          isOnSale: product.status === 1
        }
      };
    })
    .catch(error => {
      console.error('获取商品详情失败:', error);
      throw error;
    });
}
