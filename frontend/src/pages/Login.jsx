import { useState } from "react";
import { API_BASE_URL } from "../config";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";
import { setAccessToken, setRefreshToken } from "../utils/tokens";
import { jwtDecode } from "jwt-decode";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const response = await axios.post(`${API_BASE_URL}/login`, {
        email,
        password,
      });

      const accessToken = response.data.access_token;
      setAccessToken(accessToken);
      setRefreshToken(response.data.refresh_token);
      const decodedToken = jwtDecode(accessToken);
      const userRole = decodedToken.role || "freelancer";
      localStorage.setItem("role", userRole);
      alert("Login successful!");
      navigate("/");
    } catch (err) {
      setError(err.response?.data?.error || "Login failed");
    }
  };

  return (
    <div>
      <h2>Login</h2>
      {error && <div className="alert alert-danger">{error}</div>}
      <form onSubmit={handleLogin}>
        <div className="mb-3">
          <label className="form-label">Email</label>
          <input
            type="email"
            className="form-control"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="mb-3">
          <label className="form-label">Password</label>
          <input
            type="password"
            className="form-control"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="btn btn-primary">
          Login
        </button>
      </form>
      <p className="mt-3">
        Don't have an account? <Link to="/register">Sign up</Link>
      </p>
    </div>
  );
}

export default Login;
