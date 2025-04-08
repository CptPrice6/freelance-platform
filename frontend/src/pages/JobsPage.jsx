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
} from "react-bootstrap";
import "../styles/JobsPage.css";
import { Link } from "react-router-dom";

const JobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [filters, setFilters] = useState({
    type: [],
    length: [],
    rate: [],
    availability: [],
  });
  const jobsPerPage = 9;

  useEffect(() => {
    axiosInstance.get("/jobs").then((res) => {
      setJobs(res.data);
    });
  }, []);

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

  // Handle filter change for dropdowns with checkboxes
  const handleFilterChange = (e, filterName) => {
    const { value, checked } = e.target;

    setFilters((prevFilters) => {
      let updatedFilter = [...prevFilters[filterName]];

      if (checked) {
        updatedFilter.push(value);
      } else {
        updatedFilter = updatedFilter.filter((item) => item !== value);
      }

      return { ...prevFilters, [filterName]: updatedFilter };
    });
  };

  // Apply filters to the jobs list
  const filteredJobs = jobs.filter((job) => {
    // Filter by type
    if (filters.type.length > 0 && !filters.type.includes(job.type))
      return false;

    // Filter by length
    if (filters.length.length > 0 && !filters.length.includes(job.length))
      return false;

    // Filter by rate
    if (filters.rate.length > 0 && !filters.rate.includes(job.rate))
      return false;

    // Filter by availability (hours per week)
    if (
      filters.availability.length > 0 &&
      !filters.availability.includes(job.hours_per_week)
    )
      return false;

    return true;
  });

  // Pagination logic
  const indexOfLastJob = currentPage * jobsPerPage;
  const indexOfFirstJob = indexOfLastJob - jobsPerPage;
  const currentJobs = filteredJobs.slice(indexOfFirstJob, indexOfLastJob);

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 job-section-title">Available Jobs</h2>

      {/* Filters Section */}
      <Row className="mb-4">
        {/* Project Type Dropdown */}
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

        {/* Project Length Dropdown */}
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
            <Form.Check
              type="checkbox"
              label="Less than 1 month"
              value="<1"
              checked={filters.length.includes("<1")}
              onChange={(e) => handleFilterChange(e, "length")}
            />
            <Form.Check
              type="checkbox"
              label="1–3 months"
              value="1-3"
              checked={filters.length.includes("1-3")}
              onChange={(e) => handleFilterChange(e, "length")}
            />
            <Form.Check
              type="checkbox"
              label="3–6 months"
              value="3-6"
              checked={filters.length.includes("3-6")}
              onChange={(e) => handleFilterChange(e, "length")}
            />
            <Form.Check
              type="checkbox"
              label="6–12 months"
              value="6-12"
              checked={filters.length.includes("6-12")}
              onChange={(e) => handleFilterChange(e, "length")}
            />
            <Form.Check
              type="checkbox"
              label="More than 12 months"
              value="12+"
              checked={filters.length.includes("12+")}
              onChange={(e) => handleFilterChange(e, "length")}
            />
          </DropdownButton>
        </Col>

        {/* Rate Dropdown */}
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

        {/* Availability Dropdown */}
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
            <Form.Check
              type="checkbox"
              label="Less than 20 hours/week"
              value="<20"
              checked={filters.availability.includes("<20")}
              onChange={(e) => handleFilterChange(e, "availability")}
            />
            <Form.Check
              type="checkbox"
              label="20-40hours/week"
              value="20-40"
              checked={filters.availability.includes("20-40")}
              onChange={(e) => handleFilterChange(e, "availability")}
            />
            <Form.Check
              type="checkbox"
              label="40-60hours/week"
              value="40-60"
              checked={filters.availability.includes("40-60")}
              onChange={(e) => handleFilterChange(e, "availability")}
            />
            <Form.Check
              type="checkbox"
              label="60-80hours/week"
              value="60-80"
              checked={filters.availability.includes("60-80")}
              onChange={(e) => handleFilterChange(e, "availability")}
            />
            <Form.Check
              type="checkbox"
              label="More than 80 hours/week"
              value="80+"
              checked={filters.availability.includes("80+")}
              onChange={(e) => handleFilterChange(e, "availability")}
            />
          </DropdownButton>
        </Col>
      </Row>

      {/* No jobs found message */}
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
                        <strong>Project Type:</strong>
                        <span className="job-meta-value">
                          {formatProjectType(job.type)}
                        </span>
                      </p>
                      <p className="job-meta">
                        <strong>Project Length:</strong>
                        <span className="job-meta-value">
                          {formatProjectLength(job.length)}
                        </span>
                      </p>
                      <p className="job-meta">
                        <strong>Availability:</strong>
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
