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

   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   A := make([]int, N)
   B := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &A[i], &B[i])
   }
   sort.Ints(A)
   sort.Ints(B)

   var ans int64 = int64(N)
   for i := 0; i < N; i++ {
       if A[i] > B[i] {
           ans += int64(A[i])
       } else {
           ans += int64(B[i])
       }
   }
   fmt.Fprintln(writer, ans)
}
