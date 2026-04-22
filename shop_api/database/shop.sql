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
  `pay_time` TIMESTAMP NULL COMMENT '支付时间',
  `consignee` VARCHAR(50) NOT NULL COMMENT '收货人',
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `province` VARCHAR(50) NOT NULL COMMENT '省份',
  `city` VARCHAR(50) NOT NULL COMMENT '城市',
  `district` VARCHAR(50) NOT NULL COMMENT '区县',
  `address` VARCHAR(255) NOT NULL COMMENT '详细地址',
  `remark` VARCHAR(500) DEFAULT '' COMMENT '订单备注',
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
('shop_logo', '/uploads/logo.png', '店铺Logo');
