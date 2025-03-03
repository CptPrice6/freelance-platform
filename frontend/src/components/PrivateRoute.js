import { Navigate } from "react-router-dom";
import { getAccessToken } from "../utils/tokens";

const PrivateRoute = ({ element }) => {
  const accessToken = getAccessToken();

  if (!accessToken) {
    return <Navigate to="/login" replace />;
  }

  return element;
};

export default PrivateRoute;
