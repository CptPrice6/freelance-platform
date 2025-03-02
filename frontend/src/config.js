const isLocalhost = window.location.hostname === "localhost";

export const API_BASE_URL = isLocalhost
  ? "http://localhost:8080" // Adjusted to match your backend port
  : `http://${window.location.hostname}:${
      process.env.REACT_APP_BACKEND_EXTERNAL_PORT || 8080
    }`;
