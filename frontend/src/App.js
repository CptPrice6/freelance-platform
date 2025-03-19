import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import ProfileSettings from "./pages/ProfileSettings";
import PrivateRoute from "./components/PrivateRoute";
import AdminRoute from "./components/AdminRoute";
import AdminDashboard from "./pages/AdminDashboard";
import FreelancerDashboard from "./pages/FreelancerDashboard";
import ClientDashboard from "./pages/ClientDashboard";
import FreelancersPage from "./pages/FreelancersPage";
import FreelancerPublicPage from "./pages/FreelancerPublicPage";
import ClientsPage from "./pages/ClientsPage";
import ClientPublicPage from "./pages/ClientPublicPage";

function App() {
  return (
    <Router>
      <div className="App">
        {/* Wrap all routes with the Layout component */}
        <Layout>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

            <Route
              path="/freelancers"
              element={<PrivateRoute element={<FreelancersPage />} />}
            />
            <Route
              path="/freelancers/:id"
              element={<PrivateRoute element={<FreelancerPublicPage />} />}
            />
            <Route
              path="/clients"
              element={<PrivateRoute element={<ClientsPage />} />}
            />
            <Route
              path="/clients/:id"
              element={<PrivateRoute element={<ClientPublicPage />} />}
            />
            <Route
              path="/jobs"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/jobs/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/clients"
              element={<PrivateRoute element={<FreelancersPage />} />}
            />
            <Route
              path="/clients/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/freelancer/dashboard"
              element={<PrivateRoute element={<FreelancerDashboard />} />}
            />
            <Route
              path="/freelancer/applications"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/freelancer/jobs"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/freelancer/jobs/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/client/dashboard"
              element={<PrivateRoute element={<ClientDashboard />} />}
            />
            <Route
              path="/client/jobs"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/client/jobs/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/settings"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/admin/dashboard"
              element={<AdminRoute element={<AdminDashboard />} />}
            />
          </Routes>
        </Layout>
      </div>
    </Router>
  );
}

export default App;
