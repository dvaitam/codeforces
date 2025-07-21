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
   // days in months for 2013 (non-leap year)
   monthDays := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
   // offset for negative indices
   const offset = 500
   // enough days: cover from ~-100 to 365 -> total ~900
   const maxDays = 1000
   load := make([]int, maxDays)
   // helper to compute day-of-year 0-based
   prefix := make([]int, 13)
   for i := 1; i <= 12; i++ {
       prefix[i] = prefix[i-1] + monthDays[i-1]
   }
   for i := 0; i < n; i++ {
       var mi, di, pi, ti int
       fmt.Fscan(reader, &mi, &di, &pi, &ti)
       // day of year (0-based) for Olympiad day
       doy := prefix[mi-1] + (di - 1)
       // preparation from doy-ti to doy-1 inclusive
       start := doy - ti + offset
       end := doy - 1 + offset
       if start < 0 {
           start = 0
       }
       if end >= maxDays {
           end = maxDays - 1
       }
       for d := start; d <= end; d++ {
           load[d] += pi
       }
   }
   // find maximum load
   ans := 0
   for _, v := range load {
       if v > ans {
           ans = v
       }
   }
   fmt.Println(ans)
}
