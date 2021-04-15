-- Create ticker symbol table
CREATE TABLE IF NOT EXISTS Symbols (
    Ticker varchar(10),
    Exchange varchar(255),
    FullName varchar(255),
    ETF boolean,
    PRIMARY KEY(Ticker)
);

