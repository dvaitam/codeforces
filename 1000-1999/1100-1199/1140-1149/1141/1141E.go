package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var h int64
   var n int
   if _, err := fmt.Fscan(reader, &h, &n); err != nil {
       return
   }
   v := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }

   // sum stores total damage per cycle (positive value)
   var sum int64 = 0
   // mx stores the worst drop relative to start
   var mx int64 = 0
   // H tracks current health in first pass
   H := h
   // ans stores first day when health <= 0 in first cycle
   ans := int64(-1)
   for i := 0; i < n; i++ {
       sum -= v[i]
       H += v[i]
       if H <= 0 && ans == -1 {
           ans = int64(i + 1)
       }
       if sum > mx {
           mx = sum
       }
   }

   if ans != -1 {
       fmt.Fprint(writer, ans)
       return
   }
   // if no net damage per cycle, impossible to die
   if sum <= 0 {
       fmt.Fprint(writer, -1)
       return
   }

   // compute number of full cycles before approaching death
   // full = max number of full cycles such that h - full*sum > mx
   full := (h - mx) / sum
   h -= full * sum
   cnt := full * int64(n)

   // simulate remaining days
   for i := 0; ; i++ {
       h += v[i% n]
       cnt++
       if h <= 0 {
           fmt.Fprint(writer, cnt)
           return
       }
   }
}
