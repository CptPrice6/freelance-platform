import { Link, useNavigate } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";

function Header() {
  const navigate = useNavigate();
  const accessToken = localStorage.getItem("accessToken");
  const isAuthenticated = !!accessToken;

  const handleLogout = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    navigate("/"); // Redirect to home page after logout
  };

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark">
      <div className="container">
        <Link className="navbar-brand" to="/">
          Website
        </Link>
        <div className="collapse navbar-collapse">
          <ul className="navbar-nav ms-auto">
            {/* Show Login/Register if not authenticated */}
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
                {/* Show Dashboard and Logout if authenticated */}
                <li className="nav-item">
                  <Link className="nav-link" to="/user/dashboard">
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
