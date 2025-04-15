import { useEffect, useState } from "react";
import axios from "axios";

function CredentialList({ baseUrl }) {
  const [credentials, setCredentials] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [typeFilter, setTypeFilter] = useState("all");
  const [issuerFilter, setIssuerFilter] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [updating, setUpdating] = useState(false);

  const fetchCredentials = async () => {
    setLoading(true);
    setError(null);

    const url = `${baseUrl}/credentials`;
    const params = {
      chaincodeid: "basic",
      channelid: "mychannel"
    };

    try {
      const res = await axios.get(url, { params });
      if (res.data && res.data.success) {
        let parsedData;
        if (typeof res.data.data === "string") {
          try {
            parsedData = JSON.parse(res.data.data);
          } catch (parseErr) {
            parsedData = [];
          }
        } else if (Array.isArray(res.data.data)) {
          parsedData = res.data.data;
        } else {
          parsedData = [];
        }
        setCredentials(parsedData);
      } else {
        throw new Error(res.data?.error || "API returned unsuccessful response");
      }
    } catch (err) {
      if (err.response) {
        setError(`Server error: ${err.response.status} - ${err.response.data?.error || err.message}`);
      } else if (err.request) {
        setError("No response received from server. Check that the API is running.");
      } else {
        setError(err.message || "Unknown error");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleAction = async (id, action) => {
    setUpdating(true);
    try {
      await axios.put(`${baseUrl}/credentials/${id}/${action}`, null, {
        params: {
          chaincodeid: "basic",
          channelid: "mychannel"
        }
      });
      await new Promise((resolve) => setTimeout(resolve, 3000));
      await fetchCredentials();
    } catch (err) {
      let message = `Failed to ${action}`;
      if (err.response) {
        if (err.response.status === 403) {
          message = "You are not authorized to perform this action.";
        } else if (err.response.data?.error) {
          message = `${message}: ${err.response.data.error}`;
        }
      } else if (err.request) {
        message = `No response from server while trying to ${action}`;
      } else {
        message = `${message}: ${err.message}`;
      }

      alert(message);
    } finally {
      setUpdating(false);
    }
  };


  useEffect(() => {
    fetchCredentials();
  }, [baseUrl]);

  const filteredCredentials = credentials.filter((cred) => {
    const matchesType = typeFilter === "all" || cred.CredentialType === typeFilter;
    const issuer = cred.Institution || cred.Company || "";
    const matchesIssuer = issuer.toLowerCase().includes(issuerFilter.toLowerCase());
    const matchesStatus = statusFilter === "all" || cred.VerificationStatus === statusFilter;
    return matchesType && matchesIssuer && matchesStatus;
  });

  if (loading) {
    return <div className="text-center py-4">Loading credentials...</div>;
  }

  if (error) {
    return (
      <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded">
        <p className="font-bold">Error loading credentials</p>
        <p className="mb-2">{error}</p>
        <div className="flex gap-2">
          <button 
            onClick={fetchCredentials}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  const clearFilters = () => {
    setTypeFilter("all");
    setIssuerFilter("");
    setStatusFilter("all");
  };

  return (
    <div className="overflow-x-auto">
      <div className="flex flex-wrap gap-4 items-center mb-4">
        <div>
          <label className="mr-2">Filter by Type:</label>
          <select
            value={typeFilter}
            onChange={(e) => setTypeFilter(e.target.value)}
            className="p-2 border rounded"
          >
            <option value="all">All</option>
            <option value="academic">Academic</option>
            <option value="professional">Professional</option>
          </select>
        </div>
        <div>
          <label className="mr-2">Filter by Institution/Company:</label>
          <input
            type="text"
            value={issuerFilter}
            onChange={(e) => setIssuerFilter(e.target.value)}
            className="p-2 border rounded"
            placeholder="Search issuer..."
          />
        </div>
        <div>
          <label className="mr-2">Filter by Status:</label>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="p-2 border rounded"
          >
            <option value="all">All</option>
            <option value="Verified">Verified</option>
            <option value="Pending">Pending</option>
          </select>
        </div>
        <div className="flex gap-2 ml-auto">
          <button
            onClick={clearFilters}
            className="px-4 py-2 bg-red-100 text-red-800 rounded hover:bg-red-200"
          >
            Clear Filters
          </button>
          <button 
            onClick={fetchCredentials}
            className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
          >
            Refresh
          </button>
        </div>
      </div>

      {updating && (
        <div className="text-sm text-blue-600 mb-2 ml-1">
          Updating credentials...
        </div>
      )}

      <table className="w-full table-auto bg-white shadow-md rounded">
        <thead className="bg-gray-200">
          <tr>
            <th className="p-2">ID</th>
            <th>Name</th>
            <th>Type</th>
            <th>Institution</th>
            <th>Education</th>
            <th>Company</th>
            <th>Work Experience</th>
            <th>Skills</th>
            <th>Status</th>
            <th>Verified By</th>
            <th>Talent ID</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {filteredCredentials.length > 0 ? (
            filteredCredentials.map((cred) => (
              <tr key={cred.CredentialID} className="border-t">
                <td className="p-2 font-mono text-sm">{cred.CredentialID}</td>
                <td>{cred.FirstName} {cred.LastName}</td>
                <td>{cred.CredentialType}</td>
                <td>{cred.CredentialType === "academic" ? cred.Institution || "-" : "-"}</td>
                <td>{cred.CredentialType === "academic" ? cred.Education || "-" : "-"}</td>
                <td>{cred.CredentialType === "professional" ? cred.Company || "-" : "-"}</td>
                <td>{cred.CredentialType === "professional" ? cred.WorkExperience || "-" : "-"}</td>
                <td>{cred.Skills}</td>
                <td>
                  <span className={`px-2 py-1 rounded text-sm ${
                    cred.VerificationStatus === "Verified" ? "bg-green-100 text-green-800" : 
                    cred.VerificationStatus === "Pending" ? "bg-yellow-100 text-yellow-800" : 
                    "bg-gray-100 text-gray-800"
                  }`}>
                    {cred.VerificationStatus}
                  </span>
                </td>
                <td>{cred.VerifiedBy || "-"}</td>
                <td>{cred.TalentID}</td>
                <td>
                  {cred.VerificationStatus !== "Verified" && (
                    <button
                      onClick={() => handleAction(cred.CredentialID, "approve")}
                      className="px-2 py-1 bg-green-200 text-green-800 rounded hover:bg-green-300 mb-1 block"
                    >
                      Approve
                    </button>
                  )}
                  {cred.VerificationStatus !== "Revoked" && (
                    <button
                      onClick={() => handleAction(cred.CredentialID, "revoke")}
                      className="px-2 py-1 bg-red-200 text-red-800 rounded hover:bg-red-300 block"
                    >
                      Revoke
                    </button>
                  )}
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="12" className="text-center py-4 text-gray-500">
                No credentials found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default CredentialList;

