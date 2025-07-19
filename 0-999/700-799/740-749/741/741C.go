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
   total := 2 * n
   x := make([]int, total)
   chk := make([]bool, total)
   ans := make([]int, total)
   pairs := make([][2]int, n)

   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       pairs[i][0] = a
       pairs[i][1] = b
       x[a] = b
       x[b] = a
   }

   for i := 0; i < total; i++ {
       if chk[i] {
           continue
       }
       cur := i
       for !chk[cur] {
           ans[cur] = 0
           chk[cur] = true
           cur = x[cur]
           ans[cur] = 1
           cur ^= 1
       }
   }

   for i := 0; i < n; i++ {
       a, b := pairs[i][0], pairs[i][1]
       fmt.Fprintf(writer, "%d %d\n", ans[a]+1, ans[b]+1)
   }
}
