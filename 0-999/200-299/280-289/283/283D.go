package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // bad[i] == true if pair (a[i], a[i+1]) is not cool
   bad := make([]bool, n)
   for i := 0; i < n-1; i++ {
       x := a[i]
       y := a[i+1]
       // sum of y consecutive ints: y*k + y*(y-1)/2 == x
       // check (x - y*(y-1)/2) % y == 0
       var d int64
       if y%2 == 1 {
           d = 0
       } else {
           d = y / 2
       }
       if (x - d) % y != 0 {
           bad[i] = true
       }
   }

   // dp0[i]: max keeps up to i with i not kept
   // dp1[i]: max keeps up to i with i kept
   dp0 := make([]int, n+1)
   dp1 := make([]int, n+1)
   dp0[1] = 0
   dp1[1] = 1
   for i := 2; i <= n; i++ {
       if bad[i-2] {
           dp0[i] = max(dp0[i-1], dp1[i-1])
           dp1[i] = 1 + dp0[i-1]
       } else {
           mx := max(dp0[i-1], dp1[i-1])
           dp0[i] = mx
           dp1[i] = 1 + mx
       }
   }

   best := max(dp0[n], dp1[n])
   fmt.Fprintln(writer, n-best)
}
