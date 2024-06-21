CICD Development Testing
========================

This repository includes a number of github action workflows to enable functional verification of the Cloudzero Charts.

## Prerequisites

GitHub Actions is a powerful tool that allows developers to automate, customize, and execute their software development workflows directly in the GitHub repository. This guide will help you understand how to write and test GitHub Action workflows and use [docker](https://docs.docker.com/desktop/install/mac-install/) and the [act utility](https://github.com/nektos/act) for local development and testing.

## List All Workflows:


```bash
act -l
```

**Example output**

```bash
Stage  Job ID                                        Job name                                      Workflow name             Workflow file                 Events           
0      build-and-publish-chart                       build-and-publish-chart                       build-and-publish-chart   build-and-publish-chart.yml   push             
0      has_changes                                   has_changes                                   build_test_publish_image  build-test-publish-image.yml  push,pull_request
0      scanner                                       scanner                                       detection_rules           change-detector.yml           workflow_call    
1      test_chart_lint                               test_chart_lint                               build_test_publish_image  build-test-publish-image.yml  pull_request,push
1      build_test_chart_install_maybe_publish_image  build_test_chart_install_maybe_publish_image  build_test_publish_image  build-test-publish-image.yml  push,pull_request
```

## Verify a Specific Workflow (DRY-RUN mode):

```bash
act --dry-run -j build_test_chart_install_maybe_publish_image -s CLOUDZERO_API_TOKEN=$CZ_API_TOKEN -a $GITHUB_USER --secret GITHUB_TOKEN=$GITHUB_TOKEN
```

> add the `-n` or `--dry-run` flag only validates the jobs are syntactically correct

## Run a Specific Workflow:

```bash
act --dry-run -j build_test_chart_install_maybe_publish_image -s CLOUDZERO_API_TOKEN=$CZ_API_TOKEN -a $GITHUB_USER --secret GITHUB_TOKEN=$GITHUB_TOKEN
```

