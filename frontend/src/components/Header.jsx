import { Link, useNavigate } from "react-router-dom";
import { getAccessToken } from "../utils/tokens";
import { useState } from "react";

function Header() {
  const navigate = useNavigate();
  const accessToken = getAccessToken();
  const isAuthenticated = !!accessToken;
  const userRole = localStorage.getItem("role");
  const [isNavbarOpen, setIsNavbarOpen] = useState(false);

  const handleLogout = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("role");
    navigate("/");
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
        <Link className="navbar-brand" to="/">
          FreelancePlatform
        </Link>

        {/* Mobile Toggle Button */}
        <button
          className="navbar-toggler"
          type="button"
          onClick={() => setIsNavbarOpen(!isNavbarOpen)}
        >
          <span className="navbar-toggler-icon"></span>
        </button>

        {/* Navbar Links */}
        <div
          className={`collapse navbar-collapse ${isNavbarOpen ? "show" : ""}`}
        >
          <ul className="navbar-nav ms-auto">
            {/* Links for authenticated users */}
            {isAuthenticated ? (
              <>
                {/* General Links */}
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
                <li className="nav-item">
                  <Link className="nav-link" to="/clients">
                    Clients
                  </Link>
                </li>

                {/* Dashboard Link */}
                <li className="nav-item">
                  <Link className="nav-link" to={getDashboardRoute()}>
                    Dashboard
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" to="/settings">
                    Settings
                  </Link>
                </li>

                {/* Role-specific Links */}
                {renderRoleSpecificLinks()}

                {/* Logout Button */}
                <li className="nav-item">
                  <button
                    className="nav-link btn btn-link"
                    onClick={handleLogout}
                  >
                    Logout
                  </button>
                </li>
              </>
            ) : (
              <>
                {/* Links for unauthenticated users */}
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
            )}
          </ul>
        </div>
      </div>
    </nav>
  );
}

export default Header;
