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
   for i := 0; i < T; i++ {
       var A, B, C, D int
       fmt.Fscan(reader, &A, &B, &C, &D)
       // sign: negative if odd number of negatives (A+B odd)
       neg := (A+B)%2 == 1
       // has factors abs>=1 if any from box1 or box4
       hasGe := (A + D) > 0
       // has factors abs<1 if any from box2 or box3
       hasLt := (B + C) > 0

       if neg {
           // final negative: possible boxes 1 and 2
           if hasGe {
               writer.WriteString("Ya")
           } else {
               writer.WriteString("Tidak")
           }
           writer.WriteByte(' ')
           if hasLt {
               writer.WriteString("Ya")
           } else {
               writer.WriteString("Tidak")
           }
           writer.WriteString(" Tidak Tidak")
       } else {
           // final positive: possible boxes 3 and 4
           writer.WriteString("Tidak Tidak ")
           if hasLt {
               writer.WriteString("Ya")
           } else {
               writer.WriteString("Tidak")
           }
           writer.WriteByte(' ')
           if hasGe {
               writer.WriteString("Ya")
           } else {
               writer.WriteString("Tidak")
           }
       }
       writer.WriteByte('\n')
   }
}
