<template>
  <main>
    <div id="result" class="result">
      <table>
        <tr v-for="(row, rowIndex) in tableData" :key="rowIndex">
          <td v-for="(item, colIndex) in row" :key="colIndex">
            {{ item }}
          </td>
        </tr>
      </table>
    </div>
  </main>
</template>

<script setup>
import {onMounted, ref} from 'vue'
import {GetShopTable} from '../../wailsjs/go/main/App'

let tableData = ref([[]])
const headers = ref([[]])
const newRow = ref([[]])

// 获取数据
const fetchData = async () => {
  try {
    GetShopTable().then(response => {
      tableData = response
      console.log(tableData)
    })
  } catch (error) {
    console.error('获取数据失败:', error)
  }
}

// 添加新行
const addRow = async () => {
  if (newRow.value.some(item => !item.trim())) {
    alert('请填写所有字段')
    return
  }

  try {

  } catch (error) {
    console.error('添加失败:', error)
  }
}


onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.container {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

</style>