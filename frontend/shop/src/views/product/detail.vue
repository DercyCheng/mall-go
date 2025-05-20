<template>
  <div class="product-detail-container" v-loading="loading">
    <div class="container" v-if="product.id">
      <div class="product-detail">
        <!-- 商品基本信息区域 -->
        <div class="product-basic">
          <!-- 商品图片展示 -->
          <div class="product-gallery">
            <el-carousel in    };
    
    await addToCartApi({
      productId: product.value.id,
      quantity: quantity.value,
      productAttrValues: productAttrValues
    });
    
    ElMessage.success('添加成功');osition="outside" height="400px">
              <el-carousel-item v-for="(image, index) in productImages" :key="index">
                <img :src="image" :alt="product.name" class="product-image" />
              </el-carousel-item>
            </el-carousel>
          </div>

          <!-- 商品信息 -->
          <div class="product-info">
            <h1 class="product-name">{{ product.name }}</h1>
            
            <div class="product-tags" v-if="product.isNew || product.isHot">
              <el-tag type="success" v-if="product.isNew">新品</el-tag>
              <el-tag type="danger" v-if="product.isHot">热卖</el-tag>
            </div>
            
            <div class="product-price-section">
              <div class="product-price">¥{{ formatPrice(product.price) }}</div>
              <div class="product-original-price" v-if="product.originalPrice > product.price">
                <span class="original-price">原价: ¥{{ formatPrice(product.originalPrice) }}</span>
                <span class="discount">{{ calculateDiscount(product.price, product.originalPrice) }}折</span>
              </div>
            </div>
            
            <div class="product-stats">
              <div class="stat-item">
                <span class="stat-label">销量</span>
                <span class="stat-value">{{ product.sale || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">库存</span>
                <span class="stat-value">{{ product.stock || 0 }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">评价</span>
                <span class="stat-value">{{ product.commentCount || 0 }}</span>
              </div>
            </div>

            <!-- 商品规格选择 -->
            <div class="product-attrs" v-if="product.productAttributeList && product.productAttributeList.length > 0">
              <div 
                v-for="attr in product.productAttributeList" 
                :key="attr.id" 
                class="attr-item"
              >
                <div class="attr-name">{{ attr.name }}</div>
                <div class="attr-values">
                  <el-radio-group v-model="selectedAttrs[attr.id]">
                    <el-radio 
                      v-for="value in attr.values" 
                      :key="value" 
                      :label="value"
                      class="attr-value"
                    >
                      {{ value }}
                    </el-radio>
                  </el-radio-group>
                </div>
              </div>
            </div>

            <!-- 购买数量 -->
            <div class="product-quantity">
              <span class="quantity-label">数量</span>
              <el-input-number 
                v-model="quantity" 
                :min="1" 
                :max="product.stock || 999"
                size="large"
              />
              <span class="stock-info">库存{{ product.stock || 0 }}件</span>
            </div>

            <!-- 购买按钮 -->
            <div class="product-actions">
              <el-button 
                type="primary" 
                size="large" 
                @click="handleAddToCart" 
                :disabled="product.stock <= 0"
                :loading="addingToCart"
              >
                <el-icon><ShoppingCart /></el-icon>
                加入购物车
              </el-button>
              
              <el-button 
                type="danger" 
                size="large" 
                @click="buyNow" 
                :disabled="product.stock <= 0"
              >
                立即购买
              </el-button>
              
              <el-button
                :type="isFavorite ? 'warning' : 'info'"
                size="large"
                @click="toggleFavorite"
                :loading="favoriteLoading"
              >
                <el-icon>
                  <StarFilled v-if="isFavorite" />
                  <Star v-else />
                </el-icon>
                {{ isFavorite ? '已收藏' : '收藏' }}
              </el-button>
            </div>

            <!-- 服务承诺 -->
            <div class="product-services">
              <div class="service-item">
                <el-icon><Check /></el-icon>
                <span>正品保障</span>
              </div>
              <div class="service-item">
                <el-icon><Check /></el-icon>
                <span>7天无理由退换</span>
              </div>
              <div class="service-item">
                <el-icon><Check /></el-icon>
                <span>48小时发货</span>
              </div>
              <div class="service-item">
                <el-icon><Check /></el-icon>
                <span>满88元包邮</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 商品详情与评价 Tab -->
        <div class="product-detail-tabs">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="商品详情" name="detail">
              <div class="product-detail-content" v-html="product.detailContent"></div>
            </el-tab-pane>
            
            <el-tab-pane label="规格参数" name="params">
              <div class="product-params">
                <el-table :data="product.productParameterList || []" border>
                  <el-table-column prop="name" label="参数名" />
                  <el-table-column prop="value" label="参数值" />
                </el-table>
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="用户评价" name="comment">
              <div class="product-comments">
                <div v-if="productComments.length === 0" class="no-comments">
                  暂无评价
                </div>
                <div v-else class="comment-list">
                  <div 
                    v-for="comment in productComments" 
                    :key="comment.id" 
                    class="comment-item"
                  >
                    <div class="comment-user">
                      <el-avatar :size="40" :src="comment.avatar"></el-avatar>
                      <span class="username">{{ comment.username }}</span>
                    </div>
                    <div class="comment-content">
                      <div class="comment-rating">
                        <el-rate v-model="comment.star" disabled></el-rate>
                        <span class="comment-time">{{ comment.createTime }}</span>
                      </div>
                      <div class="comment-text">{{ comment.content }}</div>
                      <div class="comment-images" v-if="comment.images && comment.images.length > 0">
                        <el-image 
                          v-for="(img, index) in comment.images" 
                          :key="index" 
                          :src="img" 
                          :preview-src-list="comment.images"
                          class="comment-image"
                        ></el-image>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>
    
    <div class="not-found" v-else-if="!loading">
      <el-empty description="商品不存在或已下架"></el-empty>
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { ShoppingCart, Check, Star, StarFilled } from '@element-plus/icons-vue';
import { getProductDetail } from '@/api/product';
import { addToCart as addToCartApi } from '@/api/cart';
import { useUserStore } from '@/store/user';
import { addToFavorites, removeFromFavorites, checkFavoriteStatus } from '@/api/user';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const loading = ref(true);
const addingToCart = ref(false);
const product = ref({});
const quantity = ref(1);
const activeTab = ref('detail');
const selectedAttrs = reactive({});
const productComments = ref([]);
const isFavorite = ref(false);
const favoriteLoading = ref(false);

// 获取商品详情
onMounted(async () => {
  const productId = route.params.id;
  if (!productId) {
    router.push('/');
    return;
  }
  
  try {
    loading.value = true;
    const res = await getProductDetail(productId);
    product.value = res || {};
    
    // 初始化默认选中的规格
    if (product.value.productAttributeList) {
      product.value.productAttributeList.forEach(attr => {
        if (attr.values && attr.values.length > 0) {
          selectedAttrs[attr.id] = attr.values[0];
        }
      });
    }
    
    // 获取收藏状态
    if (userStore.isLoggedIn) {
      checkProductFavoriteStatus();
    }
    
    // 模拟评论数据
    loadMockComments();
  } catch (error) {
    console.error('获取商品详情失败:', error);
    ElMessage.error('获取商品详情失败');
  } finally {
    loading.value = false;
  }
});

// 商品图片列表
const productImages = computed(() => {
  if (!product.value || !product.value.albumPics) {
    return [product.value.pic || ''];
  }
  
  const pics = product.value.albumPics.split(',');
  if (product.value.pic && !pics.includes(product.value.pic)) {
    pics.unshift(product.value.pic);
  }
  
  return pics;
});

// 格式化价格
const formatPrice = (price) => {
  return Number(price).toFixed(2);
};

// 计算折扣
const calculateDiscount = (price, originalPrice) => {
  if (!originalPrice || originalPrice <= 0 || !price) return 10;
  return ((price / originalPrice) * 10).toFixed(1);
};

// 添加到购物车
const handleAddToCart = async () => {
  if (!userStore.isLoggedIn) {
    router.push({
      path: '/login',
      query: { redirect: route.fullPath }
    });
    return;
  }
  
  try {
    addingToCart.value = true;
    
    // 构建商品属性参数
    const productAttrValues = Object.entries(selectedAttrs).map(([key, value]) => {
      return {
        attributeId: key,
        value: value
      };
    });
    
    await addToCart({
      productId: product.value.id,
      quantity: quantity.value,
      productAttrValues: productAttrValues
    });
    
    ElMessage.success('添加成功');
  } catch (error) {
    console.error('添加购物车失败:', error);
    ElMessage.error('添加购物车失败');
  } finally {
    addingToCart.value = false;
  }
};

// 立即购买
const buyNow = () => {
  if (!userStore.isLoggedIn) {
    router.push({
      path: '/login',
      query: { redirect: route.fullPath }
    });
    return;
  }

  // TODO: 实现立即购买逻辑
  ElMessage.info('购买功能开发中...');
};

// 检查商品收藏状态
const checkProductFavoriteStatus = async () => {
  try {
    const res = await checkFavoriteStatus(product.value.id);
    isFavorite.value = res.data;
  } catch (error) {
    console.error('获取收藏状态失败:', error);
  }
};

// 切换收藏状态
const toggleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    router.push({
      path: '/login',
      query: { redirect: route.fullPath }
    });
    return;
  }

  try {
    favoriteLoading.value = true;
    if (isFavorite.value) {
      await removeFromFavorites(product.value.id);
      ElMessage.success('已取消收藏');
    } else {
      await addToFavorites(product.value.id);
      ElMessage.success('收藏成功');
    }
    isFavorite.value = !isFavorite.value;
  } catch (error) {
    console.error('收藏操作失败:', error);
    ElMessage.error('操作失败，请稍后重试');
  } finally {
    favoriteLoading.value = false;
  }
};

// 加载模拟评论数据
const loadMockComments = () => {
  productComments.value = [
    {
      id: 1,
      username: '用户1234567',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      content: '商品质量非常好，包装也很精美，物流很快，很满意的一次购物体验！',
      star: 5,
      createTime: '2025-05-10 10:23:45',
      images: [
        'https://img0.baidu.com/it/u=329862679,2733166476&fm=253&fmt=auto&app=138&f=JPEG'
      ]
    },
    {
      id: 2,
      username: '张三',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      content: '不错的产品，使用感受很好。',
      star: 4,
      createTime: '2025-05-08 16:43:12',
      images: []
    },
    {
      id: 3,
      username: '李四',
      avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
      content: '产品外观很漂亮，但是使用体验一般，希望商家能够改进。',
      star: 3,
      createTime: '2025-05-05 09:15:32',
      images: [
        'https://img0.baidu.com/it/u=329862679,2733166476&fm=253&fmt=auto&app=138&f=JPEG',
        'https://img2.baidu.com/it/u=1003272215,1878948666&fm=253&fmt=auto&app=120&f=JPEG'
      ]
    }
  ];
};
</script>

<style scoped>
.product-detail-container {
  background-color: #f5f7fa;
  padding: 20px 0 40px;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.product-detail {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* 商品基本信息区域 */
.product-basic {
  display: flex;
  padding: 30px;
  border-bottom: 1px solid #ebeef5;
}

.product-gallery {
  width: 400px;
  margin-right: 30px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.product-info {
  flex: 1;
}

.product-name {
  font-size: 24px;
  color: #303133;
  margin: 0 0 15px;
}

.product-tags {
  margin-bottom: 15px;
}

.product-tags .el-tag {
  margin-right: 10px;
}

.product-price-section {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.product-price {
  font-size: 28px;
  color: #f56c6c;
  font-weight: bold;
}

.product-original-price {
  display: flex;
  align-items: center;
  margin-top: 8px;
}

.original-price {
  color: #909399;
  text-decoration: line-through;
  margin-right: 10px;
}

.discount {
  background-color: #f56c6c;
  color: #fff;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.product-stats {
  display: flex;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px dashed #ebeef5;
}

.stat-item {
  margin-right: 30px;
}

.stat-label {
  color: #909399;
  margin-right: 5px;
}

.stat-value {
  color: #606266;
}

/* 商品规格选择 */
.product-attrs {
  margin-bottom: 20px;
}

.attr-item {
  margin-bottom: 15px;
}

.attr-name {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.attr-values {
  display: flex;
  flex-wrap: wrap;
}

.attr-value {
  margin-right: 10px;
  margin-bottom: 10px;
}

/* 购买数量 */
.product-quantity {
  display: flex;
  align-items: center;
  margin-bottom: 25px;
}

.quantity-label {
  color: #606266;
  margin-right: 10px;
}

.stock-info {
  color: #909399;
  margin-left: 15px;
  font-size: 14px;
}

/* 购买按钮 */
.product-actions {
  display: flex;
  margin-bottom: 25px;
}

.product-actions .el-button {
  margin-right: 15px;
  padding: 12px 25px;
}

/* 服务承诺 */
.product-services {
  display: flex;
  flex-wrap: wrap;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px dashed #ebeef5;
}

.service-item {
  display: flex;
  align-items: center;
  margin-right: 20px;
  color: #606266;
  font-size: 14px;
}

.service-item .el-icon {
  color: #67c23a;
  margin-right: 5px;
}

/* 商品详情 Tab */
.product-detail-tabs {
  padding: 20px 30px 30px;
}

.product-detail-content {
  min-height: 300px;
  padding: 20px 0;
}

/* 评论区样式 */
.no-comments {
  text-align: center;
  padding: 50px 0;
  color: #909399;
}

.comment-list {
  padding: 20px 0;
}

.comment-item {
  display: flex;
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.comment-user {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 80px;
  margin-right: 20px;
}

.username {
  margin-top: 10px;
  font-size: 14px;
  color: #606266;
  width: 100%;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-content {
  flex: 1;
}

.comment-rating {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.comment-time {
  color: #909399;
  font-size: 14px;
}

.comment-text {
  margin-bottom: 15px;
  line-height: 1.5;
  color: #303133;
}

.comment-images {
  display: flex;
  flex-wrap: wrap;
}

.comment-image {
  width: 80px;
  height: 80px;
  margin-right: 10px;
  margin-bottom: 10px;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
}

/* 未找到商品样式 */
.not-found {
  text-align: center;
  padding: 100px 0;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 992px) {
  .product-basic {
    flex-direction: column;
  }
  
  .product-gallery {
    width: 100%;
    margin-right: 0;
    margin-bottom: 30px;
  }
}

@media (max-width: 768px) {
  .product-actions {
    flex-direction: column;
  }
  
  .product-actions .el-button {
    margin-right: 0;
    margin-bottom: 15px;
  }
}
</style>
