package main

import (
   "bufio"
   "fmt"
   "os"
)

func readInt(r *bufio.Reader) int {
   neg := false
   var c byte
   // skip non-numeric
   for {
       b, err := r.ReadByte()
       if err != nil {
           return 0
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       neg = true
       b, _ := r.ReadByte()
       c = b
   }
   x := 0
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       c = b
   }
   if neg {
       x = -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   n := readInt(reader)
   m := readInt(reader)
   // count occurrences
   cnt := make([]int, m+3)
   for i := 0; i < n; i++ {
       x := readInt(reader)
       if x >= 1 && x <= m {
           cnt[x]++
       }
   }
   // dp[i][j][k]: max for first i with used j,k for i-2,i-1
   dp := make([][3][3]int, m+3)
   vis := make([][3][3]bool, m+3)
   vis[0][0][0] = true
   for i := 1; i <= m; i++ {
       for j := 0; j <= 2; j++ {
           for k := 0; k <= 2; k++ {
               if !vis[i-1][j][k] {
                   continue
               }
               for cur := 0; cur <= 2; cur++ {
                   // check availability
                   if cnt[i] < cur {
                       break
                   }
                   if cnt[i-1] < k+cur {
                       break
                   }
                   if i-2 >= 1 {
                       if cnt[i-2] < j+k+cur {
                           break
                       }
                   } else if cur > 0 {
                       break
                   }
                   // compute value
                   val := dp[i-1][j][k] + cur
                   if i-2 >= 1 {
                       val += (cnt[i-2] - k - j - cur) / 3
                   }
                   if !vis[i][k][cur] || dp[i][k][cur] < val {
                       dp[i][k][cur] = val
                       vis[i][k][cur] = true
                   }
               }
           }
       }
   }
   ans := 0
   best := 0
   for j := 0; j <= 2; j++ {
       for k := 0; k <= 2; k++ {
           if !vis[m][j][k] {
               continue
           }
           v := dp[m][j][k]
           if m-1 >= 1 {
               v += (cnt[m-1] - j - k) / 3
           }
           v += (cnt[m] - k) / 3
           if v > best {
               best = v
           }
       }
   }
   fmt.Println(ans + best)
}
