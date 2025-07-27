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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for t := 0; t < T; t++ {
       var n int
       fmt.Fscan(reader, &n)
       var rs, bs string
       fmt.Fscan(reader, &rs)
       fmt.Fscan(reader, &bs)
       sumR, sumB := 0, 0
       for i := 0; i < n; i++ {
           sumR += int(rs[i] - '0')
           sumB += int(bs[i] - '0')
       }
       switch {
       case sumR > sumB:
           fmt.Fprintln(writer, "RED")
       case sumR < sumB:
           fmt.Fprintln(writer, "BLUE")
       default:
           fmt.Fprintln(writer, "EQUAL")
       }
   }
}
