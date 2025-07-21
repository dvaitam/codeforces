package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // Characters on the diagonals
   diagChar := grid[0][0]
   // Characters off the diagonals (at least one exists since n>=3)
   otherChar := grid[0][1]
   if otherChar == diagChar {
       fmt.Println("NO")
       return
   }
   // Validate pattern
   for i := 0; i < n; i++ {
       row := grid[i]
       for j := 0; j < n; j++ {
           c := row[j]
           if i == j || i+j == n-1 {
               if c != diagChar {
                   fmt.Println("NO")
                   return
               }
           } else {
               if c != otherChar {
                   fmt.Println("NO")
                   return
               }
           }
       }
   }
   fmt.Println("YES")
}
