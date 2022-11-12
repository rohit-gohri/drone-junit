A plugin to Drone plugin to create stest summaryy using plugin cards.

# Usage

The following settings changes this plugin's behavior.

* param1 (optional) does something.
* param2 (optional) does something different.

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: run boringdownload/drone-junit plugin
  image: boringdownload/drone-junit:v0
  pull: if-not-exists
  settings:
    param1: foo
    param2: bar
```

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t boringdownload/drone-junit -f docker/Dockerfile .
```

# Testing

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
