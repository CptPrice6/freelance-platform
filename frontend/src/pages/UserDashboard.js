import { useEffect, useState } from "react";
import axiosInstance from "../utils/axios"; // Import the configured axios instance

function UserDashboard() {
  const [user, setUser] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    axiosInstance
      .get("/user")
      .then((response) => {
        setUser(response.data);
      })
      .catch((err) => {
        setError(err.message);
      });
  }, []);

  return (
    <div>
      <h2>User Dashboard</h2>
      {error && <div style={{ color: "red" }}>{error}</div>}
      {user ? (
        <div>
          <p>
            <strong>Email:</strong> {user.email}
          </p>
        </div>
      ) : (
        <p>Loading user data...</p>
      )}
    </div>
  );
}

export default UserDashboard;
