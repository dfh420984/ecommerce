<template>
  <div class="banners">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>轮播图列表</span>
          <el-button type="primary" @click="handleAdd">添加轮播图</el-button>
        </div>
      </template>

      <el-table :data="tableData" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="image" label="图片" width="150">
          <template #default="{ row }">
            <el-image v-if="row.image" :src="row.image" style="width: 100px; height: 50px" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column label="链接类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.link_type === 0" type="info">不跳转</el-tag>
            <el-tag v-else-if="row.link_type === 1" type="success">商品详情</el-tag>
            <el-tag v-else-if="row.link_type === 2" type="warning">分类页面</el-tag>
            <el-tag v-else-if="row.link_type === 3" type="primary">自定义链接</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="链接/目标" min-width="150">
          <template #default="{ row }">
            <span v-if="row.link_type === 1 && row.target_id">商品ID: {{ row.target_id }}</span>
            <span v-else-if="row.link_type === 2 && row.target_id">分类ID: {{ row.target_id }}</span>
            <span v-else-if="row.link_type === 3 && row.link">{{ row.link }}</span>
            <span v-else-if="row.link">{{ row.link }}</span>
            <span v-else style="color: #999">无</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
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
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑轮播图' : '添加轮播图'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="图片URL" prop="image">
          <el-input v-model="form.image" placeholder="请输入图片URL" />
        </el-form-item>
        <el-form-item label="链接地址">
          <el-input v-model="form.link" placeholder="请输入链接地址" />
        </el-form-item>
        <el-form-item label="链接类型">
          <el-select v-model="form.link_type" placeholder="请选择链接类型" style="width: 100%">
            <el-option label="不跳转" :value="0" />
            <el-option label="商品详情" :value="1" />
            <el-option label="分类页面" :value="2" />
            <el-option label="自定义链接" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标ID" v-if="form.link_type === 1 || form.link_type === 2">
          <el-input-number v-model="form.target_id" :min="0" placeholder="商品ID或分类ID" style="width: 100%" />
          <div style="font-size: 12px; color: #999; margin-top: 5px">
            {{ form.link_type === 1 ? '请输入商品ID' : '请输入分类ID' }}
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
import { getBanners, createBanner, updateBanner, deleteBanner } from '@/api/shop'

const tableData = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const editId = ref(null)

const form = ref({
  title: '',
  image: '',
  link: '',
  link_type: 0,
  target_id: 0,
  sort: 0,
  status: 1
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  image: [{ required: true, message: '请输入图片URL', trigger: 'blur' }]
}

const loadData = async () => {
  const res = await getBanners()
  tableData.value = res.data || []
}

const handleAdd = () => {
  form.value = { title: '', image: '', link: '', link_type: 0, target_id: 0, sort: 0, status: 1 }
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row) => {
  form.value = { ...row }
  isEdit.value = true
  editId.value = row.id
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该轮播图吗?', '提示', { type: 'warning' })
    await deleteBanner(row.id)
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
    if (isEdit.value) {
      await updateBanner(editId.value, form.value)
    } else {
      await createBanner(form.value)
    }
    ElMessage.success(isEdit.value ? '编辑成功' : '添加成功')
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
