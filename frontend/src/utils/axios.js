import axios from "axios";
import { API_BASE_URL } from "../config";
import {
  getAccessToken,
  setAccessToken,
  setRefreshToken,
  getRefreshToken,
} from "./tokens.js";

const axiosInstance = axios.create({
  baseURL: API_BASE_URL, // Your API base URL
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to attach access token to the Authorization header
axiosInstance.interceptors.request.use(
  (config) => {
    const accessToken = getAccessToken();
    if (accessToken) {
      config.headers["Authorization"] = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle token refresh on 401 errors
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config; // Original request that caused the error

    // Check if the error is a 401 and if a refresh token exists
    if (
      error.response.status === 401 &&
      !originalRequest._retry &&
      error.response?.data?.error?.includes("access token")
    ) {
      originalRequest._retry = true;

      const refreshToken = getRefreshToken();

      if (!refreshToken) {
        window.location.href = "/login";
        return Promise.reject(error);
      }

      try {
        const response = await axios.post(`${API_BASE_URL}/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token } = response.data;

        setAccessToken(access_token);
        setRefreshToken(refresh_token);

        originalRequest.headers["Authorization"] = `Bearer ${access_token}`;

        return axiosInstance(originalRequest);
      } catch (refreshError) {
        // If refresh failed, remove tokens and redirect to login
        localStorage.removeItem("accessToken");
        localStorage.removeItem("refreshToken");
        localStorage.removeItem("role");
        window.location.href = "/login";
        return Promise.reject(refreshError);
      }
    }

    if (error.response?.status === 401) {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("role");
      alert("Please log in again!");
      window.location.href = "/login";
      return Promise.reject(error);
    }

    if (
      error.response?.status === 403 &&
      error.response?.data?.error?.includes("banned")
    ) {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("role");
      alert("You are banned!");
      window.location.href = "/";
      return Promise.reject(error);
    }

    if (
      error.response?.status === 403 &&
      error.response?.data?.error?.includes("Admins only")
    ) {
      window.location.href = "/";
      return Promise.reject(error);
    }

    // If the error is not a 401 or another issue, reject the promise
    window.location.href = "/";
    return Promise.reject(error);
  }
);

export default axiosInstance;
