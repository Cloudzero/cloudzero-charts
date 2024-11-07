# Release Notes Template

Pro-tip: look through the github diff between the previous release to see what's changed. The commit titles should give an outline of what's happened.

## Upgrade Steps

- List out, as concretely as possible, any steps users have to take when they upgrade beyond just dumping the dependency.
- Write pseudocode that highlights what code should change and how.
- Call out if users are recommended to upgrade because of known problems with older releases.
- Preferably, there's nothing here.

## Breaking Changes

- A complete list of breaking changes (preferably there are none, unless this is a major version).

## New Features

- Describe the new feature and when/why to use it. Add some pictures! Call out any caveats/warnings? Is it a beta feature?

## Bug Fixes

- Call out any existing feature/functionality that now works as intended or expected.

## Improvements

- Improvements/enhancements to a workflow, performance, logging, error messaging, or user experience

## Other Changes

- Other miscellaneous changes that don't fit into any of the above categories. Try to leave this empty - ideally, all changes fit into the categories above

------

#### Copy and paste this template

```markdown
## [0.29.0](https://github.com/cloudzero/cloudzero-charts/compare/v0.27.0...v0.29.0) (2024-11-07)

> Description

### Upgrade Steps
* [ACTION REQUIRED]
* 

### Breaking Changes
* 
* 

### New Features
* 
* 

### Bug Fixes
* 
* 

### Performance Improvements
* 
* 

### Other Changes
* 
* 
```

Example:

```markdown
## [0.29.0](https://github.com/cloudzero/cloudzero-charts/compare/v0.27.0...v0.29.0) (2024-11-07)


### Performance Improvements

* **dependencies:** Bump dependencies  4a4ee13

### Other Changes

* **chore(conventionalChangelog):** Add Conventional Changelog  aafcdd9
* **docs(CHANGELOG):** Add changelog  e2c7435
```