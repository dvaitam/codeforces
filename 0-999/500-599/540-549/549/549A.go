package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   count := 0
   for i := 0; i+1 < n; i++ {
       for j := 0; j+1 < m; j++ {
           // Collect 2x2 letters
           a := grid[i][j]
           b := grid[i][j+1]
           c := grid[i+1][j]
           d := grid[i+1][j+1]
           // Count frequencies
           freq := [256]int{}
           freq[a]++
           freq[b]++
           freq[c]++
           freq[d]++
           if freq['f'] == 1 && freq['a'] == 1 && freq['c'] == 1 && freq['e'] == 1 {
               count++
           }
       }
   }
   fmt.Println(count)
}
