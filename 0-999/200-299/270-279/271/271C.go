package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   if n <= 3 || 3*k > n {
       fmt.Fprintln(writer, -1)
       return
   }
   ans := make([]int, n+5)
   for i := range ans {
       ans[i] = -1
   }
   i := 1
   heso := 4
   var lastval, minNotAns int
   for id := 1; id <= k-1; id++ {
       lastval = i + 3
       if heso == 4 {
           minNotAns = i + 2
           ans[i] = id
           ans[i+1] = id
           ans[i+3] = id
           heso = 2
           i += 2
       } else {
           minNotAns = i + 4
           ans[i] = id
           ans[i+2] = id
           ans[i+3] = id
           heso = 4
           i += 4
       }
   }
   if lastval < minNotAns {
       if minNotAns+2 == n {
           ans[minNotAns] = k
           ans[minNotAns+2] = k
           ans[minNotAns+1] = k - 1
           ans[minNotAns-1] = k
       } else {
           for j := minNotAns; j <= n; j++ {
               ans[j] = k
           }
           if minNotAns+1 <= n {
               ans[minNotAns+1] = 1
           }
       }
   } else {
       for j := minNotAns; j <= n; j++ {
           if ans[j] == -1 {
               ans[j] = k
           }
       }
   }
   for j := 1; j <= n; j++ {
       if j > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[j])
   }
   writer.WriteByte('\n')
}
