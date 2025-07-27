package main

import (
   "bufio"
   "fmt"
   "os"
)

type supply struct {
   idx int
   rem int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // total positive supply
       var totalSupply int64
       for i := 0; i < n; i++ {
           if a[i] > 0 {
               totalSupply += a[i]
           }
       }
       // queue of supplies
       supplies := make([]supply, 0, n)
       head := 0
       var freeFlow int64
       // scan positions
       for j := 0; j < n; j++ {
           if a[j] > 0 {
               supplies = append(supplies, supply{idx: j, rem: a[j]})
           } else if a[j] < 0 {
               need := -a[j]
               // match with earlier supplies for free
               for need > 0 && head < len(supplies) && supplies[head].idx < j {
                   cur := &supplies[head]
                   var d int64
                   if cur.rem <= need {
                       d = cur.rem
                   } else {
                       d = need
                   }
                   freeFlow += d
                   cur.rem -= d
                   need -= d
                   if cur.rem == 0 {
                       head++
                   }
               }
           }
       }
       // coins = totalSupply - freeFlow
       coins := totalSupply - freeFlow
       fmt.Fprintln(writer, coins)
   }
}
