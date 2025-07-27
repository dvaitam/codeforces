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
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       b := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       minA, minB := a[0], b[0]
       for i := 1; i < n; i++ {
           if a[i] < minA {
               minA = a[i]
           }
           if b[i] < minB {
               minB = b[i]
           }
       }
       var moves int64
       for i := 0; i < n; i++ {
           da := a[i] - minA
           db := b[i] - minB
           if da > db {
               moves += da
           } else {
               moves += db
           }
       }
       fmt.Fprintln(writer, moves)
   }
}
