import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";
import {
  Card,
  Container,
  Row,
  Col,
  Pagination,
  Form,
  DropdownButton,
  Alert,
} from "react-bootstrap";
import "../styles/JobsPage.css";
import { Link } from "react-router-dom";

const JobsPage = () => {
  const [jobs, setJobs] = useState(null);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [filters, setFilters] = useState({
    type: [],
    length: [],
    rate: [],
    availability: [],
  });
  const jobsPerPage = 9;

  useEffect(() => {
    const fetchJobs = async () => {
      try {
        const res = await axiosInstance.get("/jobs");
        setJobs(res.data || []);
      } catch (err) {
        setError("Failed to load jobs. Please try again later.");
        setJobs([]);
      }
    };

    fetchJobs();
  }, []);

  const formatHoursPerWeek = (value) => {
    if (!value) return "N/A";
    if (value === "<20") return "Less than 20 hours/week";
    if (value === "80+") return "More than 80 hours/week";
    return `${value} hours/week`;
  };

  const formatProjectType = (type) => {
    if (!type) return "N/A";
    return type === "ongoing" ? "Ongoing" : "One-Time";
  };

  const formatProjectLength = (length) => {
    switch (length) {
      case "<1":
        return "Less than 1 month";
      case "1-3":
        return "1–3 months";
      case "3-6":
        return "3–6 months";
      case "6-12":
        return "6–12 months";
      case "12+":
        return "More than 12 months";
      default:
        return "N/A";
    }
  };

  const handleFilterChange = (e, filterName) => {
    const { value, checked } = e.target;

    setFilters((prev) => {
      let updated = [...prev[filterName]];
      updated = checked
        ? [...updated, value]
        : updated.filter((item) => item !== value);
      return { ...prev, [filterName]: updated };
    });
  };

  const filteredJobs = (jobs || []).filter((job) => {
    if (filters.type.length && !filters.type.includes(job.type)) return false;
    if (filters.length.length && !filters.length.includes(job.length))
      return false;
    if (filters.rate.length && !filters.rate.includes(job.rate)) return false;
    if (
      filters.availability.length &&
      !filters.availability.includes(job.hours_per_week)
    )
      return false;
    return true;
  });

  const indexOfLastJob = currentPage * jobsPerPage;
  const indexOfFirstJob = indexOfLastJob - jobsPerPage;
  const currentJobs = filteredJobs.slice(indexOfFirstJob, indexOfLastJob);

  if (jobs === null) {
    return <p className="text-center mt-5">Loading...</p>;
  }

  if (error) {
    return (
      <Container className="py-5">
        <Alert variant="danger" className="text-center">
          {error}
        </Alert>
      </Container>
    );
  }

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 job-section-title">Available Jobs</h2>

      {/* Filters */}
      <Row className="mb-4">
        {/* Project Type Filter */}
        <Col md={3}>
          <DropdownButton
            variant="secondary"
            id="dropdown-type"
            title={
              filters.type.length > 0
                ? `Project Type: ${filters.type.join(", ")}`
                : "Select Project Type"
            }
          >
            <Form.Check
              type="checkbox"
              label="Ongoing"
              value="ongoing"
              checked={filters.type.includes("ongoing")}
              onChange={(e) => handleFilterChange(e, "type")}
            />
            <Form.Check
              type="checkbox"
              label="One-Time"
              value="one-time"
              checked={filters.type.includes("one-time")}
              onChange={(e) => handleFilterChange(e, "type")}
            />
          </DropdownButton>
        </Col>

        {/* Project Length Filter */}
        <Col md={3}>
          <DropdownButton
            variant="secondary"
            id="dropdown-length"
            title={
              filters.length.length > 0
                ? `Project Length: ${filters.length.join(", ")}`
                : "Select Project Length"
            }
          >
            {[
              { label: "Less than 1 month", value: "<1" },
              { label: "1–3 months", value: "1-3" },
              { label: "3–6 months", value: "3-6" },
              { label: "6–12 months", value: "6-12" },
              { label: "More than 12 months", value: "12+" },
            ].map(({ label, value }) => (
              <Form.Check
                key={value}
                type="checkbox"
                label={label}
                value={value}
                checked={filters.length.includes(value)}
                onChange={(e) => handleFilterChange(e, "length")}
              />
            ))}
          </DropdownButton>
        </Col>

        {/* Rate Filter */}
        <Col md={3}>
          <DropdownButton
            variant="secondary"
            id="dropdown-rate"
            title={
              filters.rate.length > 0
                ? `Rate: ${filters.rate.join(", ")}`
                : "Select Rate"
            }
          >
            <Form.Check
              type="checkbox"
              label="Hourly"
              value="hourly"
              checked={filters.rate.includes("hourly")}
              onChange={(e) => handleFilterChange(e, "rate")}
            />
            <Form.Check
              type="checkbox"
              label="Fixed"
              value="fixed"
              checked={filters.rate.includes("fixed")}
              onChange={(e) => handleFilterChange(e, "rate")}
            />
          </DropdownButton>
        </Col>

        {/* Availability Filter */}
        <Col md={3}>
          <DropdownButton
            variant="secondary"
            id="dropdown-availability"
            title={
              filters.availability.length > 0
                ? `Availability: ${filters.availability.join(", ")}`
                : "Select Availability"
            }
          >
            {[
              { label: "Less than 20 hours/week", value: "<20" },
              { label: "20–40 hours/week", value: "20-40" },
              { label: "40–60 hours/week", value: "40-60" },
              { label: "60–80 hours/week", value: "60-80" },
              { label: "More than 80 hours/week", value: "80+" },
            ].map(({ label, value }) => (
              <Form.Check
                key={value}
                type="checkbox"
                label={label}
                value={value}
                checked={filters.availability.includes(value)}
                onChange={(e) => handleFilterChange(e, "availability")}
              />
            ))}
          </DropdownButton>
        </Col>
      </Row>

      {/* Empty State */}
      {filteredJobs.length === 0 ? (
        <div className="text-center my-4">
          <p>No suitable jobs were found.</p>
        </div>
      ) : (
        <>
          {/* Job Cards */}
          <Row>
            {currentJobs.map((job) => (
              <Col key={job.id} md={4} sm={6} xs={12} className="mb-4">
                <Link to={`/jobs/${job.id}`} className="text-decoration-none">
                  <Card className="job-card">
                    <Card.Body className="text-center">
                      <h5 className="fw-bold job-title">{job.title}</h5>
                      <p className="job-comp text-capitalize">
                        {job.rate === "hourly"
                          ? `$${job.amount}/h`
                          : `$${job.amount}`}{" "}
                        • {job.rate}
                      </p>
                      <p className="job-meta">
                        <strong>Project Type:</strong>{" "}
                        <span className="job-meta-value">
                          {formatProjectType(job.type)}
                        </span>
                      </p>
                      <p className="job-meta">
                        <strong>Project Length:</strong>{" "}
                        <span className="job-meta-value">
                          {formatProjectLength(job.length)}
                        </span>
                      </p>
                      <p className="job-meta">
                        <strong>Availability:</strong>{" "}
                        <span className="job-meta-value">
                          {formatHoursPerWeek(job.hours_per_week)}
                        </span>
                      </p>
                    </Card.Body>
                  </Card>
                </Link>
              </Col>
            ))}
          </Row>

          {/* Pagination */}
          <Pagination className="justify-content-center mt-4 custom-pagination">
            {[
              ...Array(Math.ceil(filteredJobs.length / jobsPerPage)).keys(),
            ].map((number) => (
              <Pagination.Item
                key={number + 1}
                active={number + 1 === currentPage}
                onClick={() => setCurrentPage(number + 1)}
              >
                {number + 1}
              </Pagination.Item>
            ))}
          </Pagination>
        </>
      )}
    </Container>
  );
};

export default JobsPage;
