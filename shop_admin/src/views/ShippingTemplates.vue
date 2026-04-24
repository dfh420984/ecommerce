<template>
  <div class="shipping-templates">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>运费模板管理</span>
          <el-button type="primary" @click="handleAdd">添加模板</el-button>
        </div>
      </template>

      <el-table :data="tableData" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="模板名称" width="200" />
        <el-table-column label="是否默认" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_default === 1" type="success">是</el-tag>
            <el-tag v-else type="info">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="包邮类型" width="120">
          <template #default="{ row }">
            {{ getFreeShippingTypeText(row.free_shipping_type) }}
          </template>
        </el-table-column>
        <el-table-column label="包邮条件" width="150">
          <template #default="{ row }">
            <span v-if="row.free_shipping_type === 1 && row.free_amount > 0">满{{ row.free_amount }}元</span>
            <span v-else-if="row.free_shipping_type === 2 && row.free_quantity > 0">满{{ row.free_quantity }}件</span>
            <span v-else-if="row.free_shipping_type === 3">
              满{{ row.free_amount }}元或{{ row.free_quantity }}件
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="base_fee" label="基础运费" width="100">
          <template #default="{ row }">
            ¥{{ row.base_fee }}
          </template>
        </el-table-column>
        <el-table-column label="区域配置" width="100">
          <template #default="{ row }">
            {{ row.regions ? row.regions.length : 0 }}个地区
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">启用</el-tag>
            <el-tag v-else type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button 
              v-if="row.is_default !== 1" 
              type="warning" 
              link 
              @click="handleSetDefault(row)"
            >
              设为默认
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button 
              v-if="row.is_default !== 1" 
              type="danger" 
              link 
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑运费模板' : '添加运费模板'" 
      width="900px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入模板名称" />
        </el-form-item>
        
        <el-form-item label="是否默认" prop="is_default">
          <el-radio-group v-model="form.is_default">
            <el-radio :label="1">是</el-radio>
            <el-radio :label="0">否</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="包邮类型" prop="free_shipping_type">
          <el-radio-group v-model="form.free_shipping_type">
            <el-radio :label="1">满额包邮</el-radio>
            <el-radio :label="2">满件包邮</el-radio>
            <el-radio :label="3">满额或满件</el-radio>
            <el-radio :label="4">不包邮</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item 
          v-if="form.free_shipping_type === 1 || form.free_shipping_type === 3" 
          label="包邮金额" 
          prop="free_amount"
        >
          <el-input-number 
            v-model="form.free_amount" 
            :min="0" 
            :precision="2"
            placeholder="设置满多少元包邮"
          />
          <span style="margin-left: 10px; color: #909399;">元</span>
        </el-form-item>

        <el-form-item 
          v-if="form.free_shipping_type === 2 || form.free_shipping_type === 3" 
          label="包邮数量" 
          prop="free_quantity"
        >
          <el-input-number 
            v-model="form.free_quantity" 
            :min="0"
            placeholder="设置满多少件包邮"
          />
          <span style="margin-left: 10px; color: #909399;">件</span>
        </el-form-item>

        <el-form-item label="基础运费" prop="base_fee">
          <el-input-number 
            v-model="form.base_fee" 
            :min="0" 
            :precision="2"
            placeholder="未满足包邮条件时的运费"
          />
          <span style="margin-left: 10px; color: #909399;">元</span>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-divider content-position="left">特殊地区运费配置</el-divider>
        
        <el-form-item label="地区配置">
          <div style="margin-bottom: 10px;">
            <el-button type="primary" size="small" @click="addRegion">添加地区</el-button>
          </div>
          
          <el-table :data="form.regions" border size="small">
            <el-table-column label="省份" width="150">
              <template #default="{ row }">
                <el-input v-model="row.province" placeholder="如：新疆" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="城市" width="150">
              <template #default="{ row }">
                <el-input v-model="row.city" placeholder="如：乌鲁木齐市" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="区县" width="150">
              <template #default="{ row }">
                <el-input v-model="row.district" placeholder="可选" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="运费" width="150">
              <template #default="{ row }">
                <el-input-number 
                  v-model="row.fee" 
                  :min="0" 
                  :precision="2"
                  size="small"
                  style="width: 100%;"
                />
              </template>
            </el-table-column>
            <el-table-column label="包邮金额" width="150">
              <template #default="{ row }">
                <el-input-number 
                  v-model="row.free_amount" 
                  :min="0" 
                  :precision="2"
                  size="small"
                  style="width: 100%;"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ $index }">
                <el-button 
                  type="danger" 
                  link 
                  size="small"
                  @click="removeRegion($index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div style="margin-top: 10px; color: #909399; font-size: 12px;">
            提示：可以为特定地区设置不同的运费和包邮条件。如果不配置特殊地区，则使用基础运费。
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const tableData = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editId = ref(null)

const form = ref({
  name: '',
  is_default: 0,
  free_shipping_type: 1,
  free_amount: 0,
  free_quantity: 0,
  base_fee: 0,
  status: 1,
  regions: []
})

const rules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  free_shipping_type: [{ required: true, message: '请选择包邮类型', trigger: 'change' }]
}

// 获取包邮类型文本
const getFreeShippingTypeText = (type) => {
  const texts = {
    1: '满额包邮',
    2: '满件包邮',
    3: '满额或满件',
    4: '不包邮'
  }
  return texts[type] || '未知'
}

// 加载数据
const loadData = async () => {
  try {
    const res = await request.get('/admin/shipping/templates')
    tableData.value = res.data || []
  } catch (error) {
    console.error(error)
    ElMessage.error('获取运费模板列表失败')
  }
}

// 添加模板
const handleAdd = () => {
  form.value = {
    name: '',
    is_default: 0,
    free_shipping_type: 1,
    free_amount: 0,
    free_quantity: 0,
    base_fee: 0,
    status: 1,
    regions: []
  }
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑模板
const handleEdit = (row) => {
  form.value = {
    name: row.name,
    is_default: row.is_default,
    free_shipping_type: row.free_shipping_type,
    free_amount: row.free_amount || 0,
    free_quantity: row.free_quantity || 0,
    base_fee: row.base_fee || 0,
    status: row.status,
    regions: row.regions ? JSON.parse(JSON.stringify(row.regions)) : []
  }
  isEdit.value = true
  editId.value = row.id
  dialogVisible.value = true
}

// 添加地区
const addRegion = () => {
  form.value.regions.push({
    province: '',
    city: '',
    district: '',
    fee: 0,
    free_amount: 0,
    free_quantity: 0
  })
}

// 删除地区
const removeRegion = (index) => {
  form.value.regions.splice(index, 1)
}

// 设为默认
const handleSetDefault = async (row) => {
  try {
    await ElMessageBox.confirm('确定将该模板设为默认吗？', '提示', { type: 'warning' })
    await request.post(`/admin/shipping/templates/${row.id}/default`)
    ElMessage.success('设置成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('设置失败')
    }
  }
}

// 删除模板
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该运费模板吗？', '提示', { type: 'warning' })
    await request.delete(`/admin/shipping/templates/${row.id}`)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    if (isEdit.value) {
      await request.put(`/admin/shipping/templates/${editId.value}`, form.value)
      ElMessage.success('编辑成功')
    } else {
      await request.post('/admin/shipping/templates', form.value)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
    ElMessage.error(isEdit.value ? '编辑失败' : '添加失败')
  }
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
</style>
