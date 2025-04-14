import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { Card, Container, Row, Col } from "react-bootstrap";
import axiosInstance from "../utils/axios";
import "../styles/ClientPublicPage.css";

const ClientPublicPage = () => {
  const { id } = useParams();
  const [client, setClient] = useState(null);

  useEffect(() => {
    axiosInstance
      .get(`/clients/${id}`)
      .then((res) => {
        setClient(res.data);
      })
      .catch((err) => {
        if (err.response && err.response.status === 404) {
          setClient("not-found");
        } else {
          console.error("Failed to fetch client:", err);
        }
      });
  }, [id]);

  if (client === null) return <p className="text-center mt-5">Loading...</p>;
  if (client === "not-found") {
    return <p className="text-center mt-5 text-danger">Client not found.</p>;
  }

  return (
    <Container className="py-5 d-flex justify-content-center">
      <Card className="public-client-card shadow-lg p-4">
        <Card.Body>
          <h2 className="text-center">
            {client.name} {client.surname}
          </h2>
          <h5 className="text-center text-muted">
            🏢 {client.client_data.company_name || "No Company Name"}
          </h5>

          <Row className="mt-4">
            <Col md={12}>
              <p>
                <strong>Description:</strong>
              </p>
              <p className="text-muted">
                {client.client_data.description || "No description provided"}
              </p>
            </Col>
            <Col md={12}>
              <p>
                <strong>Industry:</strong>
              </p>
              <p className="text-muted">
                {client.client_data.industry || "No industry provided"}
              </p>
            </Col>
            <Col md={12}>
              <p>
                <strong>Location:</strong>
              </p>
              <p className="text-muted">
                {client.client_data.location || "No location provided"}
              </p>
            </Col>
          </Row>
        </Card.Body>
      </Card>
    </Container>
  );
};

export default ClientPublicPage;
