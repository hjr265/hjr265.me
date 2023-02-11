// The following JavaScript code are meant to be run from within the MongoDB Shell (mongosh).

// Populate the rounds collection.
db.rounds.insert({"player1": "Alice", "player2": "Bob", "winner": "Alice"})
db.rounds.insert({"player1": "Alice", "player2": "Carol", "winner": "Alice"})
db.rounds.insert({"player1": "Bob", "player2": "Carol", "winner": "Bob"})
db.rounds.insert({"player1": "Carol", "player2": "Ted", "winner": "Carol"})
db.rounds.insert({"player1": "Ted", "player2": "Bob", "winner": "Ted"})
db.rounds.insert({"player1": "Alice", "player2": "Ted", "winner": "Alice"})

// Populate the profiles collection.
db.profiles.insert({"player": "Alice", "planet": "Mars"})
db.profiles.insert({"player": "Bob", "planet": "Mars"})
db.profiles.insert({"player": "Carol", "planet": "Venus"})
db.profiles.insert({"player": "Ted", "planet": "Mars"})

// Add a unique index on the player field in the leaderboard collection.
db.leaderboard.ensureIndex({ player: 1 }, { unique: true })

// Run an aggregation pipeline for each player to aggregate round statistics.
['Alice', 'Bob', 'Carol', 'Ted'].forEach(player => {
	db.rounds.aggregate([
		{ "$match": { "winner": player } },
		{ "$count": "won" },
		{ "$addFields": { "player": player } },
		{ "$merge": { "into": "leaderboard", "on": "player", "whenMatched": "merge", "whenNotMatched": "insert" } }
	])
})

// Run an aggregation pipeline for each player to denormalize profile fields.
['Alice', 'Bob', 'Carol', 'Ted'].forEach(player => {
	db.profiles.aggregate([
		{ "$match": { "player": player } },
		{ "$project": { "_id": 0, "player": 1, "planet": 1 } },
		{ "$merge": { "into": "leaderboard", "on": "player", "whenMatched": "merge", "whenNotMatched": "insert" } }
	])
})

// Query leaderboard with rank.
db.leaderboard.aggregate([
  { 
    "$setWindowFields": {
      "sortBy": { "won": -1 },
      "output": {
        "rank":  {"$documentNumber": {}}
      },
    }
  }
])

// Query leaderboard, with rank, for players from Mars only.
db.leaderboard.aggregate([
  { "$match": { "planet": "Mars" } },
  {
    "$setWindowFields": {
      "sortBy": { "won": -1 },
      "output": {
        "rank":  {"$documentNumber": {}}
      },
    }
  }
])

// Query leaderboard with dense rank.
db.leaderboard.aggregate([
  { 
    "$setWindowFields": {
      "sortBy": { "won": -1 },
      "output": {
        "rank": { "$denseRank": {} }
      },
    }
  }
])
