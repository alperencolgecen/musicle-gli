# Contributing to MusicLe CLI Music Player

Thank you for your interest in contributing to MusicLe! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.26.1 or later
- Git
- A terminal that supports ANSI escape sequences

### Setup

1. Fork the repository
2. Clone your fork locally
3. Create a new branch for your feature/bug fix
4. Make your changes
5. Test your changes
6. Submit a pull request

```bash
git clone https://github.com/alperencolgecen/musicle.git
cd musicle
git checkout -b feature/your-feature-name
```

## Development

### Building the Project

```bash
go build -o musicle.exe ./cmd/musicle
```

### Running Tests

```bash
go test ./...
```

### Code Style

- Follow Go conventions and formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

## Project Structure

```
musicle/
├── cmd/musicle/          # Main application entry point
├── internal/
│   ├── bridge/          # Bridge between UI and engine
│   ├── fs/              # File system operations
│   ├── theme/           # UI theming
│   └── ui/              # User interface components
├── engine/              # Music player engine
└── docs/               # Documentation
```

## Types of Contributions

We welcome the following types of contributions:

### Bug Reports

- Use the issue template for bug reports
- Provide detailed steps to reproduce
- Include system information (OS, Go version, etc.)
- Add screenshots if applicable

### Feature Requests

- Use the feature request template
- Describe the use case clearly
- Consider if it fits the project's scope

### Code Contributions

- Bug fixes
- New features
- Performance improvements
- Documentation improvements
- UI/UX enhancements

### Documentation

- Improving README.md
- Adding inline code comments
- Creating tutorials or guides
- Translating documentation

## Pull Request Process

1. **Fork and Branch**: Create a feature branch from `main`
2. **Make Changes**: Implement your changes with clear commit messages
3. **Test**: Ensure your changes work and don't break existing functionality
4. **Documentation**: Update relevant documentation
5. **Submit PR**: Create a pull request with a clear title and description

### Pull Request Template

When submitting a PR, please include:

- **Description**: What changes were made and why
- **Testing**: How you tested the changes
- **Screenshots**: For UI changes (if applicable)
- **Breaking Changes**: Any breaking changes and migration steps

## Code Review Process

- All submissions require review
- Maintainers may request changes
- Be responsive to feedback
- Keep discussions constructive and professional

## Release Process

- Maintainers handle releases
- Follow semantic versioning
- Update CHANGELOG.md
- Tag releases properly

## Community Guidelines

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Help

- Check existing issues and documentation
- Ask questions in discussions
- Join our community channels (if available)

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.

Thank you for contributing to MusicLe! 🎵
