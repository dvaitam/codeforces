package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   br := bufio.NewReader(os.Stdin)
   bw := bufio.NewWriter(os.Stdout)
   defer bw.Flush()

   var t int
   if _, err := fmt.Fscan(br, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       solve(br, bw)
   }
}

func solve(br *bufio.Reader, bw *bufio.Writer) {
   var n, l, r, s int
   fmt.Fscan(br, &n, &l, &r, &s)
   l--
   r--
   k := r - l + 1
   // try minimal consecutive sequences
   found := false
   for first := 1; first + k - 1 <= n; first++ {
       // sum of arithmetic sequence from first of length k
       sum := k*first + (k*(k-1))/2
       if sum > s {
           break
       }
       extra := s - sum
       if extra <= k {
           // build answer
           ans := make([]int, n)
           used := make([]bool, n+1)
           // index to start adding extra +1
           needAdd := r - extra + 1
           ok := true
           for i := l; i <= r; i++ {
               val := first + (i - l)
               if i >= needAdd {
                   val++
               }
               if val > n {
                   ok = false
                   break
               }
               ans[i] = val
               used[val] = true
           }
           if !ok {
               continue
           }
           // fill remaining positions with smallest unused
           cur := 1
           for i := 0; i < n; i++ {
               if i >= l && i <= r {
                   continue
               }
               for cur <= n && used[cur] {
                   cur++
               }
               ans[i] = cur
               used[cur] = true
           }
           // output
           for i := 0; i < n; i++ {
               if i > 0 {
                   bw.WriteByte(' ')
               }
               fmt.Fprint(bw, ans[i])
           }
           bw.WriteByte('\n')
           found = true
           break
       }
   }
   if !found {
       fmt.Fprintln(bw, -1)
   }
}
