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
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var row string
       fmt.Fscan(reader, &row)
       grid[i] = []byte(row)
   }
   // Try placing crosses greedily
   for i := 1; i < n-1; i++ {
       for j := 1; j < n-1; j++ {
           if grid[i][j] == '#' &&
              grid[i-1][j] == '#' &&
              grid[i+1][j] == '#' &&
              grid[i][j-1] == '#' &&
              grid[i][j+1] == '#' {
               // mark as used
               grid[i][j] = '.'
               grid[i-1][j] = '.'
               grid[i+1][j] = '.'
               grid[i][j-1] = '.'
               grid[i][j+1] = '.'
           }
       }
   }
   // Check if any '#' remains
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if grid[i][j] == '#' {
               fmt.Println("NO")
               return
           }
       }
   }
   fmt.Println("YES")
}
