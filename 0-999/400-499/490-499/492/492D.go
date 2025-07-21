package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var x, y int64
   fmt.Fscan(reader, &n, &x, &y)
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       // binary search minimal p such that floor(p/x) + floor(p/y) >= a
       low, high := int64(0), a*min(x, y)
       for low < high {
           mid := (low + high) / 2
           if mid/x + mid/y < a {
               low = mid + 1
           } else {
               high = mid
           }
       }
       p := low
       // determine who hits at time p (units: p/(x*y) seconds)
       vanyaHit := p%y == 0
       vovaHit := p%x == 0
       var res string
       if vanyaHit && vovaHit {
           res = "Both"
       } else if vanyaHit {
           res = "Vanya"
       } else {
           res = "Vova"
       }
       fmt.Fprintln(writer, res)
   }
}
