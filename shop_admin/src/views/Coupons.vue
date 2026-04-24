<template>
  <div class="coupons">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>优惠券管理</span>
          <el-button type="primary" @click="handleAdd">创建优惠券</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="优惠券名称">
          <el-input v-model="queryForm.name" placeholder="请输入名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryForm.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="优惠券名称" min-width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 1 ? 'success' : row.type === 2 ? 'warning' : 'info'">
              {{ row.type === 1 ? '满减券' : row.type === 2 ? '折扣券' : '无门槛' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="优惠信息" width="150">
          <template #default="{ row }">
            <div v-if="row.type === 1 || row.type === 3">
              ¥{{ row.discount_value.toFixed(2) }}
            </div>
            <div v-else>
              {{ (row.discount_value * 10).toFixed(1) }}折
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="min_amount" label="最低消费" width="100">
          <template #default="{ row }">
            ¥{{ row.min_amount.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column label="发放情况" width="120">
          <template #default="{ row }">
            {{ row.received_count }}/{{ row.total_count === 0 ? '不限' : row.total_count }}
          </template>
        </el-table-column>
        <el-table-column prop="per_user_limit" label="每人限领" width="90" />
        <el-table-column label="有效期" width="180">
          <template #default="{ row }">
            <div v-if="row.valid_type === 1">
              {{ formatTime(row.start_time) }} ~ {{ formatTime(row.end_time) }}
            </div>
            <div v-else>
              领取后{{ row.valid_days }}天
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="is_new_user" label="新人券" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_new_user === 1 ? 'success' : 'info'" size="small">
              {{ row.is_new_user === 1 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link @click="handleGrant(row)">赠送</el-button>
            <el-button 
              :type="row.status === 1 ? 'warning' : 'success'" 
              link 
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadData"
        @current-change="loadData"
        style="margin-top: 20px; text-align: right"
      />
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑优惠券' : '创建优惠券'" width="700px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="优惠券名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入优惠券名称" />
        </el-form-item>
        
        <el-form-item label="优惠券类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :label="1">满减券</el-radio>
            <el-radio :label="2">折扣券</el-radio>
            <el-radio :label="3">无门槛券</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="优惠值" prop="discount_value">
          <el-input-number 
            v-model="form.discount_value" 
            :min="0" 
            :precision="2"
            :step="form.type === 2 ? 0.1 : 1"
          />
          <span style="margin-left: 10px; color: #999">
            {{ form.type === 2 ? '(如0.8表示8折)' : '(元)' }}
          </span>
        </el-form-item>

        <el-form-item label="最低消费" prop="min_amount">
          <el-input-number v-model="form.min_amount" :min="0" :precision="2" />
          <span style="margin-left: 10px; color: #999">(元)</span>
        </el-form-item>

        <el-form-item v-if="form.type === 2" label="最大优惠" prop="max_discount">
          <el-input-number v-model="form.max_discount" :min="0" :precision="2" />
          <span style="margin-left: 10px; color: #999">(元，0表示不限制)</span>
        </el-form-item>

        <el-form-item label="发放总量" prop="total_count">
          <el-input-number v-model="form.total_count" :min="0" />
          <span style="margin-left: 10px; color: #999">(0表示不限)</span>
        </el-form-item>

        <el-form-item label="每人限领" prop="per_user_limit">
          <el-input-number v-model="form.per_user_limit" :min="1" />
        </el-form-item>

        <el-form-item label="有效期类型" prop="valid_type">
          <el-radio-group v-model="form.valid_type">
            <el-radio :label="1">固定时间</el-radio>
            <el-radio :label="2">领取后N天</el-radio>
          </el-radio-group>
        </el-form-item>

        <template v-if="form.valid_type === 1">
          <el-form-item label="开始时间" prop="start_time">
            <el-date-picker
              v-model="form.start_time"
              type="datetime"
              placeholder="选择开始时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="结束时间" prop="end_time">
            <el-date-picker
              v-model="form.end_time"
              type="datetime"
              placeholder="选择结束时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 100%"
            />
          </el-form-item>
        </template>

        <el-form-item v-else label="有效天数" prop="valid_days">
          <el-input-number v-model="form.valid_days" :min="1" />
          <span style="margin-left: 10px; color: #999">(天)</span>
        </el-form-item>

        <el-form-item label="是否新人券">
          <el-switch v-model="form.is_new_user" :active-value="1" :inactive-value="0" />
        </el-form-item>

        <el-form-item label="使用说明">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入使用说明" />
        </el-form-item>

        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 赠送对话框 -->
    <el-dialog v-model="grantDialogVisible" title="赠送优惠券" width="600px">
      <el-form :model="grantForm" label-width="100px">
        <el-form-item label="优惠券">
          <el-input :value="grantCouponName" disabled />
        </el-form-item>
        <el-form-item label="用户ID" prop="user_ids">
          <el-input
            v-model="grantForm.user_ids_str"
            type="textarea"
            :rows="5"
            placeholder="请输入用户ID，多个ID用逗号或换行分隔"
          />
          <div style="margin-top: 5px; color: #999; font-size: 12px">
            示例：1,2,3 或每行一个ID
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="grantDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleGrantSubmit">确定赠送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getCoupons, 
  createCoupon, 
  updateCoupon, 
  deleteCoupon, 
  updateCouponStatus,
  grantCouponToUser 
} from '@/api/shop'

const tableData = ref([])
const dialogVisible = ref(false)
const grantDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editId = ref(null)
const grantCouponName = ref('')
const currentCouponId = ref(null)

const queryForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const form = ref({
  name: '',
  type: 1,
  discount_value: 0,
  min_amount: 0,
  max_discount: 0,
  total_count: 0,
  per_user_limit: 1,
  valid_type: 1,
  start_time: '',
  end_time: '',
  valid_days: 7,
  status: 1,
  is_new_user: 0,
  description: ''
})

const grantForm = reactive({
  user_ids_str: ''
})

const rules = {
  name: [{ required: true, message: '请输入优惠券名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择优惠券类型', trigger: 'change' }],
  discount_value: [{ required: true, message: '请输入优惠值', trigger: 'blur' }],
  per_user_limit: [{ required: true, message: '请设置每人限领数量', trigger: 'blur' }],
  valid_type: [{ required: true, message: '请选择有效期类型', trigger: 'change' }]
}

const loadData = async () => {
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      ...queryForm
    }
    const res = await getCoupons(params)
    tableData.value = res.data.list || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  queryForm.name = ''
  queryForm.status = ''
  handleSearch()
}

const handleAdd = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  editId.value = row.id
  form.value = {
    name: row.name,
    type: row.type,
    discount_value: row.discount_value,
    min_amount: row.min_amount,
    max_discount: row.max_discount,
    total_count: row.total_count,
    per_user_limit: row.per_user_limit,
    valid_type: row.valid_type,
    start_time: row.start_time || '',
    end_time: row.end_time || '',
    valid_days: row.valid_days,
    status: row.status,
    is_new_user: row.is_new_user,
    description: row.description || ''
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    if (isEdit.value) {
      await updateCoupon(editId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await createCoupon(form.value)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadData()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    }
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该优惠券吗？', '提示', {
      type: 'warning'
    })
    await deleteCoupon(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleToggleStatus = async (row) => {
  try {
    const newStatus = row.status === 1 ? 0 : 1
    await updateCouponStatus(row.id, { status: newStatus })
    ElMessage.success(newStatus === 1 ? '已启用' : '已禁用')
    loadData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleGrant = (row) => {
  currentCouponId.value = row.id
  grantCouponName.value = row.name
  grantForm.user_ids_str = ''
  grantDialogVisible.value = true
}

const handleGrantSubmit = async () => {
  if (!grantForm.user_ids_str.trim()) {
    ElMessage.warning('请输入用户ID')
    return
  }

  try {
    // 解析用户ID
    const userIds = grantForm.user_ids_str
      .split(/[,\n]/)
      .map(id => parseInt(id.trim()))
      .filter(id => !isNaN(id))

    if (userIds.length === 0) {
      ElMessage.warning('请输入有效的用户ID')
      return
    }

    const res = await grantCouponToUser({
      coupon_id: currentCouponId.value,
      user_ids: userIds
    })

    ElMessage.success(`成功赠送${res.data.success_count}个，失败${res.data.failed_count}个`)
    
    if (res.data.failed_users && res.data.failed_users.length > 0) {
      console.warn('失败详情:', res.data.failed_users)
    }

    grantDialogVisible.value = false
    loadData()
  } catch (error) {
    ElMessage.error('赠送失败')
  }
}

const resetForm = () => {
  form.value = {
    name: '',
    type: 1,
    discount_value: 0,
    min_amount: 0,
    max_discount: 0,
    total_count: 0,
    per_user_limit: 1,
    valid_type: 1,
    start_time: '',
    end_time: '',
    valid_days: 7,
    status: 1,
    is_new_user: 0,
    description: ''
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const formatTime = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
