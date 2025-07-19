package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, p int
   if _, err := fmt.Fscan(in, &n, &m, &p); err != nil {
       return
   }
   var av, bv int
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       if x%p != 0 {
           av = i
       }
   }
   for i := 0; i < m; i++ {
       var x int
       fmt.Fscan(in, &x)
       if x%p != 0 {
           bv = i
       }
   }
   // output sum of last non-divisible indices
   fmt.Println(av + bv)
}
