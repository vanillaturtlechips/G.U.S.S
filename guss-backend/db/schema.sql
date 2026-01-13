-- 1. 사용자 테이블 (변동 없음)
CREATE TABLE user_table (
    user_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    user_phone VARCHAR(20) NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    user_pw TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 체육관 테이블 (기구 관련 컬럼 제거)
CREATE TABLE guss_table (
    guss_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    guss_name VARCHAR(100) NOT NULL,
    guss_address TEXT,
    guss_phone VARCHAR(20),
    guss_status VARCHAR(10) DEFAULT 'open',
    guss_user_count INT DEFAULT 0,
    guss_size INT NOT NULL,
    guss_open_time VARCHAR(10),
    guss_close_time VARCHAR(10)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 기구 테이블 (정규화로 추가됨: 1:N 관계)
CREATE TABLE equipment_table (
    equip_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_guss_number BIGINT NOT NULL,
    equip_name VARCHAR(100) NOT NULL,      -- UI의 name
    equip_category VARCHAR(50),            -- UI의 category
    equip_quantity INT DEFAULT 0,          -- UI의 quantity
    equip_status VARCHAR(20) DEFAULT 'active', -- UI의 status (active/maintenance)
    purchase_date DATE,                    -- UI의 purchaseDate
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 예약 테이블 (변동 없음)
CREATE TABLE revs_table (
    revs_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_user_number BIGINT NOT NULL,
    fk_guss_number BIGINT NOT NULL,
    revs_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    revs_status VARCHAR(20) DEFAULT 'PENDING',
    FOREIGN KEY (fk_user_number) REFERENCES user_table(user_number),
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. 매출 테이블 (설계안 기반 추가)
CREATE TABLE sales_table (
    sales_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_guss_number BIGINT NOT NULL,
    sales_type VARCHAR(50),
    sales_count INT DEFAULT 0,
    sales_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;