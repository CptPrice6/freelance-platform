import { useState } from "react";
import UserControlPanel from "./UserControlPanel";
import SkillControlPanel from "./SkillControlPanel";

const AdminDashboard = () => {
  const [activePanel, setActivePanel] = useState("users");

  return (
    <div className="container mt-4">
      <div className="d-flex justify-content-between align-items-center mb-4">
        <h1>
          {activePanel === "users" ? "User Management" : "Skill Management"}
        </h1>

        <div className="btn-group">
          <button
            className={`btn ${
              activePanel === "users" ? "btn-primary" : "btn-outline-primary"
            }`}
            onClick={() => setActivePanel("users")}
          >
            Users
          </button>
          <button
            className={`btn ${
              activePanel === "skills" ? "btn-primary" : "btn-outline-primary"
            }`}
            onClick={() => setActivePanel("skills")}
          >
            Skills
          </button>
        </div>
      </div>

      {activePanel === "users" ? <UserControlPanel /> : <SkillControlPanel />}
    </div>
  );
};

export default AdminDashboard;
