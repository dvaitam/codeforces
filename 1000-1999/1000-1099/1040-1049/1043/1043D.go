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

   var N, M int
   fmt.Fscan(reader, &N, &M)
   // Read M rows of N integers
   a := make([][]int, M)
   for i := 0; i < M; i++ {
       a[i] = make([]int, N+2)
       for j := 1; j <= N; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }

   father := make([]int, N+1)
   cnt := make([]int, N+1)
   // Initialize mapping from first row
   for j := 1; j <= N; j++ {
       father[a[0][j]] = a[0][j+1]
   }
   // Refine mapping based on other rows
   for i := 1; i < M; i++ {
       for j := 1; j <= N; j++ {
           if father[a[i][j]] != a[i][j+1] {
               father[a[i][j]] = 0
           }
       }
   }
   // Set broken mappings to self
   for i := 1; i <= N; i++ {
       if father[i] == 0 {
           father[i] = i
       }
   }
   // DSU find with path compression
   var getfather func(int) int
   getfather = func(x int) int {
       if father[x] != x {
           father[x] = getfather(father[x])
       }
       return father[x]
   }
   // Count elements per chain root
   for i := 1; i <= N; i++ {
       root := getfather(i)
       if root != 0 {
           cnt[root]++
       }
   }
   // Compute answer: sum of sizes and pairs
   var ans int64 = int64(N)
   for i := 1; i <= N; i++ {
       if cnt[i] > 1 {
           ans += int64(cnt[i]) * int64(cnt[i]-1) / 2
       }
   }
   fmt.Fprintln(writer, ans)
}
