import { useEffect, useState } from "react";
import axios from "axios";

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));


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
      await sleep(3000);
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

  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this credential?")) return;

    setUpdating(true);
    try {
      await axios.delete(`${baseUrl}/credentials/${id}`, {
        params: {
          chaincodeid: "basic",
          channelid: "mychannel"
        }
      });
      await sleep(3000);
      await fetchCredentials();
    } catch (err) {
      alert("Failed to delete credential: " + (err.response?.data?.error || err.message));
    } finally {
      setUpdating(false);
    }
  };

  const handleUpdateSkills = async (id, currentSkills) => {
    const newSkills = prompt("Enter new skills:", currentSkills);
    if (!newSkills || newSkills === currentSkills) return;

    try {
      await axios.put(`${baseUrl}/credentials/${id}/skills`, { newSkills, chaincodeid: "basic", channelid: "mychannel" });
      await sleep(3000);
      await fetchCredentials();
    } catch (err) {
      alert("Failed to update skills: " + (err.response?.data?.error || err.message));
    }
  };

  const handleUpdateName = async (id, currentFirst, currentLast) => {
    const newFirst = prompt("New first name:", currentFirst);
    const newLast = prompt("New last name:", currentLast);
    if (!newFirst || !newLast || (newFirst === currentFirst && newLast === currentLast)) return;

    try {
      await axios.put(`${baseUrl}/credentials/${id}/name`, {
        newFirstName: newFirst,
        newLastName: newLast,
        chaincodeid: "basic",
        channelid: "mychannel"
      });
      await sleep(3000);
      await fetchCredentials();
    } catch (err) {
      alert("Failed to update name: " + (err.response?.data?.error || err.message));
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
    <div className="max-w-full">
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

	    <div className="overflow-x-auto w-full">
        <table className="w-full table-fixed bg-white shadow-md rounded text-sm">
          <thead className="bg-gray-200">
            <tr>
              <th className="w-16 p-2">ID</th>
              <th className="w-24">Name</th>
              <th className="w-20">Type</th>
              <th className="w-24">Institution</th>
              <th className="w-24">Education</th>
              <th className="w-24">Company</th>
              <th className="w-24">Work Exp.</th>
              <th className="w-24">Skills</th>
              <th className="w-20">Status</th>
              <th className="w-20">Verified By</th>
              <th className="w-16">Talent ID</th>
              <th className="w-32">Actions</th>
            </tr>
          </thead>
          <tbody>
            {filteredCredentials.length > 0 ? (
              filteredCredentials.map((cred) => (
                <tr key={cred.CredentialID} className="border-t">
                  <td className="p-2 font-mono text-sm whitespace-normal break-words">{cred.CredentialID}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.FirstName} {cred.LastName}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.CredentialType}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.CredentialType === "academic" ? cred.Institution || "-" : "-"}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.CredentialType === "academic" ? cred.Education || "-" : "-"}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.CredentialType === "professional" ? cred.Company || "-" : "-"}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.CredentialType === "professional" ? cred.WorkExperience || "-" : "-"}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.Skills}</td>
                  <td className="p-1 whitespace-normal break-words">
                    <span className={`px-2 py-1 rounded text-sm ${
                      cred.VerificationStatus === "Verified" ? "bg-green-100 text-green-800" : 
                      cred.VerificationStatus === "Pending" ? "bg-yellow-100 text-yellow-800" : 
                      "bg-gray-100 text-gray-800"
                    }`}>
                      {cred.VerificationStatus}
                    </span>
                  </td>
                  <td className="p-1 whitespace-normal break-words">{cred.VerifiedBy || "-"}</td>
                  <td className="p-1 whitespace-normal break-words">{cred.TalentID}</td>
                  <td className="p-1">
                    {/* Group 1: Approve and Revoke */}
                    <div className="flex gap-1 mb-1 justify-between">
                      {cred.VerificationStatus !== "Verified" ? (
                        <button
                          onClick={() => handleAction(cred.CredentialID, "approve")}
                          className="flex-1 px-1 py-0.5 bg-green-200 text-green-800 rounded text-xs"
                        >
                          Approve
                        </button>
                      ) : (
                        <div className="flex-1"></div>
                      )}
                      {cred.VerificationStatus !== "Revoked" ? (
                        <button
                          onClick={() => handleAction(cred.CredentialID, "revoke")}
                          className="flex-1 px-1 py-0.5 bg-red-200 text-red-800 rounded text-xs"
                        >
                          Revoke
                        </button>
                      ) : (
                        <div className="flex-1"></div>
                      )}
                    </div>
                    
                    {/* Group 2: Edit Skills and Edit Name */}
                    <div className="flex gap-1 mb-1 justify-between">
                      <button
                        onClick={() => handleUpdateSkills(cred.CredentialID, cred.Skills)}
                        className="flex-1 px-1 py-0.5 bg-blue-100 text-blue-800 rounded text-xs"
                      >
                        Edit Skills
                      </button>
                      <button
                        onClick={() => handleUpdateName(cred.CredentialID, cred.FirstName, cred.LastName)}
                        className="flex-1 px-1 py-0.5 bg-yellow-100 text-yellow-800 rounded text-xs"
                      >
                        Edit Name
                      </button>
                    </div>
                    
                    {/* Group 3: Delete */}
                    <div>
                      <button
                        onClick={() => handleDelete(cred.CredentialID)}
                        className="w-full px-1 py-0.5 bg-gray-200 text-gray-700 rounded text-xs"
                      >
                        Delete
                      </button>
                    </div>
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
    </div>
  );
}

export default CredentialList;
