package util

import (
	"fmt"
	"regexp"
)

// Define the regex for valid bucket names
var (
	validPattern          = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`)
	ipPattern             = regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	adjacentPeriodPattern = regexp.MustCompile(`\.\.`)
	adjacentDashPattern   = regexp.MustCompile(`--`)
)

func ValidateBucketName(bn string) error {
	// Step 1: Check length
	if len(bn) < 3 || len(bn) > 63 {
		return fmt.Errorf("bucket name length must be between 3 and 63 characters")
	}

	// Step 2: Check if formatted like an IP address
	if ipPattern.MatchString(bn) {
		return fmt.Errorf("bucket name must not be formatted as an IP address")
	}

	// Step 3: Check for adjacent periods
	if adjacentPeriodPattern.MatchString(bn) {
		return fmt.Errorf("bucket name must not contain two adjacent periods")
	}

	// Step 4: Check for adjacent dashes
	if adjacentDashPattern.MatchString(bn) {
		return fmt.Errorf("bucket name must not contain two adjacent dashes")
	}

	// Step 5: Check if it starts and ends with a letter or number and has valid characters
	if !validPattern.MatchString(bn) {
		return fmt.Errorf("bucket name must only contain lowercase letters, numbers, hyphens, and periods, and must start and end with a letter or number")
	}

	return nil
}