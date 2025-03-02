import { Navigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const AdminRoute = ({ element, ...rest }) => {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    return <Navigate to="/login" />;
  }

  try {
    const decoded = jwtDecode(accessToken);

    // If the user is not an admin, redirect them
    if (decoded.role !== "admin") {
      return <Navigate to="/" />;
    }
  } catch (err) {
    return <Navigate to="/login" />;
  }

  return element;
};

export default AdminRoute;
