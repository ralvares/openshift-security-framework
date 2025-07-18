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