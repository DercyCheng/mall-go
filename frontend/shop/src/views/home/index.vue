<template>
  <div class="home-container">
    <div class="container">
      <!-- Banner -->
      <div class="banner-section">
        <el-carousel height="400px">
          <el-carousel-item v-for="banner in banners" :key="banner.id">
            <div class="banner-item" @click="handleBannerClick(banner)">
              <img :src="banner.img" :alt="banner.title" class="banner-img" />
            </div>
          </el-carousel-item>
        </el-carousel>
      </div>

      <!-- Feature section -->
      <div class="feature-section">
        <div class="feature-item" v-for="feature in features" :key="feature.id">
          <el-icon :size="24" class="feature-icon">
            <component :is="feature.icon"></component>
          </el-icon>
          <div class="feature-info">
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.desc }}</p>
          </div>
        </div>
      </div>

      <!-- New arrivals -->
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">新品上市</h2>
          <router-link to="/category?sort=newest" class="view-more">
            查看更多 <el-icon><ArrowRight /></el-icon>
          </router-link>
        </div>
        <div class="product-grid">
          <product-card
            v-for="product in newProducts"
            :key="product.id"
            :product="product"
          ></product-card>
        </div>
      </div>

      <!-- Hot products -->
      <div class="section-container">
        <div class="section-header">
          <h2 class="section-title">热卖商品</h2>
          <router-link to="/category?sort=sales" class="view-more">
            查看更多 <el-icon><ArrowRight /></el-icon>
          </router-link>
        </div>
        <div class="product-grid">
          <product-card
            v-for="product in hotProducts"
            :key="product.id"
            :product="product"
          ></product-card>
        </div>
      </div>

      <!-- Recommendations -->
      <div class="section-container">
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getProductList } from '@/api/product';
import ProductCard from '@/components/ProductCard.vue';
import { ArrowRight, Ship, GoodsFilled, Collection, PriceTag } from '@element-plus/icons-vue';

const router = useRouter();
const banners = ref([]);
const newProducts = ref([]);
const hotProducts = ref([]);
const recommendProducts = ref([]);

// 服务特性
const features = ref([
  { 
    id: 1, 
    icon: 'Ship', 
    title: '全场包邮', 
    desc: '所有商品均免运费' 
  },
  { 
    id: 2, 
    icon: 'GoodsFilled', 
    title: '正品保障', 
    desc: '所有商品均为正品' 
  },
  { 
    id: 3, 
    icon: 'Collection', 
    title: '售后无忧', 
    desc: '7天无理由退换货' 
  },
  { 
    id: 4, 
    icon: 'PriceTag', 
    title: '优惠活动', 
    desc: '多种优惠活动进行中' 
  }
]);

onMounted(async () => {
  try {
    // 加载轮播图数据
    loadBanners();
    
    // 加载新品上市
    const newRes = await getProductList({ pageSize: 4, sort: 'newest' });
    newProducts.value = newRes?.list || [];
    
    // 加载热卖商品
    const hotRes = await getProductList({ pageSize: 4, sort: 'sales' });
    hotProducts.value = hotRes?.list || [];
    
    // 加载推荐商品
    const recRes = await getProductList({ pageSize: 8, sort: 'recommend' });
    recommendProducts.value = recRes?.list || [];
  } catch (error) {
    console.error('加载首页数据失败:', error);
  }
});

// 加载轮播图数据
const loadBanners = () => {
  // 这里可以替换为实际的 API 调用
  banners.value = [
    {
      id: 1,
      title: '新款手机发布',
      img: 'https://img0.baidu.com/it/u=2729211458,3797562168&fm=253&fmt=auto&app=120&f=JPEG',
      link: '/category?categoryId=1'
    },
    {
      id: 2,
      title: '时尚服饰上新',
      img: 'https://img0.baidu.com/it/u=2866115257,1431158626&fm=253&fmt=auto&app=138&f=JPEG',
      link: '/category?categoryId=2'
    },
    {
      id: 3,
      title: '生活家居特惠',
      img: 'https://img1.baidu.com/it/u=3019431950,4010520&fm=253&fmt=auto&app=120&f=JPEG',
      link: '/category?categoryId=3'
    }
  ];
};

// 处理轮播图点击
const handleBannerClick = (banner) => {
  router.push(banner.link);
};
</script>

<style scoped>
.home-container {
  padding: 0 0 40px;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* Banner styles */
.banner-section {
  margin-bottom: 30px;
}

.banner-item {
  height: 100%;
  cursor: pointer;
}

.banner-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

/* Feature section */
.feature-section {
  display: flex;
  justify-content: space-between;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.feature-item {
  display: flex;
  align-items: center;
  width: 25%;
  padding: 0 10px;
}

.feature-icon {
  color: #409EFF;
  margin-right: 10px;
}

.feature-info h3 {
  font-size: 16px;
  margin: 0 0 5px;
  color: #303133;
}

.feature-info p {
  font-size: 14px;
  margin: 0;
  color: #606266;
}

/* Section styles */
.section-container {
  margin-bottom: 40px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 22px;
  color: #303133;
  margin: 0;
  position: relative;
  padding-left: 15px;
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

.view-more {
  display: flex;
  align-items: center;
  color: #409EFF;
  text-decoration: none;
  font-size: 14px;
}

.view-more .el-icon {
  margin-left: 5px;
}

/* Product grid */
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
  
  .feature-section {
    flex-wrap: wrap;
  }
  
  .feature-item {
    width: 50%;
    margin-bottom: 15px;
  }
}

@media (max-width: 768px) {
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .feature-item {
    width: 100%;
  }
}
</style>
