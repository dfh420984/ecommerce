<template>
  <div class="users">
    <el-card>
      <template #header>
        <span>用户列表</span>
      </template>

      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="queryForm.username" placeholder="请输入用户名" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="queryForm.phone" placeholder="请输入手机号" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="用户类型">
          <el-select v-model="queryForm.user_type" placeholder="请选择类型" clearable style="width: 120px">
            <el-option label="全部" :value="0" />
            <el-option label="用户" :value="1" />
            <el-option label="管理员" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="请选择状态" clearable style="width: 120px">
            <el-option label="全部" :value="0" />
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="user_type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="row.user_type === 2 ? 'danger' : 'primary'">
              {{ row.user_type === 2 ? '管理员' : '用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.status === 1" 
              type="warning" 
              link 
              size="small"
              @click="handleToggleStatus(row, 0)"
            >
              禁用
            </el-button>
            <el-button 
              v-else 
              type="success" 
              link 
              size="small"
              @click="handleToggleStatus(row, 1)"
            >
              启用
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
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

    <!-- 编辑用户对话框 -->
    <el-dialog v-model="editVisible" title="编辑用户" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" disabled />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="editForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="editForm.email" placeholder="请输入邮箱" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="handleEditSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { exportToExcel } from '@/utils/export'

const tableData = ref([])
const loading = ref(false)
const exportLoading = ref(false)

const queryForm = reactive({
  username: '',
  phone: '',
  user_type: 0,
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await request({
      url: '/admin/users',
      method: 'get',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        username: queryForm.username || undefined,
        phone: queryForm.phone || undefined,
        user_type: queryForm.user_type || undefined,
        status: queryForm.status || undefined
      }
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  queryForm.username = ''
  queryForm.phone = ''
  queryForm.user_type = 0
  queryForm.status = ''
  handleSearch()
}

// 切换用户状态
const handleToggleStatus = async (row, newStatus) => {
  const statusText = newStatus === 1 ? '启用' : '禁用'
  
  try {
    await ElMessageBox.confirm(
      `确定要${statusText}用户 "${row.username || row.nickname}" 吗？`,
      '提示',
      { type: 'warning' }
    )
    
    await request({
      url: `/admin/users/${row.id}/status`,
      method: 'put',
      data: { status: newStatus }
    })
    
    ElMessage.success(`${statusText}成功`)
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error(`${statusText}失败`)
    }
  }
}

// 编辑用户
const editVisible = ref(false)
const editForm = reactive({
  id: null,
  username: '',
  nickname: '',
  phone: '',
  email: ''
})

const handleEdit = (row) => {
  editForm.id = row.id
  editForm.username = row.username
  editForm.nickname = row.nickname || ''
  editForm.phone = row.phone || ''
  editForm.email = row.email || ''
  editVisible.value = true
}

const handleEditSubmit = async () => {
  try {
    await request({
      url: `/admin/users/${editForm.id}`,
      method: 'put',
      data: {
        nickname: editForm.nickname,
        phone: editForm.phone,
        email: editForm.email
      }
    })
    
    ElMessage.success('更新成功')
    editVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
    ElMessage.error('更新失败')
  }
}

// 删除用户
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${row.username || row.nickname}" 吗？此操作不可恢复！`,
      '警告',
      { type: 'warning' }
    )
    
    await request({
      url: `/admin/users/${row.id}`,
      method: 'delete'
    })
    
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('删除失败')
    }
  }
}

// 导出Excel
const handleExport = async () => {
  exportLoading.value = true
  try {
    // 获取所有数据（不分页）
    const res = await request({
      url: '/admin/users',
      method: 'get',
      params: {
        page: 1,
        page_size: 10000,
        username: queryForm.username || undefined,
        phone: queryForm.phone || undefined,
        user_type: queryForm.user_type || undefined,
        status: queryForm.status || undefined
      }
    })
    
    const users = res.data?.list || []
    
    if (users.length === 0) {
      ElMessage.warning('没有可导出的数据')
      return
    }

    // 定义列映射，自定义导出字段和标题
    const columnMap = {
      id: 'ID',
      username: '用户名',
      nickname: '昵称',
      phone: '手机号',
      email: '邮箱',
      user_type: '用户类型',
      status: '状态',
      created_at: '注册时间'
    }

    // 转换用户类型和状态为文本
    const exportData = users.map(user => ({
      ...user,
      user_type: user.user_type === 2 ? '管理员' : '用户',
      status: user.status === 1 ? '启用' : '禁用'
    }))

    // 生成文件名（包含时间戳）
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
    const fileName = `用户列表_${timestamp}`

    // 调用导出函数
    exportToExcel(exportData, fileName, '用户列表', columnMap)
    
    ElMessage.success(`成功导出 ${users.length} 条用户数据`)
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
</style>
