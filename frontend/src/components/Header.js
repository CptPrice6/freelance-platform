import { Link, useNavigate } from "react-router-dom";
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
      case "freelancer":
        return "/freelancer/dashboard";
      default:
        return "/freelancer/dashboard";
    }
  };

  const renderRoleSpecificLinks = () => {
    switch (userRole) {
      case "freelancer":
        return (
          <>
            <li className="nav-item">
              <Link className="nav-link" to="/freelancer/applications">
                My Applications
              </Link>
            </li>
            <li className="nav-item">
              <Link className="nav-link" to="/freelancer/jobs">
                My Jobs
              </Link>
            </li>
          </>
        );
      case "client":
        return (
          <li className="nav-item">
            <Link className="nav-link" to="/client/jobs">
              My Posted Jobs
            </Link>
          </li>
        );
      default:
        return null;
    }
  };

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
      <div className="container">
        {/* Left side links */}
        <div className="navbar-nav">
          {isAuthenticated && (
            <>
              <li className="nav-item">
                <Link className="nav-link" to="/jobs">
                  Jobs
                </Link>
              </li>
              <li className="nav-item">
                <Link className="nav-link" to="/freelancers">
                  Freelancers
                </Link>
              </li>
            </>
          )}
        </div>

        {/* Right side links */}
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
                {renderRoleSpecificLinks()}
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
