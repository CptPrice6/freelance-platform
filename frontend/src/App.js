import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Layout from "./components/Layout";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import ProfileSettings from "./pages/ProfileSettings";
import PrivateRoute from "./components/PrivateRoute";
import AdminRoute from "./components/AdminRoute";
import AdminPanel from "./pages/AdminPanel";

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

            {/*Full list of Jobs and Freelancers */}
            <Route
              path="/freelancers"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/jobs"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            {/* Specific Job and Freelancer Detail Pages */}
            <Route
              path="/jobs/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/freelancers/:id"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/freelancer/dashboard"
              element={<PrivateRoute element={<ProfileSettings />} />}
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
              path="/client/dashboard"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />
            <Route
              path="/client/jobs"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/settings"
              element={<PrivateRoute element={<ProfileSettings />} />}
            />

            <Route
              path="/admin/dashboard"
              element={<AdminRoute element={<AdminPanel />} />}
            />
          </Routes>
        </Layout>
      </div>
    </Router>
  );
}

export default App;
