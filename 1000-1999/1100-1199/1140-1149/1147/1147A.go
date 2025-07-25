package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   fmt.Fscan(in, &n, &k)
   first := make([]int, n+2)
   last := make([]int, n+2)
   // initialize first occurrences to k+1, last to 0
   for i := 1; i <= n; i++ {
       first[i] = k + 1
       last[i] = 0
   }
   // read questions and record occurrences
   for i := 1; i <= k; i++ {
       var x int
       fmt.Fscan(in, &x)
       if first[x] > i {
           first[x] = i
       }
       last[x] = i
   }
   // count valid scenarios
   ans := 0
   for a := 1; a <= n; a++ {
       for d := -1; d <= 1; d++ {
           b := a + d
           if b < 1 || b > n {
               continue
           }
           // valid if last occurrence of b is before first occurrence of a
           if last[b] < first[a] {
               ans++
           }
       }
   }
   fmt.Fprintln(out, ans)
}
