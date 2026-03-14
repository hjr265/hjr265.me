---
title: "My First Fully Agentic Coding Project: GitTop"
date: 2026-03-14T14:20:00+06:00
tags:
  - Go
  - Git
  - AgenticCoding
  - TUI
  - ClaudeCode
toc: yes
---

# Building GitTop With Claude Code: A Weekend of Agentic Coding

I have been running [Toph](https://toph.co/) for over ten years. Somewhere along the way, I started wondering: what hours of the day do I actually work on it?

Commit timestamps felt like the right place to look. A quick one-off script could answer the question, and there are Git stats tools that spit out HTML reports. But I thought: this could be a good excuse to build a TUI application. Something like htop, but for a Git repository instead of system metrics.

That became GitTop.

## How It Was Built

This weekend I tried something I had not done before: fully agentic coding with Claude Code. I described what I wanted, guided it through one feature at a time, and let it write everything.

Twenty-six commits later, GitTop was done.

The stack is Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework, [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal styling, and [go-git](https://github.com/go-git/go-git) for reading the repository without shelling out to `git`. The result is a seven-page dashboard:

| # | Page | What it shows |
|---|-----|---------------|
| 1 | **Summary** | KPI cards + braille area chart |
| 2 | **Activity** | Heatmap, by-hour, by-weekday, by-month distributions |
| 3 | **Contributors** | Split panel: ranked list + per-author detail |
| 4 | **Branches** | Sortable table with ahead/behind counts |
| 5 | **Files** | Largest, most churn, most authors, stalest, language breakdown |
| 6 | **Releases** | Tag timeline and cadence chart |
| 7 | **Commits** | Scrollable log with diff viewer and fuzzy search |
{.table}

<br>

{{< image src="screenshot.png" alt="GitTop Activity page showing commit distribution by hour" caption="GitTop Activity page showing commit distribution by hour" >}}

The activity page answered my original question immediately. I commit to Toph mostly between 10:00 and 16:00, with a clear peak around noon. Apparently, I am a morning-to-afternoon programmer on that project.

## The Good Part

I went in expecting the code to work. I did not expect it to be particularly clever.

But most of the time it was.

The filter system is one example. Rather than a simple substring search, I wanted a proper DSL with structured queries like `author:"alice" and path:*.go` or `branch:main and not path:vendor`. Claude Code's first attempt was a hand-rolled parser in raw Go. I asked it to use [Participle](https://github.com/alecthomas/participle) instead, a parser combinator library. And it rebuilt the whole thing around it, compiling queries into an AST of filter nodes each implementing a `Match(*CommitInfo) bool` interface.

The charts were another surprise, though perhaps not in the way I expected. 
The summary page started as a blocky bar chart. I asked Claude Code to redo 
it using Unicode braille characters instead. What it came back with uses the 
U+2800 block, where each character cell encodes a 2×4 grid of dots, giving 
sub-character resolution across the chart. An 80-column terminal effectively 
has 160-column resolution for those charts. The bar charts across the other 
pages use the Unicode block elements (▏▎▍▌▋▊▉█) for fractional-width 
rendering, so a bar representing 3.5 units actually looks 3.5 characters 
wide rather than snapping to 3 or 4. I asked for braille. I did not expect 
it to also get the fractional bars right unprompted.

The branch filter architecture was also interesting. Adding a `branch:` field to every `CommitInfo` struct would have been expensive and messy. Instead, it pre-computes a hash set by walking all commits reachable from matching branches, then filters by hash membership. The data model stays clean. I did not suggest this approach. It arrived at it on its own.

None of this is magic. These are things a good Go programmer would do. But seeing an LLM make those choices, unprompted in some cases, across a project built in a single session, was genuinely surprising.

## The Odd Part

Here is the thing I keep turning over in my head: I do not feel like I own this project.

Every coding project I have worked on before has felt mine in a clear way. I made the decisions. I wrote the code. I hit the walls and figured out how to climb over them. Even the bad decisions are mine, and they teach me something.

With GitTop, I guided Claude Code in small steps toward exactly what I wanted. The feature ideas were mine. The direction was mine. The judgment calls, what to build, in what order, when something felt wrong, were mine. But the code was not written by me.

The result is a tool I am genuinely happy with. At around 4,800 lines of Go, it is substantial, the visual output is exactly what I had in mind, and the filter DSL is nicer than it needed to be. I would have been satisfied writing this myself over a few weekends.

Yet something about it feels different.

I am not sure what to do with that feeling yet. Maybe ownership is more about authorship than I had realized. Maybe it will shift as I use the tool and make it mine through use rather than construction. Or maybe this is just what building software with modern tools feels like, and the distinction between guiding and writing will keep getting harder to hold onto.

## Installing GitTop

``` sh
go install github.com/hjr265/gittop@latest
```

Or build from source:

``` sh
git clone https://github.com/hjr265/gittop.git
cd gittop
go build -o gittop .
```

Then point it at any repository:

``` sh
gittop               # current directory
gittop /path/to/repo
```

For the best experience with the block character bar charts, a terminal like Kitty, Alacritty, or WezTerm with a monospace font (JetBrains Mono or Iosevka work well) gives the cleanest rendering.

The source is on [GitHub](https://github.com/hjr265/gittop).
