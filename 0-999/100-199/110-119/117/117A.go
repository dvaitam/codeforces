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
   var n int
   var m int64
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // Period of full cycle: up and down
   T := 2 * (m - 1)
   for i := 0; i < n; i++ {
       var s, f, ti int64
       fmt.Fscan(in, &s, &f, &ti)
       // If already at destination
       if s == f {
           fmt.Fprintln(out, ti)
           continue
       }
       // Times (mod T) when elevator is at floor s going up (a) or down (b)
       a := s - 1
       b := (T - a) % T
       // Compute earliest boarding and arrival for a
       var ta int64
       if ti <= a {
           ta = a
       } else {
           ta = a + ((ti - a + T - 1) / T) * T
       }
       // Direction at ta: up if ta%T < m-1, else down
       dirUpA := (a < m-1)
       // Arrival time if boarding at ta
       var ansA int64
       if s < f {
           if dirUpA {
               ansA = ta + (f - s)
           } else {
               // go down to 1, then up to f
               ansA = ta + (s-1) + (f-1)
           }
       } else {
           // s > f
           if !dirUpA {
               ansA = ta + (s - f)
           } else {
               // go up to m, then down to f
               ansA = ta + (m-s) + (m-f)
           }
       }
       ans := ansA
       // If two distinct visits, compute for b
       if b != a {
           var tb int64
           if ti <= b {
               tb = b
           } else {
               tb = b + ((ti - b + T - 1) / T) * T
           }
           dirUpB := (b < m-1)
           var ansB int64
           if s < f {
               if dirUpB {
                   ansB = tb + (f - s)
               } else {
                   ansB = tb + (s-1) + (f-1)
               }
           } else {
               if !dirUpB {
                   ansB = tb + (s - f)
               } else {
                   ansB = tb + (m-s) + (m-f)
               }
           }
           if ansB < ans {
               ans = ansB
           }
       }
       fmt.Fprintln(out, ans)
   }
}
