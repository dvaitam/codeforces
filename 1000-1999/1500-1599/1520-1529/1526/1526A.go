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
   for ti := 0; ti < t; ti++ {
       var n int
       fmt.Fscan(reader, &n)
       v := make([]int, 2*n)
       for i := 0; i < 2*n; i++ {
           fmt.Fscan(reader, &v[i])
       }
       sort.Ints(v)
       l, r := 0, len(v)-1
       for l < r {
           // print smallest and largest remaining
           writer.WriteString(fmt.Sprintf("%d %d", v[l], v[r]))
           if l+1 < r {
               writer.WriteByte(' ')
           }
           l++
           r--
       }
       if ti != t-1 {
           writer.WriteByte('\n')
       }
   }
}
