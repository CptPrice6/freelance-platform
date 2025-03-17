import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";
import { Card, Container, Row, Col, Pagination } from "react-bootstrap";
import "../styles/FreelancersPage.css"; // Import CSS for extra styling

const FreelancersPage = () => {
  const [freelancers, setFreelancers] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const freelancersPerPage = 6;

  useEffect(() => {
    axiosInstance.get("/freelancers").then((res) => {
      setFreelancers(res.data);
    });
  }, []);

  // Pagination logic
  const indexOfLastFreelancer = currentPage * freelancersPerPage;
  const indexOfFirstFreelancer = indexOfLastFreelancer - freelancersPerPage;
  const currentFreelancers = freelancers.slice(
    indexOfFirstFreelancer,
    indexOfLastFreelancer
  );

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 section-title">Meet Our Freelancers</h2>

      <Row>
        {currentFreelancers.map((freelancer) => (
          <Col key={freelancer.id} md={4} sm={6} xs={12} className="mb-4">
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
