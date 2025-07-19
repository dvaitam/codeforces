package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   ones, twos := 0, 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x == 1 {
           ones++
       } else {
           twos++
       }
   }
   var res []int
   if ones == 0 {
       // all twos
       for i := 0; i < twos; i++ {
           res = append(res, 2)
       }
       printRes(res, writer)
       return
   }
   if twos == 0 {
       // all ones
       for i := 0; i < ones; i++ {
           res = append(res, 1)
       }
       printRes(res, writer)
       return
   }
   // both present
   // start with 2, then 1
   res = append(res, 2)
   twos--
   if twos == 0 {
       // no more twos, then all ones
       for i := 0; i < ones; i++ {
           res = append(res, 1)
       }
       printRes(res, writer)
       return
   }
   res = append(res, 1)
   ones--
   if ones == 0 {
       // no more ones, then all twos
       for i := 0; i < twos; i++ {
           res = append(res, 2)
       }
       printRes(res, writer)
       return
   }
   // remaining: all twos then all ones
   for i := 0; i < twos; i++ {
       res = append(res, 2)
   }
   for i := 0; i < ones; i++ {
       res = append(res, 1)
   }
   printRes(res, writer)
}

// printRes writes the result slice to writer with spaces
func printRes(res []int, w *bufio.Writer) {
   for i, x := range res {
       if i > 0 {
           w.WriteByte(' ')
       }
       w.WriteString(strconv.Itoa(x))
   }
}
