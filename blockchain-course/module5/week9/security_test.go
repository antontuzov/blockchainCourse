package week9

import (
	"strings"
	"testing"
)

func TestNewSecurityScanner(t *testing.T) {
	scanner := NewSecurityScanner()

	if scanner == nil {
		t.Error("Failed to create security scanner")
	}

	if len(scanner.Rules) == 0 {
		t.Error("Security scanner should have rules")
	}

	// Check that all expected rules are present
	expectedRules := []string{"REENTRANCY", "OVERFLOW", "GAS_LIMIT", "RANDOMNESS", "TX_ORIGIN"}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range scanner.Rules {
			if rule.ID == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule %s not found", expectedRule)
		}
	}
}

func TestScanContract(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create a contract with potential vulnerabilities
	contract := &SmartContract{
		Code: `pragma solidity ^0.5.0;

contract VulnerableContract {
    mapping(address => uint) public balances;
    
    function withdraw(uint amount) public {
        require(balances[msg.sender] >= amount);
        msg.sender.transfer(amount);
        balances[msg.sender] -= amount;
    }
    
    function play() public {
        require(msg.sender == tx.origin);
        // Some game logic
    }
    
    function random() public view returns (uint) {
        return block.timestamp;
    }
}`,
		FileName: "vulnerable.sol",
	}

	// Scan the contract
	vulnerabilities := scanner.ScanContract(contract)

	// Check that vulnerabilities are found
	if len(vulnerabilities) == 0 {
		t.Error("Vulnerabilities should be found")
	}

	// Check for specific vulnerabilities
	reentrancyFound := false
	txOriginFound := false
	randomnessFound := false

	for _, vuln := range vulnerabilities {
		switch vuln.Rule.ID {
		case "REENTRANCY":
			reentrancyFound = true
		case "TX_ORIGIN":
			txOriginFound = true
		case "RANDOMNESS":
			randomnessFound = true
		}
	}

	if !reentrancyFound {
		t.Error("Reentrancy vulnerability should be found")
	}

	if !txOriginFound {
		t.Error("TX Origin vulnerability should be found")
	}

	if !randomnessFound {
		t.Error("Randomness vulnerability should be found")
	}
}

func TestGenerateReport(t *testing.T) {
	scanner := NewSecurityScanner()

	// Create some vulnerabilities
	vulnerabilities := []Vulnerability{
		{
			Rule: VulnerabilityRule{
				ID:       "REENTRANCY",
				Name:     "Reentrancy Vulnerability",
				Severity: "HIGH",
			},
			File: "test.sol",
			Line: 10,
			Code: "msg.sender.transfer(amount);",
		},
		{
			Rule: VulnerabilityRule{
				ID:       "OVERFLOW",
				Name:     "Integer Overflow",
				Severity: "MEDIUM",
			},
			File: "test.sol",
			Line: 15,
			Code: "balance += amount;",
		},
	}

	// Generate report
	report := scanner.GenerateReport(vulnerabilities)

	if report == "" {
		t.Error("Report should not be empty")
	}

	// Check that report contains expected sections
	if !strings.Contains(report, "HIGH SEVERITY VULNERABILITIES") {
		t.Error("Report should contain high severity section")
	}

	if !strings.Contains(report, "MEDIUM SEVERITY VULNERABILITIES") {
		t.Error("Report should contain medium severity section")
	}

	// Test with no vulnerabilities
	emptyReport := scanner.GenerateReport([]Vulnerability{})
	if emptyReport != "No vulnerabilities found." {
		t.Error("Empty report should indicate no vulnerabilities found")
	}
}

func TestNewFuzzTester(t *testing.T) {
	fuzzTester := NewFuzzTester()

	if fuzzTester == nil {
		t.Error("Failed to create fuzz tester")
	}

	if len(fuzzTester.Contracts) != 0 {
		t.Error("Fuzz tester should have no contracts initially")
	}

	if len(fuzzTester.TestCases) != 0 {
		t.Error("Fuzz tester should have no test cases initially")
	}
}

func TestFuzzTesterAddContract(t *testing.T) {
	fuzzTester := NewFuzzTester()

	// Create a contract
	contract := &SmartContract{
		Code:     "pragma solidity ^0.5.0; contract Test {}",
		FileName: "test.sol",
	}

	// Add contract to fuzz tester
	fuzzTester.AddContract(contract)

	if len(fuzzTester.Contracts) != 1 {
		t.Error("Contract should be added to fuzz tester")
	}

	if fuzzTester.Contracts[0] != contract {
		t.Error("Added contract should match")
	}
}

func TestFuzzTesterAddTestCase(t *testing.T) {
	fuzzTester := NewFuzzTester()

	// Create a test case
	testCase := FuzzTestCase{
		Name:        "Test overflow",
		Input:       []interface{}{100, 200},
		ExpectedErr: "overflow",
	}

	// Add test case to fuzz tester
	fuzzTester.AddTestCase(testCase)

	if len(fuzzTester.TestCases) != 1 {
		t.Error("Test case should be added to fuzz tester")
	}

	if fuzzTester.TestCases[0].Name != "Test overflow" {
		t.Error("Test case name is incorrect")
	}
}

func TestRunFuzzTests(t *testing.T) {
	fuzzTester := NewFuzzTester()

	// Add a test case
	testCase := FuzzTestCase{
		Name:        "Test overflow",
		Input:       []interface{}{100, 200},
		ExpectedErr: "overflow",
	}
	fuzzTester.AddTestCase(testCase)

	// Run fuzz tests
	results := fuzzTester.RunFuzzTests()

	if len(results) != 1 {
		t.Error("Should have one test result")
	}

	// Check result
	result := results[0]
	if result.Name != "Test overflow" {
		t.Error("Test result name is incorrect")
	}

	// Test case with "overflow" in name should fail
	if result.Passed {
		t.Error("Test with 'overflow' in name should fail")
	}

	if result.Error != "Integer overflow detected" {
		t.Error("Error message is incorrect")
	}
}

func TestNewFormalVerifier(t *testing.T) {
	verifier := NewFormalVerifier()

	if verifier == nil {
		t.Error("Failed to create formal verifier")
	}

	if len(verifier.Contracts) != 0 {
		t.Error("Formal verifier should have no contracts initially")
	}

	if len(verifier.Specifications) != 0 {
		t.Error("Formal verifier should have no specifications initially")
	}
}

func TestFormalVerifierAddContract(t *testing.T) {
	verifier := NewFormalVerifier()

	// Create a contract
	contract := &SmartContract{
		Code:     "pragma solidity ^0.5.0; contract Test {}",
		FileName: "test.sol",
	}

	// Add contract to verifier
	verifier.AddContract(contract)

	if len(verifier.Contracts) != 1 {
		t.Error("Contract should be added to verifier")
	}

	if verifier.Contracts[0] != contract {
		t.Error("Added contract should match")
	}
}

func TestFormalVerifierAddSpecification(t *testing.T) {
	verifier := NewFormalVerifier()

	// Create a specification
	spec := Specification{
		Name:        "Balance preservation",
		Property:    "sum(balances) remains constant",
		Description: "Total balance should remain constant",
	}

	// Add specification to verifier
	verifier.AddSpecification(spec)

	if len(verifier.Specifications) != 1 {
		t.Error("Specification should be added to verifier")
	}

	if verifier.Specifications[0].Name != "Balance preservation" {
		t.Error("Specification name is incorrect")
	}
}

func TestVerify(t *testing.T) {
	verifier := NewFormalVerifier()

	// Add a specification
	spec := Specification{
		Name:        "Balance preservation",
		Property:    "sum(balances) remains constant",
		Description: "Total balance should remain constant",
	}
	verifier.AddSpecification(spec)

	// Run verification
	results := verifier.Verify()

	if len(results) != 1 {
		t.Error("Should have one verification result")
	}

	// Check result
	result := results[0]
	if result.Specification.Name != "Balance preservation" {
		t.Error("Specification name is incorrect")
	}

	if !result.Verified {
		t.Error("Specification should be verified")
	}

	if result.Proof != "Proof of correctness" {
		t.Error("Proof is incorrect")
	}
}
