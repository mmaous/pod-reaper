# Contributing to pod-reaper

First off, thank you for considering contributing to `pod-reaper`! It's people like you that make open source tools helpful for everyone.

## Getting Started

1. **Fork the repository** and clone it locally.
2. **Install Go**: Ensure you have Go 1.20+ installed on your system.
3. **Run the tests**: Verify your environment is set up properly by running the test suite:
   ```bash
   go test ./...
   ```

## Development Workflow

1. Create a branch for your feature or bug fix:
   ```bash
   git checkout -b feature/my-awesome-feature
   ```
2. Write code!
3. Add or update tests to cover your changes, particularly if you are adding new pod states to monitor in `pkg/reaper/detector_test.go`.
4. Run the tests and ensure everything passes (`go test ./...`).
5. Ensure your code is properly formatted (`go fmt ./...`).
6. Commit your changes with a descriptive commit message.

## Submitting a Pull Request

1. Push your branch to your fork on GitHub.
2. Open a Pull Request against the `main` branch.
3. Fill out the Pull Request template, providing as much context as possible.
4. Once submitted, the maintainers will review your PR and provide feedback or merge it in!

## Reporting Bugs and Requesting Features

If you encounter a bug or have a feature idea but aren't ready to write the code yourself, please open an issue! Use the provided GitHub Issue templates to ensure you include all the necessary information.
