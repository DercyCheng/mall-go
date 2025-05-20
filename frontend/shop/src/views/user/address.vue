<template>
  <div class="address-container">
    <div class="address-header">
      <h2 class="section-title">收货地址</h2>
      <el-button type="primary" @click="showAddressDialog">新增收货地址</el-button>
    </div>
    
    <div class="address-content" v-loading="loading">
      <div v-if="addresses.length === 0" class="no-address">
        <el-empty description="暂无收货地址，请添加">
          <el-button type="primary" @click="showAddressDialog">添加地址</el-button>
        </el-empty>
      </div>
      
      <div v-else class="address-list">
        <div 
          v-for="address in addresses" 
          :key="address.id"
          class="address-item"
        >
          <div class="address-info">
            <div class="address-name">
              <span class="name">{{ address.name }}</span>
              <span class="phone">{{ address.phone }}</span>
              <span class="tag" v-if="address.defaultStatus">默认地址</span>
            </div>
            <div class="address-detail">
              {{ formatAddress(address) }}
            </div>
          </div>
          <div class="address-actions">
            <el-button type="primary" link @click="editAddress(address)">
              编辑
            </el-button>
            <el-button type="danger" link @click="confirmDelete(address)">
              删除
            </el-button>
            <el-button v-if="!address.defaultStatus" type="success" link @click="setDefault(address.id)">
              设为默认
            </el-button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 地址表单弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑收货地址' : '添加收货地址'"
      width="500px"
    >
      <el-form 
        ref="addressFormRef" 
        :model="addressForm" 
        :rules="rules" 
        label-width="80px"
      >
        <el-form-item label="收货人" prop="name">
          <el-input v-model="addressForm.name" placeholder="请输入收货人姓名" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="addressForm.phone" placeholder="请输入手机号码" />
        </el-form-item>
        
        <el-form-item label="所在地区" prop="region">
          <el-cascader
            v-model="addressForm.region"
            :options="regionOptions"
            placeholder="请选择省/市/区"
          />
        </el-form-item>
        
        <el-form-item label="详细地址" prop="detailAddress">
          <el-input 
            v-model="addressForm.detailAddress" 
            type="textarea" 
            placeholder="请输入详细地址信息，如街道、门牌号等"
          />
        </el-form-item>
        
        <el-form-item label="邮政编码" prop="postCode">
          <el-input v-model="addressForm.postCode" placeholder="请输入邮政编码" />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="addressForm.defaultStatus">设为默认收货地址</el-checkbox>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveAddress">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 删除确认 -->
    <el-dialog
      v-model="deleteVisible"
      title="删除确认"
      width="400px"
    >
      <p>确定要删除这个收货地址吗？</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="deleteVisible = false">取消</el-button>
          <el-button type="danger" @click="handleDelete">确定删除</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getAddressList, addAddress, updateAddress, deleteAddress } from '@/api/user';

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入收货人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  region: [
    { required: true, message: '请选择所在地区', trigger: 'change' }
  ],
  detailAddress: [
    { required: true, message: '请输入详细地址', trigger: 'blur' }
  ],
  postCode: [
    { pattern: /^\d{6}$/, message: '邮政编码格式不正确', trigger: 'blur' }
  ]
};

// 数据状态
const loading = ref(false);
const addresses = ref([]);
const dialogVisible = ref(false);
const deleteVisible = ref(false);
const isEdit = ref(false);
const currentDeleteId = ref(null);
const addressFormRef = ref(null);

// 地址表单
const addressForm = reactive({
  id: null,
  name: '',
  phone: '',
  region: [],
  detailAddress: '',
  postCode: '',
  defaultStatus: false
});

// 简化的地区选项，实际项目中应该从接口获取
const regionOptions = [
  {
    value: '北京市',
    label: '北京市',
    children: [
      {
        value: '北京市',
        label: '北京市',
        children: [
          { value: '东城区', label: '东城区' },
          { value: '西城区', label: '西城区' }
          // 其他区县
        ]
      }
    ]
  },
  {
    value: '上海市',
    label: '上海市',
    children: [
      {
        value: '上海市',
        label: '上海市',
        children: [
          { value: '黄浦区', label: '黄浦区' },
          { value: '徐汇区', label: '徐汇区' }
          // 其他区县
        ]
      }
    ]
  }
  // 其他省份
];

// 获取地址列表
const fetchAddresses = async () => {
  loading.value = true;
  try {
    const res = await getAddressList();
    addresses.value = res.data || [];
  } catch (error) {
    console.error('获取地址列表失败:', error);
    ElMessage.error('获取地址列表失败');
  } finally {
    loading.value = false;
  }
};

// 格式化地址
const formatAddress = (address) => {
  let result = '';
  if (address.province) result += address.province;
  if (address.city) result += ' ' + address.city;
  if (address.region) result += ' ' + address.region;
  if (address.detailAddress) result += ' ' + address.detailAddress;
  return result;
};

// 显示地址弹窗
const showAddressDialog = () => {
  isEdit.value = false;
  resetAddressForm();
  dialogVisible.value = true;
};

// 编辑地址
const editAddress = (address) => {
  isEdit.value = true;
  
  // 填充表单数据
  addressForm.id = address.id;
  addressForm.name = address.name;
  addressForm.phone = address.phone;
  addressForm.region = [address.province, address.city, address.region].filter(Boolean);
  addressForm.detailAddress = address.detailAddress;
  addressForm.postCode = address.postCode;
  addressForm.defaultStatus = address.defaultStatus;
  
  dialogVisible.value = true;
};

// 重置表单
const resetAddressForm = () => {
  addressForm.id = null;
  addressForm.name = '';
  addressForm.phone = '';
  addressForm.region = [];
  addressForm.detailAddress = '';
  addressForm.postCode = '';
  addressForm.defaultStatus = false;
  
  // 如果表单引用存在，重置校验状态
  if (addressFormRef.value) {
    addressFormRef.value.resetFields();
  }
};

// 保存地址
const saveAddress = async () => {
  if (!addressFormRef.value) return;
  
  await addressFormRef.value.validate(async (valid) => {
    if (!valid) return;
    
    try {
      loading.value = true;
      
      // 构造请求数据
      const data = {
        name: addressForm.name,
        phone: addressForm.phone,
        province: addressForm.region[0],
        city: addressForm.region[1],
        region: addressForm.region[2],
        detailAddress: addressForm.detailAddress,
        postCode: addressForm.postCode,
        defaultStatus: addressForm.defaultStatus ? 1 : 0
      };
      
      if (isEdit.value) {
        // 更新地址
        await updateAddress(addressForm.id, data);
        ElMessage.success('地址更新成功');
      } else {
        // 新增地址
        await addAddress(data);
        ElMessage.success('地址添加成功');
      }
      
      // 关闭弹窗并刷新列表
      dialogVisible.value = false;
      fetchAddresses();
    } catch (error) {
      console.error('保存地址失败:', error);
      ElMessage.error('保存地址失败');
    } finally {
      loading.value = false;
    }
  });
};

// 确认删除
const confirmDelete = (address) => {
  currentDeleteId.value = address.id;
  deleteVisible.value = true;
};

// 执行删除
const handleDelete = async () => {
  if (!currentDeleteId.value) return;
  
  try {
    loading.value = true;
    await deleteAddress(currentDeleteId.value);
    ElMessage.success('地址删除成功');
    deleteVisible.value = false;
    fetchAddresses();
  } catch (error) {
    console.error('删除地址失败:', error);
    ElMessage.error('删除地址失败');
  } finally {
    loading.value = false;
  }
};

// 设置默认地址
const setDefault = async (id) => {
  try {
    loading.value = true;
    const address = addresses.value.find(item => item.id === id);
    if (!address) return;
    
    const data = { ...address, defaultStatus: 1 };
    await updateAddress(id, data);
    ElMessage.success('设置默认地址成功');
    fetchAddresses();
  } catch (error) {
    console.error('设置默认地址失败:', error);
    ElMessage.error('设置默认地址失败');
  } finally {
    loading.value = false;
  }
};

// 组件挂载时获取地址列表
onMounted(() => {
  fetchAddresses();
});
</script>

<style scoped>
.address-container {
  padding: 20px;
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: bold;
  margin: 0;
}

.address-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
}

.address-item {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  transition: all 0.3s;
}

.address-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.address-info {
  flex: 1;
}

.address-name {
  margin-bottom: 8px;
}

.name {
  font-weight: bold;
  margin-right: 12px;
}

.phone {
  color: #606266;
}

.tag {
  background-color: #409eff;
  color: white;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 8px;
}

.address-detail {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-all;
}

.address-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.no-address {
  margin: 40px 0;
}

@media (max-width: 768px) {
  .address-list {
    grid-template-columns: 1fr;
  }
  
  .address-item {
    flex-direction: column;
  }
  
  .address-actions {
    margin-top: 16px;
    flex-direction: row;
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
