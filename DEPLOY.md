# SkillPulse — Complete Deployment Guide
## Three-Tier App: Go + NGINX + PostgreSQL on Minikube

---

## 📁 Project Structure
```
skillpulse/
├── backend/          # Go REST API
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── frontend/         # HTML + NGINX
│   ├── index.html
│   ├── nginx.conf
│   └── Dockerfile
├── k8s/              # Kubernetes manifests
│   ├── namespace.yaml
│   ├── postgres.yaml
│   ├── backend.yaml
│   ├── frontend.yaml
│   └── hpa.yaml
├── terraform/        # AWS EKS (future)
│   ├── main.tf
│   └── variables.tf
├── .github/workflows/
│   └── deploy.yml    # CI/CD pipeline
├── docker-compose.yml
└── DEPLOY.md
```

---

## PHASE 1 — Run Locally with Docker Compose

### Step 1: Clone / create the project
```bash
cd skillpulse
```

### Step 2: Start everything
```bash
docker compose up --build
```

### Step 3: Test it
```bash
# Frontend
open http://localhost

# Backend API
curl http://localhost:8080/api/skills
curl http://localhost:8080/api/health
```

### Step 4: Stop
```bash
docker compose down -v
```

---

## PHASE 2 — Deploy to Minikube (Kubernetes)

### Prerequisites
```bash
minikube version   # must be installed
kubectl version    # must be installed
docker version     # must be installed
```

### Step 1: Start Minikube
```bash
minikube start --driver=docker --memory=4096 --cpus=2
```

### Step 2: Point Docker to Minikube's daemon
# This lets Minikube use locally built images (no registry needed)
```bash
eval $(minikube docker-env)
```

### Step 3: Build images INSIDE Minikube's Docker
```bash
docker build -t skillpulse-backend:latest ./backend
docker build -t skillpulse-frontend:latest ./frontend
```

### Step 4: Apply all Kubernetes manifests
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/hpa.yaml
```

### Step 5: Watch pods come up
```bash
kubectl get pods -n skillpulse -w
# Wait until all pods show STATUS = Running
```

### Step 6: Open the app
```bash
minikube service skillpulse-frontend -n skillpulse
# OR get the URL manually:
minikube ip   # e.g. 192.168.49.2
# Then open: http://192.168.49.2:30080
```

### Step 7: Check logs if anything fails
```bash
kubectl logs -l app=skillpulse-backend -n skillpulse
kubectl logs -l app=skillpulse-frontend -n skillpulse
kubectl describe pod <pod-name> -n skillpulse
```

### Step 8: Scale manually
```bash
kubectl scale deployment skillpulse-backend --replicas=3 -n skillpulse
```

---

## PHASE 3 — GitHub Actions CI/CD

### Step 1: Push project to GitHub
```bash
git init
git add .
git commit -m "feat: initial SkillPulse setup"
git remote add origin https://github.com/YOUR_USERNAME/skillpulse.git
git push -u origin main
```

### Step 2: Add GitHub Secrets
Go to: GitHub repo → Settings → Secrets and variables → Actions

| Secret Name      | Value                          |
|------------------|-------------------------------|
| DOCKERHUB_USER   | your Docker Hub username       |
| DOCKERHUB_TOKEN  | Docker Hub access token        |
| KUBE_CONFIG      | base64-encoded kubeconfig      |

Get your kubeconfig as base64:
```bash
cat ~/.kube/config | base64 -w 0
```

### Step 3: Push a change to trigger the pipeline
```bash
# Make any change, then:
git add .
git commit -m "trigger: test CI/CD pipeline"
git push
```

Watch it run at: GitHub → Actions tab

---

## PHASE 4 — Terraform (AWS EKS, future)

```bash
cd terraform
terraform init
terraform plan
terraform apply   # provisions VPC + EKS on AWS
```

Then update kubeconfig:
```bash
aws eks update-kubeconfig --name skillpulse-cluster --region us-east-1
kubectl apply -f ../k8s/
```

---

## Useful Commands Cheatsheet

```bash
# Minikube
minikube start / stop / delete
minikube dashboard          # visual UI
minikube service <name> -n skillpulse

# Kubernetes
kubectl get all -n skillpulse
kubectl get pods -n skillpulse
kubectl describe pod <name> -n skillpulse
kubectl logs <pod> -n skillpulse
kubectl exec -it <pod> -n skillpulse -- sh

# Docker
docker compose up --build
docker compose down -v
docker images
docker ps
```
