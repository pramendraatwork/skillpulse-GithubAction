# SkillPulse - Three Tier DevOps Project
# ⚡ SkillPulse

A production-grade **three-tier web application** built to demonstrate a complete DevOps pipeline — from local development to containerization, Kubernetes orchestration, and CI/CD automation.

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                    User Browser                      │
└─────────────────────┬───────────────────────────────┘
                      │ HTTP
┌─────────────────────▼───────────────────────────────┐
│             Frontend (NGINX + HTML)                  │
│                  Port 80                             │
│         Proxies /api/* → Backend                     │
└─────────────────────┬───────────────────────────────┘
                      │ HTTP :8080
┌─────────────────────▼───────────────────────────────┐
│              Backend (Go REST API)                   │
│                  Port 8080                           │
│         /api/skills   /api/health                    │
└─────────────────────┬───────────────────────────────┘
                      │ SQL :5432
┌─────────────────────▼───────────────────────────────┐
│             Database (PostgreSQL 15)                 │
│                  Port 5432                           │
└─────────────────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | HTML, CSS, JavaScript, NGINX |
| **Backend** | Go (Golang) 1.21 |
| **Database** | PostgreSQL 15 |
| **Containerization** | Docker, Docker Compose |
| **Orchestration** | Kubernetes (Kind cluster) |
| **CI/CD** | GitHub Actions |
| **Registry** | Docker Hub |
| **IaC** | Terraform (AWS EKS ready) |

---

## 📁 Project Structure

```
skillpulse/
├── backend/                  # Go REST API
│   ├── main.go               # API handlers + DB connection
│   ├── go.mod                # Go module definition
│   ├── go.sum                # Dependency lock file
│   └── Dockerfile            # Multi-stage build
├── frontend/                 # Static UI + NGINX
│   ├── index.html            # Dashboard UI
│   ├── nginx.conf            # Reverse proxy config
│   └── Dockerfile            # NGINX image
├── k8s/                      # Kubernetes manifests
│   ├── namespace.yaml        # skillpulse namespace
│   ├── postgres.yaml         # DB deployment + service + PVC
│   ├── backend.yaml          # Backend deployment + service
│   ├── frontend.yaml         # Frontend deployment + NodePort
│   └── hpa.yaml              # Horizontal Pod Autoscaler
├── terraform/                # AWS EKS infrastructure
│   ├── main.tf               # VPC + EKS cluster
│   └── variables.tf          # Input variables
├── .github/
│   └── workflows/
│       └── deploy.yml        # CI/CD pipeline
├── docker-compose.yml        # Local development
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Kind (`brew install kind` or see [kind.sigs.k8s.io](https://kind.sigs.k8s.io))
- kubectl
- Go 1.21+ (for local development)

---

### Phase 1 — Run Locally with Docker Compose

```bash
# Clone the repo
git clone https://github.com/pramendraatwork/skillpulse.git
cd skillpulse

# Start all three tiers
docker compose up --build

# Open in browser
http://localhost
```

---

### Phase 2 — Deploy to Kubernetes (Kind)

```bash
# Create Kind cluster
kind create cluster --name skillpulse

# Set kubeconfig
kind export kubeconfig --name skillpulse

# Build and load images
docker build -t skillpulse-skillpulse-backend:latest ./backend
docker build -t skillpulse-skillpulse-frontend:latest ./frontend
kind load docker-image skillpulse-skillpulse-backend:latest --name skillpulse
kind load docker-image skillpulse-skillpulse-frontend:latest --name skillpulse

# Deploy
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/

# Watch pods
kubectl get pods -n skillpulse -w

# Access the app
kubectl port-forward svc/skillpulse-frontend 8090:80 -n skillpulse
# Open http://localhost:8090
```

---

### Phase 3 — CI/CD with GitHub Actions

Every push to `main` automatically:
1. Runs Go tests
2. Builds Docker images
3. Pushes to Docker Hub

**Required GitHub Secrets:**

| Secret | Description |
|---|---|
| `DOCKERHUB_USER` | Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token (Read/Write/Delete) |

---

### Phase 4 — AWS EKS with Terraform (Cloud)

```bash
cd terraform
terraform init
terraform plan
terraform apply

# Update kubeconfig
aws eks update-kubeconfig --name skillpulse-cluster --region us-east-1

# Deploy
kubectl apply -f ../k8s/
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/health` | Health check |
| GET | `/api/skills` | Get all skills |

### Example Response

```json
[
  {"id": 1, "name": "Docker", "category": "DevOps", "level": "Intermediate"},
  {"id": 2, "name": "Kubernetes", "category": "DevOps", "level": "Beginner"},
  {"id": 3, "name": "Terraform", "category": "IaC", "level": "Advanced"}
]
```

---

## ⚙️ Kubernetes Resources

| Resource | Count | Purpose |
|---|---|---|
| Deployments | 3 | postgres, backend, frontend |
| Services | 3 | Internal + NodePort routing |
| ConfigMap | 1 | DB configuration |
| Secret | 1 | DB credentials |
| PVC | 1 | Persistent storage for DB |
| HPA | 1 | Auto-scales backend 2→5 pods |

---

## 🔧 Useful Commands

```bash
# View all resources
kubectl get all -n skillpulse

# Check logs
kubectl logs -l app=skillpulse-backend -n skillpulse
kubectl logs -l app=skillpulse-frontend -n skillpulse

# Shell into a pod
kubectl exec -it -n skillpulse <pod-name> -- sh

# Scale manually
kubectl scale deployment skillpulse-backend --replicas=3 -n skillpulse

# Delete everything
kubectl delete namespace skillpulse
kind delete cluster --name skillpulse
```

---

## 🗺️ Roadmap

- [ ] AWS EKS deployment via Terraform
- [ ] NGINX Ingress Controller
- [ ] cert-manager for HTTPS
- [ ] Prometheus + Grafana monitoring
- [ ] Helm chart packaging
- [ ] Multi-environment (staging/prod)

---

## 👨‍💻 Author

**Pramendra** — [@pramendraatwork](https://github.com/pramendraatwork)

Built as a hands-on DevOps learning project covering Docker, Kubernetes, GitHub Actions, and cloud-native deployment patterns.