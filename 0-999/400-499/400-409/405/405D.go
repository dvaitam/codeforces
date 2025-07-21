package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const s = 1000000
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   isX := make([]bool, s+1)
   var sumX int64
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= s && !isX[x] {
           isX[x] = true
           sumX += int64(x)
       }
   }
   // sumL = sum_{x in X} (x-1)
   sumL := sumX - int64(n)
   // special case: target sum zero -> pick y = s
   if sumL == 0 {
       // s cannot be in X because sumL==0 only if X contains only 1
       fmt.Fprintln(writer, 1)
       fmt.Fprintln(writer, s)
       return
   }
   rem := sumL
   var Y []int
   // greedy: pick y with largest weight w = s-y
   for y := 1; y <= s; y++ {
       if !isX[y] {
           w := int64(s - y)
           if w <= rem {
               Y = append(Y, y)
               rem -= w
               if rem == 0 {
                   break
               }
           }
       }
   }
   // output
   m := len(Y)
   fmt.Fprintln(writer, m)
   for i, y := range Y {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, y)
   }
   fmt.Fprintln(writer)
}
