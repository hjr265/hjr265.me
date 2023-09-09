---
title: My GitHub Status Is A Clock
date: 2023-09-09T00:30:00+06:00
tags:
  - GitHub
  - AlmostFunny
  - 100DaysToOffload
  - Go
  - GraphQL
toc: yes
---

On GitHub, you can set an emoji and a short message as your status. I don't browse GitHub much. But through what little I do, I see so many profiles with ":dart: Focusing".

Why not make it a clock? An almost functional clock.

That's exactly what I did.: [github.com/hjr265](https://github.com/hjr265)

{{< image src="cover.png" alt="Screenshot of hjr265's GitHub profile" captionHTML=`Screenshot of <a href="https://github.com/hjr265">my GitHub profile</a>` >}}

I wrote a Go program that I can leave running. The program updates my GitHub status with one of the clock emojis (one that is close to the current time) and a message saying something like "Twelve o'clock" or "Half past twelve".

This blog post is about how it works. You can find the [source code for this on GitHub](https://github.com/hjr265/mghsiac).

## Time to Emoji

Given a Go `time.Time` we need to be able to choose one of the 24 clock emojis. The emojis are named ":clock12:", ":clock1230:", ":clock1:", ":clock130:", and so on.

<center>🕛 🕧 🕐 🕜 🕑 🕝 🕒 🕞  🕟  🕠 🕡 🕖 🕢 🕗 🕣  🕤 🕥 🕚  🕦🕒</center>

The name follows the pattern ":clock{h}:" or ":clock{h}30:", where "{h}" is the hour.

``` go
func timeToEmoji(t time.Time) string {
  h := t.Hour()
  if h > 12 {
    h -= 12
  }
  if h == 0 {
    h = 12
  }
  m := t.Minute()
  m = (m / 30) * 30
  if m == 0 {
    return fmt.Sprintf(":clock%d:", h)
  }
  return fmt.Sprintf(":clock%d%d:", h, m)
}
```

For the hour part, we first go from 24-hour to 12-hour. Since we may end up with `h == 0`, we just change that to 12.

For the minute part, we round down to the nearest 30.

If then the minute part is 0, we return ":clock{h}:". Otherwise, we return ":clock{h}{m}:".

## Time to Message

Depending on the current time, we set one of the two messages: "{Hour} o'clock" or "Half past {hour}".

``` go
var hourWord = map[int]string{
  1:  "One",
  2:  "Two",
  // ... 3 to 10 omitted for brevity.
  11: "Eleven",
  12: "Twelve",
}

func timeToMessage(t time.Time) string {
  h := t.Hour()
  if h > 12 {
    h -= 12
  }
  if h == 0 {
    h = 12
  }
  m := t.Minute()
  m = (m / 30) * 30
  if m == 0 {
    return fmt.Sprintf("%s o'clock", hourWord[h])
  }
  return fmt.Sprintf("Half past %s", strings.ToLower(hourWord[h]))
}
```

Here we follow the math similar to timeToEmoji. Except, we return "{Hour} o'clock" when rounded down minutes is 0. Otherwise, we return "Half past {hour}".

We keep a map of twelve numbers (1 to 12) in words. We use words for hours here instead of numbers.

## Putting It Together

The program goes something like this:

- Set up a GitHub GraphQL client. We use `github.com/shurcooL/githubv4` for this.
- Loop:
  - Update status based on current time.
  - Sleep until another update is needed.

``` go
func main() {
  // Set up a GitHub GraphQL client.
  src := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
  )
  httpclient := oauth2.NewClient(context.Background(), src)
  client := githubv4.NewClient(httpclient)

  for {
    now := time.Now() // You can adjust the timezone with now.In(loc).

    // Use the current time to get the emoji and the message.
    emoji := timeToEmoji(now)
    message := timeToMessage(now)

    // Update the user's status using the GitHub GraphQL client.
    // Error handling omitted for brevity.
    m := struct {
      ChangeUserStatus struct {
        Status struct {
          Emoji   graphql.String
          Message graphql.String
        }
      } `graphql:"changeUserStatus(input: $input)"`
    }{}
    client.Mutate(context.Background(), &m, githubv4.ChangeUserStatusInput{
      Emoji:   githubv4.NewString(githubv4.String(emoji)),
      Message: githubv4.NewString(githubv4.String(message)),
    }, nil)

    // Sleep until another update is needed.
    time.Sleep(sleepDuration(now))
  }
}
```

The `sleepDuration` function looks at the `time.Time` passed to it and calculates the number of minutes left until the next 30-minute or 60-minute mark (whichever is the closest).

## Get A Clock

You can find the [source code for this on GitHub](https://github.com/hjr265/mghsiac).

You can also install this program with the `go install` command and run it like so:

``` sh
go install github.com/hjr265/mghsiac@latest
GITHUB_TOKEN=... mghsiac
```
