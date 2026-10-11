# WindoNodeModulesCleaner

Safe `node_modules` discovery and cleanup for Windows.

WindoNodeModulesCleaner is a Windows utility written in Go for discovering
`node_modules` directories beneath a user-selected directory tree.

Built by Edward Siwek.

## Project status

Early development.

The application can currently:

- scan a user-supplied directory tree;
- identify top-level `node_modules` directories;
- attempt to calculate the size of each discovered directory;
- report inaccessible paths encountered during an otherwise valid scan;
- report the results without modifying the filesystem.

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
C:\Projects\app-one\node_modules    824.3 MiB
C:\Projects\app-two\node_modules      1.1 GiB
```

If no matching directories are found:

```text
Scanning C:\Projects for node_modules directories...
No node_modules directories found.
```

The scanner does not descend into a discovered `node_modules` directory while
searching for additional cleanup candidates.

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

## Directory sizes

Directory sizes are calculated during scanning from the files visible to the
application at that time.

The displayed values use binary units such as:

```text
KiB
MiB
GiB
```

If the application cannot determine the size of a discovered directory, the
directory remains in the scan results and its size is reported as unavailable.

A size-calculation failure does not cause an otherwise discovered
`node_modules` directory to be discarded from the results.

## Scan warnings

If the application cannot inspect a path inside an otherwise valid scan root,
it reports the affected path as a warning and continues scanning where
possible.

For example:

```text
Scan completed with 1 warning(s):
C:\Projects\restricted: Access is denied.
```

A warning means that the scan may be incomplete because part of the selected
directory tree could not be inspected.

Failure to inspect the scan root itself remains a fatal error because the
application cannot reliably perform the requested scan.

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

The current implementation is read-only and does not remove filesystem
content.

Before destructive behaviour is introduced, deletion targets will require
independent validation and explicit user confirmation.

## Testing

The scanner and directory-size behaviour are tested using temporary filesystem
trees created specifically for each test run.

This allows discovery and size calculation to be verified without scanning or
modifying real development repositories.

The test suite also verifies that invalid scan roots, including missing paths
and files supplied in place of directories, are rejected.

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
- does not provide a graphical interface;
- does not delete directories;
- does not support cancellation for long-running scans;
- performs directory-size calculation sequentially;
- does not yet provide dedicated Windows-specific tests for inaccessible
  directory ACLs, junctions, or other reparse-point behaviour.

Filesystem deletion will be introduced only after the relevant validation and
safety behaviour has been implemented and tested.
