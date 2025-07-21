package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // Read chat participation matrix
   participated := make([][]bool, n)
   for i := 0; i < n; i++ {
       participated[i] = make([]bool, m)
       for j := 0; j < m; j++ {
           var x int
           fmt.Fscan(reader, &x)
           participated[i][j] = (x == 1)
       }
   }
   // Track total messages per chat and per-user per-chat
   totalMsgs := make([]int, m)
   userMsgs := make([][]int, n)
   for i := 0; i < n; i++ {
       userMsgs[i] = make([]int, m)
   }
   // Process events
   for e := 0; e < k; e++ {
       var u, chat int
       fmt.Fscan(reader, &u, &chat)
       u--
       chat--
       totalMsgs[chat]++
       userMsgs[u][chat]++
   }
   // Compute notifications per user
   res := make([]int, n)
   for i := 0; i < n; i++ {
       sum := 0
       for j := 0; j < m; j++ {
           if participated[i][j] {
               sum += totalMsgs[j] - userMsgs[i][j]
           }
       }
       res[i] = sum
   }
   // Output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range res {
       if i > 0 {
           out.WriteString(" ")
       }
       out.WriteString(fmt.Sprint(v))
   }
   out.WriteString("\n")
}
