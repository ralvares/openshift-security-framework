# 🚧 Work In Progress (WIP)

---

## Why This Project?

I’ve built an OpenShift Security Enablement Framework—a structured, role-based approach to help teams understand what security looks like in practice, not just in theory. It maps common roles (like Developers, Platform Engineers, Architects, and more) to clear responsibilities and skills, from basic to advanced levels.

The goal is to:
- Simplify conversations about OpenShift security.
- Help align expectations early (before jumping straight into ACS).
- Enable everyone — from sales to delivery — to talk about value, not just features.
- Provide a foundation for training, workshops, and customer enablement.

If you’re working with customers who struggle to understand where ACS fits, or are comparing it to “all-in-one” security tools, this framework can help reset the conversation and clarify what OpenShift secures by design, and where ACS adds value.

**Your feedback is invaluable!**

**If this framework resonates with you or sparks ideas, let’s connect and explore ways to collaborate!**

---

# OpenShift Security Role-Based Skills Framework

The **OpenShift Security Role-Based Skills Framework** is a structured, dynamic model that maps Kubernetes/OpenShift security concepts to key technical personas within an organization. It empowers teams by aligning daily responsibilities with actionable security skills, tailored by role and maturity level—from foundational knowledge to advanced leadership.

---

## Purpose

This framework helps organizations:
- Translate OpenShift/Kubernetes security capabilities into practical, role-specific learning paths.
- Enable tailored upskilling for developers, platform engineers, architects, security teams, and compliance officers.
- Drive adoption of Red Hat OpenShift’s built-in and extended security features (like RHACS, Compliance Operator, audit logs).
- Bridge the gap between security requirements and day-to-day platform operations or app development.

---

## How It Works

The framework is structured into three key layers:

### 1. Skills Library
A reusable set of security skills, each mapped to OpenShift-native capabilities (e.g., SCCs, RBAC, audit logs). Skills are grouped by:
- **Basic:** Foundational knowledge (e.g., shared responsibility, OWASP Top 10)
- **Intermediate:** Integration and operationalization (e.g., CI/CD pipelines, incident response)
- **Advanced:** Leadership, policy creation, automation

### 2. Responsibilities
Common day-to-day actions or duties per role, such as maintaining CI/CD pipelines or enforcing security controls, reused across multiple roles to avoid duplication.

### 3. Role Mapping
Each role includes:
- A short description
- A list of associated responsibilities (linked by ID)
- Skill expectations by level (Basic, Intermediate, Advanced)

---

## Who It’s For

This framework is tailored for multiple roles involved in the software delivery and security lifecycle, including:
- **Application Developers**
- **Software Developers**
- **Infrastructure Engineers**
- **Security Architects**
- **DevSecOps Engineers**
- **Network Engineers**
- **Security Engineers**
- **Security Operations Specialists**

Each role is mapped to relevant responsibilities and skill levels.

---

## Features

- **Interactive Role Selection:** Select one or more roles to view their descriptions, core responsibilities, and mapped skills by level.
- **Reusable Skills & Responsibilities:** Skills and responsibilities are referenced by ID, making the system DRY and scalable.
- **Extensible:** Add new roles, skills, or responsibilities by updating the JSON and Markdown files.
- **Sales & Learning Ready:** Use for onboarding, capability assessment, training plans, or to explain the business value of OpenShift Security to different teams.

---

## Roles at a Glance (Current Canonical Set)

Below is the concise intent of each role after applying the “cleaner ownership” model (reduced overlap, clearer accountability). Overlaps that remain are intentional for hand‑off points.

### Application Developer
Purpose: Build and ship secure workloads; integrate day‑one security hygiene (least privilege, secrets handling, secure image construction).
Responsibilities (IDs): R1, R2, R3, R4, R5, R6, R7, R37  
Focus Themes: Threat modeling at app scope (I1), secret usage (B5), secure/reproducible image builds & basic signing (I11), baseline security guardrails (B3/B9).  
Not Expected: Platform governance, runtime detection tuning, policy exception governance (kept out to minimize cognitive load).  
Advanced Skills: None (by design) – escalation path is usually into DevSecOps or Architecture for advanced domains.

### Platform Operator
Purpose: Operate, harden, and maintain compliant OpenShift clusters; enforce day‑to‑day platform guardrails.
Responsibilities: R8, R9, R10, R11, R12, R13, R14, R28, R30, R32  
Focus Themes: Node hardening, admission of trusted artifacts (indirect via pipelines), audit/log integration, compliance scans, quota & secret lifecycle operations.  
Advanced Scope: Limited to A1 (governance participation), A7 (automation/remediation), A8 (enablement). Strategic multi‑cluster & identity architecture moved to Architect.

### DevSecOps Engineer
Purpose: Embed security in delivery workflows; enforce policy gates; provide feedback loops from runtime to build; own supply chain policy execution.
Responsibilities: R13, R19, R20, R21, R25, R26, R31, R33  
Focus Themes: Secure pipelines (I5), artifact signing gates (R31), admission / policy controller implementations (R21), runtime detection rule tuning (R33), shift‑left threat modeling alignment (I1), infrastructure as code security (R20).  
Boundary: Does not own exception governance (R35) or platform quota/secret governance (R30/R32).

### Security Architect
Purpose: Set platform-wide security strategy, governance, multi‑cluster posture, identity, and exception processes.
Responsibilities: R6, R12, R15, R16, R17, R18, R21, R27, R29, R34, R35, R36  
Focus Themes: Governance & policy definition, architectural threat modeling patterns, multi‑cluster policy orchestration, workload identity design, exception lifecycle, organizational enablement.  
Operational items (quotas, secret rotation execution, supply chain gate operation) deliberately removed to reduce overlap.

### Network & Infrastructure Engineer (Specialized / Optional)
Purpose: Engineer the secure data plane: segmentation, encryption patterns, trust fabric, identity‑aware routing, traffic anomaly foundations.
Responsibilities: R22, R23, R24, R34, R36, R41  
Focus Themes: Advanced network segmentation & egress/ingress governance (I4/A3), encryption architecture implementation (B10/A9), workload/service identity runtime integration (A10), trust fabric engineering (R41).  
Not Owning: Compliance automation (R12), supply chain gates (R31), artifact provenance policy design (A6) – those stay with DevSecOps / Architect.

---

## Responsibility Reference (IDs → Meaning)

| ID | Summary |
|----|---------|
| R1 | Application threat modeling (app/workload level) |
| R2 | Access control via RBAC & service accounts |
| R3 | Secure secret usage & encryption enablement |
| R4 | Integrate scanning & testing in CI/CD |
| R5 | Maintain updated dependencies & images |
| R6 | Apply platform security best practices in workloads |
| R7 | Collaborate on audit/log monitoring & response |
| R8 | Configure SCC/pod security, RBAC & NetworkPolicies (platform) |
| R9 | Node hardening & patch management |
| R10 | Audit log aggregation & forwarding |
| R11 | Vulnerability scan + remediation (platform scope) |
| R12 | Automate compliance scanning |
| R13 | Embed controls & gates in pipelines |
| R14 | Provide security training/documentation |
| R15 | Design secure platform architectures |
| R16 | Lead governance & Zero Trust policy definition |
| R17 | Automate vulnerability assessment workflows |
| R18 | Mentor teams / maturity uplift |
| R19 | Build & manage secure deployment pipelines |
| R20 | Secure infrastructure as code artifacts |
| R21 | Enforce policies via admission / controllers |
| R22 | Engineer network infrastructure & segmentation |
| R23 | Implement encryption (transit + at rest) |
| R24 | Monitor network anomalies / DoS mitigation |
| R25 | Threat modeling + runtime monitoring + incident response (execution) |
| R26 | Implement workload isolation & Zero Trust models |
| R27 | Manage security infrastructure (IAM, SIEM, EDR, policy) |
| R28 | Maintain security policy documentation |
| R29 | Align security with product/lifecycle goals |
| R30 | Define / manage quotas & limits |
| R31 | Operate artifact signing & attestation verification gates |
| R32 | Execute secret lifecycle rotation & integration |
| R33 | Tune runtime detection rules & telemetry |
| R34 | Orchestrate multi-cluster security & compliance policy |
| R35 | Govern security control exceptions |
| R36 | Workload identity issuance & trust policies |
| R37 | Produce secure, minimal, reproducible images with basic signing |
| R41 | Engineer service-to-service trust fabric |

---

## Skill Interpretation

Skills are grouped by difficulty/impact, not years of experience.

| Tier | Meaning | Typical Output Examples |
|------|---------|-------------------------|
| Basic | Knows concepts & safely uses them | Non-root image, baseline NetworkPolicy, uses Secrets properly |
| Intermediate | Integrates & adapts controls | Adds scanning to pipeline, rotates secrets, builds reproducible images |
| Advanced | Designs strategy / automation / governance | Multi-cluster policy blueprint, exception workflow, identity trust fabric |

Each role’s `skills` object lists only the skills expected *for that role*. Absence of a skill does not mean the persona could not learn it—just that it is not baseline for the role.

---

## Reading `mapping.json`

Top-level keys:
1. `skills`: Canonical skill catalog (`<ID>` → { desc, relevance }).  
2. `responsibilities`: Canonical duties (`R#` → description).  
3. `roles`: Each persona with `description`, `responsibilities` (array of IDs), and tiered `skills`.  

To display a role:
1. Look up its responsibility IDs and render descriptions.  
2. Render Basic → Intermediate → Advanced skill sets in order.  
3. (Optional) Derive coverage: For any responsibility requiring a skill outside the role’s lists, flag a potential gap (this framework intentionally trimmed such mismatches in the “cleaner ownership” pass).  

---

## Cleaner Ownership Principles Applied

1. **One primary executor per operational domain** (e.g., secret rotation → Platform, supply chain gates → DevSecOps).  
2. **Strategy distinct from execution** (Architect sets multi-cluster & identity patterns; Platform / Network implement).  
3. **Developers focus on secure creation** (no governance overhead; advanced culture coaching removed).  
4. **DevSecOps ≠ Platform**: DevSecOps owns pipeline, policy gating, and feedback loops—not day-to-day node or quota operations.  
5. **Network role optional**: Only retained where segmentation + trust fabric complexity justifies specialization.

---

## Typical Usage Scenarios

| Scenario | How to Use the Framework |
|----------|--------------------------|
| New customer/platform assessment | Map current staff to roles → highlight missing responsibilities |
| Training plan creation | For each role list missing intermediate skills → build targeted workshops |
| Tool justification (e.g. runtime security) | Show which responsibilities require tuning (R33, R25) and which role owns them |
| Scope creep control | Reject additions that assign execution + governance to the same role |
| Hiring | Translate unowned responsibilities into a candidate profile |

---

## Extending the Framework

1. Add a new responsibility in `responsibilities` with next free ID (avoid reusing IDs).  
2. Only then assign it to exactly one primary role (add secondary later if absolutely needed).  
3. If the responsibility implies a capability not described, add a new skill ID (choose tier by impact).  
4. Re-run a quick consistency check manually: every role responsibility should map conceptually to at least one of its skill tiers.  
5. Keep Advanced additions scarce to preserve clear progression.

Naming guidance:
- Responsibilities = Outcome phrased verb-first, single focus.  
- Skills = Capability statements (what the person *can do*), with relevance kept vendor-neutral.  

---

## Example: Gap Analysis Walkthrough

Goal: “Do we have coverage for workload identity (R36)?”  
1. Find R36 in responsibilities → definition.  
2. Roles owning it: Security Architect, Network & Infrastructure Engineer.  
3. If neither role exists in your org, decide: elevate Platform Operator (add R36 + A10) or create specialization.  
4. Check skills: A10 present only where identity design is expected.  
Outcome: Clear staffing or enablement decision.

---

## Why Some Advanced Skills Are Absent in Certain Roles

When an advanced capability is missing from a role, that is intentional—it protects focus and keeps cognitive load manageable. For example, application developers absolutely contribute to producing secure, minimal, reproducible container images and to adopting supply‑chain good practices inside their build process. However, they are not expected to design or operate the governance layer for software supply chain integrity (such as signing infrastructure, attestation validation, or organization‑wide policy gates). Those activities sit with the DevSecOps function, which owns the operation and enforcement of build and deployment security controls. Developers consume the guardrails; DevSecOps engineers implement, integrate, and continually refine them.


---

## Feedback & Contribution

Open issues or PRs with:
- Proposed new responsibility (include justification + proposed primary owner)
- Skill clarifications or tier adjustments
- Real-world adoption feedback (what resonated / what was confusing)

Your input helps refine this into an industry-aligned reference rather than a single‑perspective artifact.

---

## Attribution / Disclaimer
This framework is opinionated but grounded in recurring patterns observed in platform security, DevSecOps adoption programs, and container security enablement across enterprise OpenShift / Kubernetes environments. Adapt locally where your org structure or regulatory context demands.
