package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var m int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int, 0, 2*n)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       xs = append(xs, a)
   }
   for i := 0; i < n; i++ {
       var b int
       fmt.Fscan(reader, &b)
       xs = append(xs, b)
   }
   fmt.Fscan(reader, &m)

   sort.Ints(xs)
   total := 2 * n
   // precompute factorials mod m
   fact := make([]int, total+1)
   fact[0] = 1 % m
   for i := 1; i <= total; i++ {
       fact[i] = int((int64(fact[i-1]) * int64(i)) % int64(m))
   }

   ans := 1 % m
   for i := 0; i < total; {
       j := i + 1
       for j < total && xs[j] == xs[i] {
           j++
       }
       cnt := j - i
       ans = int((int64(ans) * int64(fact[cnt])) % int64(m))
       i = j
   }
   fmt.Fprint(writer, ans)
}
