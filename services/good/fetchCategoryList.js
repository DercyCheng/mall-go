import request from '../../utils/request';

/** 获取分类列表 */
export function getCategoryList() {
  return request.get('/api/v1/products/categories')
    .then(categories => {
      return categories.map(category => ({
        groupId: category.id,
        name: category.name,
        thumbnail: category.icon,
        children: []
      }));
    })
    .catch(error => {
      console.error('获取分类列表失败:', error);
      return [];
    });
}
