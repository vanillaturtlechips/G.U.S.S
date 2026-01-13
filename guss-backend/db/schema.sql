-- 데이터베이스 생성 및 선택
CREATE DATABASE IF NOT EXISTS guss DEFAULT CHARACTER SET utf8mb4;
USE guss;

-- 1. 사용자 테이블: Bcrypt 해시값 저장을 위해 충분한 길이(255) 확보
CREATE TABLE user_table (
    user_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    user_phone VARCHAR(20) NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    user_pw VARCHAR(255) NOT NULL -- 해시된 비밀번호 저장
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 체육관 테이블: 운영 시간 컬럼(open/close_time) 포함
CREATE TABLE guss_table (
    guss_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    guss_name VARCHAR(100) NOT NULL,
    guss_address TEXT,
    guss_phone VARCHAR(20),
    guss_status VARCHAR(10) DEFAULT 'open',
    guss_user_count INT DEFAULT 0, -- 실시간 이용 인원 카운터
    guss_size INT NOT NULL,        -- 최대 수용 인원
    guss_open_time VARCHAR(10),    -- [오류 수정] 추가됨
    guss_close_time VARCHAR(10)    -- [오류 수정] 추가됨
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 기구 테이블: 1:N 관계 정규화
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
    revs_status VARCHAR(20) DEFAULT 'CONFIRMED',
    FOREIGN KEY (fk_user_number) REFERENCES user_table(user_number) ON DELETE CASCADE,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. 매출 테이블: 관리자 통계용
CREATE TABLE sales_table (
    sales_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_guss_number BIGINT NOT NULL,
    sales_type VARCHAR(50), -- 'DAILY', 'MONTHLY' 등
    sales_amount INT DEFAULT 0,
    sales_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



-- 테스트 체육관 5개 주입
INSERT INTO guss_table (guss_name, guss_address, guss_phone, guss_size, guss_open_time, guss_close_time) 
VALUES 
('명지대 MCC 체육시설', '서울 서대문구 거북골로 34', '02-300-1521', 50, '06:00', '23:00'),
('강남 비타민 피트니스', '서울 강남구 테헤란로 123', '02-555-0001', 100, '00:00', '24:00'),
('홍대 시너지 짐', '서울 마포구 와우산로 99', '02-333-7777', 30, '08:00', '22:00'),
('종로 바디빌딩 센터', '서울 종로구 인사동길 10', '02-777-1234', 40, '07:00', '23:00'),
('잠실 스포츠 콤플렉스', '서울 송파구 올림픽로 25', '02-444-5555', 80, '06:00', '24:00');

-- 1번 체육관(명지대) 샘플 기구 주입
INSERT INTO equipment_table (fk_guss_number, equip_name, equip_category, equip_quantity, equip_status, purchase_date)
VALUES 
(1, '천국의 계단', '유산소', 2, 'active', '2025-01-10'),
(1, '레그 프레스', '하체', 1, 'active', '2024-12-20'),
(1, '덤벨 세트', '프리웨이트', 10, 'active', '2025-01-05');