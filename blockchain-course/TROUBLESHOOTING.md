# Troubleshooting Guide

## Overview
This guide provides solutions to common issues students may encounter while working through the blockchain course.

## Common Setup Issues

### Go Installation Problems

**Problem**: Command not found when running `go version`
**Solution**: 
1. Verify Go is installed: `brew install go` (macOS) or `sudo apt-get install golang` (Linux)
2. Add Go to your PATH in ~/.bashrc or ~/.zshrc:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   export GOPATH=$HOME/go
   export PATH=$PATH:$GOPATH/bin
   ```
3. Reload your shell: `source ~/.bashrc` or `source ~/.zshrc`

### Module Import Issues

**Problem**: Cannot find package or import errors
**Solution**:
1. Initialize Go module in your project directory:
   ```bash
   go mod init blockchain-course
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```

### Docker Installation Issues

**Problem**: Permission denied when running Docker commands
**Solution**:
1. Add your user to the docker group:
   ```bash
   sudo usermod -aG docker $USER
   ```
2. Log out and log back in for changes to take effect

## Coding Issues

### Block Hash Mismatch

**Problem**: Blockchain validation fails with hash mismatch
**Solution**:
1. Check that all fields in the [Block](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module1/week1/block.go#L12-L21) structure are correctly set
2. Ensure [SetHash()](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module1/week1/block.go#L38-L48) is called when creating a new block
3. Verify that the [prepareData()](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module1/week2/blockchain.go#L110-L123) method includes all necessary fields in the correct order

### Proof of Work Not Working

**Problem**: Mining takes too long or never completes
**Solution**:
1. Check the `targetBits` constant - lower values mean higher difficulty
2. Verify the [Run()](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module1/week2/blockchain.go#L77-L107) method increments the nonce correctly
3. Ensure the [prepareData()](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module1/week2/blockchain.go#L110-L123) method includes the nonce in the data being hashed

### Transaction Validation Failures

**Problem**: Transactions are not validating correctly
**Solution**:
1. Check that transaction IDs are correctly calculated using the [Hash()](file:///Users/admin/Documents/Forbest/blockchainCourse/blockchain-course/module2/week3/transaction.go#L42-L50) method
2. Verify that inputs and outputs are properly structured
3. Ensure signatures are correctly generated and verified

### P2P Network Connection Issues

**Problem**: Nodes cannot connect to each other
**Solution**:
1. Check that ports are not blocked by firewall
2. Verify node addresses and ports are correctly configured
3. Ensure all nodes are running and listening on the correct ports

## Testing Issues

### Tests Not Running

**Problem**: `go test` command fails or finds no tests
**Solution**:
1. Ensure test files end with `_test.go`
2. Verify test functions start with `Test` (e.g., `func TestBlockchain(t *testing.T)`)
3. Check that test files are in the correct package

### Test Failures

**Problem**: Tests fail with unexpected results
**Solution**:
1. Use [t.Log()](file:///usr/local/go/src/testing/testing.go#L986-L990) or [t.Logf()](file:///usr/local/go/src/testing/testing.go#L997-L1002) to debug values during test execution
2. Check that test data matches expected formats
3. Verify that all setup/teardown steps are correctly implemented

## Performance Issues

### Slow Mining

**Problem**: Proof of Work mining is extremely slow
**Solution**:
1. Reduce the `targetBits` value to lower difficulty for testing
2. Add progress indicators to show mining is working:
   ```go
   if nonce%10000 == 0 {
       fmt.Printf("Mining... nonce: %d\n", nonce)
   }
   ```

### Memory Leaks

**Problem**: Program uses increasing amounts of memory
**Solution**:
1. Use `pprof` to profile memory usage:
   ```bash
   go tool pprof http://localhost:6060/debug/pprof/heap
   ```
2. Check for infinite loops or unbounded data structures
3. Ensure resources like file handles and network connections are properly closed

## Deployment Issues

### Docker Build Failures

**Problem**: Docker image build fails
**Solution**:
1. Check Dockerfile syntax and paths
2. Ensure all required files are included in the build context
3. Verify base image names are correct

### Kubernetes Deployment Issues

**Problem**: Pods fail to start or enter CrashLoopBackOff
**Solution**:
1. Check pod logs: `kubectl logs <pod-name>`
2. Describe the pod: `kubectl describe pod <pod-name>`
3. Verify config maps and secrets are correctly mounted
4. Check resource limits and requests

## Security Issues

### Invalid Signatures

**Problem**: Digital signatures fail verification
**Solution**:
1. Verify that the same data is used for signing and verification
2. Check that public and private keys are correctly paired
3. Ensure proper encoding/decoding of signature data

### Weak Cryptography

**Problem**: Concerns about cryptographic strength
**Solution**:
1. Use established libraries rather than implementing crypto yourself
2. Keep dependencies updated
3. Follow best practices for key management

## Debugging Tips

### Using Delve Debugger

1. Install Delve:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

2. Start debugging:
   ```bash
   dlv debug main.go
   ```

3. Set breakpoints and inspect variables:
   ```
   (dlv) break main.go:10
   (dlv) continue
   (dlv) print variableName
   ```

### Logging Best Practices

1. Use structured logging:
   ```go
   log.Printf("Processing block %d with hash %x", block.Index, block.Hash)
   ```

2. Add context to error messages:
   ```go
   return fmt.Errorf("failed to validate block %d: %w", block.Index, err)
   ```

## Common Error Messages

### "invalid memory address or nil pointer dereference"

**Cause**: Accessing a nil pointer
**Solution**: Check that pointers are properly initialized before use

### "slice bounds out of range"

**Cause**: Accessing array/slice with invalid indices
**Solution**: Add bounds checking before accessing elements

### "cannot assign requested address"

**Cause**: Network binding issues
**Solution**: Check that ports are available and not in use by other processes

## Getting Help

If you're still stuck after trying these solutions:

1. Check the course discussion forums
2. Review the relevant module materials
3. Ask for help during office hours
4. Create a minimal reproducible example for your issue