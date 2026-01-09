-- 1. 사용자 테이블: 회원 정보 저장
CREATE TABLE user_table (
    user_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL,
    user_phone VARCHAR(20) NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    user_pw TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 체육관 테이블: 시설 정보 및 실시간 인원 관리
CREATE TABLE guss_table (
    guss_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    guss_name VARCHAR(100) NOT NULL,
    guss_address TEXT,
    guss_phone VARCHAR(20),
    guss_status VARCHAR(10) DEFAULT 'open',
    guss_user_count INT DEFAULT 0,
    guss_size INT NOT NULL,
    guss_ma_type VARCHAR(100),   -- 관리 기구 타입
    guss_ma_count INT DEFAULT 0,  -- 관리 기구 수
    guss_ma_state VARCHAR(50)     -- 기구 상태
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 예약 테이블: 사용자와 체육관의 연결 (N:M 관계 해소)
CREATE TABLE revs_table (
    revs_number BIGINT AUTO_INCREMENT PRIMARY KEY,
    fk_user_number BIGINT NOT NULL,
    fk_guss_number BIGINT NOT NULL,
    revs_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    revs_status VARCHAR(20) DEFAULT 'PENDING',
    FOREIGN KEY (fk_user_number) REFERENCES user_table(user_number),
    FOREIGN KEY (fk_guss_number) REFERENCES guss_table(guss_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;