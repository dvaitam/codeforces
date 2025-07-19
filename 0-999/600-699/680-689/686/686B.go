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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var ops [][2]int
   mark := make([]int, n+2)
   b := true
   for b {
       b = false
       mk := func(cnt int) {
           // mark and swap pairs starting at cnt
           for i := cnt; i <= n; i += 2 {
               if i+1 > n {
                   break
               }
               if a[i] > a[i+1] {
                   mark[i], mark[i+1] = 1, 1
                   a[i], a[i+1] = a[i+1], a[i]
               } else if a[i] == a[i+1] {
                   mark[i], mark[i+1] = 2, 2
               }
           }
           // collect continuous marked segments
           st := 0
           for i := 0; i <= n+1; i++ {
               if st != 0 && (i > n || mark[i] == 0) {
                   ops = append(ops, [2]int{st, i - 1})
                   st = 0
                   b = true
               } else if st == 0 && i <= n && mark[i] == 1 {
                   st = i
               }
           }
           // reset marks
           for i := range mark {
               mark[i] = 0
           }
       }
       mk(1)
       mk(2)
   }
   for _, op := range ops {
       fmt.Fprintf(writer, "%d %d\n", op[0], op[1])
   }
}
