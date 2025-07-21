package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(in, &n, &x); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // find minimal a
   m := a[0]
   for i := 1; i < n; i++ {
       if a[i] < m {
           m = a[i]
       }
   }
   // candidate k values
   ks := []int64{m}
   if m > 0 {
       ks = append(ks, m-1)
   }
   // try candidates
   for _, k := range ks {
       // compute d = a - k
       d := make([]int64, n)
       ok := true
       for i := 0; i < n; i++ {
           d[i] = a[i] - k
           if d[i] < 0 {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       // count r: consecutive d>0 ending at x-1 backwards
       rCount := 0
       pos := x - 1
       for rCount < n && d[pos] > 0 {
           rCount++
           pos = (pos - 1 + n) % n
       }
       // r mod n (if full cycle, treat as zero)
       r := rCount % n
       // ensure there is at least one ball taken
       if k == 0 && r == 0 {
           continue
       }
       // compute start i (0-based)
       i0 := ( (x-1 - r) % n + n ) % n
       // mark segment positions
       seg := make([]bool, n)
       for t := 1; t <= r; t++ {
           j := (i0 + t) % n
           seg[j] = true
       }
       // build initial array
       init := make([]int64, n)
       base := k*int64(n) + int64(r)
       init[i0] = base
       valid := base > 0
       for j := 0; j < n && valid; j++ {
           if j == i0 {
               continue
           }
           if seg[j] {
               init[j] = a[j] - k - 1
           } else {
               init[j] = a[j] - k
           }
           if init[j] < 0 {
               valid = false
           }
       }
       if !valid {
           continue
       }
       // output result in 1-based order
       out := bufio.NewWriter(os.Stdout)
       defer out.Flush()
       for i, v := range init {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       return
   }
}
