<template>
  <div class="user-menu">
    <el-menu
      :default-active="activeIndex"
      class="user-side-menu"
      :router="true"
    >
      <el-menu-item index="/user/profile">
        <el-icon><UserFilled /></el-icon>
        <span>个人资料</span>
      </el-menu-item>
      
      <el-menu-item index="/user/orders">
        <el-icon><List /></el-icon>
        <span>我的订单</span>
      </el-menu-item>
      
      <el-menu-item index="/user/favorites">
        <el-icon><Star /></el-icon>
        <span>我的收藏</span>
      </el-menu-item>
      
      <el-menu-item index="/user/address">
        <el-icon><Location /></el-icon>
        <span>收货地址</span>
      </el-menu-item>
      
      <el-menu-item index="/user/coupons" disabled>
        <el-icon><Discount /></el-icon>
        <span>我的优惠券</span>
        <el-tag size="small" type="info">待开发</el-tag>
      </el-menu-item>
      
      <el-menu-item index="/user/messages" disabled>
        <el-icon><ChatDotRound /></el-icon>
        <span>消息中心</span>
        <el-tag size="small" type="info">待开发</el-tag>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { UserFilled, List, Star, Location, Discount, ChatDotRound } from '@element-plus/icons-vue';

const route = useRoute();

// Compute active menu item based on current route
const activeIndex = computed(() => {
  return route.path;
});
</script>

<style scoped>
.user-menu {
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.user-side-menu {
  border-right: none;
}

.el-menu-item {
  display: flex;
  align-items: center;
}

.el-icon {
  margin-right: 8px;
}

@media (max-width: 768px) {
  .user-menu {
    margin-bottom: 15px;
  }
  
  .user-side-menu {
    display: flex;
    flex-wrap: nowrap;
    overflow-x: auto;
    white-space: nowrap;
    border-bottom: 1px solid #e6e6e6;
  }
  
  .el-menu-item {
    border-right: 1px solid #e6e6e6;
    flex-shrink: 0;
  }
}
</style>
