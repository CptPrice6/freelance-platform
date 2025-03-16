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

const FreelancerDashboard = () => {
  const [editingField, setEditingField] = useState(null);
  const [formData, setFormData] = useState({
    name: "",
    surname: "",
    title: "",
    description: "",
    hourly_rate: 0,
    work_type: "remote",
    hours_per_week: 0,
    skills: [],
  });
  const [allSkills, setAllSkills] = useState([]);
  const [newSkill, setNewSkill] = useState(null);
  const [updateError, setUpdateError] = useState(null);
  const inputRef = useRef(null);

  // Function to fetch user and skills data
  const fetchFreelancer = () => {
    axiosInstance.get("/user").then((res) => {
      const userData = res.data;
      setFormData({
        name: userData.name,
        surname: userData.surname,
        title: userData.freelancer_data.title || "",
        description: userData.freelancer_data.description || "",
        hourly_rate: userData.freelancer_data.hourly_rate || 0,
        work_type: userData.freelancer_data.work_type || "remote",
        hours_per_week: userData.freelancer_data.hours_per_week || 0,
        skills: userData.freelancer_data.skills || [],
      });
    });
    axiosInstance.get("/skills").then((res) => setAllSkills(res.data));
  };

  // Fetch user data on initial render
  useEffect(() => {
    fetchFreelancer();
  }, []);

  // Adjust textarea height dynamically
  useEffect(() => {
    if (inputRef.current && editingField === "description") {
      inputRef.current.style.height = "auto"; // Reset height for auto resizing
      inputRef.current.style.height = `${inputRef.current.scrollHeight}px`; // Adjust height based on content
    }
  }, [editingField, formData.description]);

  const handleSave = (field) => {
    let value = formData[field];
    if (field === "hourly_rate" || field === "hours_per_week") {
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
        fetchFreelancer();
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
        fetchFreelancer();
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
          overflow: "hidden", // Prevent card content from overflowing
        }}
      >
        {[
          { field: "name", label: "Name" },
          { field: "surname", label: "Surname" },
          { field: "title", label: "Title" },
          { field: "hourly_rate", label: "Hourly Rate" },
          { field: "hours_per_week", label: "Hours per Week" },
        ].map(({ field, label }) => (
          <Row className="mb-4" key={field}>
            <Col md={4} className="d-flex align-items-center">
              <label className="fw-bold text-dark">{label}:</label>
            </Col>
            <Col md={8} className="d-flex align-items-center">
              {editingField === field ? (
                <input
                  ref={inputRef}
                  type={
                    field === "hourly_rate" || field === "hours_per_week"
                      ? "number"
                      : "text"
                  }
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
                  {formData[field]}
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
            <label className="fw-bold text-dark">Work Type:</label>
          </Col>
          <Col md={8}>
            <div>
              <Form.Select
                value={formData.work_type}
                onChange={(e) =>
                  setFormData({ ...formData, work_type: e.target.value })
                }
                onBlur={() => handleSave("work_type")}
                style={{
                  backgroundColor: "#e9f7fd",
                  borderRadius: "10px",
                  border: "1px solid #ddd",
                  fontWeight: 500,
                }}
              >
                <option value="on-site">On-Site</option>
                <option value="remote">Remote</option>
                <option value="hybrid">Hybrid</option>
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
                    e.preventDefault(); // Prevent the newline
                    handleSave("description"); // Save the description
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
          <h5 className="text-center mb-3 text-primary">Skills</h5>
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
      </Card>
    </Container>
  );
};

export default FreelancerDashboard;
