import { useState, useEffect, useRef } from "react";
import React from "react";
import axiosInstance from "../utils/axios";
import { Button, Card, Container, Row, Col, Alert } from "react-bootstrap";
import "../styles/Dashboard.css";
import { Link } from "react-router-dom";

const ClientDashboard = () => {
  const [editingField, setEditingField] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    surname: "",
    description: "",
    company_name: "",
    industry: "",
    location: "",
  });
  const [updateError, setUpdateError] = useState(null);
  const inputRef = useRef(null);
  const [id, setId] = useState(null);

  const fetchClient = () => {
    axiosInstance.get("/user").then((res) => {
      const userData = res.data;
      setId(userData.id);
      setFormData({
        name: userData.name,
        surname: userData.surname,
        description: userData.client_data?.description || "",
        company_name: userData.client_data?.company_name || "",
        industry: userData.client_data?.industry || "",
        location: userData.client_data?.location || "",
      });
    });
  };

  useEffect(() => {
    fetchClient();
  }, []);

  useEffect(() => {
    if (inputRef.current && editingField === "description") {
      inputRef.current.style.height = "auto";
      inputRef.current.style.height = `${inputRef.current.scrollHeight}px`;
    }
  }, [editingField, formData.description]);

  const handleSave = (field) => {
    let value = formData[field];
    if (field === "description") {
      value = value.trim();
    }

    const endpoint =
      field === "name" || field === "surname" ? "/user" : "/user/client";

    axiosInstance
      .put(endpoint, { [field]: value })
      .then(() => {
        fetchClient();
        setUpdateError(null);
        setEditingField(null);
      })
      .catch((error) => {
        fetchClient();
        setUpdateError(error.response?.data?.error || "An error occurred");
      });
  };

  return (
    <Container
      fluid
      className="d-flex flex-column align-items-center p-5"
      style={{ backgroundColor: "#f8f9fa" }}
    >
      <h2 className="mb-4 text-center text-primary">Client Dashboard</h2>

      {updateError && (
        <Alert
          variant="danger"
          onClose={() => setUpdateError(null)}
          dismissible
        >
          {updateError}
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
        {[
          { field: "name", label: "Name" },
          { field: "surname", label: "Surname" },
          { field: "company_name", label: "Company Name" },
          { field: "industry", label: "Industry" },
          { field: "location", label: "Location" },
        ].map(({ field, label }) => (
          <Row className="mb-4" key={field}>
            <Col md={4} className="d-flex align-items-center">
              <label className="fw-bold text-dark">{label}:</label>
            </Col>
            <Col md={8} className="d-flex align-items-center">
              {editingField === field ? (
                <input
                  ref={inputRef}
                  type="text"
                  className="form-control me-2"
                  value={formData[field]}
                  onChange={(e) =>
                    setFormData({ ...formData, [field]: e.target.value })
                  }
                  onBlur={() => handleSave(field)}
                  onKeyDown={(e) => e.key === "Enter" && handleSave(field)}
                  autoFocus
                />
              ) : (
                <div
                  className="me-2 flex-grow-1 d-inline-block field-box"
                  onClick={() => setEditingField(field)}
                >
                  {formData[field] || "Not Set"}
                </div>
              )}
              <Button
                variant="outline-primary btn-sm"
                onClick={() => setEditingField(field)}
                style={{ marginLeft: "10px" }}
              >
                ✏️
              </Button>
            </Col>
          </Row>
        ))}

        {/* Description Field - Dynamic Resizing */}
        <Row className="mb-4">
          <Col md={4} className="d-flex align-items-center">
            <label className="fw-bold text-dark">Description:</label>
          </Col>
          <Col md={8} className="d-flex align-items-center">
            {editingField === "description" ? (
              <textarea
                ref={inputRef}
                className="form-control"
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                onBlur={() => handleSave("description")}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    handleSave("description");
                  }
                }}
                autoFocus
                rows={4}
              />
            ) : (
              <div
                className="me-2 flex-grow-1 field-box"
                onClick={() => setEditingField("description")}
              >
                {formData.description.split("\n").map((line, index) => (
                  <React.Fragment key={index}>
                    {line}
                    <br />
                  </React.Fragment>
                ))}
              </div>
            )}
            <Button
              variant="outline-primary btn-sm"
              onClick={() => setEditingField("description")}
              style={{ marginLeft: "10px" }}
            >
              ✏️
            </Button>
          </Col>
        </Row>

        <div className="d-flex justify-content-center mt-4">
          <Link to={`/clients/${id}`} style={{ textDecoration: "none" }}>
            <Button
              variant="primary"
              size="lg"
              className="px-5 py-3 rounded-3 shadow-sm"
              style={{
                backgroundColor: "#007bff",
                color: "#fff",
                fontWeight: "bold",
                transition: "background-color 0.3s ease",
              }}
            >
              Your Public Page
            </Button>
          </Link>
        </div>
      </Card>
    </Container>
  );
};

export default ClientDashboard;
