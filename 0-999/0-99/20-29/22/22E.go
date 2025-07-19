package main

import (
   "bufio"
   "fmt"
   "os"
)

var rdr = bufio.NewReader(os.Stdin)
var wr = bufio.NewWriter(os.Stdout)

// readInt reads next integer from standard input.
func readInt() int {
   c, err := rdr.ReadByte()
   for err == nil && (c < '0' || c > '9') && c != '-' {
       c, err = rdr.ReadByte()
   }
   if err != nil {
       return 0
   }
   sig := 1
   if c == '-' {
       sig = -1
       c, _ = rdr.ReadByte()
   }
   x := 0
   for err == nil && c >= '0' && c <= '9' {
       x = x*10 + int(c - '0')
       c, err = rdr.ReadByte()
   }
   return x * sig
}

func main() {
   defer wr.Flush()
   n := readInt()
   a := make([]int, n+1)
   d := make([]int, n+1)
   vis := make([]int, n+1)
   f := make([]int, n+1)
   for i := 1; i <= n; i++ {
       a[i] = readInt()
       d[a[i]]++
   }
   var dfs func(int) int
   dfs = func(x int) int {
       if f[x] != 0 {
           return f[x]
       }
       if vis[x] != 0 {
           return 0
       }
       vis[x] = 1
       res := dfs(a[x])
       f[x] = res
       if f[x] != 0 {
           return f[x]
       }
       f[x] = x
       return x
   }
   for i := 1; i <= n; i++ {
       if d[i] == 0 || vis[i] == 0 {
           dfs(i)
       }
   }
   ans1 := make([]int, 0)
   ans2 := make([]int, 0)
   for i := 1; i <= n; i++ {
       if d[i] == 0 {
           ans1 = append(ans1, i)
           ans2 = append(ans2, f[i])
           vis[f[i]] = 2
       }
   }
   for i := 1; i <= n; i++ {
       if f[i] == i && vis[i] < 2 {
           ans1 = append(ans1, i)
           ans2 = append(ans2, i)
       }
   }
   m := len(ans1)
   if m == 1 && ans1[0] == ans2[0] {
       m = 0
   }
   fmt.Fprintln(wr, m)
   for i := 0; i < m; i++ {
       j := (i + 1) % m
       fmt.Fprintln(wr, ans2[i], ans1[j])
   }
}
