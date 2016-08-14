# NoDIce Usage

This program compares a Strat-O-Matic Primary Report for all teams and compares it to a list of actual at bats/innings pitched for each player.

If the percentage is above a level the person will get reported.  (This number is currently hard-coded in the code)

Program will run and spit out all players above the high water mark and generated a HTML page

To Gather File

From SOM Game Program
* Build Report: Statistics Menu, Team Statistics
* Select "Primary Stats"
* Select "Each Team"
* From File Menu select "Print To File" 
* Save PRT file generated into running directory

To TEST

go to the \go\src\ndusage directory
run "go test"


To RUN

go to the \go\src\ndusage directory
run "go run ndusage.go"


