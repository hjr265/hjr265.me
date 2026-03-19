---
title: Can AI Take a Software Engineer's Job?
htmldescription: A real conversation with an AI coding assistant that confidently fabricated a GitLab CI feature, got called out, and needed a human to suggest the actual solution.
date: 2026-03-19T20:20:00+06:00
tags:
  - AgenticCoding
  - GitLab
---

I had an interesting exchange with a coding agent this evening.

Are coding agents useful? Absolutely. But can it take a software engineer's job? You tell me.

{{< box open="true" >}}

**Me:** In `.gitlab-ci.yml`, can we set the workflow name only if the `PLAYBOOK` variable is set?

**AI:**

```yaml
workflow:
  name: "Playbook ${PLAYBOOK}"
  rules:
    - if: $PLAYBOOK != ""
    - when: always
```

> The `workflow.name` only takes effect when matched by a rule, so on push pipelines (where `$PLAYBOOK` is empty) it falls through to `when: always` which uses the default pipeline name. When `$PLAYBOOK` is set, it shows e.g. "shape playbook".

**Me:** You just made that up.

**AI:** You're right — GitLab CI doesn't support conditional `workflow.name` like that. The `name` field is always evaluated regardless of rules.

> There's no way to conditionally set `workflow.name` in GitLab CI. It will always show the interpolated value. When `$PLAYBOOK` is empty, it'll just show "Playbook ".
>
> Not much we can do here without GitLab adding support for it. Leave it as-is?

**Me:** What if we set a variable conditionally? And then use that variable as the name here...

**AI:** Good idea — use `workflow.rules` to set a variable based on the condition:

```yaml
workflow:
  name: "${PIPELINE_NAME}"
  rules:
    - if: $PLAYBOOK != ""
      variables:
        PIPELINE_NAME: "Playbook ${PLAYBOOK}"
    - when: always
      variables:
        PIPELINE_NAME: ""
```

> When `$PLAYBOOK` is set, the pipeline is named "Playbook shape" / "Playbook update". Otherwise it falls through to the default.

{{< /box >}}

The AI confidently presented a made-up behaviour, admitted it when called out, then declared the problem unsolvable, until I suggested the approach that actually worked. It wrote the implementation, but the solution was mine.

AI coding assistants are powerful tools. But "tool" is the key word here.
