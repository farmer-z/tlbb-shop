<template>
  <main>
    <div class="container">
      <button @click="fetchData">获取数据</button>

      <!-- 添加加载状态提示 -->
      <div v-if="loading">加载中...</div>
      <div v-else-if="error" class="error">{{ errorMessage }}</div>

      <!-- 调整表格渲染逻辑 -->
      <table v-if="tableData.length > 0 && tableData[0].length > 0">
        <tr v-for="(row, rowIndex) in tableData" :key="rowIndex">
          <td v-for="(item, colIndex) in row" :key="`${rowIndex}-${colIndex}`">
            {{ item }}
          </td>
        </tr>
      </table>
      <div v-else>暂无数据</div>
    </div>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { GetShopTable } from '../../wailsjs/go/main/App'

const tableData = ref([[]])
const loading = ref(false)
const error = ref(false)
const errorMessage = ref('')

const fetchData = async () => {
  try {
    loading.value = true
    error.value = false

    // 正确使用响应式赋值
    const response = await GetShopTable()

    // 调试输出原始数据
    console.log('原始响应数据:', response)

    // 确保数据结构为二维数组
    if (Array.isArray(response) && response.every(Array.isArray)) {
      tableData.value = response
    } else {
      throw new Error('返回数据格式不正确')
    }
  } catch (err) {
    error.value = true
    errorMessage.value = `获取数据失败: ${err.message}`
    console.error('请求错误:', err)
  } finally {
    loading.value = false
  }
}
</script>