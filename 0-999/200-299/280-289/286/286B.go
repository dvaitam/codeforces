package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // initialize permutation
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i + 1
   }
   // apply transformations f(p, k) for k = 2..n
   for k := 2; k <= n; k++ {
       for s := 0; s < n; s += k {
           L := min(k, n-s)
           if L <= 1 {
               continue
           }
           // rotate left by one on p[s : s+L]
           tmp := p[s]
           // shift
           copy(p[s:s+L-1], p[s+1:s+L])
           p[s+L-1] = tmp
       }
   }
   // print result
   for i, v := range p {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
