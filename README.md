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
- validate discovered cleanup candidates independently from discovery;
- report whether each discovered candidate is currently eligible for future
  cleanup;
- preserve rejected candidates in the scan results with the rejection reason;
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
- x86-64;
- Go 1.27.1.

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
  Cleanup eligibility: currently eligible

C:\Projects\app-two\node_modules      1.1 GiB
  Cleanup eligibility: currently eligible
```

A discovered directory that does not satisfy the cleanup-validation rules
remains visible:

```text
C:\Projects\workspace\app\node_modules    420.3 MiB
  Cleanup eligibility: rejected
  Reason: validate cleanup target path: path entry is not an ordinary directory: C:\Projects\workspace
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

## Discovery and cleanup eligibility

Discovery and cleanup eligibility are deliberately separate concepts.

Finding a `node_modules` directory establishes that the directory was
discovered during the scan. It does not establish permission to delete it.

Each discovered candidate is validated independently after discovery.

A candidate that fails cleanup validation remains in the results and is
reported as rejected together with the validation reason.

Therefore:

```text
discovered
does not imply
eligible for cleanup
```

A validation rejection does not become a fatal scan error and does not cause
the discovered candidate to disappear from the results.

An eligible result is reported as **currently eligible** because validation
represents the filesystem state observed at that point in time. It is not
durable authorization for a future destructive operation.

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

Size availability and cleanup eligibility are independent.

A candidate may therefore be currently eligible for future cleanup even when
its size could not be calculated.

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

Discovery, cleanup validation, and future deletion are separate operations.

The current implementation is read-only and does not remove filesystem
content.

A discovered directory must pass cleanup-target validation before it can be
reported as currently eligible for a future cleanup workflow.

Current eligibility is not deletion authorization.

## Cleanup-target validation

The current validation requires that:

- the target exists;
- the target's final path component is exactly `node_modules`;
- the target is lexically contained within the selected scan root;
- the target is not the selected scan root;
- every descendant path component between the selected scan root and the target
  is reported by `Lstat` as an ordinary directory;
- the resolved target is a directory;
- the resolved target remains contained within the resolved scan root after
  symbolic-link resolution.

The selected scan root is the user-established trust boundary.

Filesystem topology above the selected root, and the topology used to reach the
selected root itself, is not currently policed by cleanup validation.

Cleanup targets are rejected if the target itself, or any descendant path
component beneath the selected root, is a symbolic link or supported Windows
name-surrogate directory reparse point that `Lstat` does not report as an
ordinary directory.

`EvalSymlinks` provides an additional symbolic-link containment check.

It is not relied upon as complete Windows junction or mount-point resolution.

These checks establish a baseline safety boundary between discovery and future
destructive operations.

They do not constitute final deletion authorization.

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
  symbolic-link creation;
- rejection of cleanup targets beneath symbolic-link ancestors where the test
  environment permits symbolic-link creation.

The scanner remains responsible only for discovery and size information.

Cleanup eligibility is assessed separately by the application flow after
discovery.

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
- does not provide a cleanup confirmation workflow;
- does not support cancellation for long-running scans;
- performs directory-size calculation sequentially;
- does not yet provide dedicated Windows-specific automated coverage for
  junctions and broader reparse-point behaviour;
- does not provide complete Windows reparse-point handling;
- does not protect a future deletion operation from filesystem changes between
  validation and deletion.

Filesystem deletion will be introduced only after the relevant validation,
Windows-specific path-safety behaviour, pre-deletion revalidation, and
destructive-operation safeguards have been implemented and tested.
