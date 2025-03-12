import { useState, useEffect } from "react";
import axiosInstance from "../utils/axios";

const ITEMS_PER_PAGE = 10;

const UserControlPanel = () => {
  const [users, setUsers] = useState([]);
  const [page, setPage] = useState(1);

  useEffect(() => {
    fetchUsers();
  }, []);

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
      .then(() => {
        fetchUsers();
      })
      .catch((error) => {
        fetchUsers();
        console.error("Error updating user", error);
      });
  };

  const handleDeleteUser = (id) => {
    if (!window.confirm("Are you sure you want to delete this user?")) return;

    axiosInstance
      .delete(`/admin/users/${id}`)
      .then(() => {
        fetchUsers();
      })
      .catch((error) => {
        console.error("Error deleting user", error);
      });
  };

  const handleChange = (id, field, value) => {
    setUsers((prev) =>
      prev.map((user) => (user.id === id ? { ...user, [field]: value } : user))
    );
  };

  const handleBlur = (user) => {
    handleUpdateUser(user.id, { role: user.role, ban: user.ban });
  };

  const totalPages = Math.ceil(users.length / ITEMS_PER_PAGE);
  const paginatedUsers = users.slice(
    (page - 1) * ITEMS_PER_PAGE,
    page * ITEMS_PER_PAGE
  );

  return (
    <div className="container mt-4">
      <>
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
                    onBlur={() => handleBlur(user)}
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
                    onBlur={() => handleBlur(user)}
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
    </div>
  );
};

export default UserControlPanel;
