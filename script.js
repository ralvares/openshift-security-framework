let rawData = {};
let skills = {};
let responsibilities = {};

async function loadData() {
  const res = await fetch("data/mapping.json");
  rawData = await res.json();

  skills = rawData.skills;
  responsibilities = rawData.responsibilities;

  populateRoles();
}

function populateRoles() {
  const select = document.getElementById("role");
  select.innerHTML = "";

  for (const roleName in rawData.roles) {
    const option = document.createElement("option");
    option.value = roleName;
    option.textContent = roleName;
    select.appendChild(option);
  }

  select.selectedIndex = 0;
  renderContent();
}

function renderContent() {
  const selectedRoles = Array.from(document.getElementById("role").selectedOptions).map(opt => opt.value);

  if (selectedRoles.length === 0) {
    document.getElementById("description").innerHTML = "";
    document.getElementById("responsibilities").innerHTML = "";
    document.getElementById("mapping").innerHTML = "";
    return;
  }

  let combinedResponsibilities = new Set();
  let combinedSkills = { Basic: new Set(), Intermediate: new Set(), Advanced: new Set() };
  let descriptions = [];

  selectedRoles.forEach(roleKey => {
    const role = rawData.roles[roleKey];
    descriptions.push(`<h3>${roleKey}</h3><p>${role.description}</p>`);
    role.responsibilities.forEach(id => combinedResponsibilities.add(id));

    for (const level in role.skills) {
      role.skills[level].forEach(skillId => combinedSkills[level].add(skillId));
    }
  });

  document.getElementById("description").innerHTML = `<h2>Description</h2>` + descriptions.join("");

  const renderedResponsibilities = Array.from(combinedResponsibilities)
    .map(id => `<li>${responsibilities[id]}</li>`)
    .join("");
  document.getElementById("responsibilities").innerHTML = `<h2>Core Responsibilities</h2><ul>${renderedResponsibilities}</ul>`;

  let html = "";
  for (const level of ["Basic", "Intermediate", "Advanced"]) {
    const ids = Array.from(combinedSkills[level]).sort((a, b) => a.localeCompare(b, undefined, { numeric: true }));
    if (ids.length === 0) continue;

    html += `<h2>${level} Level</h2><table><tr><th>Skill ID</th><th>Description</th><th>OpenShift/K8s Relevance</th></tr>`;
    ids.forEach(skillId => {
      const skill = skills[skillId];
      html += `<tr><td>${skillId}</td><td>${skill.desc}</td><td>${skill.relevance}</td></tr>`;
    });
    html += `</table>`;
  }

  document.getElementById("mapping").innerHTML = html;
}

loadData();
