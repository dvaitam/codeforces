package main

import (
   "bufio"
   "fmt"
   "os"
)

var dp [20][4]int64
var digits [20]int

func dfs(pos int, limit bool, had int) int64 {
   if pos > 19 {
       return 1
   }
   if !limit && dp[pos][had] != -1 {
       return dp[pos][had]
   }
   if had == 3 {
       if !limit {
           dp[pos][had] = 1
       }
       return 1
   }
   var res int64
   maxd := 9
   if limit {
       maxd = digits[pos]
   }
   for d := 0; d <= maxd; d++ {
       nextLimit := limit && (d == maxd)
       nh := had
       if d != 0 {
           nh++
       }
       res += dfs(pos+1, nextLimit, nh)
   }
   if !limit {
       dp[pos][had] = res
   }
   return res
}

func solve(x int64) int64 {
   if x < 0 {
       return 0
   }
   // fill digits from most significant at pos 1 to pos 19
   for i := 19; i >= 1; i-- {
       digits[i] = int(x % 10)
       x /= 10
   }
   return dfs(1, true, 0)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // init dp with -1
   for i := 1; i <= 19; i++ {
       for j := 0; j < 4; j++ {
           dp[i][j] = -1
       }
   }
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var l, r int64
       fmt.Fscan(reader, &l, &r)
       ans := solve(r) - solve(l-1)
       fmt.Fprintln(writer, ans)
   }
}
