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
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, w int
       fmt.Fscan(in, &n, &w)
       diffs := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &diffs[i])
       }
       pref, minPref, maxPref := 0, 0, 0
       for _, d := range diffs {
           pref += d
           if pref < minPref {
               minPref = pref
           }
           if pref > maxPref {
               maxPref = pref
           }
       }
       // x must satisfy: -minPref <= x <= w - maxPref
       res := w - maxPref + minPref + 1
       if res < 0 {
           res = 0
       }
       fmt.Fprintln(out, res)
   }
}
