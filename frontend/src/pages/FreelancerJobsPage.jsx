import { useEffect, useState } from "react";
import axiosInstance from "../utils/axios";
import { Container, Card, Row, Col, Badge, Pagination } from "react-bootstrap";
import { Link } from "react-router-dom";
import "../styles/FreelancerJobsPage.css";

const FreelancerJobsPage = () => {
  const [jobs, setJobs] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const jobsPerPage = 10;

  useEffect(() => {
    axiosInstance
      .get("/user/freelancer/jobs")
      .then((res) => setJobs(res.data || []))
      .catch((err) => {
        console.error("Error fetching freelancer jobs:", err);
        setJobs([]);
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

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 job-section-title">Your Jobs</h2>

      {currentJobs.length === 0 ? (
        <div className="text-center my-4">
          <p>You have no jobs currently.</p>
        </div>
      ) : (
        <>
          <Row>
            {currentJobs.map((job) => (
              <Col key={job.id} xs={12} className="mb-3">
                <Link to={`/jobs/${job.id}`} className="text-decoration-none">
                  <Card className="freelancer-job-card p-4 shadow-sm">
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
                      <p className="mb-0 text-muted">
                        <strong>ID:</strong> {job.id}
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
    </Container>
  );
};

export default FreelancerJobsPage;
