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

   var n int
   var h int64
   if _, err := fmt.Fscan(reader, &n, &h); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // find position of minimum element in original array
   pos := 0
   for i := 1; i < n; i++ {
       if a[i] < a[pos] {
           pos = i
       }
   }
   // make a copy and sort
   b := make([]int64, n)
   copy(b, a)
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   // compute values
   val1 := b[n-1] + b[n-2] - b[0] - b[1]
   max12 := b[n-1] + b[n-2]
   if b[0]+b[n-1]+h > max12 {
       max12 = b[0] + b[n-1] + h
   }
   min01 := b[1] + b[2]
   if b[0]+b[1]+h < min01 {
       min01 = b[0] + b[1] + h
   }
   val2 := max12 - min01
   // choose arrangement
   if val1 < val2 {
       pos = -1
   }
   res := val1
   if val2 < res {
       res = val2
   }
   // output
   fmt.Fprintln(writer, res)
   for i := 0; i < n; i++ {
       if i == pos {
           writer.WriteString("2 ")
       } else {
           writer.WriteString("1 ")
       }
   }
}
