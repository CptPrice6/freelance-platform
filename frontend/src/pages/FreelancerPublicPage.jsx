import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { Card, Container, Row, Col, Badge } from "react-bootstrap";
import axiosInstance from "../utils/axios";
import "../styles/FreelancerPublicPage.css";

const FreelancerPublicPage = () => {
  const { id } = useParams();
  const [freelancer, setFreelancer] = useState(null);
  const [skills, setSkills] = useState([]);

  useEffect(() => {
    axiosInstance.get(`/freelancers/${id}`).then((res) => {
      const freelancerData = res.data.freelancer_data;
      setFreelancer(res.data);
      setSkills(freelancerData.skills || []);
    });
  }, [id]);

  if (!freelancer) return <p className="text-center mt-5">Loading...</p>;

  return (
    <Container className="py-5 d-flex justify-content-center">
      <Card className="public-freelancer-card shadow-lg p-4">
        <Card.Body>
          <h2 className="text-center">
            {freelancer.name} {freelancer.surname}{" "}
          </h2>
          <h5 className="text-center text-muted">
            {freelancer.freelancer_data.title || "No Title"}
          </h5>

          <Row className="mt-4 justify-content-center">
            <Col md={4} className="d-flex justify-content-center">
              <div className="circle-info">
                <p className="circle-label">Pay Rate</p>
                <p className="circle-value">
                  {freelancer.freelancer_data.hourly_rate
                    ? `$${freelancer.freelancer_data.hourly_rate}/h`
                    : "Not Set"}
                </p>
              </div>
            </Col>
            <Col md={4} className="d-flex justify-content-center">
              <div className="circle-info">
                <p className="circle-label">Availability</p>
                <p className="circle-value">
                  {freelancer.freelancer_data.hours_per_week
                    ? `${freelancer.freelancer_data.hours_per_week}h/week`
                    : "Not Set"}
                </p>
              </div>
            </Col>
          </Row>

          <Row className="mt-4">
            <Col md={12}>
              <p>
                <strong>Description:</strong>
              </p>
              <p className="text-muted">
                {freelancer.freelancer_data.description ||
                  "No description provided"}
              </p>
            </Col>
          </Row>

          <div className="mt-4">
            <h5>Skills:</h5>
            {skills.length > 0 ? (
              <div className="skills-container">
                {skills.map((skill) => (
                  <Badge key={skill.id} bg="primary" className="skill-badge">
                    {skill.name}
                  </Badge>
                ))}
              </div>
            ) : (
              <p className="text-muted">No skills listed</p>
            )}
          </div>
        </Card.Body>
      </Card>
    </Container>
  );
};

export default FreelancerPublicPage;
