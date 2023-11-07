---
title: 'Building Advanced Search With Go and MongoDB'
date: 2023-11-07T13:00:00+06:00
tags:
  - Go
  - MongoDB
  - 100DaysToOffload
toc: yes
---

I like software that allows advanced search. Advanced search is where you can use flags to indicate what you want.

Take email software as an example.

It may allow you to enter `from:Bloo to:Mac subject:Imagination` and find all emails that were sent from Bloo to Mac and has the word "Imagination" in the subject line.

But how do you implement something like this in Go? That is what this blog post is about.

## Making an Advanced Search Query Parser

You may want to use regular expressions to parse the advanced search query and use that to build a database query. But I assure you the idea will backfire as you attempt to support complicated search queries.

Instead, let's use a real parser.

Meet [Participle](https://github.com/alecthomas/participle).

Participle is a parser library for Go. You can define the language (e.g. your advanced search query) you want to parse using Go types.

``` go
// A query represents the entire advanced search query.
type Query struct {
  Fields []Field
}

// Field represents either a flag or an argument. Flags are key-value pairs, while arguments are non-flag words in the search query.
// For example:
//   the fox kind:"quick brown" action:jumps over:"lazy dog"
// Here, `the` and `fox` are two arguments, while the `kind:"quick brown"`, `action:jumps` and `over:"lazy dog"` are flags.
type Field struct {
  Flag *Flag 
  Arg  *string
}

// Flag represents a key-value pair.
type Flag struct {
  Key   string
  Value Value
}

// Value is a number, a boolean, or a string.
type Value struct {
  Number *float64
  Bool   *bool   
  String *string 
}
```

These Go structs represent our advanced search query language.

But we need to tell the parser how to parse the query string into these structures. To do that, we need to add a few field tags:

``` go
type Query struct {
  Fields []Field `@@*` // A query is a slice of fields.
}

type Field struct {
  Flag *Flag   `( @@ |`                  // A field may be a flag, or
  Arg  *string `( @String | @Ident ) )`  // ... an argument.
}

type Flag struct {
  Key   string `@Ident ":"` // A flag is composed of a key, followed by a colon
  Value Value  `@@`         // ... and then a value.
}

type Value struct {
  Number *float64 `@Float | @Int |`         // A value is a float, or an integer, or
  Bool   *bool    `( @"true" | "false" ) |` // ... a boolean, or
  String *string  `@String | @Ident`        // ... a string or an identifier.
}
```

We also need to define a lexer for the parser. To build a lexer with simple rules, we need a few regular expressions:

``` go
var queryLexer = lexer.MustSimple([]lexer.SimpleRule{
  {"whitespace", `\s+`},
  {"Float", `\d+\.\d*`},
  {"Int", `\d+`},
  {"String", `"(\\"|[^"])*"`},
  {"Ident", `[^:\s]+`},        // Colon and whitespace cannot be a part of an identifier.
  {"Punct", `[:]`},            // Colon is the only punctuation character in our advanced search query language.
}),
```

And now, we can have a parser:

``` go
var parser = participle.MustBuild[Query](
  participle.Lexer(queryLexer),
)
```

## Parsing a Query

You can now use the parser to parse your advanced search queries:

``` go
// Error handling omitted for brevity.

q, _ := parser.ParseString("", `the fox kind:quick kind:brown action:jumps over:"lazy dog"`)
```

And `q` will be as follows:

``` go
&Query{
  Fields: []Field{
    {Arg:  Value{String: &"the"}},
    {Arg:  Value{String: &"fox"}},
    {Flag: {Key: "kind",   Value: Value{String: &`quick`}}},
    {Flag: {Key: "kind",   Value: Value{String: &`brown`}}},
    {Flag: {Key: "action", Value: Value{String: &`jumps`}}},
    {Flag: {Key: "over",   Value: Value{String: &`"lazy dog"`}}},
  }
}
```

Notice that the string value for the flag "over" retains the double quotes. You can remove the double quotes using the `participle.Unquote` option with the parser:

``` go
var parser = participle.MustBuild[Query](
  participle.Lexer(queryLexer),
  participle.Unquote("String"),
)
```

The parsed query will be like so:

``` go
&Query{
  Fields: []Field{
    {Arg:  Value{String: &"the"}},
    {Arg:  Value{String: &"fox"}},
    {Flag: {Key: "kind",   Value: Value{String: &`quick`}}},
    {Flag: {Key: "kind",   Value: Value{String: &`brown`}}},
    {Flag: {Key: "action", Value: Value{String: &"jumps"}}},
    {Flag: {Key: "over",   Value: Value{String: &"lazy dog"}}},
  }
}
```

## Preparing the MongoDB Query

At this point, you can use a parsed query to create a filter document for your MongoDB query.

A function like the following can do the job:

``` go
func makeFilter(query *Query) bson.M {
  m := bson.M{}

  // Aggregate all arguments from the query for text search.
  text := []string{}
  for _, f := range query.Fields {
    if f.Arg != nil {
      text = append(text, strconv.Quote(*f.Arg))
    }
  }
  if len(text) > 0 {
    m["$text"] = bson.M{"$search": strings.Join(text, " ")}
  }

  // Group the flags by key and prepare filter elements.
  for _, key := range []string{"kind", "action", "over"} {
    values := []any{}
    for _, f := range query.Fields {
      if f.Flag == nil || f.Flag.Key != key {
        continue
      }
      switch {
      case f.Flag.Value.Bool != nil:
        values = append(values, *f.Flag.Value.Bool)
      case f.Flag.Value.Number != nil:
        values = append(values, *f.Flag.Value.Number)
      case f.Flag.Value.String != nil:
        values = append(values, *f.Flag.Value.String)
      }
    }
    if len(values) > 1 {
      m[key] = bson.M{"$all": values} // Multiple flags with the same key. Require all values to match.
    } else if len(values) == 1 {
      m[key] = values[0]
    }
  }

  return m
}
``` 

Note that the above function doesn't perform much validation of what flag values are present in the advanced search query. You should write a function that carefully maps query flags into filter elements, ensuring only permitted flag keys/values are used.

If you use the `makeFilter` function on the example query we began with:

``` go
// Error handling omitted for brevity.

q, _ := parser.ParseString("", `the fox kind:quick kind:brown action:jumps over:"lazy dog"`)
m := makeFilter(q)
```

You should get a `bson.M` like this:

``` go
bson.M{
  "$text": bson.M{"$search": []string{"the", "fox"}}
  "kind": bson.M{"$all": []any{"quick", "brown"}},
  "action": "jumps",
  "over": "lazy dog",
}
```

Now use that filter in a MongoDB `find` call and see the results of your advanced search query.

``` go
// Error handling omitted for brevity.

q, _ := parser.ParseString("", `the fox kind:quick kind:brown action:jumps over:"lazy dog"`)
filter := makeFilter(q)

cursor, _ := mongoClient.Database("db").Collection("collection").
    Find(ctx, filter)
```

## Wrap Up

There is so much more you can do with [Participle](https://github.com/alecthomas/participle) when parsing an advanced search query. The parser package is simple and well done.

Once you define the language and can parse advanced search queries, using it to build [MongoDB `find` queries](https://www.mongodb.com/docs/manual/reference/method/db.collection.find/) should be easy.

Make sure you have the proper set of indexes defined so that the advanced search queries are as fast as MongoDB can be.
