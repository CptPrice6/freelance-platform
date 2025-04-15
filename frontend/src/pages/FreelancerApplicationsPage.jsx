import { useState, useEffect } from "react";
import {
  Card,
  Container,
  Row,
  Col,
  Badge,
  Button,
  Modal,
  Form,
  Pagination,
} from "react-bootstrap";
import axiosInstance from "../utils/axios";
import moment from "moment";
import "../styles/FreelancerApplicationsPage.css";

const FreelancerApplicationsPage = () => {
  const [applications, setApplications] = useState([]);
  const [selectedApp, setSelectedApp] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [updatedDescription, setUpdatedDescription] = useState("");
  const [file, setFile] = useState(null);
  const [fileBase64, setFileBase64] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const appsPerPage = 10;

  useEffect(() => {
    axiosInstance
      .get("/user/freelancer/applications")
      .then((res) => setApplications(res.data || []))
      .catch((err) => {
        console.error("Failed to fetch applications:", err);
        setApplications([]);
      });
  }, []);

  const handleCardClick = (app) => {
    setSelectedApp(app);
    setUpdatedDescription(app.description);
    setFile(null);
    setFileBase64("");
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setSelectedApp(null);
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

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];

    if (selectedFile && selectedFile.type !== "application/pdf") {
      alert("Only PDF files are allowed.");
      return;
    }

    setFile(selectedFile);

    const reader = new FileReader();
    reader.readAsDataURL(selectedFile);
    reader.onload = () => {
      const base64String = reader.result.split(",")[1];
      setFileBase64(base64String);
    };
  };

  const handleUpdate = async () => {
    if (!updatedDescription) {
      alert("Please enter a description.");
      return;
    }

    const payload = {
      description: updatedDescription,
      file_name: file?.name || "",
      file_base64: fileBase64 || "",
    };

    try {
      await axiosInstance.put(
        `/user/freelancer/applications/${selectedApp.id}`,
        payload
      );
      alert("Application updated.");
      handleCloseModal();
      window.location.reload();
    } catch (err) {
      console.error(err);
      alert("Failed to update application.");
    }
  };

  const handleDelete = async () => {
    try {
      await axiosInstance.delete(
        `/user/freelancer/applications/${selectedApp.id}`
      );
      alert("Application deleted.");
      handleCloseModal();
      setApplications(applications.filter((app) => app.id !== selectedApp.id));
    } catch (err) {
      console.error(err);
      alert("Failed to delete application.");
    }
  };

  const indexOfLastApp = currentPage * appsPerPage;
  const indexOfFirstApp = indexOfLastApp - appsPerPage;
  const currentApplications = applications.slice(
    indexOfFirstApp,
    indexOfLastApp
  );

  const getStatusVariant = (status) => {
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

  return (
    <Container className="py-5">
      <h2 className="mb-4 text-center application-section-title">
        Your Applications
      </h2>

      {currentApplications.length === 0 ? (
        <div className="text-center my-4">
          <p>You have no applications currently.</p>
        </div>
      ) : (
        <>
          <Row>
            {applications.map((app) => (
              <Col key={app.id} xs={12} className="mb-3">
                <Card
                  onClick={() => handleCardClick(app)}
                  className="freelancer-application-card p-4 shadow-sm"
                  style={{ cursor: "pointer" }}
                >
                  <Card.Body>
                    <div className="d-flex justify-content-between align-items-center mb-2">
                      <h5 className="fw-bold application-title">
                        Job Title: {app.job_title}
                      </h5>
                      <Badge
                        bg={getStatusVariant(app.status)}
                        className="application-status-badge"
                      >
                        {app.status}
                      </Badge>
                    </div>
                    <p className="mb-0 text-muted">
                      <strong>ID:</strong> {app.id}
                    </p>
                    <p className="mb-0 text-muted">
                      <strong>Submitted:</strong>{" "}
                      {moment(app.created_at).format("LLL")}
                    </p>
                  </Card.Body>
                </Card>
              </Col>
            ))}
          </Row>

          <Pagination className="justify-content-center mt-4 custom-pagination">
            {[
              ...Array(Math.ceil(applications.length / appsPerPage)).keys(),
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

      {selectedApp && (
        <Modal show={showModal} onHide={handleCloseModal}>
          <Modal.Header closeButton>
            <Modal.Title>Application #{selectedApp.id}</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <div className="mb-3">
              <strong>Job:</strong>{" "}
              <a
                href={`/jobs/${selectedApp.job_id}`}
                target="_blank"
                rel="noopener noreferrer"
              >
                {selectedApp.job_title}
              </a>
            </div>

            <Form.Group className="mb-3">
              <Form.Label>Description</Form.Label>
              <Form.Control
                as="textarea"
                rows={4}
                disabled={selectedApp.status !== "pending"}
                value={updatedDescription}
                onChange={(e) => setUpdatedDescription(e.target.value)}
              />
            </Form.Group>

            <div className="mb-3">
              <Form.Label>Attachment</Form.Label>
              {selectedApp.attachment?.file_name ? (
                <div>
                  <Button
                    variant="link"
                    onClick={() => handleDownload(selectedApp.attachment.id)}
                  >
                    📄 {selectedApp.attachment.file_name}
                  </Button>
                </div>
              ) : (
                <p className="text-muted">No attachment</p>
              )}
              {selectedApp.status === "pending" && (
                <Form.Control
                  type="file"
                  accept="application/pdf"
                  onChange={handleFileChange}
                />
              )}
            </div>

            <div>
              <strong>Status:</strong>{" "}
              <Badge
                className="application-status-badge"
                bg={getStatusVariant(selectedApp.status)}
              >
                {selectedApp.status}
              </Badge>
            </div>
          </Modal.Body>
          <Modal.Footer>
            <Button variant="secondary" onClick={handleCloseModal}>
              Close
            </Button>
            {selectedApp.status === "pending" && (
              <>
                <Button variant="danger" onClick={handleDelete}>
                  Delete
                </Button>
                <Button variant="success" onClick={handleUpdate}>
                  Update
                </Button>
              </>
            )}
          </Modal.Footer>
        </Modal>
      )}
    </Container>
  );
};

export default FreelancerApplicationsPage;
