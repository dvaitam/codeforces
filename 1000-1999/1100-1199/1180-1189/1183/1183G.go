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

   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n int
       fmt.Fscan(in, &n)
       // counts per type: map type -> [total, good]
       cnt := make(map[int][2]int, n)
       for i := 0; i < n; i++ {
           var a, f int
           fmt.Fscan(in, &a, &f)
           v := cnt[a]
           v[0]++
           if f == 1 {
               v[1]++
           }
           cnt[a] = v
       }
       // build slice of (total, good)
       types := make([][2]int, 0, len(cnt))
       for _, v := range cnt {
           types = append(types, v)
       }
       // sort by total descending
       sort.Slice(types, func(i, j int) bool {
           return types[i][0] > types[j][0]
       })
       total := 0
       good := 0
       last := int(1e9)
       for _, v := range types {
           c := v[0]
           u := v[1]
           // pick x candies of this type
           if last <= 0 {
               break
           }
           x := c
           if last-1 < x {
               x = last - 1
           }
           if x <= 0 {
               break
           }
           total += x
           if u < x {
               good += u
           } else {
               good += x
           }
           last = x
       }
       fmt.Fprintln(out, total, good)
   }
}
