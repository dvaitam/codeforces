package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   diff := make(map[int]int)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       diff[x]++
   }
   for i := 0; i < m; i++ {
       var x int
       fmt.Fscan(reader, &x)
       diff[x]--
   }
   // collect species present
   keys := make([]int, 0, len(diff))
   for x := range diff {
       keys = append(keys, x)
   }
   sort.Ints(keys)
   // check suffix sums
   suffix := 0
   // iterate from largest to smallest
   for i := len(keys) - 1; i >= 0; i-- {
       suffix += diff[keys[i]]
       if suffix > 0 {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
