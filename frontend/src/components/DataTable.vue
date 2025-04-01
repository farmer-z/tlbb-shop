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
      <div v-else class="data-container">
        <!-- 第一层：城市 -->
        <div
            v-for="(city, cityKey) in tableData"
            :key="`city-${cityKey}`"
            class="city-card"
        >
          <h2 class="city-title">城市 {{ cityKey }}</h2>

          <!-- 第二层：区域 -->
          <div
              v-for="(region, regionKey) in city"
              :key="`region-${cityKey}-${regionKey}`"
              class="region-section"
          >
            <h3 class="region-title">区域 {{ regionKey }}</h3>

            <!-- 第三层：店铺 -->
            <div
                v-for="(shop, shopKey) in region"
                :key="`shop-${cityKey}-${regionKey}-${shopKey}`"
                class="shop-container"
            >
              <div class="shop-header">
                <h4>店铺 {{ shopKey }}</h4>
                <span class="item-count">(共 {{ shop.length }} 件商品)</span>
              </div>

              <!-- 商品展示 -->
              <div class="items-grid">
                <div
                    v-for="(item, index) in shop"
                    :key="`item-${cityKey}-${regionKey}-${shopKey}-${index}`"
                    class="item-card"
                    :class="{ 'special-item': item.itemSpecialType !== 0 }"
                >
                  <!-- 颜色标签 -->
                  <div
                      v-if="item.itemDisplayColor"
                      class="color-tag"
                      :style="{ backgroundColor: item.itemDisplayColor }"
                  ></div>

                  <div class="item-content">
                    <div class="item-name">
                      {{ item.itemName }}
                      <span v-if="item.itemDiscount < 100" class="discount-badge">
                    {{ 100 - item.itemDiscount }}% OFF
                  </span>
                    </div>

                    <div class="item-details">
                      <div class="price-section">
                    <span class="original-price" v-if="item.itemDiscount < 100">
                      ¥{{ (item.itemPrice / 100).toFixed(2) }}
                    </span>
                        <span class="discounted-price">
                      ¥{{ (item.itemPrice * item.itemDiscount / 10000).toFixed(2) }}
                    </span>
                      </div>
                      <div class="stock-info">
                        库存: {{ item.itemCount }}
                      </div>
                    </div>

                    <!-- 特殊商品标识 -->
                    <div  class="special-tag">
                      商品特殊标记：{{ item.itemSpecialType }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      </div>
  </main>
</template>

<style scoped>
/* 添加可视化样式 */
.shop-panel {
  border: 1px solid #ccc;
  padding: 1rem;
  margin: 1rem 0;
}

.menu-section {
  margin-left: 1.5rem;
}

.item-card {
  padding: 0.5rem;
  margin: 0.5rem 0;
}
</style>


<script setup>
import {reactive, ref} from 'vue'
import {GetShopTable} from '../../wailsjs/go/main/App'

let tableData = reactive(new Map())
const loading = ref(false)
const error = ref(false)
const errorMessage = ref('')
const fetchData = async () => {
  try {
    loading.value = true
    error.value = false

    // 正确使用响应式赋值
    const response = await GetShopTable()
    console.log('获取数据成功:', response)
    console.log('typeof response:', typeof response)
    const data = JSON.parse(response)
    tableData = Object.entries(data || {}).reduce((cities, [cityKey, cityVal]) => {
      cities[cityKey] = Object.entries(cityVal || {}).reduce((regions, [regionKey, regionVal]) => {
        regions[regionKey] = Object.entries(regionVal || {}).reduce((shops, [shopKey, shopVal]) => {
          // 过滤无效数据
          if (Array.isArray(shopVal)) {
            shops[shopKey] = shopVal.filter(item =>
                item?.itemId && item?.itemName
            )
          }
          return shops
        }, {})
        return regions
      }, {})
      return cities
    }, {})
    console.log('tableData:', tableData)


  } catch (err) {
    error.value = true
    errorMessage.value = `获取数据失败: ${err.message}`
    console.error('请求错误:', err)
  } finally {
    loading.value = false
  }
}
</script>

