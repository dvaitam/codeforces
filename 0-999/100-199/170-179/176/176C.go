package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, x1, y1, x2, y2 int
   if _, err := fmt.Fscan(in, &n, &m, &x1, &y1, &x2, &y2); err != nil {
       return
   }
   // Manhattan distance
   d := abs(x1 - x2) + abs(y1 - y2)
   // One-dimensional board
   if n == 1 {
       if d <= 4 {
           fmt.Println("First")
       } else {
           fmt.Println("Second")
       }
       return
   }
   if m == 1 {
       if d <= 4 {
           fmt.Println("First")
       } else {
           fmt.Println("Second")
       }
       return
   }
   // Two-dimensional board: first player can always force a win
   fmt.Println("First")
}
