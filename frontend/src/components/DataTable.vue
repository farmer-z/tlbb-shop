<template>
  <main>
    <div class="shop-container">
      <button @click="fetchData">刷新数据</button>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <span>数据加载中...</span>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="error-container">
        <div class="error-icon">⚠️</div>
        <div class="error-message">{{ errorMessage }}</div>
        <button @click="fetchData">重试</button>
      </div>

      <!-- 正常数据展示 -->
      <div v-else class="store-container">
          <div class="city-nav">
            <div
                v-for="shopId in tableData"
                :key="shopId"
                class="nav-item"
                :class="{ active: selectMenu?.shopId === shopId }"
                @click="selectMenu(shopId)"
            >
              大商店{{ shopId }}
            </div>
          </div>

          <!-- 第二层：区域列表 -->
          <div class="region-panel" v-if="menuData.length > 0">
            <div class="region-list">
              <div
                  v-for="menuId in menuData"
                  :key="menuId"
                  class="region-item"
                  :class="{ active: selectedSubMenu?.menuId === menuId }"
                  @click="selectedSubMenu(shopIdCache, menuId)"
              >
                商店{{ menuId }}
              </div>
            </div>

            <!-- 第三层：店铺列表 -->
            <div class="shop-list" v-if="subMenuData.length > 0">
              <div
                  v-for="subMenuId in subMenuData"
                  :key="subMenuId"
                  class="shop-item"
                  :class="{ active: selectedShopItem?.subMenuId === subMenuId }"
                  @click="selectedShopItem(shopIdCache, menuIdCache, subMenuId)"
              >
                店铺{{ subMenuId }}
              </div>
            </div>
          </div>

          <!-- 商品展示区 -->
          <div class="product-panel" v-if="itemData.length > 0">
            <div class="product-list">
              <div v-for="item in itemData" :key="item.Index" class="product-item">
                <div class="item-name">{{ item.itemName }}</div>
                <div class="item-price">元宝: {{ item.itemPrice }}</div>
              </div>
            </div>
            <!-- 分页控制 -->
            <div class="pagination">
              <button>前一页</button>
              <span>1/1</span>
              <button>后一页</button>
            </div>
          </div>
        </div>
      </div>
  </main>
</template>



<script setup>
import { ref} from 'vue'
import {GetShopTable, GetSubMenus, GetShopItems, GetMenuItems} from '../../wailsjs/go/main/App'
let shopIdCache = ref(1)
let menuIdCache = ref(1)
let subMenuIdCache = ref(1)
const tableData = ref([])
const menuData = ref([])
const subMenuData = ref([])
const itemData = ref([])
const loading = ref(false)
const error = ref(false)
const errorMessage = ref('')
const fetchData = async () => {
  try {
    loading.value = true
    error.value = false

    // 正确使用响应式赋值
    const response = await GetShopTable()
    console.log('tableData:', tableData)
    tableData.value = response

  } catch (err) {
    error.value = true
    errorMessage.value = `获取数据失败: ${err.message}`
    console.error('请求错误:', err)
  } finally {
    loading.value = false
  }
}
function selectMenu(shopId) {
  shopIdCache = shopId
  GetMenuItems(shopId).then((response) => {
    menuData.value=response
  })
}

function selectedSubMenu(shopId,menuId) {
  menuIdCache = menuId
  GetSubMenus(shopId,menuId).then((response) => {
    subMenuData.value=response
  })
}

function selectedShopItem(shopId,menuId,subMenuId) {
  subMenuIdCache=subMenuId
  GetShopItems(shopId,menuId,subMenuId).then((response) => {
    console.log('itemData:', response)
    itemData.value=response
  })
}
</script>

<style scoped>
.store-container {
  display: grid;
  grid-template-columns: 200px 300px 1fr;
  height: 100vh;
  background: #34495e;
}

.city-nav {
  padding: 20px;
  background: #2c3e50;
  color: black;
}

.nav-item {
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 4px;
  margin-bottom: 8px;
}

.nav-item:hover,
.nav-item.active {
  background: #34495e;
}

.region-panel {
  display: flex;
  background: #34495e;
  border-right: 1px solid #34495e;
}

.region-list {
  width: 150px;
  padding: 15px;
  border-right: 1px solid #34495e;
}

.region-item {
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.region-item:hover,
.region-item.active {
  background: #34495e;
  color: #3498db;
}

.shop-list {
  flex: 1;
  padding: 15px;
}

.shop-item {
  padding: 12px;
  cursor: pointer;
  border-radius: 4px;
  margin-bottom: 8px;
  transition: all 0.2s;
}

.shop-item:hover,
.shop-item.active {
  background: #34495e;
}

.product-panel {
  padding: 20px;
  background: black;
}

.product-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 15px;
}

.product-item {
  border: 1px solid #34495e;
  padding: 15px;
  border-radius: 6px;
  transition: transform 0.2s;
}

.product-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.item-name {
  font-weight: 500;
  margin-bottom: 8px;
}

.item-price {
  color: #e74c3c;
  font-size: 0.9em;
}

.pagination {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: center;
  padding: 15px;
  background: #34495e;
  border-radius: 6px;
}
</style>



