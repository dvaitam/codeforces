package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func add(x, y int) int {
   x += y
   if x >= MOD {
       x -= MOD
   }
   return x
}

func sub(x, y int) int {
   x -= y
   if x < 0 {
       x += MOD
   }
   return x
}

type pair struct { first, second int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k, m int
   fmt.Fscan(reader, &n, &k, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] != -1 {
           a[i]--
       }
   }
   if m == 1 {
       fmt.Fprintln(writer, 0)
       return
   }
   // track last two occurrences: for each position store two (pos, val)
   lst := make([][2]pair, n+1)
   lst[0][0] = pair{-1, -1}
   lst[0][1] = pair{-1, -1}
   for i := 0; i < n; i++ {
       lst[i+1][0] = lst[i][0]
       lst[i+1][1] = lst[i][1]
       if a[i] == -1 {
           continue
       }
       if lst[i+1][0].second == a[i] {
           lst[i+1][0].first = i
       } else {
           lst[i+1][1] = lst[i+1][0]
           lst[i+1][0] = pair{i, a[i]}
       }
   }
   // dp[i][x]: count for prefix length i and ending with x at i-1
   dp := make([][]int, n+1)
   for i := range dp {
       dp[i] = make([]int, k)
   }
   sumdp := make([]int, n+1)
   sumdp[0] = 1
   // canPut checks if in window [l, r) value x occurs at most once
   canPut := func(x, l, r int) bool {
       for j := 0; j < 2; j++ {
           p := lst[r][j]
           if p.second != x && p.first >= l {
               return false
           }
       }
       return true
   }
   for i := 0; i < n; i++ {
       for x := 0; x < k; x++ {
           if a[i] != -1 && a[i] != x {
               continue
           }
           dp[i+1][x] = sumdp[i]
           if i+1 >= m && canPut(x, i+1-m, i+1) {
               dp[i+1][x] = sub(dp[i+1][x], sub(sumdp[i+1-m], dp[i+1-m][x]))
           }
           sumdp[i+1] = add(sumdp[i+1], dp[i+1][x])
       }
   }
   fmt.Fprintln(writer, sumdp[n])
}
