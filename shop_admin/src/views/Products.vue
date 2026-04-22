<template>
  <div class="products">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>商品列表</span>
          <el-button type="primary" @click="handleAdd">添加商品</el-button>
        </div>
      </template>

      <el-form :inline="true" :model="queryForm" class="search-form">
        <el-form-item label="商品名称">
          <el-input v-model="queryForm.keyword" placeholder="请输入商品名称" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="queryForm.category_id" placeholder="请选择分类" clearable>
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.displayName || cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="category_id" label="分类" width="100">
          <template #default="{ row }">
            {{ getCategoryName(row.category_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            ¥{{ row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="sales" label="销量" width="80" />
        <el-table-column prop="is_online" label="上架" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_online === 1 ? 'success' : 'info'">
              {{ row.is_online === 1 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑商品' : '添加商品'" width="800px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.displayName || cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="原价" prop="original_price">
          <el-input-number v-model="form.original_price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="成本价" prop="cost">
          <el-input-number v-model="form.cost" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="商品图片">
          <ImageUpload v-model="form.images" />
        </el-form-item>
        <el-form-item label="商品描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="商品详情">
          <RichTextEditor v-model="form.content" />
        </el-form-item>
        <el-form-item label="上架">
          <el-switch v-model="form.is_online" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="推荐">
          <el-switch v-model="form.is_recommend" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="新品">
          <el-switch v-model="form.is_new" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getProducts, createProduct, updateProduct, deleteProduct, getCategories } from '@/api/shop'
import ImageUpload from '@/components/ImageUpload.vue'
import RichTextEditor from '@/components/RichTextEditor.vue'

const tableData = ref([])
const categories = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editId = ref(null)

const queryForm = reactive({
  keyword: '',
  category_id: ''
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

const form = ref({
  name: '',
  category_id: '',
  price: 0,
  original_price: 0,
  cost: 0,
  stock: 0,
  images: [],
  description: '',
  content: '',
  is_online: 1,
  is_recommend: 0,
  is_new: 0,
  sort: 0
})

const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }]
}

const loadCategories = async () => {
  try {
    const res = await getCategories()
    console.log('原始分类数据:', res.data)
    
    // 将树形结构展平为列表，并添加层级标识
    const flattenCategories = (cats, level = 0) => {
      let result = []
      cats.forEach(cat => {
        // 根据层级添加不同的前缀
        let prefix = ''
        if (level === 1) {
          prefix = '├── '
        } else if (level === 2) {
          prefix = '│   ├── '
        } else if (level > 2) {
          prefix = '│   '.repeat(level - 1) + '├── '
        }
        
        result.push({
          ...cat,
          displayName: prefix + cat.name
        })
        
        // 递归处理子分类
        if (cat.children && cat.children.length > 0) {
          result = result.concat(flattenCategories(cat.children, level + 1))
        }
      })
      return result
    }
    
    // 展平分类数据
    categories.value = res.data ? flattenCategories(res.data) : []
    console.log('展平后分类:', categories.value)
  } catch (error) {
    console.error('加载分类失败:', error)
    ElMessage.error('加载分类失败')
  }
}

const getCategoryName = (id) => {
  const cat = categories.value.find(c => c.id === id)
  return cat ? cat.name : '-'
}

const loadData = async () => {
  const res = await getProducts({
    page: pagination.page,
    page_size: pagination.page_size,
    keyword: queryForm.keyword,
    category_id: queryForm.category_id
  })
  tableData.value = res.data?.list || []
  pagination.total = res.data?.total || 0
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  queryForm.keyword = ''
  queryForm.category_id = ''
  handleSearch()
}

const handleAdd = () => {
  form.value = {
    name: '',
    category_id: '',
    price: 0,
    original_price: 0,
    cost: 0,
    stock: 0,
    images: [],
    description: '',
    content: '',
    is_online: 1,
    is_recommend: 0,
    is_new: 0,
    sort: 0
  }
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  form.value = {
    ...row,
    images: Array.isArray(row.images) ? row.images : []
  }
  isEdit.value = true
  editId.value = row.id
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该商品吗?', '提示', { type: 'warning' })
    await deleteProduct(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') console.error(error)
  }
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  const data = { ...form.value }

  try {
    if (isEdit.value) {
      await updateProduct(editId.value, data)
    } else {
      await createProduct(data)
    }
    ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
    dialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  loadCategories()
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
