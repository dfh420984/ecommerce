<template>
  <div class="reviews-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>评论管理</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="商品名称">
          <el-input v-model="searchForm.product_name" placeholder="请输入商品名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 150px">
            <el-option label="正常" :value="1" />
            <el-option label="隐藏" :value="0" />
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
        <el-table-column prop="user.nickname" label="用户" width="120">
          <template #default="{ row }">
            {{ row.is_anonymous === 1 ? '匿名用户' : (row.user?.nickname || '未知') }}
          </template>
        </el-table-column>
        <el-table-column prop="product.name" label="商品名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <el-rate v-model="row.rating" disabled show-score text-color="#ff9900" />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" min-width="250" show-overflow-tooltip />
        <el-table-column prop="images" label="图片" width="150">
          <template #default="{ row }">
            <el-image
              v-if="row.images && row.images.length > 0"
              :src="getImageUrl(row.images[0])"
              :preview-src-list="row.images.map(img => getImageUrl(img))"
              style="width: 50px; height: 50px"
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '正常' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="评价时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleView(row)">详情</el-button>
            <el-button 
              v-if="row.status === 1" 
              size="small" 
              type="warning"
              @click="handleHide(row)"
            >
              隐藏
            </el-button>
            <el-button 
              v-if="row.status === 0" 
              size="small" 
              type="success"
              @click="handleShow(row)"
            >
              显示
            </el-button>
            <el-button 
              v-if="!row.reply_content" 
              size="small" 
              type="primary"
              @click="handleReply(row)"
            >
              回复
            </el-button>
            <el-button 
              size="small" 
              type="danger"
              @click="handleDelete(row)"
            >
              删除
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
    <el-dialog v-model="detailVisible" title="评论详情" width="700px">
      <el-descriptions :column="2" border v-if="currentReview">
        <el-descriptions-item label="评论ID">{{ currentReview.id }}</el-descriptions-item>
        <el-descriptions-item label="订单ID">{{ currentReview.order_id }}</el-descriptions-item>
        <el-descriptions-item label="用户">
          {{ currentReview.is_anonymous === 1 ? '匿名用户' : (currentReview.user?.nickname || '未知') }}
        </el-descriptions-item>
        <el-descriptions-item label="商品">{{ currentReview.product?.name }}</el-descriptions-item>
        <el-descriptions-item label="评分">
          <el-rate v-model="currentReview.rating" disabled show-score text-color="#ff9900" />
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentReview.status === 1 ? 'success' : 'info'">
            {{ currentReview.status === 1 ? '正常' : '隐藏' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="是否匿名">
          {{ currentReview.is_anonymous === 1 ? '是' : '否' }}
        </el-descriptions-item>
        <el-descriptions-item label="评价时间" :span="2">
          {{ formatTime(currentReview.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="评论内容" :span="2">
          {{ currentReview.content }}
        </el-descriptions-item>
        <el-descriptions-item label="评价图片" :span="2" v-if="currentReview.images && currentReview.images.length > 0">
          <el-image
            v-for="(img, index) in currentReview.images"
            :key="index"
            :src="getImageUrl(img)"
            :preview-src-list="currentReview.images.map(i => getImageUrl(i))"
            style="width: 100px; height: 100px; margin-right: 10px"
            fit="cover"
          />
        </el-descriptions-item>
        <el-descriptions-item label="商家回复" :span="2" v-if="currentReview.reply_content">
          {{ currentReview.reply_content }}
        </el-descriptions-item>
        <el-descriptions-item label="回复时间" v-if="currentReview.reply_time">
          {{ formatTime(currentReview.reply_time) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 回复对话框 -->
    <el-dialog v-model="replyVisible" title="回复评论" width="600px">
      <el-form :model="replyForm" label-width="100px">
        <el-form-item label="评论内容">
          <div style="color: #666; padding: 10px; background: #f5f7fa; border-radius: 4px;">
            {{ currentReview?.content }}
          </div>
        </el-form-item>
        <el-form-item label="回复内容" required>
          <el-input
            v-model="replyForm.reply_content"
            type="textarea"
            :rows="6"
            placeholder="请输入回复内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmReply" :loading="submitting">确定</el-button>
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
const replyVisible = ref(false)
const submitting = ref(false)
const currentReview = ref(null)

const searchForm = reactive({
  product_name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const replyForm = reactive({
  reply_content: ''
})

// 加载评论列表
const loadReviews = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/reviews',
      method: 'GET',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        product_name: searchForm.product_name || undefined,
        status: searchForm.status !== '' ? searchForm.status : undefined
      }
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error(error)
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadReviews()
}

// 重置
const handleReset = () => {
  searchForm.product_name = ''
  searchForm.status = ''
  handleSearch()
}

// 查看详情
const handleView = (row) => {
  currentReview.value = row
  detailVisible.value = true
}

// 隐藏评论
const handleHide = async (row) => {
  try {
    await ElMessageBox.confirm('确定要隐藏该评论吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request({
      url: `/admin/reviews/${row.id}/status`,
      method: 'PUT',
      data: { status: 0 }
    })
    
    ElMessage.success('隐藏成功')
    loadReviews()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('操作失败')
    }
  }
}

// 显示评论
const handleShow = async (row) => {
  try {
    await request({
      url: `/admin/reviews/${row.id}/status`,
      method: 'PUT',
      data: { status: 1 }
    })
    
    ElMessage.success('显示成功')
    loadReviews()
  } catch (error) {
    console.error(error)
    ElMessage.error('操作失败')
  }
}

// 回复评论
const handleReply = (row) => {
  currentReview.value = row
  replyForm.reply_content = ''
  replyVisible.value = true
}

// 确认回复
const confirmReply = async () => {
  if (!replyForm.reply_content.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  submitting.value = true
  try {
    await request({
      url: `/admin/reviews/${currentReview.value.id}/reply`,
      method: 'POST',
      data: { reply_content: replyForm.reply_content }
    })
    
    ElMessage.success('回复成功')
    replyVisible.value = false
    loadReviews()
  } catch (error) {
    console.error(error)
    ElMessage.error('回复失败')
  } finally {
    submitting.value = false
  }
}

// 删除评论
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该评论吗？此操作不可恢复！', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error'
    })
    
    await request({
      url: `/admin/reviews/${row.id}`,
      method: 'DELETE'
    })
    
    ElMessage.success('删除成功')
    loadReviews()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('删除失败')
    }
  }
}

// 分页大小改变
const handleSizeChange = () => {
  loadReviews()
}

// 页码改变
const handlePageChange = () => {
  loadReviews()
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  const second = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hour}:${minute}:${second}`
}

onMounted(() => {
  loadReviews()
})
</script>

<style scoped>
.reviews-container {
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
