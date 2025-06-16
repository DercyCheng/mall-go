import request from '../../utils/request';

/**
 * 获取首页数据
 */
export function fetchHome() {
  return Promise.all([
    fetchBanners(),
    fetchCategories()
  ]).then(([banners, categories]) => {
    return {
      swiper: banners,
      tabList: categories
    };
  });
}

/**
 * 获取轮播图
 */
export function fetchBanners() {
  return request.get('/api/v1/banners')
    .then(banners => {
      return banners.map(banner => ({
        img: banner.image,
        text: banner.title,
        url: banner.link || ''
      }));
    })
    .catch(error => {
      console.error('获取轮播图失败:', error);
      return [];
    });
}

/**
 * 获取分类列表作为首页tab
 */
export function fetchCategories() {
  return request.get('/api/v1/products/categories')
    .then(categories => {
      // 添加"精选推荐"作为第一个tab
      const tabList = [
        { text: '精选推荐', key: 0 }
      ];

      // 只取前5个分类作为首页tab
      categories.slice(0, 5).forEach(category => {
        tabList.push({
          text: category.name,
          key: category.id
        });
      });

      return tabList;
    })
    .catch(error => {
      console.error('获取分类失败:', error);
      return [
        { text: '精选推荐', key: 0 }
      ];
    });
}
