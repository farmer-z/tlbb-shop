<template>
  <main>
    <div class="shop-container">
      <button @click="fetchData">加载商店</button>
      <button @click="showModal = true">搜索物品</button>
      <button @click="reloadBasicData">重新加载基础数据</button>

      <!-- 右上角弹窗 -->
      <teleport to="body">
        <transition name="slide-fade">
          <div
              v-if="showModal"
              class="corner-modal"
              :style="modalPosition"
          >
            <div>
              <p class="item-price">双击即可复制</p>
            </div>
            <div class="modal-header">
              <input id="itemName" type="text" v-model="itemName" placeholder="请输入物品名称,双击即可复制">
              <button @click="searchItem(itemName)">搜索</button>
              <button @click="closeModal">关闭</button>-->
            </div>
            <div class="modal-content" v-if="searchItemData.length > 0" >
              <div v-for="item in searchItemData" :key="item.itemId" class="copyable-container"  @dblclick="handleCopy(item.itemId)">
                <div class="item-price" >物品id:{{ item.itemId }}</div>
                <div class="item-price" >物品名称: {{ item.itemName }}</div>
              </div>
              <div v-if="showFeedback" class="copy-feedback">
                √ 已复制
                </div>
            </div>
          </div>
        </transition>
      </teleport>
      <teleport to="body">
        <!-- 新增物品弹窗 -->
        <transition name="slide-fade">
          <div v-if="showAddModal" class="corner-add-modal">
            <div class="modal-header">
              <h3>新增物品</h3>
              <button @click="closeAddModal">关闭</button>
            </div>
            <div class="modal-content">
              <form @submit.prevent="submitNewItem">
                <div class="form-group">
                  <label>物品ID:</label>
                  <input
                      v-model="newItem.itemId"
                      type="number"
                      required
                      @blur="fetchItemName"
                      :disabled="nameLoading"
                  >
                  <span v-if="nameLoading" class="loading-text">获取中...</span>
                </div>
                <div class="form-group">
                  <label>物品名称:</label>
                  <input
                      v-model="newItem.itemName"
                      required
                      :readonly="!!newItem.itemId"
                  >
                </div>
                <div class="form-group">
                  <label>物品数量:</label>
                  <input v-model="newItem.itemCount" type="number" required>
                </div>
                <div class="form-group">
                  <label>物品价格:</label>
                  <input v-model="newItem.itemPrice" type="number" required>
                </div>
                <button type="submit">提交</button>
              </form>
            </div>
          </div>
        </transition>
      </teleport>

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
            <div v-for="item in itemData" :key="item.index" class="product-item">
              <div class="item-name">序号:{{ item.index }}</div>
              <div class="item-name">物品id:{{ item.itemId }}</div>
              <div class="item-name">{{ item.itemName }}</div>
              <div class="item-price">元宝: {{ item.itemPrice }}</div>
            </div>
            <button @click="openAddModal" class="product-item">新增物品</button>
          </div>
        </div>



      </div>
    </div>
  </main>
</template>


<script setup>
import {ref, computed} from 'vue'
import {GetMenuItems, GetShopItems, GetShopTable, GetSubMenus, SearchItem, MatchItemName,AddShopItem,ReloadItemMap} from '../../wailsjs/go/main/App.js'

// ------------------------------基础数据-----------------------------------
let tableData = ref([])
const fetchData = async () => {
  try {
    loading.value = true
    error.value = false
    tableData = ref([])
    GetShopTable().then((response) => {
      tableData.value = response
    })


  } catch (err) {
    error.value = true
    errorMessage.value = `获取数据失败: ${err.message}`
    console.error('请求错误:', err)
  } finally {
    loading.value = false
  }
}
const reloadBasicData = async () => {
  await ReloadItemMap()
}
// ------------------------------基础数据 end-----------------------------------

// ------------------------------弹窗相关-----------------------------------
const showAddModal = ref(false)
const showModal = ref(false)
const showFeedback =ref(false)
const closeAddModal = () => {
  showAddModal.value = false
}

const closeModal = () => {
  showModal.value = false
}

const handleCopy = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    showFeedback.value = true
    setTimeout(()=> showFeedback.value = false, 2000)
  } catch {
    // 回退方案
    const input = document.createElement('input')
    input.value = text
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    showFeedback.value = true
    setTimeout(()=>showFeedback.value = false, 2000)
  }
}

// 根据触发按钮定位
const modalPosition = computed(() => {
  if (!triggerRef.value) return {}
  const rect = triggerRef.value.getBoundingClientRect()
  return {
    top: `${rect.top}px`,
    right: `${window.innerWidth - rect.right}px`
  }
})

// 打开弹窗时自动获取当前选中店铺信息
const openAddModal = () => {
  newItem.value = {
    ...newItem.value,
    shopId: shopIdCache,
    menuId: menuIdCache,
    subMenuId: subMenuIdCache,
  }
  showAddModal.value = true
}
// ------------------------------弹窗相关 end-----------------------------------

// ------------------------------加载商店数据-----------------------------------

const menuData = ref([])
const subMenuData = ref([])
const itemData = ref([])

const loading = ref(false)
const error = ref(false)
const errorMessage = ref('')

const triggerRef = ref(null)
const nameLoading = ref(false)

// 获取物品名称方法

const fetchItemName = async () => {
  if (!newItem.value.itemId) return

  try {
    nameLoading.value = true
    const res = await MatchItemName(newItem.value.itemId.toString())

    if (res.length > 0) {
      newItem.value.itemName = res
    } else {
      alert('未找到对应物品')
      newItem.value.itemName = ''
    }
  } catch (error) {
    console.error('获取物品失败:', error)
    alert('查询失败，请检查ID是否正确')
  } finally {
    nameLoading.value = false
  }
}

const newItem = ref({
  itemId: 1,
  itemName: '',
  itemPrice: 1,
  itemCount: 1,
  shopId: 1,
  menuId: 1,
  subMenuId: 1,
})



// 提交新增物品
const submitNewItem = async () => {
  try {
    console.log('提交新增物品:', newItem)
    await AddShopItem(newItem.value.shopId, newItem.value.menuId, newItem.value.subMenuId, newItem.value.itemId, newItem.value.itemCount, newItem.value.itemPrice)
    console.log('提交新增物品:', newItem.value)
    // 刷新当前店铺数据
    await fetchData()
    console.log("fetchData success")
    console.log("selectedShopItem success", shopIdCache, menuIdCache, subMenuIdCache)
    selectedShopItem(shopIdCache, menuIdCache, subMenuIdCache)
    showAddModal.value = false
    alert('新增成功！')
  } catch (error) {
    console.error('新增失败:', error)
    alert('新增失败，请检查数据格式')
  }
}

// 搜索物品
const searchItemData = ref([])
function searchItem(itemName) {
  console.log("itemName",itemName)
  SearchItem(itemName).then((response) => {
    console.log('itemData:', response)
    searchItemData.value = response
  })
}



// 商店数据
let shopIdCache = ref(0)
let menuIdCache = ref(0)
let subMenuIdCache = ref(0)

function selectMenu(shopId) {
  shopIdCache = shopId
  GetMenuItems(shopId).then((response) => {
    menuData.value = response
  })
}

function selectedSubMenu(shopId, menuId) {
  menuIdCache = menuId
  GetSubMenus(shopId, menuId).then((response) => {
    subMenuData.value = response
  })
}

function selectedShopItem(shopId, menuId, subMenuId) {
  subMenuIdCache = subMenuId
  GetShopItems(shopId, menuId, subMenuId).then((response) => {
    console.log('itemData:', response)
    itemData.value = response
  })
}
// ------------------------------加载商店数据 end-----------------------------------
</script>

<style scoped>

.copy-feedback {
  position: absolute;
  top: 30px;
  left: 50%;
  transform: translateX(-50%);
  background: #333;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  animation: fadeIn0ut 2s;
}
.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: black;
}

.form-group input {
  width: 50%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}


.copyable-container {
  position: relative;
  display: grid;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  background: white;
  transition: background 0.2s;
}

.copyable-container:hover {
  background: black;
}



@keyframes fadeInOut {
  0% { opacity: 0; }
  20% { opacity: 0; }
  80% { opacity: 1; }
  100% { opacity: 0; }
}
.corner-add-modal {
  position: fixed;
  top: 20px;      /* 距顶部距离 */
  right: 20px;    /* 距右侧距离 */
  width: 350px;   /* 固定宽度 */
  max-height: calc(100vh - 40px); /* 最大高度 */
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.corner-modal {
  position: fixed;
  top: 20px;      /* 距顶部距离 */
  left: 20px;    /* 距右侧距离 */
  width: 350px;   /* 固定宽度 */
  max-height: calc(100vh - 40px); /* 最大高度 */
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #eee;
}


.modal-content {
  padding: 16px;
  overflow-y: auto;
  max-height: 40vh;
}



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
  overflow: auto;
}

.product-item {
  border: 1px solid #34495e;
  padding: 15px;
  border-radius: 6px;
  transition: transform 0.2s;
}

.product-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.item-name {
  font-weight: 500;
  margin-bottom: 8px;
}

.item-price {
  color: #e74c3c;
  font-size: 0.9em;
}


</style>



