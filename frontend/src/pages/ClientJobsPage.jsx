import { useEffect, useState, useRef } from "react";
import axiosInstance from "../utils/axios";
import {
  Container,
  Card,
  Row,
  Col,
  Pagination,
  Badge,
  Modal,
} from "react-bootstrap";
import { Link } from "react-router-dom";
import "../styles/ClientJobsPage.css";

const ClientJobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [showModal, setShowModal] = useState(false);
  const [allSkills, setAllSkills] = useState([]);
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
  const titleRef = useRef(null);
  const descRef = useRef(null);
  const jobsPerPage = 10;

  useEffect(() => {
    axiosInstance
      .get("/user/client/jobs")
      .then((res) => {
        setJobs(res.data || []);
      })
      .catch((err) => {
        console.error("Error fetching client jobs:", err);
        setJobs([]);
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
  }, []);

  const indexOfLastJob = currentPage * jobsPerPage;
  const indexOfFirstJob = indexOfLastJob - jobsPerPage;
  const currentJobs = jobs.slice(indexOfFirstJob, indexOfLastJob);

  const getStatusVariant = (status) => {
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

  const handleCloseModal = () => setShowModal(false);
  const handleShowModal = () => setShowModal(true);

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

  const handleSubmitJob = () => {
    if (!formData.title.trim()) {
      titleRef.current.focus();
      titleRef.current.setCustomValidity("This field is required");
      titleRef.current.reportValidity();
      return;
    } else {
      titleRef.current.setCustomValidity(""); // clear message
    }

    if (!formData.description.trim()) {
      descRef.current.focus();
      descRef.current.setCustomValidity("This field is required");
      descRef.current.reportValidity();
      return;
    } else {
      descRef.current.setCustomValidity("");
    }

    axiosInstance.post("/user/client/jobs", formData).then(() => {
      handleCloseModal();
      axiosInstance.get("/user/client/jobs").then((res) => setJobs(res.data));
    });
  };
  return (
    <Container className="py-5">
      <div className="d-flex justify-content-end mb-4">
        <button className="btn btn-success" onClick={handleShowModal}>
          + Create Job
        </button>
      </div>

      <h2 className="text-center mb-4 job-section-title">Your Posted Jobs</h2>
      {currentJobs.length === 0 ? (
        <div className="text-center my-4">
          <p>You haven't posted any jobs yet.</p>
        </div>
      ) : (
        <>
          <Row>
            {currentJobs.map((job) => (
              <Col key={job.id} xs={12} className="mb-3">
                <Link
                  to={`/client/jobs/${job.id}`}
                  className="text-decoration-none"
                >
                  <Card className="client-job-card p-4 shadow-sm">
                    <Card.Body>
                      <div className="d-flex justify-content-between align-items-center mb-2">
                        <h5 className="fw-bold job-title">{job.title}</h5>
                        <Badge
                          bg={getStatusVariant(job.status)}
                          className="job-status-badge"
                        >
                          {job.status}
                        </Badge>
                      </div>
                      <p className="mb-1 text-muted">
                        <strong>ID:</strong> {job.id}
                      </p>
                      <p className="mb-0 text-muted">
                        <strong>Applications:</strong> {job.application_count}
                      </p>
                    </Card.Body>
                  </Card>
                </Link>
              </Col>
            ))}
          </Row>

          <Pagination className="justify-content-center mt-4 custom-pagination">
            {[...Array(Math.ceil(jobs.length / jobsPerPage)).keys()].map(
              (number) => (
                <Pagination.Item
                  key={number + 1}
                  active={number + 1 === currentPage}
                  onClick={() => setCurrentPage(number + 1)}
                >
                  {number + 1}
                </Pagination.Item>
              )
            )}
          </Pagination>
        </>
      )}

      <Modal show={showModal} onHide={handleCloseModal} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Create a New Job</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <div className="mb-3">
            <label className="form-label">Title (max 30 chars)</label>
            <input
              type="text"
              className="form-control"
              name="title"
              maxLength={30}
              required
              ref={titleRef}
              value={formData.title}
              onChange={handleInputChange}
            />
          </div>

          <div className="mb-3">
            <label className="form-label">Description</label>
            <textarea
              className="form-control"
              name="description"
              required
              ref={descRef}
              rows={4}
              value={formData.description}
              onChange={handleInputChange}
            />
          </div>

          <div className="row">
            <div className="col-md-6 mb-3">
              <label className="form-label">Project Type</label>
              <select
                className="form-select"
                name="type"
                value={formData.type}
                onChange={handleInputChange}
              >
                <option value="ongoing">Ongoing</option>
                <option value="one-time">One-Time</option>
              </select>
            </div>

            <div className="col-md-6 mb-3">
              <label className="form-label">Rate</label>
              <select
                className="form-select"
                name="rate"
                value={formData.rate}
                onChange={handleInputChange}
              >
                <option value="hourly">Hourly</option>
                <option value="fixed">Fixed</option>
              </select>
            </div>

            <div className="col-md-6 mb-3">
              <label className="form-label">
                Amount{" "}
                {formData.rate === "hourly"
                  ? "/h (max $1000)"
                  : "(max $1,000,000)"}
              </label>
              <input
                type="number"
                className="form-control"
                name="amount"
                min={1}
                max={formData.rate === "hourly" ? 1000 : 1000000}
                value={formData.amount}
                onChange={handleInputChange}
                required
              />
            </div>

            <div className="col-md-6 mb-3">
              <label className="form-label">Project Length</label>
              <select
                className="form-select"
                name="length"
                value={formData.length}
                onChange={handleInputChange}
              >
                {[
                  { value: "<1", label: "Less than 1 month" },
                  { value: "1-3", label: "1-3 months" },
                  { value: "3-6", label: "3-6 months" },
                  { value: "6-12", label: "6-12 months" },
                  { value: "12+", label: "More than 12 months" },
                ].map(({ value, label }) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>

            <div className="col-md-6 mb-3">
              <label className="form-label">Availability</label>
              <select
                className="form-select"
                name="hours_per_week"
                value={formData.hours_per_week}
                onChange={handleInputChange}
              >
                {[
                  { value: "<20", label: "Less than 20 hours/week" },
                  { value: "20-40", label: "20–40 hours/week" },
                  { value: "40-60", label: "40–60 hours/week" },
                  { value: "60-80", label: "60–80 hours/week" },
                  { value: "80+", label: "More than 80 hours/week" },
                ].map(({ value, label }) => (
                  <option key={value} value={value}>
                    {label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="mb-3">
            <label className="form-label">Skills</label>
            <div className="d-flex flex-wrap gap-2">
              {allSkills.length === 0 ? (
                <p>No skills available. Please try again later.</p> // Message when there are no skills
              ) : (
                allSkills.map((skill) => (
                  <button
                    type="button"
                    key={skill.id}
                    className={`btn btn-sm ${
                      formData.skills.some((s) => s.id === skill.id)
                        ? "btn btn-success"
                        : "btn btn-outline-secondary"
                    }`}
                    onClick={() => handleSkillToggle(skill.id)}
                  >
                    {skill.name}
                  </button>
                ))
              )}
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <button className="btn btn-secondary" onClick={handleCloseModal}>
            Cancel
          </button>
          <button className="btn btn-success" onClick={handleSubmitJob}>
            Create Job
          </button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
};

export default ClientJobsPage;
