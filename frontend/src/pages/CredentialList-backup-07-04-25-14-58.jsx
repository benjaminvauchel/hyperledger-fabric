import { useEffect, useState } from "react";
import axios from "axios";

function CredentialList({ baseUrl }) {
  const [credentials, setCredentials] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchCredentials = async () => {
    setLoading(true);
    setError(null);

    const url = `${baseUrl}/credentials`;
    const params = {
      chaincodeid: "basic",
      channelid: "mychannel"
    };

    console.log(`Fetching credentials from: ${url}`);
    console.log("Request params:", params);

    try {
      const res = await axios.get(url, { params });
      console.log("API response received:", res);

      if (res.data && res.data.success) {
        let parsedData;
        if (typeof res.data.data === "string") {
          try {
            parsedData = JSON.parse(res.data.data);
          } catch (parseErr) {
            console.error("Error parsing response data:", parseErr);
            parsedData = [];
          }
        } else if (Array.isArray(res.data.data)) {
          parsedData = res.data.data;
        } else {
          console.warn("Unexpected data format:", res.data.data);
          parsedData = [];
        }

        console.log("Processed credential data:", parsedData);
        setCredentials(parsedData);
      } else {
        throw new Error(res.data?.error || "API returned unsuccessful response");
      }
    } catch (err) {
      console.error("Fetch error details:", err);

      if (err.response) {
        console.error("Error response data:", err.response.data);
        console.error("Error response status:", err.response.status);
        console.error("Error response headers:", err.response.headers);
        setError(`Server error: ${err.response.status} - ${err.response.data?.error || err.message}`);
      } else if (err.request) {
        console.error("Error request:", err.request);
        setError("No response received from server. Check that the API is running.");
      } else {
        console.error("Error message:", err.message);
        setError(err.message || "Unknown error");
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCredentials();
  }, [baseUrl]);

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
          <button
            onClick={() => console.log("Current baseUrl:", baseUrl)}
            className="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
          >
            Log API URL
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
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
          </tr>
        </thead>
        <tbody>
          {credentials && credentials.length > 0 ? (
            credentials.map((cred) => (
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
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="11" className="text-center py-4 text-gray-500">
                No credentials found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
      <div className="mt-4 text-right">
        <button 
          onClick={fetchCredentials}
          className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
        >
          Refresh
        </button>
      </div>
    </div>
  );
}

export default CredentialList;

