<template>
  <div class="cart-container">
    <div class="container">
      <div class="cart-header">
        <h1 class="cart-title">购物车</h1>
        <el-button type="primary" @click="$router.push('/home')">继续购物</el-button>
      </div>
      
      <div class="cart-content" v-loading="loading">
        <!-- 空购物车 -->
        <div v-if="!loading && cartItems.length === 0" class="empty-cart">
          <el-empty description="购物车空空如也，去选购商品吧！">
            <el-button type="primary" @click="$router.push('/home')">去购物</el-button>
          </el-empty>
        </div>
        
        <!-- 购物车列表 -->
        <div v-else class="cart-list">
          <div class="cart-list-header">
            <el-checkbox 
              v-model="allSelected" 
              :indeterminate="isIndeterminate"
              @change="handleSelectAllChange"
            >
              全选
            </el-checkbox>
            <div class="header-item">商品信息</div>
            <div class="header-item">单价</div>
            <div class="header-item">数量</div>
            <div class="header-item">小计</div>
            <div class="header-item">操作</div>
          </div>
          
          <div class="cart-list-body">
            <div 
              v-for="item in cartItems" 
              :key="item.id" 
              class="cart-item"
            >
              <el-checkbox 
                v-model="item.selected" 
                @change="handleSelectChange"
              ></el-checkbox>
              
              <div class="item-info">
                <div class="item-image" @click="goToProduct(item.productId)">
                  <img :src="item.productPic" :alt="item.productName" />
                </div>
                <div class="item-details">
                  <div class="item-name" @click="goToProduct(item.productId)">
                    {{ item.productName }}
                  </div>
                  <div class="item-attrs" v-if="item.productAttr">
                    {{ item.productAttr }}
                  </div>
                </div>
              </div>
              
              <div class="item-price">
                <div class="price-amount">¥{{ formatPrice(item.price) }}</div>
                <div class="price-original" v-if="item.productPrice > item.price">
                  ¥{{ formatPrice(item.productPrice) }}
                </div>
              </div>
              
              <div class="item-quantity">
                <el-input-number 
                  v-model="item.quantity" 
                  :min="1" 
                  :max="item.maxQuantity || 99"
                  @change="(val) => handleQuantityChange(item, val)"
                ></el-input-number>
              </div>
              
              <div class="item-subtotal">
                ¥{{ formatPrice(item.price * item.quantity) }}
              </div>
              
              <div class="item-actions">
                <el-button 
                  type="danger" 
                  size="small" 
                  @click="handleDelete(item.id)"
                >删除</el-button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 购物车底部结算 -->
        <div v-if="cartItems.length > 0" class="cart-footer">
          <div class="footer-left">
            <el-checkbox 
              v-model="allSelected" 
              :indeterminate="isIndeterminate"
              @change="handleSelectAllChange"
            >
              全选
            </el-checkbox>
            <el-button 
              type="text" 
              @click="handleBatchDelete"
              :disabled="selectedItems.length === 0"
            >
              删除选中商品
            </el-button>
            <el-button type="text" @click="handleClearInvalid" v-if="hasInvalidItems">
              清除失效商品
            </el-button>
          </div>
          
          <div class="footer-right">
            <div class="checkout-info">
              <div class="selected-count">
                已选择 <span class="highlight">{{ selectedItems.length }}</span> 件商品
              </div>
              <div class="total-price">
                合计: <span class="highlight">¥{{ formatPrice(totalPrice) }}</span>
              </div>
            </div>
            
            <el-button 
              type="danger" 
              size="large" 
              :disabled="selectedItems.length === 0"
              @click="handleCheckout"
            >
              去结算
            </el-button>
          </div>
        </div>
      </div>
      
      <!-- 推荐商品 -->
      <div class="recommend-section" v-if="recommendProducts.length > 0">
        <div class="section-header">
          <h2 class="section-title">猜你喜欢</h2>
        </div>
        <div class="product-grid">
          <product-card
            v-for="product in recommendProducts"
            :key="product.id"
            :product="product"
          ></product-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { getCartList, deleteCartItem, updateCartQuantity, clearCart } from '@/api/cart';
import { getProductList } from '@/api/product';
import ProductCard from '@/components/ProductCard.vue';
import { useUserStore } from '@/store/user';

const router = useRouter();
const userStore = useUserStore();

const loading = ref(false);
const cartItems = ref([]);
const recommendProducts = ref([]);

// 全选状态
const allSelected = ref(false);
const isIndeterminate = ref(false);

// 是否有失效商品
const hasInvalidItems = computed(() => {
  return cartItems.value.some(item => item.invalid);
});

// 选中的商品
const selectedItems = computed(() => {
  return cartItems.value.filter(item => item.selected && !item.invalid);
});

// 总价
const totalPrice = computed(() => {
  return selectedItems.value.reduce((sum, item) => {
    return sum + (item.price * item.quantity);
  }, 0);
});

// 初始化
onMounted(async () => {
  // 检查用户是否登录
  if (!userStore.isLoggedIn) {
    router.push('/login?redirect=/cart');
    return;
  }
  
  await loadCartData();
  loadRecommendProducts();
});

// 加载购物车数据
const loadCartData = async () => {
  try {
    loading.value = true;
    
    const result = await getCartList();
    cartItems.value = (result || []).map(item => ({
      ...item,
      selected: true,
      invalid: item.stock <= 0
    }));
    
    // 更新全选状态
    updateSelectAllStatus();
  } catch (error) {
    console.error('加载购物车数据失败:', error);
    ElMessage.error('加载购物车失败');
  } finally {
    loading.value = false;
  }
};

// 加载推荐商品
const loadRecommendProducts = async () => {
  try {
    const result = await getProductList({ pageSize: 4, sort: 'recommend' });
    recommendProducts.value = result?.list || [];
  } catch (error) {
    console.error('加载推荐商品失败:', error);
  }
};

// 格式化价格
const formatPrice = (price) => {
  return Number(price).toFixed(2);
};

// 更新全选状态
const updateSelectAllStatus = () => {
  const validItems = cartItems.value.filter(item => !item.invalid);
  if (validItems.length === 0) {
    allSelected.value = false;
    isIndeterminate.value = false;
    return;
  }
  
  const selectedCount = validItems.filter(item => item.selected).length;
  allSelected.value = selectedCount === validItems.length;
  isIndeterminate.value = selectedCount > 0 && selectedCount < validItems.length;
};

// 全选/取消全选
const handleSelectAllChange = (val) => {
  cartItems.value.forEach(item => {
    if (!item.invalid) {
      item.selected = val;
    }
  });
  updateSelectAllStatus();
};

// 单个商品选择状态变化
const handleSelectChange = () => {
  updateSelectAllStatus();
};

// 修改商品数量
const handleQuantityChange = async (item, value) => {
  try {
    await updateCartQuantity(item.id, value);
    // 更新小计
    item.quantity = value;
  } catch (error) {
    console.error('更新数量失败:', error);
    ElMessage.error('更新数量失败');
    // 恢复原数量
    loadCartData();
  }
};

// 删除商品
const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要从购物车中删除该商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await deleteCartItem(id);
    ElMessage.success('删除成功');
    
    // 重新加载购物车数据
    loadCartData();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
};

// 批量删除选中商品
const handleBatchDelete = async () => {
  if (selectedItems.value.length === 0) {
    return;
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedItems.value.length} 件商品吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    
    const selectedIds = selectedItems.value.map(item => item.id);
    const promises = selectedIds.map(id => deleteCartItem(id));
    
    await Promise.all(promises);
    ElMessage.success('删除成功');
    
    // 重新加载购物车数据
    loadCartData();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error);
      ElMessage.error('批量删除失败');
    }
  }
};

// 清除失效商品
const handleClearInvalid = async () => {
  try {
    const invalidItems = cartItems.value.filter(item => item.invalid);
    if (invalidItems.length === 0) {
      return;
    }
    
    await ElMessageBox.confirm('确定要清除所有失效商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    const invalidIds = invalidItems.map(item => item.id);
    const promises = invalidIds.map(id => deleteCartItem(id));
    
    await Promise.all(promises);
    ElMessage.success('清除成功');
    
    // 重新加载购物车数据
    loadCartData();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清除失效商品失败:', error);
      ElMessage.error('清除失效商品失败');
    }
  }
};

// 去结算
const handleCheckout = () => {
  if (selectedItems.value.length === 0) {
    return;
  }
  
  // 跳转到结算页面
  router.push('/checkout');
};

// 跳转到商品详情
const goToProduct = (productId) => {
  router.push(`/product/${productId}`);
};
</script>

<style scoped>
.cart-container {
  background-color: #f5f7fa;
  padding: 20px 0 40px;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* 购物车头部 */
.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.cart-title {
  font-size: 24px;
  color: #303133;
  margin: 0;
}

/* 购物车内容 */
.cart-content {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
  min-height: 200px;
}

/* 空购物车 */
.empty-cart {
  padding: 50px 0;
  text-align: center;
}

/* 购物车列表头部 */
.cart-list-header {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  color: #909399;
}

.cart-list-header .el-checkbox {
  margin-right: 80px;
}

.header-item {
  flex: 1;
  text-align: center;
}

.header-item:nth-child(2) {
  flex: 3;
  text-align: left;
}

/* 购物车列表项 */
.cart-item {
  display: flex;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

.cart-item .el-checkbox {
  margin-right: 20px;
}

.item-info {
  flex: 3;
  display: flex;
  align-items: center;
}

.item-image {
  width: 80px;
  height: 80px;
  margin-right: 15px;
  border: 1px solid #ebeef5;
  cursor: pointer;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.item-details {
  flex: 1;
}

.item-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  cursor: pointer;
}

.item-name:hover {
  color: #409EFF;
}

.item-attrs {
  font-size: 12px;
  color: #909399;
}

.item-price, .item-quantity, .item-subtotal, .item-actions {
  flex: 1;
  text-align: center;
}

.price-amount {
  font-size: 14px;
  color: #303133;
}

.price-original {
  font-size: 12px;
  color: #909399;
  text-decoration: line-through;
}

.item-subtotal {
  color: #F56C6C;
  font-weight: bold;
}

/* 购物车底部 */
.cart-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-top: 1px solid #ebeef5;
}

.footer-left {
  display: flex;
  align-items: center;
}

.footer-left .el-checkbox {
  margin-right: 15px;
}

.footer-right {
  display: flex;
  align-items: center;
}

.checkout-info {
  margin-right: 20px;
  text-align: right;
}

.selected-count {
  margin-bottom: 5px;
  font-size: 14px;
  color: #606266;
}

.total-price {
  font-size: 16px;
  color: #606266;
}

.highlight {
  color: #F56C6C;
  font-weight: bold;
}

/* 推荐商品区域 */
.recommend-section {
  margin-top: 30px;
}

.section-header {
  margin-bottom: 20px;
}

.section-title {
  font-size: 22px;
  margin: 0;
  position: relative;
  padding-left: 15px;
  color: #303133;
}

.section-title::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background-color: #409EFF;
  border-radius: 2px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 992px) {
  .product-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .cart-list-header, .cart-item {
    font-size: 14px;
  }
}

@media (max-width: 768px) {
  .cart-list-header, .cart-item {
    display: block;
    position: relative;
    padding: 15px;
  }
  
  .cart-list-header .header-item {
    display: none;
  }
  
  .cart-list-header .el-checkbox {
    margin-right: 0;
  }
  
  .cart-item .el-checkbox {
    position: absolute;
    top: 15px;
    left: 15px;
    margin-right: 0;
  }
  
  .item-info {
    margin-left: 30px;
    margin-bottom: 15px;
  }
  
  .item-price, .item-quantity, .item-subtotal {
    display: inline-block;
    width: 30%;
    margin-bottom: 15px;
    text-align: left;
  }
  
  .item-actions {
    text-align: left;
  }
  
  .cart-footer {
    flex-direction: column;
  }
  
  .footer-left, .footer-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .footer-left {
    margin-bottom: 15px;
  }
  
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .cart-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .cart-title {
    margin-bottom: 15px;
  }
}
</style>
