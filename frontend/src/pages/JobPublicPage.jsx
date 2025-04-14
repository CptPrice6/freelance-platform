import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import {
  Card,
  Container,
  Row,
  Col,
  Badge,
  Button,
  Modal,
} from "react-bootstrap";
import axiosInstance from "../utils/axios";
import "../styles/JobPublicPage.css";

const JobPublicPage = () => {
  const { id } = useParams();
  const [job, setJob] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isFreelancer, setIsFreelancer] = useState(false);
  const [isApplied, setIsApplied] = useState(false);
  const [proposalDescription, setProposalDescription] = useState("");
  const [file, setFile] = useState(null);
  const [fileBase64, setFileBase64] = useState("");

  useEffect(() => {
    const role = localStorage.getItem("role");
    if (role === "admin") setIsAdmin(true);
    if (role === "freelancer") setIsFreelancer(true);

    axiosInstance
      .get(`/jobs/${id}`)
      .then((res) => {
        setJob(res.data);
        if (role === "freelancer" && res.data.application_id !== 0) {
          setIsApplied(true);
        }
      })
      .catch((err) => {
        if (err.response && err.response.status === 404) {
          setJob("not-found");
        } else {
          console.error("Failed to fetch job:", err);
        }
      });
  }, [id]);

  if (job === null) return <p className="text-center mt-5">Loading...</p>;
  if (job === "not-found") {
    return <p className="text-center mt-5 text-danger">Job not found.</p>;
  }

  const handleApplyClick = () => {
    // Show the modal for applying to the job
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
  };

  const handleApplyJob = () => {
    if (!proposalDescription) {
      alert("Please write a proposal before submitting.");
      return;
    }

    const payload = {
      job_id: job.id,
      description: proposalDescription,
      file_name: file ? file.name : "",
      file_base64: fileBase64,
    };

    axiosInstance
      .post(`/user/freelancer/applications`, payload, {
        headers: { "Content-Type": "application/json" },
      })
      .then(() => {
        setIsApplied(true);
        setShowModal(false);
        alert("Application submitted successfully!");
        setTimeout(() => {
          window.location.href = "/freelancer/applications";
        }, 1000); // waits 1 second before redirecting
      })
      .catch((err) => {
        console.error(err);
        alert("Failed to submit application.");
      });
  };

  const handleDeleteJob = () => {
    // Delete job request (for admin)
    axiosInstance
      .delete(`/admin/jobs/${id}`)
      .then(() => {
        alert("Job deleted successfully!");
        // Redirect to the job listings or some other page
        window.location.href = "/jobs";
      })
      .catch((err) => console.log(err));
  };

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];

    if (selectedFile && selectedFile.type !== "application/pdf") {
      alert("Only PDF files are allowed.");
      return;
    }

    setFile(selectedFile);

    const reader = new FileReader();
    reader.readAsDataURL(selectedFile);
    reader.onload = () => {
      const base64String = reader.result.split(",")[1];
      setFileBase64(base64String);
    };
  };

  const handleClientProfileClick = () => {
    // Redirect to the client's profile page
    window.location.href = `/clients/${job.client_id}`;
  };

  // Formatters for displaying job info
  const formatHoursPerWeek = (value) => {
    if (!value) return "N/A";
    if (value === "<20") return "Less than 20 hours/week";
    if (value === "80+") return "More than 80 hours/week";
    return `${value} hours/week`;
  };

  const formatProjectType = (type) => {
    if (!type) return "N/A";
    if (type === "ongoing") return "Ongoing";
    if (type === "one-time") return "One-Time";
    return type;
  };

  const formatProjectLength = (length) => {
    if (!length) return "N/A";
    if (length === "<1") return "Less than 1 month";
    if (length === "1-3") return "1–3 months";
    if (length === "3-6") return "3–6 months";
    if (length === "6-12") return "6–12 months";
    if (length === "12+") return "More than 12 months";
    return length;
  };

  const formatRate = (rate) => {
    if (!rate) return "N/A";
    return rate.charAt(0).toUpperCase() + rate.slice(1); // Capitalize the rate
  };

  const formatAmount = (amount, rate) => {
    if (!amount) return "N/A";
    return rate === "hourly" ? `$${amount}/h` : `$${amount}`; // Add '/h' if hourly
  };

  return (
    <Container className="py-5 d-flex justify-content-center">
      <Card className="public-job-card shadow-lg p-4">
        <Card.Body>
          <h2 className="text-center">{job.title}</h2>
          <h5 className="text-center text-muted">
            {formatProjectType(job.type) || "No Type Provided"} Project
          </h5>

          {/* Posted By Section */}
          <div className="mt-4 d-flex justify-content-center align-items-center">
            <span className="me-3 text-muted">Posted By:</span>
            <div
              className="client-avatar-circle"
              onClick={handleClientProfileClick}
              style={{ cursor: "pointer" }}
            >
              <div className="default-avatar">C</div>
            </div>
          </div>

          <Row className="mt-4">
            <Col md={12}>
              <p>
                <strong>Description:</strong>
              </p>
              <p className="text-muted">
                {job.description || "No description provided"}
              </p>
            </Col>
            <Col md={6}>
              <p>
                <strong>Rate:</strong>
              </p>
              <p className="text-muted">
                {formatRate(job.rate) || "No rate provided"}
              </p>
            </Col>
            <Col md={6}>
              <p>
                <strong>Amount:</strong>
              </p>
              <p className="text-muted">
                {formatAmount(job.amount, job.rate) || "No amount provided"}
              </p>
            </Col>
            <Col md={6}>
              <p>
                <strong>Project Length:</strong>
              </p>
              <p className="text-muted">
                {formatProjectLength(job.length) || "No length provided"}
              </p>
            </Col>
            <Col md={6}>
              <p>
                <strong>Availability:</strong>
              </p>
              <p className="text-muted">
                {formatHoursPerWeek(job.hours_per_week) || "No hours provided"}
              </p>
            </Col>
          </Row>

          <div className="mt-4">
            <h5>Skills:</h5>
            {job.skills && job.skills.length > 0 ? (
              <div className="skills-container">
                {job.skills.map((skill) => (
                  <Badge key={skill.id} bg="primary" className="skill-badge">
                    {skill.name}
                  </Badge>
                ))}
              </div>
            ) : (
              <p className="text-muted">No skills listed</p>
            )}
          </div>

          {/* Conditional buttons */}
          <div className="mt-4 d-flex justify-content-center">
            {isFreelancer && !isApplied && (
              <Button
                className="me-2"
                variant="success"
                onClick={handleApplyClick}
              >
                Apply
              </Button>
            )}
            {isFreelancer && isApplied && (
              <Button className="me-2" disabled>
                Already Applied
              </Button>
            )}
            {isAdmin && (
              <Button
                className="me-2"
                variant="danger"
                onClick={handleDeleteJob}
              >
                Delete
              </Button>
            )}
          </div>

          {/* Modal for applying to the job */}
          <Modal show={showModal} onHide={handleCloseModal}>
            <Modal.Header closeButton>
              <Modal.Title>Apply for Job</Modal.Title>
            </Modal.Header>
            <Modal.Body>
              <div className="mb-3">
                <label htmlFor="proposal" className="form-label">
                  Proposal
                </label>
                <textarea
                  id="proposal"
                  className="form-control"
                  rows={4}
                  placeholder="Describe your application..."
                  value={proposalDescription}
                  onChange={(e) => setProposalDescription(e.target.value)}
                />
              </div>
              <div className="mb-3">
                <label htmlFor="fileUpload" className="form-label">
                  Attach PDF (optional)
                </label>
                <input
                  id="fileUpload"
                  type="file"
                  className="form-control"
                  accept="application/pdf"
                  onChange={handleFileChange}
                />
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button variant="secondary" onClick={handleCloseModal}>
                Close
              </Button>
              <Button variant="success" onClick={handleApplyJob}>
                Apply
              </Button>
            </Modal.Footer>
          </Modal>
        </Card.Body>
      </Card>
    </Container>
  );
};

export default JobPublicPage;
