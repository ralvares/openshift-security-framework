<!-- ============================================================= -->
# Kubernetes & OpenShift Security Field Playbook (Layer 1)
Fast, opinionated guide for day-to-day execution. Read this layer to act; see Layer 2 Annex for all framework IDs & detailed evidence crosswalk.

### Minimal Lifecycle (Mental Model)
Build → Deploy → Run

| Stage | Verbs | Goal |
|-------|-------|------|
| Build | Sign · Scan · Harden | Produce trusted, minimal artifacts |
| Deploy | Gate · Segment · Pin | Let in only what meets bars |
| Run | Detect · Contain · Prove | Spot drift, respond, preserve evidence |

### Core Security Questions (Answer “Yes” Without Hesitation)
Below: 9 core questions. For each: checklist, a Quick Win, and a Reality Check.

#### 1. Are all images trusted & verifiable? (→ Annex)
Checklist:
- [ ] All running images scanned in last 24h
- [ ] Digests (no mutable :latest)
- [ ] Signatures verified at admission (pilot → broad)
- [ ] Unapproved registries blocked
Quick Win: Enforce digest pinning + block unknown registries.
Reality Check: Unsigned / unscanned image = silent Trojan delivery path.

#### 2. Are risky configs blocked before they reach the cluster? (→ Annex)
Checklist:
- [ ] Privileged, host mounts, run as root all blocked (not just warned)
- [ ] Resource limits required
- [ ] Non‑root + read‑only FS (where feasible) policies active
- [ ] Policy set versioned (Git, signed)
Quick Win: Flip “Privileged Container” policy from warn → block.
Reality Check: One lingering privileged pod multiplies blast radius.

#### 3. Is RBAC least‑privilege and stable? (→ Annex)
Checklist:
- [ ] Single cluster‑admin group (no direct users)
- [ ] Quarterly RBAC diff reviewed & signed
- [ ] Service accounts scoped (no wildcards *)
- [ ] Escalation verbs (bind, escalate) restricted
Quick Win: Remove all direct cluster-admin user bindings—leave only group.
Reality Check: Excess admin bindings turn any stolen token into cluster compromise.

#### 4. Is east‑west traffic segmented by default? (→ Annex)
Checklist:
- [ ] Every namespace has deny-all ingress + egress baseline
- [ ] Allow policies reviewed & stored in Git
- [ ] Coverage metric = 100% workloads with both ingress & egress policies
- [ ] Secondary networks (Multus/UDN) inventoried & owned
Quick Win: Apply namespace deny-all baseline manifests cluster-wide.
Reality Check: No baseline = attacker lateral exploration playground.

#### 5. Are fixable vulnerabilities gated & remediated on time? (→ Annex)
Checklist:
- [ ] Block fixable Critical CVEs at deploy
- [ ] High CVE block date scheduled
- [ ] SLA: Critical ≤7d, High ≤30d tracked
- [ ] Report shows ≥95% Critical in SLA
Quick Win: Turn on “block fixable Critical” policy today; log first denied deploy.
Reality Check: Aging fixable CVE = accepted, compounding risk.

#### 6. Will runtime anomalies be caught & contained? (→ Annex)
Checklist:
- [ ] 100% sensor / collector node coverage
- [ ] At least 2 high‑confidence runtime policies enforced (kill/scale or block)
- [ ] Notifier path (→ ticket / SIEM) tested
- [ ] Alert precision measured (noise <20%)
Quick Win: Enable crypto‑miner & reverse shell runtime policies with notifier.
Reality Check: No runtime containment = dwell time escalates silently.

#### 7. Are secrets handled safely (no plain leaks)? (→ Annex)
Checklist:
- [ ] Secret-in-env pattern detection on
- [ ] Zero new secret violations week over week
- [ ] High-risk secrets sourced from vault (not static YAML)
- [ ] Rotation evidence available (external)
Quick Win: Enable secret‑in‑env policy + remediate top offenders.
Reality Check: Hard‑coded credential = instant lateral foothold.

#### 8. Is evidence tamper‑evident & exceptions under control? (→ Annex)
Checklist:
- [ ] Nightly compliance/policy export + SHA‑256 digest
- [ ] Hash index stored in WORM/object‑lock
- [ ] Exception register empty or all entries < expiry
- [ ] Admission webhook error rate monitored
Quick Win: Automate nightly export + digest + upload to object‑lock bucket.
Reality Check: Missing immutable evidence invites audit challenge & gap claims.

#### 9. Will noisy neighbors or runaway pods be stopped? (→ Annex)
Checklist:
- [ ] CPU & memory limits on all workloads (policy enforced)
- [ ] Namespace quotas sized & reviewed
- [ ] Alert on “missing limits” older than 24h
- [ ] Capacity vs quota variance tracked
Quick Win: Enforce “missing resource limits” policy cluster‑wide.
Reality Check: Unbounded pod can starve critical services → cascading outage.

### Quick Wins (Top 10 One-Liners)
| # | Action | Why |
|---|--------|-----|
| 1 | Enforce privileged container block | Removes largest lateral escalation vector |
| 2 | Enforce digest pinning + ban :latest | Stops image swap & rollback ambiguity |
| 3 | Block fixable Critical CVEs | Immediate reduction of exploitable risk |
| 4 | Apply namespace deny-all (ingress+egress) | Establishes segmentation floor |
| 5 | Remove extra cluster-admin subjects | Shrinks blast radius of credential theft |
| 6 | Enable runtime crypto‑miner & reverse shell policies | Catches common early-stage abuse |
| 7 | Turn on secret-in-env detection | Ends accidental credential leakage |
| 8 | Automate nightly compliance export + hash | Creates tamper-evident audit trail |
| 9 | Require resource limits policy | Prevents noisy neighbor exhaustion |
| 10 | Add exception register with expiry alert | Stops silent indefinite risk acceptance |

### Minimal Daily Triage Loop
1. Check vulnerability SLA dashboard (Critical aging / breaches).
2. Review new blocked deploys (why blocked? legitimate? fix or exception?).
3. Inspect runtime high-severity alerts (noise vs signal; tune if noisy).
4. Confirm network policy coverage = 100% (no regressions).
5. Scan exception register for upcoming / expired items.

---
# ---
# Control Mapping Annex (Layer 2)
Authoritative reference: framework control IDs, responsibility designations, and evidence artifacts. No new prescriptive actions here—only mappings for audit & assurance.

## A. Consolidated Mapping Overview
| Domain (Core Question) | CIS (v8) | NIST 800-53 (Primary) | NIST 800-190 | DORA (Indicative*) | Responsibility (OCP / RHACS / External) | Evidence Examples |
|------------------------|----------|-----------------------|--------------|--------------------|------------------------------------------|------------------|
| Image Trust & Provenance (Q1) | 2.1, 2.2, 2.3, 2.5, 11.1 | CM-2, CM-6, CM-7, RA-5 (scan), SI-7, SR-11 | 4.1.1–4.1.4, 4.1.10–4.1.12 | Art. 6(2) ICT risk; Art. 8 monitoring | OCP:P (admission sign verify), RHACS:C (scan/gate), External:E (SBOM gen, key custody) | Policy JSON, blocked unsigned log, SBOM + digest, signer key SOP |
| Baseline Config & Hardening (Q2) | 4.1, 4.6, 11.1 | CM-2, CM-3, CM-6, CM-7 | 4.1.3, 4.1.8–4.1.9, 4.1.13, 4.2.7 | Art. 6 governance controls | OCP:C (SCC/PodSec), RHACS:P (misconfig detect), External:E (hardening std) | MachineConfig diff, misconfig policy export, drift trend |
| RBAC & Least Privilege (Q3) | 4.8, 5.1, 3.3 | AC-2, AC-3, AC-6, CM-5 | 4.2.1, 4.2.4 | Art. 6 access mgmt | OCP:C (RBAC), RHACS:P (visibility), External:E (IAM lifecycle) | RBAC diff signed, privilege anomaly report |
| Network Segmentation (Q4) | 12.1, 12.4 | SC-7(+3/4), AC-3 | 4.2.2 | Art. 11 (resilience), Art. 12 (testing) | OCP:C (NetworkPolicy), RHACS:P (coverage/flows), External:E (L7/mTLS governance) | Deny-all baseline manifests, coverage %, flow graph |
| Vulnerability Gating & SLA (Q5) | 2.2, 6.2 | RA-5, SI-2 | 4.1.4, 4.4.2 | Art. 6(6) continuous monitoring | OCP:P (digest pin), RHACS:C (scan/block), External:E (rebuild process) | Blocked CVE event, SLA report, rebuild pipeline log |
| Runtime Threat Detection (Q6) | 2.6, 2.7, 8.2, 17.2 | SI-4, IR-4(5), IR-5, IR-6(1), AU-12 | 4.5.1, 4.5.2 | Art. 17 incident handling | OCP:P (audit/log surface), RHACS:C (runtime detect/contain), External:E (IR runbooks) | Runtime alert → ticket, policy action log |
| Secrets Handling (Q7) | 13.1, 3.3 | SI-7, SI-7(1), (SC-28 ext) | 4.1.7, 4.2.3 | Art. 6 data security | OCP:P (secret storage), RHACS:P (pattern detect), External:E (vault, rotation) | Secret violation trend, vault rotation report |
| Evidence & Exceptions (Q8) | 7.1, 7.3, 8.2, 16.12, 10.4 | AU-6, AU-12, AU-9(ext), IR-6(1) | 4.2.6 | Art. 8 monitoring & reporting | OCP:P (audit emit), RHACS:P (export), External:E (WORM retention, correlation) | Export + hash chain, exception register, SIEM ingestion dashboard |
| Resource Governance (Q9) | 4.6 | SC-6, CM-7 | 4.2.5 | Art. 11 operational resilience | OCP:C (quotas/limits), RHACS:P (missing limits detect), External:E (capacity planning) | Quota manifests, limits violation report |

*Indicative DORA references: high-level alignment only; consult official DORA Articles for authoritative scope (e.g., Articles 5–15 ICT risk management & resilience, 17 incident reporting, 21 testing).

## B. Evidence Cheat Table (What to Show Fast)
| Domain | Minimum Pair (Platform+RHACS) | External Companion |
|--------|-------------------------------|--------------------|
| Image Trust | Blocked unsigned deploy log + policy JSON | SBOM file + signer key SOP |
| Baseline Config | MachineConfig/SCC export + misconfig policy trend | Hardening standard change record |
| RBAC | Signed RBAC diff + privileged binding removal ticket | IAM access recertification report |
| Segmentation | Deny-all manifests + coverage % graph | Mesh mTLS / firewall ACL set |
| Vulnerabilities | Blocked Critical CVE event + SLA dashboard | Rebuild pipeline log (new digest) |
| Runtime Threats | Runtime alert → ticket + auto action log | IR runbook reference & closure notes |
| Secrets | Secret violation trend + policy export | Vault rotation/lease report |
| Evidence & Exceptions | Compliance export + hash index | Object-lock config / retention policy |
| Resource Governance | Limits policy violation trend + quota manifests | Capacity forecast vs actual |

## C. Responsibility Legend
| Mark | Meaning |
|------|---------|
| C | Fully enforced/evidenced in that layer |
| P | Partial contribution (shared responsibility or evidentiary assist) |
| E | Entirely external (tracked in External Control Register) |

## D. Per-Domain Detail (Concise Crosswalk)
| Domain | Intent | Core Action (Do This) | Minimal Evidence | Framework IDs (CIS / 800-53 / 800-190 / DORA) | Risk if Weak |
|--------|--------|-----------------------|------------------|-----------------------------------------------|--------------|
| Image Trust & Provenance (Q1) | Only trusted, signed, scanned, pinned images | Enforce digest pin + signature verify + block unapproved registries | Policy export; blocked unsigned log; SBOM + external digest | CIS 2.1/2.2/2.3/2.5/11.1; CM-2/6/7, RA-5, SI-7, SR-11; 4.1.1–4.1.4, 4.1.10–12; DORA Art.6(2) | Supply chain tampering unnoticed |
| Baseline Config & Hardening (Q2) | Declarative least-function posture | Block privileged/root/host mounts; require limits | MachineConfig diff; misconfig trend; signed policy commit | CIS 4.1/4.6/11.1; CM-2/3/6/7; 4.1.3, 4.1.8–9, 4.1.13, 4.2.7; DORA Art.6 | Privilege creep & attack surface growth |
| RBAC & Least Privilege (Q3) | Minimal standing privilege | Single cluster-admin group; quarterly RBAC diff | Signed RBAC diff; scope report; escalation violation | CIS 4.8/5.1/3.3; AC-2/3/6, CM-5; 4.2.1, 4.2.4; DORA Art.6 | Lateral movement via overbroad rights |
| Network Segmentation (Q4) | Deny-by-default east-west | Namespace deny-all + curated allowlist + coverage metric | Deny-all manifests; coverage %; flow graph | CIS 12.1/12.4; SC-7(+3/4), AC-3; 4.2.2; DORA Art.11/12 | Rapid lateral propagation |
| Vulnerability Gating & SLA (Q5) | Timely fix of exploitable risk | Block fixable Critical; schedule High block; track SLA | Blocked deploy log; SLA dashboard; rebuild digest diff | CIS 2.2/6.2; RA-5, SI-2; 4.1.4, 4.4.2; DORA Art.6(6) | Accumulating exploitable backlog |
| Runtime Threat Detection (Q6) | Detect & contain anomalies | Enforce high-signal runtime policies + ensure 100% sensor coverage | Alert→ticket; action log; sensor coverage report | CIS 2.6/2.7/8.2/17.2; SI-4, IR-4(5), IR-5, IR-6(1), AU-12; 4.5.1–4.5.2; DORA Art.17 | Silent post-compromise activity |
| Secrets Handling (Q7) | Prevent secret leakage | Enable secret pattern detection + migrate to vault | Secret violation trend; policy export; vault rotation report | CIS 13.1/3.3; SI-7, SI-7(1), (SC-28 ext); 4.1.7, 4.2.3; DORA Art.6 | Credential exposure → escalation |
| Evidence & Exceptions (Q8) | Tamper-evident operations record | Nightly export + hash index + active exception register | Export + hash chain; exception register; SIEM ingest metrics | CIS 7.1/7.3/8.2/16.12/10.4; AU-6/12/9(ext), IR-6(1); 4.2.6; DORA Art.8 | Inability to prove operation |
| Resource Governance (Q9) | Prevent noisy neighbor exhaustion | Enforce resource limit policy + quotas & monitor variance | Quota manifests; limits violation report; usage vs quota | CIS 4.6; SC-6, CM-7; 4.2.5; DORA Art.11 | Resource exhaustion & instability |

Abbrev: 800-190 sections numeric; DORA references indicative. Responsibility & expanded evidence examples appear in earlier Annex tables; this table is the normalized crosswalk.

---
# (Existing Detailed Appendices Retained Below)

---
## 2. Build Trust & Baseline
Goal: Produce minimal, attestable, signed, scanned artifacts plus hardened declarative configuration & least privilege before anything is considered for deployment.

Focus Areas (merged from prior Image, Baseline Config, RBAC, part Resource):
- Trusted Sources: Approved registries only; unknown registry = block.
- Integrity & Authenticity: Sign images (progressively) and verify at admission; pin by digest (ban :latest).
- Minimal Surface: Enforce non-root, no dangerous capabilities, read-only FS where feasible, resource limits set.
- RBAC Hygiene: Single cluster-admin group; service accounts least privilege; quarterly diff review.
- Policy-as-Code: All security policies stored & versioned (Git) with signed commits.

Key Actions (Condensed):
1. Inventory registries → configure allowlist + block policy.
2. Turn on signature verification (pilot namespace → expand).
3. Enforce digest pinning & ban :latest; add admission check.
4. Enable core misconfiguration policies (privileged, host mounts, run as root, missing limits) in warn → block progression.
5. Quarterly RBAC diff & service account scope review.
6. Store SBOM + policy exports with external SHA-256 digest (generation external to RHACS).

Evidence Bundle (example): policy export (signed commit), blocked unsigned image log, RBAC diff approval, SBOM digest record.
Risk if skipped: Untrusted or mutable images slip in; privilege creep; baseline drift; unverifiable provenance.

---
## 3. Gate & Segment
Goal: Enforce only what meets security bars; tightly constrain ingress/egress & lateral movement at deploy time.

Components:
- Admission Gating: Block fixable Critical CVEs; warn then block High later.
- Network Baseline: Deny-all ingress + deny-all egress per namespace; explicit allow rules only.
- Coverage Analytics: Track % workloads with both ingress & egress policies (target 100%).
- Shadow Paths: Identify Multus / UDN networks; ensure each has an owner (External Register row) or treat as gap.
- Classification (Optional Advanced): Label & taint nodes by data sensitivity; enforce scheduling + stricter policies for higher zones.

Key Actions:
1. Apply namespace deny baseline manifests (kept in Git).
2. Generate candidate allow policies from observed flows → tighten → review → apply.
3. Monitor for missing policy coverage or sudden new connections.
4. Enforce privileged/host mount/hostNetwork denials.
5. For secondary networks/UDN: add owner + review cadence in External Register.

Evidence Bundle: NetworkPolicy set + coverage metric, sample blocked privileged deploy, segmentation drift alert, exception (if any) with expiry.

---
## 4. Runtime & Secrets
Goal: Detect & contain unexpected runtime behavior and residual secret exposures.

Runtime Detection:
- Baselines processes & connections; flags new/abnormal.
- High-confidence patterns (crypto miner, reverse shell hints, package installs at runtime).
- Optionally auto-contain (kill/scale) for specific policies.
Secrets Safety Net:
- Detect obvious secret patterns in env/config (heuristic only) while migrating to vault-backed references.
- Track trend of secret violations → aim for zero new after initial cleanup.

Resource Governance:
- Enforce CPU/memory limits & quotas to prevent noisy neighbor amplification of incidents.

Key Actions:
1. Enable top runtime policies + notifier integration test.
2. Validate 100% collector coverage (no blind nodes).
3. Set alert precision target; tune noise quarterly.
4. Enable secret-in-env detection; remediate high-risk hits first.
5. Enforce resource limit policy; alert on unbounded workloads.

Evidence Bundle: Runtime alert → ticket, collector coverage report, secret violation trend chart, resource quota vs usage export.

---
## 5. Vulnerability Remediation (Lifecycle Deep Dive)
Goal: Keep fixable risk within SLA and prove rebuild over ignore.

Lifecycle Pattern:
1. Discover (scan on build + rescan running images <24h age).
2. Gate (block new deploys with fixable Critical; plan High block date).
3. Rebuild (pipeline rebuild on base image update; produce new digest).
4. Verify (policy re-evaluates; violation clears).
5. Report (daily summary; SLA dashboard; external hash for report integrity).

Core Metrics: % Critical in SLA (≥95%), median time to remediate Critical (downward trend), count of exceptions & their age.

Evidence Bundle: Blocked deploy log, vulnerability summary (hashed), rebuilt image digest comparison, exception register entries (if any) with expiry.

---
## 6. Evidence & Governance
Goal: Immutable, trustworthy audit trail + disciplined exception management + explicit external ownership.

Foundational Practices:
- Nightly compliance & policy export → compute SHA-256 → store in WORM/object-lock.
- Forward runtime & deploy alerts to SIEM; alert on pipeline delivery failures.
- Hash Chain / Index: Maintain chronological index of evidence digests (tamper-evidence). Immutability handled by external storage controls.
- Exception Register: ID, rationale, risk rating, approver, creation & expiry; zero overdue.
- External Control Register: Host hardening, IAM/MFA, key management, secret rotation, log retention, WAF, backup/DR, etc. Each has owner + review cadence.

Failure Modes (Sample): admission webhook timeout (treat as deny for critical policies), unscanned image aging, missing collector node, notifier delivery failure, missing daily export. Monitor & alert accordingly.
Evidence Bundle: Export hash index snippet, SIEM ingestion screenshot, exception register report (no overdue), external register snapshot with statuses.

---

# Kubernetes & Container Security Compliance Guide (End-User Focus)

Practical guidance for structuring and evidencing container/Kubernetes security controls with **Red Hat Advanced Cluster Security (RHACS / StackRox)** and **OpenShift**. This is *enablement material* (not a formal attestation) and must be paired with organizational policies, procedures, and broader platform controls.

---
## 0. How to Use & Scope
This guide normalizes overlapping framework language into actionable security “themes”. For each theme you get: intent, risk, platform + RHACS capabilities, key actions, and incremental evidence. Use the quick reference + appendices to translate into specific control IDs.

### 0.1 Lifecycle Alignment: Build → Deploy → Run → Detect & Respond (NIST-Inspired)
If the 9 Theme model feels heavy, you can operate day-to-day with a streamlined lifecycle lens. These phases loosely map to NIST CSF Core Functions (Identify/Protect/Detect/Respond) while preserving all underlying controls. Use this as your “quick mental model,” then dive into Theme sections when you need depth.

| Lifecycle Phase | NIST CSF Emphasis (Primary) | Source Themes (Detail) | Core Questions | Representative Actions (Examples) | Key Evidence (Minimal Bundle) |
|-----------------|-----------------------------|------------------------|----------------|-----------------------------------|--------------------------------|
| Build | Identify / Protect | 1 (image trust), 2 (baseline config), 3 (RBAC), 5 (resource defaults) | Is the artifact trustworthy & minimal? | Enforce signed & scanned images; codify policy-as-code; least-priv RBAC; set default limits | Policy repo commit (signed); image provenance scan report; RBAC diff; quota/limit manifest |
| Deploy (Gate) | Protect | 1,2,4 (initial segmentation), 5,6 (CVE gate), 8 (secret leak check) | Should this be allowed in? | Block fixable Critical CVEs; deny privileged/host mounts; digest pinning; namespace deny-all base | Blocked deployment log; admission policy export; network deny base YAML; secret violation trend |
| Run (Operate) | Protect / Detect | 4 (segmentation maintenance), 5 (resource fairness), 7 (runtime), 8 (late secrets) | Is behavior within approved boundaries? | Maintain 100% NetworkPolicy & sensor coverage; observe resource anomalies; prune stale privileges | NetworkPolicy coverage metric; runtime sensor coverage; resource quota vs usage report |
| Detect & Respond | Detect / Respond | 6 (aging vuln deltas), 7 (runtime alerts/actions), 9 (logging & evidence), Exceptions Register | Did something unexpected happen & did we act? | Runtime anomaly → ticket; SLA vulnerability aging; automated response (kill/scale); exception expiry enforcement | Runtime alert + ticket link; vuln SLA dashboard; response action log; nightly compliance export + hash |

Condensed “Top 10 Moves” (mapped to phases):
1. Enforce digest pinning & ban `:latest` (Build/Deploy)
2. Require signatures (progressively) for high-risk namespaces (Build/Deploy)
3. Block privileged & host mount misconfigs (Deploy)
4. Deny-by-default ingress & egress (Deploy → Run)
5. Block deploys with fixable Critical CVEs (Deploy)
6. Reduce cluster-admin to one group + quarterly diff (Build)
7. Achieve & monitor 100% runtime sensor/node coverage (Run / Detect)
8. Secret-in-env detection + migration to vault (Deploy / Run)
9. Nightly compliance export + external SHA-256 + WORM/object lock (Detect & Respond)
10. Exception register with expiry + alert on overdue (Detect & Respond)

“5 Fast Health Metrics” (executive snapshot):
| Metric | Phase Signal | Healthy Target | Why It Matters |
|--------|--------------|---------------|----------------|
| % Running images scanned <24h | Build/Deploy hygiene | ≥98% | Fresh vuln intel pre & post deploy |
| Fixable Critical CVEs in SLA | Deploy/Detect | ≥95% in SLA | Prevent aging high risk |
| Namespace ingress+egress baseline coverage | Deploy/Run | 100% | Lateral movement constraint |
| Runtime sensor coverage (nodes) | Run/Detect | 100% | No blind operational zones |
| Overdue exceptions | Detect & Respond | 0 | Governance integrity |

Lifecycle ↔ Theme ↔ Framework Reference:
- Build: NIST families CM, AC (scoping), SR (authenticity), partial SI (integrity); CIS 2.x (software inventory), 4.1/4.6, 11.1.
- Deploy: NIST CM, SC-7 (deny baseline), RA-5/SI-2 (gating risk), SI-7 (signature), AC-3 enforcement; CIS 2.3, 2.5, 4.6, 12.1, 12.4, 6.2, 13.1.
- Run: NIST SC-7 continuous segmentation, SC-6 resource, SI-4 monitoring; CIS 12.1/12.4 (ongoing), 4.6 least functionality, 2.6/2.7 runtime allowlist facets.
- Detect & Respond: NIST SI-4, IR-4(5), IR-5, IR-6(1), AU-6/AU-12, RA-5 trend, AU-9 (external retention); CIS 7.1, 7.3, 8.2, 16.12, 10.4, 6.2 remediation evidence.

Use this lifecycle table for rapid onboarding & stakeholder decks; keep the deeper Theme sections authoritative for audits and appendices crosswalks. When updating controls, modify the Theme section first—then adjust the lifecycle summary if the change is material (>policy shift or new enforcement capability).


> Coverage Model Clarification: The scope combines (a) OpenShift / RHCOS platform primitives ("OCP" – SCC/Pod Security, RBAC, NetworkPolicy, MachineConfig/OSTree, ClusterImagePolicy & signature admission, Compliance & Security Profiles Operators, ingress TLS, optional mesh mTLS) and (b) Red Hat Advanced Cluster Security ("RHACS") overlay capabilities (image & component scanning, deploy misconfig & vuln gating, runtime anomaly detection, secret pattern detection, policy & compliance evidence exports). Tables now show three columns (OCP | RHACS | External). A blank cell means “no substantive contribution.” The External column lists items outside the combined in‑scope boundary that must be evidenced via Appendix E (External Control Register).

### Covered vs External Responsibilities
| Category | RHACS Primary | RHACS Partial (Evidence Component) | External / Platform (Document Separately) |
|----------|---------------|------------------------------------|-------------------------------------------|
| Image / Supply Chain Policy | Scan, policy gate, risk scoring | SBOM association, signature policy tie‑in | Signing infra, build provenance chain (Cosign/Sigstore, pipeline attestation) |
| Runtime Detection & Response | Process / network anomaly, policy enforcement | Alert forwarding / correlation | Full SIEM correlation, SOAR workflows, WAF/RASP |
| Vulnerability Management | Prioritization, fixable metrics, gating | SLA tracking exports | Patch orchestration, inventory governance (CMDB) |
| Access / RBAC Hygiene | Visibility, cluster-admin minimization check | Mapping service accounts to namespace scope | Enterprise IAM (MFA, SSO session controls, PAM) |
| Network Segmentation | NetworkPolicy coverage analytics | Flow visualization for validation | East/West deep inspection, service mesh mTLS policy authority |
| Secrets Exposure | Secret-in-env detection | Partial detection of embedded credentials | Enterprise vault, key lifecycle management |
| Logging / Evidence | Policy + violation export, compliance summaries | Supplemental security event stream | Immutable log store, retention, anti‑tamper, TRA artifacts |

*Clarification (SBOM & Signature Scope – concise):* With a Signature Integration, RHACS verifies Cosign signatures (public key / cert / keyless) and, if enabled, Rekor transparency log inclusion on discovery and periodic (~4h) re-checks, and can block unverified images ("Not verified by trusted image signers"). It does **not** generate or sign SBOMs, manage long‑term keys / Fulcio roots, or build full SLSA / in‑toto provenance beyond the configured signature + optional Rekor check. Pipelines (RHTAP + RHTAS) supply SBOM + attestation + key lifecycle; RHACS enforces signed & pinned images and exports verification evidence.

### Baseline Evidence Pattern (Applies Unless Theme Lists “Additional Evidence”)
Unless a theme explicitly lists “Additional Evidence”, capture this baseline set:
1. Current (date-stamped) compliance report excerpt for relevant checks
2. Policy configuration snapshot (JSON export or signed Git version)
3. Sample (sanitized) prior violation + its remediation closure proof
4. External SIEM log / ticket reference linking alert → action
5. Change history (Git / ticket ID) for material control adjustments

### Enforcement Failure Modes & Resilience (Consolidated)
Understanding how controls behave under partial failure prevents silent coverage gaps. The table below summarizes typical failure / degradation scenarios and recommended guardrails.

| Stage / Function | Component(s) | Potential Failure / Condition | Typical Default Behavior | Risk | Recommended Control / Setting | Monitoring Signal |
|------------------|-------------|--------------------------------|---------------------------|------|-------------------------------|-------------------|
| Build / CI Policy Evaluation | CI plugin / API call to Central | Central/API unreachable, latency, auth error | Pipeline step fails (hard fail) or is skipped (if not enforced) | Unsanctioned image proceeds | Fail pipeline on evaluation error (treat as deny) | CI job logs; alert on evaluation error count >0 |
| Image Scanning | Scanner deployment / external scanner integration | Scanner pod crash, version drift, registry credential failure | Image marked unscanned / stale scan retained | Blind to new CVEs | Alert if any deployed image lacks a scan < N hours old | Compliance “unscanned images” delta |
| Admission Enforcement (Deploy Stage) | RHACS validating webhook + OpenShift admission chain | Webhook timeout / DNS / cert expiry | Kubernetes fail-open by default unless failurePolicy=Fail | Risky deploy allowed | Set failurePolicy=Fail for critical policies (e.g., unsigned image, critical CVE) | Admission controller error rate metric |
| Admission Ordering | Multiple controllers (SCC, PodSecurity, Gatekeeper/Kyverno, RHACS) | Conflicting deny reasons / mutation ordering | Inconsistent error surfaced to user | Mis-triage & bypass attempts | Define ownership matrix; avoid duplicate rules across controllers | Admission audit logs; compare rejection sources |
| Runtime Collection | Sensor / Collector DaemonSet | Pod eviction / version mismatch / network partition | Gaps in runtime events; policies still appear “configured” | Undetected runtime anomaly / incident | Monitor collector heartbeat & connected nodes; alert on <100% coverage | Heartbeat metric; node coverage dashboard |
| Notifier Delivery | Slack / PagerDuty / SIEM forwarding | Credential rotation, endpoint outage | Alerts buffered or dropped; silent failure | Delayed response, lost evidence | Health check synthetic alert daily; alert on notifier error backlog | Notifier failure counter |
| Vulnerability Export / Reports | Scheduled job / API script | Script error / auth token expired | Missing daily evidence artifact | Evidence gap for audit period | Compute & store external SHA-256 hash chain of daily exports; alert on missing date in sequence | Object store listing gap detection |
| Logging Pipeline | Forwarder / SIEM ingestion | Buffer full, parse errors, TLS failure | Partial ingest; some events lost | Incomplete forensic trail | Enable pipeline failure alerting & health checks | SIEM ingestion error dashboard |
| Policy Exception Expiry | Exception register | Exception passes expiry unnoticed | Control gap persists | Risk acceptance indefinite | Automated job flags exceptions past due date | Daily exception aging report |

> Principle: Treat *inability to evaluate* the same as *deny* for critical gates; fail-closed where business impact is acceptable, fail-open only with compensating detection and explicit, time-bound exception.

### Capability Boundaries & Disclaimers
RHACS provides *container & Kubernetes workload* focused security evidence. It does **not**:
- Replace host OS / kernel hardening (CIS benchmark, kernel module integrity, eBPF constraint outside sensor scope)
- Enforce file-level FIM (File Integrity Monitoring) for host paths
- Provide full WAF, RASP, or API schema validation (pair with ingress/WAF layer)
- Manage IAM / MFA / SSO life-cycle (external IdP / IAM system authoritative)
- Offer data-at-rest encryption or key lifecycle management (delegate to platform KMS / vault)
- Guarantee retention / immutability (external WORM/object-lock store required)

Where such controls appear in mapping tables, RHACS contributes *partial* (P) evidence only (telemetry, detection triggers, or gating context) and relies on external systems for full compliance.

### Enforcement Modes & Fallback Nuances
Add these context notes anywhere you operationalize policies so expectations stay realistic (see also the Enforcement Failure Modes & Resilience table above for consolidated behavior references):
- **Progressive Enforcement:** Most teams start high-impact policies (privileged container, unsigned image, critical CVE) in “alert / warn (no block)” mode for 1–2 sprints before switching to enforce. Document the promotion decision (date, risk rationale) for audit.
- **Admission vs Sensor Fallback:** RHACS admission webhook can time out (network, DNS, cert expiry). If `failurePolicy=Ignore` (fail-open) the deploy may proceed; sensor-side (“deploy-time”) enforcement may still catch some conditions but timing differs. Critical policies should usually set `failurePolicy=Fail` + alert on webhook error rate.
- **Hard vs Soft Actions:** Some runtime policies use “scale-to-zero” or alert-only actions—these are *soft* responses compared to an admission block. Mark each policy with its action class in your policy register.
- **Break-Glass / Exceptions:** Temporary allowance (e.g., adding a dangerous capability) must reference an exception ID and expiry. Avoid ad‑hoc manual toggles; prefer Git-based changes.
- **Latency & Race Windows:** A very rapid deploy after image push may momentarily lack latest scan results; mitigate by scanning in CI (pre-push) and failing build on unacceptable issues.
- **Version / Feature Gating:** Some advanced enforcement (e.g., signature policy integration, specific runtime kill actions) depends on cluster + RHACS version parity; annotate policies with minimum version where relevant.

> Tag policy YAML (or JSON export) with labels: `enforcementPhase=warn|block`, `criticality=high|medium|low`, `failurePolicy=Fail|Ignore` for clarity.

### Multi-Controller Policy Interplay (Brief Note)
If you run multiple admission / policy controllers (e.g., SCC / PodSecurity, RHACS admission webhook, Gatekeeper/Kyverno, Sigstore verification), document for each enforced rule WHICH controller is authoritative. Undocumented overlap creates:
- Confusing or duplicated deny messages
- Race conditions / ordering ambiguity (different failurePolicy settings)
- Hidden gaps (assumed “other controller” covers it)

Minimal recommended doc set (kept in Git): controller inventory, per-control owner, failurePolicy, deny message prefix convention, change approval process. Remove redundant enforcement—prefer one authoritative source and treat others as visibility only.

### Policy Bypass / Exception Audit Requirements
All bypasses or temporary relaxations must produce *auditable artifacts*:
1. **Exception Register Entry:** ID, control/policy name, rationale, risk rating, approver, creation date, expiry date.
2. **Mechanism:** Prefer Git-managed policy change (PR includes exception metadata) over ad-hoc UI toggles.
3. **Annotations (Optional):** If using Kubernetes annotations to tag exceptioned workloads (pattern `security.exception/<policy-id>=<exception-id>`), log and export these in inventory reports.
4. **Expiry Enforcement:** Scheduled job checks for past-due exceptions; generate alert + auto-create ticket for review.
5. **Evidence:** Retain original denied manifest (redacted if needed) + post-remediation manifest proving closure.

### External Control Register (Summary)
Controls marked (E) in mappings require explicit tracking. Maintain a register (sample in Appendix E) with: Control Domain, External Owner/System, Evidence Artifact Type, Review Cadence, Last Verified Date. Link each external domain to authoritative documentation (e.g., SOC2 policies, platform hardening guides). This prevents “silent gaps” where a dependency was assumed but never evidenced.

Key External Domains (illustrative): Host OS Hardening, Kernel Patching, Node CIS Benchmark, IAM & MFA, Key Management (KMS/Vault), Data Encryption (at rest / in transit), WAF / API Gateway, SAST/DAST, License Compliance & Legal Review, Backup & DR, Log Retention / WORM, SIEM Correlation Rules, Secrets Lifecycle (rotation), Incident Runbooks, Targeted Risk Analyses, Data Masking / Tokenization.

> Action: Add an *External Owner* column to internal audit prep spreadsheets; absence of a named owner flags a governance risk.

### Reading Order Recommendation
Review Themes 1–9 in sequence for operational rollout; use Section 11 (Control Mapping Quick Reference) for framework control ID cross-references, and Appendices for detailed per-framework translation.

---
## 1. Image Provenance & Supply Chain Integrity (Build & Trust)
**Representative Control Families:** NIST CM‑2 / CM‑6 / (RA‑5 – vuln scanning aspect) / SI‑7 / SR‑11 (component authenticity); NIST 800‑190 4.1.1–4.1.4, 4.1.10–4.1.12 (supply chain / signature / digest pinning) (4.1.10, 4.1.11, 4.1.12 partial; SBOM generation external).  (Note: RA‑5 primary stewardship lives in Theme 6 for remediation SLAs; included here for provenance scanning gate. SI‑7 limited to signature / integrity enforcement portions.)
**CIS (v8) Representative Safeguards:** 2.1, 2.2, 2.3, 2.5, 11.1 (trusted images / authorized software / provenance policy) – see Appendix A for scope & evidence nuances.

### Intent
Assure only approved, scanned, signed, minimal, immutable images from trusted sources reach deploy.

### Risk (If Ignored)
Tampered or vulnerable images deliver exploitable components early; signature gaps reduce provenance confidence; drift invites silent privilege or dependency expansion.

### RHACS Levers
- Image & component scanning; CVE severity & fixability classification
- Policy gating (e.g., disallow unsigned images, disallowed registries, fixable critical CVEs)
- Detection of risky config in manifests (privileged, root user, mutable FS) *before* deploy
- Report assertion: “All deployed images have scanner + registry coverage”

### OpenShift / Platform Levers
- Signature & attestation verification (Cosign/Sigstore admission integration)
- ImageContentSourcePolicy for controlled registry mirrors
- Build pipeline isolation + SBOM generation via Red Hat Trusted Application Pipeline (RHTAP) + GitOps deployment pinning by digest

> Sigstore / ClusterImagePolicy Precondition: Ensure keys, root-of-trust configuration, and any required MachineConfig or operator enablement are completed; signature verification is **not** implicitly active. Document the signing key custody & rotation process.

### Key Actions
1. Inventory registries → integrate all in RHACS; block unknown registries.
2. Enforce “no :latest tag” & digest pinning (policy + manifest review).
3. Require signatures/attestations for high-risk namespaces (progressively roll out).
4. Enable and enforce policies for disallowed critical CVEs & unsigned images.
5. Generate SBOM at build; store artifact + externally computed SHA-256 digest. RHACS does not *generate* SBOMs; external pipeline tooling (e.g., Red Hat Trusted Application Pipeline) should create, sign, and (optionally) hash them. Current RHACS correlation is vulnerability-centric—treat SBOM retention, hashing, and attestation verification as an external control.
6. Track mean time from image build → deploy for provenance freshness metric.

### Additional Evidence
- Signed SBOM artifact (external SHA-256 digest + timestamp)
- Example signature verification admission log (success + rejection)
 - Exception (if any) showing controlled temporary fallback from block→warn with expiry
 - (If applicable) Secure coding pipeline evidence (SAST/DAST report + external SHA-256 digest) for public-facing apps

---
## 2. Baseline Configuration & Drift Control
**Representative Controls:** NIST CM‑2 / CM‑3 / CM‑6 / CM‑7 (least functionality) (CM‑7(1) partial); NIST 800‑190 4.1.3, 4.1.8–4.1.9, 4.1.13; 4.2.7 (admission & enforcement). (CM‑7 surfaced implicitly before—now made explicit.)
**CIS (v8) Representative Safeguards:** 4.1, 4.6, 11.1 (secure configuration baseline, least functionality, policy-as-code) – Appendix A.

### Intent
Codify & continuously enforce hardened deployment settings; surface deviations quickly.

### RHACS Levers
- Deploy-stage misconfiguration policies (privileged, host mounts, escalation, absent limits)
- “Unresolved deploy violations” feed for live drift awareness

### OpenShift Levers
- SCC & Pod Security profiles (restrict privilege + capabilities)
- Admission controllers enforcing resource limits & forbidding host networking

### Key Actions
1. Map internal hardening standard → RHACS policy set (clone, label, commit to Git).
2. Enforce (initially warn, then block) top 5 riskiest misconfigs.
3. Daily triage of new high/critical deploy violations (≤24h closure goal).
4. Quarterly review: prune obsolete custom policies & document rationale changes.

### Additional Evidence
- Drift metrics: count of high-severity misconfigs over trailing 30 days (trend downward)

---
## 3. Least Privilege & RBAC Governance
**Representative Controls:** NIST AC‑2 / AC‑3 / AC‑6 / CM‑5 (AC‑6(1) enhancement implied); NIST 800‑190 4.2.1 / 4.2.4 (access / privilege limit) (partial where external IAM workflows apply).
**CIS (v8) Representative Safeguards:** 4.8, 5.1, 3.3 (admin privilege restriction, service account governance, secret/config access) – Appendix A.

### Intent
Restrict administrative & broad-impact permissions; ensure explicit approvals & periodic review.

### RHACS Levers
- RBAC visualization; detection of multiple cluster‑admin subjects
- Policies detecting privilege escalation vectors at container runtime

### OpenShift Levers
- Granular ClusterRoles + namespace RoleBindings; group-based binding strategy
- SCC layering to enforce default non‑privileged runtime contexts

### Key Actions
1. Consolidate cluster-admin to one group; remove direct USER bindings.
2. Quarterly RBAC diff review (export → compare → sign-off in ticket).
3. Enforce policy on privilege escalation (no additional capabilities, disallow escalate). 
4. Service account scope minimization: restrict * verbs & delete wildcards.

### Additional Evidence
- Signed RBAC diff report (before/after) for quarterly review cycle

---
## 4. Network Segmentation & Boundary Protection
**Representative Controls:** NIST SC‑7 / SC‑7(3) / SC‑7(4) (boundary, deny-by-default, segmentation) + AC‑3 (enforcement tie); NIST 800‑190 4.2.2 (segmentation), 4.4.1 (registry access – partial provenance linkage), Multus/UDN operational governance (external owner). (AC‑3 added for authorization enforcement context.)
**CIS (v8) Representative Safeguards:** 12.1, 12.4 (segmentation & boundary defenses) – Appendix A.

### Intent
Enforce explicit ingress/egress flows; deny-by-default to limit lateral movement.

### RHACS Levers
- Coverage checks: deployments missing ingress and/or egress NetworkPolicies
- Network graph to validate allowed vs observed flows (lateral movement visualization)
- Suggested NetworkPolicy generation from current observed traffic (candidate baseline)
- Post-deployment drift detection: unexpected new connections after baseline

### OpenShift Levers
- NetworkPolicy: Primary L3/L4 segmentation primitive inside the cluster
- Namespaces: Provide administrative scoping only — no isolation unless combined with NetworkPolicy
- Service Mesh (optional): Adds mTLS identity and L7 authorization policies (external to RHACS)
- Multus: Enables secondary network interfaces; traffic on those interfaces bypasses primary cluster NetworkPolicy controls
- User Defined Networks (UDN): Extends OVN-Kubernetes to support multiple logical networks:
	- A UDN may serve as the alternate primary network for a namespace (only one primary) or as a secondary attachment
	- Backends can be:
		- localnet – VLAN-backed segment bridging into physical infrastructure
		- Overlay (L2/L3 VRF) – logical networks isolated from other overlays
		- Routed L3 fabric segment
	- Security stance: treat each UDN as a separate security zone. Maintain an inventory of workloads per UDN, define ACL/firewall policies, and document all cross-UDN flows as explicit “inter-zone” rules.

### Key Actions
1. Apply a deny-all ingress + deny-all egress NetworkPolicy in every namespace.
2. Use RHACS to generate candidate NetworkPolicies from current traffic; review, tighten selectors, store in Git, and only then apply.
3. Simulate coverage and policy changes before enforcing; record approvals as evidence.
4. Flag and document any use of hostNetwork, hostPID, or hostIPC.
5. Weekly: measure % of workloads with both ingress and egress policies (target = 100%).

> Segmentation Clarification & Governance: Namespaces alone do not isolate traffic. True segmentation begins only when NetworkPolicies (or service mesh authz rules) explicitly deny by default and allow required flows. Multus and UDN attachments create parallel paths outside the default pod network—treat each as its own security zone. For every Multus secondary network or UDN logical network, create an External Control Register row (owner, firewall/ACL policy scope, change approval workflow, review cadence). A missing owner constitutes a segmentation compliance gap.

> Policy Generation Caveat: RHACS policies are based on observed traffic. Cold-start or low-traffic services may omit legitimate flows. Stage in warn-only mode, monitor denied traffic alerts, then promote. Periodically regenerate and diff to detect real architecture changes vs anomalous lateral communication.

### Segmentation Scope & Limitations
Kubernetes NetworkPolicies operate at L3/L4 (namespace/pod/port). They do not:
- Inspect payloads or enforce application protocol semantics
- Provide DPI / IDS / IPS capabilities
- Perform data classification / DLP

Use:
- Service Mesh for L7 identity + mTLS authorization
- IDS/IPS or eBPF platforms for deep packet east-west threat detection

Document each extended control (owner + evidence) in the External Control Register (Appendix E).

### Additional Evidence
- NetworkPolicy coverage percentage over time (e.g., last 8 weeks)
- Sample generated NetworkPolicy YAML + review ticket approval + before/after coverage diff
- Example drift detection alert showing unexpected new connection

### Workload Classification & Node Placement ("Compute Zones")
When multiple data sensitivity or regulatory classifications (e.g., Public, Internal, Confidential, Restricted) must coexist on a single cluster, NetworkPolicies alone do not mitigate all residual risks (kernel escape, side-channel, noisy neighbor, forensic contamination). Introduce explicit compute zones that combine node-level segregation, scheduling constraints, and policy enforcement. Treat unapproved co-residency as a violation.

Key Elements:
1. Taxonomy: Publish ordered classification levels with examples + handling rules.
2. Node Segmentation: Label & taint nodes per zone (`classification=restricted`, taint `classification=restricted:NoSchedule`).
3. Scheduling Controls: Require pod label `data-classification=<level>` AND nodeSelector / affinity matching that label; higher classification pods tolerate only their zone taint.
4. Admission / Policy Guardrails: RHACS deploy-time custom policy (or Gatekeeper/Kyverno – choose one authoritative) to enforce presence & consistency of classification labels, forbid privileged/hostNetwork in high zones.
5. Namespace Strategy: Separate namespaces per classification (e.g., `apps-restricted`) plus deny-all ingress/egress; only explicit inter-zone NetworkPolicies allowed (justify each exception).
6. Differential Enforcement: Stricter runtime actions (block vs alert) and shorter vuln SLAs for higher zones (e.g., Critical fix ≤48h for restricted, ≤7d baseline elsewhere).
7. Secrets Handling: Enforce external vault references; block plain env secrets in restricted zone.
8. Drift Detection: Daily job enumerates pods where `data-classification` label mismatches node label; zero tolerance—auto ticket.
9. Residual Risk Register: Document shared kernel exposure & trigger conditions for migrating a zone to its own cluster (e.g., inability to meet accelerated patch SLA, regulatory mandate).
10. Exception Workflow: Temporary co-residency requires exception ID, risk rationale, expiry, and approval (tracked in Exception Register).

> Example enforcement logic (illustrative pseudocode – adapt to actual policy engine):
> IF namespace matches /(apps-confidential|apps-restricted)/ THEN
>  require label data-classification present AND
>  require node selector key classification == data-classification label AND
>  forbid privileged OR hostNetwork=true for data-classification in (restricted)
> VIOLATION if any condition fails
> (Store actual JSON export in Git; reference commit hash in evidence.)

Additional Evidence for Compute Zones:
- Node label & taint inventory export (hash + timestamp)
- RHACS classification enforcement policy export
- Daily drift report (pod↔node classification mismatch) with uninterrupted date chain
- Inter-zone flow matrix (approved NetworkPolicy exceptions) + ticket links
- Vulnerability SLA matrix per zone + sample accelerated remediation proof
- Exception register entries (if any) governing temporary deviations

Escalate to Separate Clusters When:
- Regulatory / contractual requirement for isolation beyond logical segmentation
- Inability to consistently meet hardened SLA / patch cadence for shared nodes
- Frequent contention or noisy neighbor undermining zone guarantees

Document the decision criteria so auditors see a rational progression plan from single-cluster multi-zone to multi-cluster architecture if/when triggers occur.

---
## 5. Resource Governance & Availability
**Representative Controls:** NIST SC‑6 (resource availability) / CM‑7 (least functionality applied via enforced limits) (CM‑7 partial); NIST 800‑190 4.2.5 (resource controls).
**CIS (v8) Representative Safeguards:** 4.6 (least functionality / limiting unnecessary resource exposure) – Appendix A.

### Intent
Prevent noisy-neighbor risk & resource exhaustion through enforced CPU/memory boundaries.

### RHACS Levers
- Policies: missing resource limits / requests

### OpenShift Levers
- LimitRange + ResourceQuota for namespaces

### Key Actions
1. Enforce policy requiring both CPU & memory limits.
2. Add namespace quotas aligned to capacity planning assumptions.
3. Alert on deployments lacking limits >24h after introduction.

### Additional Evidence
- Namespace quota report + variance to actual usage

---
## 6. Vulnerability Remediation Lifecycle (Fix & Prove Closure)
**Representative Controls:** NIST RA‑5 / SI‑2 (flaw remediation) (RA‑5 scanning + gating, SI‑2 remediation metrics); NIST 800‑190 4.1.4 (scanning – in scope), 4.1.6 (rebuild vs patch – External), 4.1.14 (license compliance – External), 4.4.2 (build integrity – partial). (Explicitly distinguishing external-only refs 4.1.6 & 4.1.14.)
**CIS (v8) Representative Safeguards:** 2.2, 6.2 (supported software / vulnerability remediation process) – Appendix A.

### Intent
Quickly identify and block the riskiest (fixable) vulnerabilities and prove you are rebuilding images instead of letting risk age out.

### RHACS Capabilities (Focus Only)
- Continuous image & component scanning (all connected registries)
- Policy gating: block deploy/build if image has fixable Critical (and later High) CVEs
- Severity + fixable filtering & age views; exportable reports / API
- Notifier-driven alert when a vulnerability breaches SLA

### OpenShift / Pipeline Capabilities
- Automated image rebuild on updated base image
- GitOps promotion restricted to images that passed RHACS policy (digest pinning)

### Simple Action Pattern
1. Publish a minimal SLA (Critical 7 days, High 30 days). Medium/Low = track only.
2. Enforce: block new images with fixable Critical CVEs; warn on High (plan date to move High → block).
3. Daily export a vulnerability summary (keep last 30 days + external SHA-256 digests for tamper-evidence) – optional but useful.
4. Rebuild & redeploy images failing policy; verify new digest shows “no fixable Critical”.
5. Track two metrics: (a) % fixable Critical within SLA (aim ≥95%), (b) Median days to fix Critical (TTRc) trending down.

### Evidence (Lightweight)
- Policy export (showing Critical=block)
- Sample blocked deployment (log or RHACS violation) with timestamp
- 30‑day vulnerability summary snippet (counts new/fixed/remaining Critical)

### Notes
- RHACS enforces & measures; it does not patch—your pipeline rebuilds.
- Any accepted exception must have an expiry (see Exception Register section).

---
## 7. Runtime Threat Detection & Automated Response
**Representative Controls:** NIST SI‑4 / IR‑4(5) / IR‑5 / IR‑6(1) / AU‑12 (audit generation); (SI‑3 partial – limited suspicious process patterns, not full malware engine); NIST 800‑190 4.5.1 / 4.5.2.
**CIS (v8) Representative Safeguards:** 2.6, 2.7, 8.2, 17.2 (allowlist runtime binaries/scripts, security event alerting, automated incident response) – Appendix A.

### Intent
Detect anomalous or malicious runtime activity and (optionally) apply automated containment.

### RHACS Levers
- Runtime process & network baseline anomalies, exec into container, crypto miner patterns
- Policy actions (scale-to-zero / block / alert) + notifier integrations

### Detection (Simple View)
RHACS does two things:
1. Baseline: learns normal processes / connections; anything new is flagged (new ≠ automatically bad, just unexpected).
2. Prebuilt risky patterns: detects obvious attacker / abuse behaviors (crypto miner names, curl|wget pipe to shell, package manager installs, reverse shell hints, privilege escalation attempts).

Not in scope: deep packet inspection, syscall tracing, full lateral movement analysis. Use other tools for those (list them as external controls).
> Caveat: RHACS uses kernel instrumentation (eBPF-based collection) to observe process executions and network connections, but its detection logic operates at the process/command + connection abstraction layer (baseline anomalies, known risky patterns) rather than exposing arbitrary raw syscall sequence rule authoring or deep packet payload inspection.

### OpenShift / Platform Levers
- Cluster audit logs for correlated identity context
- Network isolation reducing noise + containment domain

### Key Actions
1. Enable top critical runtime policies; attach at least one high-urgency notifier.
2. Test alert → ticket workflow end-to-end (document timing metrics).
3. Consider selective enforcement for high-confidence miner / privilege escalation.
4. Quarterly tune false positives (measure alert precision & drop noise >20%).

### Additional Evidence
- Alert precision metric (true positive / total high severity alerts) over last 30 days

---
## 8. Secrets & Sensitive Data Exposure Prevention
**Representative Controls:** NIST SI‑7 / SI‑7(1) (software/information integrity – secret tamper & exposure detection limited) (partial), SC‑28 (protection at rest – External for vault/KMS); NIST 800‑190 4.1.7 (secrets), 4.2.3 (secret management – External rotation & lifecycle), 4.4.1 (registry provenance – indirect integrity linkage). (SC‑28 added as external dependency; rotation explicitly external.)
**CIS (v8) Representative Safeguards:** 13.1, 3.3 (data protection / access control to secrets) – Appendix A.

### Intent
Prevent embedding or accidental leakage of secrets inside images or environment variables.

### RHACS Levers
- Secret pattern detection in env vars / config
- Deploy/build-stage blocking policy for explicit secret strings

### OpenShift / Platform Levers
- External secret operators (vault integration) & sealed secrets
- Encrypted storage for secret data at rest (platform managed)

### Key Actions
1. Enable secret-in-env detection; whitelist benign tokens.
2. Enforce policy for high-sensitivity keys (e.g., private keys) at deploy.
3. Migrate static credentials to external vault references; remove from Git.

### Detection Limitations
Secret pattern detection is heuristic/string-pattern based. It may *miss*:
- Encrypted or base64-obfuscated sensitive blobs masquerading as benign strings
- Secrets stored inside binary layers or compressed archives
- Proprietary token formats not matching default regexes

Additional Caveat: RHACS does **not** analyze the cryptographic strength, rotation interval, or entropy of values stored inside Kubernetes Secret objects; weak or long-lived keys must be governed by external secret management and rotation processes.
Explicit Out-of-Scope: Entropy assessment, key age tracking, automatic rotation enforcement, and revocation workflows all sit outside RHACS; treat these as External Control Register entries (Secrets Lifecycle & Rotation).

Do **not** rely on this as primary control—treat it as a compensating “last line” safety net. Primary controls: external vault, short-lived credentials, automated rotation.

### Additional Evidence
- Reduction count of secrets flagged in last 90 days

---
## 9. Logging, Reporting & Continuous Evidence
**Representative Controls:** NIST SI‑4 (monitoring contribution) / IR‑5 / IR‑6(1) / AU‑6 / AU‑12; (AU‑9 integrity & retention = External with partial tamper-evidence via hashes). (Clarifies AU‑9 externalization.)
**CIS (v8) Representative Safeguards:** 7.1, 7.3, 8.2, 16.12, 10.4 (logging establishment, centralization, alerting, provider log collection, secure backups) – Appendix A.

### Intent
Maintain immutable, correlated, reviewable evidence of control operation & exceptions.

### RHACS Levers
- Scheduled compliance & policy exports
- Alert forwarding to SIEM / ticketing

### OpenShift / Platform Levers
- External SIEM pipeline, log integrity (external hashing + WORM), retention policy enforcement
- Time sync (NTP/chrony) for consistent event ordering

### Key Actions
1. Nightly compliance export (compute external SHA-256 digest + store artifact; immutable/WORM storage is external).
2. Forward policy + runtime alerts to SIEM; alert on pipeline failures.
3. Implement log integrity verification (externally maintained SHA-256 hash chain / object lock). 
4. Quarterly Targeted Risk Analysis (TRA) if deviating from default review cadence.

### Additional Evidence
- Log pipeline health check report + failure alert test case
 - Admission webhook availability SLO report (ties to enforcement reliability)
 - Statement/evidence of external immutable storage (object lock / WORM) since in-cluster logging stacks are not inherently immutable

---
## 10. Quick Start Checklist
| Objective | Action (Do This) | Proof to Capture (Simple Evidence) |
|-----------|------------------|------------------------------------|
| Full scan coverage | Add all registries & enable RHACS scanner; rescan running images | Screenshot/export: 0 unscanned running images |
| Baseline config enforced | Turn on core misconfig policies (privileged, host mount, run as root, no limits) in enforce mode | Policy export showing Enforced=true + zero critical deploy violations |
| Vulnerability gate active | Enforce block on fixable Critical CVEs (warn on High initially) | Blocked deployment log + policy JSON (Critical=block) |
| Runtime visibility working | Verify collector healthy; enable 2–3 runtime policies; trigger safe test event | Runtime alert + notifier delivery record |
| Network segmentation started | Apply namespace default deny (ingress & egress) + first allow rules | NetworkPolicy manifests + coverage screenshot |
| RBAC hygiene | Reduce cluster-admin to one group; remove direct user bindings | Before/after clusterrolebinding diff (only one group) |
| Secret leak prevention | Enable secret-in-env detection; fix flagged env vars | Before/after secret violation count trend (→0) |
| Evidence automation | Schedule nightly compliance export & forward alerts/logs to SIEM | Stored report + external SHA-256 digest + SIEM entry with RHACS alert |

Note: Each row maps to detailed sections below. Non-experts can ignore framework/control IDs; auditors can use the Control Mapping table.

---
## 11. Control Mapping Quick Reference
High-level crosswalk of each theme to the principal framework control families. This is an orientation aid only—see Appendices A (CIS Controls v8 mapping), B (NIST 800-53) & D (NIST 800-190) for granular control-by-control coverage, notes, and evidence nuances.

| Theme | NIST 800-53 (Primary Technical Controls) | NIST 800-190 (Key Refs) |
|-------|-------------------------------------------|-------------------------|
| 1 Image Provenance & Supply Chain | CM-2, CM-6, RA-5, SI-7, SR-11 | 4.1.1–4.1.4, 4.1.11–4.1.12 |
| 2 Baseline Config & Drift | CM-2, CM-3, CM-6, CM-7, (CM-8) | 4.1.2, 4.1.3, 4.1.8, 4.2.7 |
| 3 Least Privilege & RBAC | AC-2, AC-3, AC-6, CM-5 | 4.2.1, 4.2.4 |
| 4 Network Segmentation | SC-7 (+ SC-7(3)/(4)), AC-3 (enforcement tie) | 4.2.2 |
| 5 Resource Governance | SC-6, CM-7, (CP-10 External) | 4.2.5 |
| 6 Vulnerability Lifecycle | RA-5, SI-2 | 4.1.4, 4.1.6, 4.1.14 |
| 7 Runtime Detection & Response | SI-4, IR-4(5), IR-5, IR-6(1), AU-12, (SI-3) | 4.5.1, 4.5.2 |
| 8 Secrets Protection | SI-7, SI-7(1), (SC-28) | 4.1.7, 4.2.3 |
| 9 Logging & Evidence | AU-6, AU-12, AU-9 (partial), IR-6(1), SI-4 | 4.2.6 |

Tri-Column Coverage Model (applies to all control mapping & evidence tables): Columns enumerate OCP (OpenShift/RHCOS primitives), RHACS (security overlay), External (out-of-scope systems/governance). Per column: C = fully enforced/evidenced within that layer; P = partial contribution (shared responsibility or evidentiary assist); blank = negligible/no substantive contribution; External column uses E when entirely outside OCP+RHACS scope. See Section 0 for model rationale.

Notes:
- CP-10 (System Recovery) explicitly categorized External: DR plan execution & recovery testing lie outside platform/RHACS evidentiary scope.
- AU-9 listed under Theme 9 for thematic alignment; detailed mappings classify it External (artifact emission only; tamper-resistance & immutability external).
- AU-12 included for Runtime Detection: runtime alerts contribute to the security event corpus.

> Interpretation Nuance: If an intended "C" capability is temporarily not enforced (e.g., policy in warn, webhook fail-open), treat it operationally as downgraded (manage via exception register) until restored—do not silently leave as C in internal audit prep artifacts.

---
## Appendix A – CIS Controls v8.1 Mapping (Container Platform Focus)
This appendix distills the CIS Controls v8.1 safeguards most materially impacted by OpenShift (OCP) + RHACS capabilities in a container/Kubernetes context. It is a *scoped subset* (not the full CIS catalog). Source: internal reduction of provided CIS→NIST mapping CSV; selection criteria:
1. Technical enforceability or evidentiary contribution by OCP/RHACS (vs purely programmatic / policy).  
2. Relevance to container workload security (image, runtime, network, secrets, logging, vulnerability, RBAC).  
3. Mappable to existing Themes 1–9 without creating a new theme.  

Legend: OCP, RHACS, External columns retain C (Complete), P (Partial), blank (none), E (External only – used in External column). “Primary Theme” = dominant theme; some safeguards span multiple themes (listed comma-separated). “NIST Crosswalk” uses family identifiers (subset) for orientation; refer to Appendix B for depth.  

| CIS Control (v8) | Safeguard (Short) | Primary Theme(s) | OCP | RHACS | External | NIST Crosswalk (High-Level) | Evidence Pointer |
|------------------|------------------|------------------|-----|-------|----------|---------------------------|------------------|
| 1.1 Asset Inventory (Enterprise Assets) | Node / cluster component inventory (scoped) | 2,5 | P |  | E | CM-8, CM-2 | MachineConfig & node list + external CMDB ref |
| 1.2 Address Unauthorized Assets | Detect & act on unauthorized nodes/components | 2 | P |  | E | CM-8(3), SI-4 | Cluster node list diff & exception ticket |
| 2.1 Software Inventory (Container Images) | Authorized image registry scope | 1 | P | P | E | CM-8, CM-7 | Allowed registry policy + blocked unknown image log |
| 2.2 Supported Software Only | Remove unsupported base images | 1,6 |  | P | E | SI-2, RA-5 | Vulnerability report highlighting EOL base image |
| 2.3 Address Unauthorized Software | Block disallowed packages / images | 1,2 | P | C |  | CM-7, RA-5 | Policy blocking disallowed image + violation log |
| 2.5 Allowlist Authorized Software | Enforce only trusted signed images | 1 | P | P | E | SI-7, SR-11 | Signature verification config + blocked unsigned deploy |
| 2.6 Allowlist Libraries (Runtime Binaries) | Restrict unexpected binaries/processes | 7 |  | P | E | SI-4, CM-7 | Runtime anomaly alert showing new binary |
| 2.7 Allowlist Scripts | Restrict unauthorized scripts in build/runtime | 1,7 |  | P | E | SI-7(1), CM-7 | Policy or runtime detection event |
| 3.3 Data Access Control Lists (scoped to secrets/config) | Limit secret/config access | 3,8 | C | P | E | AC-3, AC-6 | RBAC diff + secret detection trend |
| 3.4 Enforce Data Retention (logs/security evidence) | Retention of security artifacts | 9 |  | P | E | AU-11, SI-12 | Compliance export + external retention policy |
| 4.1 Secure Configuration (Container Platform) | Baseline hardened config | 2 | C | P | E | CM-2, CM-6 | SCC / policy set export + hardening doc |
| 4.6 Least Functionality | Remove unnecessary privileges/capabilities | 2,5,7 | C | P |  | CM-7, AC-6 | Blocked privileged deploy log |
| 4.8 Controlled Use of Administrative Privileges | Constrain cluster-admin | 3 | C | P | E | AC-2, AC-6 | Quarterly RBAC review sign-off |
| 5.1 Account Management (Service Accounts) | Manage k8s service accounts | 3 | C | P | E | AC-2, AC-3 | Service account inventory + anomaly report |
| 6.2 Maintain Vulnerability Remediation Process | Apply fixable image remediation | 6 |  | C | E | RA-5, SI-2 | Blocked Critical CVE deploy + rebuilt digest |
| 7.1 Establish & Maintain Logging | Generate security event logs | 9 | P | P | E | AU-12, AU-6 | Alert forwarding config + sample event |
| 7.3 Centralize Log Management | Central aggregation (external stack) | 9 |  | P | E | AU-6, AU-9 | External SIEM ingestion dashboard |
| 8.2 Security Event Alerting | Runtime / deploy / vuln alerts | 1,6,7,9 |  | C | E | SI-4, IR-5 | Alert + ticket linkage evidence |
| 8.4 Conduct Security Event Correlation | Cross-source correlation | 9 |  |  | E | AU-6(3), IR-5 | SIEM correlation rule pack diff |
| 10.4 Secure Backups (Configuration & Evidence Artifacts) | Secure backup of config/policy exports | 9 |  | P | E | CP-9, SI-12 | Backup job log + external hash chain |
| 11.1 Standard Secure Config Process (Policy-as-Code) | Git-managed policy baselines | 1,2,6 | P | P | E | CM-3, CM-5 | Policy repo commit (signed) + enforcement diff |
| 12.1 Network Segmentation | Deny-by-default pod traffic | 4 | C | P | E | SC-7 | NetworkPolicy coverage % trend |
| 12.4 Boundary Defenses (Ingress/Egress Control) | Control external egress exposure | 4 | C | P | E | SC-7(3) | Egress restrict policy manifest + review ticket |
| 13.1 Data Protection (Secrets Handling) | Secret exposure prevention | 8 | P | P | E | SI-7, SC-28 | Secret violation trend & vault rotation report |
| 16.12 Collect Service Provider Logs (Container Security Tool) | RHACS export consumption | 9 |  | C | E | AU-12 | Compliance export + external digest |
| 17.2 Remediate Detected Incidents (Automated Containment) | Automated runtime response | 7 |  | C |  | IR-4(5) | Runtime kill/scale action log |

Notes:
1. Some CIS safeguards encompass broader enterprise scope (asset management across all endpoints). Only the container-platform slice is represented here; remainder is External (E) or Out-of-Scope.  
2. “Allowlist” safeguards in container context map to enforcing: signed/trusted images, restricted capabilities, and runtime baseline process sets.  
3. Data classification & retention beyond security evidence artifacts remain external governance responsibilities (policy, DLP tooling).  
4. Evidence pointer examples should be adapted to your actual artifact naming & storage conventions (append SHA-256 digests where feasible).  

### A.1 Thematic Coverage Density
| Theme | Key CIS Safeguards (Representative) |
|-------|--------------------------------------|
| 1 Image Provenance | 2.1, 2.2, 2.3, 2.5, 11.1 |
| 2 Baseline Config & Drift | 4.1, 4.6, 11.1 |
| 3 Least Privilege & RBAC | 4.8, 5.1, 3.3 |
| 4 Network Segmentation | 12.1, 12.4 |
| 5 Resource Governance | 4.6 (capacity/limit facets) |
| 6 Vulnerability Lifecycle | 2.2, 6.2 |
| 7 Runtime Detection & Response | 2.6, 2.7, 8.2, 17.2 |
| 8 Secrets & Sensitive Data | 13.1, 3.3 |
| 9 Logging & Evidence | 7.1, 7.3, 8.2, 16.12, 10.4 |

### A.2 Methodology Excerpt
The CSV mapping enumerates CIS <-> NIST 800-53 relationships (subset/superset/equivalent). We normalized by:
1. Collapsing multi-row CIS–NIST expansions into a single row per safeguard with dominant NIST family set.  
2. Filtering out purely procedural safeguards (e.g., policy governance with no direct OCP/RHACS artifact).  
3. Tagging Partial (P) where OCP/RHACS generate some evidence (e.g., asset visibility, log emission) but authoritative inventory, correlation, or retention is external.  
4. Aligning each safeguard to the pre-existing Theme whose intent most closely matches its risk reduction objective (prefer minimal theme proliferation).  

### A.3 Gap & Externalization Highlights
- Asset & Software Inventories (CIS 1.x / 2.1) are only partially met: cluster & image visibility ≠ authoritative enterprise inventory (External CMDB/SBOM system).  
- Retention / Correlation (7.3, 8.4) rely on external SIEM and WORM storage; platform exports are necessary but insufficient.  
- Script & Library allowlisting (2.6, 2.7) limited to anomaly/risky pattern detection—not full allow-by-exception engine; reflect as Partial (P).  
- Data retention & protection (3.4, 13.1) beyond secrets exposure is external (encryption, lifecycle, classification tooling).  

### A.4 Minimal CIS Evidence Bundles (Examples)
| CIS Safeguard | OCP Artifact | RHACS Artifact | External Artifact | Rationale |
|--------------|-------------|---------------|------------------|-----------|
| 2.5 Trusted Images | ClusterImagePolicy + signature config | Blocked unsigned image violation | Key custody / signer rotation SOP | Shows trust policy + enforcement + governance |
| 4.6 Least Functionality | SCC / PodSecurity profile list | Policy violations trend (privileged, capabilities) | Hardening standard excerpt | Validates config baseline + active enforcement |
| 6.2 Vulnerability Remediation | (optional) Digest pinning manifest | Blocked Critical CVE deploy log + vuln trend | Rebuild pipeline log | Demonstrates detection → gate → fix chain |
| 12.1 Segmentation | NetworkPolicy manifests + coverage % | Missing policy alert trend | Mesh mTLS / perimeter ACL (if used) | Layered segmentation evidence |
| 7.1 / 7.3 Logging & Centralization | Audit forwarder config | Scheduled compliance export + hash | SIEM ingestion dashboard + retention policy | Emission + integrity + centralized retention |

---
## Appendix B – NIST 800‑53 Runtime / Incident Subset
| Control | Theme | OCP | RHACS | External | Notes |
|---------|-------|-----|-------|----------|-------|
| SI‑4 | 7 |  | C |  | Runtime process & network telemetry (RHACS baseline + anomaly). |
| IR‑4(5) | 7 |  | C |  | Automated response via runtime policy actions & notifiers. |
| IR‑5 | 7 |  | C |  | Continuous runtime incident monitoring. |
| IR‑6(1) | 7 / 9 |  | C |  | Automated reporting to external systems via notifiers. |
| AU‑6 / AU‑12 | 9 | P | P | E | OCP audit/log events (P) + RHACS security events (P); full centralized correlation & retention external. |

### B.1 Expanded NIST 800‑53 Mapping (Selected High‑Relevance Controls)
Focused on Moderate baseline (Rev 5) control families most often cited in platform/container security audits. Not exhaustive; omit families with minimal direct technical tie (e.g., PE – Physical) or purely programmatic (e.g., PM) where RHACS offers no evidence. Use this as a *translation accelerator*, not a replacement for a formal System Security Plan (SSP).

| Control (Rev5) | Control Title (Abbrev) | Primary Theme(s) | OCP | RHACS | External | Notes / Contribution & Boundaries |
|----------------|------------------------|------------------|-----|-------|----------|------------------------------------------|
| AC‑2 / AC‑2(1) | Account Management / Automated Disable | 3 | C | P | E | OCP RBAC enforces; RHACS highlights cluster-admin subjects; lifecycle (provision/disable) external IAM. |
| AC‑3 | Access Enforcement | 3 / 4 | C | P |  | RBAC & NetworkPolicy in OCP; RHACS evidences usage/drift. |
| AC‑6 / AC‑6(1) | Least Privilege / Authorizations | 3 | C | P | E | OCP roles & SCC enforce; RHACS detects over-privilege; approval workflow external. |
| AC‑17 / AC‑17(2) | Remote Access / AuthN Strength | 3 / External |  |  | E | MFA / remote access controls external (IdP, bastion). |
| AC‑19 | Access Control for Mobile / BYOD | External |  |  | E | Not applicable to in-cluster workloads. |
| AU‑2 / AU‑2(3) | Event Logging / Central Review | 9 | P | P | E | OCP audit events + RHACS security events; central correlation & full audit set external. |
| AU‑6 / AU‑6(3) | Audit Review / Correlation | 9 | P | P | E | Requires SIEM correlation externally. |
| AU‑8 | Time Stamps | 9 | P |  | E | Node/cluster NTP (platform); correlation governance external. |
| AU‑9 / AU‑9(2) | Audit Protection / Tamper Resistance | 9 |  |  | E | Tamper resistance & immutability entirely external (WORM/object-lock + hash chain). Platform/RHACS only emit artifacts. |
| AU‑12 | Audit Generation | 9 | P | P | E | Partial security telemetry only; full audit scope external. |
| CA‑7 | Continuous Monitoring | 1–9 | P | P | E | Contributes technical signals; org-wide monitoring strategy external. |
| CM‑2 / CM‑2(2) | Baseline Configuration / Automation | 1 / 2 | C | P | E | OCP declarative config (MachineConfig, SCC); RHACS drift/misconfig detection; baseline approval external. |
| CM‑3 | Configuration Change Control | 1 / 2 / 6 | P | P | E | Evidence diffs (P); formal CAB external. |
| CM‑5 | Access Restrictions for Changes | 3 | C | P | E | RBAC gating (C); RHACS visibility (P); Git approval external. |
| CM‑6 | Configuration Settings | 1 / 2 | C | P | E | OCP enforces via SCC/Policies; RHACS policy mapping; hardening catalog external. |
| CM‑7 / CM‑7(1) | Least Functionality / Prevent Unauthorized Software | 1 / 2 / 7 | P | P | E | Platform restricts privilege/capabilities (P); RHACS anomaly/risky binary detect (P); allowlist governance external. |
| CP‑9 | Information System Backup | External |  |  | E | Backup & restore validation external. |
| CP‑10 | System Recovery | External |  |  | E | DR exercises external. |
| IR‑4 / IR‑4(5) | Incident Handling / Automated Response | 7 |  | C |  | Runtime action policies + notifiers. |
| IR‑5 | Incident Monitoring | 7 |  | C |  | Continuous runtime observation. |
| IR‑6 / IR‑6(1) | Incident Reporting / Automated Reporting | 7 / 9 |  | C |  | Automated forwarding to SIEM/ticket. |
| IR‑8 | Incident Response Plan | External |  |  | E | Human process & documentation external. |
| MA‑4 | Nonlocal Maintenance | External |  |  | E | Platform ops domain external. |
| RA‑5 / RA‑5(2) | Vulnerability Monitoring / Update Mechanisms | 6 | P | C | E | RHACS scans/gating (C); OCP assists via image pinning (P); host/non-container assets external. |
| SA‑11 | Developer Security Testing | 1 / 6 |  | P | E | Post-build gating only; SAST/DAST external. |
| SA‑15 | Development Process / Standards | External |  |  | E | Secure SDLC governance external. |
| SC‑6 | Resource Availability Protection | 5 | C | P |  | OCP quotas/limits (C); RHACS missing limits detection (P). |
| SC‑7 / SC‑7(3)/(4) | Boundary Protection / Segmentation / Deny by Default | 4 | C | P | E | NetworkPolicy enforcement (C); RHACS coverage/flow viz (P); L7/WAF/mTLS governance external. |
| SC‑8 / SC‑8(1) | Transmission Confidentiality & Integrity | 4 | P | P | E | OCP provides TLS ingress + optional mesh mTLS/IPsec capabilities (partial until universal enforcement evidenced). RHACS may surface plaintext endpoint findings (partial only). Key lifecycle & cipher policy external. |
| SC‑13 | Cryptographic Protection (At Rest) | External |  |  | E | Storage/etcd encryption external. |
| SC‑28 | Protection of Information at Rest | External |  |  | E | Volume/database encryption external. |
| SI‑2 | Flaw Remediation | 6 | P | C | E | RHACS detects vulnerable images & enforces gates (C); rebuild orchestration external (E). |
| SI‑3 | Malicious Code Protection | 7 |  | P | E | Suspicious process patterns only; traditional AV external. |
| SI‑4 / SI‑4(2)/(4) | System Monitoring / Indicators / Traffic Anomalies | 7 / 9 | P | C | E | OCP provides baseline audit/network constructs (P); RHACS anomaly detection (C); deep packet/IDS external. |
| SI‑5 | Security Alerts / Advisories | 6 / 9 |  | P | E | CVE/advisory ingestion; enterprise advisory program external. |
| SI‑7 | Software / Information Integrity | 1 / 8 | P | P | E | RHACS enforces trusted-signer policies & secret pattern detection (partial); OCP admission verifies signatures (partial). Key custody, Rekor/attestation chain, SBOM integrity & rotation external. |
| SI‑10 | Information Input Validation | External |  |  | E | App/WAF layer control external. |
| SR‑11 | Component Authenticity | 1 | P | P | E | Signature presence policies (RHACS) + admission config (OCP) (P/P); attestation chain external. |

Legend Recap (Appendix B): Column-specific. OCP: OpenShift/RHCOS primitives. RHACS: overlay detection/enforcement/evidence. External: out-of-scope systems or governance. C = fully enforced/evidenced in that column; P = partial contribution; blank = negligible; E (External column only) = entirely external responsibility.

*SC‑8 Clarification:* Marked Partial because OpenShift natively terminates and serves TLS for Routes/Ingress, can enable encrypted node overlay (IPsec depending on network configuration/version), and (optionally) Service Mesh supplies mTLS for east‑west traffic. RHACS itself does not generate, rotate, or validate certificates or cipher policies—capture evidence via: Ingress Controller TLS config/certificate inventory, mesh PeerAuthentication / DestinationRule (or equivalent) showing STRICT mTLS, and (if applicable) cluster network encryption status documentation. If none of these platform features are enabled yet, downgrade SC‑8 to External (E) until cryptographic controls are operational. Mark External (E) as well if you cannot produce evidence of consistent cluster-wide TLS/mTLS/IPsec enforcement.

> Implementation Tip: When building an SSP, cite this table and then link each (P) / (E) control to either (a) platform configuration export (e.g., NetworkPolicy manifests, SCC profiles, mesh mTLS policy) or (b) governance artifacts (CAB approvals, IR plan version). For (C) items, embed RHACS policy JSON export plus its externally computed SHA-256 digest + sample violation or compliance report line item (digest provides tamper-evidence, not immutability).

### B.2 Tailoring & Gaps
1. Tailor out controls not applicable to container platform scope (e.g., AC‑19) to prevent artificial gap listings.
2. For each (E) control, ensure an owner appears in the External Control Register (Appendix E) — absence indicates governance risk.
3. For mixed controls (C/P), define an internal rule: treat an unmet prerequisite (e.g., admission webhook fail‑open) as temporary downgrade → mark exception with expiry.
4. Maintain a delta log: when RHACS adds functionality narrowing a (P) control toward (C), update this appendix and version the change (auditors appreciate traceability).
 5. Mapping Correction (2025-10-03): Adjusted over-crediting of RHACS/OCP for AU-9, SC-8, SI-7, and NIST 800-190 refs 4.1.6, 4.1.7, 4.1.11, 4.1.12 to align with strict enforcement vs evidentiary vs external process boundaries.

### B.3 Minimal Evidence Bundles (Examples)
| Control Focus | OCP Artifact | RHACS Artifact | External Artifact | Sufficiency Rationale |
|---------------|-------------|---------------|-------------------|-----------------------|
| RA‑5 (Vuln Monitoring) | (optional) Image digest pinning manifest | Vulnerability report export (timestamp + external digest) | Pipeline rebuild log referencing digest | Correlates detected risk → enforced gate → rebuild action. |
| SC‑7 (Segmentation) | NetworkPolicy manifest set + coverage % | Coverage trend graph | Mesh mTLS policy export + firewall ACL | Validates layered L3/L4 deny + L7/mTLS + perimeter segmentation. |
| CM‑2 (Baseline Config) | MachineConfig & SCC profile list | Policy set JSON (signed commit) | Hardening standard doc version | Links declarative baseline → enforcement → approved standard. |
| IR‑4(5) (Automated Response) | (N/A) | Runtime policy kill/scale action log | Incident ticket with closure notes | Shows automated containment tied to formal IR follow-up. |
| AU‑9 (Audit Protection) | Audit forwarder config checksum | External hash chain index of daily exports | Object store WORM policy export | Demonstrates end-to-end tamper-evidence + immutable retention chain. |


---
<!-- Appendix C removed per scope reduction request (all health privacy clauses out of scope). -->

---
## Appendix D – NIST SP 800‑190 Section 4.1.x
| Ref | Theme | OCP | RHACS | External | Intent |
|-----|-------|-----|-------|----------|--------|
| 4.1.1 | 1 | P | P | E | Minimal base surface via image sourcing + policy detection. |
| 4.1.2 | 1 | P | C | E | Trusted registries (policy enforcement primarily RHACS); OCP admission config (P). |
| 4.1.3 | 2 | P | P | E | Remove unnecessary components (policy + build guidance externally governed). |
| 4.1.4 | 1 / 6 |  | C |  | Pre-deploy scanning (RHACS). |
| 4.1.6 | 6 |  |  | E | Rebuild vs patch workflow entirely external (pipeline responsibility); RHACS only observes vulnerable images (no rebuild enforcement). |
| 4.1.7 | 8 | P | P |  | Secret detection is heuristic (partial); comprehensive secret management & rotation external. |
| 4.1.8 | 2 | C | P |  | Non-root enforcement (SCC) + RHACS detection. |
| 4.1.9 | 2 | C | P |  | Read-only FS (SCC/Pod settings) + policy detection. |
| 4.1.10 | 1 | P | P | E | SBOM correlation partial; generation external. |
| 4.1.11 | 1 | P | P | E | Signature/integrity gating shared: OCP admission verifies; RHACS policy checks presence (partial). |
| 4.1.12 | 1 | P | P |  | Digest pinning detection/enforcement shared (partial) – policies detect but do not guarantee universal pinning. |
| 4.1.13 | 2 | C | P |  | Remove escalation paths (SCC + detection). |
| 4.1.14 | 1 |  |  | E | License compliance external. |
| 4.1.15 | 1 | P | P | E | Secure build pipeline integrity evidence external; policy gating partial. |

Legend: * items with platform or process focus outside RHACS core are marked (P) Partial or (E) External.

### D.1 Additional 800‑190 Sections (Selected)
| Ref | Area | Theme(s) | OCP | RHACS | External | Notes |
|-----|------|----------|-----|-------|----------|-------|
| 4.2.1 | Orchestrator Access Control | 3 | C | P | E | RBAC enforcement (C); RHACS visibility (P); strong authN external. |
| 4.2.2 | Segmentation & Network Policy | 4 | C | P | E | NetworkPolicy enforcement (C); coverage analytics (P); L7/mTLS governance external. |
| 4.2.3 | Secret Management | 8 | P | P | E | Leak detection (RHACS) + basic secret object handling (OCP P); vault & rotation external. |
| 4.2.4 | Limit Privileges | 2 / 3 | C | P | E | SCC/Pod Security (C); RHACS detection (P); host hardening external. |
| 4.2.5 | Resource Controls | 5 | C | P |  | Quotas/limits (C); missing limits detection (P). |
| 4.2.6 | Logging & Monitoring | 7 / 9 | P | P | E | Partial security events; full central logging external. |
| 4.2.7 | Admission & Policy Enforce | 1 / 2 | P | C | E | RHACS gating (C); OCP admission ordering & SCC (P); fail-closed config required; some policy roots external. |
| 4.3.1 | Host Hardening | External |  |  | E | CIS benchmark outside scope. |
| 4.3.2 | Host Vulnerabilities | External |  |  | E | Host scanning agents external. |
| 4.4.1 | Registry Security | 1 | P | P | E | Allowed registry policy (RHACS + OCP admission partial); registry RBAC external. |
| 4.4.2 | Build Integrity | 1 / 6 | P | P | E | Gating on outputs; provenance attestations external. |
| 4.5.1 | Runtime Threat Detection | 7 |  | C | E | Runtime baselines & patterns; deep forensics external. |
| 4.5.2 | IR Integration | 7 / 9 |  | C | E | Notifiers feed IR; runbooks external. |

### D.2 800‑190 Evidence Bundles
| Focus | OCP Artifact | RHACS Artifact | External Artifact | Narrative |
|-------|-------------|---------------|-------------------|----------|
| Segmentation | NetworkPolicy manifests + coverage % | Coverage % + flow graph | Firewall/mesh policy export | Layered segmentation (deny baseline, L7/mTLS, perimeter). |
| Secrets | (optional) External secret operator config reference | Secret violation trend | Vault rotation report | Detection complementing managed secret lifecycle & rotation. |
| Build Integrity | ClusterImagePolicy / admission verify config | Blocked unsigned image log | Pipeline attestation (SLSA/in‑toto) | Trust chain: config → enforcement → provenance proof. |
| Runtime Detection | (N/A) | Runtime alert → ticket | IR ticket with closure | Detection-to-response validation. |

### D.3 Quick Gap Check
1. Image from unapproved registry? (4.4.1)
2. Namespace missing deny-all baseline? (4.2.2)
3. Privileged container present? (4.2.4)
4. Missing attestation for critical service image? (4.4.2)
5. Collector coverage <100% nodes? (4.5.1 risk)

### D.4 Tailoring Note
Document that deep packet inspection, full memory forensics, and persistent packet capture are out-of-scope; list compensating tools in External Control Register.

## Appendix E – Clarification Index (Nuanced Partial / Shared Controls)
| Topic / Control Aspect | OCP Provides (Evidence / Enforcement) | RHACS Provides (Evidence / Enforcement) | External Responsibility | Recommended Evidence Bundle |
|------------------------|----------------------------------------|----------------------------------------|------------------------|-----------------------------|
| SBOM Association & Signature Policy | Admission signature verification config (ClusterImagePolicy), digest pinning | Policy gating on unsigned images / metadata; violation & compliance reports | SBOM generation, signing, storage (RHTAP), key custody | Policy JSON, admission config, blocked deploy log, SBOM file + cosign verify output |
| Vulnerability SLA Enforcement | (Optional) Image pinning manifests | Detects fixable CVEs, blocks severity, age metrics | Rebuild workflows, base image maintenance, change approvals | SLA matrix doc, vuln trend export, rebuild pipeline logs |
| Host / Node Hardening | MachineConfig enforcement, SCC restricting host access | Detection of privileged/host mount attempts | CIS benchmark, kernel params, firewall, patch cadence | CIS report, MCO diff, absence of privileged containers evidence |
| Network Segmentation | NetworkPolicy enforcement, namespace isolation baseline | Gap detection (missing policies), flow visualization | L7 authZ, DPI, IDS/IPS, mesh mTLS governance | NetworkPolicy coverage report, mesh policy export, IDS alert sample |
| Secrets Exposure | Platform secret objects, external secret operator integration | Pattern-based env/config secret detection | Vault storage, rotation, short-lived creds | Secret violation trend, vault rotation report, exception register |
| Runtime Threat Detection | (Baseline isolation reducing noise) | Process/network anomaly policies, notifier evidence | Full IR runbooks, forensics, SIEM correlation rules | Runtime alert sample + ticket, IR runbook version, SIEM correlated event |
| Logging & Integrity | Audit log emission & forwarding config | Alert/log export events, compliance scheduling, external hash chain | WORM storage, retention, central correlation | External hash chain index, object lock config, SIEM ingestion dashboards |
| License Compliance (800-190 4.1.14 External) | (N/A) | (Indirect) package inventory via scans | License analysis, legal approval workflow | License scan diff, approval tickets, component report snapshot |
| Signature / Attestation Chain | Signature verification enforcement (admission) | Policy check for presence of signatures/labels | Key custody, Rekor transparency validation | Key management SOP, cosign verify log, policy pass report |
| Policy Exceptions | (N/A) | Violation visibility, enforcement phase tracking | Governance workflow (approvals, expiry) | Exception register, policy diff, closure ticket |
| Admission Reliability | Admission ordering & failurePolicy on platform webhooks | Webhook evaluation results & error metrics | HA config, cert rotation, DNS/network reliability | Webhook SLO dashboard, error metrics, blocked vs allowed stats |
| Resource Governance | Quotas & LimitRanges enforce ceilings | Missing limit detection policies | Capacity planning & autoscaling strategy | Limits compliance pass, quota manifests, utilization vs quota report |
| Data-in-Transit Security | Ingress TLS termination, optional mesh mTLS/IPsec | (Optional) Detection of plaintext endpoints (custom) | Certificate lifecycle, cipher policy management | Mesh cert inventory, gateway TLS config, detection result |
| Data-at-Rest Integrity | Encrypted etcd/storage (platform config) | Enforce signed/immutable images, non‑root | Storage encryption keys, snapshot protection | KMS config, snapshot immutability proof, non-root policy pass |
| SBOM vs Component Inventory | (N/A) | Vulnerability-derived component listing | Formal SPDX/CycloneDX & license context | Component export, SBOM file external digest, license report |

> Use this index to pre-empt auditor “scope inflation” questions: for each shared control, you present split responsibilities plus cohesive evidence chain.

---
## Appendix F – Representative Controls Validation & Corrections (Added 2025-10-03)
This appendix validates the “Representative Controls” declared in Themes 1–9 against the detailed crosswalks (Appendix A – CIS Controls v8, Appendix B – NIST 800‑53, Appendix D – NIST 800‑190) and internal coverage model. It highlights omissions, over-attributions, and clarifies External vs Partial boundaries. Use this to defend scoping decisions and avoid inadvertent over-crediting during audits.

### F.1 Validation Summary Table
| Theme | Declared (Original) | Key Verified Control Families (Adjusted) | Additions Made | Over-Credited / Moved to External | Rationale for Adjustment |
|-------|---------------------|------------------------------------------|----------------|-----------------------------------|--------------------------|
| 1 Image Provenance | CM-2, CM-6, RA-5; 800-190 4.1.x | + SI-7, SR-11; refine 4.1.x to 4.1.1–4.1.4, 4.1.10–4.1.12 | SI-7, SR-11 | SBOM gen (4.1.10) external portion; attestation chain external | Aligns with Quick Reference (SI-7, SR-11); narrows 4.1.x to actually enforced subrefs; prevents implying SBOM creation. |
| 2 Baseline Config | CM-2, CM-3, CM-6; 800-190 hardening | + CM-7; 4.1.3, 4.1.8–9, 4.1.13, 4.2.7 | CM-7 | — | Least functionality (CM‑7) is materially enforced via SCC & policy gating. |
| 3 Least Privilege | AC-2, AC-3, AC-6, CM-5 | + AC-6(1) implied; 4.2.1, 4.2.4 | 4.2.1/4.2.4 refs | — | Adds orchestrator & privilege limit refs for completeness; doesn’t alter scope. |
| 4 Segmentation | SC-7 (+variants); 800-190 isolation | + SC-7(3)/(4), AC-3, 4.2.2 | SC-7(3)/(4), AC-3 | — | Explicit deny-by-default & enforcement ties; consistent with Quick Reference nuance. |
| 5 Resource Gov | SC-6 | + CM-7, 4.2.5 | CM-7 | — | Enforced limits reduce functionality surface (CM‑7). |
| 6 Vuln Lifecycle | RA-5; 4.1.4 / 4.1.6 / 4.1.14 | + SI-2; 4.4.2 (partial) | SI-2 | 4.1.6, 4.1.14 external | Distinguishes detection (RA‑5) vs remediation measurement (SI‑2); license & rebuild workflow external. |
| 7 Runtime Threat | SI-4, IR-4(5), IR-5, IR-6(1) | + AU-12; (SI-3 partial) | AU-12, SI-3(partial) | — | Runtime alert generation contributes to audit generation. |
| 8 Secrets | SI-7 | + SI-7(1), SC-28 (external), 4.1.7, 4.2.3 | SI-7(1), SC-28(ext) | Rotation & KMS external | Clarifies integrity enhancement & at-rest protection dependency on external vault/KMS. |
| 9 Logging & Evidence | SI-4 / IR-5 / IR-6(1) / AU-6 / AU-12 | + AU-9 (external integrity) | AU-9 (ext) | AU-9 internal enforcement | AU‑9 clarified as External (immutability & retention), partial tamper-evidence only. |

Legend: “Additions Made” = newly surfaced families/refs now explicit. “Over-Credited” lists items originally implied as in-scope but actually External per tri-column model.

### F.2 Notable Adjustments & Justifications
1. RA‑5 Duplication (Themes 1 & 6): Kept in Theme 1 (image provenance scanning gate) but primary remediation metric ownership remains Theme 6; label RA‑5 in Theme 1 as scanning-centric.
2. SI‑7 & SR‑11 in Theme 1: These integrity/authenticity controls materially apply to signature verification & component trust; omission created a gap between Quick Reference and theme declaration.
3. CM‑7 Surfacing (Themes 2 & 5): Least functionality principles are enforced via removal of privileged constructs and mandatory resource limits; surfaced to avoid under-scoping baseline.
4. Externalization of 4.1.6 & 4.1.14: Rebuild (supply chain remediation pipeline) and license compliance are process/tooling domains outside RHACS/OpenShift enforcement; previously unqualified listing could imply full coverage.
5. AU‑9 Clarification (Theme 9): Platform + RHACS deliver emission + tamper-evidence (hash). Immutability/retention (object lock/WORM) is external; reclassified to prevent over-credit.
6. SC‑28 & Secrets: At-rest encryption & key lifecycle governed externally; Theme 8 now explicitly distinguishes detection vs cryptographic control responsibilities.
7. SI‑3 Partial (Theme 7): Runtime suspicious process detection ≠ full malware engine; partial attribution prevents scope inflation.

### F.3 Guidance for Future Additions
- Any new representative control must cite: (a) enforcement locus (OCP, RHACS, External) and (b) evidence artifact category. Append to this appendix with date + commit hash.
- When a Partial (P) capability graduates to Complete (C) (e.g., new RHACS feature adds stronger integrity assertion), update both the Theme line and Appendix tables; keep a “delta” note for one review cycle.

### F.4 Evidence Alignment Checklist (Per Theme)
| Theme | Minimum Evidence Pair (Platform + RHACS) | External Evidence Needed (if any) |
|-------|-------------------------------------------|-----------------------------------|
| 1 | Signature verification config + blocked unsigned deploy log | SBOM file + signed digest, key custody SOP |
| 2 | MachineConfig/SCC export + misconfig policy JSON | Hardening standard document (approved) |
| 3 | RBAC diff + privilege anomaly report | IAM approval workflow logs |
| 4 | NetworkPolicy coverage trend + flow graph | Mesh mTLS policy / firewall ACL |
| 5 | Quota & LimitRange manifests + “missing limits” violation trend | Capacity plan (utilization vs forecast) |
| 6 | Vulnerability gate policy + blocked Critical CVE event | Rebuild pipeline log / SLA exception TRA |
| 7 | Runtime alert + automated action (scale/kill) log | IR ticket closure summary |
| 8 | Secret violation trend + enforced secret policy config | Vault rotation report / key policy |
| 9 | Compliance export + external hash chain index | WORM retention config / SIEM ingestion health |

### F.5 Version Note
Initial validation pass (2025‑10‑03) introduced no scope *reductions*; only clarifications and explicit additions. Future deltas must reference this appendix section number (H.5) for traceability.

---
## Appendix G – Scope & Boundary Declaration
This appendix defines the authoritative scope boundaries for the evidence and control coverage described in this guide. Use it verbatim (with environment-specific substitutions) at the start of an audit cycle to suppress mis-scoping and to direct auditors to the correct authoritative systems.

### G.1 In-Scope Technical Components (Split: Platform Baseline vs RHACS Overlay)

#### G.1a Platform Baseline (OCP / RHCOS)
| Component | Scope Description | Control Surface (Representative) | Primary Evidence Artifacts |
|-----------|-------------------|----------------------------------|-----------------------------|
| OpenShift Kubernetes Control Plane (API Server, Scheduler, Controller Manager, etcd) | Cluster orchestration & API | RBAC, admission ordering context, audit/event emission | RBAC diff exports, selected audit log excerpts (forwarded), API server config (sanitized) |
| OpenShift RHCOS Nodes (Transactional OS) | Container execution substrate | OSTree signed images, MachineConfig-managed state, SELinux enforcing | MachineConfig diff history, OSTree commit IDs, SELinux enforcing status sample |
| Machine Config Operator (MCO) | Declarative node config reconciliation | Kernel args, file drop-ins, kubelet config | Signed MachineConfig YAML commits, MCO status reports |
| Compliance Operator (if deployed) | Benchmark scanning | CIS / custom profile rule evaluation | Scan summary report, pass/fail delta trend |
| Security Profiles Operator (SPO) (if deployed) | Seccomp / SELinux profile lifecycle | Custom profile distribution & attachment | Profile YAML + attachment manifests, audit log entries proving enforcement |
| Gatekeeper / OPA (if deployed) | Constraint-based policy | Naming/label/OPA invariants, custom admission checks | ConstraintTemplate & Constraint YAML (signed), violation events summary |
| Image Signature / Verification Config (ClusterImagePolicy, Quay, Sigstore) | Image trust & provenance enforcement | Signature & certificate trust roots, policy binding | ClusterImagePolicy export, verification logs, signer key inventory |
| Service Mesh / Ingress TLS Config (if enabled) | Transport security (north-south & east-west) | mTLS policy, TLS cipher config, route/ingress termination | Mesh PeerAuthentication/DestinationRule (STRICT), Ingress TLS cert inventory |
| Network & Segmentation Primitives | L3/L4 isolation & namespace scoping | NetworkPolicy enforcement, Multus/UDN segmentation | NetworkPolicy manifests, coverage % report, UDN inventory mapping |
| Resource Governance (Quotas/LimitRange) | Capacity & DoS resilience | Namespace quotas, limit ranges | Quota manifests + utilization report |

#### G.1b RHACS Overlay
| Component | Scope Description | Control Surface (Representative) | Primary Evidence Artifacts |
|-----------|-------------------|----------------------------------|-----------------------------|
| RHACS Central | Policy brain & API | Policy evaluation, compliance reporting, vuln data | Policy JSON exports (signed); external digests of compliance reports |
| RHACS Scanner / Scanner DB | Image & component analysis | Vulnerability & component inventory | Scan result exports, vuln trend metrics |
| RHACS Admission Controller (Validating Webhook) | Deploy-time gating | Misconfig, vuln, signature, risk-based deny/warn | Admission denial events, failurePolicy config snapshot |
| RHACS Sensor & Collector | Runtime & deploy telemetry | Process/network baselines, runtime policy triggers | Runtime alert logs, connected node coverage report |
| RHACS Notifiers (SIEM, Ticket, Chat) | External evidence & IR linkage | Alert forwarding & ticket creation | Notifier config snapshot, sample forwarded alert IDs |
| Secret Pattern Detection | Last-line exposure detection | Environment/config secret pattern matches | Secret violation trend report |
| Exception / Policy Phase Tracking | Progressive enforcement governance | Warn→Block transitions, time-bound exceptions | Exception register excerpt, policy phase labels export |
| Signature / Attestation Policy Integration | Enforce trusted images | Unsigned / untrusted image deny | Blocked deploy event, policy config + signer list |
| Vulnerability SLA Metrics | Risk reduction measurement | Age, fixability, SLA breach detection | Vulnerability SLA dashboard export, blocked image evidence |

### G.2 Explicit Out-of-Scope Areas
These areas are acknowledged dependencies or complementary controls but **not** in evidentiary scope for RHACS/OpenShift security enforcement in this guide. They must produce their own artifacts (tracked in External Control Register / Appendix E):
| Domain | Out-of-Scope Rationale | Primary System(s) | Required External Evidence |
|--------|-----------------------|-------------------|----------------------------|
| Enterprise IAM / MFA / SSO | User identity lifecycle & strong auth handled upstream | IdP (Keycloak, Okta, AAD, etc.) | MFA policy doc, IdP config snapshot, access review reports |
| Non-RHCOS Host OS Hardening / Bare Metal / Ancillary VMs | Guide centers on managed RHCOS nodes; other OS baselines differ | RHEL, Windows Server, Hypervisor layer | CIS benchmark reports, patch cadence, hardening scripts |
| Vault / Key Lifecycle Management | Secrets storage, rotation, escrow, key destruction | HashiCorp Vault, KMS (AWS KMS, Azure Key Vault), HSM | Key policy JSON, rotation logs, vault audit log excerpt |
| SAST / DAST / IaC Scanning | Application & infrastructure code analysis outside runtime/deploy gating | CI/CD security tools (SonarQube, CodeQL, Checkov, Trivy IaC) | Scan reports (external digest), remediation tickets, pipeline run IDs |
| Software Composition Analysis License Governance | Legal & license risk not enforced in RHACS | Dependency/license scanner | License scan diff, approval register |
| Backup & Disaster Recovery | Data/state resilience, restore validation | Backup platform, DR orchestration | DR test report, RPO/RTO metrics, backup integrity log |
| Business Continuity / BIA | Organizational process domain | GRC tooling | BIA document, review approval |
| WAF / API Gateway / L7 Threat Mitigation | Application-layer security beyond L3/L4 policy | API Gateway, WAF, CDN | WAF policy export, sampled blocked request logs |
| IDS / IPS / Deep Packet Inspection | Packet payload & advanced signature analysis | Network IDS/IPS, eBPF sensors | Alert sample, rule pack version, coverage map |
| SIEM Correlation & Advanced Analytics | Cross-domain event normalization & correlation logic | SIEM / UEBA platform | Correlation rule pack diff, suppression list, dashboard screenshot |
| Central Log Retention, WORM Storage | Long-term immutable storage & legal hold | Object store (S3 Object Lock, GCS), SIEM archive | Retention policy export, object lock configuration, external hash chain index |
| Data Encryption (At Rest & In Transit) Beyond Cluster Defaults | TLS termination, database/storage encryption lifecycle | Mesh, Ingress Controller, DB/KMS | TLS cipher policy, cert inventory, encryption enablement evidence |
| Incident Response Runbooks & Forensic Procedures | Human process & deep forensic tooling | IR platform, playbook repository | Runbook version hash, tabletop exercise report |
| Advanced Forensics & Memory Analysis | Memory/disk timeline, packet capture beyond RHACS telemetry | Forensics suite / EDR | Memory dump procedure, forensic artifact chain of custody |
| Host Vulnerability Management Outside Container Images | Kernel & package CVEs on host OS beyond image layer | Host scanning agents | Host vuln report, remediation SLA metrics |

### G.3 RHCOS Transactional (Controlled) vs “Immutable” Clarification
RHCOS (Red Hat Enterprise Linux CoreOS) is a **transactional, controlled** operating system managed via OSTree + MachineConfig Operator, not strictly immutable. Key assurance anchors:
1. **Signed Content:** OS updates delivered as signed OSTree commits (Red Hat content trust chain).
2. **Declarative State:** Desired node configuration declared via MachineConfig objects; drift visible and reconciled.
3. **Controlled Update Pipeline:** Cluster version operator orchestrates staged, verified rollouts (supports change evidence via version + commit IDs).
4. **SELinux Enforcing:** Mandatory access control assures workload confinement.
5. **No Assumption of Absolute Immutability:** Local alterations outside MachineConfig (emergency debug) must be treated as *exceptions* and remediated (or codified) quickly; evidence = diff + closure ticket.

Evidence Bundle (Example):
- MachineConfig YAML (signed commit) + associated OSTree commit IDs.
- `oc adm release info` output (release image signature) captured; external SHA-256 digest recorded.
- SELinux enforcing status sample across nodes.
- Exception log (if any) for manual node changes with remediation.

### G.4 Scope Statement (Sample Language)
Use this statement in audit introductions:
> The scope of container platform security evidence covers workload and cluster security controls enforced and/or evidenced by RHACS, OpenShift control plane components, transactional RHCOS nodes (via MachineConfig & OSTree), and designated Red Hat-supported security operators (Compliance Operator, Security Profiles Operator, Gatekeeper where deployed). Controls outside this boundary (enterprise IAM/MFA, key lifecycle, application-layer security testing, DR/backup, WAF/IDS, SIEM correlation logic, long-term log retention, data encryption lifecycle) are provided and evidenced by external systems referenced in the External Control Register.

### G.5 Boundary Validation Checklist (Quarterly)
| Check | Method | Pass Criterion | Exception Handling |
|-------|--------|---------------|-------------------|
| All nodes on expected OSTree commit set | Compare reported node OSTree commit IDs against approved release manifest | 100% match (allow controlled canary subset) | Log deviation → investigate → reconcile MachineConfig |
| MachineConfig drift | Review current rendered MachineConfig state vs version-controlled baseline | No unmanaged node file changes | Create remediation PR or exception entry |
| SELinux enforcing everywhere | Sample representative nodes to confirm SELinux enforcement status | All = Enforcing | Investigate node; restore enforcing & document |
| Unsupported manual changes (out-of-band node edits) | MachineConfig Operator rendered-state comparison plus (optional) compliance file rule and targeted node inspection | No unmanaged file drift; all nodes conform to rendered MachineConfig | Immediate cordon & investigate → revert or codify via MachineConfig; raise exception ticket (time‑bound); treat manual change as policy violation |
| Operator baselines intact (Compliance/SPO) | Confirm operators healthy; review compliance suites & remediations status; verify active security profiles match approved content digests (Git commit SHAs) in repo; ensure no failed checks or unmanaged local-only profiles | All relevant operator components healthy; every profile matches approved baseline; zero failed compliance checks | If drift or failure detected: raise exception, restore profile from source or update baseline via approved review, document closure |

### G.6 Handling Out-of-Scope Auditor Requests
Provide a polite redirect pattern:
| Request Type | Response Template |
|--------------|-------------------|
| “Show MFA enrollment statistics” | Outside RHACS/OpenShift scope; refer to IAM evidence bundle (IdP config + MFA policy + enrollment report). |
| “Provide WAF blocked request sample” | WAF is external; see External Control Register (WAF domain) for policy export & log sample location. |
| “Demonstrate database encryption keys” | Key management is external; provide KMS key policy & rotation logs referenced in External Control Register. |

### G.7 Change Control & Versioning of Scope
Version this Appendix (G) independently; any addition/removal of in-scope components requires:
1. Update Appendix G table(s)
2. Commit with signed tag referencing change ticket
3. Regenerate Appendix E entries if new external dependencies introduced
4. Notify audit preparation distribution list

> Principle: Scope drift without explicit versioning erodes evidence credibility—treat scope like code.

### G.8 Hashing Clarification (Tamper-Evidence vs Immutability)
In this guide, references to “hashing” or “hash digests” mean the external process of applying a cryptographic one-way function (e.g., SHA-256) to exported artifacts (compliance reports, logs, SBOMs, policy JSON, configuration snapshots).

- **Purpose:** Produces a unique fingerprint of each artifact. Any subsequent alteration changes the digest, providing **tamper-evidence**.
- **Limitation:** A digest alone does **not** make evidence immutable; it only detects modification after the fact.
- **Immutability Requirement:** True evidentiary immutability + retention (e.g., for AU-9 expectations) must be furnished by external controls: object-lock / WORM-capable storage, SIEM archival tiers, or compliant records management systems. RHACS and OpenShift *emit* artifacts; they do not perform or manage hashing or write-once retention internally.
- **Audit Mapping (AU-9):** External cryptographic digests satisfy the “integrity / tamper-detection” element (Partial = P). Immutable retention & legal hold are fully External (E) responsibilities.
- **Operational Practice:** Automate digest generation (e.g., pipeline job computing SHA-256, updating a signed index manifest) and store both artifact and digest in WORM/object-lock storage. Periodically reconcile index vs stored objects; alert on gaps or digest mismatches.

---
## 12. Extending & Maturing
| Phase | Focus | Additions |
|-------|-------|----------|
| Foundational | Visibility + Gate | Critical vuln & privileged config enforcement |
| Progressive | Segmentation + Least Privilege | 100% NetworkPolicy + RBAC diff reviews |
| Advanced | Automation & Response | Runtime enforcement, ticket & ChatOps hooks |
| Optimized | Metrics & Predictive | SLA trend KPIs, policy-as-code pipelines |

---
## 13. Common Pitfalls & Remedies
| Pitfall | Impact | Remedy |
|---------|--------|--------|
| Missing registry integration | Blind image coverage gaps | Add registry & rescan backlog |
| Policies only in “alert” mode | Drift & vulnerable images ship | Gradually enable enforcement (risk-ranked) |
| Excess cluster-admin bindings | High lateral compromise blast radius | Consolidate & implement quarterly review |
| Sparse NetworkPolicies | Lateral movement & exfil risk | Default deny + iterative allow modeling |
| Alert fatigue | Missed true positives | Measure precision; tune or suppress noisy patterns |
| No resource limits | Availability degradation | Enforce limits + quotas |
| Secrets in env vars | Credential theft risk | Vault mapping & policy block |

---
## 14. Maintaining the Guide
- Quarterly review for framework revisions (NIST & CIS updates)
- Regenerate mapping tables when new internal checks/policies added
- Treat policy set as code (signed commits, mandatory review)
- Track improvement KPIs: TTRc, NetworkPolicy coverage %, RBAC admin subject count, alert precision

---
## 15. Summary
Combining RHACS (visibility, policy gating, runtime telemetry, evidence exports) with OpenShift (SCC, NetworkPolicy, signature verification, admission & RBAC primitives) yields a continuously validated control stack covering major container security expectations across NIST 800‑53 and NIST 800‑190 (with alignment to CIS benchmarks). Focus on measurable reduction (vuln backlog, misconfig drift, alert noise) while maintaining tamper‑resistant, cryptographically verifiable evidence.

> Sustained compliance emerges from disciplined engineering feedback loops: enforce baselines, measure risk reduction, automate evidence, iterate.

---
