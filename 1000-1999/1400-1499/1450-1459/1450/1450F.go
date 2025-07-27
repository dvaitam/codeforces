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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tt := 0; tt < t; tt++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // count breaks where adjacent tags equal
       breakCnt := 0
       for i := 0; i+1 < n; i++ {
           if a[i] == a[i+1] {
               breakCnt++
           }
       }
       // number of segments
       k := breakCnt + 1
       // count palindromic segments (same tag at ends)
       palCount := make(map[int]int)
       // traverse segments
       l := 0
       for i := 0; i < n; i++ {
           if i+1 < n && a[i] != a[i+1] {
               continue
           }
           // segment from l to i
           if a[l] == a[i] {
               palCount[a[l]]++
           }
           l = i + 1
       }
       // check feasibility
       maxPal := (k + 1) / 2
       ok := true
       for _, cnt := range palCount {
           if cnt > maxPal {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Fprintln(writer, -1)
       } else {
           // minimal cost is segments-1
           fmt.Fprintln(writer, k-1)
       }
   }
}
