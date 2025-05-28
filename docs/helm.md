# Steps to Create a Helm Chart Repository
> To convert your existing Helm chart directory into a repository that can be used as a Helm chart repo (so others can install it via helm repo add and helm install)

If you **don’t want to use GitHub Pages** and just want to use your GitHub repo as a Helm repository **directly from the repo files** (for example, the `charts/` directory in your main/master branch), you can still use the raw file URLs to serve your Helm repo.

Here’s how you can do it:

---

## 1. Package your chart and create the index

```bash
cd helm/distributed-job-scheduler-operator
helm package .
mkdir -p ../../charts
mv distributed-job-scheduler-operator-*.tgz ../../charts/
cd ../../charts
helm repo index . --url https://raw.githubusercontent.com/vishalanarase/openinnovationai/master/charts
```

This will generate an `index.yaml` pointing to the raw GitHub URLs.

---

## 2. Push the `charts/` directory and its contents (`.tgz` and `index.yaml`) to your repo

```bash
git add charts/
git commit -m "Add Helm chart and index"
git push
```

---

## 3. Add the repo using the **raw GitHub URL**:

```bash
helm repo add openinnovationai https://raw.githubusercontent.com/vishalanarase/openinnovationai/master/charts
helm repo update
helm search repo openinnovationai
helm install my-job-scheduler openinnovationai/distributed-job-scheduler-operator
```

---

### Notes
- **raw.githubusercontent.com** serves static files from your repo, so `index.yaml` and `.tgz` files are accessible.
- This method works for **public repos only**.
- This is the most common approach for “Helm chart repos” that don’t use GitHub Pages or a custom web server.

---

### Summary

- Use `https://raw.githubusercontent.com/<owner>/<repo>/<branch>/charts` as your Helm repo URL
- Make sure `index.yaml` and your packaged chart are present in that directory and committed

---