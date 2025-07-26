# Codeforces Solutions

This repository collects solutions for Codeforces problems. We expect to have all solutions and verifiers within a week. The directory tree mirrors the contest numbers (`0-999`, `1000-1999`, `2000-2999`, ...). Inside each contest folder you will typically find:

- `problemX.txt` – the problem statement or notes for problem X.
- Solution files in Go or other languages (`1234A.go`, `solB.cpp`, etc.).
- Optional verifier programs such as `verifierA.go` that compile a submission and run it against sample tests.

## Utilities

Several helper tools live in the repository root:

- `create.go` – scaffold a new contest/problem directory.
- `auto.go` – convert existing C++ solutions to Go using the Codex CLI.
- `auto_so.go` – generate Go solutions from problem statements with Codex.
- `webserver.go` – simple HTTP server for browsing problems and submitting solutions locally.
- Other helper utilities may be added for brute-force testing or instrumentation.
- `eval.go` – evaluate AI-generated solutions. Supports flags like `-model`,
  `-provider`, `-db` and `-timeout` to control the HTTP request timeout.

All Go utilities build with the standard tooling. Example:

```bash
cd 1000-1999/1900-1999/1990-1999/1994
go run verifierA.go ./1994A
```

This repository aims to provide a complete archive. We expect to have all solutions and verifiers within a week. Feel free to use the code or utilities as a reference for your own workflow.

Some contests include verifier programs with deterministic test cases. For
example, contest 90 provides verifiers for problems A and B:

```bash
cd 0-999/0-99/90-99/90
go run verifierA.go ./mySolutionA
go run verifierB.go ./mySolutionB
```

## Running the local webserver

The `webserver.go` program lets you browse contests and test solutions
directly from your browser. Launch it from the repository root:

```bash
go run webserver.go
```

Open <http://localhost:8081> to see the list of contests. Each problem page
allows you to paste code or upload a file in C, C++, Go, Rust, Java or Python.
If a verifier is present in the contest directory it will run automatically
after compilation.

### Required compilers

The web server relies on external compilers/interpreters. Install the
following tools so every language option works:

- `gcc` and `g++`
- `javac`/`java`
- `go`
- `rustc`
- `python3`

Below are minimal installation commands for common platforms.

**Linux (Debian/Ubuntu)**

```bash
sudo apt update
sudo apt install build-essential openjdk-17-jdk golang rustc python3
```

**macOS (Homebrew)**

```bash
brew install gcc openjdk go rust python
```

**Windows**

- Install [MSYS2](https://www.msys2.org/) or enable the
  [Windows Subsystem for Linux](https://learn.microsoft.com/windows/wsl/).
- With MSYS2 you can run:

```bash
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-gcc-libgfortran \
  mingw-w64-x86_64-go mingw-w64-x86_64-rust mingw-w64-x86_64-python
```

Alternatively use the official installers for
[MinGW-w64](https://www.mingw-w64.org/),
[Adoptium JDK](https://adoptium.net/),
[Go](https://go.dev/dl/),
[Rust](https://rustup.rs/) and
[Python](https://www.python.org/downloads/).
