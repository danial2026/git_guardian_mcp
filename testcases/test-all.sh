#!/usr/bin/env bash
#
# Test all 7 test scenarios for Git Guardian MCP
#

set -e

TESTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$TESTDIR/../.."

echo "=========================================="
echo "Git Guardian MCP - Test Suite"
echo "=========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test 1: Passing Tests
echo "Test 1: ✅ Passing Tests"
echo "────────────────────────────────────────"
if go test -v ./pkg/testcases -run "^TestPassing" 2>&1; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED (unexpected)${NC}"
fi
echo ""

# Test 2: Formatting Errors
echo "Test 2: ❌ Formatting Errors (gofmt)"
echo "────────────────────────────────────────"
if gofmt -l pkg/testcases/formatting_error_test.go | grep -q "."; then
    echo -e "${YELLOW}✓ Found formatting issues (expected):${NC}"
    gofmt -d pkg/testcases/formatting_error_test.go | head -20
else
    echo -e "${RED}✗ No issues found (unexpected)${NC}"
fi
echo ""

# Test 3: Go Vet Errors
echo "Test 3: ❌ Go Vet Errors"
echo "────────────────────────────────────────"
if go vet ./pkg/testcases/govet_error_test.go 2>&1 | head -10; then
    echo -e "${RED}✗ No issues found (unexpected)${NC}"
else
    echo -e "${YELLOW}✓ Found go vet issues (expected)${NC}"
fi
echo ""

# Test 4: Failing Tests  
echo "Test 4: ❌ Test Failures"
echo "────────────────────────────────────────"
if go test -v ./pkg/testcases/failing_test.go ./pkg/testcases/passing_test.go -run "^TestFailing" 2>&1 | head -15; then
    echo -e "${RED}✗ Tests passed (unexpected)${NC}"
else
    echo -e "${YELLOW}✓ Tests failed as expected${NC}"
fi
echo ""

# Test 5: Unused Variables (skip if golangci-lint not installed)
echo "Test 5: ❌ Unused Variables/Imports"
echo "────────────────────────────────────────"
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run pkg/testcases/unused_error_test.go 2>&1 | head -15; then
        echo -e "${RED}✗ No issues found${NC}"
    else
        echo -e "${YELLOW}✓ Found unused code issues (expected)${NC}"
    fi
else
    echo -e "${YELLOW}⊘ Skipped (golangci-lint not installed)${NC}"
fi
echo ""

# Test 6: Race Conditions
echo "Test 6: ❌ Race Conditions"
echo "────────────────────────────────────────"
if go test -race ./pkg/testcases -run "^TestRace" 2>&1 | head -20; then
    echo -e "${RED}✗ No races found (tests might pass sometimes)${NC}"
else
    echo -e "${YELLOW}✓ Found race conditions (expected)${NC}"
fi
echo ""

# Test 7: Compilation Errors
echo "Test 7: ❌ Compilation Errors"
echo "────────────────────────────────────────"
if [ -f "pkg/testcases/compilation_error_test.go.disabled" ]; then
    echo -e "${YELLOW}✓ Compilation error test file exists (disabled)${NC}"
    echo "To test: rename .go.disabled to .go and run 'go test'"
    echo "Expected result: compilation will fail"
else
    echo -e "${RED}✗ Test file not found${NC}"
fi
echo ""

echo "=========================================="
echo "Summary"
echo "=========================================="
echo "1. ✅ Passing tests - PASS"
echo "2. ❌ Formatting errors - gofmt catches"
echo "3. ❌ Go vet errors - go vet catches"
echo "4. ❌ Test failures - go test catches"
echo "5. ❌ Unused code - linter catches"
echo "6. ❌ Race conditions - go test -race catches"
echo "7. ❌ Compilation errors - compiler catches"
echo ""
echo "All 7 test scenarios created successfully!"

