package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   // number of sequences equals number of end markers (-1)
   n := 0
   for _, v := range b {
       if v == -1 {
           n++
       }
   }
   // prepare sequences
   seqs := make([][]int, n)
   // active sequence indices
   active := make([]int, n)
   for i := 0; i < n; i++ {
       active[i] = i
   }
   idx := 0
   // reconstruct by rounds
   for len(active) > 0 && idx < m {
       next := make([]int, 0, len(active))
       for _, si := range active {
           if idx >= m {
               break
           }
           v := b[idx]
           idx++
           if v != -1 {
               seqs[si] = append(seqs[si], v)
               next = append(next, si)
           }
           // skip ended sequences (v == -1)
       }
       active = next
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, n)
   for i := 0; i < n; i++ {
       fmt.Fprint(out, len(seqs[i]))
       for _, v := range seqs[i] {
           fmt.Fprint(out, " ", v)
       }
       fmt.Fprintln(out)
   }
}
