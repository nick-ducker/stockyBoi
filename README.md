# stockyBoi

A way for me to not forget about all my "very" important investments.

Some key commands so I don't forget:

**Build the docker image for ARM64 arch and push to the github registry**
    
    `docker buildx build --platform linux/amd64,linux/arm64 -t "ghcr.io/nick-ducker/stockyboi" . --push`

**Restart the deployment so k8s will pull the new image**

    `kubectl rollout restart deploy <name>`

**Figure out which port the dang thing is running on**

    `kubectl get svc <deployment name>`