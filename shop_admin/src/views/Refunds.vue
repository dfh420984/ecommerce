<template>
  <div class="refunds-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>退款管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 150px">
            <el-option label="待审核" :value="1" />
            <el-option label="已通过" :value="2" />
            <el-option label="已拒绝" :value="3" />
            <el-option label="退款中" :value="4" />
            <el-option label="已退款" :value="5" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order.order_no" label="订单号" width="180" />
        <el-table-column prop="user.nickname" label="用户" width="120" />
        <el-table-column prop="refund_type" label="退款类型" width="100">
          <template #default="{ row }">
            {{ row.refund_type === 'refund_only' ? '仅退款' : '换货' }}
          </template>
        </el-table-column>
        <el-table-column prop="refund_amount" label="退款金额" width="120">
          <template #default="{ row }">
            ¥{{ row.refund_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="退款原因" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button 
              v-if="row.status === 1" 
              size="small" 
              type="success"
              @click="handleApprove(row)"
            >
              通过
            </el-button>
            <el-button 
              v-if="row.status === 1" 
              size="small" 
              type="danger"
              @click="handleReject(row)"
            >
              拒绝
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="退款详情" width="700px">
      <el-descriptions :column="2" border v-if="currentRefund">
        <el-descriptions-item label="退款单号">{{ currentRefund.id }}</el-descriptions-item>
        <el-descriptions-item label="订单号">{{ currentRefund.order?.order_no }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentRefund.user?.nickname }}</el-descriptions-item>
        <el-descriptions-item label="退款类型">
          {{ currentRefund.refund_type === 'refund_only' ? '仅退款' : '换货' }}
        </el-descriptions-item>
        <el-descriptions-item label="退款金额">¥{{ currentRefund.refund_amount?.toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentRefund.status)">
            {{ getStatusText(currentRefund.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="申请时间" :span="2">
          {{ formatTime(currentRefund.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="退款原因" :span="2">
          {{ currentRefund.reason }}
        </el-descriptions-item>
        <el-descriptions-item label="凭证图片" :span="2" v-if="currentRefund.images && currentRefund.images.length > 0">
          <el-image
            v-for="(img, index) in currentRefund.images"
            :key="index"
            :src="getImageUrl(img)"
            :preview-src-list="currentRefund.images.map(i => getImageUrl(i))"
            style="width: 100px; height: 100px; margin-right: 10px"
            fit="cover"
          />
        </el-descriptions-item>
        <el-descriptions-item label="审核意见" :span="2" v-if="currentRefund.handler_reply">
          {{ currentRefund.handler_reply }}
        </el-descriptions-item>
        <el-descriptions-item label="处理人" v-if="currentRefund.handler">
          {{ currentRefund.handler.nickname }}
        </el-descriptions-item>
        <el-descriptions-item label="退款完成时间" v-if="currentRefund.refunded_at">
          {{ formatTime(currentRefund.refunded_at) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog v-model="approveVisible" title="审核退款" width="500px">
      <el-form :model="approveForm" label-width="100px">
        <el-form-item label="审核意见">
          <el-input
            v-model="approveForm.handler_reply"
            type="textarea"
            :rows="4"
            placeholder="请输入审核意见"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmApprove" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 拒绝对话框 -->
    <el-dialog v-model="rejectVisible" title="拒绝退款" width="500px">
      <el-form :model="rejectForm" label-width="100px">
        <el-form-item label="拒绝原因" required>
          <el-input
            v-model="rejectForm.handler_reply"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmReject" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { getImageUrl } from '@/utils/image'

const loading = ref(false)
const tableData = ref([])
const detailVisible = ref(false)
const approveVisible = ref(false)
const rejectVisible = ref(false)
const submitting = ref(false)
const currentRefund = ref(null)

const searchForm = reactive({
  order_no: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const approveForm = reactive({
  handler_reply: '审核通过，正在为您办理退款'
})

const rejectForm = reactive({
  handler_reply: ''
})

let currentRefundId = null

// 获取退款列表
const fetchRefunds = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...searchForm
    }
    
    const res = await request.get('/admin/refunds', { params })
    tableData.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error(error.message || '获取退款列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchRefunds()
}

// 重置
const handleReset = () => {
  searchForm.order_no = ''
  searchForm.status = ''
  pagination.page = 1
  fetchRefunds()
}

// 查看详情
const handleView = (row) => {
  currentRefund.value = row
  detailVisible.value = true
}

// 审核通过
const handleApprove = (row) => {
  currentRefundId = row.id
  approveForm.handler_reply = '审核通过，正在为您办理退款'
  approveVisible.value = true
}

// 确认通过
const confirmApprove = async () => {
  if (!approveForm.handler_reply) {
    ElMessage.warning('请输入审核意见')
    return
  }

  submitting.value = true
  try {
    await request.post(`/admin/refunds/${currentRefundId}/approve`, approveForm)
    ElMessage.success('审核通过')
    approveVisible.value = false
    fetchRefunds()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 拒绝
const handleReject = (row) => {
  currentRefundId = row.id
  rejectForm.handler_reply = ''
  rejectVisible.value = true
}

// 确认拒绝
const confirmReject = async () => {
  if (!rejectForm.handler_reply) {
    ElMessage.warning('请输入拒绝原因')
    return
  }

  submitting.value = true
  try {
    await request.post(`/admin/refunds/${currentRefundId}/reject`, rejectForm)
    ElMessage.success('已拒绝退款申请')
    rejectVisible.value = false
    fetchRefunds()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 分页
const handleSizeChange = () => {
  fetchRefunds()
}

const handlePageChange = () => {
  fetchRefunds()
}

// 状态文本
const getStatusText = (status) => {
  const map = {
    1: '待审核',
    2: '已通过',
    3: '已拒绝',
    4: '退款中',
    5: '已退款'
  }
  return map[status] || '未知'
}

// 状态类型
const getStatusType = (status) => {
  const map = {
    1: 'warning',
    2: 'success',
    3: 'danger',
    4: 'info',
    5: ''
  }
  return map[status] || ''
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchRefunds()
})
</script>

<style scoped>
.refunds-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
