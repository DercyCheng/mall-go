<template>
  <div class="product-card" @click="goToDetail">
    <div class="product-image">
      <img :src="product.pic" :alt="product.name" />
      <div class="product-tags" v-if="product.isNew || product.isHot">
        <span class="tag new" v-if="product.isNew">新品</span>
        <span class="tag hot" v-if="product.isHot">热卖</span>
      </div>
    </div>
    <div class="product-info">
      <h3 class="product-name">{{ product.name }}</h3>
      <p class="product-description">{{ product.description }}</p>
      <div class="product-price-row">
        <span class="product-price">¥{{ formatPrice(product.price) }}</span>
        <span class="product-original-price" v-if="product.originalPrice && product.originalPrice > product.price">
          ¥{{ formatPrice(product.originalPrice) }}
        </span>
      </div>
      <div class="product-footer">
        <div class="product-sales">已售 {{ product.sale || 0 }}</div>
        <el-button size="small" type="primary" @click.stop="handleAddToCart">加入购物车</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import { addToCart } from '@/api/cart';
import { ElMessage } from 'element-plus';

const router = useRouter();
const props = defineProps({
  product: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      name: '',
      pic: '',
      description: '',
      price: 0,
      originalPrice: 0,
      isNew: false,
      isHot: false,
      sale: 0
    })
  }
});

// 格式化价格
const formatPrice = (price) => {
  return Number(price).toFixed(2);
};

// 跳转到商品详情
const goToDetail = () => {
  router.push(`/product/${props.product.id}`);
};

// 添加到购物车
const handleAddToCart = async (event) => {
  // 阻止冒泡到父元素
  event.stopPropagation();
  
  try {
    await addToCart({
      productId: props.product.id,
      quantity: 1
    });
    
    ElMessage({
      message: '添加成功',
      type: 'success'
    });
  } catch (error) {
    ElMessage({
      message: '添加失败',
      type: 'error'
    });
  }
};
</script>

<style scoped>
.product-card {
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px 0 rgba(0,0,0,0.15);
}

.product-image {
  height: 200px;
  position: relative;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.product-card:hover .product-image img {
  transform: scale(1.1);
}

.product-tags {
  position: absolute;
  top: 10px;
  left: 10px;
  display: flex;
}

.tag {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  margin-right: 5px;
  color: #fff;
}

.tag.new {
  background-color: #409EFF;
}

.tag.hot {
  background-color: #F56C6C;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 16px;
  margin: 0 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #303133;
}

.product-description {
  font-size: 14px;
  color: #606266;
  margin: 0 0 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  height: 42px;
}

.product-price-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.product-price {
  font-size: 18px;
  font-weight: bold;
  color: #F56C6C;
  margin-right: 10px;
}

.product-original-price {
  font-size: 14px;
  color: #909399;
  text-decoration: line-through;
}

.product-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-sales {
  font-size: 12px;
  color: #909399;
}
</style>
