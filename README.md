# ctxmerge

`ctxmerge` is a lightweight command-line tool for managing Kubernetes configurations (`kubeconfig`) efficiently. It allows you to rename contexts, merge kubeconfigs, and update cluster IPs seamlessly. Built with Go, `ctxmerge` provides a secure and minimal Dockerized deployment using Distroless.

---

## Features

- **Rename Kubernetes Contexts**: Easily rename the current context in your kubeconfig.
- **Merge Kubeconfigs**: Combine multiple kubeconfigs into a single configuration.
- **Update Cluster Server IPs**: Modify the server IP address for existing clusters during a merge.
- **Multi-Platform Support**: Build binaries for Linux, macOS, and Windows.
- **Dockerized Deployment**: Deploy securely using a lightweight Distroless container.

---

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/your-github-username/ctxmerge.git
   cd ctxmerge
   ```
2. Build the application:
   ```bash
   make build
   ```
3. Run the executable from the `build/` directory.

### Using Docker

Pull the latest Docker image:
```bash
docker pull your-docker-repo/ctxmerge:latest
```

---

## Usage

### Rename Context
Rename the default context in a kubeconfig:
```bash
cat ~/.kube/config | ./ctxmerge new-context-name
```

### Merge Kubeconfigs
Merge an additional kubeconfig into your default kubeconfig:
```bash
cat additional-kubeconfig | ./ctxmerge new-context-name
```

### Update Cluster Server IP
Update a cluster's server IP while merging kubeconfigs:
```bash
cat additional-kubeconfig | ./ctxmerge new-context-name https://new-server-ip
```

### Docker Usage
Run `ctxmerge` in a Docker container:
```bash
cat ~/.kube/config | docker run --rm -i your-docker-repo/ctxmerge:latest new-context-name
```

---

## Building from Source

### Prerequisites
- [Go](https://golang.org/dl/) 1.20 or newer
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/)

### Build Commands

1. **Run Tests**:
   ```bash
   make test
   ```

2. **Build Binaries**:
   ```bash
   make build
   ```

3. **Build Docker Image**:
   ```bash
   make docker-build
   ```

4. **Release (GitHub)**:
   ```bash
   make release
   ```

---

## Development

### Folder Structure

- `main.go`: Entry point of the application.
- `main_test.go`: Unit tests for the application.
- `Dockerfile`: Multi-stage Dockerfile for building and deploying `ctxmerge`.
- `Makefile`: Automates building, testing, and releasing.

### Running Locally
Run the tool locally:
```bash
go run main.go new-context-name
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new feature branch.
3. Commit your changes and open a pull request.

---

## Acknowledgments

`ctxmerge` is inspired by the need for efficient and secure Kubernetes configuration management. Special thanks to the Kubernetes and Go communities for providing excellent tools and documentation.
