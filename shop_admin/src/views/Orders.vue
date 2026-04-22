<template>
  <div class="orders">
    <el-card>
      <template #header>
        <span>订单列表</span>
      </template>

      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="queryForm.order_no" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="queryForm.status" placeholder="请选择状态" clearable>
            <el-option label="待付款" :value="1" />
            <el-option label="待发货" :value="2" />
            <el-option label="配送中" :value="3" />
            <el-option label="已收货" :value="4" />
            <el-option label="已完成" :value="5" />
            <el-option label="已取消" :value="6" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column label="订单金额" width="120">
          <template #default="{ row }">
            ¥{{ row.pay_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="order_status" label="订单状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.order_status)">
              {{ getStatusText(row.order_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pay_status" label="支付状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.pay_status === 1 ? 'success' : 'info'">
              {{ row.pay_status === 1 ? '已支付' : '未支付' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="consignee" label="收货人" width="100" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="address" label="收货地址" />
        <el-table-column prop="created_at" label="下单时间" width="180" />
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'

const tableData = ref([])
const queryForm = reactive({
  order_no: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const getStatusType = (status) => {
  const types = { 1: 'warning', 2: 'primary', 3: 'primary', 4: 'success', 5: 'success', 6: 'info' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 1: '待付款', 2: '待发货', 3: '配送中', 4: '已收货', 5: '已完成', 6: '已取消', 7: '退款中', 8: '已退款' }
  return texts[status] || '未知'
}

const loadData = async () => {
  const res = await request({
    url: '/user/orders',
    method: 'get',
    params: {
      page: pagination.page,
      page_size: pagination.page_size,
      status: queryForm.status
    }
  })
  tableData.value = res.data?.list || []
  pagination.total = res.data?.total || 0
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  queryForm.order_no = ''
  queryForm.status = ''
  handleSearch()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>
