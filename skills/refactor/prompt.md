Analyze the code and suggest refactoring improvements.

## Code
{{.code}}

{{if .language}}
## Language
{{.language}}
{{end}}

{{if .focus}}
## Focus Area
{{.focus}}
{{end}}

## Analysis Areas

1. **Code Smells**: Long methods, large classes, duplicate code, dead code
2. **SOLID Principles**: Single responsibility, open/closed, Liskov, interface segregation, dependency inversion
3. **DRY**: Don't repeat yourself
4. **Readability**: Naming, structure, comments
5. **Performance**: Unnecessary operations, N+1 queries, memory leaks

## Output Format

# Refactoring Analysis

## Summary
Brief overview of code quality

## Issues Found

### 1. [Issue Name]
- **Type**: Code smell / SOLID violation / Performance
- **Location**: Where in the code
- **Problem**: What's wrong
- **Solution**: How to fix

### 2. [Issue Name]
...

## Refactored Code
```
// Improved version with changes applied
```

## Recommendations
- Priority improvements
- Long-term suggestions
