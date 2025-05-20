<template>
  <div class="orders-container">
    <div class="orders-header">
      <h2 class="section-title">我的订单</h2>
      
      <div class="filter-actions">
        <el-tabs v-model="activeTab" @tab-click="handleTabChange">
          <el-tab-pane label="全部" name="all"></el-tab-pane>
          <el-tab-pane label="待付款" name="unpaid"></el-tab-pane>
          <el-tab-pane label="待发货" name="unshipped"></el-tab-pane>
          <el-tab-pane label="待收货" name="shipped"></el-tab-pane>
          <el-tab-pane label="已完成" name="completed"></el-tab-pane>
          <el-tab-pane label="已取消" name="cancelled"></el-tab-pane>
        </el-tabs>
        
        <div class="search-input">
          <el-input 
            v-model="searchKeyword" 
            placeholder="订单号/商品名称" 
            clearable
            @keyup.enter="searchOrders"
          >
            <template #append>
              <el-button @click="searchOrders">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>
    </div>
    
    <div class="orders-content" v-loading="loading">
      <!-- 无订单时显示 -->
      <div v-if="!loading && orders.length === 0" class="empty-orders">
        <el-empty description="暂无相关订单">
          <el-button type="primary" @click="$router.push('/home')">去购物</el-button>
        </el-empty>
      </div>
      
      <!-- 订单列表 -->
      <div v-else class="order-list">
        <div v-for="order in orders" :key="order.id" class="order-item">
          <div class="order-header">
            <div class="order-info">
              <span class="order-time">{{ formatDate(order.createTime) }}</span>
              <span class="order-number">订单号: {{ order.orderSn }}</span>
            </div>
            <div class="order-status">
              <el-tag :type="getOrderStatusType(order.status)">
                {{ getOrderStatusText(order.status) }}
              </el-tag>
            </div>
          </div>
          
          <div class="order-products">
            <div v-for="item in order.orderItemList" :key="item.id" class="product-item">
              <div class="product-image" @click="goToProduct(item.productId)">
                <img :src="item.productPic" :alt="item.productName" />
              </div>
              <div class="product-info">
                <div class="product-name" @click="goToProduct(item.productId)">
                  {{ item.productName }}
                </div>
                <div class="product-attrs" v-if="item.productAttr">
                  {{ item.productAttr }}
                </div>
              </div>
              <div class="product-price">¥{{ formatPrice(item.price) }}</div>
              <div class="product-quantity">x{{ item.quantity }}</div>
              <div class="product-subtotal">¥{{ formatPrice(item.price * item.quantity) }}</div>
            </div>
          </div>
          
          <div class="order-footer">
            <div class="order-amount">
              共 {{ getTotalQuantity(order) }} 件商品，
              <span class="total-label">实付款：</span>
              <span class="total-price">¥{{ formatPrice(order.payAmount) }}</span>
            </div>
            
            <div class="order-actions">
              <!-- 不同状态显示不同按钮 -->
              <template v-if="order.status === 0">
                <el-button type="primary" size="small" @click="goToPay(order)">去付款</el-button>
                <el-button size="small" @click="cancelOrder(order.id)">取消订单</el-button>
              </template>
              
              <template v-if="order.status === 1">
                <el-button size="small" @click="viewOrderDetail(order.id)">查看详情</el-button>
                <el-button size="small" @click="reminderShipment(order.id)">提醒发货</el-button>
              </template>
              
              <template v-if="order.status === 2">
                <el-button type="primary" size="small" @click="confirmReceive(order.id)">确认收货</el-button>
                <el-button size="small" @click="viewOrderDetail(order.id)">查看详情</el-button>
                <el-button size="small" @click="viewLogistics(order.id)">查看物流</el-button>
              </template>
              
              <template v-if="order.status === 3">
                <el-button type="primary" size="small" @click="reviewOrder(order.id)">评价晒单</el-button>
                <el-button size="small" @click="buyAgain(order)">再次购买</el-button>
                <el-button size="small" @click="viewOrderDetail(order.id)">查看详情</el-button>
              </template>
              
              <template v-if="order.status === 4">
                <el-button size="small" @click="deleteOrder(order.id)">删除订单</el-button>
                <el-button size="small" @click="viewOrderDetail(order.id)">查看详情</el-button>
              </template>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <div class="pagination-section" v-if="orders.length > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Search } from '@element-plus/icons-vue';
import { getOrderList, cancelOrder as apiCancelOrder, deleteOrder as apiDeleteOrder } from '@/api/order';

const router = useRouter();

const loading = ref(false);
const orders = ref([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const activeTab = ref('all');
const searchKeyword = ref('');

// 监听标签页变化
watch(activeTab, () => {
  currentPage.value = 1;
  loadOrders();
});

// 初始化
onMounted(() => {
  loadOrders();
});

// 加载订单列表
const loadOrders = async () => {
  try {
    loading.value = true;
    
    const params = {
      pageNum: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value
    };
    
    // 根据标签设置状态筛选
    if (activeTab.value !== 'all') {
      params.status = getStatusByTab(activeTab.value);
    }
    
    const result = await getOrderList(params);
    orders.value = result.list || [];
    total.value = result.total || 0;
  } catch (error) {
    console.error('加载订单列表失败:', error);
    ElMessage.error('加载订单列表失败');
  } finally {
    loading.value = false;
  }
};

// 格式化价格
const formatPrice = (price) => {
  return Number(price).toFixed(2);
};

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '';
  
  const date = new Date(dateStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

// 根据标签获取订单状态
const getStatusByTab = (tab) => {
  switch (tab) {
    case 'unpaid': return 0;
    case 'unshipped': return 1;
    case 'shipped': return 2;
    case 'completed': return 3;
    case 'cancelled': return 4;
    default: return null;
  }
};

// 获取订单状态文本
const getOrderStatusText = (status) => {
  switch (status) {
    case 0: return '待付款';
    case 1: return '待发货';
    case 2: return '待收货';
    case 3: return '已完成';
    case 4: return '已取消';
    default: return '未知状态';
  }
};

// 获取订单状态对应的标签类型
const getOrderStatusType = (status) => {
  switch (status) {
    case 0: return 'warning';
    case 1: return '';
    case 2: return 'success';
    case 3: return 'success';
    case 4: return 'info';
    default: return '';
  }
};

// 计算订单商品总数
const getTotalQuantity = (order) => {
  if (!order.orderItemList) return 0;
  
  return order.orderItemList.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);
};

// 标签页切换
const handleTabChange = () => {
  // 重置页码
  currentPage.value = 1;
  // 加载订单
  loadOrders();
};

// 每页数量变化
const handleSizeChange = (val) => {
  pageSize.value = val;
  loadOrders();
};

// 页码变化
const handleCurrentChange = (val) => {
  currentPage.value = val;
  loadOrders();
};

// 搜索订单
const searchOrders = () => {
  currentPage.value = 1;
  loadOrders();
};

// 跳转到商品详情
const goToProduct = (productId) => {
  router.push(`/product/${productId}`);
};

// 跳转到支付页面
const goToPay = (order) => {
  router.push(`/payment?orderId=${order.id}&orderSn=${order.orderSn}`);
};

// 取消订单
const cancelOrder = async (orderId) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await apiCancelOrder(orderId);
    
    ElMessage.success('订单已取消');
    
    // 重新加载订单
    loadOrders();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error);
      ElMessage.error('取消订单失败');
    }
  }
};

// 删除订单
const deleteOrder = async (orderId) => {
  try {
    await ElMessageBox.confirm('确定要删除该订单吗？删除后无法恢复。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await apiDeleteOrder(orderId);
    
    ElMessage.success('订单已删除');
    
    // 重新加载订单
    loadOrders();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除订单失败:', error);
      ElMessage.error('删除订单失败');
    }
  }
};

// 查看订单详情
const viewOrderDetail = (orderId) => {
  router.push(`/order/detail/${orderId}`);
};

// 提醒发货
const reminderShipment = (orderId) => {
  ElMessage.success('已提醒商家尽快发货');
};

// 确认收货
const confirmReceive = async (orderId) => {
  try {
    await ElMessageBox.confirm('确认已收到商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    // 这里应该调用确认收货的API
    // 模拟API调用成功
    ElMessage.success('确认收货成功');
    
    // 重新加载订单
    loadOrders();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认收货失败:', error);
      ElMessage.error('确认收货失败');
    }
  }
};

// 查看物流
const viewLogistics = (orderId) => {
  // 模拟跳转到物流页面
  router.push(`/order/logistics/${orderId}`);
};

// 评价订单
const reviewOrder = (orderId) => {
  // 模拟跳转到评价页面
  router.push(`/order/review/${orderId}`);
};

// 再次购买
const buyAgain = (order) => {
  // 模拟重新加入购物车
  ElMessage.success('已将商品加入购物车');
  router.push('/cart');
};
</script>

<style scoped>
.orders-container {
  padding: 10px;
}

.orders-header {
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  color: #303133;
  margin: 0 0 20px;
}

.filter-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
}

.search-input {
  width: 300px;
}

/* 订单列表 */
.order-list {
  margin-bottom: 20px;
}

.order-item {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  overflow: hidden;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.order-info {
  display: flex;
  align-items: center;
  color: #606266;
  font-size: 14px;
}

.order-time {
  margin-right: 20px;
}

/* 订单商品 */
.order-products {
  padding: 15px;
  border-bottom: 1px solid #ebeef5;
}

.product-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px dashed #ebeef5;
}

.product-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.product-image {
  width: 80px;
  height: 80px;
  margin-right: 15px;
  border: 1px solid #ebeef5;
  cursor: pointer;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  margin-bottom: 8px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.product-name:hover {
  color: #409EFF;
}

.product-attrs {
  color: #909399;
  font-size: 12px;
}

.product-price, .product-quantity, .product-subtotal {
  width: 80px;
  text-align: center;
  color: #606266;
}

.product-subtotal {
  color: #F56C6C;
  font-weight: bold;
}

/* 订单底部 */
.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
}

.order-amount {
  color: #606266;
  font-size: 14px;
}

.total-label {
  color: #303133;
  margin-left: 10px;
}

.total-price {
  color: #F56C6C;
  font-size: 18px;
  font-weight: bold;
}

.order-actions {
  display: flex;
}

.order-actions .el-button {
  margin-left: 10px;
}

/* 分页 */
.pagination-section {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 空订单 */
.empty-orders {
  padding: 40px 0;
  text-align: center;
}

@media (max-width: 768px) {
  .filter-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-input {
    width: 100%;
    margin-top: 15px;
  }
  
  .product-item {
    flex-wrap: wrap;
  }
  
  .product-info {
    width: calc(100% - 95px);
    margin-bottom: 10px;
  }
  
  .product-price, .product-quantity, .product-subtotal {
    width: 33.333%;
    margin-left: 95px;
  }
  
  .order-footer {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .order-amount {
    margin-bottom: 15px;
  }
  
  .order-actions {
    display: flex;
    flex-wrap: wrap;
  }
  
  .order-actions .el-button {
    margin-left: 0;
    margin-right: 10px;
    margin-bottom: 10px;
  }
}
</style>
