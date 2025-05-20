<template>
  <div class="layout-container">
    <!-- Header -->
    <header class="header">
      <div class="container">
        <div class="header-left">
          <div class="logo" @click="goToHome">
            <img src="/vite.svg" alt="商城" class="logo-img" />
            <span class="logo-text">微服务商城</span>
          </div>
          <div class="search-box">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索商品"
              class="search-input"
              @keyup.enter="onSearch"
            >
              <template #append>
                <el-button @click="onSearch">
                  <el-icon><Search /></el-icon>
                </el-button>
              </template>
            </el-input>
          </div>
        </div>
        <div class="header-right">
          <div class="nav-links">
            <router-link to="/cart" class="nav-item">
              <el-badge :value="cartCount" :hidden="cartCount <= 0">
                <el-icon :size="20"><ShoppingCart /></el-icon>
              </el-badge>
              <span>购物车</span>
            </router-link>
            <template v-if="userStore.isLoggedIn">
              <router-link to="/user" class="nav-item">
                <el-avatar :size="24" :src="userAvatar"></el-avatar>
                <span>{{ userStore.username }}</span>
              </router-link>
              <span class="nav-item logout" @click="logout">退出</span>
            </template>
            <template v-else>
              <router-link to="/login" class="nav-item">登录</router-link>
              <router-link to="/register" class="nav-item">注册</router-link>
            </template>
          </div>
        </div>
      </div>
    </header>

    <!-- Category Menu -->
    <div class="category-menu">
      <div class="container">
        <div class="categories">
          <router-link
            v-for="category in categories"
            :key="category.id"
            :to="`/category?id=${category.id}`"
            class="category-item"
          >
            {{ category.name }}
          </router-link>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <main class="main-content">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h3>关于我们</h3>
            <p>微服务商城是一个基于微服务架构的现代电商平台</p>
          </div>
          <div class="footer-section">
            <h3>客户服务</h3>
            <ul>
              <li>帮助中心</li>
              <li>联系我们</li>
              <li>配送信息</li>
              <li>退换政策</li>
            </ul>
          </div>
          <div class="footer-section">
            <h3>商城服务</h3>
            <ul>
              <li>支付方式</li>
              <li>配送方式</li>
              <li>服务保证</li>
              <li>常见问题</li>
            </ul>
          </div>
          <div class="footer-section">
            <h3>关注我们</h3>
            <div class="social-icons">
              <span class="social-icon"><el-icon><ChatDotRound /></el-icon></span>
              <span class="social-icon"><el-icon><Promotion /></el-icon></span>
              <span class="social-icon"><el-icon><Place /></el-icon></span>
            </div>
          </div>
        </div>
        <div class="copyright">
          <p>&copy; 2025 微服务商城 版权所有</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { getCategories } from '@/api/product';
import { getCartList } from '@/api/cart';
import { Search, ShoppingCart, ChatDotRound, Promotion, Place } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const searchKeyword = ref('');
const categories = ref([]);
const cartCount = ref(0);

const userAvatar = computed(() => {
  return userStore.userInfo?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png';
});

onMounted(async () => {
  try {
    // 获取商品分类
    const categoryData = await getCategories();
    categories.value = categoryData || [];
    
    // 如果用户已登录，获取购物车数量和用户信息
    if (userStore.isLoggedIn) {
      loadUserData();
    }
  } catch (error) {
    console.error('初始化数据失败:', error);
  }
});

// 加载用户数据
const loadUserData = async () => {
  try {
    // 获取用户信息
    if (!userStore.userInfo || !userStore.userInfo.id) {
      await userStore.getUserInfo();
    }
    
    // 获取购物车数量
    const cartData = await getCartList();
    cartCount.value = cartData?.length || 0;
  } catch (error) {
    console.error('加载用户数据失败:', error);
  }
};

// 搜索商品
const onSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({
      path: '/category',
      query: {
        keyword: searchKeyword.value
      }
    });
  }
};

// 返回首页
const goToHome = () => {
  router.push('/');
};

// 退出登录
const logout = () => {
  userStore.logout();
  router.push('/login');
};
</script>

<style scoped>
.layout-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* Header styles */
.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 15px 0;
}

.header .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 30px;
  cursor: pointer;
}

.logo-img {
  height: 40px;
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: #409EFF;
}

.search-box {
  width: 400px;
}

.header-right {
  display: flex;
  align-items: center;
}

.nav-links {
  display: flex;
  align-items: center;
}

.nav-item {
  display: flex;
  align-items: center;
  margin-left: 20px;
  color: #606266;
  text-decoration: none;
  cursor: pointer;
}

.nav-item span {
  margin-left: 5px;
}

.nav-item:hover {
  color: #409EFF;
}

.logout {
  color: #F56C6C;
}

/* Category menu styles */
.category-menu {
  background-color: #f5f7fa;
  border-bottom: 1px solid #e4e7ed;
}

.categories {
  display: flex;
  padding: 12px 0;
  overflow-x: auto;
}

.category-item {
  padding: 0 15px;
  color: #303133;
  text-decoration: none;
  white-space: nowrap;
}

.category-item:hover {
  color: #409EFF;
}

/* Main content styles */
.main-content {
  flex: 1;
  padding: 20px 0;
  background-color: #f5f7fa;
}

/* Footer styles */
.footer {
  background-color: #303133;
  color: #c0c4cc;
  padding: 40px 0 20px;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
}

.footer-section {
  width: 25%;
  padding: 0 15px;
  margin-bottom: 20px;
}

.footer-section h3 {
  color: #fff;
  margin-bottom: 15px;
  font-size: 16px;
}

.footer-section ul {
  list-style: none;
  padding: 0;
}

.footer-section li {
  margin-bottom: 8px;
  cursor: pointer;
}

.footer-section li:hover {
  color: #fff;
}

.social-icons {
  display: flex;
}

.social-icon {
  background-color: #767a82;
  color: #303133;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  cursor: pointer;
}

.social-icon:hover {
  background-color: #fff;
}

.copyright {
  border-top: 1px solid #424242;
  padding-top: 20px;
  margin-top: 20px;
  text-align: center;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .footer-section {
    width: 50%;
  }
  
  .search-box {
    width: 200px;
  }
}

@media (max-width: 576px) {
  .footer-section {
    width: 100%;
  }
  
  .header-left {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .search-box {
    margin-top: 10px;
    width: 100%;
  }
}
</style>
