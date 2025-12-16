Generate comprehensive tests for the provided code.

## Code to Test
{{.code}}

{{if .framework}}
## Test Framework
{{.framework}}
{{end}}

## Coverage Type
{{.coverage}}

## Test Categories

1. **Happy Path**: Normal expected behavior
2. **Edge Cases**: Boundary conditions, empty inputs, max values
3. **Error Cases**: Invalid inputs, exceptions, error handling
4. **Integration**: Component interactions (if applicable)

## Guidelines

- **Arrange-Act-Assert**: Clear test structure
- **Descriptive Names**: Test names describe behavior
- **One Assertion**: Focus each test on single behavior
- **Mocking**: Mock external dependencies
- **Coverage**: Aim for meaningful coverage, not 100%

## Output Format

```language
// Test file with comprehensive test cases

describe('ComponentName', () => {
  describe('methodName', () => {
    it('should do expected behavior when given valid input', () => {
      // Arrange
      // Act
      // Assert
    });

    it('should handle edge case', () => {
      // ...
    });

    it('should throw error when invalid input', () => {
      // ...
    });
  });
});
```

Generate complete, runnable test code.
