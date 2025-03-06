import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { getAccessToken } from "../utils/tokens";
import axiosInstance from "../utils/axios";

const PrivateRoute = ({ element }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const accessToken = getAccessToken();

  useEffect(() => {
    if (accessToken) {
      axiosInstance
        .get("/user/auth")
        .then((response) => {
          setIsAuthenticated(true);
          const userRole = response.data.role;
          localStorage.setItem("role", userRole);
        })
        .catch(() => {
          setIsAuthenticated(false);
        });
    } else {
      setIsAuthenticated(false);
    }
  }, [accessToken]);

  if (isAuthenticated === null) {
    return <p>Loading...</p>;
  }

  if (!isAuthenticated) {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("role");
    return <Navigate to="/login" replace />;
  }

  return element;
};

export default PrivateRoute;
