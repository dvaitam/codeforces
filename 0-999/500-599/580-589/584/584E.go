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
   var n int
   fmt.Fscan(reader, &n)
   p := make([]int, n+1)
   pp := make([]int, n+1)
   s := make([]int, n+1)
   ps := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
       pp[p[i]] = i
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &s[i])
       ps[s[i]] = i
   }
   // map p elements to their positions in s
   for i := 1; i <= n; i++ {
       p[i] = ps[p[i]]
   }
   var mc int64
   for x := 1; x <= n; x++ {
       diff := pp[x] - ps[x]
       if diff < 0 {
           diff = -diff
       }
       mc += int64(diff)
   }
   // perform swaps
   swaps := make([][2]int, 0)
   for {
       flag := false
       for x := 1; x <= n; x++ {
           if p[x] >= x+1 {
               for y := p[x]; y > x; y-- {
                   if p[y] <= x {
                       p[x], p[y] = p[y], p[x]
                       swaps = append(swaps, [2]int{x, y})
                       flag = true
                       break
                   }
               }
           }
       }
       if !flag {
           break
       }
   }
   // output results
   fmt.Fprintln(writer, mc/2)
   fmt.Fprintln(writer, len(swaps))
   for _, sw := range swaps {
       fmt.Fprintln(writer, sw[0], sw[1])
   }
}
