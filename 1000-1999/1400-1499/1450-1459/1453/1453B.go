package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // total operations without change
       var total int64
       for i := 1; i < n; i++ {
           total += abs(a[i] - a[i-1])
       }
       // compute maximum gain by changing one element
       var maxGain int64
       if n >= 2 {
           // endpoints
           gainFirst := abs(a[1] - a[0])
           gainLast := abs(a[n-1] - a[n-2])
           if gainFirst > maxGain {
               maxGain = gainFirst
           }
           if gainLast > maxGain {
               maxGain = gainLast
           }
       }
       // middle elements
       for i := 1; i < n-1; i++ {
           left := abs(a[i] - a[i-1])
           right := abs(a[i+1] - a[i])
           skip := abs(a[i+1] - a[i-1])
           gain := left + right - skip
           if gain > maxGain {
               maxGain = gain
           }
       }
       // result is total operations minus best gain
       result := total - maxGain
       fmt.Fprintln(writer, result)
   }
}
