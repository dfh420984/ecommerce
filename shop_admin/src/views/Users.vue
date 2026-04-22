<template>
  <div class="users">
    <el-card>
      <template #header>
        <span>用户列表</span>
      </template>

      <el-table :data="tableData" style="width: 100%">
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

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const loadData = async () => {
  const res = await request({
    url: '/admin/users',
    method: 'get',
    params: {
      page: pagination.page,
      page_size: pagination.page_size
    }
  })
  tableData.value = res.data?.list || []
  pagination.total = res.data?.total || 0
}

onMounted(() => {
  loadData()
})
</script>
