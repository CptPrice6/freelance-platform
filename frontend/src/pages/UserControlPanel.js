import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";

const ITEMS_PER_PAGE = 10;

const UserControlPanel = () => {
  const [users, setUsers] = useState([]);
  const [filteredUsers, setFilteredUsers] = useState([]);
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState("");
  const [updateError, setUpdateError] = useState("");
  const [updateSuccess, setUpdateSuccess] = useState("");

  useEffect(() => {
    fetchUsers();
  }, []);

  useEffect(() => {
    if (!searchQuery.trim()) {
      setFilteredUsers(users);
    } else {
      setFilteredUsers(
        users.filter((user) =>
          user.email.toLowerCase().includes(searchQuery.toLowerCase())
        )
      );
    }
    setPage(1);
  }, [users, searchQuery]);

  const fetchUsers = () => {
    axiosInstance
      .get("/admin/users")
      .then((response) => {
        const nonAdminUsers = response.data.filter(
          (user) => user.role !== "admin"
        );
        setUsers(nonAdminUsers);
      })
      .catch((error) => {
        console.error("Error fetching users", error);
      });
  };

  const handleUpdateUser = (id, updatedUser) => {
    axiosInstance
      .put(`/admin/users/${id}`, updatedUser)
      .then((response) => {
        fetchUsers();
        const updateMessage = response?.data?.message;
        setUpdateError("");
        setUpdateSuccess(updateMessage);
      })
      .catch((error) => {
        fetchUsers();
        console.error("Error updating user", error);
        const errorMessage = error.response?.data?.error || error.message;
        setUpdateSuccess("");
        setUpdateError(errorMessage);
      });
  };

  const handleDeleteUser = (id) => {
    if (!window.confirm("Are you sure you want to delete this user?")) return;

    axiosInstance
      .delete(`/admin/users/${id}`)
      .then((response) => {
        fetchUsers();
        const updateMessage = response?.data?.message;
        setUpdateError("");
        setUpdateSuccess(updateMessage);
      })
      .catch((error) => {
        console.error("Error deleting user", error);
        const errorMessage = error.response?.data?.error || error.message;
        setUpdateSuccess("");
        setUpdateError(errorMessage);
      });
  };

  const handleChange = (id, field, value) => {
    const updatedUsers = users.map((user) =>
      user.id === id ? { ...user, [field]: value } : user
    );
    setUsers(updatedUsers);

    const updatedUser = updatedUsers.find((user) => user.id === id);
    handleUpdateUser(id, { role: updatedUser.role, ban: updatedUser.ban });
  };

  const totalPages = Math.ceil(filteredUsers.length / ITEMS_PER_PAGE);
  const paginatedUsers = filteredUsers.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE
  );

  return (
    <div className="container mt-4">
      {/* Display Success/Error Messages */}
      {updateError && (
        <div className="alert alert-danger mt-3" role="alert">
          {updateError}
        </div>
      )}
      {updateSuccess && (
        <div className="alert alert-success mt-3" role="alert">
          {updateSuccess}
        </div>
      )}

      <div className="mb-3">
        <input
          type="text"
          className="form-control"
          placeholder="Search by email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <table className="table table-striped">
        <thead>
          <tr>
            <th>ID</th>
            <th>Email</th>
            <th>Role</th>
            <th>Ban</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {paginatedUsers.map((user) => (
            <tr key={user.id}>
              <td>{user.id}</td>
              <td>{user.email}</td>
              <td>
                <select
                  className="form-select"
                  value={user.role}
                  onChange={(e) =>
                    handleChange(user.id, "role", e.target.value)
                  }
                >
                  <option value="freelancer">Freelancer</option>
                  <option value="client">Client</option>
                  <option value="admin">Admin</option>
                </select>
              </td>
              <td>
                <select
                  className="form-select"
                  value={user.ban.toString()}
                  onChange={(e) =>
                    handleChange(user.id, "ban", e.target.value === "true")
                  }
                >
                  <option value="true">Banned</option>
                  <option value="false">Active</option>
                </select>
              </td>
              <td>
                <button
                  className="btn btn-danger"
                  onClick={() => handleDeleteUser(user.id)}
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
              <button className="page-link" onClick={() => setPage(index + 1)}>
                {index + 1}
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  );
};

export default UserControlPanel;
