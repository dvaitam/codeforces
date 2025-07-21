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
   var p int
   fmt.Fscan(reader, &p)
   for i := 0; i < p; i++ {
       var lvl int
       fmt.Fscan(reader, &lvl)
       if lvl >= 1 && lvl <= n {
           seen[lvl] = true
       }
   }
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var lvl int
       fmt.Fscan(reader, &lvl)
       if lvl >= 1 && lvl <= n {
           seen[lvl] = true
       }
   }
   for i := 1; i <= n; i++ {
       if !seen[i] {
           fmt.Println("Oh, my keyboard!")
           return
       }
   }
   fmt.Println("I become the guy.")
}
