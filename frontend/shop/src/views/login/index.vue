<template>
  <div class="login-container">
    <div class="login-form-container">
      <div class="login-header">
        <h2 class="login-title">{{ isLogin ? '用户登录' : '用户注册' }}</h2>
        <div class="login-tabs">
          <span 
            :class="['tab-item', { active: isLogin }]" 
            @click="isLogin = true"
          >
            登录
          </span>
          <span 
            :class="['tab-item', { active: !isLogin }]" 
            @click="isLogin = false"
          >
            注册
          </span>
        </div>
      </div>
      
      <el-form 
        ref="formRef" 
        :model="formData" 
        :rules="rules" 
        label-position="top"
        class="login-form"
      >
        <!-- 登录表单 -->
        <template v-if="isLogin">
          <el-form-item prop="username" label="用户名">
            <el-input 
              v-model="formData.username" 
              placeholder="请输入用户名" 
              :prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="password" label="密码">
            <el-input 
              v-model="formData.password" 
              type="password" 
              placeholder="请输入密码" 
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <div class="form-actions">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-link type="primary">忘记密码?</el-link>
          </div>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn" 
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </template>
        
        <!-- 注册表单 -->
        <template v-else>
          <el-form-item prop="username" label="用户名">
            <el-input 
              v-model="formData.username" 
              placeholder="请输入用户名" 
              :prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item prop="password" label="密码">
            <el-input 
              v-model="formData.password" 
              type="password" 
              placeholder="请输入密码" 
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="confirmPassword" label="确认密码">
            <el-input 
              v-model="formData.confirmPassword" 
              type="password" 
              placeholder="请再次输入密码" 
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item prop="phone" label="手机号">
            <el-input 
              v-model="formData.phone" 
              placeholder="请输入手机号" 
              :prefix-icon="Iphone"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="submit-btn" 
              :loading="loading"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
        </template>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { User, Lock, Iphone } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import { login, register } from '@/api/user';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const formRef = ref(null);
const isLogin = ref(true);
const loading = ref(false);
const rememberMe = ref(false);

// 表单数据
const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  phone: ''
});

// 验证规则
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== formData.password) {
    callback(new Error('两次输入的密码不一致'));
  } else {
    callback();
  }
};

const rules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ]
});

// 处理登录
const handleLogin = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    loading.value = true;
    
    const loginData = {
      username: formData.username,
      password: formData.password
    };
    
    await userStore.login(loginData);
    
    // 如果有重定向地址则跳转到该地址，否则跳转到首页
    const redirect = route.query.redirect || '/';
    router.replace(redirect);
    
    ElMessage({
      message: '登录成功',
      type: 'success'
    });
  } catch (error) {
    console.error('登录失败:', error);
    ElMessage({
      message: error.message || '登录失败，请检查用户名和密码',
      type: 'error'
    });
  } finally {
    loading.value = false;
  }
};

// 处理注册
const handleRegister = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    loading.value = true;
    
    const registerData = {
      username: formData.username,
      password: formData.password,
      phone: formData.phone
    };
    
    await register(registerData);
    
    ElMessage({
      message: '注册成功，请登录',
      type: 'success'
    });
    
    // 注册成功后切换到登录页
    isLogin.value = true;
  } catch (error) {
    console.error('注册失败:', error);
    ElMessage({
      message: error.message || '注册失败，请稍后再试',
      type: 'error'
    });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  padding: 50px 0;
}

.login-form-container {
  width: 420px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
}

.login-tabs {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.tab-item {
  padding: 10px 20px;
  margin: 0 10px;
  cursor: pointer;
  color: #606266;
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.tab-item.active {
  color: #409EFF;
  border-bottom-color: #409EFF;
}

.login-form {
  margin-top: 20px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.submit-btn {
  width: 100%;
  padding: 12px 0;
  font-size: 16px;
}

@media (max-width: 576px) {
  .login-form-container {
    width: 90%;
    padding: 20px;
  }
}
</style>
