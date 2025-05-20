<template>
  <div class="profile-container">
    <div class="profile-header">
      <h2 class="section-title">个人资料</h2>
    </div>
    
    <div class="profile-content" v-loading="loading">
      <el-form 
        ref="formRef" 
        :model="userForm" 
        :rules="rules" 
        label-width="100px"
        class="profile-form"
      >
        <div class="avatar-section">
          <el-avatar :size="100" :src="userForm.avatar || defaultAvatar"></el-avatar>
          <el-upload
            class="avatar-uploader"
            action="#"
            :http-request="uploadAvatar"
            :show-file-list="false"
            accept="image/*"
          >
            <el-button type="primary" size="small">更换头像</el-button>
          </el-upload>
        </div>
        
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" disabled />
        </el-form-item>
        
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="userForm.nickname" />
        </el-form-item>
        
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="userForm.phone" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" />
        </el-form-item>
        
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="userForm.gender">
            <el-radio :label="0">保密</el-radio>
            <el-radio :label="1">男</el-radio>
            <el-radio :label="2">女</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="生日" prop="birthday">
          <el-date-picker
            v-model="userForm.birthday"
            type="date"
            placeholder="选择生日"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="个人简介" prop="bio">
          <el-input 
            v-model="userForm.bio" 
            type="textarea" 
            :rows="3" 
            placeholder="介绍一下自己吧"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="updateProfile" :loading="submitting">保存修改</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 安全设置 -->
    <div class="profile-header">
      <h2 class="section-title">安全设置</h2>
    </div>
    
    <div class="security-content">
      <div class="security-item">
        <div class="security-info">
          <div class="security-title">登录密码</div>
          <div class="security-desc">定期修改密码有助于保护账号安全</div>
        </div>
        <div class="security-action">
          <el-button plain @click="showPasswordDialog">修改</el-button>
        </div>
      </div>
      
      <div class="security-item">
        <div class="security-info">
          <div class="security-title">绑定手机</div>
          <div class="security-desc">
            {{ userForm.phone ? `已绑定：${maskPhone(userForm.phone)}` : '绑定手机号码可以用于登录和找回密码' }}
          </div>
        </div>
        <div class="security-action">
          <el-button plain @click="showPhoneDialog">
            {{ userForm.phone ? '更换' : '绑定' }}
          </el-button>
        </div>
      </div>
      
      <div class="security-item">
        <div class="security-info">
          <div class="security-title">绑定邮箱</div>
          <div class="security-desc">
            {{ userForm.email ? `已绑定：${maskEmail(userForm.email)}` : '绑定邮箱可以用于登录和接收通知' }}
          </div>
        </div>
        <div class="security-action">
          <el-button plain @click="showEmailDialog">
            {{ userForm.email ? '更换' : '绑定' }}
          </el-button>
        </div>
      </div>
    </div>
    
    <!-- 修改密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="500px">
      <el-form 
        ref="passwordFormRef" 
        :model="passwordForm" 
        :rules="passwordRules" 
        label-width="100px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input 
            v-model="passwordForm.oldPassword" 
            type="password" 
            show-password
            placeholder="请输入原密码" 
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="passwordForm.newPassword" 
            type="password" 
            show-password
            placeholder="请输入新密码" 
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="passwordForm.confirmPassword" 
            type="password" 
            show-password
            placeholder="请再次输入新密码" 
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdatePassword" :loading="pwdSubmitting">
          确认修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import { getUserInfo, updateUserInfo, updatePassword as updatePasswordApi } from '@/api/user';

const userStore = useUserStore();

const loading = ref(false);
const submitting = ref(false);
const pwdSubmitting = ref(false);
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png';
const formRef = ref(null);
const passwordFormRef = ref(null);
const passwordDialogVisible = ref(false);

// 用户表单数据
const userForm = reactive({
  id: '',
  username: '',
  nickname: '',
  phone: '',
  email: '',
  gender: 0,
  birthday: '',
  bio: '',
  avatar: ''
});

// 密码表单数据
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 表单验证规则
const rules = {
  nickname: [
    { max: 20, message: '昵称不能超过20个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' }
  ],
  bio: [
    { max: 200, message: '简介不能超过200个字符', trigger: 'blur' }
  ]
};

// 密码表单验证规则
const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.newPassword) {
    callback(new Error('两次输入的密码不一致'));
  } else {
    callback();
  }
};

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
};

// 初始化
onMounted(async () => {
  await loadUserInfo();
});

// 加载用户信息
const loadUserInfo = async () => {
  try {
    loading.value = true;
    
    // 如果已有用户信息，直接使用
    if (userStore.userInfo && userStore.userInfo.id) {
      Object.assign(userForm, userStore.userInfo);
      return;
    }
    
    // 否则从API获取
    const result = await getUserInfo();
    Object.assign(userForm, result || {});
  } catch (error) {
    console.error('加载用户信息失败:', error);
    ElMessage.error('加载用户信息失败');
  } finally {
    loading.value = false;
  }
};

// 上传头像
const uploadAvatar = async (options) => {
  try {
    const file = options.file;
    
    // 这里应该调用上传文件的API
    // 模拟上传成功
    ElMessage.success('头像上传成功');
    
    // 更新头像URL（实际应该是后端返回的URL）
    userForm.avatar = URL.createObjectURL(file);
  } catch (error) {
    console.error('上传头像失败:', error);
    ElMessage.error('上传头像失败');
  }
};

// 保存修改
const updateProfile = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    submitting.value = true;
    
    await updateUserInfo(userForm);
    
    // 更新本地用户信息
    userStore.userInfo = { ...userForm };
    
    ElMessage.success('保存成功');
  } catch (error) {
    console.error('保存失败:', error);
    ElMessage.error('保存失败');
  } finally {
    submitting.value = false;
  }
};

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  
  // 重新加载用户信息
  loadUserInfo();
};

// 显示修改密码对话框
const showPasswordDialog = () => {
  passwordDialogVisible.value = true;
  
  // 清空密码表单
  passwordForm.oldPassword = '';
  passwordForm.newPassword = '';
  passwordForm.confirmPassword = '';
};

// 修改密码
const handleUpdatePassword = async () => {
  if (!passwordFormRef.value) return;
  
  try {
    await passwordFormRef.value.validate();
    
    pwdSubmitting.value = true;
    
    await updatePasswordApi({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    });
    
    ElMessage.success('密码修改成功');
    passwordDialogVisible.value = false;
  } catch (error) {
    console.error('修改密码失败:', error);
    ElMessage.error('修改密码失败');
  } finally {
    pwdSubmitting.value = false;
  }
};

// 显示绑定手机对话框
const showPhoneDialog = () => {
  // 实现绑定/更换手机的逻辑
  ElMessage.info('暂未实现绑定/更换手机功能');
};

// 显示绑定邮箱对话框
const showEmailDialog = () => {
  // 实现绑定/更换邮箱的逻辑
  ElMessage.info('暂未实现绑定/更换邮箱功能');
};

// 手机号码脱敏
const maskPhone = (phone) => {
  if (!phone) return '';
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2');
};

// 邮箱脱敏
const maskEmail = (email) => {
  if (!email) return '';
  const parts = email.split('@');
  if (parts.length !== 2) return email;
  
  let name = parts[0];
  const domain = parts[1];
  
  if (name.length <= 2) {
    name = name.substring(0, 1) + '**';
  } else {
    name = name.substring(0, 2) + '***';
  }
  
  return `${name}@${domain}`;
};
</script>

<style scoped>
.profile-container {
  padding: 10px;
}

.profile-header {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
}

.section-title {
  font-size: 18px;
  color: #303133;
  margin: 0;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 30px;
}

.avatar-uploader {
  margin-top: 15px;
}

.profile-form {
  max-width: 500px;
  margin: 0 auto;
}

/* 安全设置 */
.security-content {
  margin-bottom: 30px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.security-title {
  font-size: 16px;
  color: #303133;
  margin-bottom: 8px;
}

.security-desc {
  color: #909399;
  font-size: 14px;
}

@media (max-width: 576px) {
  .profile-form {
    max-width: 100%;
  }
  
  .security-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .security-action {
    margin-top: 15px;
  }
}
</style>
