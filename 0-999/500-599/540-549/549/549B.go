package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // read contact matrix
   mat := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       mat[i] = []byte(s)
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // cnt[j]: number of messages employee j receives from current attendees
   cnt := make([]int, n)
   alive := make([]bool, n)
   for i := 0; i < n; i++ {
       alive[i] = true
       for j := 0; j < n; j++ {
           if mat[i][j] == '1' {
               cnt[j]++
           }
       }
   }
   // queue of removals
   q := make([]int, 0, n)
   inQ := make([]bool, n)
   for i := 0; i < n; i++ {
       if cnt[i] == a[i] {
           q = append(q, i)
           inQ[i] = true
       }
   }
   // process removals
   for head := 0; head < len(q); head++ {
       i := q[head]
       if !alive[i] || cnt[i] != a[i] {
           continue
       }
       // remove i from attendees
       alive[i] = false
       // update counts
       for j := 0; j < n; j++ {
           if mat[i][j] == '1' {
               cnt[j]--
               if alive[j] && !inQ[j] && cnt[j] == a[j] {
                   q = append(q, j)
                   inQ[j] = true
               }
           }
       }
   }
   // collect result
   res := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if alive[i] {
           res = append(res, i+1)
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   // output number of attendees
   fmt.Fprintln(w, len(res))
   // output attendee indices (empty line if none)
   if len(res) == 0 {
       fmt.Fprintln(w)
   } else {
       for i, v := range res {
           if i > 0 {
               w.WriteByte(' ')
           }
           fmt.Fprint(w, v)
       }
       w.WriteByte('\n')
   }
}
