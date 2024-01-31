# OminousPositivity

## Overview

A toy project inspired by a funny picture with a quote purporting "Ominous Positivity", e.g "You will be okay, you have no choice."

The goal of making this project was to reinforce fundamentals and practice with the following concepts/tech

* Go
  * Unit testing
  * Integration testing
* JavaScript (React)
  * Unit testing with Jest
* Terraform
* AWS
  * Lambda APIGW
  * DynamoDB
  * CloudFront
  * IAM
  * CloudWatch
  * X-Ray
  * S3
  * Route53
  * ACM

* Logging
* Metrics
* Tracing
* Content Delivery Networks 
* Infrastructure as Code (Automation)
* Serverless Concepts
* NoSQL/Single Table DynamoDB Design

The project includes a SAM template, allowing for local development of a 'serverless' application meant to be hosted in AWS APIGW.

DynamoDB Local can be ran in a Docker container or via NoSQL WorkBench.

Additionally the project has a lot of key focus on DevOps/SRE principles.
> The choice of architecture allows it to handle large volumes of traffic without any manual intervention

The end result is something that is:
* Scalable
* Highly Available
* Resilient
* Reliable
* Monitored
* Observable


## Considerations

There is a large amount of additional configuration that could be added.
For example the Terraform modules could have outputs or great flexibility.

> Running the Terraform may result in billing costs