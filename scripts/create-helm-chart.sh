#!/bin/bash

# Define variables
OPERATOR_SDK_DIR=$PWD/../distributed-job-scheduler-operator
HELM_CHART_DIR=$PWD/../helm
HELM_CHART_NAME="distributed-job-scheduler-operator"

# Create the Helm chart structure
mkdir -p "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates"
mkdir -p "${HELM_CHART_DIR}/${HELM_CHART_NAME}/crds"

# Copy CRDs from Operator SDK project to Helm chart CRDs directory
cp "${OPERATOR_SDK_DIR}/config/crd/bases/"*.yaml "${HELM_CHART_DIR}/${HELM_CHART_NAME}/crds/"

# Copy RBAC manifests
cp "${OPERATOR_SDK_DIR}/config/rbac/"*.yaml "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/"

# Copy deployment manifests
cp "${OPERATOR_SDK_DIR}/config/manager/manager.yaml" "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/deployment.yaml"

# Replace static values with Helm template placeholders
sed -i '' 's|replicas: [0-9]*|replicas: {{ .Values.replicaCount }}|' "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/deployment.yaml"
sed -i '' 's|image:.*|image: {{ .Values.image.repository }} : {{ .Values.image.tag }}|' "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/deployment.yaml"
# Add new resources block with placeholders
sed -i '' '/image: {{ .Values.image.repository }}:{{ .Values.image.tag }}/a\
\        resources:\
\          limits:\
\            cpu: {{ .Values.resources.limits.cpu }}\
\            memory: {{ .Values.resources.limits.memory }}\
\          requests:\
\            cpu: {{ .Values.resources.requests.cpu }}\
\            memory: {{ .Values.resources.requests.memory }}' "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/deployment.yaml"

# Copy service account manifests
cp "${OPERATOR_SDK_DIR}/config/manager/service_account.yaml" "${HELM_CHART_DIR}/${HELM_CHART_NAME}/templates/serviceaccount.yaml"

# Create Helm's Chart.yaml file
cat <<EOF > "${HELM_CHART_DIR}/${HELM_CHART_NAME}/Chart.yaml"
apiVersion: v2
name: ${HELM_CHART_NAME}
description: A Helm chart for deploying ${HELM_CHART_NAME} operator
version: 0.1.1
appVersion: "1.0.0"
EOF

# Create Helm's values.yaml file
cat <<EOF > "${HELM_CHART_DIR}/${HELM_CHART_NAME}/values.yaml"
replicaCount: 1

image:
  repository: your-operator-image
  tag: latest
  pullPolicy: Always

serviceAccount:
  create: true

rbac:
  create: true

resources: {}
EOF

# Inform the user that the script has completed
echo "Helm chart for ${HELM_CHART_NAME} has been created at ${HELM_CHART_DIR}/${HELM_CHART_NAME}"
