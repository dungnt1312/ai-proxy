Generate a conventional commit message from the git diff.

## Git Diff
{{.diff}}

{{if .type}}
## Commit Type
{{.type}}
{{end}}

## Conventional Commit Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting (no code change)
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

## Guidelines

1. **Subject**: Max 50 chars, imperative mood, no period
2. **Body**: Explain what and why (not how), wrap at 72 chars
3. **Scope**: Optional, indicates section of codebase
4. **Breaking Changes**: Add `BREAKING CHANGE:` in footer

## Output

Provide ONLY the commit message, no explanation:

```
type(scope): subject

Body explaining the changes.
```
