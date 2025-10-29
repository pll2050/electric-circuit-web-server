# Electric Circuit Web - API 명세서

## 개요

Electric Circuit Web 백엔드 API는 전자 회로 설계 웹 애플리케이션을 위한 RESTful API입니다.
Clean Architecture 패턴으로 구성되어 있으며, Firebase를 백엔드 서비스로 사용합니다.

- **Base URL**: `http://localhost:8080/api`
- **Protocol**: HTTP/HTTPS
- **Data Format**: JSON
- **Authentication**: Firebase ID Token

## 공통 응답 형식

### 성공 응답
```json
{
  "success": true,
  "message": "작업이 성공적으로 완료되었습니다",
  "data": {...}
}
```

### 에러 응답
```json
{
  "success": false,
  "error": "에러 메시지",
  "message": "추가 설명"
}
```

## 인증 방식

이 API는 Firebase Authentication을 사용하여 사용자 인증을 처리합니다.

### 권장 인증 흐름 (Firebase SDK 방식)

Firebase SDK를 사용한 클라이언트 측 인증이 권장되며 일반적인 방법입니다:

1. **클라이언트에서 Firebase SDK로 직접 로그인**
   ```javascript
   import { signInWithEmailAndPassword, getAuth } from 'firebase/auth';
   
   const auth = getAuth();
   const userCredential = await signInWithEmailAndPassword(auth, email, password);
   const idToken = await userCredential.user.getIdToken();
   ```

2. **ID 토큰을 서버로 전송하여 검증**
   ```javascript
   const response = await fetch('/api/auth/verify', {
     method: 'POST',
     headers: {
       'Content-Type': 'application/json',
     },
     body: JSON.stringify({ token: idToken })
   });
   ```

3. **이후 API 요청에 토큰 포함**
   ```
   Authorization: Bearer <firebase_id_token>
   ```

### 인증이 필요한 API
대부분의 API 엔드포인트는 Firebase ID Token을 통한 인증이 필요합니다. (`/health`, `/auth/verify`, `/auth/create-user` 제외)

### 토큰 갱신
Firebase ID 토큰은 1시간 후 만료되므로, 클라이언트에서 자동으로 토큰을 갱신해야 합니다:
```javascript
// Firebase가 자동으로 토큰을 갱신합니다
const freshToken = await user.getIdToken(true);
```

---

## 1. 시스템 API

### 1.1 헬스 체크

서버 상태를 확인합니다.

**Endpoint**: `GET /health`  
**인증**: 불필요  

#### 응답
```json
{
  "status": "healthy"
}
```

---

## 2. 인증 API

Firebase 인증 관련 기능을 제공합니다.

### 2.1 토큰 검증 (권장 방법)

**클라이언트에서 Firebase SDK로 로그인한 후, ID 토큰을 서버에서 검증합니다.**

**Endpoint**: `POST /auth/verify`  
**인증**: 불필요  
**용도**: Firebase SDK 인증 후 서버 측 검증

#### 사용 시나리오
1. 클라이언트에서 Firebase SDK로 로그인
2. 획득한 ID 토큰을 이 엔드포인트로 전송
3. 서버에서 토큰 유효성 검증 및 사용자 정보 반환

#### 요청 본문
```json
{
  "token": "firebase_id_token"
}
```

#### 응답
```json
{
  "success": true,
  "message": "Token verified successfully",
  "user": {
    "uid": "user_id",
    "email": "user@example.com",
    "emailVerified": true,
    "displayName": "User Name"
  }
}
```

### 2.2 사용자 생성 (서버 측 생성)

**서버에서 직접 사용자를 생성합니다. (관리자용 또는 특별한 경우)**

**Endpoint**: `POST /auth/create-user`  
**인증**: 불필요  
**참고**: 일반적으로는 클라이언트에서 Firebase SDK로 계정 생성을 권장  

#### 요청 본문
```json
{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "User Name",
  "photo_url": "https://example.com/photo.jpg"
}
```

#### 응답
```json
{
  "success": true,
  "message": "User created successfully",
  "user": {
    "uid": "new_user_id",
    "email": "user@example.com",
    "displayName": "User Name"
  }
}
```

### 2.3 사용자 정보 조회

사용자 정보를 조회합니다.

**Endpoint**: `GET /auth/get-user?uid={user_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `uid` (required): 사용자 ID

#### 응답
```json
{
  "success": true,
  "message": "User found",
  "user": {
    "uid": "user_id",
    "email": "user@example.com",
    "displayName": "User Name",
    "photoURL": "https://example.com/photo.jpg"
  }
}
```

### 2.4 사용자 정보 수정

사용자 정보를 수정합니다.

**Endpoint**: `PUT /auth/update-user`  
**인증**: 필요  

#### 요청 본문
```json
{
  "display_name": "New Display Name",
  "photo_url": "https://example.com/new-photo.jpg"
}
```

#### 응답
```json
{
  "success": true,
  "message": "User updated successfully",
  "user": {
    "uid": "user_id",
    "displayName": "New Display Name",
    "photoURL": "https://example.com/new-photo.jpg"
  }
}
```

### 2.5 사용자 삭제

사용자를 삭제합니다.

**Endpoint**: `DELETE /auth/delete-user?uid={user_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `uid` (required): 삭제할 사용자 ID

#### 응답
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

### 2.6 커스텀 클레임 설정

사용자에게 커스텀 클레임을 설정합니다.

**Endpoint**: `POST /auth/set-custom-claims`  
**인증**: 필요 (관리자 권한)  

#### 요청 본문
```json
{
  "uid": "user_id",
  "custom_claims": {
    "role": "admin",
    "permissions": ["read", "write"]
  }
}
```

#### 응답
```json
{
  "success": true,
  "message": "Custom claims set successfully"
}
```

---

## 3. 프로젝트 API

회로 설계 프로젝트 관리 기능을 제공합니다.

### 3.1 프로젝트 목록 조회

사용자의 모든 프로젝트를 조회합니다.

**Endpoint**: `GET /projects`  
**인증**: 필요  

#### 응답
```json
{
  "success": true,
  "message": "Projects retrieved successfully",
  "projects": [
    {
      "id": "project_1",
      "name": "LED 회로 설계",
      "description": "기본적인 LED 회로 설계 프로젝트",
      "user_id": "user_id",
      "status": "active",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-16T14:20:00Z"
    }
  ]
}
```

### 3.2 프로젝트 생성

새로운 프로젝트를 생성합니다.

**Endpoint**: `POST /projects/create`  
**인증**: 필요  

#### 요청 본문
```json
{
  "name": "새 프로젝트",
  "description": "프로젝트 설명"
}
```

#### 응답
```json
{
  "success": true,
  "message": "Project created successfully",
  "project": {
    "id": "new_project_id",
    "name": "새 프로젝트",
    "description": "프로젝트 설명",
    "user_id": "user_id",
    "status": "active",
    "created_at": "2024-01-17T09:00:00Z",
    "updated_at": "2024-01-17T09:00:00Z"
  }
}
```

### 3.3 프로젝트 상세 조회

특정 프로젝트의 상세 정보를 조회합니다.

**Endpoint**: `GET /projects/get?projectId={project_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `projectId` (required): 프로젝트 ID

#### 응답
```json
{
  "success": true,
  "message": "Project retrieved successfully",
  "project": {
    "id": "project_id",
    "name": "LED 회로 설계",
    "description": "기본적인 LED 회로 설계 프로젝트",
    "user_id": "user_id",
    "status": "active",
    "settings": {
      "grid_size": 10,
      "snap_to_grid": true
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z"
  }
}
```

### 3.4 프로젝트 수정

프로젝트 정보를 수정합니다.

**Endpoint**: `PUT /projects/update`  
**인증**: 필요  

#### 요청 본문
```json
{
  "project_id": "project_id",
  "name": "수정된 프로젝트 이름",
  "description": "수정된 설명"
}
```

#### 응답
```json
{
  "success": true,
  "message": "Project updated successfully",
  "project": {
    "id": "project_id",
    "name": "수정된 프로젝트 이름",
    "description": "수정된 설명",
    "updated_at": "2024-01-17T10:00:00Z"
  }
}
```

### 3.5 프로젝트 삭제

프로젝트를 삭제합니다.

**Endpoint**: `DELETE /projects/delete?projectId={project_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `projectId` (required): 삭제할 프로젝트 ID

#### 응답
```json
{
  "success": true,
  "message": "Project deleted successfully"
}
```

### 3.6 프로젝트 복제

기존 프로젝트를 복제합니다.

**Endpoint**: `POST /projects/duplicate`  
**인증**: 필요  

#### 요청 본문
```json
{
  "project_id": "original_project_id",
  "name": "복제된 프로젝트"
}
```

#### 응답
```json
{
  "success": true,
  "message": "Project duplicated successfully",
  "project": {
    "id": "new_duplicated_project_id",
    "name": "복제된 프로젝트",
    "description": "원본 프로젝트 설명",
    "user_id": "user_id",
    "status": "active",
    "created_at": "2024-01-17T11:00:00Z"
  }
}
```

---

## 4. 회로 API

전자 회로 설계 및 관리 기능을 제공합니다.

### 4.1 프로젝트 회로 목록 조회

특정 프로젝트의 모든 회로를 조회합니다.

**Endpoint**: `GET /circuits?projectId={project_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `projectId` (required): 프로젝트 ID

#### 응답
```json
{
  "success": true,
  "message": "Circuits retrieved successfully",
  "circuits": [
    {
      "id": "circuit_1",
      "name": "메인 회로",
      "description": "프로젝트의 주요 회로",
      "project_id": "project_id",
      "user_id": "user_id",
      "version": 1,
      "is_template": false,
      "tags": ["led", "resistor"],
      "created_at": "2024-01-15T11:00:00Z",
      "updated_at": "2024-01-16T15:30:00Z"
    }
  ]
}
```

### 4.2 회로 생성

새로운 회로를 생성합니다.

**Endpoint**: `POST /circuits/create`  
**인증**: 필요  

#### 요청 본문
```json
{
  "project_id": "project_id",
  "name": "새 회로",
  "description": "회로 설명",
  "data": {
    "elements": [
      {
        "type": "resistor",
        "value": "1k",
        "position": {"x": 100, "y": 100}
      }
    ],
    "connections": []
  }
}
```

#### 응답
```json
{
  "success": true,
  "message": "Circuit created successfully",
  "circuit": {
    "id": "new_circuit_id",
    "name": "새 회로",
    "description": "회로 설명",
    "project_id": "project_id",
    "user_id": "user_id",
    "data": {
      "elements": [...],
      "connections": []
    },
    "version": 1,
    "created_at": "2024-01-17T12:00:00Z"
  }
}
```

### 4.3 회로 상세 조회

특정 회로의 상세 정보를 조회합니다.

**Endpoint**: `GET /circuits/get?circuitId={circuit_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `circuitId` (required): 회로 ID

#### 응답
```json
{
  "success": true,
  "message": "Circuit retrieved successfully",
  "circuit": {
    "id": "circuit_id",
    "name": "메인 회로",
    "description": "프로젝트의 주요 회로",
    "project_id": "project_id",
    "user_id": "user_id",
    "data": {
      "elements": [
        {
          "id": "R1",
          "type": "resistor",
          "value": "1k",
          "position": {"x": 100, "y": 100}
        },
        {
          "id": "LED1",
          "type": "led",
          "color": "red",
          "position": {"x": 200, "y": 100}
        }
      ],
      "connections": [
        {
          "from": "R1.pin2",
          "to": "LED1.pin1"
        }
      ]
    },
    "version": 1,
    "tags": ["led", "resistor"],
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-16T15:30:00Z"
  }
}
```

### 4.4 회로 수정

회로 정보 및 설계 데이터를 수정합니다.

**Endpoint**: `PUT /circuits/update`  
**인증**: 필요  

#### 요청 본문
```json
{
  "circuit_id": "circuit_id",
  "name": "수정된 회로",
  "description": "수정된 설명",
  "data": {
    "elements": [...],
    "connections": [...]
  }
}
```

#### 응답
```json
{
  "success": true,
  "message": "Circuit updated successfully",
  "circuit": {
    "id": "circuit_id",
    "name": "수정된 회로",
    "description": "수정된 설명",
    "version": 2,
    "updated_at": "2024-01-17T13:00:00Z"
  }
}
```

### 4.5 회로 삭제

회로를 삭제합니다.

**Endpoint**: `DELETE /circuits/delete?circuitId={circuit_id}`  
**인증**: 필요  

#### 쿼리 파라미터
- `circuitId` (required): 삭제할 회로 ID

#### 응답
```json
{
  "success": true,
  "message": "Circuit deleted successfully"
}
```

### 4.6 회로 템플릿 목록 조회

사용 가능한 회로 템플릿을 조회합니다.

**Endpoint**: `GET /circuits/templates`  
**인증**: 필요  

#### 응답
```json
{
  "success": true,
  "message": "Templates feature not yet implemented",
  "templates": []
}
```

### 4.7 템플릿으로부터 회로 생성

기존 템플릿을 사용하여 새 회로를 생성합니다.

**Endpoint**: `POST /circuits/create-from-template`  
**인증**: 필요  

#### 요청 본문
```json
{
  "project_id": "project_id",
  "template_id": "template_id",
  "name": "템플릿 기반 회로"
}
```

#### 응답
```json
{
  "success": true,
  "message": "Circuit created from template successfully",
  "circuit": {
    "id": "new_circuit_id",
    "name": "템플릿 기반 회로",
    "project_id": "project_id",
    "created_at": "2024-01-17T14:00:00Z"
  }
}
```

---

## 5. 스토리지 API

파일 업로드 및 관리 기능을 제공합니다.

### 5.1 파일 업로드

파일을 서버에 업로드합니다.

**Endpoint**: `POST /storage/upload`  
**인증**: 필요  
**Content-Type**: `multipart/form-data`

#### 요청 본문 (FormData)
- `file` (required): 업로드할 파일
- `folder` (optional): 저장할 폴더명

#### 응답
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "download_url": "https://storage.example.com/files/file_id",
  "file_path": "/uploads/user_id/filename.ext",
  "file_name": "filename.ext",
  "size": 1024
}
```

### 5.2 파일 URL 조회

저장된 파일의 다운로드 URL을 조회합니다.

**Endpoint**: `GET /storage/url?filePath={file_path}`  
**인증**: 필요  

#### 쿼리 파라미터
- `filePath` (required): 파일 경로

#### 응답
```json
{
  "success": true,
  "message": "File URL retrieved successfully",
  "download_url": "https://storage.example.com/files/file_id"
}
```

### 5.3 파일 삭제

저장된 파일을 삭제합니다.

**Endpoint**: `DELETE /storage/delete?filePath={file_path}`  
**인증**: 필요  

#### 쿼리 파라미터
- `filePath` (required): 삭제할 파일 경로

#### 응답
```json
{
  "success": true,
  "message": "File deleted successfully"
}
```

### 5.4 파일 목록 조회

사용자의 파일 목록을 조회합니다.

**Endpoint**: `GET /storage/list?folder={folder_name}`  
**인증**: 필요  

#### 쿼리 파라미터
- `folder` (optional): 조회할 폴더명

#### 응답
```json
{
  "success": true,
  "message": "Files listed successfully",
  "files": [
    {
      "name": "circuit_image.png",
      "path": "/uploads/user_id/circuit_image.png",
      "size": 2048,
      "type": "image/png",
      "url": "https://storage.example.com/files/file_id",
      "created_at": "2024-01-17T15:00:00Z"
    }
  ]
}
```

### 5.5 회로 이미지 업로드

회로 이미지를 업로드합니다.

**Endpoint**: `POST /storage/upload-circuit-image`  
**인증**: 필요  
**Content-Type**: `multipart/form-data`

#### 요청 본문 (FormData)
- `image` (required): 회로 이미지 파일
- `circuitId` (required): 연결할 회로 ID

#### 응답
```json
{
  "success": true,
  "message": "Circuit image uploaded successfully",
  "download_url": "https://storage.example.com/circuits/circuit_id/image.png",
  "file_path": "/circuits/circuit_id/image.png"
}
```

---

## HTTP 상태 코드

| 코드 | 설명 |
|------|------|
| 200 | 성공 |
| 201 | 생성됨 |
| 400 | 잘못된 요청 |
| 401 | 인증 필요 |
| 403 | 권한 없음 |
| 404 | 리소스 없음 |
| 405 | 허용되지 않는 메서드 |
| 500 | 서버 내부 오류 |

## 에러 처리

모든 API는 일관된 에러 응답 형식을 사용합니다:

```json
{
  "success": false,
  "error": "구체적인 에러 메시지",
  "message": "사용자에게 표시할 메시지"
}
```

### 일반적인 에러 메시지

- `Token is required`: 인증 토큰이 필요합니다
- `User not authenticated`: 사용자 인증이 실패했습니다
- `Invalid JSON`: 요청 본문이 올바른 JSON 형식이 아닙니다
- `Method not allowed`: 지원하지 않는 HTTP 메서드입니다
- `Project ID is required`: 프로젝트 ID가 필요합니다
- `Circuit ID is required`: 회로 ID가 필요합니다
- `File path is required`: 파일 경로가 필요합니다

## 개발 정보

- **서버 포트**: 8080
- **CORS**: 모든 엔드포인트에서 지원
- **요청 제한**: 현재 제한 없음 (추후 Rate Limiting 적용 예정)
- **파일 업로드 제한**: 32MB

## 변경 사항

### v1.0.0 (2024-01-17)
- 초기 API 명세서 작성
- Clean Architecture 기반 엔드포인트 구현
- Firebase 인증 통합
- 프로젝트, 회로, 스토리지 API 구현