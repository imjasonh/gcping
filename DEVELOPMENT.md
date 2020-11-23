# Development

Deploy Cloud Run site using Terraform

### Prerequisites

Install [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli) and [`ko`](https://github.com/google/ko)

### Deploy using Terraform

```
$ terraform init
$ terraform apply -var image=$(ko publish ./cmd/ping/)
```

This deploys the ping service to all Cloud Run regions and configures a global HTTPS Load Balancer with Google-managed SSL certificate for `global.gcping.com`.

### Run locally

```
docker run -p 8080:8080 $(ko publish ./cmd/ping/)
```

And browse to http://localhost:8080/
