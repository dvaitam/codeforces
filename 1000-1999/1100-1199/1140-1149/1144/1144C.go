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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   count := make([]int, 200001)
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       count[x]++
   }
   for _, c := range count {
       if c > 2 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }
   inc := make([]int, 0, n)
   dec := make([]int, 0, n)
   for x, c := range count {
       if c >= 1 {
           inc = append(inc, x)
       }
       if c == 2 {
           dec = append(dec, x)
       }
   }
   sort.Ints(inc)
   sort.Sort(sort.Reverse(sort.IntSlice(dec)))
   fmt.Fprintln(writer, "YES")
   // increasing sequence
   fmt.Fprintln(writer, len(inc))
   for i, v := range inc {
       if i > 0 {
           writer.WriteString(" ")
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteString("\n")
   // decreasing sequence
   fmt.Fprintln(writer, len(dec))
   for i, v := range dec {
       if i > 0 {
           writer.WriteString(" ")
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteString("\n")
}
