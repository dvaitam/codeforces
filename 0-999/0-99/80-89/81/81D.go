package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   v := make([]int, m+1)
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &v[i])
   }
   ans := make([]int, n)
   last := -1
   for i := 0; i < n; i++ {
       id := 0
       for j := 1; j <= m; j++ {
           if j == last || (i == n-1 && j == ans[0]) {
               continue
           }
           if v[j] > v[id] {
               id = j
           } else if v[j] == v[id] && j == ans[0] {
               id = j
           }
       }
       if v[id] == 0 {
           fmt.Println(-1)
           return
       }
       v[id]--
       ans[i] = id
       last = id
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteString(" ")
       }
       fmt.Fprint(writer, ans[i])
   }
   writer.WriteByte('\n')
}
