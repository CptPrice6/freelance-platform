import React, { useEffect, useState } from "react";
import {
  Container,
  Card,
  Row,
  Col,
  Badge,
  Button,
  Modal,
  Tabs,
  Tab,
} from "react-bootstrap";
import { useParams } from "react-router-dom";
import axiosInstance from "../utils/axios";
import "../styles/ClientJobPage.css";
import moment from "moment";
import { Link } from "react-router-dom";

const ClientJobPage = () => {
  const { id } = useParams();
  const [job, setJob] = useState(null);
  const [allSkills, setAllSkills] = useState([]);
  const [activeTab, setActiveTab] = useState("job");
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedApp, setSelectedApp] = useState(null);
  const [showAppModal, setShowAppModal] = useState(false);
  const [rejectionReason, setRejectionReason] = useState("");
  const [newStatus, setNewStatus] = useState("");

  const [formData, setFormData] = useState({
    title: "",
    description: "",
    type: "ongoing",
    rate: "hourly",
    amount: 1,
    length: "<1",
    hours_per_week: "<20",
    skills: [],
  });

  useEffect(() => {
    axiosInstance
      .get(`/user/client/jobs/${id}`)
      .then((res) => {
        setJob(res.data);
        setFormData({
          ...res.data,
          skills: res.data.skills || [],
        });
      })
      .catch((err) => {
        console.error("Error fetching job:", err);
        setJob(null);
      });

    axiosInstance
      .get("/skills")
      .then((res) => {
        setAllSkills(res.data || []);
      })
      .catch((err) => {
        console.error("Error fetching skills:", err);
        setAllSkills([]);
      });
  }, [id]);

  const getJobStatusVariant = (status) => {
    switch (status) {
      case "open":
        return "primary";
      case "in-progress":
        return "warning";
      case "completed":
        return "success";
      default:
        return "secondary";
    }
  };

  const getApplicationStatusVariant = (status) => {
    switch (status) {
      case "pending":
        return "warning";
      case "accepted":
        return "success";
      case "rejected":
        return "danger";
      default:
        return "secondary";
    }
  };

  const handleAppClick = (app) => {
    setSelectedApp(app);
    setNewStatus(app.status);
    setRejectionReason(app.rejection_reason || "");
    setShowAppModal(true);
  };

  const handleCloseAppModal = () => {
    setSelectedApp(null);
    setShowAppModal(false);
  };

  const handleRespond = () => {
    axiosInstance
      .post(`/user/client/jobs/applications/${selectedApp.id}`, {
        status: newStatus,
        rejection_reason: newStatus === "rejected" ? rejectionReason : "",
      })
      .then(() => {
        const updated = job.applications.map((app) =>
          app.id === selectedApp.id
            ? { ...app, status: newStatus, rejection_reason: rejectionReason }
            : app
        );
        setJob({ ...job, applications: updated });
        handleCloseAppModal();

        axiosInstance
          .get(`/user/client/jobs/${id}`)
          .then((res) => {
            setJob(res.data);
            setFormData({
              ...res.data,
              skills: res.data.skills || [],
            });
          })
          .catch((err) => {
            console.error("Error fetching job:", err);
            setJob(null);
          });
      })

      .catch((err) => {
        console.error("Error updating application status:", err);
        alert("Failed to update application.");
      });
  };

  const handleDelete = () => {
    if (!window.confirm("Are you sure you want to delete this job?")) return;

    axiosInstance
      .delete(`/user/client/jobs/${id}`)
      .then(() => {
        window.location.href = "/client/jobs";
      })
      .catch((err) => {
        console.error("Error deleting job:", err);
        alert("Failed to delete job.");
      });
  };

  const handleComplete = () => {
    axiosInstance
      .post(`/user/client/jobs/${id}/complete`)
      .then(() => {
        axiosInstance
          .get(`/user/client/jobs/${id}`)
          .then((res) => setJob(res.data))
          .catch((err) => {
            console.error("Error refreshing job:", err);
          });
      })
      .catch((err) => {
        console.error("Error marking job complete:", err);
        alert("Failed to mark as complete.");
      });
  };

  const handleDownload = async (attachmentId) => {
    try {
      const response = await axiosInstance.get(
        `/user/attachments/${attachmentId}`,
        { responseType: "blob" }
      );

      const contentDisposition = response.headers["content-disposition"];
      let fileName = `attachment_${attachmentId}.pdf`;

      if (contentDisposition) {
        const fileNameMatch = contentDisposition.match(
          /filename[^;=\n]*=(['"]?)([^;"\n]+)\1/
        );
        if (fileNameMatch?.[2]) {
          fileName = fileNameMatch[2].trim();
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", fileName);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error("Download error:", error);
      alert("Failed to download file.");
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;

    if (name === "amount") {
      const numericValue = Math.max(1, Number(value));
      const maxLimit = formData.rate === "hourly" ? 1000 : 1000000;
      setFormData((prev) => ({
        ...prev,
        amount: Math.min(numericValue, maxLimit),
      }));
    } else if (name === "rate") {
      const maxLimit = value === "hourly" ? 1000 : 1000000;
      setFormData((prev) => ({
        ...prev,
        rate: value,
        amount: Math.min(prev.amount, maxLimit),
      }));
    } else {
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const handleSkillToggle = (id) => {
    setFormData((prev) => {
      const exists = prev.skills.some((skill) => skill.id === id);
      return {
        ...prev,
        skills: exists
          ? prev.skills.filter((skill) => skill.id !== id)
          : [...prev.skills, { id }],
      };
    });
  };

  const handleUpdateJob = () => {
    axiosInstance
      .put(`/user/client/jobs/${id}`, formData)
      .then(() => {
        axiosInstance
          .get(`/user/client/jobs/${id}`)
          .then((res) => {
            setJob(res.data);
            handleCloseEdit();
          })
          .catch((err) => {
            console.error("Error refreshing job after update:", err);
          });
      })
      .catch((err) => {
        console.error("Error updating job:", err);
        alert("Failed to update job.");
      });
  };

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
    return rate.charAt(0).toUpperCase() + rate.slice(1);
  };

  const formatAmount = (amount, rate) => {
    if (!amount) return "N/A";
    return rate === "hourly" ? `$${amount}/h` : `$${amount}`;
  };

  const handleEdit = () => setShowEditModal(true);
  const handleCloseEdit = () => setShowEditModal(false);

  if (!job) return <p>Loading...</p>;

  return (
    <Container className="py-5">
      <Tabs
        activeKey={activeTab}
        onSelect={(k) => setActiveTab(k)}
        className="mb-3"
      >
        <Tab eventKey="job" title="Job">
          <Container className="d-flex justify-content-center">
            <Card className="job-card shadow-lg p-4">
              <Card.Body>
                <h2 className="text-center">{job.title}</h2>
                <h5 className="text-center text-muted">
                  {formatProjectType(job.type)} Project
                </h5>

                {/* Posted By Section (Optional — You can skip or reuse if needed) */}
                <div className="mt-4 d-flex justify-content-center align-items-center">
                  <span className="me-3 text-muted">Status:</span>
                  <Badge
                    bg={getJobStatusVariant(job.status)}
                    className="px-3 py-2"
                  >
                    {job.status.toUpperCase()}
                  </Badge>
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
                      {formatAmount(job.amount, job.rate) ||
                        "No amount provided"}
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
                      {formatHoursPerWeek(job.hours_per_week) ||
                        "No availability"}
                    </p>
                  </Col>
                </Row>

                <div className="mt-4">
                  <h5>Skills:</h5>
                  {job.skills?.length > 0 ? (
                    <div className="skills-container">
                      {job.skills.map((skill) => (
                        <Badge
                          key={skill.id}
                          bg="primary"
                          className="skill-badge me-2"
                        >
                          {skill.name}
                        </Badge>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted">No skills listed</p>
                  )}
                </div>

                <div className="d-flex  justify-content-center mt-2">
                  <Link to={`/jobs/${job.id}`}>
                    <Button
                      variant="primary"
                      size="sm"
                      className="px-4 py-2 rounded-3 shadow-sm"
                    >
                      Your Public Page
                    </Button>
                  </Link>
                </div>

                {/* Action Buttons */}
                <div className="mt-3 d-flex justify-content-center gap-3">
                  {job.status === "open" && (
                    <>
                      <Button variant="warning" onClick={handleEdit}>
                        Edit
                      </Button>
                      <Button variant="danger" onClick={handleDelete}>
                        Delete
                      </Button>
                    </>
                  )}
                  {job.status === "in-progress" && (
                    <Button variant="success" onClick={handleComplete}>
                      Mark as Completed
                    </Button>
                  )}
                  {job.status === "completed" && (
                    <Button variant="secondary" disabled>
                      Completed
                    </Button>
                  )}
                </div>
              </Card.Body>
            </Card>
          </Container>
        </Tab>

        <Tab eventKey="applications" title="Applications">
          {job.applications.length === 0 ? (
            <p>No applications yet.</p>
          ) : (
            <Row>
              {job.applications.map((app) => (
                <Col xs={12} key={app.id} className="mb-3">
                  <Card
                    onClick={() => handleAppClick(app)}
                    className="freelancer-application-card p-3"
                    style={{ cursor: "pointer" }}
                  >
                    <Card.Body>
                      <div className="d-flex justify-content-between align-items-center">
                        <h5 className="application-title">
                          Applicant #{app.user_id}
                        </h5>
                        <Badge
                          bg={getApplicationStatusVariant(app.status)}
                          className="application-status-badge"
                        >
                          {app.status.toUpperCase()}
                        </Badge>
                      </div>
                      <p className="mb-1 text-muted">
                        <strong>Submitted at:</strong>{" "}
                        {moment(app.created_at).format("LLL")}
                      </p>
                    </Card.Body>
                  </Card>
                </Col>
              ))}
            </Row>
          )}

          {/* Application Modal */}
          {selectedApp && (
            <Modal show={showAppModal} onHide={handleCloseAppModal} size="lg">
              <Modal.Header closeButton>
                <Modal.Title>Application #{selectedApp.id}</Modal.Title>
              </Modal.Header>
              <Modal.Body>
                <div className="mb-2">
                  <strong>Submitted at:</strong>{" "}
                  {moment(selectedApp.created_at).format("LLL")}
                </div>

                <div className="mb-2">
                  <strong>Submitted by:</strong>{" "}
                  <a
                    href={`/freelancers/${selectedApp.user_id}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    Freelancer #{selectedApp.user_id}
                  </a>
                </div>

                <div className="mb-3">
                  <strong>Description:</strong>
                  <p className="mt-1">{selectedApp.description}</p>
                </div>

                <div className="mb-3">
                  <strong>Attachment:</strong>{" "}
                  {selectedApp.attachment?.file_name ? (
                    <div>
                      <Button
                        variant="link"
                        onClick={() =>
                          handleDownload(selectedApp.attachment.id)
                        }
                      >
                        📄 {selectedApp.attachment.file_name}
                      </Button>
                    </div>
                  ) : (
                    <p className="text-muted">No attachment</p>
                  )}
                </div>

                <div className="mb-4 text-center">
                  <strong>Status:</strong>
                  <div className="d-flex justify-content-center mt-2 gap-3">
                    {["rejected", "pending", "accepted"].map((status) => (
                      <Button
                        key={status}
                        variant={
                          newStatus === status
                            ? getApplicationStatusVariant(status)
                            : `outline-${getApplicationStatusVariant(status)}`
                        }
                        onClick={() =>
                          selectedApp.status === "pending" &&
                          setNewStatus(status)
                        }
                        disabled={selectedApp.status !== "pending"}
                      >
                        {status.charAt(0).toUpperCase() + status.slice(1)}
                      </Button>
                    ))}
                  </div>
                </div>

                {(selectedApp.status === "rejected" ||
                  newStatus === "rejected") && (
                  <div className="mb-3">
                    <strong>Rejection Reason:</strong>
                    <textarea
                      className="form-control"
                      rows={3}
                      value={rejectionReason}
                      onChange={(e) => setRejectionReason(e.target.value)}
                      placeholder="Explain why the application was rejected..."
                      disabled={selectedApp.status !== "pending"}
                    />
                  </div>
                )}
              </Modal.Body>
              <Modal.Footer>
                <Button variant="secondary" onClick={handleCloseAppModal}>
                  Close
                </Button>
                {selectedApp.status === "pending" && (
                  <Button variant="success" onClick={handleRespond}>
                    Respond
                  </Button>
                )}
              </Modal.Footer>
            </Modal>
          )}
        </Tab>
      </Tabs>

      {/* Edit Modal */}
      <Modal show={showEditModal} onHide={handleCloseEdit} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Edit Job</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {/* Reuse the form UI here */}
          <div className="mb-3">
            <label className="form-label">Title</label>
            <input
              type="text"
              className="form-control"
              name="title"
              maxLength={30}
              value={formData.title}
              onChange={handleInputChange}
            />
          </div>
          <div className="mb-3">
            <label className="form-label">Description</label>
            <textarea
              className="form-control"
              name="description"
              rows={4}
              value={formData.description}
              onChange={handleInputChange}
            />
          </div>
          <Row>
            <Col md={6}>
              <label>Project Type</label>
              <select
                className="form-select"
                name="type"
                value={formData.type}
                onChange={handleInputChange}
              >
                <option value="ongoing">Ongoing</option>
                <option value="one-time">One-Time</option>
              </select>
            </Col>
            <Col md={6}>
              <label>Rate</label>
              <select
                className="form-select"
                name="rate"
                value={formData.rate}
                onChange={handleInputChange}
              >
                <option value="hourly">Hourly</option>
                <option value="fixed">Fixed</option>
              </select>
            </Col>
            <Col md={6}>
              <label>Amount</label>
              <input
                type="number"
                className="form-control"
                name="amount"
                value={formData.amount}
                onChange={handleInputChange}
              />
            </Col>
            <Col md={6}>
              <label>Length</label>
              <select
                className="form-select"
                name="length"
                value={formData.length}
                onChange={handleInputChange}
              >
                <option value="<1">Less than 1 month</option>
                <option value="1-3">1–3 months</option>
                <option value="3-6">3–6 months</option>
                <option value="6-12">6–12 months</option>
                <option value="12+">12+ months</option>
              </select>
            </Col>
            <Col md={6}>
              <label>Availability</label>
              <select
                className="form-select"
                name="hours_per_week"
                value={formData.hours_per_week}
                onChange={handleInputChange}
              >
                <option value="<20">Less than 20 hrs/week</option>
                <option value="20-40">20–40 hrs/week</option>
                <option value="40-60">40–60 hrs/week</option>
                <option value="60-80">60–80 hrs/week</option>
                <option value="80+">80+ hrs/week</option>
              </select>
            </Col>
          </Row>
          <div className="mt-3">
            <label>Skills</label>
            <div className="d-flex flex-wrap gap-2">
              {allSkills.map((skill) => (
                <Button
                  key={skill.id}
                  size="sm"
                  variant={
                    formData.skills.some((s) => s.id === skill.id)
                      ? "success"
                      : "outline-secondary"
                  }
                  onClick={() => handleSkillToggle(skill.id)}
                >
                  {skill.name}
                </Button>
              ))}
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleCloseEdit}>
            Cancel
          </Button>
          <Button variant="success" onClick={handleUpdateJob}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
};

export default ClientJobPage;
