# SudokGo

A high speed sudoku solver written in Go.
This rules based solver also describes the steps to solve to aid in learning.

(This includes a very slow generator that is a work in progress.)

## Usage - Solve

The solver takes the puzzle as a string input, a string of 81 characters
that are either the numbers 1-9 or a dash for a blank.

```
$ sudokgo solve -9-6-4---54-78-2-978--291-----9--8-3----6-------8---7---4--6-5---5247---1----8-2-   

 9 |6 4|   
54 |78 |2 9
78 | 29|1  
-----------
   |9  |8 3
   | 6 |   
   |8  | 7 
-----------
  4|  6| 5 
  5|247|   
1  |  8| 2 
D5 must be 4 it is the only possible place in the column
D7 must be 1 it is the only possible place in the column
E4 must be 7 it is the only possible place in the column
C5 must be 8 it is the only possible place in the column
B5 must be 7 it is the only possible place in the row
C9 must be 7 it is the only possible place in the column
C6 must be 9 it is the only possible place in the column
C4 cannot be 1 as it must appear elsewhere in the column
C4 cannot be 6 as it must appear elsewhere in the column
C4 2 is the last possibility
A5 3 is the last possibility
A1 2 is the last possibility
B7 must be 2 it is the only possible place in the row
B4 subset 15 cycle on col, eliminate extras
B6 subset 15 cycle on col, eliminate extras
H4 subset 46 cycle on row, eliminate extras
A8 subset 89 cycle on col, eliminate extras
H8 subset 19 cycle on col, eliminate extras
H1 must be 8 it is the only possible place in the column
G1 cannot be 3 as it must appear elsewhere in the column
E1 subset 13 cycle on row, eliminate extras
D3 must be 5 it is the only possible place in the box
D9 3 is the last possibility
E7 9 is the last possibility
A7 8 is the last possibility
I7 7 is the last possibility
I1 5 is the last possibility
G1 7 is the last possibility
G7 3 is the last possibility
A8 9 is the last possibility
G8 6 is the last possibility
B8 3 is the last possibility
H8 1 is the last possibility
H5 9 is the last possibility
G5 5 is the last possibility
G6 4 is the last possibility
H4 6 is the last possibility
H2 3 is the last possibility
F2 1 is the last possibility
E1 3 is the last possibility
C1 1 is the last possibility
C2 6 is the last possibility
C3 3 is the last possibility
H3 4 is the last possibility
I3 6 is the last possibility
A4 4 is the last possibility
F4 5 is the last possibility
B4 1 is the last possibility
F5 2 is the last possibility
I5 1 is the last possibility
A6 6 is the last possibility
B6 5 is the last possibility
E6 1 is the last possibility
F6 3 is the last possibility
I6 2 is the last possibility
I8 8 is the last possibility
B9 6 is the last possibility
E9 5 is the last possibility
G9 9 is the last possibility
I9 4 is the last possibility
Finished in  582µs
Difficulty score was 24041 (diabolical).
291|634|785
546|781|239
783|529|146
-----------
412|975|863
378|462|591
659|813|472
-----------
824|196|357
935|247|618
167|358|924

```

## Usage - Generate

The generator is currently very simple and takes an optional second argument
that is the difficulty to produce. The default difficulty (if omitted) is
moderate.

```
$ sudokgo generate simple
Generated grid with difficulty simple after 1 attempts.
5  |3  | 82
 42|961|57 
9 7|25 |64 
-----------
 58|  9|13 
 79|435|2 8
3 6|18 |4  
-----------
 93| 16|8  
214|  3|756
685|74 | 19


```

