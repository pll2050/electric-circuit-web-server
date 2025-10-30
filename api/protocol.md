# Electric Circuit Web API Protocols

이 문서는 Electric Circuit Web 서버의 API 프로토콜, 인증 방식, 데이터 포맷, 에러 처리, 스펙 파일 활용법을 설명합니다.

---

## 1. 인증 방식

- **Firebase Authentication**
  - 모든 주요 API는 Firebase ID Token을 통한 인증이 필요합니다.
  - 클라이언트에서 Firebase SDK로 로그인 후, ID 토큰을 획득하여 서버에 전송합니다.
  - 인증 헤더 예시:
    ```http
    Authorization: Bearer <firebase_id_token>
    ```

- **토큰 검증 엔드포인트**
  - `/auth/verify` : 클라이언트에서 획득한 토큰을 서버에서 검증

---

## 2. 데이터 포맷

- **JSON**
  - 모든 요청/응답은 기본적으로 JSON 포맷을 사용합니다.
  - 파일 업로드는 `multipart/form-data`를 사용합니다.

---

## 3. API 엔드포인트 예시

- `/api/projects` : 프로젝트 목록 조회
- `/api/circuits` : 회로 목록 및 상세 조회
- `/api/storage/upload` : 파일 업로드
- `/api/auth/verify` : 토큰 검증

---

## 4. 에러 처리

- 모든 에러 응답은 아래와 같은 형식을 따릅니다:
  ```json
  {
    "success": false,
    "error": "에러 메시지",
    "message": "추가 설명"
  }
  ```
- HTTP 상태 코드와 함께 상세 에러 메시지 제공

---

## 5. 주요 스펙 파일 설명

- **openapi.yaml** : OpenAPI 3.0 전체 명세 (Swagger/Redoc 등에서 시각화)
- **swagger.json** : Swagger 2.0 예시 스펙 (일부 툴 호환용)
- **schema.json** : 주요 데이터(JSON) 스키마 (유효성 검증, 타입 자동 생성)
- **README.md** : 폴더 및 파일 설명

---

## 6. 활용 방법

- API 문서 자동화 및 시각화: Swagger UI, Redoc 등에서 `openapi.yaml`을 불러와 사용
- 데이터 검증: `schema.json`을 활용해 입력/출력 데이터의 유효성 체크
- 협업 및 문서화: 이 문서와 README를 참고해 API 설계 및 개발 진행

---

## 7. 참고

- 모든 파일은 버전 관리되며, 변경 시 반드시 문서화 및 테스트를 병행하세요.
- API 설계 변경 시, 스펙 파일과 프로토콜 문서를 반드시 최신화하세요.
