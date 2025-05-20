<template>
  <div class="register-container">
    <h2>用户注册</h2>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="username">用户名</label>
        <input
          id="username"
          v-model="form.username"
          type="text"
          required
          placeholder="请输入用户名"
        />
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input
          id="password"
          v-model="form.password"
          type="password"
          required
          placeholder="请输入密码"
        />
      </div>
      <div class="form-group">
        <label for="confirmPassword">确认密码</label>
        <input
          id="confirmPassword"
          v-model="form.confirmPassword"
          type="password"
          required
          placeholder="请再次输入密码"
        />
      </div>
      <button type="submit" class="submit-btn">注册</button>
    </form>
  </div>
</template>

<script>
export default {
  name: 'Register',
  data() {
    return {
      form: {
        username: '',
        password: '',
        confirmPassword: ''
      }
    }
  },
  methods: {
    handleSubmit() {
      // 这里添加注册逻辑
      console.log('注册表单提交:', this.form)
    }
  }
}
</script>

<style scoped>
.register-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
}

input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.submit-btn {
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.submit-btn:hover {
  background-color: #3aa876;
}
</style>