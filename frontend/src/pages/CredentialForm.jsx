import { useState } from "react";
import axios from "axios";

function CredentialForm({ baseUrl }) {
  const [type, setType] = useState("academic");
  const [form, setForm] = useState({
    id: "",
    firstName: "",
    lastName: "",
    institution: "",
    company: "",
    education: "",
    workExperience: "",
    skills: "",
    talentID: ""
  });

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async () => {
    try {
      const endpoint = `${baseUrl}/credentials/${type}`;
      const payload = {
        credential: {
          CredentialID: form.id,
          FirstName: form.firstName,
          LastName: form.lastName,
          Skills: form.skills,
          TalentID: form.talentID,
          ...(type === "academic"
            ? {
                Institution: form.institution,
                Education: form.education
              }
            : {
                Company: form.company,
                WorkExperience: form.workExperience
              })
        },
        chaincodeid: "basic",
        channelid: "mychannel"
      };

      const res = await axios.post(endpoint, payload);
      alert(`Success: ${res.data.message}`);
    } catch (err) {
      alert(`Error: ${err.response?.data?.error || "Unknown error"}`);
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <div className="flex gap-4 mb-4">
        <button
          onClick={() => setType("academic")}
          className={`px-4 py-2 rounded ${type === "academic" ? "bg-blue-600 text-white" : "bg-white border"}`}
        >
          Academic
        </button>
        <button
          onClick={() => setType("professional")}
          className={`px-4 py-2 rounded ${type === "professional" ? "bg-blue-600 text-white" : "bg-white border"}`}
        >
          Professional
        </button>
      </div>

      <input type="text" name="id" placeholder="Credential ID" className="input" onChange={handleChange} />
      <input type="text" name="firstName" placeholder="First Name" className="input" onChange={handleChange} />
      <input type="text" name="lastName" placeholder="Last Name" className="input" onChange={handleChange} />
      <input type="text" name="talentID" placeholder="Talent ID" className="input" onChange={handleChange} />

      {type === "academic" ? (
        <>
          <input type="text" name="institution" placeholder="Institution" className="input" onChange={handleChange} />
          <input type="text" name="education" placeholder="Education" className="input" onChange={handleChange} />
        </>
      ) : (
        <>
          <input type="text" name="company" placeholder="Company" className="input" onChange={handleChange} />
          <input type="text" name="workExperience" placeholder="Work Experience" className="input" onChange={handleChange} />
        </>
      )}

      <input type="text" name="skills" placeholder="Skills (comma-separated)" className="input" onChange={handleChange} />

      <button onClick={handleSubmit} className="mt-4 w-full bg-blue-600 text-white py-2 rounded">
        Submit {type} Credential
      </button>
    </div>
  );
}

export default CredentialForm;

