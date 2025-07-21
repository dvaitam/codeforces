package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   base := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       var v int64
       fmt.Fscan(reader, &v)
       base[i] = v
   }
   var add int64

   for i := 0; i < m; i++ {
       var t int
       fmt.Fscan(reader, &t)
       switch t {
       case 1:
           var idx int
           var x int64
           fmt.Fscan(reader, &idx, &x)
           // store value adjusted by current global addition
           base[idx] = x - add
       case 2:
           var y int64
           fmt.Fscan(reader, &y)
           add += y
       case 3:
           var q int
           fmt.Fscan(reader, &q)
           // output actual value
           fmt.Fprintln(writer, base[q]+add)
       }
   }
}
