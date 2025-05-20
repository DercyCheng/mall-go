<template>
  <div class="address-manager">
    <!-- 地址列表 -->
    <div class="address-list" v-if="!showAddressForm && addresses.length > 0">
      <div v-for="address in addresses" :key="address.id" class="address-item">
        <div class="address-info">
          <div class="address-name">
            <span class="name">{{ address.name }}</span>
            <span class="phone">{{ address.phone }}</span>
            <span class="tag" v-if="address.defaultStatus">默认</span>
          </div>
          <div class="address-detail">
            {{ formatAddress(address) }}
          </div>
        </div>
        <div class="address-actions">
          <el-button type="primary" link @click="editAddress(address)">编辑</el-button>
          <el-button type="primary" link @click="deleteAddress(address.id)">删除</el-button>
          <el-button 
            type="primary" 
            link 
            v-if="!address.defaultStatus"
            @click="setDefault(address.id)"
          >
            设为默认
          </el-button>
        </div>
      </div>
    </div>
    
    <div class="empty-address" v-if="!showAddressForm && addresses.length === 0">
      <el-empty description="暂无收货地址">
        <el-button type="primary" @click="addAddress">添加地址</el-button>
      </el-empty>
    </div>
    
    <!-- 地址表单 -->
    <div class="address-form" v-if="showAddressForm">
      <h3 class="form-title">{{ currentAddress.id ? '编辑地址' : '新增地址' }}</h3>
      
      <el-form 
        ref="formRef" 
        :model="currentAddress" 
        :rules="rules" 
        label-width="100px"
      >
        <el-form-item label="收货人" prop="name">
          <el-input v-model="currentAddress.name" placeholder="请输入收货人姓名" />
        </el-form-item>
        
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="currentAddress.phone" placeholder="请输入联系电话" />
        </el-form-item>
        
        <el-form-item label="所在地区" prop="regions">
          <el-cascader
            v-model="regions"
            :options="regionOptions"
            placeholder="请选择所在地区"
            @change="handleRegionChange"
          />
        </el-form-item>
        
        <el-form-item label="详细地址" prop="detailAddress">
          <el-input 
            v-model="currentAddress.detailAddress" 
            type="textarea" 
            :rows="2"
            placeholder="请输入详细地址，如街道、门牌号等" 
          />
        </el-form-item>
        
        <el-form-item label="邮政编码" prop="postCode">
          <el-input v-model="currentAddress.postCode" placeholder="请输入邮政编码" />
        </el-form-item>
        
        <el-form-item>
          <el-checkbox v-model="currentAddress.defaultStatus">设为默认收货地址</el-checkbox>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="saveAddress">保存</el-button>
          <el-button @click="cancelEdit">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 底部按钮 -->
    <div class="address-footer" v-if="!showAddressForm && addresses.length > 0">
      <el-button type="primary" @click="addAddress">新增收货地址</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

const props = defineProps({
  addresses: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['save-address', 'delete-address', 'set-default']);

const showAddressForm = ref(false);
const formRef = ref(null);
const currentAddress = reactive({
  id: '',
  name: '',
  phone: '',
  province: '',
  city: '',
  region: '',
  detailAddress: '',
  postCode: '',
  defaultStatus: false
});

const regions = ref([]);

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入收货人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  regions: [
    { required: true, message: '请选择所在地区', trigger: 'change' }
  ],
  detailAddress: [
    { required: true, message: '请输入详细地址', trigger: 'blur' },
    { min: 5, max: 100, message: '长度在 5 到 100 个字符', trigger: 'blur' }
  ],
  postCode: [
    { pattern: /^\d{6}$/, message: '邮政编码格式不正确', trigger: 'blur' }
  ]
};

// 地区数据（这里使用简化版，实际应用中应该从API获取完整的省市区数据）
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
          { value: '西城区', label: '西城区' },
          { value: '朝阳区', label: '朝阳区' },
          { value: '海淀区', label: '海淀区' },
          { value: '丰台区', label: '丰台区' }
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
          { value: '徐汇区', label: '徐汇区' },
          { value: '长宁区', label: '长宁区' },
          { value: '静安区', label: '静安区' },
          { value: '普陀区', label: '普陀区' }
        ]
      }
    ]
  },
  {
    value: '广东省',
    label: '广东省',
    children: [
      {
        value: '广州市',
        label: '广州市',
        children: [
          { value: '越秀区', label: '越秀区' },
          { value: '荔湾区', label: '荔湾区' },
          { value: '海珠区', label: '海珠区' },
          { value: '天河区', label: '天河区' },
          { value: '白云区', label: '白云区' }
        ]
      },
      {
        value: '深圳市',
        label: '深圳市',
        children: [
          { value: '福田区', label: '福田区' },
          { value: '罗湖区', label: '罗湖区' },
          { value: '南山区', label: '南山区' },
          { value: '宝安区', label: '宝安区' },
          { value: '龙岗区', label: '龙岗区' }
        ]
      }
    ]
  }
];

// 格式化地址
const formatAddress = (address) => {
  return `${address.province} ${address.city} ${address.region} ${address.detailAddress}`;
};

// 地区选择变化
const handleRegionChange = (value) => {
  if (value && value.length === 3) {
    currentAddress.province = value[0];
    currentAddress.city = value[1];
    currentAddress.region = value[2];
  }
};

// 添加地址
const addAddress = () => {
  resetForm();
  showAddressForm.value = true;
};

// 编辑地址
const editAddress = (address) => {
  // 复制地址对象，避免直接修改原对象
  Object.assign(currentAddress, address);
  
  // 设置地区选择器的值
  regions.value = [address.province, address.city, address.region];
  
  showAddressForm.value = true;
};

// 删除地址
const deleteAddress = async (addressId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个地址吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    emit('delete-address', addressId);
  } catch (error) {
    // 用户取消操作
  }
};

// 设为默认地址
const setDefault = (addressId) => {
  emit('set-default', addressId);
};

// 保存地址
const saveAddress = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    emit('save-address', { ...currentAddress });
    
    // 关闭表单
    showAddressForm.value = false;
  } catch (error) {
    console.error('表单验证失败:', error);
  }
};

// 取消编辑
const cancelEdit = () => {
  showAddressForm.value = false;
};

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  
  // 清空当前地址数据
  Object.assign(currentAddress, {
    id: '',
    name: '',
    phone: '',
    province: '',
    city: '',
    region: '',
    detailAddress: '',
    postCode: '',
    defaultStatus: false
  });
  
  // 清空地区选择
  regions.value = [];
};
</script>

<style scoped>
.address-manager {
  padding: 10px 0;
}

/* 地址列表 */
.address-list {
  margin-bottom: 20px;
}

.address-item {
  display: flex;
  justify-content: space-between;
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}

.address-item:last-child {
  border-bottom: none;
}

.address-info {
  flex: 1;
}

.address-name {
  margin-bottom: 8px;
}

.name {
  font-weight: bold;
  color: #303133;
}

.phone {
  margin-left: 15px;
  color: #606266;
}

.tag {
  display: inline-block;
  padding: 2px 6px;
  background-color: #409EFF;
  color: #fff;
  font-size: 12px;
  border-radius: 2px;
  margin-left: 10px;
}

.address-detail {
  color: #606266;
  line-height: 1.5;
}

.address-actions {
  display: flex;
  align-items: center;
}

/* 空地址 */
.empty-address {
  padding: 30px 0;
  text-align: center;
}

/* 地址表单 */
.address-form {
  padding: 10px 0;
}

.form-title {
  font-size: 18px;
  color: #303133;
  margin-bottom: 20px;
}

/* 底部按钮 */
.address-footer {
  padding-top: 20px;
  text-align: center;
  border-top: 1px solid #ebeef5;
}
</style>
