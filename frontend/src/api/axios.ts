import axios, { InternalAxiosRequestConfig } from 'axios';

const api = axios.create({
  // Vite 환경 변수 읽기 (타입 단언 추가)
  baseURL: (import.meta.env.VITE_API_URL as string) || 'http://localhost:9000',
});

// 요청 인터셉터에 타입 지정
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default api;