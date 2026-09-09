# Issue Acceptance Criteria Template

Every issue linked to a PR must contain acceptance criteria before the PR is opened.

## Minimal template
```markdown
## Goal
[One sentence: what problem does this solve?]

## Acceptance Criteria
- [ ] AC1: [specific, testable condition]
- [ ] AC2: [specific, testable condition]
- [ ] AC3: [specific, testable condition]

## Out of scope
[Explicitly list what will NOT be done in this issue]
```

## In review: how to check
- `has_acceptance_criteria` field in PR audit data must be `true`
- If `false`, flag as HIGH spec gap and block merge until AC is added to the issue
- A single prose sentence does NOT count as acceptance criteria
- BDD Given/When/Then, numbered requirements with 'must/should/shall', or checkbox lists all count