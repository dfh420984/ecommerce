<template>
  <div class="rich-text-editor">
    <Toolbar
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      mode="default"
      class="toolbar"
    />
    <Editor
      :defaultConfig="editorConfig"
      v-model="contentValue"
      mode="default"
      class="editor"
      @onCreated="handleCreated"
      @onChange="handleChange"
    />
  </div>
</template>

<script setup>
import { ref, shallowRef, onBeforeUnmount } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import '@wangeditor/editor/dist/css/style.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  height: {
    type: String,
    default: '500px'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
const contentValue = ref(props.modelValue)

// 工具栏配置
const toolbarConfig = {
  excludeKeys: [
    'fullScreen',  // 排除全屏按钮
  ]
}

// 编辑器配置
const editorConfig = {
  placeholder: '请输入商品详情...',
  MENU_CONF: {
    uploadImage: {
      // 图片上传配置
      server: '/api/upload',
      fieldName: 'file',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      // 自定义插入图片
      customInsert(res, insertFn) {
        if (res.code === 0) {
          insertFn(res.data.url, res.data.url, res.data.url)
        }
      }
    },
    uploadVideo: {
      // 视频上传配置
      server: '/api/upload/video',
      fieldName: 'file',
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      // 自定义插入视频
      customInsert(res, insertFn) {
        if (res.code === 0) {
          insertFn(res.data.url, res.data.url)
        }
      }
    }
  }
}

// 编辑器创建完成
const handleCreated = (editor) => {
  editorRef.value = editor
}

// 内容变化
const handleChange = (editor) => {
  const html = editor.getHtml()
  emit('update:modelValue', html)
}

// 组件销毁时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

// 监听外部值变化
const updateContent = (val) => {
  contentValue.value = val
}

defineExpose({
  updateContent
})
</script>

<style scoped>
.rich-text-editor {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.toolbar {
  border-bottom: 1px solid #dcdfe6;
}

.editor {
  height: v-bind(height);
  overflow-y: hidden;
}

:deep(.w-e-text-container) {
  background-color: #fff;
}

:deep(.w-e-text-placeholder) {
  color: #c0c4cc;
}
</style>
