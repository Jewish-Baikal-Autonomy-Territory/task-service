# Contributing to [task-service]

Thank you for your interest in contributing to **task-service**! We welcome contributions from the community to help improve the project. Here are some guidelines to help you get started:

## Commit message Guidelines

We use [<ins>Conventional Commits</ins>](https://www.conventionalcommits.org/en/v1.0.0/). All PRs must follow this format:

`<type>[optional-scope]: <description>`

- **feat**: A new feature (e.g. `feat: add support for kafka consumer`)
- **fix**: A bug fix (e.g. `fix: correct memory leak in <class method>`)
- **perf**: A code change that improves performance (e.g. `perf: optimize object allocation algorithm for better performace`)
- **refactor**: A code change that neither fixes a bug nor adds a feature (e.g. `refactor: simplify <struct> interface`)
- **ci**: Changes to our CI/CD configuration or scripts (e.g. `ci: update Github Actions workflow for better testing`)
- **docs**: Documentation only changes (e.g. `docs: update README with new <struct> examples`)
- **test**: Adding missing tests or correcting existing tests (e.g. `test: add unit tests for <struct>`)
- **chore**: Other changes that don't modify src or test files (e.g. `chore: update dependencies`)

> **Note:** Commits that do not follow this format will be rejected by the CI/CD linter.