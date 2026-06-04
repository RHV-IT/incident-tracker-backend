# Contributing to Issue Tracker

Thank you for considering contributing to the Issue Tracker project! We welcome contributions from the community.

## How to Contribute

There are many ways to contribute to this project:

1. **Report Bugs**: If you find a bug, please open an issue with detailed steps to reproduce
2. **Suggest Features**: Have an idea for improvement? Open an issue to discuss it
3. **Improve Documentation**: Help us make our documentation better and more comprehensive
4. **Fix Bugs**: Look through the open issues and submit pull requests for bug fixes
5. **Add Features**: Implement new features from the issue tracker or your own ideas

## Getting Started

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/issueTracker.git
   cd issueTracker
   ```

3. Set up the development environment:
   ```bash
   # Install dependencies (if any beyond standard library)
   go mod download
   
   # Set up the database
   docker-compose up -d
   ./createtables.sh
   
   # Run the application
   go run ./cmd/main.go
   # Or for live reload:
   air
   ```

### Making Changes

1. Create a new branch for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bug-fix
   ```

2. Make your changes, following the coding standards below
3. Test your changes thoroughly
4. Commit your changes:
   ```bash
   git add .
   git commit -m "Descriptive commit message"
   ```
5. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
6. Open a Pull Request against the main repository's `main` branch

## Coding Standards

### Go Code Style

- Follow [Go's formatting standards](https://golang.org/doc/effective_go.html#formatting)
- Use `go fmt` before committing
- Use `go vet` to check for common mistakes
- Write clear, concise comments explaining why, not what
- Keep functions focused and reasonably sized
- Handle errors explicitly, don't ignore them

### Commit Messages

- Use clear, descriptive messages
- Format: `<type>: <description>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Examples:
  - `feat: add user disable/enable endpoints`
  - `fix: correct password validation in login`
  - `docs: update API documentation in README`
  - `refactor: simplify middleware logic`
  - `test: add unit tests for user model`

### Pull Request Process

1. Ensure your code passes `go fmt` and `go vet`
2. Make sure your changes don't break existing functionality
3. Update documentation if needed (README, comments, etc.)
4. Keep PRs focused on a single change or feature
5. Write a clear PR description explaining what and why
6. Reference any related issues in the PR description
7. Be responsive to feedback and questions from maintainers

## Development Workflow

### Database Changes

If you need to modify the database schema:
1. Update `tables.sql` with the new schema
2. Update the corresponding model in `/internal/db/`
3. Update any handlers that use the changed schema
4. Test with a fresh database using `./createtables.sh`
5. Note in your PR that database migrations are needed

### Adding New Endpoints

When adding new API endpoints:
1. Define the handler function in the appropriate file under `/cmd/`
2. Add the route to `routes.js` in the appropriate group
3. Apply middleware as needed (`authMiddleware()` for protected routes)
4. Define request/response types in `types.go` if needed
5. Update the README with documentation for the new endpoint
6. Add appropriate validation and error handling

### Security Considerations

When making changes that affect security:
- Always validate and sanitize input
- Use parameterized queries to prevent SQL injection
- Never log sensitive information (passwords, tokens, etc.)
- Follow the principle of least privilege
- Consider authentication and authorization requirements
- Passwords must always be hashed with bcrypt before storage

## Testing

While this project doesn't have a comprehensive test suite yet, please consider:

1. **Manual Testing**: Test your changes thoroughly with various inputs
2. **Edge Cases**: Think about and test edge cases
3. **Error Conditions**: Ensure error handling works correctly
4. **Security Testing**: Verify that security controls work as expected
5. **Performance**: Consider performance implications of your changes

### Running Tests

If tests exist:
```bash
go test ./...
```

For specific packages:
```bash
go test ./internal/db/
```

## Reporting Issues

When reporting issues, please include:

1. **Clear Description**: What is the problem?
2. **Steps to Reproduce**: How can we reproduce the issue?
3. **Expected Behavior**: What should happen?
4. **Actual Behavior**: What actually happens?
5. **Environment**: Go version, Docker version, OS, etc.
6. **Logs**: Any relevant error logs or output
7. **Screenshots**: If applicable, include screenshots

## Getting Help

If you need help with your contribution:
1. Check the existing documentation (README.md, ARCHITECTURE.md)
2. Look at existing code for patterns and conventions
3. Ask questions in the issue tracker
4. Look for similar implementations in the codebase

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for more details.

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

Thank you again for your contributions! Together we can make this project better for everyone.