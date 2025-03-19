import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";
import { Card, Container, Row, Col, Pagination } from "react-bootstrap";
import "../styles/ClientsPage.css";
import { Link } from "react-router-dom";

const ClientsPage = () => {
  const [clients, setClients] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const clientsPerPage = 9;

  useEffect(() => {
    axiosInstance.get("/clients").then((res) => {
      setClients(res.data);
    });
  }, []);

  const indexOfLastClient = currentPage * clientsPerPage;
  const indexOfFirstClient = indexOfLastClient - clientsPerPage;
  const currentClients = clients.slice(indexOfFirstClient, indexOfLastClient);

  return (
    <Container className="py-5">
      <h2 className="text-center mb-4 section-title">Our Valued Clients</h2>

      <Row>
        {currentClients.map((client) => (
          <Col key={client.id} md={4} sm={6} xs={12} className="mb-4">
            <Link to={`/clients/${client.id}`} className="text-decoration-none">
              <Card className="client-card">
                <Card.Body className="text-center">
                  <h5 className="fw-bold client-name">
                    {client.name} {client.surname}
                  </h5>
                  <p className="client-company">
                    🏢 {client.company_name || "No Company Name"}
                  </p>
                </Card.Body>
              </Card>
            </Link>
          </Col>
        ))}
      </Row>

      {/* Pagination */}
      <Pagination className="justify-content-center mt-4 custom-pagination">
        {[...Array(Math.ceil(clients.length / clientsPerPage)).keys()].map(
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
    </Container>
  );
};

export default ClientsPage;
