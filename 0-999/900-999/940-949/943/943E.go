package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   for i, v := range a {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(v))
   }
   writer.WriteByte('\n')
}
