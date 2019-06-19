# parsepgn

Used to analyse chess PGN files.

## Introduction
The goal here is to get an overview of the number of moves of the game for reaching draw, win, loss for each player.
This is to see if one player wins sooner then the other player.
Code used specifically for lc0 tests to see what PRs could help against trolling in endgames.

## Input
The input of the program is a pgn file containing a match between two players.

## Output explanation
The output is a printed list containing:
movenumber, wins player A, wins player B, draws

For example:
* 47, 2, 0, 0
* 48, 1, 2, 3

The above two lines indicate that for the input pgn given, there were 2 games having 47  ply, and 6 games ending at ply 48.
Of those ending at 47 ply those were both won by player A.
There were also 6 games ending at ply 48, those were 1 win for player A, 2 wins for player B and 3 draws.

Note:from this out you can not see if player A or B played white or black; this is because some games can be adjudicated as win, loss or draw by
the adjudicator (Arena or Cutechess) because of 3 fold draw scoring or adjudicating for a certain amount of moves with CP point difference.




