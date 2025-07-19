package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   vals := make([]int64, 2*n)
   for i := range vals {
       fmt.Fscan(reader, &vals[i])
   }
   sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
   allEqual := true
   for i := 0; i < len(vals)-1; i++ {
       if vals[i] != vals[i+1] {
           allEqual = false
           break
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if allEqual {
       fmt.Fprint(writer, -1)
   } else {
       for i, x := range vals {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
   }
}
