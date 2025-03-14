import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";

const ITEMS_PER_PAGE = 10;

const SkillControlPanel = () => {
  const [skills, setSkills] = useState([]);
  const [filteredSkills, setFilteredSkills] = useState([]);
  const [newSkill, setNewSkill] = useState("");
  const [showModal, setShowModal] = useState(false);
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [errorMessageModal, setErrorMessageModal] = useState("");
  const [errorMessageMenu, setErrorMessageMenu] = useState("");
  const [successMessageMenu, setSuccessMessageMenu] = useState("");
  const [editingSkillId, setEditingSkillId] = useState(null);
  const [editedSkillName, setEditedSkillName] = useState("");

  useEffect(() => {
    fetchSkills();
  }, []);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredSkills(skills);
    } else {
      setFilteredSkills(
        skills.filter((skill) =>
          skill.name.toLowerCase().includes(searchQuery.toLowerCase())
        )
      );
    }
    setPage(1);
  }, [skills, searchQuery]);

  const fetchSkills = () => {
    axiosInstance
      .get("/skills")
      .then((response) => {
        setSkills(response.data);
        setFilteredSkills(response.data);
      })
      .catch((error) => {
        console.error("Error fetching skills", error);
      });
  };

  const handleEditClick = (skill) => {
    setEditingSkillId(skill.id);
    setEditedSkillName(skill.name);
  };

  const handleSave = (id) => {
    if (!editedSkillName.trim()) {
      setSuccessMessageMenu("");
      setErrorMessageMenu("Please enter skill name");
      return;
    }

    axiosInstance
      .put(`/admin/skills/${id}`, { skill_name: editedSkillName.trim() })
      .then((response) => {
        fetchSkills();
        setEditingSkillId(null);
        setEditedSkillName("");
        setErrorMessageMenu("");
        setSuccessMessageMenu(response?.data?.message);
      })
      .catch((error) => {
        console.error("Error updating skill", error);
        setSuccessMessageMenu("");
        setErrorMessageMenu(error.response?.data?.error || error.message);
      });
  };

  const handleKeyDown = (e, id) => {
    if (e.key === "Enter") {
      handleSave(id);
    }
  };

  const handleDeleteSkill = (id) => {
    if (!window.confirm("Are you sure you want to delete this skill?")) return;

    axiosInstance
      .delete(`/admin/skills/${id}`)
      .then((response) => {
        fetchSkills();
        const updateMessage = response?.data?.message;
        setErrorMessageMenu("");
        setSuccessMessageMenu(updateMessage);
      })
      .catch((error) => {
        console.error("Error deleting skill", error);
        const errorMessage = error.response?.data?.error || error.message;
        setSuccessMessageMenu("");
        setErrorMessageMenu(errorMessage);
      });
  };

  const handleAddSkill = (e) => {
    e.preventDefault();
    if (!newSkill.trim()) return;

    axiosInstance
      .post("/admin/skills", { skill_name: newSkill.trim() })
      .then((response) => {
        setNewSkill("");
        setErrorMessageModal("");
        setShowModal(false);
        fetchSkills();

        const updateMessage = response?.data?.message;
        setErrorMessageMenu("");
        setSuccessMessageMenu(updateMessage);
      })
      .catch((error) => {
        if (
          error.response &&
          error.response.data.error.includes("already exists")
        ) {
          setErrorMessageModal(error.response?.data?.error);
        } else {
          console.error("Error adding skill", error);
        }
      });
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setErrorMessageModal("");
    setNewSkill("");
  };

  const totalPages = Math.ceil(filteredSkills.length / ITEMS_PER_PAGE);
  const paginatedSkills = filteredSkills.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE
  );

  return (
    <div className="container mt-4">
      <>
        {/* Display Success/Error Messages */}
        {errorMessageMenu && (
          <div className="alert alert-danger mt-3" role="alert">
            {errorMessageMenu}
          </div>
        )}
        {successMessageMenu && (
          <div className="alert alert-success mt-3" role="alert">
            {successMessageMenu}
          </div>
        )}

        {/* Search Input */}
        <div className="mb-3">
          <input
            type="text"
            className="form-control"
            placeholder="Search by name..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

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
              <tr key={skill.id}>
                <td>{skill.id}</td>
                <td>
                  {editingSkillId === skill.id ? (
                    <input
                      type="text"
                      className="form-control"
                      value={editedSkillName}
                      onChange={(e) => setEditedSkillName(e.target.value)}
                      onKeyDown={(e) => handleKeyDown(e, skill.id)}
                    />
                  ) : (
                    skill.name
                  )}
                </td>
                <td>
                  {editingSkillId === skill.id ? (
                    <button
                      className="btn btn-primary me-2"
                      onClick={() => handleSave(skill.id)}
                    >
                      Save
                    </button>
                  ) : (
                    <button
                      className="btn btn-warning me-2"
                      onClick={() => handleEditClick(skill)}
                    >
                      Edit
                    </button>
                  )}
                  <button
                    className="btn btn-danger"
                    onClick={() => handleDeleteSkill(skill.id)}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Pagination and Add Skill Button */}
        <div className="d-flex justify-content-between align-items-center">
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

          <button
            className="btn btn-success"
            onClick={() => setShowModal(true)}
          >
            +
          </button>
        </div>
      </>

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
                {errorMessageModal && (
                  <div className="text-danger mt-2">{errorMessageModal}</div>
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

      {showModal && <div className="modal-backdrop show"></div>}
    </div>
  );
};

export default SkillControlPanel;
