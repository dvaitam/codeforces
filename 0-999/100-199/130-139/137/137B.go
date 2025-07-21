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
   seen := make([]bool, n+1)
   distinct := 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= n && !seen[x] {
           seen[x] = true
           distinct++
       }
   }
   // Minimum changes = missing numbers count
   fmt.Println(n - distinct)
}
