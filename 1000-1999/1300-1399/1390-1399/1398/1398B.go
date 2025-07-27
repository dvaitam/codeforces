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
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       var runs []int
       cnt := 0
       for i := 0; i < len(s); i++ {
           if s[i] == '1' {
               cnt++
           } else {
               if cnt > 0 {
                   runs = append(runs, cnt)
                   cnt = 0
               }
           }
       }
       if cnt > 0 {
           runs = append(runs, cnt)
       }
       sort.Sort(sort.Reverse(sort.IntSlice(runs)))
       score := 0
       for i := 0; i < len(runs); i += 2 {
           score += runs[i]
       }
       fmt.Fprintln(writer, score)
   }
}
