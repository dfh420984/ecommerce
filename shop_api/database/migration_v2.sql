-- ============================================
-- 电商系统数据库升级脚本 V2.0
-- 功能：订单超时取消、退款管理、评价系统、收藏、浏览历史
-- 执行时间：2026-04-24
-- ============================================

-- 1. 订单表增强 - 增加过期时间和预占库存字段
ALTER TABLE orders 
ADD COLUMN expire_time DATETIME COMMENT '订单过期时间' AFTER created_at,
ADD COLUMN stock_reserved INT DEFAULT 0 COMMENT '预占库存数量' AFTER stock;

-- 添加索引优化超时订单查询性能
ALTER TABLE orders 
ADD INDEX idx_expire_status (expire_time, order_status, pay_status);

-- 2. 商品表增强 - 增加评分相关字段
ALTER TABLE products 
ADD COLUMN avg_rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '平均评分' AFTER sales,
ADD COLUMN review_count INT DEFAULT 0 COMMENT '评价数量' AFTER avg_rating;

-- 3. 创建退款申请表
CREATE TABLE IF NOT EXISTS refund_applications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL COMMENT '订单ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    refund_type VARCHAR(50) NOT NULL DEFAULT 'refund_only' COMMENT '退款类型:refund_only/exchange',
    reason VARCHAR(500) NOT NULL COMMENT '退款原因',
    images JSON COMMENT '凭证图片',
    refund_amount DECIMAL(10,2) NOT NULL COMMENT '退款金额',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态:1待审核 2已通过 3已拒绝 4退款中 5已退款',
    handler_id BIGINT COMMENT '处理人ID',
    handler_reply TEXT COMMENT '审核意见',
    refund_trade_no VARCHAR(100) COMMENT '退款交易号',
    refunded_at DATETIME COMMENT '退款完成时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_order (order_id),
    INDEX idx_user (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款申请表';

-- 4. 创建商品评价表
CREATE TABLE IF NOT EXISTS product_reviews (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL COMMENT '订单ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    order_item_id BIGINT NOT NULL COMMENT '订单项ID',
    rating TINYINT NOT NULL COMMENT '评分1-5',
    content TEXT COMMENT '评价内容',
    images JSON COMMENT '评价图片',
    reply_content TEXT COMMENT '商家回复',
    reply_time DATETIME COMMENT '回复时间',
    is_anonymous TINYINT DEFAULT 0 COMMENT '是否匿名',
    status TINYINT DEFAULT 1 COMMENT '状态:1正常 0隐藏',
    helpful_count INT DEFAULT 0 COMMENT '点赞数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_product (product_id),
    INDEX idx_user (user_id),
    UNIQUE KEY uk_order_item (order_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';

-- 5. 创建用户收藏表
CREATE TABLE IF NOT EXISTS user_favorites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_product (user_id, product_id),
    INDEX idx_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户收藏表';

-- 6. 创建浏览历史表
CREATE TABLE IF NOT EXISTS browse_histories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_time (user_id, viewed_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='浏览历史表';

-- ============================================
-- 数据迁移说明
-- ============================================
-- 1. 对于已存在的订单，设置过期时间为创建时间+30分钟
UPDATE orders 
SET expire_time = DATE_ADD(created_at, INTERVAL 30 MINUTE)
WHERE expire_time IS NULL AND order_status = 1 AND pay_status = 0;

-- 2. 初始化商品评分字段（如果已有评价数据）
UPDATE products p
LEFT JOIN (
    SELECT product_id, AVG(rating) as avg_rating, COUNT(*) as review_count
    FROM product_reviews
    WHERE status = 1
    GROUP BY product_id
) r ON p.id = r.product_id
SET p.avg_rating = COALESCE(r.avg_rating, 0),
    p.review_count = COALESCE(r.review_count, 0);

-- ============================================
-- 验证脚本
-- ============================================
-- 检查字段是否添加成功
SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'orders' 
  AND COLUMN_NAME IN ('expire_time', 'stock_reserved');

SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'products' 
  AND COLUMN_NAME IN ('avg_rating', 'review_count');

-- 检查表是否创建成功
SHOW TABLES LIKE 'refund_applications';
SHOW TABLES LIKE 'product_reviews';
SHOW TABLES LIKE 'user_favorites';
SHOW TABLES LIKE 'browse_histories';
