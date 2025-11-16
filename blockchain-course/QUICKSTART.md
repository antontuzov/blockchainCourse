# Quick Start Guide

## Overview
This guide provides a fast path to get you started with the blockchain course.

## Prerequisites

1. **Go 1.19+** - Download from [golang.org](https://golang.org/dl/)
2. **Git** - Version control system
3. **Docker** - Containerization platform
4. **VS Code** (recommended) - Code editor with Go extensions

## Installation Steps

### 1. Install Go
```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt-get update
sudo apt-get install golang

# Windows
# Download installer from https://golang.org/dl/
```

### 2. Install Docker
```bash
# macOS
brew install docker

# Ubuntu/Debian
sudo apt-get install docker docker-compose

# Windows
# Download Docker Desktop from https://www.docker.com/products/docker-desktop
```

### 3. Clone the Repository
```bash
git clone <repository-url>
cd blockchain-course
```

### 4. Initialize Go Modules
```bash
go mod tidy
```

## Running Your First Example

### 1. Navigate to Module 1, Week 2
```bash
cd module1/week2
```

### 2. Run the Basic Blockchain
```bash
go run main.go
```

### 3. Interact with the CLI
```
Blockchain CLI
1. Add block
2. Print blockchain
3. Validate blockchain
4. Exit
Enter choice: 
```

## Week-by-Week Quick Start

### Week 1: Cryptography Basics
```bash
cd module1/week1
go run main.go
```

### Week 2: Basic Blockchain
```bash
cd module1/week2
go run main.go
```

### Week 3: Transactions
```bash
cd module2/week3
go run main.go
```

### Week 4: Smart Contracts
```bash
cd module2/week4
go run main.go
```

### Week 5: P2P Network
```bash
cd module3/week5
go run main.go
```

### Week 6: Consensus
```bash
cd module3/week6
go run main.go
```

## Testing Your Code

### Run All Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Specific Module Tests
```bash
cd module1/week1
go test -v
```

## Building Docker Images

### Build Main Application
```bash
docker build -t blockchain-app .
```

### Run in Container
```bash
docker run -p 8080:8080 blockchain-app
```

## Common Commands

### Code Formatting
```bash
go fmt ./...
```

### Dependency Management
```bash
# Add a new dependency
go get github.com/some/package

# Update dependencies
go mod tidy
```

### Benchmarking
```bash
go test -bench=.
```

## Project Structure

```
blockchain-course/
├── module1/           # Fundamentals
│   ├── week1/         # Cryptography
│   └── week2/         # Basic Blockchain
├── module2/           # Advanced Features
│   ├── week3/         # Transactions
│   └── week4/         # Smart Contracts
├── module3/           # Networking & Consensus
│   ├── week5/         # P2P Network
│   └── week6/         # Consensus Algorithms
├── module4/           # Enterprise Blockchain
│   ├── week7/         # Permissioned Chains
│   └── week8/         # Scalability
├── module5/           # Production & Security
│   ├── week9/         # Security
│   └── week10/        # Deployment
└── capstone-projects/ # Final Projects
```

## Getting Help

1. **Check Documentation**: Each module has a README.md with detailed instructions
2. **Review Tutorials**: See [TUTORIALS.md](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/TUTORIALS.md) for step-by-step guides
3. **Troubleshooting**: See [TROUBLESHOOTING.md](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/TROUBLESHOOTING.md) for common issues
4. **Architecture Overview**: See [ARCHITECTURE.md](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/ARCHITECTURE.md) for system diagrams

## Next Steps

1. **Complete Week 1**: Start with basic cryptography
2. **Build the Core**: Work through modules sequentially
3. **Experiment**: Modify code to understand concepts better
4. **Test**: Run all tests to ensure your implementation works
5. **Capstone**: Choose a capstone project in Week 10

## Useful Resources

- [Go Documentation](https://golang.org/doc/)
- [Blockchain Whitepapers](https://bitcoin.org/bitcoin.pdf)
- [Course Architecture Diagrams](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/ARCHITECTURE.md)
- [Interactive Tutorials](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/TUTORIALS.md)