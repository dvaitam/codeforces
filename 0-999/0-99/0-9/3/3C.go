package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   grid := make([][]rune, 3)
   var xCount, oCount int
   for i := 0; i < 3; i++ {
       line, _, err := reader.ReadLine()
       if err != nil {
           fmt.Println("illegal")
           return
       }
       if len(line) != 3 {
           fmt.Println("illegal")
           return
       }
       grid[i] = make([]rune, 3)
       for j, ch := range string(line) {
           grid[i][j] = ch
           if ch == 'X' {
               xCount++
           } else if ch == '0' {
               oCount++
           } else if ch != '.' {
               fmt.Println("illegal")
               return
           }
       }
   }
   // Check counts
   if !(xCount == oCount || xCount == oCount+1) {
       fmt.Println("illegal")
       return
   }
   xWin := win(grid, 'X')
   oWin := win(grid, '0')
   // Both cannot win
   if xWin && oWin {
       fmt.Println("illegal")
       return
   }
   // X wins
   if xWin {
       if xCount == oCount+1 {
           fmt.Println("the first player won")
       } else {
           fmt.Println("illegal")
       }
       return
   }
   // O wins
   if oWin {
       if xCount == oCount {
           fmt.Println("the second player won")
       } else {
           fmt.Println("illegal")
       }
       return
   }
   // No wins
   if xCount+oCount == 9 {
       fmt.Println("draw")
       return
   }
   // Next turn
   if xCount == oCount {
       fmt.Println("first")
   } else {
       fmt.Println("second")
   }
}

func win(grid [][]rune, p rune) bool {
   // Rows and columns
   for i := 0; i < 3; i++ {
       if grid[i][0] == p && grid[i][1] == p && grid[i][2] == p {
           return true
       }
       if grid[0][i] == p && grid[1][i] == p && grid[2][i] == p {
           return true
       }
   }
   // Diagonals
   if grid[0][0] == p && grid[1][1] == p && grid[2][2] == p {
       return true
   }
   if grid[0][2] == p && grid[1][1] == p && grid[2][0] == p {
       return true
   }
   return false
}
