# Steps to Create a Helm Chart Repository
> To convert your existing Helm chart directory into a repository that can be used as a Helm chart repo (so others can install it via helm repo add and helm install)

---

## 1. Package your Helm chart

Navigate to your chart directory and use `helm package`:

```bash
cd helm/distributed-job-scheduler-operator
helm package .
```
This creates a `.tgz` file (e.g., `distributed-job-scheduler-operator-<version>.tgz`).

---

## 2. Move the chart package to a charts directory (optional but recommended)

Create a `charts` directory at the root of your repo and move the `.tgz` file there:

```bash
mkdir -p ../../charts
mv distributed-job-scheduler-operator-*.tgz ../../charts/
cd ../../charts
```

---

## 3. Generate an index.yaml for your chart repo

Use the `helm repo index` command to generate or update the Helm repository index file:

```bash
helm repo index . --url https://github.com/vishalanarase/openinnovationai/releases/latest/download
```

- If you want to serve from the `main` branch via GitHub Pages, use:
  ```
  helm repo index . --url https://vishalanarase.github.io/openinnovationai/charts
  ```

---

## 4. Push the `charts/` directory (containing `.tgz` and `index.yaml`) to your GitHub repository

Add, commit, and push:

```bash
git add charts/
git commit -m "Add packaged helm chart and repo index"
git push
```

---

## 5. (Recommended) Serve your chart via GitHub Pages

- Go to your repo’s settings.
- Under "Pages", set the source to the `charts/` directory on the `main` branch.
- Your Helm repo URL will then be:  
  `https://vishalanarase.github.io/openinnovationai/charts`

---

## 6. Add your Helm repo and install the chart

On any machine:

```bash
helm repo add openinnovationai https://vishalanarase.github.io/openinnovationai/charts
helm repo update
helm search repo openinnovationai
helm install my-job-scheduler openinnovationai/distributed-job-scheduler-operator
```

---

### **Summary**

- Package your chart (`helm package .`)
- Create `charts/`, move package there
- Generate `index.yaml` (`helm repo index . --url ...`)
- Push to GitHub, enable Pages if desired
- Add repo with `helm repo add ...`, then install

---