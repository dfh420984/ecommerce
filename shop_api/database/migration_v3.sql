-- ============================================
-- 电商系统数据库升级脚本 V3.0
-- 功能：点赞、收藏、评论功能增强
-- 执行时间：2026-04-26
-- ============================================

-- 1. 创建商品点赞表
CREATE TABLE IF NOT EXISTS `product_likes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
  `product_id` BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_user` (`product_id`, `user_id`),
  INDEX `idx_product_id` (`product_id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品点赞表';

-- 2. 商品表添加点赞数字段
ALTER TABLE `products` 
ADD COLUMN `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数' AFTER `review_count`;

-- 3. 插入初始数据（可选，用于测试）
-- 无

-- ============================================
-- 验证脚本
-- ============================================
-- 检查表是否创建成功
SHOW TABLES LIKE 'product_likes';

-- 检查字段是否添加成功
SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = DATABASE() 
  AND TABLE_NAME = 'products' 
  AND COLUMN_NAME = 'like_count';
