# Drone Junit

[![Release](https://github.com/rohit-gohri/drone-junit/actions/workflows/release.yaml/badge.svg)](https://github.com/rohit-gohri/drone-junit/actions/workflows/release.yaml)
[![Test Build](https://cloud.drone.io/api/badges/rohit-gohri/drone-junit/status.svg?ref=refs/heads/main)](https://cloud.drone.io/rohit-gohri/drone-junit)

A Drone plugin to parse Junit test reports and create tests summary using [plugin cards](https://docs.drone.io/plugins/adaptive_cards/).

![Preview of junit report](./docs/image.jpg)

## Usage

The following settings changes this plugin's behavior.

* paths (required) - Pass a glob pattern to match all xml junit files
* report_name (optional) - Customize the name of the report
* total (optional, default true) - Combine and show sum of all reports

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:

# Run your tests here and generate report
- name: tests
  image: golang
  commands:
    - go install github.com/jstemmer/go-junit-report/v2@latest
    - go test -v 2>&1 ./... | go-junit-report -set-exit-code > report.xml

- name: junit-reports
  image: ghcr.io/rohit-gohri/drone-junit:v0
  settings:
    paths: report.xml
  when:
    status: [
      'success',
      'failure',
    ]
```

### GHCR

Images are published on Github Container Registry - <https://github.com/rohit-gohri/drone-junit/pkgs/container/drone-junit>

```yaml
- name: junit-reports
  image: ghcr.io/rohit-gohri/drone-junit:v0
  settings:
    paths: report.xml
  when:
    status: [
      'success',
      'failure',
    ]
```

### Docker Hub

Images are also published on Docker Hub - <https://hub.docker.com/r/boringdownload/drone-junit>

```yaml
- name: junit-reports
  image: boringdownload/drone-junit:v0
  settings:
    paths: report.xml
```

## Development

### Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t boringdownload/drone-junit -f docker/Dockerfile .
```

### Testing

Execute the plugin from your current working directory:

```text
docker run --rm -e PLUGIN_PATHS="report.xml" -e PLUGIN_REPORT_NAME="drone-junit" \
  -e DRONE_COMMIT_SHA=8f51ad7884c5eb69c11d260a31da7a745e6b78e2 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_BUILD_NUMBER=43 \
  -e DRONE_BUILD_STATUS=success \
  -w /drone/src \
  -v $(pwd):/drone/src \
  boringdownload/drone-junit
```
