import { useEffect, useState } from "react";
import axiosInstance from "../utils/axios"; // Import the configured axios instance

// Include the Bootstrap CDN in your public HTML or install via npm/yarn in your project
function UserDashboard() {
  const [user, setUser] = useState(null);
  const [error, setError] = useState(""); // General error
  const [email, setEmail] = useState(""); // State for the email input field
  const [password, setPassword] = useState(""); // State for the old password
  const [newPassword, setNewPassword] = useState(""); // State for the new password
  const [updateError, setUpdateError] = useState(""); // Error for update process
  const [updateSuccess, setUpdateSuccess] = useState(""); // Success message

  // Fetch user data when the component mounts
  useEffect(() => {
    axiosInstance
      .get("/user")
      .then((response) => {
        setUser(response.data);
        setEmail(response.data.email); // Pre-fill email input with current email
      })
      .catch((err) => {
        const errorMessage = err.response?.data?.error || err.message;
        setError(errorMessage); // Set error state with either the response's error or the generic message
      });
  }, []);

  // Handle email update
  const handleEmailSubmit = (e) => {
    e.preventDefault();

    if (!email || email === user.email) {
      setUpdateSuccess(""); // Clear any previous success messages
      setUpdateError("No changes were made to the email.");
      return;
    }

    const updatedData = { email };

    // Send PUT request to update email
    axiosInstance
      .put("/user", updatedData)
      .then(() => {
        setUser((prevUser) => ({
          ...prevUser,
          email: email, // Update email in the user state with the new value
        }));
        setUpdateSuccess("Email updated successfully.");
        setUpdateError(""); // Clear any previous errors
      })
      .catch((err) => {
        setUpdateError(
          err.response?.data?.error || "An unknown error occurred."
        );
        setUpdateSuccess(""); // Clear any previous success messages
      });
  };

  // Handle password update
  const handlePasswordSubmit = (e) => {
    e.preventDefault();

    if (!password || !newPassword) {
      setUpdateSuccess(""); // Clear any previous success messages
      setUpdateError("Please provide both old and new passwords.");
      return;
    }

    const updatedData = { password, new_password: newPassword };

    // Send PUT request to update password
    axiosInstance
      .put("/user", updatedData)
      .then(() => {
        setUpdateSuccess("Password updated successfully.");
        setUpdateError(""); // Clear any previous errors
        setPassword(""); // Clear password fields after success
        setNewPassword(""); // Clear password fields after success
      })
      .catch((err) => {
        setUpdateError(
          err.response?.data?.error || "An unknown error occurred."
        );
        setUpdateSuccess(""); // Clear any previous success messages
      });
  };

  return (
    <div className="container mt-5">
      <h2>User Dashboard</h2>
      {error && (
        <div className="alert alert-danger" role="alert">
          {error}
        </div>
      )}

      {user ? (
        <div>
          <p>
            <strong>Email:</strong> {user.email}
          </p>

          {/* Email Update Section */}
          <section>
            <h3>Update Email</h3>
            <form onSubmit={handleEmailSubmit}>
              <div className="mb-3">
                <label htmlFor="email" className="form-label">
                  Change Email:
                </label>
                <input
                  type="email"
                  id="email"
                  className="form-control"
                  value={email} // Controlled input value
                  onChange={(e) => setEmail(e.target.value)} // Update email state
                />
              </div>
              <button type="submit" className="btn btn-primary">
                Update Email
              </button>
            </form>
          </section>

          {/* Password Update Section */}
          <section className="mt-5">
            <h3>Update Password</h3>
            <form onSubmit={handlePasswordSubmit}>
              <div className="mb-3">
                <label htmlFor="password" className="form-label">
                  Old Password:
                </label>
                <input
                  type="password"
                  id="password"
                  className="form-control"
                  value={password} // Controlled input value
                  onChange={(e) => setPassword(e.target.value)} // Update old password state
                />
              </div>
              <div className="mb-3">
                <label htmlFor="new_password" className="form-label">
                  New Password:
                </label>
                <input
                  type="password"
                  id="new_password"
                  className="form-control"
                  value={newPassword} // Controlled input value
                  onChange={(e) => setNewPassword(e.target.value)} // Update new password state
                />
              </div>
              <button type="submit" className="btn btn-primary">
                Update Password
              </button>
            </form>
          </section>

          {/* Display Success/Error Messages */}
          {updateError && (
            <div className="alert alert-danger mt-3" role="alert">
              {updateError}
            </div>
          )}
          {updateSuccess && (
            <div className="alert alert-success mt-3" role="alert">
              {updateSuccess}
            </div>
          )}
        </div>
      ) : (
        <p>Loading user data...</p>
      )}
    </div>
  );
}

export default UserDashboard;
