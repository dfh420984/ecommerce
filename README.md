# 电商小程序项目

一套完整的电商解决方案，包含微信小程序用户端、VUE3管理后台和Go后端API。

## 项目结构

```
ecommerce/
├── shop_api/          # Go后端API (Gin框架)
│   ├── config/        # 配置模块
│   ├── database/     # 数据库连接
│   ├── handlers/      # HTTP处理器
│   ├── middleware/    # 中间件
│   ├── models/        # 数据模型
│   ├── routes/        # 路由配置
│   ├── services/      # 业务服务(支付等)
│   ├── utils/         # 工具函数
│   ├── uploads/       # 上传文件目录
│   ├── main.go        # 程序入口
│   ├── go.mod         # Go模块
│   └── config.yaml    # 配置文件
│
├── shop_admin/        # VUE3管理后台
│   ├── src/
│   │   ├── api/       # API请求
│   │   ├── router/    # 路由配置
│   │   ├── store/     # 状态管理
│   │   ├── utils/     # 工具函数
│   │   └── views/     # 页面组件
│   ├── package.json
│   └── vite.config.js
│
├── shop_miniapp/      # 微信小程序
│   ├── pages/         # 页面
│   ├── utils/         # 工具函数
│   ├── app.js         # 小程序入口
│   └── app.json       # 小程序配置
│
└── database/          # 数据库设计
    └── shop.sql       # 数据库脚本
```

## 技术栈

### 后端 (shop_api)
- **框架**: Gin v1.9
- **数据库**: MySQL 8.0 + GORM
- **缓存**: Redis
- **认证**: JWT
- **支付**: 微信支付/支付宝(可扩展)

### 管理后台 (shop_admin)
- **框架**: Vue3 + Composition API
- **UI**: Element Plus
- **构建**: Vite
- **状态管理**: Pinia
- **请求**: Axios

### 小程序 (shop_miniapp)
- **框架**: 原生微信小程序
- **API请求**: Promise封装

## 功能模块

### 用户端(小程序)
- [x] 首页轮播图
- [x] 商品分类浏览
- [x] 商品列表/详情
- [x] 购物车管理
- [x] 下单流程
- [x] 微信支付
- [x] 地址管理
- [x] 用户中心
- [x] 订单管理

### 管理后台
- [x] 登录/退出
- [x] 分类管理(树形结构)
- [x] 商品管理(CRUD)
- [x] 轮播图管理
- [x] 订单管理
- [x] 用户管理

### 后端API
- [x] 用户注册/登录(JWT)
- [x] 商品CRUD
- [x] 分类管理
- [x] 购物车管理
- [x] 订单流程
- [x] 微信支付统一下单/回调/退款
- [x] 支付宝支付接口(可扩展)
- [x] 库存管理
- [x] 地址管理
- [x] 操作日志

## 快速开始

### 1. 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+
- 微信开发者工具

### 2. 数据库初始化

```bash
mysql -u root -p < shop_api/database/shop.sql
```

或导入 `shop_api/database/shop.sql` 文件到MySQL。

### 3. 后端启动

```bash
cd shop_api

# 修改配置文件 config.yaml
# 设置数据库、Redis、微信支付等参数

# 下载依赖
go mod tidy

# 启动服务
go run main.go
```

服务默认运行在 `http://localhost:8686`

### 4. 管理后台启动

```bash
cd shop_admin

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build
```

后台访问: `http://localhost:8000`
默认账号: admin / 123456

### 5. 小程序配置

1. 打开微信开发者工具
2. 导入 `shop_miniapp` 目录
3. 修改 `utils/request.js` 中的 `apiBase` 为你的后端地址
4. 在微信公众平台配置合法域名

## 配置文件说明

### shop_api/config.yaml

```yaml
app:
  name: shop_api
  host: 0.0.0.0
  port: 8686
  mode: debug          # debug/release
  jwt_secret: your-jwt-secret-key-change-in-production
  jwt_expire: 72       # Token过期时间(小时)

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your-password
  dbname: shop_db
  charset: utf8mb4

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

wechat:
  appid: your-wechat-appid
  secret: your-wechat-secret
  mchid: your-wechat-mchid
  apikey: your-wechat-apikey
  notify_url: http://your-domain.com/api/notify/wechat

alipay:
  appid: your-alipay-appid
  private_key: your-private-key
  public_key: alipay-public-key
  notify_url: http://your-domain.com/api/notify/alipay
```

## API接口文档

### 小程序端接口 (/api/miniapp)
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /miniapp/register | 用户注册 |
| POST | /miniapp/login | 用户登录 |
| POST | /miniapp/wechat_login | 微信登录 |
| GET | /miniapp/banners | 获取轮播图 |
| GET | /miniapp/categories | 获取分类 |
| GET | /miniapp/products | 获取商品列表 |
| GET | /miniapp/products/recommend | 获取推荐商品 |
| GET | /miniapp/products/new | 获取新品 |
| GET | /miniapp/products/:id | 获取商品详情 |

### 用户接口 (/api/user) - 需要Token
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /user/info | 获取用户信息 |
| PUT | /user/info | 更新用户信息 |
| PUT | /user/password | 修改密码 |
| GET | /user/addresses | 获取地址列表 |
| POST | /user/addresses | 添加地址 |
| PUT | /user/addresses/:id | 更新地址 |
| DELETE | /user/addresses/:id | 删除地址 |
| GET | /user/cart | 获取购物车 |
| POST | /user/cart | 加入购物车 |
| PUT | /user/cart/:id | 更新购物车 |
| DELETE | /user/cart/:id | 删除购物车 |
| GET | /user/orders | 获取订单列表 |
| POST | /user/orders | 创建订单 |
| PUT | /user/orders/:id/cancel | 取消订单 |
| PUT | /user/orders/:id/confirm | 确认收货 |
| POST | /user/pay | 获取支付链接 |
| POST | /user/pay/refund/:id | 申请退款 |

### 管理端接口 (/api/admin) - 需要Admin Token
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/categories | 获取分类 |
| POST | /admin/categories | 创建分类 |
| PUT | /admin/categories/:id | 更新分类 |
| DELETE | /admin/categories/:id | 删除分类 |
| GET | /admin/products | 获取商品 |
| POST | /admin/products | 创建商品 |
| PUT | /admin/products/:id | 更新商品 |
| DELETE | /admin/products/:id | 删除商品 |
| GET | /admin/banners | 获取轮播图 |
| POST | /admin/banners | 创建轮播图 |
| PUT | /admin/banners/:id | 更新轮播图 |
| DELETE | /admin/banners/:id | 删除轮播图 |

## 数据库表结构

- **users** - 用户表
- **categories** - 商品分类表
- **products** - 商品表
- **banners** - 轮播图表
- **addresses** - 收货地址表
- **cart** - 购物车表
- **orders** - 订单表
- **order_items** - 订单明细表
- **pay_logs** - 支付日志表
- **configs** - 系统配置表
- **operation_logs** - 操作日志表

## 部署说明

### 后端部署

1. 安装Go环境
2. 配置MySQL和Redis
3. 修改config.yaml配置
4. 编译运行:
```bash
go build -o shop_api main.go
./shop_api
```

### 管理后台部署

1. 编译:
```bash
npm run build
```

2. 将dist目录部署到Nginx:
```nginx
server {
    listen 80;
    server_name admin.your-domain.com;
    root /path/to/dist;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }
    location /api {
        proxy_pass http://localhost:8686;
    }
}
```

### 小程序发布

1. 在微信开发者工具中上传
2. 在微信公众平台提交审核
3. 配置合法域名

## 安全注意事项

1. **密钥保管**: 所有密钥(appid, apikey, jwt_secret等)必须保密
2. **生产环境**: 生产环境务必修改默认密码
3. **HTTPS**: 生产环境必须使用HTTPS
4. **域名配置**: 小程序需配置request合法域名
5. **参数校验**: 所有用户输入必须进行校验

## 开发指南

### 添加新的支付方式

在 `services/pay.go` 中添加新的支付服务类，实现统一接口。

### 添加新的API

1. 在 `handlers/` 添加处理器
2. 在 `routes/router.go` 注册路由
3. 添加对应的中间件(如需要认证)

## License

MIT License
