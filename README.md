# Electric Circuit Web - .NET Core API Server

이 프로젝트는 전기회로 설계 및 시뮬레이터의 .NET Core 백엔드 API 서버입니다.

## 기술 스택

- **.NET 9.0**
- **ASP.NET Core Web API**
- **Entity Framework Core 9.0**
- **PostgreSQL** (Npgsql)
- **Firebase Admin SDK** (인증)

## 프로젝트 구조

```
server-dotnet/
├── src/
│   └── ElectricCircuitWeb.API/
│       ├── Controllers/        # API 컨트롤러
│       ├── Services/           # 비즈니스 로직
│       ├── Repositories/       # 데이터 액세스 레이어
│       ├── Models/             # 데이터 모델
│       ├── Data/               # DbContext
│       ├── Middleware/         # 커스텀 미들웨어
│       └── Config/             # 설정 클래스
├── appsettings.json           # 애플리케이션 설정
├── appsettings.Development.json
├── Program.cs                 # 애플리케이션 진입점
└── README.md
```

## 시작하기

### 사전 요구사항

1. **.NET 9.0 SDK** 설치
2. **PostgreSQL** 데이터베이스 실행 중
3. (선택) **Firebase 프로젝트** 설정

### 데이터베이스 설정

Docker를 사용하여 PostgreSQL 실행:

```bash
docker run --name electric-circuit-db -e POSTGRES_PASSWORD=q1w2e3r4 -p 5432:5432 -d postgres
```

### 설정 파일 구성

`appsettings.json` 파일에서 데이터베이스 연결 문자열을 확인하세요:

```json
{
  "ConnectionStrings": {
    "DefaultConnection": "Host=localhost;Port=5432;Database=electric_circuit_db;Username=postgres;Password=q1w2e3r4"
  }
}
```

Firebase를 사용하는 경우:

```json
{
  "Firebase": {
    "ProjectId": "your-project-id",
    "ServiceAccountKeyPath": "path/to/serviceAccountKey.json",
    "DatabaseUrl": "https://your-project.firebaseio.com"
  }
}
```

### 데이터베이스 마이그레이션

Entity Framework Core 마이그레이션 생성:

```bash
# 마이그레이션 생성
dotnet ef migrations add InitialCreate

# 데이터베이스 업데이트
dotnet ef database update
```

> **참고**: 개발 환경에서는 애플리케이션 시작 시 자동으로 마이그레이션이 적용됩니다.

### 애플리케이션 실행

```bash
# 개발 모드로 실행
dotnet run

# 또는
dotnet watch run  # 핫 리로드 활성화
```

서버는 기본적으로 다음 주소에서 실행됩니다:
- HTTPS: `https://localhost:5001`
- HTTP: `http://localhost:5000`

### API 문서

애플리케이션 실행 후 Swagger UI에 접근:
- `http://localhost:5000/swagger`

## 주요 명령어

```bash
# 의존성 복원
dotnet restore

# 빌드
dotnet build

# 테스트 실행
dotnet test

# 릴리스 빌드
dotnet publish -c Release -o ./publish
```

## API 엔드포인트

### Health Check
- `GET /api/health` - 서버 상태 확인

### Authentication
- `POST /api/auth/verify` - Firebase ID 토큰 검증
- `POST /api/auth/signup` - 새 사용자 등록

## 개발 가이드

### 새 컨트롤러 추가

1. `src/ElectricCircuitWeb.API/Controllers/` 폴더에 새 컨트롤러 생성
2. `[ApiController]` 및 `[Route]` 특성 추가
3. 필요한 서비스를 의존성 주입으로 추가

### 새 서비스 추가

1. `src/ElectricCircuitWeb.API/Services/` 폴더에 인터페이스와 구현 클래스 생성
2. `Program.cs`에서 서비스 등록:
   ```csharp
   builder.Services.AddScoped<IYourService, YourService>();
   ```

### 새 모델 추가

1. `src/ElectricCircuitWeb.API/Models/` 폴더에 모델 클래스 생성
2. `ApplicationDbContext`에 DbSet 추가
3. 마이그레이션 생성 및 적용

## 환경 변수

다음 환경 변수를 설정할 수 있습니다:

- `ASPNETCORE_ENVIRONMENT` - 환경 설정 (Development, Staging, Production)
- `ConnectionStrings__DefaultConnection` - 데이터베이스 연결 문자열
- `Firebase__ProjectId` - Firebase 프로젝트 ID
- `Firebase__ServiceAccountKeyPath` - Firebase 서비스 계정 키 경로

## 보안 고려사항

1. **프로덕션 환경**에서는 반드시 강력한 데이터베이스 비밀번호 사용
2. Firebase 서비스 계정 키 파일을 Git에 커밋하지 말 것
3. CORS 정책을 프로덕션 환경에 맞게 조정
4. HTTPS 사용 권장

## 트러블슈팅

### 데이터베이스 연결 오류

PostgreSQL이 실행 중인지 확인:
```bash
docker ps | grep electric-circuit-db
```

### 마이그레이션 오류

EF Core 도구 설치:
```bash
dotnet tool install --global dotnet-ef
```

## 라이선스

이 프로젝트는 비공개 프로젝트입니다.

## 기여

프로젝트 팀원만 기여할 수 있습니다.
