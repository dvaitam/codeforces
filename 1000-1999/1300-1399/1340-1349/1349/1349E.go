package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct{ first, second int }

var (
   n, m    int
   t       []int64
   rg      []pair
   dp, to  [][]int
   res     []int
   writer  *bufio.Writer
)

func calc1(x int64) int64 {
   return x * (x - 1) / 2
}

func clac2(x, y int64) int64 {
   return x*y + calc1(y)
}

func check(L, R int, s int64) bool {
   if s < 0 {
       return false
   }
   low, high, ans := 0, R-L+1, 0
   for high-low > 5 {
       mid := (low + high) >> 1
       if clac2(int64(L), int64(mid)) <= s {
           low, ans = mid, mid
       } else {
           high = mid
       }
   }
   for i := low; i < high; i++ {
       if clac2(int64(L), int64(i)) <= s {
           ans = i
           break
       }
   }
   sl := clac2(int64(L), int64(ans))
   sr := clac2(int64(R-ans+1), int64(ans))
   return sl <= s && s <= sr
}

func dfs(l, r int, s int64) {
   if l > r {
       return
   }
   if check(l+1, r, s) {
       dfs(l+1, r, s)
   } else {
       res[l] = 1
       dfs(l+1, r, s-int64(l))
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   t = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   // reverse t[1..n]
   for i, j := 1, n; i < j; i, j = i+1, j-1 {
       t[i], t[j] = t[j], t[i]
   }
   rg = make([]pair, n+2)
   // build segments
   for i := 1; i <= n; i++ {
       if t[i] != 0 {
           if m > 0 && t[i] <= t[rg[m].second] {
               rg[m].first = i
           } else {
               m++
               rg[m] = pair{i, i}
           }
       }
   }
   if m == 0 {
       for i := 1; i <= n; i++ {
           writer.WriteByte('0')
       }
       writer.WriteByte('\n')
       return
   }
   // reverse rg[1..m]
   for i, j := 1, m; i < j; i, j = i+1, j-1 {
       rg[i], rg[j] = rg[j], rg[i]
   }
   // count in first segment
   cnt, cntt := 0, 0
   for i := rg[1].first; i <= rg[1].second; i++ {
       if t[i] != 0 {
           if t[i] == t[rg[1].second] {
               cnt++
           } else {
               cntt++
           }
       }
   }
   dp = make([][]int, m+3)
   to = make([][]int, m+3)
   for i := range dp {
       dp[i] = make([]int, 2)
       to[i] = make([]int, 2)
   }
   dp[m+1][1] = n + 1
   // DP
   for i := m; i >= 2; i-- {
       for j1 := 0; j1 <= 1; j1++ {
           if dp[i+1][j1] != 0 {
               for j2 := 0; j2 <= 1; j2++ {
                   if j2 == 1 || t[rg[i].first] == t[rg[i].second] {
                       lst := t[rg[i+1].first] + int64(j1) - 1
                       limit := dp[i][j2]
                       for k := rg[i].first; k > rg[i-1].second && k > limit; k-- {
                           if j2 == 0 && k == rg[i].first {
                               continue
                           }
                           if j2 == 1 && k < rg[i].first {
                               break
                           }
                           if check(rg[i].second+1, dp[i+1][j1]-1, t[rg[i].first]+int64(j2)-1-int64(k)-lst) {
                               dp[i][j2] = k
                               to[i][j2] = j1
                           }
                       }
                   }
               }
           }
       }
   }
   res = make([]int, n+2)
   // reconstruct
   for i := 0; i <= 1; i++ {
       if dp[2][i] != 0 {
           lst := t[rg[2].first] + int64(i) - 1
           for j := 0; j <= rg[1].second; j++ {
               if cnt != 1 && t[j] == t[rg[1].second] {
                   continue
               }
               if cntt != 0 && t[j] != t[rg[1].second]-1 {
                   continue
               }
               now := t[j]
               if now == 0 {
                   now = t[rg[1].second] - 1
               }
               if check(rg[1].second+1, dp[2][i]-1, now-lst-int64(j)) {
                   res[j] = 1
                   dfs(rg[1].second+1, dp[2][i]-1, now-lst-int64(j))
                   cur := i
                   for k := 2; k <= m; k++ {
                       nxt := to[k][cur]
                       lstt := t[rg[k+1].first] + int64(nxt) - 1
                       noww := t[rg[k].first] + int64(cur) - 1
                       res[dp[k][cur]] = 1
                       dfs(rg[k].second+1, dp[k+1][nxt]-1, noww-lstt-int64(dp[k][cur]))
                       cur = nxt
                   }
                   for i2 := 1; i2 <= n; i2++ {
                       if res[i2] == 1 {
                           writer.WriteByte('1')
                       } else {
                           writer.WriteByte('0')
                       }
                   }
                   writer.WriteByte('\n')
                   return
               }
           }
       }
   }
   // fallback
   fmt.Fprintf(writer, "%d %d\n", dp[2][0], dp[2][1])
}
