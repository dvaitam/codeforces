package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   m := (1 << uint(k)) - 1
   a := make([]int, n+1)
   // prefix xor with minimization
   for i := 1; i <= n; i++ {
       var bi int
       fmt.Fscan(reader, &bi)
       xi := a[i-1] ^ bi
       if t := xi ^ m; t < xi {
           a[i] = t
       } else {
           a[i] = xi
       }
   }
   sort.Ints(a)
   ln := int64(n + 1)
   ans := ln * (ln - 1) / 2
   for i := 0; i < len(a); {
       j := i + 1
       for j < len(a) && a[j] == a[i] {
           j++
       }
       cnt := int64(j - i)
       c1 := cnt / 2
       c2 := cnt - c1
       ans -= c1 * (c1 - 1) / 2
       ans -= c2 * (c2 - 1) / 2
       i = j
   }
   fmt.Println(ans)
}
