import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";

const ITEMS_PER_PAGE = 10;

const SkillControlPanel = () => {
  const [skills, setSkills] = useState([]);
  const [newSkill, setNewSkill] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    fetchSkills();
  }, []);

  const fetchSkills = () => {
    setLoading(true);
    axiosInstance
      .get("/skills")
      .then((response) => {
        setSkills(response.data);
      })
      .catch((error) => {
        console.error("Error fetching skills", error);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const handleUpdateSkill = (id, updatedName) => {
    axiosInstance
      .put(`/admin/skills/${id}`, { skill_name: updatedName })
      .then(() => {
        fetchSkills();
      })
      .catch((error) => {
        console.error("Error updating skill", error);
      });
  };

  const handleDeleteSkill = (id) => {
    if (!window.confirm("Are you sure you want to delete this skill?")) return;

    axiosInstance
      .delete(`/admin/skills/${id}`)
      .then(() => {
        fetchSkills();
      })
      .catch((error) => {
        console.error("Error deleting skill", error);
      });
  };

  const handleAddSkill = (e) => {
    e.preventDefault();
    if (!newSkill.trim()) return;

    axiosInstance
      .post("/admin/skills", { skill_name: newSkill.trim() })
      .then(() => {
        setNewSkill("");
        setErrorMessage("");
        setShowModal(false);
        fetchSkills();
      })
      .catch((error) => {
        if (
          error.response &&
          error.response.data.error.includes("unique constraint")
        ) {
          setErrorMessage("Skill already exists!");
        } else {
          console.error("Error adding skill", error);
        }
      });
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setErrorMessage("");
    setNewSkill("");
  };

  const totalPages = Math.ceil(skills.length / ITEMS_PER_PAGE);
  const paginatedSkills = skills.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE
  );

  return (
    <div className="container mt-4">
      <h1 className="mb-4 d-flex justify-content-between align-items-center">
        <button className="btn btn-success" onClick={() => setShowModal(true)}>
          +
        </button>
      </h1>

      {loading ? (
        <p>Loading skills...</p>
      ) : (
        <>
          <table className="table table-striped">
            <thead>
              <tr>
                <th>ID</th>
                <th>Skill Name</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {paginatedSkills.map((skill) => (
                <tr key={skill.Id}>
                  <td>{skill.Id}</td>
                  <td>
                    <input
                      type="text"
                      className="form-control"
                      value={skill.Name}
                      onChange={(e) => {
                        const updatedName = e.target.value;
                        setSkills((prev) =>
                          prev.map((s) =>
                            s.Id === skill.Id ? { ...s, Name: updatedName } : s
                          )
                        );
                      }}
                      onBlur={() => handleUpdateSkill(skill.Id, skill.Name)}
                    />
                  </td>
                  <td>
                    <button
                      className="btn btn-danger"
                      onClick={() => handleDeleteSkill(skill.Id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {/* Pagination */}
          <nav>
            <ul className="pagination">
              {[...Array(totalPages)].map((_, index) => (
                <li
                  key={index}
                  className={`page-item ${page === index + 1 ? "active" : ""}`}
                >
                  <button
                    className="page-link"
                    onClick={() => setPage(index + 1)}
                  >
                    {index + 1}
                  </button>
                </li>
              ))}
            </ul>
          </nav>
        </>
      )}

      {/* Add Skill Modal */}
      {showModal && (
        <div className="modal show d-block" tabIndex="-1">
          <div className="modal-dialog">
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">Add New Skill</h5>
                <button
                  type="button"
                  className="btn-close"
                  onClick={handleCloseModal}
                ></button>
              </div>
              <div className="modal-body">
                <input
                  type="text"
                  className="form-control"
                  placeholder="Enter skill name"
                  value={newSkill}
                  onChange={(e) => setNewSkill(e.target.value)}
                />
                {errorMessage && (
                  <div className="text-danger mt-2">{errorMessage}</div>
                )}
              </div>
              <div className="modal-footer">
                <button
                  type="button"
                  className="btn btn-secondary"
                  onClick={handleCloseModal}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  className="btn btn-primary"
                  onClick={handleAddSkill}
                  disabled={!newSkill.trim()}
                >
                  Add Skill
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Modal Backdrop */}
      {showModal && <div className="modal-backdrop show"></div>}
    </div>
  );
};

export default SkillControlPanel;
