<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <img src="/images/404.svg" alt="Page Not Found" class="not-found-image" />
      <h1 class="not-found-title">404</h1>
      <h2 class="not-found-subtitle">Page Not Found</h2>
      <p class="not-found-message">
        Sorry, the page you are looking for does not exist or has been moved.
      </p>
      <el-button type="primary" @click="goHome" class="not-found-button">
        Back to Home
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f9f9f9;
  padding: 20px;
}

.not-found-content {
  text-align: center;
  max-width: 600px;
}

.not-found-image {
  width: 250px;
  margin-bottom: 20px;
}

.not-found-title {
  font-size: 80px;
  font-weight: bold;
  color: var(--el-color-primary);
  margin: 0;
  line-height: 1;
}

.not-found-subtitle {
  font-size: 30px;
  color: #5a5a5a;
  margin: 10px 0;
}

.not-found-message {
  color: #8c8c8c;
  margin-bottom: 30px;
  font-size: 16px;
}

.not-found-button {
  padding: 12px 30px;
  font-size: 16px;
}

@media (max-width: 768px) {
  .not-found-image {
    width: 180px;
  }

  .not-found-title {
    font-size: 60px;
  }

  .not-found-subtitle {
    font-size: 24px;
  }
}
</style>
