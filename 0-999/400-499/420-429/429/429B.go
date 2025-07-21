package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([][]int64, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // dp1: from (0,0) to (i,j)
   dp1 := make([][]int64, n)
   // dp2: from (n-1,m-1) to (i,j)
   dp2 := make([][]int64, n)
   // dp3: from (n-1,0) to (i,j)
   dp3 := make([][]int64, n)
   // dp4: from (0,m-1) to (i,j)
   dp4 := make([][]int64, n)
   for i := 0; i < n; i++ {
       dp1[i] = make([]int64, m)
       dp2[i] = make([]int64, m)
       dp3[i] = make([]int64, m)
       dp4[i] = make([]int64, m)
   }
   // compute dp1
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           dp1[i][j] = a[i][j]
           if i > 0 {
               dp1[i][j] = max(dp1[i][j], dp1[i-1][j]+a[i][j])
           }
           if j > 0 {
               dp1[i][j] = max(dp1[i][j], dp1[i][j-1]+a[i][j])
           }
       }
   }
   // compute dp2
   for i := n - 1; i >= 0; i-- {
       for j := m - 1; j >= 0; j-- {
           dp2[i][j] = a[i][j]
           if i+1 < n {
               dp2[i][j] = max(dp2[i][j], dp2[i+1][j]+a[i][j])
           }
           if j+1 < m {
               dp2[i][j] = max(dp2[i][j], dp2[i][j+1]+a[i][j])
           }
       }
   }
   // compute dp3
   for i := n - 1; i >= 0; i-- {
       for j := 0; j < m; j++ {
           dp3[i][j] = a[i][j]
           if i+1 < n {
               dp3[i][j] = max(dp3[i][j], dp3[i+1][j]+a[i][j])
           }
           if j > 0 {
               dp3[i][j] = max(dp3[i][j], dp3[i][j-1]+a[i][j])
           }
       }
   }
   // compute dp4
   for i := 0; i < n; i++ {
       for j := m - 1; j >= 0; j-- {
           dp4[i][j] = a[i][j]
           if i > 0 {
               dp4[i][j] = max(dp4[i][j], dp4[i-1][j]+a[i][j])
           }
           if j+1 < m {
               dp4[i][j] = max(dp4[i][j], dp4[i][j+1]+a[i][j])
           }
       }
   }
   var ans int64
   // iterate possible meeting points (i,j)
   for i := 1; i < n-1; i++ {
       for j := 1; j < m-1; j++ {
           // two crossing patterns
           v1 := dp1[i][j-1] + dp2[i][j+1] + dp3[i+1][j] + dp4[i-1][j]
           v2 := dp1[i-1][j] + dp2[i+1][j] + dp3[i][j-1] + dp4[i][j+1]
           if v1 > ans {
               ans = v1
           }
           if v2 > ans {
               ans = v2
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
