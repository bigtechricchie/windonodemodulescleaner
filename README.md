# WindoNodeModulesCleaner

Safe `node_modules` discovery and cleanup for Windows.

WindoNodeModulesCleaner is a Windows utility written in Go for discovering
`node_modules` directories beneath a user-selected directory tree.

Built by Edward Siwek.

The application is being designed around a controlled cleanup workflow:

```text
select root
    ↓
scan
    ↓
review discovered targets
    ↓
select targets
    ↓
confirm deletion
    ↓
validate
    ↓
delete
    ↓
report results
```

## Project status

Early development.

The initial implementation will focus on establishing correct filesystem
discovery before introducing destructive operations.

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

Before a destructive operation is attempted, the selected target must still
satisfy the application's cleanup rules.

## Current limitations

The application is in early development and does not yet perform filesystem
cleanup.

Do not interpret the current repository as a completed release.
