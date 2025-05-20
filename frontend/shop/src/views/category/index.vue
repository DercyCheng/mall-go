<template>
  <div class="category-container">
    <div class="container">
      <!-- 分类标题与筛选 -->
      <div class="category-header">
        <div class="category-title">
          <h1>{{ categoryName }}</h1>
          <div class="search-results" v-if="keyword">
            搜索 "{{ keyword }}" 的结果
          </div>
        </div>
        <div class="category-filter">
          <el-select v-model="sortType" placeholder="排序方式" size="large">
            <el-option label="默认排序" value="default" />
            <el-option label="价格从低到高" value="price-asc" />
            <el-option label="价格从高到低" value="price-desc" />
            <el-option label="销量优先" value="sales" />
            <el-option label="最新上架" value="newest" />
          </el-select>
          <el-select v-model="pageSize" placeholder="每页数量" size="large">
            <el-option :label="12" :value="12" />
            <el-option :label="24" :value="24" />
            <el-option :label="36" :value="36" />
            <el-option :label="48" :value="48" />
          </el-select>
        </div>
      </div>
      
      <!-- 分类筛选器 -->
      <div class="filter-section" v-if="filters.length > 0">
        <div class="filter-item" v-for="filter in filters" :key="filter.key">
          <div class="filter-title">{{ filter.name }}</div>
          <div class="filter-values">
            <el-check-tag
              v-for="option in filter.options"
              :key="option.value"
              :checked="isFilterSelected(filter.key, option.value)"
              @change="toggleFilter(filter.key, option.value)"
            >
              {{ option.label }}
            </el-check-tag>
          </div>
        </div>
      </div>
      
      <!-- 产品列表 -->
      <div class="product-section" v-loading="loading">
        <div v-if="products.length === 0 && !loading" class="empty-result">
          <el-empty description="没有找到符合条件的商品"></el-empty>
        </div>
        <div v-else class="product-grid">
          <product-card
            v-for="product in products"
            :key="product.id"
            :product="product"
          ></product-card>
        </div>
      </div>
      
      <!-- 分页 -->
      <div class="pagination-section" v-if="products.length > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getProductList, getCategories, getCategoryDetail } from '@/api/product';
import ProductCard from '@/components/ProductCard.vue';

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const products = ref([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(24);
const sortType = ref('default');
const filters = ref([]);
const selectedFilters = ref({});
const allCategories = ref([]);
const currentCategory = ref(null);
const keyword = ref('');

// 分类名称
const categoryName = computed(() => {
  if (keyword.value) {
    return '商品搜索';
  }
  return currentCategory.value?.name || '全部商品';
});

// 监听路由变化，更新查询条件
watch(
  () => route.query,
  (newQuery) => {
    // 更新关键词
    keyword.value = newQuery.keyword || '';
    
    // 更新分类ID
    const categoryId = newQuery.categoryId;
    if (categoryId) {
      loadCategoryDetail(categoryId);
    } else {
      currentCategory.value = null;
    }
    
    // 更新排序方式
    sortType.value = newQuery.sort || 'default';
    
    // 重置分页
    currentPage.value = 1;
    
    // 加载数据
    loadProducts();
  },
  { immediate: true }
);

// 监听排序和分页变化，重新加载数据
watch([sortType, currentPage, pageSize], () => {
  loadProducts();
  
  // 更新URL查询参数，但不触发路由重新加载
  const query = { ...route.query };
  
  if (sortType.value !== 'default') {
    query.sort = sortType.value;
  } else {
    delete query.sort;
  }
  
  if (currentPage.value !== 1) {
    query.page = currentPage.value;
  } else {
    delete query.page;
  }
  
  if (pageSize.value !== 24) {
    query.pageSize = pageSize.value;
  } else {
    delete query.pageSize;
  }
  
  router.replace({ query });
});

// 初始化
onMounted(async () => {
  await loadCategories();
  initializeFilters();
});

// 加载所有分类
const loadCategories = async () => {
  try {
    const categories = await getCategories();
    allCategories.value = categories || [];
  } catch (error) {
    console.error('加载分类失败:', error);
  }
};

// 加载分类详情
const loadCategoryDetail = async (categoryId) => {
  try {
    // 如果已经有分类数据，先从中查找
    const category = allCategories.value.find(cat => cat.id == categoryId);
    if (category) {
      currentCategory.value = category;
      return;
    }
    
    // 否则从API获取详情
    const detail = await getCategoryDetail(categoryId);
    currentCategory.value = detail;
  } catch (error) {
    console.error('加载分类详情失败:', error);
    currentCategory.value = null;
  }
};

// 初始化筛选器
const initializeFilters = () => {
  // 这里可以根据实际业务来设置不同的筛选条件
  filters.value = [
    {
      key: 'priceRange',
      name: '价格区间',
      options: [
        { label: '0-100元', value: '0-100' },
        { label: '100-300元', value: '100-300' },
        { label: '300-500元', value: '300-500' },
        { label: '500-1000元', value: '500-1000' },
        { label: '1000元以上', value: '1000-' }
      ]
    },
    {
      key: 'inStock',
      name: '库存状态',
      options: [
        { label: '有货', value: 'true' }
      ]
    },
    {
      key: 'promotion',
      name: '优惠活动',
      options: [
        { label: '特价商品', value: 'special' },
        { label: '新品', value: 'new' },
        { label: '促销', value: 'promotion' }
      ]
    }
  ];
};

// 判断筛选条件是否被选中
const isFilterSelected = (key, value) => {
  if (!selectedFilters.value[key]) {
    return false;
  }
  return selectedFilters.value[key].includes(value);
};

// 切换筛选条件
const toggleFilter = (key, value) => {
  if (!selectedFilters.value[key]) {
    selectedFilters.value[key] = [];
  }
  
  const index = selectedFilters.value[key].indexOf(value);
  if (index === -1) {
    selectedFilters.value[key].push(value);
  } else {
    selectedFilters.value[key].splice(index, 1);
  }
  
  // 重置分页
  currentPage.value = 1;
  loadProducts();
};

// 加载商品列表
const loadProducts = async () => {
  try {
    loading.value = true;
    
    // 构建查询参数
    const params = {
      pageNum: currentPage.value,
      pageSize: pageSize.value
    };
    
    // 添加关键词
    if (keyword.value) {
      params.keyword = keyword.value;
    }
    
    // 添加分类ID
    if (currentCategory.value) {
      params.categoryId = currentCategory.value.id;
    }
    
    // 添加排序方式
    if (sortType.value !== 'default') {
      if (sortType.value === 'price-asc') {
        params.sort = 'price_asc';
      } else if (sortType.value === 'price-desc') {
        params.sort = 'price_desc';
      } else {
        params.sort = sortType.value;
      }
    }
    
    // 添加筛选条件
    for (const [key, values] of Object.entries(selectedFilters.value)) {
      if (values.length > 0) {
        params[key] = values.join(',');
      }
    }
    
    const result = await getProductList(params);
    products.value = result.list || [];
    total.value = result.total || 0;
  } catch (error) {
    console.error('加载商品失败:', error);
  } finally {
    loading.value = false;
  }
};

// 每页数量变化
const handleSizeChange = (val) => {
  pageSize.value = val;
};

// 页码变化
const handleCurrentChange = (val) => {
  currentPage.value = val;
};
</script>

<style scoped>
.category-container {
  background-color: #f5f7fa;
  padding: 20px 0 40px;
}

.container {
  width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

/* 分类标题与筛选 */
.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.category-title h1 {
  font-size: 22px;
  margin: 0 0 5px;
  color: #303133;
}

.search-results {
  font-size: 14px;
  color: #909399;
}

.category-filter {
  display: flex;
}

.category-filter .el-select {
  margin-left: 15px;
}

/* 筛选器样式 */
.filter-section {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.filter-item {
  margin-bottom: 15px;
}

.filter-item:last-child {
  margin-bottom: 0;
}

.filter-title {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.filter-values {
  display: flex;
  flex-wrap: wrap;
}

.filter-values .el-check-tag {
  margin-right: 10px;
  margin-bottom: 10px;
  padding: 5px 12px;
}

/* 产品列表 */
.product-section {
  margin-bottom: 20px;
  min-height: 400px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.empty-result {
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fff;
  padding: 50px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

/* 分页 */
.pagination-section {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
  }
}

@media (max-width: 992px) {
  .product-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .category-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .category-filter {
    margin-top: 15px;
  }
  
  .category-filter .el-select:first-child {
    margin-left: 0;
  }
  
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
  .category-filter {
    width: 100%;
    flex-direction: column;
  }
  
  .category-filter .el-select {
    margin-left: 0;
    margin-top: 10px;
    width: 100%;
  }
  
  .category-filter .el-select:first-child {
    margin-top: 0;
  }
}
</style>
