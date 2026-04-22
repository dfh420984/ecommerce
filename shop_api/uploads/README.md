# 上传文件目录

此目录用于存储用户上传的文件，包括：

## 文件类型

- **轮播图**: banner1.jpg, banner2.jpg, banner3.jpg (800x400)
- **商品图**: product1.jpg, product2.jpg, product3.jpg (400x400)
- **Logo**: logo.png (200x200)

## 访问方式

### 开发环境
通过 Vite 代理访问：
```
http://localhost:8000/uploads/banner1.jpg
```

### 生产环境
直接访问后端服务：
```
http://your-domain.com/uploads/banner1.jpg
```

## 生成占位图片

运行以下命令生成测试用的占位图片：

```bash
cd shop_api
go run generate_images.go
```

## 注意事项

1. 确保此目录有写入权限
2. 生产环境建议配置 CDN
3. 定期清理无用文件
4. 建议设置文件大小限制
