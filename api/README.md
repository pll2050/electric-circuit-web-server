# API 스펙 및 프로토콜 폴더

이 폴더에는 Electric Circuit Web 서버의 API 스펙, JSON 스키마, 프로토콜 정의 파일들이 포함됩니다.

## 파일 설명

- **openapi.yaml**
  - Electric Circuit Web API의 전체 명세(OpenAPI 3.0)
  - 모든 엔드포인트, 요청/응답, 데이터 모델, 인증 방식 등 상세 명세 포함
  - Swagger UI, Redoc 등에서 시각화 및 테스트 가능

- **swagger.json**
  - Swagger 2.0 형식의 API 명세 예시
  - 일부 툴에서 호환이 필요할 때 사용

- **schema.json**
  - 주요 데이터 구조(프로젝트, 회로 등)에 대한 JSON 스키마
  - 데이터 유효성 검증, 프론트엔드/백엔드 타입 자동 생성 등에 활용

- **protocol.md**
  - API 인증 방식, 데이터 포맷, 에러 응답 등 프로토콜 설명
  - 각 스펙 파일의 역할 및 사용법 안내

## 활용 예시

- API 문서 자동화 및 시각화: Swagger UI, Redoc 등에서 `openapi.yaml`을 불러와 사용
- 데이터 검증: `schema.json`을 활용해 입력/출력 데이터의 유효성 체크
- 협업 및 문서화: `protocol.md`와 이 README를 참고해 API 설계 및 개발 진행

## 참고

- 실제 서비스에서는 `openapi.yaml`만으로도 Swagger UI, Redoc 등에서 API 문서 시각화 및 테스트가 가능합니다.
- 모든 파일은 버전 관리되며, 변경 시 반드시 문서화 및 테스트를 병행하세요.
