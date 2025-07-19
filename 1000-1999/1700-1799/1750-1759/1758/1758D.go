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

   var tc int
   fmt.Fscan(in, &tc)
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(in, &n)
       solve(n, out)
   }
}

func solve(n int, out *bufio.Writer) {
   nn := int64(n)
   for sq := nn + 100; ; sq++ {
       mx := sq + 1
       ans := make([]int64, n)
       var sum int64
       for i := 0; i < n; i++ {
           if i == n-1 {
               ans[i] = mx
           } else {
               ans[i] = int64(i + 1)
           }
           sum += ans[i]
       }
       rem := sq*sq - sum
       inc := rem / nn
       mod := rem % nn
       if mod > nn-2 {
           continue
       }
       for i := n - 1; i >= 0; i-- {
           ans[i] += inc
           if i > 0 && i+1 < n && mod > 0 {
               mod--
               ans[i]++
           }
       }
       for i, v := range ans {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
       return
   }
}
