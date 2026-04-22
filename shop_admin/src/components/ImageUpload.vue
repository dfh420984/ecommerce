<template>
  <div class="image-upload">
    <el-upload
      action="/api/upload"
      :headers="uploadHeaders"
      list-type="picture-card"
      :file-list="fileList"
      :on-success="handleSuccess"
      :on-error="handleError"
      :before-upload="beforeUpload"
      :on-remove="handleRemove"
      :on-preview="handlePreview"
      accept="image/*"
    >
      <el-icon><Plus /></el-icon>
    </el-upload>

    <el-dialog v-model="dialogVisible" title="图片预览">
      <img w-full :src="previewUrl" alt="预览图片" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: [Array, String],
    default: () => []
  },
  max: {
    type: Number,
    default: 9
  }
})

const emit = defineEmits(['update:modelValue'])

const fileList = ref([])
const dialogVisible = ref(false)
const previewUrl = ref('')

// 上传请求头（携带 token）
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 初始化文件列表
const initFileList = (urls) => {
  if (urls && urls.length > 0) {
    return urls.map((url, index) => ({
      uid: Date.now() + index, // 使用时间戳确保唯一性
      name: `image-${index}.jpg`,
      url: url
    }))
  }
  return []
}

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (val) => {
    let urls = []
    
    // 处理字符串（单图）或数组（多图）
    if (typeof val === 'string') {
      urls = val ? [val] : []
    } else if (Array.isArray(val)) {
      urls = val
    }
    
    // 只在初始化或外部数据变化时更新，避免与上传操作冲突
    if (fileList.value.length === 0 || urls.length !== fileList.value.length) {
      fileList.value = initFileList(urls)
    }
  },
  { immediate: true }
)

// 上传前验证
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  
  // 检查上传数量限制
  if (fileList.value.length >= props.max) {
    ElMessage.error(`最多只能上传 ${props.max} 张图片!`)
    return false
  }
  
  return true
}

// 上传成功
const handleSuccess = (response, uploadFile, fileListParam) => {
  console.log('上传成功响应:', response)
  console.log('上传文件:', uploadFile)
  
  if (response.code === 0) {
    // 使用 fileListParam 更新本地 fileList
    const newFileList = fileListParam.map(file => ({
      uid: file.uid,
      name: file.name,
      url: file.response?.data?.url || file.url || response.data.url
    }))
    
    fileList.value = newFileList
    
    // 通知父组件更新
    const urls = newFileList.map(file => file.url).filter(url => url)
    console.log('更新图片列表:', urls)
    
    // 根据 max 属性决定返回类型
    if (props.max === 1) {
      // 单图模式：返回字符串
      emit('update:modelValue', urls[0] || '')
    } else {
      // 多图模式：返回数组
      emit('update:modelValue', urls)
    }
    
    ElMessage.success('上传成功')
  } else {
    ElMessage.error(response.msg || '上传失败')
  }
}

// 上传失败
const handleError = (error) => {
  console.error('上传错误:', error)
  ElMessage.error('上传失败，请重试')
}

// 删除图片
const handleRemove = (uploadFile) => {
  const urls = fileList.value
    .filter(file => file.uid !== uploadFile.uid)
    .map(file => file.url)
    .filter(url => url)
  
  // 根据 max 属性决定返回类型
  if (props.max === 1) {
    // 单图模式：返回字符串
    emit('update:modelValue', urls[0] || '')
  } else {
    // 多图模式：返回数组
    emit('update:modelValue', urls)
  }
}

// 预览图片
const handlePreview = (uploadFile) => {
  previewUrl.value = uploadFile.url
  dialogVisible.value = true
}
</script>

<style scoped>
.image-upload :deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
}

.image-upload :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 100px;
  height: 100px;
}
</style>
