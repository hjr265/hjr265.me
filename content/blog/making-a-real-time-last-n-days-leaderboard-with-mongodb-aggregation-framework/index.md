---
title: 'Making a Real-time "Last N Days" Leaderboard with MongoDB Aggregation Framework'
date: 2023-10-22T18:45:00+06:00
tags:
  - MongoDB
  - AggregationFramework
  - 100DaysToOffload
toc: yes
---

On Toph, there is a [leaderboard of top solvers](https://toph.co/leaderboard). Without any filters, this leaderboard shows the list of programmers who solved the most programming problems in the last seven days. Toph updates the leaderboard in real-time.

{{< image src="overview.png" alt="Real-time \"Last N Days\" Leaderboard with MongoDB Aggregation Framework" >}}

There are a few ways to build a leaderboard like this one.

## The Naive Way

The easy way is to run a daily or hourly cron that aggregates all the solutions submitted over the last seven days and regenerates the leaderboard. However, the leaderboard will be far from real-time with this approach.

If you are to use cron, you are limiting updates to specific intervals. The leaderboard will stay stale until the next update.

If you make the cron trigger too frequently, you risk wasting resources by running the leaderboard calculation even when no update is needed.

On the other hand, a cron that triggers at a shorter interval than it takes for the calculation to complete, then you are now facing a different challenge.

This approach works only when you do not have that much data to process. For example, this was okay in Toph several years ago when we didn't get many submissions every hour (or even a day).

## Trigger Updates Only As Needed

The better approach is to trigger an update only when needed. For example, in Toph, we recalculate an entry in the leaderboard after that programmer makes a new submission.

At first, this sounds like a good plan. But it comes with a challenge.

Imagine a leaderboard that has ten people on it right now. Eight remain active over the next several days, but two do not.

After a day, you now have two stale entries on the leaderboard because you update the leaderboard entries only when there is any activity.

## The Solution

The solution to this challenge is to look forward.

Instead of calculating the score for today, based on the last seven days, calculate seven scores:

- For today, based on the last seven days
- For tomorrow, based on the last six days
- For the day after tomorrow, based on the last five days
- And so on

The idea is that you will be calculating seven scores against the dates you will be using them ahead of time.

You will recalculate these scores every time there is an activity.

If there is no activity, and you are querying the leaderboard the next day, you will use the sum that accounts for the date change and lack of programmer activity.

## Example

Let's go through an example. 

{{< image src="solutions.png" alt="Number of solutions by a programmer in the last seven days" caption="Number of solutions by a programmer in the last seven days" >}}

A programmer has solved five programming problems today, three yesterday, three the day before that, zero the day before that, one the day before that, two the day before that and four the day before that.

On the leaderboard today, this programmer's score would be 18.

If this programmer didn't solve any other programming problems now, the programmer's score tomorrow would be 14. The day after tomorrow, the programmer's score would be 12. And so on.

By calculating this information ahead of time, we can do away with the cron entirely.

{{< image src="scores.png" alt="Scores of a programmer for the next seven days" caption="Scores of a programmer for the next seven days" >}}

## Updating the Leaderboard Entry with MongoDB Aggregation Framework

Let's assume we have the daily solution count of the programmers in the `solutionCounts` collection. And it looks something like this:

``` json
{
  "programmer": "Alice",
  "daily_counts": {
    // ...
    "18275": 4,
    "18276": 2,
    "18277": 1,
    "18278": 0,
    "18279": 3,
    "18280": 3,
    "18281": 5
  }
}
```

Here, in the "solutionCounts" object, each key represents the nth day since the Unix epoch. "0" would be January 1, 1970.

Let's assume that today is the 18281th day since the Unix epoch.

To update the leaderboard, run an aggregation pipeline like this:

``` json
[
  { "$match": { "programmer": "Alice" } },
  {
    "$project": {
      "_id": 0,
      "sums_of_daily_counts": {
        "18281": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280", "$daily_counts.18279", "$daily_counts.18278", "$daily_counts.18277", "$daily_counts.18276", "$daily_counts.18275" ] }, // Last 7 days
        "18282": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280", "$daily_counts.18279", "$daily_counts.18278", "$daily_counts.18277", "$daily_counts.18276" ] }, // Last 6 days
        "18283": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280", "$daily_counts.18279", "$daily_counts.18278", "$daily_counts.18277" ] }, // Last 5 days
        "18284": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280", "$daily_counts.18279", "$daily_counts.18278" ] }, // Last 4 days
        "18285": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280", "$daily_counts.18279" ] }, // Last 3 days
        "18286": { "$sum": [ "$daily_counts.18281", "$daily_counts.18280" ] }, // Last 2 days
        "18287": { "$sum": [ "$daily_counts.18281" ] } // Last 1 day
      },
    },
  }
]
```

This pipeline is calculating seven scores for the seven days starting from today.

## Querying the Leaderboard

If you were to render the leaderboard today, which we are assuming is the 18281st since epoch, and sort the results by decreasing scores, then the query would look like this:

``` json {linenos=false}
{ "sums_of_daily_counts.18281": { "$gt": 0 } }
```

And sort by the `sums_of_daily_counts.18281` field:

``` json {linenos=false}
{ "sums_of_daily_counts.18281": -1 }
```

If you were to query the leaderboard tomorrow, then you would query:

``` json {linenos=false}
{ "sums_of_daily_counts.18282": { "$gt": 0 } }
```

And sort by the `sums_of_daily_counts.18282` field:

``` json {linenos=false}
{ "sums_of_daily_counts.18282": -1 }
```

And so on.

## Indexing the Leaderboard

How do you index these fields? Use [MongoDB wildcard index](https://www.mongodb.com/docs/manual/core/indexes/index-types/index-wildcard/).

``` json {linenos=false}
{ "sums_of_daily_counts.$**": -1 }
```

## Wrap Up

That is how you can maintain a real-time "last N days" leaderboard with the MongoDB aggregation framework without relying on any cron.

I hope you found this blog post interesting. If you have any questions, please ask in the comments area below.

_Earlier this year, I wrote a blog post on implementing a [ranked leaderboard with the MongoDB aggregation framework](/blog/making-ranked-leaderboards-with-mongodb-aggregation-framework/)._
