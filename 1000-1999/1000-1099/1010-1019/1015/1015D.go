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

   var N, M, K int64
   if _, err := fmt.Fscan(reader, &N, &M, &K); err != nil {
       return
   }
   if K < M || (N-1)*M < K {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")

   ans := make([]int64, M+2)
   P := K / M
   for i := int64(1); i <= M; i++ {
       if i%2 == 1 {
           ans[i] = 1 + P
       } else {
           ans[i] = 1
       }
   }
   r := K % M
   if r%2 == 1 {
       for i := int64(1); i <= r-1; i++ {
           if i%2 == 1 {
               ans[i]++
           }
       }
       if M%2 == 1 {
           ans[M]++
       } else {
           ans[M-1] = 1 + P + 1
           ans[M] = 2
       }
   } else {
       for i := int64(1); i <= r; i++ {
           if i%2 == 1 {
               ans[i]++
           }
       }
   }
   for i := int64(1); i <= M; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
   writer.WriteByte('\n')
}
