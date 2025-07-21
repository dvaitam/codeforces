package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, p int
   if _, err := fmt.Fscan(in, &n, &m, &p); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   // Build target counts
   target := make(map[int]int)
   for _, v := range b {
       target[v]++
   }
   ukeys := len(target)

   var res []int
   // Process each modulo class
   for r := 0; r < p; r++ {
       if r >= n {
           break
       }
       // group length
       K := (n-1-r)/p + 1
       if K < m {
           continue
       }
       // build group values
       vals := make([]int, K)
       idx := r
       for i := 0; i < K; i++ {
           vals[i] = a[idx]
           idx += p
       }
       // initialize count map for this group
       cnt := make(map[int]int, len(target))
       for k, v := range target {
           cnt[k] = v
       }
       missing := ukeys
       // helper for add and remove
       add := func(x int) {
           c := cnt[x] // zero if not exist
           c--
           cnt[x] = c
           if c == 0 {
               missing--
           } else if c == -1 {
               missing++
           }
       }
       remove := func(x int) {
           c := cnt[x]
           c++
           cnt[x] = c
           if c == 0 {
               missing--
           } else if c == 1 {
               missing++
           }
       }
       // initial window
       for i := 0; i < m; i++ {
           add(vals[i])
       }
       if missing == 0 {
           // starting at j=0
           q := r + 1
           res = append(res, q)
       }
       // slide
       for i := m; i < K; i++ {
           // window moves: remove vals[i-m], add vals[i]
           remove(vals[i-m])
           add(vals[i])
           if missing == 0 {
               // start index j = i-m+1
               q := r + (i-m+1)*p + 1
               res = append(res, q)
           }
       }
   }
   sort.Ints(res)
   // output count and positions
   fmt.Fprintln(out, len(res))
   // print positions or empty line
   for i, q := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, q)
   }
   fmt.Fprintln(out)
}
