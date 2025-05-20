<script setup>
// No imports needed at the root level
</script>

<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

.app-container {
  min-height: 100vh;
  width: 100%;
}
</style>
