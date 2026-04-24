-- 电商数据库设计
-- 数据库创建
CREATE DATABASE IF NOT EXISTS shop_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE shop_db;

-- 用户表
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
  `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码(加密)',
  `nickname` VARCHAR(100) DEFAULT '' COMMENT '昵称',
  `avatar` VARCHAR(500) DEFAULT '' COMMENT '头像URL',
  `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
  `email` VARCHAR(100) DEFAULT '' COMMENT '邮箱',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0禁用,1启用',
  `user_type` TINYINT NOT NULL DEFAULT 1 COMMENT '用户类型:1普通用户,2管理员',
  `openid` VARCHAR(100) DEFAULT '' COMMENT '微信openid',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_username` (`username`),
  INDEX `idx_openid` (`openid`),
  INDEX `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 分类表
CREATE TABLE `categories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
  `name` VARCHAR(100) NOT NULL COMMENT '分类名称',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类ID',
  `level` TINYINT NOT NULL DEFAULT 1 COMMENT '层级',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `icon` VARCHAR(500) DEFAULT '' COMMENT '图标',
  `image` VARCHAR(500) DEFAULT '' COMMENT '图片',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0禁用,1启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_parent_id` (`parent_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- 商品表
CREATE TABLE `products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '商品ID',
  `name` VARCHAR(200) NOT NULL COMMENT '商品名称',
  `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
  `price` DECIMAL(10,2) NOT NULL COMMENT '价格',
  `original_price` DECIMAL(10,2) DEFAULT 0 COMMENT '原价',
  `cost` DECIMAL(10,2) DEFAULT 0 COMMENT '成本价',
  `stock` INT NOT NULL DEFAULT 0 COMMENT '库存',
  `sales` INT NOT NULL DEFAULT 0 COMMENT '销量',
  `images` TEXT COMMENT '商品图片JSON数组',
  `description` TEXT COMMENT '商品描述',
  `content` TEXT COMMENT '商品详情(HTML)',
  `specs` TEXT COMMENT '规格JSON',
  `is_online` TINYINT NOT NULL DEFAULT 1 COMMENT '是否上架:0下架,1上架',
  `is_recommend` TINYINT NOT NULL DEFAULT 0 COMMENT '是否推荐',
  `is_new` TINYINT NOT NULL DEFAULT 0 COMMENT '是否新品',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0删除,1正常',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_is_online` (`is_online`),
  INDEX `idx_is_recommend` (`is_recommend`),
  INDEX `idx_status` (`status`),
  FULLTEXT `ft_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 轮播图表
CREATE TABLE `banners` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '轮播图ID',
  `title` VARCHAR(200) DEFAULT '' COMMENT '标题',
  `image` VARCHAR(500) NOT NULL COMMENT '图片URL',
  `link` VARCHAR(500) DEFAULT '' COMMENT '链接地址',
  `link_type` TINYINT NOT NULL DEFAULT 1 COMMENT '链接类型:1商品,2分类,3链接',
  `target_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '目标ID(商品/分类ID)',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0禁用,1启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_status` (`status`),
  INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 地址表
CREATE TABLE `addresses` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '地址ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `consignee` VARCHAR(50) NOT NULL COMMENT '收货人',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `province` VARCHAR(50) NOT NULL COMMENT '省份',
  `city` VARCHAR(50) NOT NULL COMMENT '城市',
  `district` VARCHAR(50) NOT NULL COMMENT '区县',
  `address` VARCHAR(255) NOT NULL COMMENT '详细地址',
  `postal_code` VARCHAR(20) DEFAULT '' COMMENT '邮政编码',
  `is_default` TINYINT NOT NULL DEFAULT 0 COMMENT '是否默认:0否,1是',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';

-- 购物车表
CREATE TABLE `cart` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '购物车ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `sku_id` BIGINT UNSIGNED DEFAULT 0 COMMENT 'SKU ID',
  `quantity` INT NOT NULL DEFAULT 1 COMMENT '数量',
  `selected` TINYINT NOT NULL DEFAULT 1 COMMENT '是否选中:0否,1是',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY `uk_user_product_sku` (`user_id`, `product_id`, `sku_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 订单表
CREATE TABLE `orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '订单ID',
  `order_no` VARCHAR(32) NOT NULL UNIQUE COMMENT '订单号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `order_status` TINYINT NOT NULL DEFAULT 1 COMMENT '订单状态:1待付款,2待发货,3配送中,4已收货,5已完成,6已取消,7退款中,8已退款',
  `pay_status` TINYINT NOT NULL DEFAULT 0 COMMENT '支付状态:0未支付,1已支付,2已退款',
  `pay_type` TINYINT NOT NULL DEFAULT 0 COMMENT '支付方式:0未支付,1微信,2支付宝,3银行卡',
  `total_amount` DECIMAL(10,2) NOT NULL COMMENT '订单总金额',
  `discount_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '优惠金额',
  `freight_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '运费金额',
  `pay_amount` DECIMAL(10,2) NOT NULL COMMENT '实付金额',
  `coupon_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用的优惠券ID',
  `coupon_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '优惠券优惠金额',
  `pay_time` TIMESTAMP NULL COMMENT '支付时间',
  `consignee` VARCHAR(50) NOT NULL COMMENT '收货人',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `province` VARCHAR(50) NOT NULL COMMENT '省份',
  `city` VARCHAR(50) NOT NULL COMMENT '城市',
  `district` VARCHAR(50) NOT NULL COMMENT '区县',
  `address` VARCHAR(255) NOT NULL COMMENT '详细地址',
  `remark` VARCHAR(500) DEFAULT '' COMMENT '订单备注',
  `express_company` VARCHAR(100) DEFAULT '' COMMENT '快递公司',
  `express_no` VARCHAR(100) DEFAULT '' COMMENT '快递单号',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `cancel_time` TIMESTAMP NULL COMMENT '取消时间',
  `complete_time` TIMESTAMP NULL COMMENT '完成时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_order_no` (`order_no`),
  INDEX `idx_order_status` (`order_status`),
  INDEX `idx_pay_status` (`pay_status`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 订单明细表
CREATE TABLE `order_items` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '订单明细ID',
  `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `product_name` VARCHAR(200) NOT NULL COMMENT '商品名称',
  `product_image` VARCHAR(500) NOT NULL COMMENT '商品图片',
  `price` DECIMAL(10,2) NOT NULL COMMENT '单价',
  `quantity` INT NOT NULL COMMENT '数量',
  `sku_id` BIGINT UNSIGNED DEFAULT 0 COMMENT 'SKU ID',
  `sku_name` VARCHAR(200) DEFAULT '' COMMENT 'SKU规格',
  `subtotal` DECIMAL(10,2) NOT NULL COMMENT '小计金额',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';

-- 支付日志表
CREATE TABLE `pay_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '支付日志ID',
  `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
  `order_no` VARCHAR(32) NOT NULL COMMENT '订单号',
  `trade_no` VARCHAR(64) DEFAULT '' COMMENT '第三方交易号',
  `pay_type` TINYINT NOT NULL COMMENT '支付方式:1微信,2支付宝,3银行卡',
  `pay_status` TINYINT NOT NULL COMMENT '支付状态:0失败,1成功,2处理中',
  `pay_amount` DECIMAL(10,2) NOT NULL COMMENT '支付金额',
  `notify_data` TEXT COMMENT '回调数据JSON',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_order_no` (`order_no`),
  INDEX `idx_trade_no` (`trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付日志表';

-- 系统配置表
CREATE TABLE `configs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '配置ID',
  `name` VARCHAR(100) NOT NULL UNIQUE COMMENT '配置名称',
  `value` TEXT COMMENT '配置值',
  `description` VARCHAR(500) DEFAULT '' COMMENT '描述',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 操作日志表
CREATE TABLE `operation_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '日志ID',
  `user_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '操作用户ID',
  `username` VARCHAR(50) DEFAULT '' COMMENT '用户名',
  `action` VARCHAR(100) NOT NULL COMMENT '操作动作',
  `target` VARCHAR(100) DEFAULT '' COMMENT '操作对象',
  `target_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '对象ID',
  `content` TEXT COMMENT '操作内容',
  `ip` VARCHAR(50) DEFAULT '' COMMENT 'IP地址',
  `user_agent` VARCHAR(500) DEFAULT '' COMMENT 'UserAgent',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_action` (`action`),
  INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 帮助中心分类表
CREATE TABLE `help_categories` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '分类ID',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0禁用,1启用',
  `description` VARCHAR(200) DEFAULT '' COMMENT '描述',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_status` (`status`),
  INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帮助中心分类表';

-- 帮助中心问题表
CREATE TABLE `help_questions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '问题ID',
  `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
  `title` VARCHAR(200) NOT NULL COMMENT '问题标题',
  `answer` TEXT NOT NULL COMMENT '问题答案',
  `sort` INT NOT NULL DEFAULT 0 COMMENT '排序',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0禁用,1启用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='帮助中心问题表';

-- 插入默认数据
INSERT INTO `users` (`username`, `password`, `nickname`, `user_type`, `status`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '管理员', 2, 1);

INSERT INTO `categories` (`name`, `parent_id`, `level`, `sort`, `status`) VALUES
('水果', 0, 1, 1, 1),
('蔬菜', 0, 1, 2, 1),
('肉类', 0, 1, 3, 1),
('海鲜', 0, 1, 4, 1),
('蛋奶', 0, 1, 5, 1);

INSERT INTO `categories` (`name`, `parent_id`, `level`, `sort`, `status`) VALUES
('苹果', 1, 2, 1, 1),
('香蕉', 1, 2, 2, 1),
('橙子', 1, 2, 3, 1);

INSERT INTO `banners` (`title`, `image`, `link`, `link_type`, `sort`, `status`) VALUES
('热门商品', '/uploads/banner1.jpg', '', 0, 1, 1),
('新品上市', '/uploads/banner2.jpg', '', 0, 2, 1),
('限时优惠', '/uploads/banner3.jpg', '', 0, 3, 1);

INSERT INTO `products` (`name`, `category_id`, `price`, `original_price`, `stock`, `sales`, `images`, `is_online`, `is_recommend`, `status`) VALUES
('红富士苹果 5斤装', 6, 29.90, 39.90, 100, 50, '["/uploads/product1.jpg"]', 1, 1, 1),
('香蕉 3斤装', 7, 15.90, 19.90, 80, 30, '["/uploads/product2.jpg"]', 1, 1, 1),
('新鲜橙子 5斤装', 8, 24.90, 34.90, 60, 20, '["/uploads/product3.jpg"]', 1, 0, 1);

INSERT INTO `configs` (`name`, `value`, `description`) VALUES
('shop_name', '电商小程序', '店铺名称'),
('shop_logo', '/uploads/logo.png', '店铺Logo'),
('shop_phone', '400-123-4567', '客服电话'),
('shop_address', '北京市朝阳区xxx路xxx号', '店铺地址'),
('min_order_amount', '99', '最小订单金额'),
('free_shipping_amount', '199', '包邮金额'),
('return_policy', '7天无理由退换货', '退换货政策'),
('business_hours', '09:00-21:00', '营业时间'),
('about_intro', '这是一款基于 Go + Vue3 + 微信小程序的电商系统，提供完整的商品展示、购物车、订单管理等功能。', '关于我们简介'),
('about_tech_stack', '后端：Go + Gin + GORM + MySQL\n前端：Vue3 + Element Plus\n小程序：原生微信小程序', '技术栈说明'),
('about_contact', '邮箱：contact@example.com\n电话：400-123-4567\n地址：北京市朝阳区xxx路xxx号', '联系方式'),
('app_version', '1.0.0', '应用版本号'),
('service_time', '09:00-21:00', '客服服务时间'),
('customer_service_phone', '400-123-4567', '客服电话');

-- 插入帮助中心默认数据
INSERT INTO `help_categories` (`name`, `sort`, `status`, `description`) VALUES
('购物指南', 1, 1, '关于购物的常见问题'),
('配送说明', 2, 1, '配送相关问题'),
('售后服务', 3, 1, '退换货和售后问题'),
('账户安全', 4, 1, '账户和安全问题');

INSERT INTO `help_questions` (`category_id`, `title`, `answer`, `sort`, `status`) VALUES
(1, '如何下单购买？', '1. 浏览商品，点击商品进入详情页\n2. 选择规格和数量，点击“加入购物车”\n3. 进入购物车，勾选商品\n4. 点击“结算”，选择收货地址\n5. 提交订单并完成支付', 1, 1),
(1, '如何修改收货地址？', '在订单提交前，可以在确认订单页面点击“收货地址”进行修改。\n\n如果订单已提交但未发货，请联系客服修改。\n\n如果订单已发货，则无法修改地址。', 2, 1),
(1, '如何使用优惠券？', '1. 在“我的-优惠券”中查看可用优惠券\n2. 下单时，系统会自动匹配可用优惠券\n3. 也可以在结算页面手动选择优惠券\n4. 注意查看优惠券的使用条件和有效期', 3, 1),
(1, '支持哪些支付方式？', '目前支持以下支付方式：\n• 微信支付\n• 支付宝支付\n• 银行卡支付\n\n部分商品可能不支持某些支付方式，请以实际为准。', 4, 1),
(2, '配送范围有哪些？', '我们目前支持全国大部分地区配送。\n\n偏远地区（如西藏、新疆等）可能需要额外配送时间，具体以结算页面显示为准。', 1, 1),
(2, '多久能收到货？', '• 同城配送：1-2天\n• 省内配送：2-3天\n• 跨省配送：3-5天\n• 偏远地区：5-7天\n\n具体时间以物流公司实际配送为准。', 2, 1),
(2, '运费是多少？', '• 订单满199元包邮\n• 不满包邮金额，收取固定运费10元\n• 特殊商品（如生鲜、大件）可能有额外运费\n\n具体运费以结算页面显示为准。', 3, 1),
(2, '如何查询物流信息？', '1. 进入“我的-我的订单”\n2. 找到对应订单，点击“查看物流”\n3. 即可看到详细的物流跟踪信息\n\n也可以通过物流公司官网查询。', 4, 1),
(3, '退换货政策是什么？', '我们提供7天无理由退换货服务。\n\n退换货条件：\n• 商品未使用、未拆封\n• 包装完整、配件齐全\n• 在退换货有效期内\n\n以下情况不支持退换货：\n• 定制类商品\n• 鲜活易腐类商品\n• 数字化商品', 1, 1),
(3, '如何申请退款？', '1. 进入“我的-我的订单”\n2. 找到需要退款的订单\n3. 点击“申请退款”\n4. 选择退款原因，提交申请\n5. 等待商家审核\n6. 审核通过后，退款将原路返回\n\n退款到账时间：1-7个工作日', 2, 1),
(3, '退款多久到账？', '退款审核通过后：\n• 微信支付：1-3个工作日\n• 支付宝：1-3个工作日\n• 银行卡：3-7个工作日\n\n具体到账时间以银行处理为准。', 3, 1),
(3, '商品有质量问题怎么办？', '如果收到商品有质量问题：\n1. 请在签收后24小时内联系客服\n2. 提供商品照片和视频\n3. 客服会为您安排退换货\n4. 来回运费由我们承担\n\n请保留好商品包装和配件。', 4, 1),
(4, '如何修改密码？', '1. 进入“我的-设置”\n2. 点击“修改密码”\n3. 输入原密码和新密码\n4. 确认修改\n\n如果忘记密码，可以点击“忘记密码”通过手机号重置。', 1, 1),
(4, '如何绑定手机号？', '1. 进入“我的-设置”\n2. 点击“绑定手机”\n3. 输入手机号和验证码\n4. 完成绑定\n\n一个手机号只能绑定一个账号。', 2, 1),
-- 修改订单表，添加优惠券相关字段
ALTER TABLE `orders` ADD COLUMN `coupon_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用的优惠券ID';
ALTER TABLE `orders` ADD COLUMN `coupon_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '优惠券抵扣金额';

-- 优惠券模板表
CREATE TABLE `coupons` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '优惠券ID',
  `name` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
  `type` TINYINT NOT NULL DEFAULT 1 COMMENT '类型:1-满减券,2-折扣券,3-无门槛券',
  `discount_value` DECIMAL(10,2) NOT NULL COMMENT '优惠值(满减金额/折扣率)',
  `min_amount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '最低消费金额',
  `max_discount` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '最大优惠金额(折扣券用)',
  `total_count` INT NOT NULL DEFAULT 0 COMMENT '发放总量(0表示不限)',
  `received_count` INT NOT NULL DEFAULT 0 COMMENT '已领取数量',
  `used_count` INT NOT NULL DEFAULT 0 COMMENT '已使用数量',
  `per_user_limit` INT NOT NULL DEFAULT 1 COMMENT '每人限领数量',
  `valid_type` TINYINT NOT NULL DEFAULT 1 COMMENT '有效期类型:1-固定时间,2-领取后N天',
  `start_time` TIMESTAMP NULL COMMENT '开始时间',
  `end_time` TIMESTAMP NULL COMMENT '结束时间',
  `valid_days` INT NOT NULL DEFAULT 0 COMMENT '有效天数(领取后)',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:0-禁用,1-启用',
  `is_new_user` TINYINT NOT NULL DEFAULT 0 COMMENT '是否新人券:0-否,1-是',
  `description` TEXT COMMENT '使用说明',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_status` (`status`),
  INDEX `idx_time` (`start_time`, `end_time`),
  INDEX `idx_is_new_user` (`is_new_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券模板表';

-- 用户优惠券表
CREATE TABLE `user_coupons` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '用户优惠券ID',
  `coupon_id` BIGINT UNSIGNED NOT NULL COMMENT '优惠券模板ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态:1-未使用,2-已使用,3-已过期',
  `receive_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
  `use_time` TIMESTAMP NULL COMMENT '使用时间',
  `order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用的订单ID',
  `expire_time` TIMESTAMP NOT NULL COMMENT '过期时间',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX `idx_user_status` (`user_id`, `status`),
  INDEX `idx_expire` (`expire_time`),
  INDEX `idx_coupon` (`coupon_id`),
  FOREIGN KEY (`coupon_id`) REFERENCES `coupons`(`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- 插入默认数据（示例优惠券）
INSERT INTO `coupons` (`name`, `type`, `discount_value`, `min_amount`, `max_discount`, `total_count`, `per_user_limit`, `valid_type`, `start_time`, `end_time`, `valid_days`, `status`, `is_new_user`, `description`) VALUES
('新人专享券', 1, 10.00, 50.00, 0, 0, 1, 2, NULL, NULL, 7, 1, 1, '新用户注册即送，满50元减10元'),
('满减优惠券', 1, 20.00, 100.00, 0, 1000, 2, 1, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 0, 1, 0, '全场通用，满100元减20元'),
('8折折扣券', 2, 0.80, 50.00, 50.00, 500, 1, 1, NOW(), DATE_ADD(NOW(), INTERVAL 15 DAY), 0, 1, 0, '商品8折优惠，最高优惠50元');

