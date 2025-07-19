package main

import (
   "bufio"
   "fmt"
   "os"
)

const n = 6

const (
   UP = iota
   DOWN
   LEFT
   RIGHT
)

// Move represents a rotation: dir at position p
type Move struct {
   p   int
   dir int
}

var mat [][]int
var table [][]int
var answer []Move

// goMoves rotates row or column p in direction dir by times steps
func goMoves(p, dir, times int) {
   p %= n
   // copy mat to ret
   ret := make([][]int, n)
   for i := 0; i < n; i++ {
       ret[i] = make([]int, n)
       copy(ret[i], mat[i])
   }
   // determine delta steps
   delta := 0
   switch dir {
   case UP, LEFT:
       delta = times
   case DOWN, RIGHT:
       delta = -times
   }
   delta %= n
   if delta < 0 {
       delta += n
   }
   if dir == UP || dir == DOWN {
       // rotate column p
       for i := 0; i < n; i++ {
           ret[i][p] = mat[(i+delta)%n][p]
       }
   } else {
       // rotate row p
       for j := 0; j < n; j++ {
           ret[p][j] = mat[p][(j+delta)%n]
       }
   }
   // record moves
   for t := 0; t < times; t++ {
       answer = append(answer, Move{p, dir})
   }
   mat = ret
}

// move2 performs a specific sequence of moves in the last row algorithm
func move2(i, j int) {
   goMoves(j+2, DOWN, 1)
   goMoves(i, RIGHT, 1)
   goMoves(j+2, UP, 1)
   goMoves(i, RIGHT, 1)
   goMoves(j+2, DOWN, 1)
   goMoves(i, LEFT, 2)
   goMoves(j+2, UP, 1)
}

func main() {
   // initialize target table
   table = make([][]int, n)
   for i := 0; i < n; i++ {
       table[i] = make([]int, n)
       for j := 0; j < n; j++ {
           table[i][j] = i*n + j
       }
   }
   // read input matrix
   mat = make([][]int, n)
   scanner := bufio.NewScanner(os.Stdin)
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           return
       }
       line := scanner.Text()
       mat[i] = make([]int, n)
       for j := 0; j < n && j < len(line); j++ {
           c := line[j]
           if c >= '0' && c <= '9' {
               mat[i][j] = int(c - '0')
           } else {
               mat[i][j] = int(c - 'A' + 10)
           }
       }
   }
   // main solving loop: all but last row
   for i := 0; i+1 < n; i++ {
       for j := 0; j < n; j++ {
           // find current position of desired element
           want := table[i][j]
           row, col := -1, -1
           for r := 0; r < n && row == -1; r++ {
               for c := 0; c < n && col == -1; c++ {
                   if mat[r][c] == want {
                       row, col = r, c
                   }
               }
           }
           if row == i {
               goMoves(col, DOWN, 1)
               row++
               goMoves(row, RIGHT, 1)
               goMoves(col, UP, 1)
               col = (col + 1) % n
           }
           if col == j {
               goMoves(row, RIGHT, 1)
               col = (col + 1) % n
           }
           goMoves(j, DOWN, row-i)
           if col > j {
               goMoves(row, LEFT, col-j)
           }
           if col < j {
               goMoves(row, RIGHT, j-col)
           }
           goMoves(j, UP, row-i)
       }
   }
   // fix last row
   i := n - 1
   // align first element
   for k := 0; k < n; k++ {
       if mat[i][k] == table[i][0] {
           goMoves(i, LEFT, k)
           break
       }
   }
   // align remaining elements
   for j := 1; j < n; j++ {
       if mat[i][j] == table[i][j] {
           continue
       }
       // find current column
       col := -1
       for k := 0; k < n; k++ {
           if mat[i][k] == table[i][j] {
               col = k
               break
           }
       }
       for mat[i][j] != table[i][j] {
           move2(i, col)
           col = (col + 2) % n
           if col < 2 {
               goMoves(i, RIGHT, 1)
               col++
           }
       }
   }
   // output moves
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, len(answer))
   dirs := []string{"U", "D", "L", "R"}
   for _, mv := range answer {
       fmt.Fprintf(out, "%s%d\n", dirs[mv.dir], mv.p+1)
   }
}
