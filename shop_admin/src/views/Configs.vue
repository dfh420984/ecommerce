<template>
  <div class="configs">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统配置</span>
          <el-button type="primary" @click="handleAdd">添加配置</el-button>
        </div>
      </template>

      <el-table :data="tableData" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="配置名称" width="200" />
        <el-table-column prop="value" label="配置值" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.value }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑配置' : '添加配置'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="配置名称" prop="name">
          <el-input 
            v-model="form.name" 
            placeholder="请输入配置名称（英文，如：site_name）" 
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <!-- 文本输入 -->
          <el-input 
            v-model="form.value" 
            type="textarea" 
            :rows="4"
            placeholder="请输入配置值（支持文本和图片URL）" 
            style="margin-bottom: 10px;"
          />
          <!-- 图片上传 -->
          <div style="margin-top: 10px;">
            <div style="margin-bottom: 8px; color: #606266; font-size: 14px;">上传图片：</div>
            <ImageUpload 
              v-model="uploadedImages"
              :max="9"
            />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              提示：上传图片后，图片URL会自动追加到配置值中
            </div>
          </div>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入配置描述" 
          />
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
import { getConfigs, createConfig, updateConfig, deleteConfig } from '@/api/config'
import ImageUpload from '@/components/ImageUpload.vue'

const tableData = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editId = ref(null)
const uploadedImages = ref([]) // 上传的图片列表

const form = ref({
  name: '',
  value: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入配置名称', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '配置名称只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  value: [{ required: true, message: '请输入配置值', trigger: 'blur' }]
}

const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit', 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit' 
  })
}

const loadData = async () => {
  try {
    const res = await getConfigs()
    tableData.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

const handleAdd = () => {
  form.value = { name: '', value: '', description: '' }
  uploadedImages.value = []
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  form.value = {
    name: row.name,
    value: row.value,
    description: row.description
  }
  // 提取已有的图片URL
  extractImagesFromValue(row.value)
  isEdit.value = true
  editId.value = row.id
  dialogVisible.value = true
}

// 从配置值中提取图片URL
const extractImagesFromValue = (value) => {
  if (!value) {
    uploadedImages.value = []
    return
  }
  
  // 匹配图片URL（支持常见图片格式）
  const imageRegex = /(https?:\/\/[^\s<>"']+(?:\.jpg|\.jpeg|\.png|\.gif|\.webp|\.bmp))/gi
  const matches = value.match(imageRegex) || []
  
  uploadedImages.value = matches.map(url => ({
    uid: Date.now() + Math.random(),
    name: url.split('/').pop(),
    url: url
  }))
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该配置吗?', '提示', { type: 'warning' })
    await deleteConfig(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') console.error(error)
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    // 如果有新上传的图片，将图片URL追加到配置值中
    let finalValue = form.value.value
    if (uploadedImages.value.length > 0) {
      const newImageUrls = uploadedImages.value
        .map(img => img.url)
        .filter(url => url && !finalValue.includes(url))
      
      if (newImageUrls.length > 0) {
        // 如果配置值不为空，添加换行符
        if (finalValue) {
          finalValue += '\n' + newImageUrls.join('\n')
        } else {
          finalValue = newImageUrls.join('\n')
        }
      }
    }

    const data = {
      name: form.value.name,
      value: finalValue,
      description: form.value.description
    }

    if (isEdit.value) {
      await updateConfig(editId.value, data)
      ElMessage.success('编辑成功')
    } else {
      await createConfig(data)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
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
