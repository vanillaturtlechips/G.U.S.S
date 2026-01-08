# GUSS Frontend

헬스장 관리 시스템 프론트엔드 프로젝트입니다.

## 기술 스택

- **npm**: 10.9.3
- **React**: 19.2.3
- **TypeScript**: ~5.6.2
- **Tailwind CSS**: 4.1.18
- **Vite**: ^6.0.3

## 프로젝트 구조

```
frontend/
├── src/
│   ├── pages/          # 페이지 컴포넌트
│   │   ├── admin.tsx   # 관리자 페이지
│   │   ├── guss.tsx    # 헬스장 현황 페이지
│   │   └── register.tsx # 회원가입 페이지
│   ├── App.tsx         # 메인 앱 컴포넌트
│   └── main.tsx        # 진입점
├── public/             # 정적 파일
└── package.json        # 프로젝트 설정
```

## 설치 및 실행 방법

### 1. 저장소 클론

```bash
git clone [저장소 URL]
cd frontend
```

### 2. 의존성 설치

```bash
npm install
```

### 3. 개발 서버 실행

```bash
npm run dev
```

개발 서버가 실행되면 브라우저에서 `http://localhost:5173` (또는 표시된 주소)로 접속하세요.

### 4. 프로덕션 빌드

```bash
npm run build
```

빌드 결과물은 `dist/` 폴더에 생성됩니다.

### 5. 빌드 미리보기

```bash
npm run preview
```

## 사용 가능한 스크립트

- `npm run dev` - 개발 서버 실행
- `npm run build` - 프로덕션 빌드
- `npm run preview` - 빌드 결과 미리보기
- `npm run lint` - 코드 린팅

## 주요 페이지

- **회원가입 페이지** (`register.tsx`) - 신규 회원 등록
- **헬스장 현황 페이지** (`guss.tsx`) - 실시간 헬스장 현황 및 예약
- **관리자 페이지** (`admin.tsx`) - 기구 관리, 예약 현황, 매출 로그

## 주의사항

- 현재 백엔드 연동은 구현되지 않았습니다.
- 페이지 간 이동은 `src/App.tsx`의 `page` 상태를 변경하여 확인할 수 있습니다.
- Tailwind CSS는 CDN과 빌드 방식을 함께 사용하고 있습니다.

## 개발 환경

- Node.js 버전: npm 10.9.3 이상 권장
- React 버전: 정확히 19.2.3 사용
- Tailwind CSS 버전: 정확히 4.1.18 사용

## 라이선스

Private Project
