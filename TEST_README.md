# Test configuration and scripts for electric-circuit-web server

## Running Tests

### All Tests
```bash
go test ./...
```

### Specific Package Tests
```bash
# Handler tests
go test ./internal/handlers/tests

# Controller tests  
go test ./internal/controllers/tests

# Service tests
go test ./internal/services/tests

# Firebase package tests
go test ./pkg/firebase/tests
```

### Test with Coverage
```bash
go test -cover ./...
```

### Detailed Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Verbose Test Output
```bash
go test -v ./...
```

## Test Structure

### Unit Tests
- **Handler Tests**: Test HTTP request/response handling
- **Controller Tests**: Test business logic coordination
- **Service Tests**: Test business logic implementation
- **Repository Tests**: Test data access layer

### Integration Tests
- **Full Stack Tests**: Test complete request flow
- **Firebase Integration**: Test with Firebase emulator

### Mock Objects
- Mock controllers for handler testing
- Mock services for controller testing
- Mock repositories for service testing

## Test Conventions

### File Naming
- Test files: `*_test.go`
- Test functions: `TestFunctionName_Scenario_ExpectedResult`
- Example: `TestAuthHandler_HandleVerifyToken_Success`

### Test Structure (AAA Pattern)
```go
func TestFunction_Scenario_ExpectedResult(t *testing.T) {
    // Arrange - Set up test data and mocks
    
    // Act - Execute the function being tested
    
    // Assert - Verify the results
}
```

### Mock Naming
- Mock structs: `Mock{InterfaceName}`
- Example: `MockAuthController`, `MockAuthService`

## Firebase Testing

For Firebase-related tests, you can use the Firebase emulator:

### Start Firebase Emulator
```bash
firebase emulators:start --only auth,firestore
```

### Run Tests with Emulator
```bash
export FIRESTORE_EMULATOR_HOST=localhost:8080
export FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
go test ./...
```

## Continuous Integration

Add to your CI pipeline:
```yaml
- name: Run Tests
  run: |
    go test -v -cover ./...
    
- name: Generate Coverage Report
  run: |
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out
```

## Performance Tests

For performance testing:
```bash
go test -bench=. ./...
```

## Test Data

Place test data files in:
- `testdata/` directories next to test files
- Use `//go:embed` for embedding test files

## Common Test Patterns

### Table-Driven Tests
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```