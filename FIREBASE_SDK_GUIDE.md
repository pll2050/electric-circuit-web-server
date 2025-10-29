# Firebase SDK 인증 가이드

이 문서는 Electric Circuit Web API와 Firebase SDK를 연동하는 방법을 설명합니다.

## 개요

Electric Circuit Web API는 Firebase Authentication을 사용하며, **클라이언트에서 Firebase SDK로 직접 인증하는 방식**을 권장합니다.

### 왜 Firebase SDK 방식인가?

1. **보안**: 서버에서 비밀번호를 직접 처리하지 않음
2. **편의성**: Firebase가 제공하는 다양한 인증 방법 활용 가능
3. **확장성**: 소셜 로그인, 다중 인증 등 쉽게 추가 가능
4. **표준화**: Firebase의 권장 패턴

## 빠른 시작

### 1. Firebase 프로젝트 설정

```bash
npm install firebase
```

### 2. Firebase 초기화

```javascript
// firebase-config.js
import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';

const firebaseConfig = {
  apiKey: "your-api-key",
  authDomain: "your-project.firebaseapp.com",
  projectId: "your-project-id",
  // ... 기타 설정
};

const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
```

### 3. 인증 서비스 구현

```javascript
// auth-service.js
import { 
  signInWithEmailAndPassword, 
  createUserWithEmailAndPassword,
  signOut,
  onAuthStateChanged 
} from 'firebase/auth';
import { auth } from './firebase-config.js';

class AuthService {
  constructor(apiBaseUrl = 'http://localhost:8080/api') {
    this.apiBaseUrl = apiBaseUrl;
    this.token = null;
    
    // 인증 상태 변화 감시
    onAuthStateChanged(auth, this.handleAuthStateChange.bind(this));
  }

  async handleAuthStateChange(user) {
    if (user) {
      // 사용자 로그인됨 - 토큰 획득
      this.token = await user.getIdToken();
      console.log('사용자 로그인됨:', user.email);
    } else {
      // 사용자 로그아웃됨
      this.token = null;
      console.log('사용자 로그아웃됨');
    }
  }

  // 로그인
  async login(email, password) {
    try {
      const userCredential = await signInWithEmailAndPassword(auth, email, password);
      const user = userCredential.user;
      this.token = await user.getIdToken();
      
      // 서버에서 토큰 검증 (선택사항)
      await this.verifyToken();
      
      return user;
    } catch (error) {
      throw new Error(`로그인 실패: ${error.message}`);
    }
  }

  // 회원가입
  async register(email, password) {
    try {
      const userCredential = await createUserWithEmailAndPassword(auth, email, password);
      const user = userCredential.user;
      this.token = await user.getIdToken();
      
      return user;
    } catch (error) {
      throw new Error(`회원가입 실패: ${error.message}`);
    }
  }

  // 로그아웃
  async logout() {
    try {
      await signOut(auth);
      this.token = null;
    } catch (error) {
      throw new Error(`로그아웃 실패: ${error.message}`);
    }
  }

  // 서버에서 토큰 검증
  async verifyToken() {
    if (!this.token) {
      throw new Error('토큰이 없습니다');
    }

    const response = await fetch(`${this.apiBaseUrl}/auth/verify`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ token: this.token })
    });

    if (!response.ok) {
      throw new Error('토큰 검증 실패');
    }

    return await response.json();
  }

  // API 호출용 헤더 생성
  getAuthHeaders() {
    return {
      'Content-Type': 'application/json',
      ...(this.token && { 'Authorization': `Bearer ${this.token}` })
    };
  }

  // 현재 사용자
  getCurrentUser() {
    return auth.currentUser;
  }

  // 토큰 갱신
  async refreshToken() {
    const user = auth.currentUser;
    if (user) {
      this.token = await user.getIdToken(true);
      return this.token;
    }
    throw new Error('사용자가 로그인되지 않았습니다');
  }
}

export default AuthService;
```

### 4. API 클라이언트 구현

```javascript
// api-client.js
class ApiClient {
  constructor(baseUrl, authService) {
    this.baseUrl = baseUrl;
    this.authService = authService;
  }

  async request(endpoint, options = {}) {
    const url = `${this.baseUrl}${endpoint}`;
    const config = {
      ...options,
      headers: {
        ...this.authService.getAuthHeaders(),
        ...options.headers
      }
    };

    try {
      const response = await fetch(url, config);
      
      // 401 에러시 토큰 갱신 후 재시도
      if (response.status === 401) {
        await this.authService.refreshToken();
        config.headers = {
          ...this.authService.getAuthHeaders(),
          ...options.headers
        };
        return await fetch(url, config);
      }

      if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API 요청 실패:', error);
      throw error;
    }
  }

  // GET 요청
  async get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  }

  // POST 요청
  async post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }

  // PUT 요청
  async put(endpoint, data) {
    return this.request(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data)
    });
  }

  // DELETE 요청
  async delete(endpoint) {
    return this.request(endpoint, { method: 'DELETE' });
  }
}

export default ApiClient;
```

### 5. 실제 사용 예제

```javascript
// app.js
import AuthService from './auth-service.js';
import ApiClient from './api-client.js';

// 서비스 초기화
const authService = new AuthService();
const apiClient = new ApiClient('http://localhost:8080/api', authService);

// 사용 예제
async function example() {
  try {
    // 1. 로그인
    await authService.login('user@example.com', 'password123');
    
    // 2. 프로젝트 목록 조회
    const projects = await apiClient.get('/projects');
    console.log('프로젝트 목록:', projects);
    
    // 3. 새 프로젝트 생성
    const newProject = await apiClient.post('/projects/create', {
      name: '새 프로젝트',
      description: '프로젝트 설명'
    });
    console.log('새 프로젝트:', newProject);
    
  } catch (error) {
    console.error('에러 발생:', error);
  }
}

example();
```

## React에서 사용하기

```jsx
// useAuth.js (React Hook)
import { useState, useEffect, useContext, createContext } from 'react';
import AuthService from './auth-service';

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const authService = new AuthService();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setLoading(false);
    });

    return unsubscribe;
  }, []);

  const login = async (email, password) => {
    return await authService.login(email, password);
  };

  const register = async (email, password) => {
    return await authService.register(email, password);
  };

  const logout = async () => {
    return await authService.logout();
  };

  const value = {
    user,
    login,
    register,
    logout,
    authService
  };

  return (
    <AuthContext.Provider value={value}>
      {!loading && children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
```

```jsx
// App.jsx
import { useAuth } from './useAuth';

function App() {
  const { user, login, logout } = useAuth();

  if (!user) {
    return <LoginForm onLogin={login} />;
  }

  return (
    <div>
      <h1>환영합니다, {user.email}!</h1>
      <button onClick={logout}>로그아웃</button>
      <ProjectList />
    </div>
  );
}
```

## Vue.js에서 사용하기

```javascript
// composables/useAuth.js
import { ref, onMounted } from 'vue';
import { onAuthStateChanged } from 'firebase/auth';
import { auth } from '../firebase-config';
import AuthService from '../auth-service';

export function useAuth() {
  const user = ref(null);
  const loading = ref(true);
  const authService = new AuthService();

  onMounted(() => {
    onAuthStateChanged(auth, (firebaseUser) => {
      user.value = firebaseUser;
      loading.value = false;
    });
  });

  const login = async (email, password) => {
    return await authService.login(email, password);
  };

  const register = async (email, password) => {
    return await authService.register(email, password);
  };

  const logout = async () => {
    return await authService.logout();
  };

  return {
    user,
    loading,
    login,
    register,
    logout,
    authService
  };
}
```

## 주의사항

### 토큰 만료 처리
Firebase ID 토큰은 1시간 후 만료됩니다. 위의 예제에서는 자동으로 토큰을 갱신하지만, 추가 처리가 필요할 수 있습니다.

### 오프라인 처리
네트워크가 불안정한 환경에서는 적절한 에러 처리와 재시도 로직이 필요합니다.

### 보안
- API 키는 환경 변수로 관리하세요
- HTTPS를 사용하세요
- Firebase Security Rules를 적절히 설정하세요

## 문제 해결

### 일반적인 에러들

1. **"Firebase: Error (auth/network-request-failed)"**
   - 네트워크 연결을 확인하세요
   - Firebase 프로젝트 설정을 확인하세요

2. **"Firebase: Error (auth/invalid-api-key)"**
   - Firebase 설정의 API 키를 확인하세요

3. **401 Unauthorized from API**
   - 토큰이 만료되었거나 유효하지 않습니다
   - 토큰 갱신을 시도하세요

4. **CORS 에러**
   - 서버에서 CORS 설정을 확인하세요
   - Firebase 프로젝트의 인증된 도메인을 확인하세요

이 가이드를 따라하시면 Firebase SDK와 Electric Circuit Web API를 성공적으로 연동하실 수 있습니다!