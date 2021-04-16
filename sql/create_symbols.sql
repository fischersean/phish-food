-- Create ticker symbol table
CREATE TABLE "Symbols" (
	"Ticker"	TEXT UNIQUE,
	"Exchange"	TEXT,
	"FullName"	TEXT,
	"ETF"	boolean,
	PRIMARY KEY("Ticker")
)
