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

   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   A := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i])
   }
   stock := make([]struct{ cnt, id int }, N)
   ans := make([]int, N)
   lastid := 1
   for i := 0; i < N; i++ {
       k := (N - 1) - A[i]
       if stock[k].cnt <= 0 {
           ans[i] = lastid
           stock[k].cnt = k
           stock[k].id = lastid
           lastid++
       } else {
           ans[i] = stock[k].id
           stock[k].cnt--
       }
   }
   for i := 0; i < N; i++ {
       if stock[i].cnt != 0 {
           fmt.Fprintln(writer, "Impossible")
           return
       }
   }
   fmt.Fprintln(writer, "Possible")
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
