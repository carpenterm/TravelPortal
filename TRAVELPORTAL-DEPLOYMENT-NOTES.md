# Travel Portal rename and deployment notes

Travel Portal was originally deployed as `todo-app`. The repo now uses `travelportal` for the built executable, container image, Helm chart, and default Helm release.

## Changes in the repo

| Area | Change |
| --- | --- |
| Helm chart | Moved the chart out of the `helm-charts/` root into `helm-charts/travelportal/`. The folder matches `name: travelportal` in [Chart.yaml](helm-charts/travelportal/Chart.yaml). The chart version is now `0.2.0`. The move includes the templates, default and Harness values, and the regenerated sample manifest. |
| Go | Changed the module and internal imports from `github.com/antonkazachenko/go-todo-list-api` to `github.com/skandpurohit/TravelPortal`. The integration tests accept the renamed port and database variables. |
| Docker | The [Dockerfile](Dockerfile) builds and launches `/app/travelportal`. Its image source label now points to the TravelPortal repository. |
| Deployment defaults | The [Makefile](Makefile) and [deploy script](deploy/helm-deploy.sh) default to image and release `travelportal`, chart path `helm-charts/travelportal`, and namespace `default`. The default image is `quay.io/rahulk10/travelportal:latest`. The [README](README.md) uses the new names. |
| Harness | The build/push and deployment artifact paths in [pipeline.yaml](pipeline.yaml) now end in `/travelportal`. The [Harness values](helm-charts/travelportal/values-harness.yaml) also pull `/travelportal` from the packages registry. The build still tags the image with the pipeline execution ID. |
| Istio preview | [istio-preview-template.yaml](istio-preview-template.yaml) routes to `travelportal.default.svc.cluster.local` on port `80`. |

## Runtime configuration

| Setting | Preferred variable | Legacy fallback | Default |
| --- | --- | --- | --- |
| HTTP port | `TRAVELPORTAL_PORT` | `TODO_PORT` | `7540` |
| SQLite file | `TRAVELPORTAL_DBFILE` | `TODO_DBFILE` | `scheduler.db` |
| Sign-in password | `TRAVELPORTAL_PASSWORD` | `TODO_PASSWORD` | The Docker image supplies a demo default through the legacy variable. The Go app can also read the Vault-injected file at `/vault/secrets/config.txt`. |

The legacy environment variables apply when the corresponding new variable is unset. Helm uses `env.travelportalPassword` and accepts `env.todoPassword` when the new value is empty. The chart creates `<fullname>-auth` with key `travelportalpassword` and injects it as `TRAVELPORTAL_PASSWORD`. With the default release, that Secret is `travelportal-auth`. The Makefile and deploy script pass image overrides only. To set a password at deploy time, use the Helm command in the README or configure the values in your deployment system.

## Checks completed locally

- Built the Go executable with Go 1.25. Started it with a temporary SQLite file and verified `/health`, the branded login page, and successful sign-in. Checked both the new settings and the legacy fallbacks, including precedence when both were supplied.
- Ran Helm lint for default and Harness values, rendered the chart with both password keys, and packaged it as `travelportal-0.2.0.tgz`. Checked that the Makefile and deploy script use the renamed chart path.
- The only failing Go integration test was `TestApp`, which also failed against the original revision. Its unauthenticated request for `index.html` follows the redirect to the login page and compares that response with the dashboard file.
- Docker image build and container runtime verification were blocked by an unavailable Docker daemon. A live Kubernetes rollout was unverified.

## If Docker fails

Run these from the repository root. The local run uses the image's demo password default.

```bash
docker info
make docker-build REGISTRY=localhost TAG=local
docker run --rm -p 7540:7540 localhost/travelportal:local
# In another terminal
curl -fsS http://localhost:7540/health
```

- For build failures, check the first failing Dockerfile step. The builder uses Go 1.25 and needs registry access, Go module downloads through `GOPROXY`, and Alpine build packages for CGO and SQLite.
- For startup failures, inspect the container logs. Check that the image contains `/app/travelportal` and that the non-root container user can write to the directory holding the SQLite file, especially if you mounted a volume.

## If Kubernetes fails

Render the manifests with a placeholder first. For Harness, repeat with `-f helm-charts/travelportal/values-harness.yaml` and check the image and tag after Harness resolves them. Run the cluster commands in the intended kube context and set `namespace` and `release` to the values used for deployment.

```bash
helm lint ./helm-charts/travelportal --set-string env.travelportalPassword=render-only
helm template travelportal ./helm-charts/travelportal --set-string env.travelportalPassword=render-only

namespace=default
release=travelportal
helm status "$release" -n "$namespace"
kubectl get deploy,svc,pods -n "$namespace" -l "app.kubernetes.io/instance=$release"
kubectl get events -n "$namespace" --sort-by=.lastTimestamp
kubectl logs -n "$namespace" -l "app.kubernetes.io/instance=$release" --tail=100
```

- **Image pull errors.** Verify the exact `travelportal:<tag>` image was pushed to the registry the pod uses. Check registry permissions and `imagePullSecrets` for private images. If Harness still references `helm-charts/`, check the manifest path configured in the external service `spservice`.
- **Pod startup or sign-in errors.** Inspect pod events and check that the Deployment references the generated Secret and its `travelportalpassword` key. Check that the database path is writable. The Service targets container port `7540` and probes use `/health`. A custom app port also needs corresponding changes in the Kubernetes manifests.
- **Preview routing errors.** Check the actual Service name and namespace, especially if you changed the release name or set a chart name override. The preview route assumes Service `travelportal` in namespace `default`.
- **Existing `todo-app` release.** Installing with the new default release name creates a separate release. Point traffic at the new Service after checking readiness and plan to reuse or migrate any persistent SQLite data before retiring the old release.
