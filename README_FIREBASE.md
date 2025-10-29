# Firebase Integration Setup Guide

## Firebase 프로젝트 설정

### 1. Firebase 프로젝트 생성
1. [Firebase Console](https://console.firebase.google.com/)에 접속
2. "프로젝트 추가" 클릭
3. 프로젝트 이름 입력 (예: electric-circuit-web)
4. Google Analytics 설정 (선택사항)

### 2. Firebase 서비스 활성화

#### Authentication 설정
1. Firebase Console → Authentication → Sign-in method
2. 이메일/비밀번호 활성화
3. 필요시 다른 제공업체도 활성화 (Google, GitHub 등)

#### Firestore Database 설정
1. Firebase Console → Firestore Database
2. "데이터베이스 만들기" 클릭
3. 보안 규칙 모드 선택:
   - 테스트 모드: 개발용 (30일간 모든 읽기/쓰기 허용)
   - 프로덕션 모드: 보안 규칙 직접 설정
4. 지역 선택 (asia-northeast3 권장)

#### Realtime Database 설정 (선택사항)
1. Firebase Console → Realtime Database
2. "데이터베이스 만들기" 클릭
3. 보안 규칙 설정

### 3. 서비스 계정 키 생성
1. Firebase Console → 프로젝트 설정 (톱니바퀴 아이콘)
2. 서비스 계정 탭
3. "새 비공개 키 생성" 클릭
4. JSON 파일 다운로드
5. 서버 디렉토리에 저장 (예: `server/config/firebase-service-account.json`)

### 4. 환경 변수 설정
`.env` 파일 생성 (`.env.example` 참고):

```bash
# Firebase Configuration
FIREBASE_PROJECT_ID=your-firebase-project-id
FIREBASE_SERVICE_ACCOUNT_KEY=./config/firebase-service-account.json
FIREBASE_DATABASE_URL=https://your-project-id-default-rtdb.firebaseio.com/
```

## API 엔드포인트

### 인증 API
- `POST /api/auth/verify` - Firebase ID 토큰 검증
- `POST /api/auth/create-user` - 새 사용자 생성

### 프로젝트 API (인증 필요)
- `GET /api/projects` - 사용자 프로젝트 목록 조회
- `POST /api/projects/create` - 새 프로젝트 생성
- `PUT /api/projects/update?id={projectId}` - 프로젝트 업데이트
- `DELETE /api/projects/delete?id={projectId}` - 프로젝트 삭제

### 회로 API (인증 필요)
- `GET /api/circuits/firebase?projectId={projectId}` - 프로젝트 회로 목록
- `POST /api/circuits/create` - 새 회로 생성

## 보안 규칙 예시

### Firestore 보안 규칙
```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // 프로젝트 컬렉션: 사용자는 자신의 프로젝트만 접근 가능
    match /projects/{projectId} {
      allow read, write: if request.auth != null && resource.data.userId == request.auth.uid;
      allow create: if request.auth != null && request.resource.data.userId == request.auth.uid;
    }
    
    // 회로 컬렉션: 프로젝트 소유자만 접근 가능
    match /circuits/{circuitId} {
      allow read, write: if request.auth != null && 
        exists(/databases/$(database)/documents/projects/$(resource.data.projectId)) &&
        get(/databases/$(database)/documents/projects/$(resource.data.projectId)).data.userId == request.auth.uid;
      allow create: if request.auth != null && 
        exists(/databases/$(database)/documents/projects/$(request.resource.data.projectId)) &&
        get(/databases/$(database)/documents/projects/$(request.resource.data.projectId)).data.userId == request.auth.uid;
    }
  }
}
```

### Realtime Database 보안 규칙
```json
{
  "rules": {
    "users": {
      "$uid": {
        ".read": "$uid === auth.uid",
        ".write": "$uid === auth.uid"
      }
    },
    "projects": {
      "$projectId": {
        ".read": "auth != null && data.child('userId').val() === auth.uid",
        ".write": "auth != null && data.child('userId').val() === auth.uid"
      }
    }
  }
}
```

## 클라이언트 통합

### Firebase SDK 설치 (클라이언트)
```bash
npm install firebase
```

### Firebase 설정 파일 (클라이언트)
```javascript
// firebase.config.js
import { initializeApp } from 'firebase/app'
import { getAuth } from 'firebase/auth'
import { getFirestore } from 'firebase/firestore'

const firebaseConfig = {
  apiKey: "your-api-key",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project-id",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "123456789",
  appId: "your-app-id"
}

const app = initializeApp(firebaseConfig)
export const auth = getAuth(app)
export const db = getFirestore(app)
```

## 테스트

### 서버 실행
```bash
cd server
go run cmd/app/main.go
```

### API 테스트 예시
```bash
# 토큰 검증
curl -X POST http://localhost:8080/api/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"token": "your-firebase-id-token"}'

# 프로젝트 조회 (인증 필요)
curl -X GET http://localhost:8080/api/projects \
  -H "Authorization: Bearer your-firebase-id-token"
```

## 주의사항

1. **서비스 계정 키 보안**: JSON 키 파일을 GitHub에 커밋하지 마세요
2. **보안 규칙**: 프로덕션에서는 반드시 적절한 보안 규칙을 설정하세요
3. **CORS 설정**: 클라이언트에서 API 호출시 CORS 설정이 필요할 수 있습니다
4. **요금**: Firebase 사용량을 모니터링하여 예상치 못한 요금이 발생하지 않도록 주의하세요