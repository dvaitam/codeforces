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

   var n, k, x int
   fmt.Fscan(reader, &n, &k, &x)
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n/k > x {
       fmt.Fprint(writer, -1)
       return
   }
   // dp arrays: g for previous, f for current
   g := make([]int64, n+2)
   f := make([]int64, n+2)
   // initial dp for l=1
   for i := 1; i <= k && i <= n; i++ {
       g[i] = int64(a[i])
   }
   // temporary deque for sliding max
   dui := make([]int, n+2)
   // build dp for l = 2..x
   for l := 2; l <= x; l++ {
       // reset current dp
       for i := 1; i <= n; i++ {
           f[i] = 0
       }
       t1, t2 := 0, -1
       // positions where dp values can be non-zero: i in [l, min(n, l*k)]
       upper := l * k
       if upper > n {
           upper = n
       }
       for i := l; i <= upper; i++ {
           // push index i-1 into deque, maintain descending g[]
           for t1 <= t2 && g[i-1] >= g[dui[t2]] {
               t2--
           }
           t2++
           dui[t2] = i - 1
           // pop front if out of window size k
           for t1 <= t2 && dui[t1] < i-k {
               t1++
           }
           if t1 <= t2 {
               f[i] = g[dui[t1]] + int64(a[i])
           }
       }
       // swap f and g
       for i := 1; i <= n; i++ {
           g[i] = f[i]
       }
   }
   // answer is max over last k positions
   var ans int64
   start := max(1, n-k+1)
   for i := start; i <= n; i++ {
       if g[i] > ans {
           ans = g[i]
       }
   }
   fmt.Fprint(writer, ans)
}
