import { Navigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import { getAccessToken } from "../utils/tokens";

const AdminRoute = ({ element }) => {
  const accessToken = getAccessToken();

  if (!accessToken) {
    return <Navigate to="/login" replace />;
  }

  try {
    const decoded = jwtDecode(accessToken);

    // If the user is not an admin, redirect them
    if (decoded.role !== "admin") {
      return <Navigate to="/" replace />;
    }
  } catch (err) {
    return <Navigate to="/login" replace />;
  }

  return element;
};

export default AdminRoute;
