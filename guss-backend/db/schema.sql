-- 데이터베이스 생성 및 선택
CREATE DATABASE IF NOT EXISTS guss DEFAULT CHARACTER SET utf8mb4;
USE guss;

-- 1. 사용자 테이블: 비밀번호 해시(Bcrypt) 저장을 위해 TEXT 타입 유지
CREATE TABLE user_table (
    user_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    user_phone VARCHAR(20) NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    user_pw VARCHAR(255) NOT NULL -- 해시된 비밀번호 저장
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 체육관 테이블: 기구 정보는 별도 테이블로 분리됨
CREATE TABLE guss_table (
    guss_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    guss_name VARCHAR(100) NOT NULL,
    guss_address TEXT,
    guss_phone VARCHAR(20),
    guss_status VARCHAR(10) DEFAULT 'open', -- 'open' / 'close'
    guss_user_count INT DEFAULT 0,          -- 실시간 이용 인원
    guss_size INT NOT NULL,                 -- 수용 가능 최대 인원
    guss_open_time VARCHAR(10),
    guss_close_time VARCHAR(10)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 기구 테이블: 1:N 관계 (체육관별 기구 관리)
CREATE TABLE equipment_table (
    equip_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_guss_number BIGINT NOT NULL,
    equip_name VARCHAR(100) NOT NULL,
    equip_category VARCHAR(50),
    equip_quantity INT DEFAULT 0,
    equip_status VARCHAR(20) DEFAULT 'active', -- 'active' / 'maintenance'
    purchase_date DATE,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 예약 테이블: 노쇼 방지 및 실시간 상태 관리
CREATE TABLE revs_table (
    revs_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_user_number BIGINT NOT NULL,
    fk_guss_number BIGINT NOT NULL,
    revs_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    revs_status VARCHAR(20) DEFAULT 'CONFIRMED', -- 'CONFIRMED' / 'CANCELED'
    FOREIGN KEY (fk_user_number) REFERENCES user_table(user_number) ON DELETE CASCADE,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. 매출 테이블: 관리자 통계용
CREATE TABLE sales_table (
    sales_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_guss_number BIGINT NOT NULL,
    sales_type VARCHAR(50), -- 'DAILY', 'MONTHLY', 'PT' 등
    sales_amount INT DEFAULT 0,
    sales_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- [테스트 데이터 주입]
INSERT INTO guss_table (guss_name, guss_address, guss_size) VALUES ('명지 체육시설', '서울 서대문구', 50);