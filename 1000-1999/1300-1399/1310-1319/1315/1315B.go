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
   for ; t > 0; t-- {
       var a, b int64
       var p int64
       var s string
       fmt.Fscan(reader, &a, &b, &p)
       fmt.Fscan(reader, &s)
       n := len(s)
       var costSum int64
       // last seen segment type
       var last byte = 0
       ans := 1
       // traverse edges from n-1 downto 1 (0-based: s[0..n-2])
       for i := n - 2; i >= 0; i-- {
           if s[i] != last {
               // new segment
               if s[i] == 'A' {
                   costSum += a
               } else {
                   costSum += b
               }
               last = s[i]
           }
           if costSum > p {
               // need to start after this edge
               ans = i + 2 // convert to 1-based crossroad
               break
           }
       }
       // if total cost within p, can start at 1
       if costSum <= p {
           ans = 1
       }
       fmt.Fprintln(writer, ans)
   }
}
