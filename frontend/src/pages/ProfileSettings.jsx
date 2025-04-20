import { useEffect, useState, useRef } from "react";
import axiosInstance from "../utils/axios";
import {
  Container,
  Card,
  Row,
  Col,
  Form,
  Button,
  Alert,
} from "react-bootstrap";
import "../styles/Dashboard.css";

function ProfileSettings() {
  const [user, setUser] = useState(null);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [updateError, setUpdateError] = useState("");
  const [updateSuccess, setUpdateSuccess] = useState("");
  const [editingEmail, setEditingEmail] = useState(false);
  const emailRef = useRef(null);

  useEffect(() => {
    fetchSettings();
  }, []);

  const fetchSettings = () => {
    axiosInstance
      .get("/user")
      .then((res) => {
        setUser(res.data);
        setEmail(res.data.email);
      })
      .catch((err) => {
        setUpdateError(
          err.response?.data?.error || "Failed to load user data."
        );
      });
  };

  const handleEmailSave = () => {
    if (email === user.email) {
      setEditingEmail(false);
      return;
    }

    axiosInstance
      .put("/user", { email })
      .then(() => {
        fetchSettings();
        setUpdateError("");
        setUser((prev) => ({ ...prev, email }));
      })
      .catch((err) => {
        fetchSettings();
        setUpdateError(err.response?.data?.error || "An error occurred.");
        setUpdateSuccess("");
      })
      .finally(() => setEditingEmail(false));
  };

  const handlePasswordSubmit = (e) => {
    e.preventDefault();
    if (!password || !newPassword)
      return setUpdateError("Both fields are required.");

    axiosInstance
      .put("/user", { password, new_password: newPassword })
      .then(() => {
        fetchSettings();
        setUpdateSuccess("Password updated successfully.");
        setUpdateError("");
        setPassword("");
        setNewPassword("");
      })
      .catch((err) => {
        fetchSettings();
        setUpdateError(err.response?.data?.error || "An error occurred.");
        setUpdateSuccess("");
      });
  };

  const handleDeleteAccount = () => {
    if (window.confirm("Are you sure you want to delete your account?")) {
      axiosInstance
        .delete("/user")
        .then(() => {
          localStorage.clear();
          window.location.replace("/");
        })
        .catch((err) => {
          setUpdateError(
            err.response?.data?.error || "Account deletion failed."
          );
        });
    }
  };

  return (
    <Container
      fluid
      className="d-flex flex-column align-items-center p-5"
      style={{ backgroundColor: "#f8f9fa" }}
    >
      <h2 className="mb-4 text-center text-primary">Profile Settings</h2>

      {(updateError || updateSuccess) && (
        <Alert
          variant={updateError ? "danger" : "success"}
          onClose={() => {
            setUpdateError("");
            setUpdateSuccess("");
          }}
          dismissible
        >
          {updateError || updateSuccess}
        </Alert>
      )}

      <Card
        className="p-4 shadow-lg w-100"
        style={{
          maxWidth: "800px",
          borderRadius: "12px",
          backgroundColor: "#fff",
          overflow: "hidden",
        }}
      >
        {user && (
          <>
            {/* Email Row */}
            <Row className="align-items-center mb-4">
              <Col md={4} className="d-flex align-items-center">
                <label className="fw-bold text-dark">Email:</label>
              </Col>
              <Col md={8} className="d-flex align-items-center">
                {editingEmail ? (
                  <Form.Control
                    type="email"
                    className="form-control"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    onBlur={handleEmailSave}
                    onKeyDown={(e) => e.key === "Enter" && handleEmailSave()}
                    ref={emailRef}
                    autoFocus
                    required
                  />
                ) : (
                  <div
                    className="me-2 flex-grow-1 field-box"
                    onClick={() => setEditingEmail(true)}
                  >
                    {email}
                  </div>
                )}
                <Button
                  variant="outline-primary btn-sm"
                  onClick={() => setEditingEmail(true)}
                  style={{ marginLeft: "10px" }}
                >
                  ✏️
                </Button>
              </Col>
            </Row>

            {/* Password Update */}
            <h5 className="section-heading mt-4 mb-3">Update Password</h5>
            <Form onSubmit={handlePasswordSubmit}>
              <Row className="mb-4 align-items-center">
                <Col md={4}>
                  <label className="fw-bold text-dark">Old Password:</label>
                </Col>
                <Col md={8}>
                  <Form.Control
                    type="password"
                    className="form-control"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  />
                </Col>
              </Row>

              <Row className="mb-4 align-items-center">
                <Col md={4}>
                  <label className="fw-bold text-dark">New Password:</label>
                </Col>
                <Col md={8}>
                  <Form.Control
                    type="password"
                    className="form-control"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    required
                  />
                </Col>
              </Row>

              <Row className="mb-4">
                <Col md={{ span: 8, offset: 4 }}>
                  <Button variant="primary" type="submit">
                    Update
                  </Button>
                </Col>
              </Row>
            </Form>

            <hr />

            {/* Delete Account */}
            <Row className="mt-3">
              <Col md={4}>
                <Button variant="outline-danger" onClick={handleDeleteAccount}>
                  Delete Account
                </Button>
              </Col>
            </Row>
          </>
        )}
      </Card>
    </Container>
  );
}

export default ProfileSettings;
