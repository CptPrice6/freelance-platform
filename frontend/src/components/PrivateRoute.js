import { Navigate } from "react-router-dom";

const PrivateRoute = ({ element }) => {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    return <Navigate to="/login" replace />;
  }

  return element;
};

export default PrivateRoute;
