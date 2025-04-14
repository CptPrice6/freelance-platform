import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";
import { Card, Container, Row, Col, Pagination, Alert } from "react-bootstrap";
import "../styles/FreelancersPage.css";
import { Link } from "react-router-dom";

const FreelancersPage = () => {
  const [freelancers, setFreelancers] = useState(null);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const freelancersPerPage = 9;

  useEffect(() => {
    const fetchFreelancers = async () => {
      try {
        const res = await axiosInstance.get("/freelancers");
        setFreelancers(res.data || []);
      } catch (err) {
        setError("Failed to load freelancers. Please try again later.");
        setFreelancers([]);
      }
    };

    fetchFreelancers();
  }, []);

  if (freelancers === null) {
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

  if (freelancers.length === 0) {
    return <p className="text-center mt-5 text-muted">No freelancers found.</p>;
  }

  const indexOfLastFreelancer = currentPage * freelancersPerPage;
  const indexOfFirstFreelancer = indexOfLastFreelancer - freelancersPerPage;
  const currentFreelancers = freelancers.slice(
    indexOfFirstFreelancer,
    indexOfLastFreelancer
  );

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 freelancer-section-title">
        Meet Our Freelancers
      </h2>

      <Row>
        {currentFreelancers.map((freelancer) => (
          <Col key={freelancer.id} md={4} sm={6} xs={12} className="mb-4">
            <Link
              to={`/freelancers/${freelancer.id}`}
              className="text-decoration-none"
            >
              <Card className="freelancer-card">
                <Card.Body className="text-center">
                  <h5 className="fw-bold freelancer-name">
                    {freelancer.name} {freelancer.surname}
                  </h5>
                  <p className="freelancer-title">
                    {freelancer.title || "No Title"}
                  </p>
                  <p className="freelancer-rate">
                    {freelancer.hourly_rate
                      ? `$${freelancer.hourly_rate}/h`
                      : "Rate Not Set"}
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
          ...Array(Math.ceil(freelancers.length / freelancersPerPage)).keys(),
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
    </Container>
  );
};

export default FreelancersPage;
