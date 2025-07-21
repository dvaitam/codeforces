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

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // star matrix dimensions
   N := n - 2
   M := m - 2
   if N <= 0 || M <= 0 {
       fmt.Fprint(writer, 0)
       return
   }
   // build star matrix
   star := make([][]int, N)
   for i := 0; i < N; i++ {
       star[i] = make([]int, M)
       for j := 0; j < M; j++ {
           gi := i + 1
           gj := j + 1
           if grid[gi][gj] == '1' && grid[gi-1][gj] == '1' && grid[gi+1][gj] == '1' && grid[gi][gj-1] == '1' && grid[gi][gj+1] == '1' {
               star[i][j] = 1
           }
       }
   }
   colSum := make([]int, M)
   S := make([]int, M+1)
   var ans int64
   // iterate over row ranges in star matrix
   for l := 0; l < N; l++ {
       // reset column sums
       for j := 0; j < M; j++ {
           colSum[j] = 0
       }
       for r := l; r < N; r++ {
           // add row r
           for j := 0; j < M; j++ {
               colSum[j] += star[r][j]
           }
           // prefix sums over columns
           S[0] = 0
           for j := 1; j <= M; j++ {
               S[j] = S[j-1] + colSum[j-1]
           }
           // two-pointer over prefix to count weighted column ranges
           var innerSum int64
           c := 0
           for R := 1; R <= M; R++ {
               if S[R] < k {
                   continue
               }
               target := S[R] - k
               for c <= R-1 && S[c] <= target {
                   c++
               }
               t := c
               tri := int64(t) * int64(t+1) / 2
               innerSum += tri * int64(M+1-R)
           }
           if innerSum != 0 {
               A := int64(l+1) * int64(N-r)
               ans += A * innerSum
           }
       }
   }
   fmt.Fprint(writer, ans)
}
