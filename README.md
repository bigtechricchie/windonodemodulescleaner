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
- allow currently eligible candidates to be explicitly selected by invocation-
  local candidate index;
- allow the exact selected candidate set to be explicitly confirmed for the
  current invocation;
- freshly validate every confirmed candidate before any future destructive
  stage;
- report selection, confirmation, and fresh-validation state without modifying
  the filesystem.

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
- Go 1.27.2.

## Usage

Scan a directory tree:

```powershell
go run . -root C:\Projects
```

Example output:

```text
Scanning C:\Projects for node_modules directories...

Found 3 node_modules directories:

[1] C:\Projects\app-one\node_modules    824.3 MiB
    Cleanup eligibility: currently eligible

[2] C:\Projects\app-two\node_modules      1.1 GiB
    Cleanup eligibility: currently eligible

[3] C:\Projects\linked\node_modules     420.3 MiB
    Cleanup eligibility: rejected
    Reason: validate cleanup target path: path entry is not an ordinary directory: C:\Projects\linked
```

Currently eligible candidates can be selected using their displayed indexes:

```powershell
go run . -root C:\Projects -select 1,2
```

Example selection output:

```text
Selected 2 candidates:

[1] C:\Projects\app-one\node_modules
[2] C:\Projects\app-two\node_modules

No filesystem changes were made.
```

The exact selected set can be explicitly confirmed:

```powershell
go run . -root C:\Projects -select 1,2 -confirm
```

Example confirmation output:

```text
Selected 2 candidates:

[1] C:\Projects\app-one\node_modules
[2] C:\Projects\app-two\node_modules

Selection confirmed.

Fresh validation passed for 2 candidates.

No filesystem changes were made.
```

The `-confirm` option requires `-select`. Supplying `-confirm` without a
selection request is invalid and is rejected before filesystem scanning begins.

The `-select` option accepts a comma-separated list of positive decimal
candidate indexes.

For example:

```text
1
1,2
1,2,5
```

Candidate indexes are valid only for the current invocation. They are
presentation and input handles, not persistent candidate identities.

All discovered candidates receive an index, including candidates that are
rejected by cleanup validation.

Rejected candidates remain visible but cannot be selected or confirmed.

Malformed, duplicate, nonexistent, or rejected candidate indexes cause the
entire selection request to fail. Selection is never partially accepted.

Omitting `-select` performs a normal inspection-only scan.

Explicitly supplying an empty `-select` value is invalid.

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

## Discovery, eligibility, selection, and confirmation

Discovery, cleanup eligibility, selection, and confirmation are deliberately
separate concepts.

Finding a `node_modules` directory establishes that the directory was
discovered during the scan. It does not establish permission to delete it.

Each discovered candidate is validated independently after discovery.

A candidate that fails cleanup validation remains in the results and is
reported as rejected together with the validation reason.

An eligible result is reported as **currently eligible** because validation
represents the filesystem state observed at that point in time.

Selection is a separate explicit user action.

Confirmation is an explicit acknowledgement of the exact successfully resolved
selected candidate set for the current invocation.

The lifecycle is:

```text
discovered
≠
currently eligible
≠
selected
≠
confirmed
≠
freshly validated for deletion
≠
deleted
```

Only candidates discovered in the current invocation and assessed as currently
eligible can be selected.

Selection does not constitute confirmation, durable authorization, or
permission to delete filesystem content.

Confirmation also does not constitute durable authorization or permission to
delete filesystem content. It records explicit user acknowledgement of the
current selected set only.

After confirmation, every selected candidate is validated again against the
current filesystem state. If any candidate fails this fresh validation, the
entire gate fails and the workflow exits without filesystem modification.

A successful fresh-validation result is still point-in-time information rather
than durable authorization for later deletion.

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

Discovery, cleanup validation, selection, confirmation, future fresh validation,
and future deletion are separate workflow stages.

The current implementation is read-only and does not remove filesystem
content.

A discovered directory must pass cleanup-target validation before it can be
reported as currently eligible for selection.

A current eligibility result is an informational snapshot, not durable
authorization.

Selection records user intent for a later workflow stage.

Confirmation records explicit acknowledgement of the exact selected candidate
set for the current invocation.

Neither selection nor confirmation authorizes deletion.

After confirmation, the application performs fresh validation over the exact
selected candidate set. Any failure aborts the whole fresh-validation gate.

This fresh-validation stage remains read-only. When deletion is introduced,
validation must remain immediately adjacent to the destructive operation rather
than relying on an earlier successful result.

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

The scanner, directory-size behaviour, cleanup-target validation, and candidate
selection policy are tested using test-owned data and temporary filesystem trees.

The selection tests cover:

- parsing a single candidate index;
- parsing multiple candidate indexes;
- rejection of empty selection input;
- rejection of zero and negative indexes;
- rejection of malformed input;
- rejection of unsupported whitespace syntax;
- rejection of duplicate indexes;
- rejection of nonexistent indexes;
- rejection of currently ineligible candidates;
- all-or-nothing selection behaviour;
- selection failure when no candidates exist.

Confirmation currently adds only a thin CLI orchestration rule: `-confirm`
requires `-select`. It does not introduce a separate confirmation algorithm or
persistent confirmation state.

Fresh-validation tests cover:

- acceptance of valid confirmed candidates;
- acceptance of multiple valid confirmed candidates;
- rejection of an empty selected-candidate set;
- rejection of a candidate that disappears before fresh validation;
- rejection of changed cleanup-path topology where the test environment permits
  symbolic-link creation;
- preservation of invocation-local candidate indexes in fresh-validation
  failures.

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
- uses invocation-local numeric indexes as the temporary CLI selection
  interface;
- uses an invocation-local boolean `-confirm` flag as the temporary CLI
  confirmation interface;
- does not persist selection or confirmation state;
- does not provide a graphical interface;
- does not delete directories;
- does not support cancellation for long-running scans;
- performs directory-size calculation sequentially;
- does not yet provide dedicated Windows-specific automated coverage for
  junctions and broader reparse-point behaviour;
- does not provide complete Windows reparse-point handling;
- performs a read-only fresh-validation gate after confirmation but does not
  yet make that validation immediately adjacent to an actual deletion
  operation;
- does not yet protect a future deletion operation from filesystem changes
  between its final validation and filesystem modification.

Filesystem deletion will be introduced only after the relevant validation,
Windows-specific path-safety behaviour, explicit confirmation, fresh
pre-deletion validation, and destructive-operation safeguards have been
implemented and tested.
