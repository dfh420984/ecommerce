# 图片上传组件使用说明

## 组件位置
`src/components/ImageUpload.vue`

## 功能特性

✅ **多图上传**：支持一次上传多张图片
✅ **图片预览**：点击已上传图片可放大预览
✅ **文件验证**：自动验证文件类型和大小
✅ **拖拽上传**：支持拖拽文件到上传区域
✅ **进度显示**：显示上传进度
✅ **删除功能**：可删除已上传的图片

## 使用方法

### 1. 基础用法

```vue
<template>
  <ImageUpload v-model="imageList" />
</template>

<script setup>
import { ref } from 'vue'
import ImageUpload from '@/components/ImageUpload.vue'

const imageList = ref([])
</script>
```

### 2. 在表单中使用

```vue
<el-form-item label="商品图片">
  <ImageUpload v-model="form.images" />
</el-form-item>
```

## Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 图片URL数组 | Array | [] |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| update:modelValue | 图片列表变化时触发 | (value: string[]) |

## 示例

### 商品管理
```vue
<template>
  <el-form :model="form">
    <el-form-item label="商品图片">
      <ImageUpload v-model="form.images" />
    </el-form-item>
  </el-form>
</template>

<script setup>
const form = ref({
  images: []  // 存储图片URL数组
})
</script>
```

### 轮播图管理
```vue
<template>
  <ImageUpload v-model="bannerImages" />
</template>

<script setup>
const bannerImages = ref([
  '/uploads/banner1.jpg',
  '/uploads/banner2.jpg'
])
</script>
```

## 技术实现

### 1. 上传配置
- **接口地址**: `/api/upload`
- **请求方法**: POST
- **Content-Type**: multipart/form-data
- **文件大小限制**: 5MB
- **支持格式**: jpg, jpeg, png, gif, webp

### 2. 认证方式
自动从 localStorage 获取 token 并添加到请求头：
```javascript
Authorization: Bearer ${token}
```

### 3. 响应格式
成功响应：
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "url": "/uploads/xxx-xxx-xxx.jpg"
  }
}
```

## 注意事项

1. **Token 认证**：确保用户已登录，localStorage 中有 token
2. **文件大小**：单张图片不超过 5MB
3. **图片格式**：只支持常见图片格式
4. **代理配置**：开发环境已通过 vite.config.js 配置代理
5. **后端服务**：确保后端服务正常运行在 8686 端口

## 样式定制

可以通过 CSS 变量或 deep selector 自定义样式：

```css
.image-upload :deep(.el-upload--picture-card) {
  width: 120px;
  height: 120px;
}
```

## 常见问题

### Q: 上传失败怎么办？
A: 检查以下几点：
1. 后端服务是否运行
2. Token 是否有效
3. 文件大小是否超限
4. 文件格式是否正确

### Q: 如何修改上传大小限制？
A: 修改 `ImageUpload.vue` 中的 `beforeUpload` 函数：
```javascript
const isLt5M = file.size / 1024 / 1024 < 10  // 改为10MB
```

### Q: 如何添加更多图片格式支持？
A: 在后端 `handlers/upload.go` 中添加：
```go
allowedExts := map[string]bool{
    "jpg": true, "jpeg": true, "png": true, 
    "gif": true, "webp": true, "bmp": true,  // 添加 bmp
}
```
