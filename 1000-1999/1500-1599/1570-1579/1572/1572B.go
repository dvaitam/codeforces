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

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       // Check feasibility: total XOR must be 0 and at least one zero
       totalXor := 0
       hasZero := false
       for _, v := range a {
           totalXor ^= v
           if v == 0 {
               hasZero = true
           }
       }
       if totalXor != 0 || !hasZero {
           fmt.Fprintln(out, "NO")
           continue
       }
       // Prepare counts of ones in each triple
       cnt := make([]int, n-2)
       for i := 0; i < n-2; i++ {
           cnt[i] = a[i] + a[i+1] + a[i+2]
       }
       // BFS queue for triples with 1 or 2 ones
       queue := make([]int, 0, n)
       inq := make([]bool, n-2)
       for i := 0; i < n-2; i++ {
           if cnt[i] == 1 || cnt[i] == 2 {
               queue = append(queue, i)
               inq[i] = true
           }
       }
       ops := make([]int, 0, n)
       head := 0
       for head < len(queue) && len(ops) <= n {
           i := queue[head]
           head++
           if i < 0 || i >= n-2 {
               continue
           }
           c := cnt[i]
           if c != 1 && c != 2 {
               continue
           }
           // perform operation at i: set a[i..i+2] = parity
           p := c & 1
           ops = append(ops, i+1) // 1-based
           for j := i; j < i+3; j++ {
               a[j] = p
           }
           // update neighboring counts
           for k := i - 2; k <= i+2; k++ {
               if k >= 0 && k < n-2 {
                   cnt[k] = a[k] + a[k+1] + a[k+2]
                   if !inq[k] && (cnt[k] == 1 || cnt[k] == 2) {
                       queue = append(queue, k)
                       inq[k] = true
                   }
               }
           }
       }
       // check all zeros
       ok := true
       for _, v := range a {
           if v != 0 {
               ok = false
               break
           }
       }
       if !ok || len(ops) > n {
           fmt.Fprintln(out, "NO")
       } else {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, len(ops))
           for idx, v := range ops {
               if idx > 0 {
                   fmt.Fprint(out, ' ')
               }
               fmt.Fprint(out, v)
           }
           fmt.Fprintln(out)
       }
   }
}
