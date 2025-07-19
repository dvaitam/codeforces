package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   m := len(s)
   // 1-based arrays
   a := make([]int, m+2)
   for i := 1; i <= m; i++ {
       a[i] = int(s[i-1] - '0')
   }
   Next := make([]int, m+2)
   // build KMP next
   j := 0
   for i := 2; i <= m; i++ {
       for j > 0 && a[j+1] != a[i] {
           j = Next[j]
       }
       if a[j+1] == a[i] {
           j++
       }
       Next[i] = j
   }
   // used positions: -1 means free, 0 or 1 forced
   used := make([]int, n+2)
   // dp array F[i][k][ok]
   F := make([][][]int64, n+2)
   for i := 0; i <= n; i++ {
       F[i] = make([][]int64, m)
       for k := 0; k < m; k++ {
           F[i][k] = make([]int64, 2)
       }
   }
   tmp := make([][]int64, m)
   for k := 0; k < m; k++ {
       tmp[k] = make([]int64, 2)
   }
   var ans int64
   // clear dp
   clear := func() {
       for i := 0; i <= n; i++ {
           for k := 0; k < m; k++ {
               F[i][k][0] = 0
               F[i][k][1] = 0
           }
       }
   }
   // dp from st to ed inclusive
   var dp func(st, ed int)
   dp = func(st, ed int) {
       for i := st; i <= ed; i++ {
           for k := 0; k < m; k++ {
               for ok := 0; ok < 2; ok++ {
                   cnt := F[i-1][k][ok]
                   if cnt == 0 {
                       continue
                   }
                   for cur := 0; cur < 2; cur++ {
                       if used[i] != -1 && used[i] != cur {
                           continue
                       }
                       j := k
                       for j > 0 && a[j+1] != cur {
                           j = Next[j]
                       }
                       if a[j+1] == cur {
                           j++
                       }
                       if j == m {
                           F[i][0][1] += cnt
                       } else {
                           F[i][j][ok] += cnt
                       }
                   }
               }
           }
       }
   }
   // first: full string
   for i := 1; i <= n; i++ {
       used[i] = -1
   }
   clear()
   F[0][0][0] = 1
   dp(1, n)
   for k := 0; k < m; k++ {
       ans += F[n][k][1]
   }
   // rotations
   for length := 1; length < m; length++ {
       for i := 1; i <= n; i++ {
           used[i] = -1
       }
       // force rotated pattern
       for i := 1; i <= length; i++ {
           used[n-length+i] = a[i]
       }
       for i := length + 1; i <= m; i++ {
           used[i-length] = a[i]
       }
       clear()
       start := n - length + 1
       F[start][0][0] = 1
       dp(start+1, n)
       for k := 0; k < m; k++ {
           for ok := 0; ok < 2; ok++ {
               tmp[k][ok] = F[n][k][ok]
           }
       }
       clear()
       for k := 0; k < m; k++ {
           for ok := 0; ok < 2; ok++ {
               F[0][k][ok] = tmp[k][ok]
           }
       }
       dp(1, n)
       for k := 0; k < m; k++ {
           ans += F[n][k][0]
       }
   }
   fmt.Println(ans)
}
