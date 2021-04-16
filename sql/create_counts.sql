-- Create count table
CREATE TABLE "Counts" (
	"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Subreddit"	TEXT NOT NULL,
	"FormatedDate"	TEXT NOT NULL,
	"Hour"	INTEGER NOT NULL,
	"Ticker"	TEXT NOT NULL,
	"PostScore"	INTEGER,
	"CommentScore"	INTEGER,
	"TotalScore"	REAL,
	"PostMentions"	INTEGER,
	"CommentMentions"	INTEGER
)
