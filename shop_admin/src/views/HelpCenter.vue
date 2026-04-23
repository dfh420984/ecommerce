<template>
  <div class="help-center">
    <!-- 统计卡片 -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409eff;">📚</div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.category_count || 0 }}</div>
              <div class="stat-label">分类数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67c23a;">❓</div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.question_count || 0 }}</div>
              <div class="stat-label">问题总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e6a23c;">✅</div>
            <div class="stat-info">
              <div class="stat-value">{{ statistics.active_question_count || 0 }}</div>
              <div class="stat-label">已启用问题</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-tabs v-model="activeTab" type="border-card">
      <!-- 分类管理 -->
      <el-tab-pane label="分类管理" name="categories">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>帮助中心分类</span>
              <el-button type="primary" @click="handleAddCategory">添加分类</el-button>
            </div>
          </template>

          <el-table :data="categories" border>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="分类名称" width="150" />
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditCategory(row)">编辑</el-button>
                <el-button type="danger" link @click="handleDeleteCategory(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- 问题管理 -->
      <el-tab-pane label="问题管理" name="questions">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>常见问题列表</span>
              <div>
                <el-select 
                  v-model="filterCategoryId" 
                  placeholder="选择分类" 
                  clearable
                  style="width: 150px; margin-right: 10px;"
                  @change="loadQuestions"
                >
                  <el-option 
                    v-for="cat in categories" 
                    :key="cat.id" 
                    :label="cat.name" 
                    :value="cat.id" 
                  />
                </el-select>
                <el-button type="primary" @click="handleAddQuestion">添加问题</el-button>
              </div>
            </div>
          </template>

          <el-table :data="questions" border>
            <el-table-column type="selection" width="55" />
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="category.name" label="分类" width="120" />
            <el-table-column prop="title" label="问题标题" min-width="250" show-overflow-tooltip />
            <el-table-column prop="answer" label="答案" min-width="300" show-overflow-tooltip>
              <template #default="{ row }">
                <div style="max-height: 60px; overflow: hidden; text-overflow: ellipsis;">
                  {{ row.answer }}
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditQuestion(row)">编辑</el-button>
                <el-button type="danger" link @click="handleDeleteQuestion(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div style="margin-top: 20px;">
            <el-button type="success" @click="handleBatchEnable">批量启用</el-button>
            <el-button type="warning" @click="handleBatchDisable">批量禁用</el-button>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 分类对话框 -->
    <el-dialog v-model="categoryDialogVisible" :title="isEditCategory ? '编辑分类' : '添加分类'" width="600px">
      <el-form :model="categoryForm" :rules="categoryRules" ref="categoryFormRef" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="categoryForm.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="categoryForm.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入分类描述" 
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="categoryForm.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="categoryForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitCategory">确定</el-button>
      </template>
    </el-dialog>

    <!-- 问题对话框 -->
    <el-dialog v-model="questionDialogVisible" :title="isEditQuestion ? '编辑问题' : '添加问题'" width="800px">
      <el-form :model="questionForm" :rules="questionRules" ref="questionFormRef" label-width="100px">
        <el-form-item label="所属分类" prop="category_id">
          <el-select v-model="questionForm.category_id" placeholder="请选择分类" style="width: 100%;">
            <el-option 
              v-for="cat in categories" 
              :key="cat.id" 
              :label="cat.name" 
              :value="cat.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="问题标题" prop="title">
          <el-input v-model="questionForm.title" placeholder="请输入问题标题" />
        </el-form-item>
        <el-form-item label="问题答案" prop="answer">
          <el-input 
            v-model="questionForm.answer" 
            type="textarea" 
            :rows="10"
            placeholder="请输入问题答案（支持多行文本）" 
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="questionForm.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="questionForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="questionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitQuestion">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getHelpCategories,
  createHelpCategory,
  updateHelpCategory,
  deleteHelpCategory,
  getHelpQuestions,
  createHelpQuestion,
  updateHelpQuestion,
  deleteHelpQuestion,
  batchUpdateQuestionsStatus,
  getHelpStatistics
} from '@/api/help'

const activeTab = ref('categories')
const statistics = ref({})
const categories = ref([])
const questions = ref([])
const filterCategoryId = ref(null)

// 分类相关
const categoryDialogVisible = ref(false)
const isEditCategory = ref(false)
const categoryFormRef = ref(null)
const editCategoryId = ref(null)
const categoryForm = ref({
  name: '',
  description: '',
  sort: 0,
  status: 1
})

const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

// 问题相关
const questionDialogVisible = ref(false)
const isEditQuestion = ref(false)
const questionFormRef = ref(null)
const editQuestionId = ref(null)
const questionForm = ref({
  category_id: null,
  title: '',
  answer: '',
  sort: 0,
  status: 1
})

const questionRules = {
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  title: [{ required: true, message: '请输入问题标题', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入问题答案', trigger: 'blur' }]
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

// 加载统计数据
const loadStatistics = async () => {
  try {
    const res = await getHelpStatistics()
    statistics.value = res.data || {}
  } catch (error) {
    console.error(error)
  }
}

// 加载分类
const loadCategories = async () => {
  try {
    const res = await getHelpCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

// 加载问题
const loadQuestions = async () => {
  try {
    const params = filterCategoryId.value ? { category_id: filterCategoryId.value } : {}
    const res = await getHelpQuestions(params)
    questions.value = res.data || []
  } catch (error) {
    console.error(error)
  }
}

// 分类 - 添加
const handleAddCategory = () => {
  categoryForm.value = { name: '', description: '', sort: 0, status: 1 }
  isEditCategory.value = false
  categoryDialogVisible.value = true
}

// 分类 - 编辑
const handleEditCategory = (row) => {
  categoryForm.value = {
    name: row.name,
    description: row.description,
    sort: row.sort,
    status: row.status
  }
  isEditCategory.value = true
  editCategoryId.value = row.id
  categoryDialogVisible.value = true
}

// 分类 - 删除
const handleDeleteCategory = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？如果该分类下有問題，将无法删除。', '提示', { type: 'warning' })
    await deleteHelpCategory(row.id)
    ElMessage.success('删除成功')
    loadCategories()
    loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

// 分类 - 提交
const handleSubmitCategory = async () => {
  const valid = await categoryFormRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    if (isEditCategory.value) {
      await updateHelpCategory(editCategoryId.value, categoryForm.value)
      ElMessage.success('编辑成功')
    } else {
      await createHelpCategory(categoryForm.value)
      ElMessage.success('添加成功')
    }
    categoryDialogVisible.value = false
    loadCategories()
    loadStatistics()
  } catch (error) {
    console.error(error)
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

// 问题 - 添加
const handleAddQuestion = () => {
  questionForm.value = { category_id: null, title: '', answer: '', sort: 0, status: 1 }
  isEditQuestion.value = false
  questionDialogVisible.value = true
}

// 问题 - 编辑
const handleEditQuestion = (row) => {
  questionForm.value = {
    category_id: row.category_id,
    title: row.title,
    answer: row.answer,
    sort: row.sort,
    status: row.status
  }
  isEditQuestion.value = true
  editQuestionId.value = row.id
  questionDialogVisible.value = true
}

// 问题 - 删除
const handleDeleteQuestion = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该问题吗？', '提示', { type: 'warning' })
    await deleteHelpQuestion(row.id)
    ElMessage.success('删除成功')
    loadQuestions()
    loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('删除失败')
    }
  }
}

// 问题 - 提交
const handleSubmitQuestion = async () => {
  const valid = await questionFormRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    if (isEditQuestion.value) {
      await updateHelpQuestion(editQuestionId.value, questionForm.value)
      ElMessage.success('编辑成功')
    } else {
      await createHelpQuestion(questionForm.value)
      ElMessage.success('添加成功')
    }
    questionDialogVisible.value = false
    loadQuestions()
    loadStatistics()
  } catch (error) {
    console.error(error)
    ElMessage.error(error.response?.data?.message || '操作失败')
  }
}

// 批量启用
const handleBatchEnable = async () => {
  const selectedRows = questions.value.filter(q => q.checked)
  if (selectedRows.length === 0) {
    ElMessage.warning('请先选择要启用的问题')
    return
  }

  try {
    await ElMessageBox.confirm(`确定启用选中的 ${selectedRows.length} 个问题吗？`, '提示', { type: 'warning' })
    const ids = selectedRows.map(q => q.id)
    await batchUpdateQuestionsStatus(ids, 1)
    ElMessage.success('启用成功')
    loadQuestions()
    loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('操作失败')
    }
  }
}

// 批量禁用
const handleBatchDisable = async () => {
  const selectedRows = questions.value.filter(q => q.checked)
  if (selectedRows.length === 0) {
    ElMessage.warning('请先选择要禁用的问题')
    return
  }

  try {
    await ElMessageBox.confirm(`确定禁用选中的 ${selectedRows.length} 个问题吗？`, '提示', { type: 'warning' })
    const ids = selectedRows.map(q => q.id)
    await batchUpdateQuestionsStatus(ids, 0)
    ElMessage.success('禁用成功')
    loadQuestions()
    loadStatistics()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
      ElMessage.error('操作失败')
    }
  }
}

onMounted(() => {
  loadStatistics()
  loadCategories()
  loadQuestions()
})
</script>

<style scoped>
.help-center {
  padding: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  color: white;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
