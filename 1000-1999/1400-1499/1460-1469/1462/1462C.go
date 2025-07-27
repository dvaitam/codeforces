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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x > 45 {
           fmt.Fprintln(writer, -1)
           continue
       }
       var digits []int
       for d := 9; d >= 1 && x > 0; d-- {
           if x >= d {
               digits = append(digits, d)
               x -= d
           }
       }
       sort.Ints(digits)
       for _, d := range digits {
           fmt.Fprint(writer, d)
       }
       fmt.Fprintln(writer)
   }
}
