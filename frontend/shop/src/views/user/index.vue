<template>
  <div class="user-container">
    <div class="container">
      <el-row :gutter="20">
        <!-- 侧边栏 -->
        <el-col :xs="24" :sm="8" :md="6" :lg="5">
          <div class="sidebar">
            <div class="user-info">
              <el-avatar :size="80" :src="userInfo.avatar || defaultAvatar"></el-avatar>
              <div class="username">{{ userInfo.username }}</div>
              <div class="member-level">
                <el-tag type="warning" effect="dark">
                  {{ userInfo.memberLevelName || '普通会员' }}
                </el-tag>
              </div>
            </div>
            
            <!-- User Menu Component -->
            <UserMenu />
          </div>
        </el-col>
        
        <!-- 主内容 -->
        <el-col :xs="24" :sm="16" :md="18" :lg="19">
          <div class="content-container">
            <router-view />
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useUserStore } from '@/store/user';
import { getUserInfo } from '@/api/user';
import UserMenu from '@/components/UserMenu.vue';

const route = useRoute();
const userStore = useUserStore();

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png';
const userInfo = ref({});
const loading = ref(false);

onMounted(async () => {
  await loadUserInfo();
});

// 加载用户信息
const loadUserInfo = async () => {
  try {
    loading.value = true;
    
    // 如果已有用户信息，直接使用
    if (userStore.userInfo && userStore.userInfo.id) {
      userInfo.value = userStore.userInfo;
      return;
    }
    
    // 否则从API获取
    const result = await userStore.getUserInfo();
    userInfo.value = result || {};
  } catch (error) {
    console.error('加载用户信息失败:', error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.user-container {
  background-color: #f5f7fa;
  padding: 20px 0 40px;
  min-height: calc(100vh - 200px);
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* 侧边栏 */
.sidebar {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.user-info {
  padding: 30px 20px;
  text-align: center;
  background: linear-gradient(to right, #40a9ff, #409EFF);
  color: #fff;
}

.username {
  margin: 15px 0 10px;
  font-size: 18px;
  font-weight: bold;
}

.member-level {
  margin-bottom: 10px;
}

.nav-menu {
  padding: 10px 0;
}

.user-menu {
  border-right: none;
}

/* 主内容 */
.content-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
  min-height: 600px;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .sidebar {
    margin-bottom: 20px;
  }
  
  .user-info {
    padding: 20px;
  }
}
</style>
