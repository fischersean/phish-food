-- Create count table
CREATE TABLE IF NOT EXISTS Counts (
    Id varchar(25),
    CountDate datetime,
    Ticker varchar(10),
    PostScore int,
    CommentScore int,
    TotalScore float,
    PostMentions int,
    CommentMentions int,
    PRIMARY KEY(Id)
);

