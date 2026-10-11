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
- provide independent cleanup-target validation for future deletion workflows;
- reject linked cleanup targets that do not appear as ordinary directories;
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

Before a discovered directory can become eligible for future deletion, it must
pass validation independently from discovery.

## Cleanup-target validation

Deletion is not implemented yet.

The current validation requires that:

- the target exists;
- the target's final path component is exactly `node_modules`;
- the target is lexically contained within the selected scan root;
- the target is not the selected scan root;
- the final target entry is an ordinary directory rather than a symbolic link
  or supported Windows name-surrogate directory reparse point;
- the resolved target is a directory;
- the resolved target remains contained within the resolved scan root after
  symbolic-link resolution.

Cleanup targets that are themselves symbolic links or Windows name-surrogate
directory reparse points, such as junctions or mounted folders, are not
eligible for cleanup.

The `Lstat` check applies to the final cleanup-target entry. Resolved-path
containment remains necessary because an ancestor directory can still redirect
the complete path through a symbolic link, junction, or other supported
name-surrogate path component.

These checks establish a baseline safety boundary between discovery and future
destructive operations.

They do not constitute final deletion authorization.

Broader Windows reparse-point policy and protection against filesystem changes
between validation and deletion require additional work before deletion is
introduced.

## Testing

The scanner, directory-size behaviour, and cleanup-target validation are tested
using temporary filesystem trees created specifically for each test run.

This allows discovery, size calculation, and path-validation behaviour to be
verified without scanning or modifying real development repositories.

The validation tests cover:

- valid `node_modules` targets beneath the selected root;
- directories with the wrong final name;
- missing targets;
- files named `node_modules`;
- targets outside the selected root;
- path-prefix confusion such as `projects` versus `projects-other`;
- rejection of the selected scan root itself;
- rejection of symbolic-link cleanup targets where the test environment permits
  symbolic-link creation.

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
- does not yet provide dedicated Windows-specific automated coverage for
  junctions and broader reparse-point behaviour;
- does not yet provide a complete Windows reparse-point deletion policy;
- does not yet protect a future deletion operation from filesystem changes
  between validation and deletion.

Filesystem deletion will be introduced only after the relevant validation,
Windows-specific path-safety behaviour, and destructive-operation safeguards
have been implemented and tested.
