<template>
  <div class="orders">
    <el-card>
      <template #header>
        <span>订单列表</span>
      </template>

      <el-form :inline="true" :model="queryForm" class="search-form">
      <el-form-item label="订单号">
        <el-input v-model="queryForm.order_no" placeholder="请输入订单号" clearable style="width: 200px" />
      </el-form-item>
      <el-form-item label="订单状态">
        <el-select v-model="queryForm.order_status" placeholder="请选择状态" clearable style="width: 120px">
          <el-option label="全部" :value="0" />
          <el-option label="待付款" :value="1" />
          <el-option label="待发货" :value="2" />
          <el-option label="配送中" :value="3" />
          <el-option label="已收货" :value="4" />
          <el-option label="已完成" :value="5" />
          <el-option label="已取消" :value="6" />
          <el-option label="退款中" :value="7" />
          <el-option label="已退款" :value="8" />
        </el-select>
      </el-form-item>
      <el-form-item label="支付状态">
        <el-select v-model="queryForm.pay_status" placeholder="请选择支付状态" clearable style="width: 120px">
          <el-option label="未支付" :value="1" />
          <el-option label="已支付" :value="2" />
          <el-option label="已退款" :value="3" />
        </el-select>
      </el-form-item>
        <br/><br/>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
          <el-button type="success" @click="handleExport" :loading="exportLoading">导出Excel</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="order-expand">
              <div class="expand-item">
                <span class="expand-label">收货地址：</span>
                <span class="expand-value">{{ row.province }}{{ row.city }}{{ row.district }}{{ row.address }}</span>
              </div>
              <div class="expand-item" v-if="row.express_company">
                <span class="expand-label">快递公司：</span>
                <span class="expand-value">{{ row.express_company }}</span>
              </div>
              <div class="expand-item" v-if="row.express_no">
                <span class="expand-label">快递单号：</span>
                <span class="expand-value">{{ row.express_no }}</span>
              </div>
              <div class="expand-item" v-if="row.remark">
                <span class="expand-label">备注：</span>
                <span class="expand-value">{{ row.remark }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column label="用户信息" width="120">
          <template #default="{ row }">
            <div v-if="row.user">
              <div>{{ row.user.username || row.user.nickname }}</div>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="订单金额" width="100">
          <template #default="{ row }">
            ¥{{ row.pay_amount?.toFixed(2) || '0.00' }}
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
              {{ getPayStatusText(row.pay_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="consignee" label="收货人" width="100" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="created_at" label="下单时间" width="180" />
        <el-table-column label="操作" width="250">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button 
              v-if="row.order_status === 2 && row.pay_status === 1" 
              type="success" 
              link 
              @click="handleShip(row)"
            >
              发货
            </el-button>
            <el-button 
              v-if="row.order_status === 3" 
              type="warning" 
              link 
              @click="handleUpdateStatus(row, 4)"
            >
              已收货
            </el-button>
            <el-button 
              v-if="row.order_status === 4" 
              type="info" 
              link 
              @click="handleUpdateStatus(row, 5)"
            >
              完成
            </el-button>
          </template>
        </el-table-column>
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

    <!-- 订单详情对话框 -->
    <el-dialog v-model="detailVisible" title="订单详情" width="800px">
      <el-descriptions :column="2" border v-if="currentOrder">
        <el-descriptions-item label="订单号">{{ currentOrder.order_no }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusType(currentOrder.order_status)">
            {{ getStatusText(currentOrder.order_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付状态">
          <el-tag :type="currentOrder.pay_status === 1 ? 'success' : 'info'">
            {{ getPayStatusText(currentOrder.pay_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">
          {{ getPayTypeText(currentOrder.pay_type) }}
        </el-descriptions-item>
        <el-descriptions-item label="订单金额">¥{{ currentOrder.total_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="实付金额">¥{{ currentOrder.pay_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="收货人">{{ currentOrder.consignee }}</el-descriptions-item>
        <el-descriptions-item label="手机号">{{ currentOrder.phone }}</el-descriptions-item>
        <el-descriptions-item label="收货地址" :span="2">
          {{ currentOrder.province }}{{ currentOrder.city }}{{ currentOrder.district }}{{ currentOrder.address }}
        </el-descriptions-item>
        <el-descriptions-item label="快递公司" v-if="currentOrder.express_company">
          {{ currentOrder.express_company }}
        </el-descriptions-item>
        <el-descriptions-item label="快递单号" v-if="currentOrder.express_no">
          {{ currentOrder.express_no }}
        </el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ currentOrder.created_at }}</el-descriptions-item>
        <el-descriptions-item label="支付时间" v-if="currentOrder.pay_time">
          {{ currentOrder.pay_time }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2" v-if="currentOrder.remark">
          {{ currentOrder.remark }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider>商品列表</el-divider>
      <el-table :data="currentOrder?.items || []" border>
        <el-table-column prop="product_name" label="商品名称" />
        <el-table-column label="商品图片" width="100">
          <template #default="{ row }">
            <el-image 
              v-if="row.product_image" 
              :src="row.product_image" 
              style="width: 60px; height: 60px" 
              fit="cover" 
            />
          </template>
        </el-table-column>
        <el-table-column prop="price" label="单价" width="100">
          <template #default="{ row }">¥{{ row.price?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column prop="subtotal" label="小计" width="100">
          <template #default="{ row }">¥{{ row.subtotal?.toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 发货对话框 -->
    <el-dialog v-model="shipVisible" title="订单发货" width="500px">
      <el-form :model="shipForm" label-width="100px">
        <el-form-item label="快递公司" required>
          <el-input v-model="shipForm.express_company" placeholder="请输入快递公司名称" />
        </el-form-item>
        <el-form-item label="快递单号" required>
          <el-input v-model="shipForm.express_no" placeholder="请输入快递单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" @click="handleShipSubmit">确定发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrders, getOrder, shipOrder, updateOrderStatus } from '@/api/shop'
import { exportToExcel } from '@/utils/export'

const tableData = ref([])
const loading = ref(false)
const exportLoading = ref(false)
const queryForm = reactive({
  order_no: '',
  order_status: 0,
  pay_status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

// 订单详情
const detailVisible = ref(false)
const currentOrder = ref(null)

// 发货
const shipVisible = ref(false)
const shipForm = reactive({
  express_company: '',
  express_no: ''
})
const currentOrderId = ref(null)

const getStatusType = (status) => {
  const types = { 1: 'warning', 2: 'primary', 3: 'primary', 4: 'success', 5: 'success', 6: 'info', 7: 'danger', 8: 'info' }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = { 1: '待付款', 2: '待发货', 3: '配送中', 4: '已收货', 5: '已完成', 6: '已取消', 7: '退款中', 8: '已退款' }
  return texts[status] || '未知'
}

const getPayStatusText = (status) => {
  const texts = { 0: '未支付', 1: '已支付', 2: '已退款' }
  return texts[status] || '未知'
}

const getPayTypeText = (type) => {
  const texts = { 0: '未支付', 1: '微信支付', 2: '支付宝', 3: '银行卡' }
  return texts[type] || '未知'
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getOrders({
      page: pagination.page,
      page_size: pagination.page_size,
      order_no: queryForm.order_no || undefined,
      order_status: queryForm.order_status || undefined,
      pay_status: queryForm.pay_status || undefined
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error(error)
    ElMessage.error('获取订单列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  queryForm.order_no = ''
  queryForm.order_status = ''
  queryForm.pay_status = ''
  handleSearch()
}

// 查看订单详情
const handleView = async (row) => {
  try {
    const res = await getOrder(row.id)
    currentOrder.value = res.data
    detailVisible.value = true
  } catch (error) {
    console.error(error)
    ElMessage.error('获取订单详情失败')
  }
}

// 发货
const handleShip = (row) => {
  currentOrderId.value = row.id
  shipForm.express_company = ''
  shipForm.express_no = ''
  shipVisible.value = true
}

const handleShipSubmit = async () => {
  if (!shipForm.express_company || !shipForm.express_no) {
    ElMessage.warning('请填写快递公司和快递单号')
    return
  }

  try {
    await shipOrder(currentOrderId.value, shipForm)
    ElMessage.success('发货成功')
    shipVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
    ElMessage.error('发货失败')
  }
}

// 更新订单状态
const handleUpdateStatus = async (row, status) => {
  try {
    await ElMessageBox.confirm(
      `确定要将订单状态更新为"${getStatusText(status)}"吗？`,
      '提示',
      { type: 'warning' }
    )
    
    await updateOrderStatus(row.id, { order_status: status })
    ElMessage.success('更新成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('更新失败')
    }
  }
}

// 导出Excel
const handleExport = async () => {
  exportLoading.value = true
  try {
    // 获取所有数据（不分页）
    const res = await getOrders({
      page: 1,
      page_size: 10000,
      order_no: queryForm.order_no || undefined,
      order_status: queryForm.order_status || undefined,
      pay_status: queryForm.pay_status || undefined
    })
    
    const orders = res.data?.list || []
    
    if (orders.length === 0) {
      ElMessage.warning('没有可导出的数据')
      return
    }

    // 定义列映射，自定义导出字段和标题
    const columnMap = {
      id: 'ID',
      order_no: '订单号',
      consignee: '收货人',
      phone: '手机号',
      total_amount: '订单金额',
      pay_amount: '实付金额',
      order_status: '订单状态',
      pay_status: '支付状态',
      province: '省份',
      city: '城市',
      district: '区县',
      address: '详细地址',
      express_company: '快递公司',
      express_no: '快递单号',
      created_at: '下单时间',
      pay_time: '支付时间'
    }

    // 转换订单状态和支付状态为文本
    const exportData = orders.map(order => ({
      ...order,
      order_status: getStatusText(order.order_status),
      pay_status: getPayStatusText(order.pay_status),
      total_amount: order.total_amount ? `¥${order.total_amount.toFixed(2)}` : '¥0.00',
      pay_amount: order.pay_amount ? `¥${order.pay_amount.toFixed(2)}` : '¥0.00'
    }))

    // 生成文件名（包含时间戳）
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
    const fileName = `订单列表_${timestamp}`

    // 调用导出函数
    exportToExcel(exportData, fileName, '订单列表', columnMap)
    
    ElMessage.success(`成功导出 ${orders.length} 条订单数据`)
  } catch (error) {
    console.error(error)
    ElMessage.error('导出失败')
  } finally {
    exportLoading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 0;
}

.order-expand {
  padding: 20px 40px;
  background-color: #fafafa;
}

.expand-item {
  display: flex;
  align-items: flex-start;
  padding: 8px 0;
  border-bottom: 1px dashed #e8e8e8;
}

.expand-item:last-child {
  border-bottom: none;
}

.expand-label {
  flex-shrink: 0;
  min-width: 100px;
  font-weight: 600;
  color: #606266;
  font-size: 14px;
}

.expand-value {
  flex: 1;
  color: #303133;
  font-size: 14px;
  line-height: 1.6;
  word-break: break-all;
}
</style>
