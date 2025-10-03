# ---------- helpers -------------------------------------------------
def ns_default:
  (.objectRef.namespace // "") == "default";

def fmt($msg):
  {alert:$msg,
   ts:.requestReceivedTimestamp,
   user:(.user.username // "N/A"),
   verb:.verb,
   res:.objectRef.resource,
   sub:.objectRef.subresource,
   name:(.objectRef.name // ""),
   uri:.requestURI};

# ---------- rules ---------------------------------------------------
# 0. Pod created in default
def demo_pod_create:
  select(ns_default and .objectRef.resource=="pods" and .verb=="create")
  | fmt("DEMO – Pod created (attacker foothold)");

# 1. Exec into pod
def demo_exec:
  select(ns_default and .objectRef.subresource=="exec" and .verb=="get")
  | fmt("DEMO – Interactive exec in pod");

# 2. ClusterRoleBinding escalation
def demo_crb_escalate:
  select(.objectRef.resource=="clusterrolebindings" and
         (.verb|test("^(create|patch|update)$")))
  | fmt("DEMO – SA escalated to cluster-admin");

# 3a. API list/get via ServiceAccount
def demo_sa_api_abuse:
  select(.user.username|startswith("system:serviceaccount:default:") and
         (.verb|test("^(get|list)$")))
  | fmt("DEMO – SA bulk API access");

# 3b. Secrets list
def demo_secrets_list:
  select(ns_default and .objectRef.resource=="secrets" and .verb=="list")
  | fmt("DEMO – Secrets enumeration");

# 4. CronJob creation
def demo_cronjob:
  select(ns_default and .objectRef.resource=="cronjobs" and
         (.verb|test("^(create|patch|update)$")))
  | fmt("DEMO – Persistence via CronJob");

# 5. Port-forward
def demo_portforward:
  select(ns_default and .objectRef.subresource=="portforward" and .verb=="get")
  | fmt("DEMO – Port-forward established");

# ---------- dispatcher ----------------------------------------------
(demo_pod_create,
 demo_exec,
 demo_crb_escalate,
 demo_sa_api_abuse,
 demo_secrets_list,
 demo_cronjob,
 demo_portforward)