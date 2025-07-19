package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   MaxN = 555
   MaxK = 11
)

var reader = bufio.NewReader(os.Stdin)
var g [MaxN][10]int
var costArr [MaxN]int64
var d [MaxN][MaxN][MaxK]int64
var used [MaxN][MaxN]bool
var G, k int

func readInt() int {
   sign := 1
   c, err := reader.ReadByte()
   for err == nil && (c < '0' || c > '9') && c != '-' {
       c, err = reader.ReadByte()
   }
   if err != nil {
       return 0
   }
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   val := 0
   for err == nil && c >= '0' && c <= '9' {
       val = val*10 + int(c-'0')
       c, err = reader.ReadByte()
   }
   return sign * val
}

func readString() string {
   c, err := reader.ReadByte()
   for err == nil && c <= ' ' {
       c, err = reader.ReadByte()
   }
   if err != nil {
       return ""
   }
   var sb []byte
   for err == nil && c > ' ' {
       sb = append(sb, c)
       c, err = reader.ReadByte()
   }
   return string(sb)
}

func max64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func dfs(x, p, dip int) {
   if used[x][p] {
       return
   }
   used[x][p] = true
   for i := 0; i < 10; i++ {
       y := g[x][i]
       if y == 0 {
           continue
       }
       dfs(y, p, dip+1)
       for j := k; j >= 0; j-- {
           for h := 0; h <= j; h++ {
               cand := d[x][p][j-h] + d[y][p][h]
               if d[x][p][j] < cand {
                   d[x][p][j] = cand
               }
           }
       }
   }
   if dip != p {
       dfs(x, dip, dip)
       for i := 1; i <= k; i++ {
           cand := d[x][p][i-1] + int64(dip-p)*costArr[x]
           if d[x][p][i] < cand {
               d[x][p][i] = cand
           }
       }
   }
}

func main() {
   n := readInt()
   k = readInt()
   for i := 0; i < n; i++ {
       s := readString()
       cost := int64(readInt())
       x := 0
       for _, ch := range s {
           d := int(ch - '0')
           if g[x][d] == 0 {
               G++
               g[x][d] = G
           }
           x = g[x][d]
           costArr[x] += cost
       }
   }
   dfs(0, 0, 0)
   var ans int64
   for i := 1; i <= k; i++ {
       if ans < d[0][0][i] {
           ans = d[0][0][i]
       }
   }
   ans = -ans
   for i := 1; i <= G; i++ {
       ans += costArr[i]
   }
   fmt.Println(ans)
}
