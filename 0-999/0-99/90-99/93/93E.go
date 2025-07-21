package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var n uint64
var k int
var a []uint64
var res int64

// dfs applies inclusion-exclusion to count numbers divisible by any subset of a
func dfs(idx int, prod uint64, depth int) {
   for i := idx; i < k; i++ {
       ai := a[i]
       if prod > n/ai {
           continue
       }
       newProd := prod * ai
       cnt := int64(n / newProd)
       if depth%2 == 0 {
           res += cnt
       } else {
           res -= cnt
       }
       dfs(i+1, newProd, depth+1)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a = make([]uint64, k)
   for i := 0; i < k; i++ {
       if _, err := fmt.Fscan(reader, &a[i]); err != nil {
           return
       }
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // if any indicator is 1, all numbers are divisible
   if k > 0 && a[0] == 1 {
       fmt.Fprint(writer, 0)
       return
   }
   dfs(0, 1, 0)
   damage := int64(n) - res
   fmt.Fprint(writer, damage)
}
