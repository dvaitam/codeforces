package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   maxV := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxV {
           maxV = a[i]
       }
   }
   // presence array
   present := make([]bool, maxV+1)
   for _, v := range a {
       present[v] = true
   }
   // last[i] = largest value <= i that is present
   last := make([]int, maxV+1)
   for i := 1; i <= maxV; i++ {
       if present[i] {
           last[i] = i
       } else {
           last[i] = last[i-1]
       }
   }
   ans := 0
   // for each possible divisor x
   for x := 1; x <= maxV; x++ {
       if !present[x] {
           continue
       }
       // check segments [j-x, j-1] for multiples j of x, starting at 2*x for a_i >= a_j
       for j := 2 * x; j <= maxV + x; j += x {
           lo := j - x
           if lo > maxV {
               break
           }
           hi := j - 1
           if hi > maxV {
               hi = maxV
           }
           y := last[hi]
           if y >= lo {
               // y % x
               r := y - lo
               if r > ans {
                   ans = r
               }
               // maximum possible mod is x-1
               if ans == x-1 {
                   break // no larger mod for this x
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, ans)
}
