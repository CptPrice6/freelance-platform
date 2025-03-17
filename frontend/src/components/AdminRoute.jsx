import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { getAccessToken } from "../utils/tokens";
import axiosInstance from "../utils/axios";

const AdminRoute = ({ element }) => {
  const [isAdmin, setIsAdmin] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const accessToken = getAccessToken();

  useEffect(() => {
    if (accessToken) {
      axiosInstance
        .get("/user/auth")
        .then((response) => {
          setIsAuthenticated(true);
          const userRole = response.data.role;
          setIsAdmin(userRole === "admin");
          localStorage.setItem("role", userRole);
        })
        .catch(() => {
          setIsAuthenticated(false);
          setIsAdmin(false);
        });
    } else {
      setIsAuthenticated(false);
      setIsAdmin(false);
    }
  }, [accessToken]);

  if (isAdmin === null || isAuthenticated === null) {
    return <p>Loading...</p>;
  }

  if (!isAuthenticated) {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("role");
    return <Navigate to="/login" replace />;
  }

  if (!isAdmin) {
    return <Navigate to="/" replace />;
  }

  return element;
};

export default AdminRoute;
