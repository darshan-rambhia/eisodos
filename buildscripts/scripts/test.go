package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	testOutputDir    = ".out/test"
	junitFile        = "junit.xml"
	coverageFile     = "coverage.out"
	minCoverageLimit = 80.0
	testTimeout      = "30s"
)

func main() {
	// Create test output directory
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		fmt.Printf("Error creating test output directory: %v\n", err)
		os.Exit(1)
	}

	// Construct gotestsum command
	args := []string{
		"--junitfile", filepath.Join(testOutputDir, junitFile),
		"--format", "testname",
		"--",
		"-v",
		"-race",
		"-timeout", testTimeout,
		"-coverprofile=" + filepath.Join(testOutputDir, coverageFile),
		"github.com/darshan-rambhia/eisodos/cmd/eisodos",
		"github.com/darshan-rambhia/eisodos/internal/backend",
		"github.com/darshan-rambhia/eisodos/internal/serverpool",
		"github.com/darshan-rambhia/eisodos/config",
	}

	// Run gotestsum
	cmd := exec.Command("gotestsum", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}

	// Check coverage
	coverageFile := filepath.Join(testOutputDir, coverageFile)
	coverage, err := getCoverage(coverageFile)
	if err != nil {
		fmt.Printf("Error checking coverage: %v\n", err)
		os.Exit(1)
	}

	// Print package-wise coverage
	fmt.Printf("\nPackage Coverage:\n")
	fmt.Printf("================\n")
	packageCoverage, err := getPackageCoverage(coverageFile)
	if err != nil {
		fmt.Printf("Error getting package coverage: %v\n", err)
	} else {
		for pkg, cov := range packageCoverage {
			fmt.Printf("%-50s %.1f%%\n", pkg, cov)
		}
	}

	fmt.Printf("\nTotal Coverage: %.2f%%\n", coverage)
	if coverage < minCoverageLimit {
		fmt.Printf("Coverage %.2f%% is below the minimum required %.2f%%\n", coverage, minCoverageLimit)
		os.Exit(1)
	}

	fmt.Printf("\nTest reports generated in %s:\n", testOutputDir)
	fmt.Printf("- JUnit XML: %s\n", junitFile)
	fmt.Printf("- Coverage: %s\n", coverageFile)
}

func getCoverage(coverageFile string) (float64, error) {
	// Run go tool cover to get coverage statistics
	cmd := exec.Command("go", "tool", "cover", "-func", coverageFile)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get coverage: %v", err)
	}

	// Parse the output to get total coverage
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "total:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				coverage := strings.TrimSuffix(fields[len(fields)-1], "%")
				return strconv.ParseFloat(coverage, 64)
			}
		}
	}

	return 0, fmt.Errorf("could not find total coverage in output")
}

func getPackageCoverage(coverageFile string) (map[string]float64, error) {
	// Run go tool cover to get coverage statistics
	cmd := exec.Command("go", "tool", "cover", "-func", coverageFile)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get coverage: %v", err)
	}

	// Parse the output to get package coverage
	packageCoverage := make(map[string]float64)
	var currentPkg string

	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.Contains(line, "total:") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		filePath := parts[0]
		coverage := strings.TrimSuffix(parts[len(parts)-1], "%")
		cov, err := strconv.ParseFloat(coverage, 64)
		if err != nil {
			continue
		}

		// Extract package name from file path
		pkgPath := strings.TrimPrefix(filePath, "github.com/darshan-rambhia/eisodos/")
		pkg := pkgPath[:strings.LastIndex(pkgPath, "/")]

		if pkg != currentPkg {
			currentPkg = pkg
			packageCoverage[pkg] = cov
		} else {
			// Average the coverage for the package
			packageCoverage[pkg] = (packageCoverage[pkg] + cov) / 2
		}
	}

	return packageCoverage, nil
}
