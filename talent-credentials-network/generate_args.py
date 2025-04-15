import json
import random

skills = ["Python", "Java", "Go", "Rust", "C++", "Kubernetes", "React", "ML", "Data Science"]
degrees = ["B.Sc.", "M.Sc.", "PhD"]
schools = ["Concordia University", "McGill", "Université de Montréal", "Polytechnique"]

data = []
n_entries = 500

for i in range(n_entries):
    tid = f"talent{i:04}"
    first = f"User{i}"
    last = f"Test{i}"
    skillset = ", ".join(random.sample(skills, random.randint(2, 5)))
    degree = random.choice(degrees)
    school = random.choice(schools)

    # Create array of strings for each entry to match the expected format
    data.append([tid, first, last, skillset, degree, school])

with open("args.json", "w") as f:
    json.dump(data, f, indent=2)

print(f"✅ Generated args.json with {n_entries} entries")
