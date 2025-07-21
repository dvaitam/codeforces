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

   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   chapters := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &chapters[i])
   }
   sort.Ints(chapters)
   var total int64
   for i, c := range chapters {
       t := x - i
       if t < 1 {
           t = 1
       }
       total += int64(t) * int64(c)
   }
   fmt.Fprint(writer, total)
}
