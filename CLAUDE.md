# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview
Tig is an ncurses-based text-mode interface for Git, functioning as a repository browser, staging tool, and pager for Git commands. The codebase is written in C with a modular architecture centered around views and Git integration.

## Build Commands
- **Build**: `make` or `make all` - builds the main executable `src/tig`
- **Debug build**: `make all-debug` - builds with debug symbols and optimizations disabled
- **Clean**: `make clean` - removes build artifacts
- **Dist clean**: `make distclean` - removes build artifacts and configuration files
- **Configure**: `./configure` (if building from Git, run `make configure` first)

## Testing
- **Run all tests**: `make test`
- **Run specific test**: `make test/<path-to-test>` (e.g., `make test/main/default-test`)
- **Run TODO tests**: `make test-todo`
- **Test with coverage**: `make test-coverage`
- **Test with address sanitizer**: `make test-address-sanitizer`
- **Test options**: Set `TEST_OPTS` environment variable with options like:
  - `TEST_OPTS=verbose` - show all test results
  - `TEST_OPTS=debugger=lldb` - run with debugger
  - `TEST_OPTS=filter=*:*default` - filter tests
  - `TEST_OPTS=valgrind` - run with Valgrind

## Installation
- **Install**: `make install` (installs to `$HOME/bin` by default)
- **Install to prefix**: `make prefix=/usr/local install`
- **Install docs**: `make install-doc`
- **Uninstall**: `make uninstall`

## Architecture
The codebase follows a view-based architecture with these key components:

### Core Structure
- **Views**: Main, diff, log, reflog, tree, blob, blame, refs, status, stage, stash, grep, pager, help
- **Modules**: Each view has corresponding `.c` and `.h` files in `src/` and `include/tig/`
- **Utilities**: Common functionality in string.c, util.c, io.c, options.c, etc.

### Key Files
- `src/tig.c`: Main entry point and initialization
- `src/main.c`: Main view implementation
- `src/display.c` + `src/draw.c`: UI rendering and display management
- `src/view.c`: View management and navigation
- `src/options.c`: Configuration and option parsing
- `src/keys.c`: Key binding and input handling
- `src/io.c`: Git command execution and I/O handling

### Graph Rendering
- `src/graph.c`: Graph visualization for commit history
- `src/graph-v1.c` + `src/graph-v2.c`: Different graph rendering algorithms

### Git Integration
- `src/git.c`: Git command interface
- `src/repo.c`: Repository management
- `src/refdb.c`: Reference database operations
- `src/refs.c`: Reference handling (branches, tags, etc.)

### Configuration
- `tigrc`: Default configuration file (compiled into binary)
- Configuration parsed by `src/options.c` and `src/parse.c`

## Development Setup
1. Dependencies: `git-core`, `ncurses` (with wide char support), `iconv`
2. Optional: `readline`, `pcre2`, `autoconf`, `asciidoc`, `xmlto`
3. Build: `make`
4. Test: `make test`

## File Organization
- `src/`: Main source code
- `include/tig/`: Header files
- `test/`: Test suite with subdirectories for each view type
- `doc/`: Documentation source files
- `tools/`: Build and utility scripts
- `compat/`: Compatibility layer for different systems
- `contrib/`: Configuration examples and contributed files