---
title: "Where Is This Screenshot From? Asking Three AIs"
date: 2026-10-21
tags:
  - AI
  - Claude
  - ChatGPT
  - Toph
images:
  - screenshot.png
---

I ran a quick experiment with a few AI services.

I took a screenshot of [Toph](https://toph.co)'s submissions page and pasted it into three of them: Claude, Google Search (not Gemini), and ChatGPT. Each in a fresh chat, each with the same question: "Can you tell me where this screenshot is from?"

I built Toph, so I knew the answer. What I wanted to know was whether the page was recognizable to anyone, or anything, that doesn't stare at it every day.

This is the screenshot I used:

{{< image src="screenshot.png" alt="The screenshot of Toph's submissions page" caption="The screenshot of Toph's submissions page" >}}

Nothing special: a table of recent submissions with problem names, usernames, languages, and verdicts. No logo, no URL bar.

Then I asked them how they figured it out.

## Google Search

{{< image src="google.png" alt="Google Search's answer" caption="Google Search's answer" >}}

Google got it right on the first try. It identified the page as Toph's submissions page and picked out a few of the problem names and verdicts.

{{< image src="google-how.png" alt="Google Search explaining how it figured it out" caption="Google Search explaining how it figured it out" >}}

When I asked how, it pointed to three things: the problem names ("Clock Math", "Where is the Ghost"), the layout and dark theme, and the submission IDs. The problem names make sense. They're searchable, and they lead straight back to Toph.

The submission IDs are a stranger claim. It said the 7-digit IDs starting with 21383xx "perfectly align" with the volume of submissions Toph handles. The number is about right: submission IDs on Toph are sequential, and Toph [crossed one million submissions in 2022](https://blog.toph.co/spotlight/one-million-submissions/). But did Google actually know that, or did it reason backwards from the answer?

## Claude

{{< image src="claude.png" alt="Claude's answer" caption="Claude's answer" >}}

Claude got it right too, and explained its reasoning without being asked. When I did ask, it described the whole thing as a process of elimination.

{{< image src="claude-how.png" alt="Claude explaining how it guessed" caption="Claude explaining how it guessed" >}}

The layout told it this was an online judge. The question was which one. And the strongest clue was the text "CPU limit exceeded". Most judges say "Time limit exceeded". Toph doesn't.

That one is on me. When we named Toph's verdicts, I went with "CPU limit exceeded" because that's what we actually limit: how much CPU time a program uses, not how much wall-clock time passes. A program that sits idle waiting on I/O isn't burning CPU. The verdict says what was measured.

The rest of the details only had to fit: the verdict format ("on test 13 (1.0 GB)"), the Bangladeshi flavour of the usernames and problem titles, and the `vjudge` accounts, which are submissions relayed through Virtual Judge.

And it hedged: no single clue settled it, and a screenshot can be misleading.

## ChatGPT

{{< image src="chatgpt.png" alt="ChatGPT's answer" caption="ChatGPT's answer" >}}

ChatGPT... well, it did something. Mr. Altman's toy seems to have had a bad day. It produced a long answer about building a "constraint-native computation engine", complete with a diagram, that had nothing to do with my question.

{{< image src="chatgpt-recovery.png" alt="ChatGPT recovering" caption="ChatGPT recovering" >}}

When I asked what it was even talking about, it apologized and described the screenshot correctly, verdicts and all. It had clearly read the image. It even listed "CPU limit exceeded", but didn't treat it as a clue. It couldn't name the site until I said yes to letting it dig further. Then it searched for "Where is the Ghost" and found Toph.

So in the end, ChatGPT took the same route as Google: search for a problem name.

Remember, this was a fresh chat. The off-topic answer reads like the continuation of some other conversation.

## Clever Indeed

Two routes to the same answer. Google and ChatGPT matched text they could search for. Claude recognized what was unusual about the wording.

One caveat: when you ask a model how it figured something out, you get an explanation, not a recording of what actually happened inside it. Google's bit about submission IDs may be an example of that. Still, the explanations line up with what's on the screen.

And it's a neat reminder that the small choices you make when building something, like what to call a verdict, end up becoming its fingerprint.
