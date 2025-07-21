package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, d int
   if _, err := fmt.Fscan(reader, &n, &m, &d); err != nil {
       return
   }
   total := n * m
   vals := make([]int, 0, total)
   var modBase int
   first := true
   for i := 0; i < total; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if first {
           modBase = a % d
           first = false
       }
       if a % d != modBase {
           fmt.Fprintln(writer, -1)
           return
       }
       vals = append(vals, a / d)
   }
   sort.Ints(vals)
   // median minimizes absolute deviations
   median := vals[total/2]
   var moves int64
   for _, v := range vals {
       if v > median {
           moves += int64(v - median)
       } else {
           moves += int64(median - v)
       }
   }
   fmt.Fprintln(writer, moves)
}
