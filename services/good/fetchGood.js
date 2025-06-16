import request from '../../utils/request';

/** 获取商品详情 */
export function fetchGood(goodsId) {
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
          isOnSale: product.status === 1,
          etitle: product.sub_title || ''
        }
      };
    })
    .catch(error => {
      console.error('获取商品详情失败:', error);
      throw error;
    });
}
