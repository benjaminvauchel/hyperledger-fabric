import { useState } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import Home from "./pages/Home";
import CredentialForm from "./pages/CredentialForm";
import CredentialList from "./pages/CredentialList";

function App() {
  const [org, setOrg] = useState("Org1"); // Org1 or Org2
  const baseUrl = org === "Org1" ? "http://localhost:3000" : "http://localhost:3001";

  return (
    <Router>
      <div className="min-h-screen bg-gray-100 text-gray-900">
        <header className="bg-white shadow p-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">Talent Credentials Network</h1>
          <select
            className="border p-1 rounded"
            value={org}
            onChange={(e) => setOrg(e.target.value)}
          >
            <option value="Org1">Institution / Company</option>
            <option value="Org2">Talent</option>
          </select>
        </header>

        <nav className="bg-gray-200 p-4 flex gap-4">
          <Link to="/" className="hover:underline">Home</Link>
          <Link to="/create" className="hover:underline">Create Credential</Link>
          <Link to="/list" className="hover:underline">View Credentials</Link>
        </nav>

        <main className="container mx-auto p-4 max-w-7xl overflow-hidden">
          <Routes>
            <Route path="/" element={<Home org={org} />} />
            <Route path="/create" element={<CredentialForm baseUrl={baseUrl} />} />
            <Route path="/list" element={<CredentialList baseUrl={baseUrl} />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;

