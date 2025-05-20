<template>
  <div class="favorites-container">
                      <span class="original-price">¥{{ item.originalPrice.toFixed(2) }}</span>
                </div>
                <div class="product-actions-bottom">
                  <el-button type="primary" size="small" @click="handleAddToCart(item.id)">Add to Cart</el-button>
                </div>
              </div>
            </div>d class="favorites-card">
      <template #header>
        <div class="favorites-header">
          <h2>My Favorites</h2>
          <div class="header-actions">
            <el-select v-model="sortBy" placeholder="Sort by" size="small">
              <el-option label="Date Added (Newest)" value="dateDesc" />
              <el-option label="Date Added (Oldest)" value="dateAsc" />
              <el-option label="Price (Low to High)" value="priceAsc" />
              <el-option label="Price (High to Low)" value="priceDesc" />
            </el-select>
          </div>
        </div>
      </template>
      
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      
      <div v-else-if="favorites.length === 0" class="empty-favorites">
        <el-empty description="No favorites yet">
          <template #image>
            <i class="el-icon-star-off" style="font-size: 60px; color: #909399;"></i>
          </template>
          <el-button type="primary" @click="goToShop">Browse Products</el-button>
        </el-empty>
      </div>
      
      <div v-else class="favorites-list">
        <el-row :gutter="20">
          <el-col v-for="item in favorites" :key="item.id" :xs="24" :sm="12" :md="8" :lg="6" class="product-col">
            <div class="product-item">
              <div class="product-actions">
                <el-button 
                  type="danger" 
                  circle 
                  icon="Delete"
                  size="small"
                  @click="handleRemoveFromFavorites(item.id)"
                  title="Remove from favorites"
                />
              </div>
              <div class="product-image" @click="viewProductDetail(item.id)">
                <img :src="item.image" :alt="item.name" />
              </div>
              <div class="product-info">
                <div class="product-name" @click="viewProductDetail(item.id)">{{ item.name }}</div>
                <div class="product-price">
                  <span class="current-price">¥{{ item.price.toFixed(2) }}</span>
                  <span v-if="item.originalPrice > item.price" class="original-price">¥{{ item.originalPrice.toFixed(2) }}</span>
                </div>
                <div class="product-actions-bottom">
                  <el-button type="primary" size="small" @click="addToCart(item.id)">Add to Cart</el-button>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
      
      <el-pagination
        v-if="total > pageSize"
        layout="prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="currentPage"
        @current-change="handlePageChange"
        class="pagination"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getFavoritesList, removeFromFavorites as removeFromFavoritesApi } from '@/api/user';
import { addToCart as addToCartApi } from '@/api/cart';

const router = useRouter();
const loading = ref(true);
const favorites = ref([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(12);
const sortBy = ref('dateDesc');

// Fetch favorites data
const fetchFavorites = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      sort: sortBy.value,
    };
    
    const res = await getFavoritesList(params);
    favorites.value = res.data.items;
    total.value = res.data.total;
  } catch (error) {
    console.error('Failed to fetch favorites:', error);
    ElMessage.error('Failed to load favorites list');
  } finally {
    loading.value = false;
  }
};

// Handle pagination
const handlePageChange = (page) => {
  currentPage.value = page;
  fetchFavorites();
};

// Navigate to product detail
const viewProductDetail = (productId) => {
  router.push(`/product/${productId}`);
};

// Navigate to shop/category page
const goToShop = () => {
  router.push('/category');
};

// Add product to cart
const handleAddToCart = async (productId) => {
  try {
    await addToCartApi({
      productId,
      quantity: 1
    });
    ElMessage.success('Added to cart successfully');
  } catch (error) {
    console.error('Failed to add to cart:', error);
    ElMessage.error('Failed to add to cart');
  }
};

// Remove product from favorites
const handleRemoveFromFavorites = async (productId) => {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to remove this item from favorites?',
      'Confirm Removal',
      {
        confirmButtonText: 'Remove',
        cancelButtonText: 'Cancel',
        type: 'warning',
      }
    );
    
    await removeFromFavoritesApi(productId);
    ElMessage.success('Removed from favorites');
    fetchFavorites(); // Reload the list
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove from favorites:', error);
      ElMessage.error('Failed to remove from favorites');
    }
  }
};

// Watch for sort change
watch(sortBy, () => {
  fetchFavorites();
});

onMounted(() => {
  fetchFavorites();
});
</script>

<style scoped>
.favorites-container {
  padding: 20px 0;
}

.favorites-card {
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.favorites-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.favorites-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.loading-container {
  padding: 20px 0;
}

.empty-favorites {
  padding: 40px 0;
  text-align: center;
}

.product-col {
  margin-bottom: 20px;
}

.product-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
  transition: all 0.3s;
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.product-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
}

.product-actions {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
  opacity: 0;
  transition: opacity 0.3s;
}

.product-item:hover .product-actions {
  opacity: 1;
}

.product-image {
  height: 200px;
  overflow: hidden;
  cursor: pointer;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.product-image:hover img {
  transform: scale(1.05);
}

.product-info {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.product-name {
  font-size: 14px;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  cursor: pointer;
  color: #303133;
}

.product-name:hover {
  color: var(--el-color-primary);
}

.product-price {
  margin-bottom: 12px;
}

.current-price {
  font-weight: bold;
  color: #f56c6c;
  font-size: 16px;
  margin-right: 8px;
}

.original-price {
  color: #909399;
  text-decoration: line-through;
  font-size: 12px;
}

.product-actions-bottom {
  margin-top: auto;
}

.pagination {
  text-align: center;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .favorites-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-actions {
    margin-top: 10px;
    width: 100%;
  }
  
  .header-actions .el-select {
    width: 100%;
  }
  
  .product-image {
    height: 180px;
  }
}
</style>
