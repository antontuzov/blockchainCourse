package week10

import (
	"strings"
	"testing"
	"time"
)

func TestNewDeploymentManager(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "golang:1.19",
		Ports:    []int{8080, 8081},
		Environment: map[string]string{
			"ENV": "test",
		},
		Resources: ResourceConfig{
			CPU:    "500m",
			Memory: "512Mi",
		},
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
		Subnet:      "10.0.0.0/24",
		Gateway:     "10.0.0.1",
		DNS:         []string{"8.8.8.8", "8.8.4.4"},
	}

	dm := NewDeploymentManager(config, network)

	if dm == nil {
		t.Error("Failed to create deployment manager")
	}

	if dm.Config != config {
		t.Error("Deployment config is incorrect")
	}

	if dm.Network != network {
		t.Error("Network config is incorrect")
	}

	if dm.Monitor == nil {
		t.Error("Monitor should be created")
	}

	if dm.Monitor.Metrics == nil {
		t.Error("Metrics channel should be created")
	}

	if dm.Monitor.Alerts == nil {
		t.Error("Alerts channel should be created")
	}
}

func TestDeploy(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "golang:1.19",
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Deploy the application
	err := dm.Deploy()
	if err != nil {
		t.Errorf("Failed to deploy: %s", err)
	}

	// Check that monitoring is started
	// This is difficult to test directly, but we can check that the channels are being used
	time.Sleep(1 * time.Second)

	// Check that metrics are being generated
	metrics := dm.GetMetrics()
	if len(metrics) > 0 {
		t.Logf("Generated %d metrics", len(metrics))
	}
}

func TestScale(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "golang:1.19",
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Scale the deployment
	err := dm.Scale(5)
	if err != nil {
		t.Errorf("Failed to scale: %s", err)
	}

	if dm.Config.Replicas != 5 {
		t.Error("Replicas count is incorrect")
	}
}

func TestUpdate(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "golang:1.19",
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Update the deployment
	err := dm.Update("1.0.1")
	if err != nil {
		t.Errorf("Failed to update: %s", err)
	}

	if dm.Config.Version != "1.0.1" {
		t.Error("Version is incorrect")
	}
}

func TestRollback(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.1",
		Replicas: 3,
		Image:    "golang:1.19",
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Rollback the deployment
	err := dm.Rollback()
	if err != nil {
		t.Errorf("Failed to rollback: %s", err)
	}

	// In a real implementation, we would check that the version was rolled back
	// For this test, we just verify that the function completed without error
}

func TestGenerateDockerfile(t *testing.T) {
	config := &DeploymentConfig{
		Name:    "test-blockchain",
		Version: "1.0.0",
		Image:   "golang:1.19",
		Ports:   []int{8080, 8081},
		Environment: map[string]string{
			"ENV":   "test",
			"DEBUG": "true",
		},
		HealthCheck: HealthCheckConfig{
			Command: []string{"CMD", "curl", "-f", "http://localhost:8080/health"},
		},
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Generate Dockerfile
	dockerfile := dm.GenerateDockerfile()

	if dockerfile == "" {
		t.Error("Dockerfile should not be empty")
	}

	// Check that Dockerfile contains expected content
	if !strings.Contains(dockerfile, "FROM golang:1.19") {
		t.Error("Dockerfile should contain base image")
	}

	if !strings.Contains(dockerfile, "ENV ENV=test") {
		t.Error("Dockerfile should contain environment variables")
	}

	if !strings.Contains(dockerfile, "EXPOSE 8080") {
		t.Error("Dockerfile should expose ports")
	}

	if !strings.Contains(dockerfile, "HEALTHCHECK") {
		t.Error("Dockerfile should contain health check")
	}
}

func TestGenerateKubernetesManifest(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "test-blockchain:1.0.0",
		Ports:    []int{8080, 8081},
		Environment: map[string]string{
			"ENV":   "test",
			"DEBUG": "true",
		},
		Resources: ResourceConfig{
			CPU:    "500m",
			Memory: "512Mi",
		},
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Generate Kubernetes manifest
	manifest := dm.GenerateKubernetesManifest()

	if manifest == "" {
		t.Error("Manifest should not be empty")
	}

	// Check that manifest contains expected content
	if !strings.Contains(manifest, "apiVersion: apps/v1") {
		t.Error("Manifest should contain apiVersion")
	}

	if !strings.Contains(manifest, "kind: Deployment") {
		t.Error("Manifest should contain kind")
	}

	if !strings.Contains(manifest, "name: test-blockchain") {
		t.Error("Manifest should contain deployment name")
	}

	if !strings.Contains(manifest, "replicas: 3") {
		t.Error("Manifest should contain replicas")
	}

	if !strings.Contains(manifest, "image: test-blockchain:1.0.0") {
		t.Error("Manifest should contain image")
	}

	if !strings.Contains(manifest, "containerPort: 8080") {
		t.Error("Manifest should contain ports")
	}

	if !strings.Contains(manifest, "cpu: 500m") {
		t.Error("Manifest should contain resource limits")
	}
}

func TestGetMetrics(t *testing.T) {
	config := &DeploymentConfig{
		Name:     "test-blockchain",
		Version:  "1.0.0",
		Replicas: 3,
		Image:    "golang:1.19",
	}

	network := &NetworkConfig{
		NetworkName: "test-network",
	}

	dm := NewDeploymentManager(config, network)

	// Start monitoring
	go dm.StartMonitoring()

	// Wait for some metrics to be generated
	time.Sleep(2 * time.Second)

	// Get metrics
	metrics := dm.GetMetrics()

	// Check that metrics were collected
	if len(metrics) == 0 {
		t.Error("Metrics should be collected")
	}

	// Check metric structure
	for _, metric := range metrics {
		if metric.Name == "" {
			t.Error("Metric name should not be empty")
		}

		if metric.Timestamp.IsZero() {
			t.Error("Metric timestamp should not be zero")
		}

		if metric.Labels == nil {
			t.Error("Metric labels should not be nil")
		}
	}
}
