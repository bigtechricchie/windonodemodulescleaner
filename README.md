# WindoNodeModulesCleaner

Safe `node_modules` discovery and cleanup for Windows.

WindoNodeModulesCleaner is a Windows utility written in Go for discovering
`node_modules` directories beneath a user-selected directory tree.

Built by Edward Siwek.

## Project status

Early development.

The application can currently scan a user-supplied directory tree and report top-level `node_modules` directories without modifying the filesystem.

Deletion is not implemented yet.

## Scope

The initial product is intentionally limited to directories named:

```text
node_modules
```

It is not intended to be a generic filesystem cleaner.

Additional cleanup targets are outside the initial scope.

## Platform

Initial target:

- Windows 10;
- Windows 11;
- x86-64.

## Usage

Scan a directory tree by providing the root directory:

```powershell
go run . -root C:\Projects
```

Example output:

```text
Scanning C:\Projects for node_modules directories...
Found 2 node_modules directories:
C:\Projects\app-one\node_modules
C:\Projects\app-two\node_modules
```

If no matching directories are found:

```text
Scanning C:\Projects for node_modules directories...
No node_modules directories found.
```

The scanner does not descend into a discovered `node_modules` directory.

For example, given:

```text
C:\Projects\app\
└── node_modules\
    └── package-a\
        └── node_modules\
```

the scan reports only:

```text
C:\Projects\app\node_modules
```

The nested dependency directory is not reported separately because it is
already contained within the outer `node_modules` tree.

## Engineering principles

The project favours:

- clear and descriptive Go code;
- standard-library solutions where practical;
- explicit error handling;
- narrow destructive operations;
- least privilege;
- testable filesystem behaviour;
- minimal dependencies;
- architecture driven by demonstrated requirements.

## Safety model

Discovery and deletion are separate operations.

Finding a directory gives the application permission to present it for review,
not permission to delete it.

The current scanner is read-only and does not remove filesystem content.

Before destructive behaviour is introduced, deletion targets will require
independent validation and explicit user confirmation.

## Testing

The scanner is tested using temporary filesystem trees created specifically for
the test run.

This allows discovery behaviour to be verified without scanning or modifying
real development repositories.

Run the test suite with:

```powershell
go test ./...
```

Additional engineering checks:

```powershell
go fmt ./...
go vet ./...
go test ./...
go build ./...
```

## Current limitations

The application currently:

- requires the scan root to be supplied through the command line;
- reports discovered directory paths only;
- does not calculate directory sizes;
- does not provide a graphical interface;
- does not delete directories;
- does not yet provide detailed reporting for partially inaccessible directory trees.

Filesystem deletion will be introduced only after the relevant validation and
safety behaviour has been implemented and tested.
