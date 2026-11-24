# Travel Portal (presentation/demo app)

Lightweight Go + SQLite demo app with a static frontend branded as "Travel Portal." This repo contains a Dockerfile and a Helm chart so you can run locally in Docker or deploy to Kubernetes.

Quick notes:
- Default app port: `7540` (configurable via `TRAVELPORTAL_PORT`)
- Default password (container/CI): `test12345` (configurable via `TRAVELPORTAL_PASSWORD`)
- Health endpoint: `GET /health` (used for probes)

## Requirements
- Docker (build & run)
- Helm (for Kubernetes deploy)
- Kubernetes cluster (optional)

## Build & Run (local)

Build and run with Make (image will be built locally):

```bash
make docker-build REGISTRY=localhost TAG=local
docker run --rm -p 7540:7540 -e TRAVELPORTAL_PASSWORD=test12345 localhost/travelportal:local
```

Open http://localhost:7540 in your browser and sign in using the password in `TRAVELPORTAL_PASSWORD`.

Notes for CI systems: the Dockerfile is designed to be CI-friendly (sets `GOPROXY`, enables modules, caches `go mod download`, and builds in a multi-stage image). The runtime image runs as a non-root `app` user.

## Deploy to Kubernetes

The Helm chart is in `helm-charts/travelportal`. Set `env.travelportalPassword` to populate the `travelportalpassword` key in the chart-generated Secret, named `<fullname>-auth`.

Recommended deploy (example using `helm` flags):

```bash
helm upgrade --install travelportal ./helm-charts/travelportal \
	--namespace demo --create-namespace \
	--set image.repository=<registry>/<image> \
	--set image.tag=<tag> \
	--set-string env.travelportalPassword="your-secure-password"
```

If you deployed the previous chart as `todo-app`, the command above installs a separate release. Plan the switch to the new Service and carry over any persistent storage you use before retiring the old release.

The app still accepts legacy `TODO_PORT`, `TODO_DBFILE`, and `TODO_PASSWORD` variables when their corresponding `TRAVELPORTAL_*` variables are unset. The chart accepts `env.todoPassword` when `env.travelportalPassword` is empty.

Persistence: the app uses SQLite with `scheduler.db` by default. Override the file path with `TRAVELPORTAL_DBFILE`. For production-like runs, mount a PVC using `volumeMounts` and `volumes` in the chart values.

Probes: liveness/readiness probe paths are set to `/health` in `helm-charts/travelportal/values.yaml`.

Security & runtime recommendations:
- Provide a non-default, strong password through the Helm values.
- Set resource requests/limits in `values.yaml`.
- Consider running the container with a read-only root filesystem and explicit UID/GID via `securityContext` in values.

## CI integration

The Harness configuration in `pipeline.yaml` builds and deploys the `travelportal` image. For other CI systems, the basic steps are:

```bash
# build
docker build -t <registry>/<image>:<tag> .
# push
docker push <registry>/<image>:<tag>
# optional: helm upgrade --install (requires kubeconfig in CI)
```

If you want, I can add a sample `gitlab-ci.yml`, `drone.yml`, or a minimal Jenkinsfile tailored to your CI provider.

## Where to look
- Dockerfile: `Dockerfile`
- Helm chart: `helm-charts/travelportal/`
- Web UI: `web/` (login at `login.html`, main app at `/`)

## Updates 
- Enhance login
- Enable Secret Injection
