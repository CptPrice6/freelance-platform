import { useState, useEffect, useRef } from "react";
import React from "react";
import axiosInstance from "../utils/axios";
import {
  Button,
  Form,
  Card,
  Container,
  Row,
  Col,
  Alert,
} from "react-bootstrap";
import "../styles/Dashboard.css";
import { Link } from "react-router-dom";

const FreelancerDashboard = () => {
  const [editingField, setEditingField] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    surname: "",
    title: "",
    description: "",
    hourly_rate: 0,
    hours_per_week: "",
    skills: [],
  });
  const [allSkills, setAllSkills] = useState([]);
  const [newSkill, setNewSkill] = useState(null);
  const [updateError, setUpdateError] = useState(null);
  const inputRef = useRef(null);
  const [id, setId] = useState(null);

  const fetchFreelancer = () => {
    axiosInstance.get("/user").then((res) => {
      const userData = res.data;
      setId(userData.id);
      setFormData({
        name: userData.name,
        surname: userData.surname,
        title: userData.freelancer_data?.title || "",
        description: userData.freelancer_data?.description || "",
        hourly_rate: userData.freelancer_data?.hourly_rate || 0,
        hours_per_week: userData.freelancer_data?.hours_per_week || "",
        skills: userData.freelancer_data?.skills || [],
      });
    });
    axiosInstance.get("/skills").then((res) => setAllSkills(res.data));
  };

  useEffect(() => {
    fetchFreelancer();
  }, []);

  useEffect(() => {
    if (inputRef.current && editingField === "description") {
      inputRef.current.style.height = "auto";
      inputRef.current.style.height = `${inputRef.current.scrollHeight}px`;
    }
  }, [editingField, formData.description]);

  const handleSave = (field) => {
    let value = formData[field];
    if (field === "hourly_rate") {
      value = parseInt(value, 10);
    }

    const endpoint =
      field === "name" || field === "surname" ? "/user" : "/user/freelancer";

    axiosInstance
      .put(endpoint, { [field]: value })
      .then(() => {
        fetchFreelancer();
        setUpdateError(null);
        setEditingField(null);
      })
      .catch((error) => {
        fetchFreelancer();
        setUpdateError(error.response?.data?.error || "An error occurred");
      });
  };

  const addSkill = () => {
    if (newSkill === null) return;

    axiosInstance
      .post("/user/freelancer/skills", { skill_id: parseInt(newSkill) })
      .then(() => {
        fetchFreelancer();
        setUpdateError(null);
        setNewSkill(null);
      })
      .catch((error) => {
        setUpdateError(error.response?.data?.error || "An error occurred");
      });
  };

  const removeSkill = (skillId) => {
    axiosInstance
      .delete(`/user/freelancer/skills`, { data: { skill_id: skillId } })
      .then(() => {
        fetchFreelancer();
        setUpdateError(null);
      })
      .catch((error) => {
        setUpdateError(error.response?.data?.error || "An error occurred");
      });
  };

  const filteredSkills = allSkills.filter(
    (skill) => !formData.skills.some((userSkill) => userSkill.id === skill.id)
  );

  return (
    <Container
      fluid
      className="d-flex flex-column align-items-center p-5"
      style={{ backgroundColor: "#f8f9fa" }}
    >
      <h2 className="mb-4 text-center text-primary">Freelancer Dashboard</h2>

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
          { field: "title", label: "Title" },
          { field: "hourly_rate", label: "Hourly Rate" },
        ].map(({ field, label }) => (
          <Row className="mb-4" key={field}>
            <Col md={4} className="d-flex align-items-center">
              <label className="fw-bold text-dark">{label}:</label>
            </Col>
            <Col md={8} className="d-flex align-items-center">
              {editingField === field ? (
                <input
                  ref={inputRef}
                  type={field === "hourly_rate" ? "number" : "text"}
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
                  className="me-2 flex-grow-1 field-box"
                  onClick={() => setEditingField(field)}
                >
                  {/* Format hourly_rate  */}
                  {field === "hourly_rate"
                    ? formData[field]
                      ? `${formData[field]}$/h`
                      : "Not Set"
                    : formData[field] || "Not Set"}{" "}
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

        <Row className="mb-4">
          <Col md={4} className="d-flex align-items-center">
            <label className="fw-bold text-dark">Availability:</label>
          </Col>
          <Col md={8}>
            <div>
              <Form.Select
                value={formData.hours_per_week}
                onChange={(e) =>
                  setFormData({ ...formData, hours_per_week: e.target.value })
                }
                onBlur={() => handleSave("hours_per_week")}
                style={{
                  backgroundColor: "#e9f7fd",
                  borderRadius: "10px",
                  border: "1px solid #ddd",
                  fontWeight: 600,
                }}
              >
                <option value="<20">Less than 20 hours/week</option>
                <option value="20-40">20-40 hours/week</option>
                <option value="40-60">40-60 hours/week</option>
                <option value="60-80">60-80 hours/week</option>
                <option value="80+">More than 80 hours/week</option>
              </Form.Select>
            </div>
          </Col>
        </Row>

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
                  // If Enter is pressed without Shift, save the description
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    handleSave("description");
                  }
                }}
                autoFocus
                rows={4}
                style={{
                  resize: "none",
                  maxWidth: "100%",
                  wordWrap: "break-word",
                  overflowWrap: "break-word",
                  wordBreak: "break-all",
                }}
              />
            ) : (
              // Convert newlines to <br /> for display
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

        <div className="mt-4">
          <h5 className="text-center mb-3 skill-text">Skills</h5>
          <ul className="list-group mb-4">
            {formData.skills.map((skill) => (
              <li
                key={skill.id}
                className="list-group-item d-flex justify-content-between align-items-center skill-item"
              >
                <span className="skill-name">{skill.name}</span>
                <Button
                  variant="danger"
                  size="sm"
                  onClick={() => removeSkill(skill.id)}
                  className="delete-skill-btn"
                >
                  ❌
                </Button>
              </li>
            ))}
          </ul>

          <div className="d-flex justify-content-between align-items-center mt-4">
            {/* Skill Select Dropdown */}
            <Form.Select
              className="skill-select me-3"
              value={newSkill || ""}
              onChange={(e) => setNewSkill(e.target.value)}
              aria-label="Select a skill"
            >
              <option value="">Select a skill</option>
              {filteredSkills.map((skill) => (
                <option key={skill.id} value={skill.id}>
                  {skill.name}
                </option>
              ))}
            </Form.Select>

            {/* Add Skill Button */}
            <Button
              variant="success"
              onClick={addSkill}
              disabled={newSkill === null}
              className="add-skill-btn d-flex align-items-center justify-content-center"
            >
              <span className="me-2">➕</span> Add Skill
            </Button>
          </div>
        </div>

        <div className="d-flex justify-content-center mt-4">
          <Link to={`/freelancers/${id}`} style={{ textDecoration: "none" }}>
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

export default FreelancerDashboard;
