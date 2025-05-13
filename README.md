# Istio release
This repository contains release-ready Istio helm chart tarballs used for downstream builds of Sail operator and tooling to build them.

It's used in cases when we need helm chart changes which are not yet available upstream or not available for given Istio version.

## Usage
1. Create a new directory containing `manifest.yaml`
   - Directory name should follow `<istio_version>-redhat` pattern, e.g. `1.24.5-redhat`
   - existing `manifest.yaml` can be copied over and edited to match the directory version
1. Make sure the referenced repo/branch in `manifest.yaml` contains all expected changes before building the helm charts
1. Build the charts (you need to pass the directory name)
   - `go run . 1.24.5-redhat`
1. tarballs will be created in `helm` directory, e.g. `1.24.5-redhat/helm/`
1. Push the changes
1. Update the Sail operator to use the new tarball versions, e.g. see this [PR](https://github.com/openshift-service-mesh/sail-operator/pull/306)