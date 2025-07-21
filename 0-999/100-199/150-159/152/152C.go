package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   names := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i])
   }
   result := 1
   for col := 0; col < m; col++ {
       var seen [26]bool
       count := 0
       for i := 0; i < n; i++ {
           c := names[i][col]
           idx := c - 'A'
           if !seen[idx] {
               seen[idx] = true
               count++
           }
       }
       result = result * count % MOD
   }
   fmt.Println(result)
}
