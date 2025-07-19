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

   var n, alice int
   if _, err := fmt.Fscan(reader, &n, &alice); err != nil {
       return
   }
   const MAXC = 1000001
   cars := make([]int, MAXC)
   threshold := 0
   for i := 0; i < n; i++ {
       var tmp int
       fmt.Fscan(reader, &tmp)
       if tmp == alice {
           threshold++
           continue
       }
       if cars[tmp] < threshold {
           cars[tmp] = -1
       } else {
           cars[tmp]++
       }
   }
   for i := 1; i < MAXC; i++ {
       if i == alice {
           continue
       }
       if cars[i] < threshold {
           continue
       }
       fmt.Fprintln(writer, i)
       return
   }
   fmt.Fprintln(writer, -1)
}
