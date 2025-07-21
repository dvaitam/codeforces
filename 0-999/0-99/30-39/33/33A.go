package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // caps[r] holds the maximum times row r can be used
   caps := make([]int64, m+1)
   // initialize with large value
   const inf = int64(1) << 60
   for i := 1; i <= m; i++ {
       caps[i] = inf
   }
   // read each tooth's row and viability
   for i := 0; i < n; i++ {
       var r int
       var c int64
       if _, err := fmt.Fscan(reader, &r, &c); err != nil {
           return
       }
       if c < caps[r] {
           caps[r] = c
       }
   }
   // sum of capacities
   var total int64
   for i := 1; i <= m; i++ {
       total += caps[i]
   }
   // cannot exceed k
   if total > k {
       total = k
   }
   fmt.Println(total)
}
