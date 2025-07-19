package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

var (
   N, M int
   Y    int64
   Adj  [][]int
   InE  []int
)

func hasCycle(u int, visited, recStack []bool) bool {
   visited[u] = true
   recStack[u] = true
   for _, v := range Adj[u] {
       if !visited[v] {
           if hasCycle(v, visited, recStack) {
               return true
           }
       } else if recStack[v] {
           return true
       }
   }
   recStack[u] = false
   return false
}

func Checked() bool {
   visited := make([]bool, N)
   recStack := make([]bool, N)
   for i := 0; i < N; i++ {
       if !visited[i] {
           if hasCycle(i, visited, recStack) {
               return false
           }
       }
   }
   return true
}

// Find number of topological orders consistent with Given
func Find(Given []int) int64 {
   size := 1 << N
   dp := make([]int64, size)
   dp[0] = 1
   for mask := 0; mask < size-1; mask++ {
       if dp[mask] == 0 {
           continue
       }
       c := bits.OnesCount(uint(mask))
       if Given[c] == -1 {
           for j := 0; j < N; j++ {
               if (mask&(1<<j))==0 && (mask&InE[j])==InE[j] {
                   dp[mask|(1<<j)] += dp[mask]
               }
           }
       } else {
           j := Given[c]
           if (mask&InE[j])==InE[j] {
               dp[mask|(1<<j)] += dp[mask]
           }
       }
   }
   return dp[size-1]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &N, &Y, &M)
   Y -= 2000
   Adj = make([][]int, N)
   InE = make([]int, N)
   for i := 0; i < M; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       Adj[u] = append(Adj[u], v)
       InE[v] |= 1 << u
   }
   if !Checked() {
       fmt.Fprintln(writer, "The times have changed")
       return
   }
   // prepare arrays
   Given := make([]int, N)
   for i := 0; i < N; i++ {
       Given[i] = -1
   }
   Used := make([]bool, N)
   result := make([]int, N)
   for i := 0; i < N; i++ {
       flag := false
       for j := 0; j < N; j++ {
           if !Used[j] {
               Given[i] = j
               Used[j] = true
               cnt := Find(Given)
               if cnt >= Y {
                   result[i] = j
                   flag = true
                   break
               }
               Y -= cnt
               Given[i] = -1
               Used[j] = false
           }
       }
       if !flag {
           fmt.Fprintln(writer, "The times have changed")
           return
       }
   }
   for i, v := range result {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v+1)
   }
   writer.WriteByte('\n')
}
