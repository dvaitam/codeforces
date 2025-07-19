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

   var n int
   var d, x int64
   if _, err := fmt.Fscan(in, &n, &d); err != nil {
       return
   }
   // Prepare slices
   maxn := make([]int64, n)
   minn := make([]int64, n)
   maxnn := make([]int, n)
   minnn := make([]int, n)
   ans := make([]int64, n)
   ans1 := make([]int, n)

   // Read first element
   fmt.Fscan(in, &x)
   tail := 0
   maxn[0], minn[0] = x, x
   maxnn[0], minnn[0] = 1, 1

   // Process remaining
   for i := 2; i <= n; i++ {
       fmt.Fscan(in, &x)
       if minn[tail]+d <= x || maxn[tail]-d >= x {
           tail++
           maxn[tail], minn[tail] = x, x
           maxnn[tail], minnn[tail] = i, i
       } else if tail > 0 {
           // try updating previous segment bounds
           if minn[tail-1]+d <= x || maxn[tail-1]-d >= x {
               if maxn[tail] < x {
                   maxn[tail] = x
                   maxnn[tail] = i
               }
               if minn[tail] > x {
                   minn[tail] = x
                   minnn[tail] = i
               }
           }
       } else {
           // tail == 0, update current
           if maxn[tail] < x {
               maxn[tail] = x
               maxnn[tail] = i
           }
           if minn[tail] > x {
               minn[tail] = x
               minnn[tail] = i
           }
       }
   }

   // Number of segments is tail+1
   fmt.Fprintln(out, tail+1)
   // Reconstruct answer indices
   ans[tail] = minn[tail]
   ans1[tail] = minnn[tail]
   for i := tail - 1; i >= 0; i-- {
       if minn[i]+d <= ans[i+1] {
           ans[i] = minn[i]
           ans1[i] = minnn[i]
       } else {
           ans[i] = maxn[i]
           ans1[i] = maxnn[i]
       }
   }
   // Output indices
   for i := 0; i <= tail; i++ {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, ans1[i])
   }
   fmt.Fprintln(out)
}
