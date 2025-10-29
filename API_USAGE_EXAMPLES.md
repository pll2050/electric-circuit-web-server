# API Usage Examples

이 문서는 Electric Circuit Web API의 실제 사용 예제를 제공합니다.

## 환경 설정

### Firebase 클라이언트 설정
```javascript
// firebase-config.js
import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';

const firebaseConfig = {
  apiKey: "your-api-key",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project-id",
  storageBucket: "your-project.appspot.com",
  messagingSenderId: "123456789",
  appId: "your-app-id"
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
```

### 환경 변수
```bash
API_BASE_URL=http://localhost:8080/api
FIREBASE_API_KEY=your_firebase_api_key
```

### JavaScript/TypeScript 환경

#### API 클라이언트 설정
```typescript
// api-client.ts
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080/api';

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setAuthToken(token: string) {
    this.token = token;
  }

  clearAuthToken() {
    this.token = null;
  }

  private async request(endpoint: string, options: RequestInit = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers,
    };

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
  }

  async get(endpoint: string) {
    return this.request(endpoint, { method: 'GET' });
  }

  async post(endpoint: string, data: any) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put(endpoint: string, data: any) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete(endpoint: string) {
    return this.request(endpoint, { method: 'DELETE' });
  }

  async uploadFile(endpoint: string, file: File, folder?: string) {
    const formData = new FormData();
    formData.append('file', file);
    if (folder) {
      formData.append('folder', folder);
    }

    const headers = {
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
    };

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers,
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload Error: ${response.status}`);
    }

    return response.json();
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
```

## 사용 예제

### 1. Firebase SDK 인증 (권장 방법)

#### Firebase 로그인 및 API 연동
```typescript
import { 
  signInWithEmailAndPassword, 
  createUserWithEmailAndPassword, 
  signOut, 
  onAuthStateChanged, 
  updateProfile,
  User 
} from 'firebase/auth';
import { auth } from './firebase-config';

class AuthService {
  private apiClient: ApiClient;

  constructor(apiClient: ApiClient) {
    this.apiClient = apiClient;
  }

  // Firebase SDK로 로그인
  async login(email: string, password: string) {
    try {
      // 1. Firebase에서 로그인
      const userCredential = await signInWithEmailAndPassword(auth, email, password);
      const user = userCredential.user;
      
      // 2. ID 토큰 획득
      const idToken = await user.getIdToken();
      
      // 3. API 클라이언트에 토큰 설정
      this.apiClient.setAuthToken(idToken);
      
      // 4. 서버에서 토큰 검증 (선택사항)
      const response = await this.apiClient.post('/auth/verify', {
        token: idToken
      });
      
      console.log('로그인 성공:', response.user);
      return { user, token: idToken };
      
    } catch (error) {
      console.error('로그인 실패:', error);
      throw error;
    }
  }

  // Firebase SDK로 회원가입
  async register(email: string, password: string, displayName?: string) {
    try {
      // 1. Firebase에서 계정 생성
      const userCredential = await createUserWithEmailAndPassword(auth, email, password);
      const user = userCredential.user;
      
      // 2. 프로필 업데이트 (선택사항)
      if (displayName) {
        await updateProfile(user, { displayName });
      }
      
      // 3. ID 토큰 획득
      const idToken = await user.getIdToken();
      
      // 4. API 클라이언트에 토큰 설정
      this.apiClient.setAuthToken(idToken);
      
      console.log('회원가입 성공:', user);
      return { user, token: idToken };
      
    } catch (error) {
      console.error('회원가입 실패:', error);
      throw error;
    }
  }

  // 로그아웃
  async logout() {
    try {
      await signOut(auth);
      this.apiClient.clearAuthToken();
      console.log('로그아웃 성공');
    } catch (error) {
      console.error('로그아웃 실패:', error);
      throw error;
    }
  }

  // 인증 상태 감시
  onAuthStateChanged(callback: (user: User | null) => void) {
    return onAuthStateChanged(auth, async (user) => {
      if (user) {
        // 사용자가 로그인된 경우 토큰 갱신
        const token = await user.getIdToken();
        this.apiClient.setAuthToken(token);
      } else {
        // 사용자가 로그아웃된 경우 토큰 클리어
        this.apiClient.clearAuthToken();
      }
      callback(user);
    });
  }

  // 토큰 자동 갱신
  async refreshToken() {
    const user = auth.currentUser;
    if (user) {
      const token = await user.getIdToken(true); // force refresh
      this.apiClient.setAuthToken(token);
      return token;
    }
    throw new Error('사용자가 로그인되지 않았습니다.');
  }
}

// 사용 예제
const authService = new AuthService(apiClient);

// 로그인
await authService.login('user@example.com', 'password123');

// 인증 상태 감시
authService.onAuthStateChanged((user) => {
  if (user) {
    console.log('사용자 로그인됨:', user.email);
  } else {
    console.log('사용자 로그아웃됨');
  }
});
```

#### 토큰 자동 갱신이 포함된 API 클라이언트
```typescript
// enhanced-api-client.ts
class EnhancedApiClient extends ApiClient {
  private authService: AuthService;

  constructor(baseURL: string, authService: AuthService) {
    super(baseURL);
    this.authService = authService;
  }

  // 토큰 만료 시 자동 갱신
  private async requestWithTokenRefresh(endpoint: string, options: RequestInit = {}) {
    try {
      return await this.request(endpoint, options);
    } catch (error) {
      // 401 에러 시 토큰 갱신 후 재시도
      if (error.status === 401) {
        try {
          await this.authService.refreshToken();
          return await this.request(endpoint, options);
        } catch (refreshError) {
          // 토큰 갱신 실패 시 로그아웃
          await this.authService.logout();
          throw new Error('인증이 만료되었습니다. 다시 로그인해주세요.');
        }
      }
      throw error;
    }
  }

  // 기존 메서드들을 토큰 갱신 로직으로 오버라이드
  async get(endpoint: string) {
    return this.requestWithTokenRefresh(endpoint, { method: 'GET' });
  }

  async post(endpoint: string, data: any) {
    return this.requestWithTokenRefresh(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}
```
```

### 2. 프로젝트 관리

#### 프로젝트 목록 조회
```typescript
async function getProjectList() {
  try {
    const response = await apiClient.get('/projects');
    console.log('프로젝트 목록:', response.projects);
    return response.projects;
  } catch (error) {
    console.error('프로젝트 조회 실패:', error);
    throw error;
  }
}
```

#### 새 프로젝트 생성
```typescript
async function createProject(name: string, description?: string) {
  try {
    const response = await apiClient.post('/projects/create', {
      name,
      description: description || ''
    });

    console.log('프로젝트 생성 성공:', response.project);
    return response.project;
  } catch (error) {
    console.error('프로젝트 생성 실패:', error);
    throw error;
  }
}
```

#### 프로젝트 수정
```typescript
async function updateProject(projectId: string, updates: { name?: string; description?: string }) {
  try {
    const response = await apiClient.put('/projects/update', {
      project_id: projectId,
      ...updates
    });

    console.log('프로젝트 수정 성공:', response.project);
    return response.project;
  } catch (error) {
    console.error('프로젝트 수정 실패:', error);
    throw error;
  }
}
```

#### 프로젝트 복제
```typescript
async function duplicateProject(originalProjectId: string, newName: string) {
  try {
    const response = await apiClient.post('/projects/duplicate', {
      project_id: originalProjectId,
      name: newName
    });

    console.log('프로젝트 복제 성공:', response.project);
    return response.project;
  } catch (error) {
    console.error('프로젝트 복제 실패:', error);
    throw error;
  }
}
```

### 3. 회로 설계

#### 회로 목록 조회
```typescript
async function getCircuitList(projectId: string) {
  try {
    const response = await apiClient.get(`/circuits?projectId=${projectId}`);
    console.log('회로 목록:', response.circuits);
    return response.circuits;
  } catch (error) {
    console.error('회로 조회 실패:', error);
    throw error;
  }
}
```

#### 새 회로 생성
```typescript
interface CircuitElement {
  id: string;
  type: string;
  value?: string;
  position: { x: number; y: number };
}

interface CircuitConnection {
  from: string;
  to: string;
}

async function createCircuit(
  projectId: string, 
  name: string, 
  description: string,
  elements: CircuitElement[] = [],
  connections: CircuitConnection[] = []
) {
  try {
    const response = await apiClient.post('/circuits/create', {
      project_id: projectId,
      name,
      description,
      data: {
        elements,
        connections
      }
    });

    console.log('회로 생성 성공:', response.circuit);
    return response.circuit;
  } catch (error) {
    console.error('회로 생성 실패:', error);
    throw error;
  }
}
```

#### 회로 저장 (설계 데이터 업데이트)
```typescript
async function saveCircuitDesign(
  circuitId: string,
  elements: CircuitElement[],
  connections: CircuitConnection[]
) {
  try {
    const response = await apiClient.put('/circuits/update', {
      circuit_id: circuitId,
      data: {
        elements,
        connections
      }
    });

    console.log('회로 저장 성공:', response.circuit);
    return response.circuit;
  } catch (error) {
    console.error('회로 저장 실패:', error);
    throw error;
  }
}
```

#### 복잡한 회로 예제
```typescript
async function createLEDCircuit(projectId: string) {
  const elements: CircuitElement[] = [
    {
      id: 'V1',
      type: 'voltage_source',
      value: '5V',
      position: { x: 50, y: 100 }
    },
    {
      id: 'R1',
      type: 'resistor',
      value: '220Ω',
      position: { x: 150, y: 100 }
    },
    {
      id: 'LED1',
      type: 'led',
      value: 'red',
      position: { x: 250, y: 100 }
    }
  ];

  const connections: CircuitConnection[] = [
    { from: 'V1.positive', to: 'R1.pin1' },
    { from: 'R1.pin2', to: 'LED1.anode' },
    { from: 'LED1.cathode', to: 'V1.negative' }
  ];

  return await createCircuit(
    projectId,
    'LED 기본 회로',
    '저항과 LED를 이용한 기본 회로',
    elements,
    connections
  );
}
```

### 4. 파일 스토리지

#### 이미지 업로드
```typescript
async function uploadCircuitImage(file: File, projectId: string) {
  try {
    const response = await apiClient.uploadFile('/storage/upload', file, `projects/${projectId}/images`);
    
    console.log('이미지 업로드 성공:', response.download_url);
    return response;
  } catch (error) {
    console.error('이미지 업로드 실패:', error);
    throw error;
  }
}
```

#### 파일 URL 조회
```typescript
async function getFileUrl(filePath: string) {
  try {
    const response = await apiClient.get(`/storage/url?filePath=${encodeURIComponent(filePath)}`);
    return response.download_url;
  } catch (error) {
    console.error('파일 URL 조회 실패:', error);
    throw error;
  }
}
```

### 5. 에러 처리

#### 통합 에러 핸들러
```typescript
function handleApiError(error: any) {
  if (error.response) {
    // API 에러 응답
    const { status, data } = error.response;
    
    switch (status) {
      case 401:
        // 인증 만료
        console.error('인증이 만료되었습니다. 다시 로그인해주세요.');
        // 로그인 페이지로 리디렉션
        break;
      case 403:
        console.error('권한이 없습니다.');
        break;
      case 404:
        console.error('요청하신 리소스를 찾을 수 없습니다.');
        break;
      case 500:
        console.error('서버 내부 오류가 발생했습니다.');
        break;
      default:
        console.error('API 오류:', data?.error || '알 수 없는 오류');
    }
  } else if (error.request) {
    // 네트워크 오류
    console.error('네트워크 오류가 발생했습니다.');
  } else {
    // 기타 오류
    console.error('오류:', error.message);
  }
}
```

### 6. React 컴포넌트 예제

#### 프로젝트 목록 컴포넌트
```tsx
import React, { useEffect, useState } from 'react';

interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

const ProjectList: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    try {
      setLoading(true);
      const projectList = await getProjectList();
      setProjects(projectList);
    } catch (error) {
      handleApiError(error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateProject = async (name: string, description: string) => {
    try {
      const newProject = await createProject(name, description);
      setProjects(prev => [...prev, newProject]);
    } catch (error) {
      handleApiError(error);
    }
  };

  const handleDeleteProject = async (projectId: string) => {
    try {
      await apiClient.delete(`/projects/delete?projectId=${projectId}`);
      setProjects(prev => prev.filter(p => p.id !== projectId));
    } catch (error) {
      handleApiError(error);
    }
  };

  if (loading) {
    return <div>로딩 중...</div>;
  }

  return (
    <div>
      <h1>프로젝트 목록</h1>
      {projects.map(project => (
        <div key={project.id} className="project-card">
          <h3>{project.name}</h3>
          <p>{project.description}</p>
          <small>생성일: {new Date(project.created_at).toLocaleDateString()}</small>
          <button onClick={() => handleDeleteProject(project.id)}>
            삭제
          </button>
        </div>
      ))}
    </div>
  );
};
```

### 7. Vue.js 컴포저블 예제

#### 프로젝트 관리 컴포저블
```typescript
// composables/useProjects.ts
import { ref, computed } from 'vue';

export function useProjects() {
  const projects = ref<Project[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const projectCount = computed(() => projects.value.length);

  const loadProjects = async () => {
    try {
      loading.value = true;
      error.value = null;
      const data = await getProjectList();
      projects.value = data;
    } catch (err) {
      error.value = '프로젝트를 불러오는데 실패했습니다.';
      console.error(err);
    } finally {
      loading.value = false;
    }
  };

  const createProject = async (name: string, description: string) => {
    try {
      const newProject = await createProject(name, description);
      projects.value.push(newProject);
      return newProject;
    } catch (err) {
      error.value = '프로젝트 생성에 실패했습니다.';
      throw err;
    }
  };

  return {
    projects: readonly(projects),
    loading: readonly(loading),
    error: readonly(error),
    projectCount,
    loadProjects,
    createProject
  };
}
```

## 테스트 예제

### Jest를 사용한 API 테스트
```typescript
// tests/api.test.ts
import { apiClient } from '../api-client';

describe('API Client Tests', () => {
  beforeEach(() => {
    // 테스트용 토큰 설정
    apiClient.setAuthToken('test_token');
  });

  test('should get project list', async () => {
    const projects = await getProjectList();
    expect(Array.isArray(projects)).toBe(true);
  });

  test('should create project', async () => {
    const project = await createProject('Test Project', 'Test Description');
    expect(project).toHaveProperty('id');
    expect(project.name).toBe('Test Project');
  });

  test('should handle API errors', async () => {
    apiClient.setAuthToken('invalid_token');
    
    await expect(getProjectList()).rejects.toThrow();
  });
});
```

이 예제들을 참고하여 Electric Circuit Web API를 효과적으로 활용하실 수 있습니다.