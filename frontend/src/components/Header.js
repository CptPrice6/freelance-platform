import { Link, useNavigate } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";
import axiosInstance from "../utils/axios";
import { getAccessToken } from "../utils/tokens";

function Header() {
  const navigate = useNavigate();
  const accessToken = getAccessToken();
  const isAuthenticated = !!accessToken;
  const userRole = localStorage.getItem("role");

  const handleLogout = async () => {
    try {
      if (isAuthenticated) {
        await axiosInstance.post("/user/logout");
      }
    } catch (error) {
      console.error("Logout failed:", error);
    } finally {
      // Clear all authentication-related data
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("role");
      navigate("/");
    }
  };
  const getDashboardRoute = () => {
    switch (userRole) {
      case "admin":
        return "/admin/dashboard";
      case "client":
        return "/client/dashboard";
      case "contractor":
        return "/contractor/dashboard";
      default:
        return "/contractor/dashboard";
    }
  };

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
      <div className="container">
        <Link className="navbar-brand" to="/">
          CodeHire
        </Link>
        <div className="collapse navbar-collapse">
          <ul className="navbar-nav ms-auto">
            {!isAuthenticated ? (
              <>
                <li className="nav-item">
                  <Link className="nav-link" to="/login">
                    Login
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" to="/register">
                    Register
                  </Link>
                </li>
              </>
            ) : (
              <>
                <li className="nav-item">
                  <Link className="nav-link" to={getDashboardRoute()}>
                    Dashboard
                  </Link>
                </li>
                <li className="nav-item">
                  <button
                    className="nav-link btn btn-link"
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </li>
              </>
            )}
          </ul>
        </div>
      </div>
    </nav>
  );
}

export default Header;
