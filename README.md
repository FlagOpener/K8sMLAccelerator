# K8sMLAccelerator: Kubernetes Native Accelerator for Big Data Processing in Machine Learning Applications

## Why?

Machine Learning (ML)/Deep Learning (DL) training using datasets stored at S3/GCS/Azure can experience rate limiting and suboptimal downloading throughput. K8sMLAccelerator autorewrites training job's Pod spec and reroutes S3/GCS/Azure requests to a local cache, thus bolstering performance and scalability.

## How It Works

This is achieved via a mutating webhook that modifies the Pod spec. The webhook deploys host aliases to Pod spec once it identifies the `app` label in the deployments or jobs. Upon running 