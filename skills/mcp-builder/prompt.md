Create an MCP (Model Context Protocol) server for the specified service.

## Service to Integrate
{{.service}}

## Language
{{.language}}

## Transport
{{.transport}}

## Process

### Phase 1: Research and Planning

1. **Understand the API**: Review service's API documentation, authentication, data models
2. **Tool Selection**: Prioritize comprehensive API coverage. List endpoints to implement.
3. **Tool Naming**: Use consistent prefixes (e.g., `github_create_issue`, `github_list_repos`)

### Phase 2: Implementation

**Project Structure (TypeScript)**:
```
mcp-server-name/
├── src/
│   ├── index.ts      # Entry point
│   ├── tools/        # Tool implementations
│   └── types.ts      # Type definitions
├── package.json
└── tsconfig.json
```

**Project Structure (Python)**:
```
mcp_server_name/
├── src/
│   ├── __init__.py
│   ├── server.py     # Main server
│   └── tools.py      # Tool implementations
├── pyproject.toml
└── README.md
```

### Key Principles

1. **Clear Tool Names**: Descriptive, action-oriented naming
2. **Concise Descriptions**: Help agents find right tools quickly
3. **Actionable Errors**: Guide agents toward solutions
4. **Focused Results**: Return relevant data, support pagination

### Phase 3: Testing

1. Test each tool individually
2. Verify error handling
3. Test with actual MCP client

Create a complete, working MCP server implementation.
