<template>
  <div class="checkout-container">
    <div class="container">
      <div class="checkout-header">
        <h1 class="checkout-title">订单结算</h1>
      </div>
      
      <div class="checkout-content" v-loading="loading">
        <!-- 收货地址 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">收货地址</h2>
            <el-button type="primary" link @click="showAddressDialog">
              {{ addresses.length > 0 ? '管理收货地址' : '添加收货地址' }}
            </el-button>
          </div>
          
          <div class="section-content">
            <div v-if="addresses.length === 0" class="no-address">
              <el-empty description="暂无收货地址，请添加">
                <el-button type="primary" @click="showAddressDialog">添加地址</el-button>
              </el-empty>
            </div>
            
            <div v-else class="address-list">
              <div 
                v-for="address in addresses" 
                :key="address.id"
                :class="['address-item', { active: address.id === selectedAddressId }]"
                @click="selectAddress(address.id)"
              >
                <div class="address-info">
                  <div class="address-name">
                    <span class="name">{{ address.name }}</span>
                    <span class="phone">{{ address.phone }}</span>
                  </div>
                  <div class="address-detail">
                    <span class="tag" v-if="address.defaultStatus">默认</span>
                    {{ formatAddress(address) }}
                  </div>
                </div>
                <div class="address-check" v-if="address.id === selectedAddressId">
                  <el-icon><Check /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 支付方式 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">支付方式</h2>
          </div>
          
          <div class="section-content">
            <div class="payment-list">
              <div 
                v-for="payment in paymentMethods" 
                :key="payment.id"
                :class="['payment-item', { active: payment.id === selectedPaymentId }]"
                @click="selectPayment(payment.id)"
              >
                <div class="payment-info">
                  <img :src="payment.icon" :alt="payment.name" class="payment-icon" />
                  <span class="payment-name">{{ payment.name }}</span>
                </div>
                <div class="payment-check" v-if="payment.id === selectedPaymentId">
                  <el-icon><Check /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 配送方式 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">配送方式</h2>
          </div>
          
          <div class="section-content">
            <div class="delivery-list">
              <div 
                v-for="delivery in deliveryMethods" 
                :key="delivery.id"
                :class="['delivery-item', { active: delivery.id === selectedDeliveryId }]"
                @click="selectDelivery(delivery.id)"
              >
                <div class="delivery-info">
                  <span class="delivery-name">{{ delivery.name }}</span>
                  <span class="delivery-price">¥{{ formatPrice(delivery.price) }}</span>
                </div>
                <div class="delivery-desc">
                  {{ delivery.desc }}
                </div>
                <div class="delivery-check" v-if="delivery.id === selectedDeliveryId">
                  <el-icon><Check /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 商品清单 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">商品清单</h2>
            <router-link to="/cart" class="back-to-cart">返回购物车</router-link>
          </div>
          
          <div class="section-content">
            <div class="product-list">
              <el-table :data="orderItems" style="width: 100%">
                <el-table-column label="商品信息" min-width="400">
                  <template #default="scope">
                    <div class="product-info">
                      <div class="product-image">
                        <img :src="scope.row.productPic" :alt="scope.row.productName" />
                      </div>
                      <div class="product-details">
                        <div class="product-name">{{ scope.row.productName }}</div>
                        <div class="product-attrs" v-if="scope.row.productAttr">
                          {{ scope.row.productAttr }}
                        </div>
                      </div>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column prop="price" label="单价" width="120">
                  <template #default="scope">
                    ¥{{ formatPrice(scope.row.price) }}
                  </template>
                </el-table-column>
                
                <el-table-column prop="quantity" label="数量" width="100" />
                
                <el-table-column label="小计" width="120">
                  <template #default="scope">
                    ¥{{ formatPrice(scope.row.price * scope.row.quantity) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </div>
        
        <!-- 订单备注 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">订单备注</h2>
          </div>
          
          <div class="section-content">
            <el-input
              v-model="orderNote"
              type="textarea"
              placeholder="请输入订单备注，选填"
              :rows="3"
              maxlength="200"
              show-word-limit
            />
          </div>
        </div>
        
        <!-- 结算信息 -->
        <div class="checkout-section">
          <div class="section-header">
            <h2 class="section-title">结算信息</h2>
          </div>
          
          <div class="section-content">
            <div class="order-summary">
              <div class="summary-item">
                <span class="item-label">商品总额：</span>
                <span class="item-value">¥{{ formatPrice(productTotal) }}</span>
              </div>
              <div class="summary-item">
                <span class="item-label">运费：</span>
                <span class="item-value">¥{{ formatPrice(deliveryFee) }}</span>
              </div>
              <div class="summary-item">
                <span class="item-label">优惠：</span>
                <span class="item-value discount">-¥{{ formatPrice(discount) }}</span>
              </div>
              <div class="summary-divider"></div>
              <div class="summary-item total">
                <span class="item-label">应付总额：</span>
                <span class="item-value">¥{{ formatPrice(orderTotal) }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 提交订单 -->
        <div class="checkout-footer">
          <div class="footer-summary">
            <span class="summary-count">
              共 <span class="highlight">{{ totalQuantity }}</span> 件商品，
              合计（含运费）：
              <span class="highlight total-price">¥{{ formatPrice(orderTotal) }}</span>
            </span>
          </div>
          
          <el-button 
            type="danger" 
            size="large" 
            :disabled="!canSubmit" 
            :loading="submitting"
            @click="handleSubmitOrder"
          >
            提交订单
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 地址管理对话框 -->
    <el-dialog
      v-model="addressDialogVisible"
      title="管理收货地址"
      width="600px"
    >
      <address-manager
        :addresses="addresses"
        @save-address="handleSaveAddress"
        @delete-address="handleDeleteAddress"
        @set-default="handleSetDefault"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Check } from '@element-plus/icons-vue';
import { generateConfirmOrder, submitOrder as submitOrderApi } from '@/api/order';
import { useUserStore } from '@/store/user';
import AddressManager from '@/components/AddressManager.vue';

const router = useRouter();
const userStore = useUserStore();

const loading = ref(false);
const submitting = ref(false);
const addresses = ref([]);
const orderItems = ref([]);
const selectedAddressId = ref('');
const selectedPaymentId = ref('1'); // 默认支付方式：在线支付
const selectedDeliveryId = ref('1'); // 默认配送方式：普通配送
const orderNote = ref('');
const addressDialogVisible = ref(false);
const discount = ref(0);

// 支付方式
const paymentMethods = ref([
  {
    id: '1',
    name: '在线支付',
    icon: 'https://img0.baidu.com/it/u=1085102667,3823302930&fm=253&fmt=auto&app=138&f=PNG',
  },
  {
    id: '2',
    name: '货到付款',
    icon: 'https://img0.baidu.com/it/u=1085102667,3823302930&fm=253&fmt=auto&app=138&f=PNG',
  }
]);

// 配送方式
const deliveryMethods = ref([
  {
    id: '1',
    name: '普通配送',
    price: 0,
    desc: '订单满88元包邮，预计1-3天送达'
  },
  {
    id: '2',
    name: '次日达',
    price: 15,
    desc: '工作日15点前下单，可次日送达'
  },
  {
    id: '3',
    name: '加急配送',
    price: 30,
    desc: '2小时内送达，仅限市区支持'
  }
]);

// 商品总价
const productTotal = computed(() => {
  return orderItems.value.reduce((sum, item) => {
    return sum + (item.price * item.quantity);
  }, 0);
});

// 运费
const deliveryFee = computed(() => {
  const delivery = deliveryMethods.value.find(item => item.id === selectedDeliveryId.value);
  
  // 满88包邮（仅限普通配送）
  if (selectedDeliveryId.value === '1' && productTotal.value >= 88) {
    return 0;
  }
  
  return delivery ? delivery.price : 0;
});

// 订单总价
const orderTotal = computed(() => {
  return productTotal.value + deliveryFee.value - discount.value;
});

// 商品总数量
const totalQuantity = computed(() => {
  return orderItems.value.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);
});

// 是否可以提交订单
const canSubmit = computed(() => {
  return selectedAddressId.value && orderItems.value.length > 0;
});

// 初始化
onMounted(async () => {
  // 检查用户是否登录
  if (!userStore.isLoggedIn) {
    router.push('/login?redirect=/checkout');
    return;
  }
  
  await loadOrderData();
});

// 加载订单数据
const loadOrderData = async () => {
  try {
    loading.value = true;
    
    const result = await generateConfirmOrder();
    
    // 设置地址列表
    addresses.value = result.addressList || [];
    if (addresses.value.length > 0) {
      // 优先选择默认地址
      const defaultAddress = addresses.value.find(address => address.defaultStatus);
      selectedAddressId.value = defaultAddress ? defaultAddress.id : addresses.value[0].id;
    }
    
    // 设置订单商品
    orderItems.value = result.cartItemList || [];
    
    // 设置优惠金额
    discount.value = result.discount || 0;
  } catch (error) {
    console.error('加载订单数据失败:', error);
    ElMessage.error('加载订单数据失败');
  } finally {
    loading.value = false;
  }
};

// 格式化价格
const formatPrice = (price) => {
  return Number(price).toFixed(2);
};

// 格式化地址
const formatAddress = (address) => {
  return `${address.province} ${address.city} ${address.region} ${address.detailAddress}`;
};

// 选择地址
const selectAddress = (addressId) => {
  selectedAddressId.value = addressId;
};

// 选择支付方式
const selectPayment = (paymentId) => {
  selectedPaymentId.value = paymentId;
};

// 选择配送方式
const selectDelivery = (deliveryId) => {
  selectedDeliveryId.value = deliveryId;
};

// 显示地址管理对话框
const showAddressDialog = () => {
  addressDialogVisible.value = true;
};

// 保存地址
const handleSaveAddress = async (address) => {
  try {
    // 这里应该调用保存地址的API
    // 模拟API调用
    if (address.id) {
      // 更新地址
      const index = addresses.value.findIndex(item => item.id === address.id);
      if (index !== -1) {
        addresses.value[index] = address;
      }
    } else {
      // 新增地址
      address.id = Date.now().toString();
      addresses.value.push(address);
    }
    
    // 如果是默认地址，更新其他地址的默认状态
    if (address.defaultStatus) {
      addresses.value.forEach(item => {
        if (item.id !== address.id) {
          item.defaultStatus = false;
        }
      });
      
      // 自动选择默认地址
      selectedAddressId.value = address.id;
    }
    
    ElMessage.success('保存地址成功');
  } catch (error) {
    console.error('保存地址失败:', error);
    ElMessage.error('保存地址失败');
  }
};

// 删除地址
const handleDeleteAddress = async (addressId) => {
  try {
    // 这里应该调用删除地址的API
    // 模拟API调用
    addresses.value = addresses.value.filter(item => item.id !== addressId);
    
    // 如果删除的是当前选中的地址，重新选择地址
    if (selectedAddressId.value === addressId) {
      selectedAddressId.value = addresses.value.length > 0 ? addresses.value[0].id : '';
    }
    
    ElMessage.success('删除地址成功');
  } catch (error) {
    console.error('删除地址失败:', error);
    ElMessage.error('删除地址失败');
  }
};

// 设置默认地址
const handleSetDefault = async (addressId) => {
  try {
    // 这里应该调用设置默认地址的API
    // 模拟API调用
    addresses.value.forEach(item => {
      item.defaultStatus = item.id === addressId;
    });
    
    // 自动选择默认地址
    selectedAddressId.value = addressId;
    
    ElMessage.success('设置默认地址成功');
  } catch (error) {
    console.error('设置默认地址失败:', error);
    ElMessage.error('设置默认地址失败');
  }
};

// 提交订单
const handleSubmitOrder = async () => {
  if (!canSubmit.value) {
    return;
  }
  
  try {
    await ElMessageBox.confirm('确定要提交订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    submitting.value = true;
    
    const selectedAddress = addresses.value.find(address => address.id === selectedAddressId.value);
    
    const orderData = {
      memberReceiveAddressId: selectedAddressId.value,
      payType: selectedPaymentId.value,
      deliveryMethod: selectedDeliveryId.value,
      note: orderNote.value,    orderItemList: orderItems.value.map(item => ({
      productId: item.productId,
      quantity: item.quantity,
      productAttr: item.productAttr
    }))
  };
  
  const result = await submitOrderApi(orderData);
  
  ElMessage.success('订单提交成功');
    
    // 跳转到支付页面或订单详情页
    router.push(`/payment?orderId=${result.orderId}&orderSn=${result.orderSn}`);
  } catch (error) {
    if (error !== 'cancel') {
      console.error('提交订单失败:', error);
      ElMessage.error('提交订单失败');
    }
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
.checkout-container {
  background-color: #f5f7fa;
  padding: 20px 0 40px;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* 结算头部 */
.checkout-header {
  margin-bottom: 20px;
}

.checkout-title {
  font-size: 24px;
  color: #303133;
  margin: 0;
}

/* 结算内容 */
.checkout-content {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

/* 结算区块通用样式 */
.checkout-section {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
}

.checkout-section:last-child {
  border-bottom: none;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  color: #303133;
  margin: 0;
}

.back-to-cart {
  color: #409EFF;
  text-decoration: none;
}

/* 地址列表 */
.no-address {
  padding: 30px 0;
  text-align: center;
}

.address-list {
  display: flex;
  flex-wrap: wrap;
}

.address-item {
  width: calc(33.333% - 15px);
  margin-right: 15px;
  margin-bottom: 15px;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  cursor: pointer;
  position: relative;
  transition: all 0.3s;
}

.address-item:hover {
  border-color: #409EFF;
}

.address-item.active {
  border-color: #409EFF;
  background-color: #ecf5ff;
}

.address-info {
  line-height: 1.5;
}

.address-name {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.address-name .phone {
  margin-left: 15px;
  font-weight: normal;
  color: #606266;
}

.address-detail {
  color: #606266;
}

.address-detail .tag {
  display: inline-block;
  padding: 2px 6px;
  background-color: #409EFF;
  color: #fff;
  font-size: 12px;
  border-radius: 2px;
  margin-right: 6px;
}

.address-check {
  position: absolute;
  bottom: 10px;
  right: 10px;
  color: #409EFF;
}

/* 支付方式 */
.payment-list {
  display: flex;
}

.payment-item {
  width: 200px;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 20px;
  position: relative;
  transition: all 0.3s;
}

.payment-item:hover {
  border-color: #409EFF;
}

.payment-item.active {
  border-color: #409EFF;
  background-color: #ecf5ff;
}

.payment-info {
  display: flex;
  align-items: center;
}

.payment-icon {
  width: 30px;
  height: 30px;
  margin-right: 10px;
}

.payment-name {
  color: #303133;
}

.payment-check {
  position: absolute;
  bottom: 10px;
  right: 10px;
  color: #409EFF;
}

/* 配送方式 */
.delivery-list {
  display: flex;
  flex-wrap: wrap;
}

.delivery-item {
  width: 250px;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 20px;
  margin-bottom: 15px;
  position: relative;
  transition: all 0.3s;
}

.delivery-item:hover {
  border-color: #409EFF;
}

.delivery-item.active {
  border-color: #409EFF;
  background-color: #ecf5ff;
}

.delivery-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.delivery-name {
  color: #303133;
  font-weight: bold;
}

.delivery-price {
  color: #F56C6C;
}

.delivery-desc {
  color: #909399;
  font-size: 14px;
}

.delivery-check {
  position: absolute;
  bottom: 10px;
  right: 10px;
  color: #409EFF;
}

/* 商品列表 */
.product-info {
  display: flex;
  align-items: center;
}

.product-image {
  width: 60px;
  height: 60px;
  margin-right: 15px;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.product-name {
  margin-bottom: 5px;
}

.product-attrs {
  font-size: 12px;
  color: #909399;
}

/* 结算信息 */
.order-summary {
  width: 350px;
  margin-left: auto;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

.item-label {
  color: #606266;
}

.item-value {
  color: #303133;
}

.item-value.discount {
  color: #F56C6C;
}

.summary-divider {
  height: 1px;
  background-color: #ebeef5;
  margin: 15px 0;
}

.summary-item.total {
  font-size: 18px;
  font-weight: bold;
}

.summary-item.total .item-value {
  color: #F56C6C;
}

/* 结算底部 */
.checkout-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-top: 1px solid #ebeef5;
  border-radius: 0 0 8px 8px;
}

.footer-summary {
  margin-right: 20px;
  font-size: 14px;
  color: #606266;
}

.highlight {
  color: #F56C6C;
  font-weight: bold;
}

.total-price {
  font-size: 20px;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 992px) {
  .address-item {
    width: calc(50% - 15px);
  }
  
  .payment-item, .delivery-item {
    width: calc(50% - 15px);
  }
}

@media (max-width: 768px) {
  .address-item {
    width: 100%;
    margin-right: 0;
  }
  
  .payment-item, .delivery-item {
    width: 100%;
    margin-right: 0;
  }
  
  .order-summary {
    width: 100%;
  }
  
  .checkout-footer {
    flex-direction: column;
  }
  
  .footer-summary {
    margin-right: 0;
    margin-bottom: 15px;
  }
}
</style>
