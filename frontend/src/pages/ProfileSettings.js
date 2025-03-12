import { useEffect, useState } from "react";
import axiosInstance from "../utils/axios";

function ProfileSettings() {
  const [user, setUser] = useState(null);
  const [error, setError] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [updateError, setUpdateError] = useState("");
  const [updateSuccess, setUpdateSuccess] = useState("");

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = () => {
    axiosInstance
      .get("/user")
      .then((response) => {
        setUser(response.data);
        setEmail(response.data.email);
      })
      .catch((err) => {
        const errorMessage = err.response?.data?.error || err.message;
        setError(errorMessage);
      });
  };

  // Handle email update
  const handleEmailSubmit = (e) => {
    e.preventDefault();

    if (!email || email === user.email) {
      setUpdateSuccess("");
      setUpdateError("No changes were made to the email.");
      return;
    }

    const updatedData = { email };

    axiosInstance
      .put("/user", updatedData)
      .then(() => {
        setUser((prevUser) => ({
          ...prevUser,
          email: email,
        }));
        setUpdateSuccess("Email updated successfully.");
        setUpdateError("");
      })
      .catch((err) => {
        fetchSettings();
        setUpdateError(
          err.response?.data?.error || "An unknown error occurred."
        );
        setUpdateSuccess("");
      });
  };

  const handlePasswordSubmit = (e) => {
    e.preventDefault();

    if (!password || !newPassword) {
      setUpdateSuccess("");
      setUpdateError("Please provide both old and new passwords.");
      return;
    }

    const updatedData = { password, new_password: newPassword };

    axiosInstance
      .put("/user", updatedData)
      .then(() => {
        setUpdateSuccess("Password updated successfully.");
        setUpdateError("");
        setPassword("");
        setNewPassword("");
      })
      .catch((err) => {
        setUpdateError(
          err.response?.data?.error || "An unknown error occurred."
        );
        setUpdateSuccess("");
      });
  };

  const handleDeleteAccount = () => {
    if (window.confirm("Are you sure you want to delete your account?")) {
      axiosInstance
        .delete("/user")
        .then(() => {
          localStorage.removeItem("accessToken");
          localStorage.removeItem("refreshToken");
          localStorage.removeItem("role");

          window.location.replace("/");

          setUser(null);

          alert("Your account has been deleted.");
        })
        .catch((err) => {
          setError(
            err.response?.data?.error ||
              "An unknown error occurred during deletion."
          );
        });
    }
  };

  return (
    <div className="container mt-5">
      <h2>Profile Settings</h2>
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
                  value={email}
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
                  value={password}
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
                  value={newPassword}
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

          {/* Delete Account Section */}
          <section className="mt-5">
            <button className="btn btn-danger" onClick={handleDeleteAccount}>
              Delete Account
            </button>
          </section>
        </div>
      ) : (
        <p>Loading user data...</p>
      )}
    </div>
  );
}

export default ProfileSettings;
