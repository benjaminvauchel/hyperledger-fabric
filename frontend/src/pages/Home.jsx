function Home({ org }) {
  return (
    <div className="text-center">
      <h2 className="text-2xl font-semibold mb-2">Welcome, {org}</h2>
      {org === "Org1" ? (
        <p>You can create, approve, or revoke credentials.</p>
      ) : (
        <p>You can submit and view your credentials.</p>
      )}
    </div>
  );
}

export default Home;

