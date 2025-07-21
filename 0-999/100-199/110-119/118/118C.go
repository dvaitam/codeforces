package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   bytes := []byte(s)
   bestCost := math.MaxInt64
   bestResult := ""
   // try target digit d
   for td := byte('0'); td <= byte('9'); td++ {
       // count existing
       cnt := 0
       for i := 0; i < n; i++ {
           if bytes[i] == td {
               cnt++
           }
       }
       need := k - cnt
       cost := 0
       // copy original
       b := make([]byte, n)
       copy(b, bytes)
       if need > 0 {
           // first decrease larger digits to td
           for dlt := 1; dlt <= 9 && need > 0; dlt++ {
               for i := 0; i < n && need > 0; i++ {
                   if b[i] > td && int(b[i]-td) == dlt {
                       cost += dlt
                       b[i] = td
                       need--
                   }
               }
               // then increase smaller digits to td, from end for lexicographic minimality
               for i := n - 1; i >= 0 && need > 0; i-- {
                   if b[i] < td && int(td-b[i]) == dlt {
                       cost += dlt
                       b[i] = td
                       need--
                   }
               }
           }
       }
       result := string(b)
       // update best
       if cost < bestCost || (cost == bestCost && result < bestResult) {
           bestCost = cost
           bestResult = result
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, bestCost)
   fmt.Fprintln(writer, bestResult)
}
