package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   type fig struct { cnt, cost int64 }
   figs := make([]fig, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &figs[i].cnt, &figs[i].cost)
   }
   var t int
   fmt.Fscan(reader, &t)
   p := make([]int64, t)
   for i := 0; i < t; i++ {
       fmt.Fscan(reader, &p[i])
   }
   // sort by cost ascending
   sort.Slice(figs, func(i, j int) bool {
       return figs[i].cost < figs[j].cost
   })
   // total figures
   var total int64
   for _, f := range figs {
       total += f.cnt
   }
   // build segments of kills for each factor
   seg := make([]int64, 0, t+1)
   var prev int64
   for i := 0; i < t; i++ {
       if prev >= total {
           break
       }
       up := p[i]
       if up > total {
           up = total
       }
       cur := up - prev
       if cur > 0 {
           seg = append(seg, cur)
           prev += cur
       }
   }
   if prev < total {
       seg = append(seg, total-prev)
   }
   // assign figures to segments
   var ans int64
   si, rem := 0, int64(0)
   if len(seg) > 0 {
       rem = seg[0]
   }
   for _, f := range figs {
       cnt := f.cnt
       for cnt > 0 && si < len(seg) {
           take := cnt
           if take > rem {
               take = rem
           }
           ans += take * f.cost * int64(si+1)
           cnt -= take
           rem -= take
           if rem == 0 {
               si++
               if si < len(seg) {
                   rem = seg[si]
               }
           }
       }
   }
   fmt.Fprintln(writer, ans)
