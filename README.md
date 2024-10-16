# k8sdemo (WIP)

## Running with Docker

### Building demo app images

You can build individual images using:

```bash
make build demoapp1
make build demoapp2
make build demoapp3
```

Alternatively, you can build all demo images at once with:

```bash
make build-all
```

### Running demo app images

```bash
make up-all
```

This command will start all demo apps defined as Docker Compose services, along with Prometheus.

### Accessing demo apps

```bash
    curl http://localhost:8081/demoapp1
    curl http://localhost:8082/demoapp2
    curl http://localhost:8083/demoapp3
```

Each demo app is accessible through its own TCP port, as defined in the respective Docker Compose services.


