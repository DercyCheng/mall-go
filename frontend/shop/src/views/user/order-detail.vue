<template>
  <div class="order-detail-container" v-loading="loading">
    <div v-if="order.id">
      <el-card class="order-card">
        <div class="order-header">
          <h2 class="order-title">订单详情</h2>
          <div class="order-status">
            <el-tag :type="getStatusTagType(order.status)">{{ getStatusText(order.status) }}</el-tag>
          </div>
        </div>
        
        <div class="order-info">
          <div class="info-item">
            <span class="label">订单编号：</span>
            <span class="value">{{ order.orderSn }}</span>
            <el-button type="text" @click="copyOrderSn">复制</el-button>
          </div>
          <div class="info-item">
            <span class="label">下单时间：</span>
            <span class="value">{{ formatDate(order.createTime) }}</span>
          </div>
          <div class="info-item">
            <span class="label">支付方式：</span>
            <span class="value">{{ getPaymentText(order.payType) }}</span>
          </div>
        </div>
      </el-card>
      
      <!-- 物流信息 -->
      <el-card class="order-card logistics-card" v-if="order.status >= 2">
        <template #header>
          <div class="card-header">
            <h3>物流信息</h3>
          </div>
        </template>
        
        <div class="logistics-info">
          <div class="info-item">
            <span class="label">物流公司：</span>
            <span class="value">{{ order.deliveryCompany || '暂无' }}</span>
          </div>
          <div class="info-item">
            <span class="label">物流单号：</span>
            <span class="value">{{ order.deliverySn || '暂无' }}</span>
            <el-button type="text" v-if="order.deliverySn" @click="copyTrackingNumber">复制</el-button>
          </div>
        </div>
        
        <!-- 物流跟踪时间线 -->
        <div class="logistics-timeline" v-if="logisticsData.length > 0">
          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in logisticsData"
              :key="index"
              :timestamp="activity.time"
            >
              {{ activity.context }}
            </el-timeline-item>
          </el-timeline>
        </div>
        <div v-else class="no-logistics">
          <el-empty description="暂无物流跟踪信息"></el-empty>
        </div>
      </el-card>
      
      <!-- 收货信息 -->
      <el-card class="order-card">
        <template #header>
          <div class="card-header">
            <h3>收货信息</h3>
          </div>
        </template>
        
        <div class="address-info">
          <div class="info-item">
            <span class="label">收货人：</span>
            <span class="value">{{ order.receiverName }}</span>
          </div>
          <div class="info-item">
            <span class="label">联系电话：</span>
            <span class="value">{{ maskPhone(order.receiverPhone) }}</span>
          </div>
          <div class="info-item">
            <span class="label">收货地址：</span>
            <span class="value">{{ getFullAddress() }}</span>
          </div>
        </div>
      </el-card>
      
      <!-- 商品信息 -->
      <el-card class="order-card">
        <template #header>
          <div class="card-header">
            <h3>商品信息</h3>
          </div>
        </template>
        
        <div class="product-list">
          <div class="product-item" v-for="item in order.orderItemList" :key="item.id">
            <div class="product-img" @click="viewProduct(item.productId)">
              <img :src="item.productPic" :alt="item.productName">
            </div>
            <div class="product-info">
              <div class="product-name" @click="viewProduct(item.productId)">{{ item.productName }}</div>
              <div class="product-sku" v-if="item.productAttr">{{ item.productAttr }}</div>
            </div>
            <div class="product-price">¥{{ item.price.toFixed(2) }}</div>
            <div class="product-quantity">x{{ item.quantity }}</div>
            <div class="product-subtotal">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
          </div>
        </div>
        
        <!-- 订单金额信息 -->
        <div class="order-amount">
          <div class="amount-item">
            <span class="label">商品总价：</span>
            <span class="value">¥{{ order.totalAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-item">
            <span class="label">运费：</span>
            <span class="value">¥{{ order.freightAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-item" v-if="order.couponAmount > 0">
            <span class="label">优惠券：</span>
            <span class="value">-¥{{ order.couponAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-item" v-if="order.promotionAmount > 0">
            <span class="label">促销优惠：</span>
            <span class="value">-¥{{ order.promotionAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-item" v-if="order.integrationAmount > 0">
            <span class="label">积分抵扣：</span>
            <span class="value">-¥{{ order.integrationAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-item total">
            <span class="label">实付金额：</span>
            <span class="value price">¥{{ order.payAmount.toFixed(2) }}</span>
          </div>
        </div>
      </el-card>
      
      <!-- 操作按钮 -->
      <div class="order-actions">
        <el-button 
          v-if="order.status === 0" 
          type="danger" 
          @click="handleCancelOrder"
        >
          取消订单
        </el-button>
        <el-button 
          v-if="order.status === 0" 
          type="primary" 
          @click="payOrder"
        >
          去支付
        </el-button>
        <el-button 
          v-if="order.status === 2" 
          type="primary" 
          @click="confirmReceive"
        >
          确认收货
        </el-button>
        <el-button 
          v-if="order.status === 3" 
          type="primary" 
          @click="viewComment"
        >
          评价商品
        </el-button>
        <el-button 
          v-if="order.status === 4 || order.status === 5" 
          type="primary" 
          @click="handleDeleteOrder"
        >
          删除订单
        </el-button>
        <el-button @click="back">返回列表</el-button>
      </div>
    </div>
    
    <div v-else-if="!loading" class="not-found">
      <el-empty description="订单不存在或已被删除"></el-empty>
      <el-button type="primary" @click="back">返回订单列表</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getOrderDetail, cancelOrder as cancelOrderApi, deleteOrder as deleteOrderApi, confirmReceiveOrder } from '@/api/order';

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const order = ref({});
const logisticsData = ref([]);

onMounted(async () => {
  const orderId = route.params.id;
  if (!orderId) {
    back();
    return;
  }
  
  await fetchOrderDetail(orderId);
});

// 获取订单详情
const fetchOrderDetail = async (orderId) => {
  try {
    loading.value = true;
    const res = await getOrderDetail(orderId);
    order.value = res.data || {};
    
    // 如果有物流信息，则获取物流详情
    if (order.value.status >= 2 && order.value.deliverySn) {
      await fetchLogisticsInfo();
    }
  } catch (error) {
    console.error('获取订单详情失败:', error);
    ElMessage.error('获取订单详情失败');
  } finally {
    loading.value = false;
  }
};

// 获取物流信息
const fetchLogisticsInfo = async () => {
  try {
    // 模拟物流数据
    logisticsData.value = [
      { time: '2023-05-15 18:30:00', context: '【收货地址】已签收，签收人：本人' },
      { time: '2023-05-15 11:20:00', context: '【配送员】正在为您配送，配送员：张三，联系电话：1380000****' },
      { time: '2023-05-14 20:30:00', context: '【配送站】已出库，准备配送' },
      { time: '2023-05-13 09:10:00', context: '【仓库】包裹已出库' },
      { time: '2023-05-12 21:00:00', context: '【仓库】商品已打包完成' },
      { time: '2023-05-12 14:20:00', context: '【系统】订单已进入物流系统' }
    ];
  } catch (error) {
    console.error('获取物流信息失败:', error);
    ElMessage.error('获取物流信息失败');
  }
};

// 返回订单列表
const back = () => {
  router.push('/user/orders');
};

// 查看商品详情
const viewProduct = (productId) => {
  router.push(`/product/${productId}`);
};

// 复制订单编号
const copyOrderSn = () => {
  navigator.clipboard.writeText(order.value.orderSn)
    .then(() => {
      ElMessage.success('订单编号已复制');
    })
    .catch(() => {
      ElMessage.error('复制失败，请手动复制');
    });
};

// 复制物流单号
const copyTrackingNumber = () => {
  navigator.clipboard.writeText(order.value.deliverySn)
    .then(() => {
      ElMessage.success('物流单号已复制');
    })
    .catch(() => {
      ElMessage.error('复制失败，请手动复制');
    });
};

// 获取完整地址
const getFullAddress = () => {
  const { receiverProvince, receiverCity, receiverRegion, receiverDetailAddress } = order.value;
  return [receiverProvince, receiverCity, receiverRegion, receiverDetailAddress]
    .filter(item => item)
    .join(' ');
};

// 手机号脱敏处理
const maskPhone = (phone) => {
  if (!phone) return '';
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2');
};

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  return date.toLocaleString();
};

// 获取订单状态文本
const getStatusText = (status) => {
  const statusMap = {
    0: '待付款',
    1: '待发货',
    2: '已发货',
    3: '已完成',
    4: '已关闭',
    5: '无效订单'
  };
  return statusMap[status] || '未知状态';
};

// 获取订单状态标签类型
const getStatusTagType = (status) => {
  const typeMap = {
    0: 'warning',
    1: 'info',
    2: 'primary',
    3: 'success',
    4: 'info',
    5: 'danger'
  };
  return typeMap[status] || 'info';
};

// 获取支付方式文本
const getPaymentText = (payType) => {
  const payMap = {
    1: '支付宝',
    2: '微信支付',
    3: '银联支付',
    4: '货到付款'
  };
  return payMap[payType] || '未知方式';
};

// 取消订单
const handleCancelOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await cancelOrderApi(order.value.id);
    ElMessage.success('订单已取消');
    await fetchOrderDetail(order.value.id);
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error);
      ElMessage.error('取消订单失败');
    }
  }
};

// 去支付
const payOrder = () => {
  router.push({
    path: '/payment',
    query: { 
      orderId: order.value.id,
      orderSn: order.value.orderSn,
      amount: order.value.payAmount
    }
  });
};

// 确认收货
const confirmReceive = async () => {
  try {
    await ElMessageBox.confirm('确认已收到商品吗？', '确认收货', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await confirmReceiveOrder(order.value.id);
    ElMessage.success('确认收货成功');
    await fetchOrderDetail(order.value.id);
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认收货失败:', error);
      ElMessage.error('确认收货失败');
    }
  }
};

// 评价商品
const viewComment = () => {
  router.push({
    path: '/user/comment',
    query: { orderId: order.value.id }
  });
};

// 删除订单
const handleDeleteOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要删除该订单吗？删除后不可恢复！', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await deleteOrderApi(order.value.id);
    ElMessage.success('订单已删除');
    back();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除订单失败:', error);
      ElMessage.error('删除订单失败');
    }
  }
};
</script>

<style scoped>
.order-detail-container {
  padding: 20px 0;
}

.order-card {
  margin-bottom: 20px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.order-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.info-item {
  display: flex;
  margin-bottom: 12px;
  font-size: 14px;
  line-height: 1.5;
}

.info-item .label {
  color: #909399;
  width: 100px;
  flex-shrink: 0;
}

.info-item .value {
  color: #303133;
  flex: 1;
}

.product-list {
  margin-bottom: 20px;
}

.product-item {
  display: flex;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}

.product-item:last-child {
  border-bottom: none;
}

.product-img {
  width: 80px;
  height: 80px;
  margin-right: 15px;
  overflow: hidden;
  border-radius: 4px;
  cursor: pointer;
}

.product-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  cursor: pointer;
}

.product-name:hover {
  color: var(--el-color-primary);
}

.product-sku {
  font-size: 12px;
  color: #909399;
}

.product-price,
.product-quantity,
.product-subtotal {
  width: 100px;
  text-align: center;
  font-size: 14px;
  color: #303133;
}

.product-subtotal {
  font-weight: bold;
}

.order-amount {
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
  padding-top: 15px;
}

.amount-item {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
  font-size: 14px;
}

.amount-item .label {
  color: #909399;
  margin-right: 15px;
}

.amount-item .value {
  width: 120px;
  text-align: right;
  color: #303133;
}

.amount-item.total {
  font-size: 16px;
  font-weight: bold;
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #ebeef5;
}

.amount-item .price {
  color: #f56c6c;
}

.logistics-timeline {
  margin-top: 20px;
  padding: 0 10px;
}

.no-logistics {
  margin: 30px 0;
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  gap: 10px;
}

.not-found {
  text-align: center;
  padding: 50px 0;
}

@media (max-width: 768px) {
  .order-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .order-status {
    margin-top: 10px;
  }
  
  .product-item {
    flex-wrap: wrap;
  }
  
  .product-img {
    width: 60px;
    height: 60px;
  }
  
  .product-info {
    width: calc(100% - 75px);
  }
  
  .product-price,
  .product-quantity,
  .product-subtotal {
    width: auto;
    margin-right: 15px;
    margin-top: 10px;
  }
  
  .product-subtotal {
    margin-right: 0;
  }
  
  .info-item {
    flex-wrap: wrap;
  }
  
  .info-item .label {
    width: 100%;
    margin-bottom: 5px;
  }
  
  .info-item .value {
    width: 100%;
  }
  
  .order-actions {
    flex-wrap: wrap;
  }
  
  .order-actions .el-button {
    margin-bottom: 10px;
  }
}
</style>
