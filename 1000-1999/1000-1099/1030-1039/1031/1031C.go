package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, b int
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   total := a + b
   // find maximum k such that k*(k+1)/2 <= total
   // approximate via sqrt, then adjust
   k := int(math.Sqrt(2 * float64(total)))
   // adjust down
   for k*(k+1)/2 > total {
       k--
   }
   // adjust up
   for (k+1)*(k+2)/2 <= total {
       k++
   }

   day1 := make([]int, 0, k)
   day2 := make([]int, 0, k)
   rem := a
   for i := k; i >= 1; i-- {
       if rem >= i {
           day1 = append(day1, i)
           rem -= i
       } else {
           day2 = append(day2, i)
       }
   }
   // print day1 and day2
   sort.Ints(day1)
   sort.Ints(day2)
   fmt.Fprintln(writer, len(day1))
   if len(day1) > 0 {
       for i, v := range day1 {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   } else {
       fmt.Fprintln(writer)
   }
   fmt.Fprintln(writer, len(day2))
   if len(day2) > 0 {
       for i, v := range day2 {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   } else {
       fmt.Fprintln(writer)
   }
}
