import axios from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
});

api.defaults.withCredentials = true;

export const setUpRequestInterceptors = (getAccessToken: () => string | null) => {
  api.interceptors.request.use(
    (config) => {
      const token = getAccessToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
    (error) => Promise.reject(error),
  );
};

// automatically refresh token on request fail
export const setUpResponseInterceptors = (setAccessToken: (token: string | null) => void) => {
  let isRefreshing = false;

  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      // Prevent infinite loop by checking if the failing request is the refresh endpoint itself
      const isRefreshEndpoint = error.config.url === "users/refresh";
      
      if (error.response?.status == 401 && !isRefreshing && !isRefreshEndpoint) {
        try {
          isRefreshing = true;
          const res = await api.post("users/refresh");
          isRefreshing = false;

          setAccessToken(res.data.access);
          error.config.headers.Authorization = `Bearer ${res.data.access}`;
          return api.request(error.config);
        } catch (refreshError) {
          isRefreshing = false;
          setAccessToken(null);
          window.location.href = "/login";
          return Promise.reject(refreshError);
        }
      }
      return Promise.reject(error);
    },
  );
};