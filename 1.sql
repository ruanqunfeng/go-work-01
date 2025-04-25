

CREATE TABLE `eth_block` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `block_number` BIGINT NOT NULL UNIQUE COMMENT '区块号(唯一索引)',
  `block_hash` VARCHAR(66) NOT NULL COMMENT '区块哈希',
  `parent_hash` VARCHAR(66) NOT NULL COMMENT '父区块哈希',
  `miner_address` VARCHAR(42) NOT NULL COMMENT '矿工地址',
  `gas_used` BIGINT NOT NULL COMMENT '消耗Gas总量',
  `gas_limit` BIGINT NOT NULL COMMENT 'Gas限制',
  `timestamp` BIGINT NOT NULL COMMENT '区块时间戳',
  `transaction_count` INT NOT NULL COMMENT '交易数量',
  `block_data` JSON NOT NULL COMMENT '完整区块数据(JSON格式)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_block_number` (`block_number`),
  UNIQUE KEY `idx_block_hash` (`block_hash`),
  KEY `idx_timestamp` (`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='区块主表';


CREATE TABLE `eth_block_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `tag_name` VARCHAR(20) NOT NULL COMMENT '标签(head/finalized/safe)',
  `block_number` BIGINT NOT NULL COMMENT '对应区块号',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='区块标签映射表';


CREATE TABLE `eth_tx` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `tx_hash` VARCHAR(66) NOT NULL COMMENT '交易哈希(主键)',
  `block_number` BIGINT NOT NULL COMMENT '所属区块号',
  `from_address` VARCHAR(42) NOT NULL COMMENT '发送方地址',
  `to_address` VARCHAR(42) DEFAULT NULL COMMENT '接收方地址',
  `value` DECIMAL(38,18) NOT NULL COMMENT '交易金额',
  `gas_price` BIGINT NOT NULL COMMENT 'Gas价格',
  `gas_limit` BIGINT NOT NULL COMMENT 'Gas限制',
  `gas` BIGINT NOT NULL COMMENT 'Gas消耗',
  `input_data` TEXT COMMENT '合约调用数据',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tx_hash` (`tx_hash`),
  KEY `idx_block_number` (`block_number`),
  KEY `idx_from_address` (`from_address`),
  KEY `idx_to_address` (`to_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='交易主表';


CREATE TABLE `eth_tx_receipt` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `tx_hash` VARCHAR(66) NOT NULL COMMENT '交易哈希',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '交易状态(0失败 1成功)',
  `gas_used` BIGINT NOT NULL COMMENT '实际消耗Gas',
  `contract_address` VARCHAR(42) DEFAULT NULL COMMENT '合约地址',
  `logs` JSON NOT NULL COMMENT '事件日志',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tx_hash` (`tx_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='交易收据表';
