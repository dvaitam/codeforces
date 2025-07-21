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

   var a [5]int
   var s string
   // Read calorie costs for strips 1 to 4
   fmt.Fscan(reader, &a[1], &a[2], &a[3], &a[4])
   // Read the sequence of appearing squares
   fmt.Fscan(reader, &s)

   total := 0
   for _, ch := range s {
       // ch is '1'...'4'; subtract '0' to get 1..4 index
       total += a[int(ch-'0')]
   }
   fmt.Fprint(writer, total)
}
