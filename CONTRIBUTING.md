# Contributing to cobra2snooty

Thanks for your interest in contributing to `cobra2snooty`,
The following is a set of guidelines for contributing to this project.
These are just guidelines, not rules. Use your best judgment, and
feel free to propose changes to this document in a pull request.

Note that `cobra2snooty` is an evolving project, so expect things to change over
time as the team learns, listens and refines how we work with the community.

## What should I know before I get started?

### Code of Conduct

This project has adopted the [MongoDB Code of Conduct](https://www.mongodb.com/community-code-of-conduct).
By participating, you are expected to uphold this code.
If you see any violations of the above or have any other concerns or questions please contact us
using the following email alias: [community-conduct@mongodb.com](mailto:community-conduct@mongodb.com).

## How Can I Contribute?

### Development setup

#### Prerequisite Tools
- [Git](https://git-scm.com/)
- [Go (at least Go 1.18)](https://golang.org/dl/)

#### Environment
- Fork the repository.
- Clone your forked repository locally.
- We use Go Modules to manage dependencies, so you can develop outside your `$GOPATH`.
- We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint our code, you can install it locally via `make setup`.
- For pull requests to be accepted, contributors must sign [MongoDB CLA](https://www.mongodb.com/legal/contributor-agreement).

### Building and testing

The following is a short list of commands that can be run in the root of the project directory:

- Run `make test` to run all unit tests.
- Run `make lint` to validate against our linting rules.

We provide a git pre-commit hook to format and check the code, to install it run `make link-git-hooks`

## Maintainer's Guide

Reviewers, please ensure that the CLA has been signed by referring to [the contributors tool](https://contributors.corp.mongodb.com/) (internal link).
