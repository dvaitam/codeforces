package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var c byte
   var x int
   var neg bool
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       neg = true
   } else {
       x = int(c - '0')
   }
   for {
       b, err := reader.ReadByte()
       if err != nil {
           break
       }
       if b < '0' || b > '9' {
           break
       }
       x = x*10 + int(b-'0')
   }
   if neg {
       return -x
   }
   return x
}

func main() {
   defer writer.Flush()
   T := readInt()
   for T > 0 {
       T--
       n := readInt()
       m := readInt()
       K := readInt()
       Siz := n * (K - 1)
       N := Siz * 2
       adj := make([][]int, N+1)

       getind := func(i, j, k int) int {
           return (j-2)*n + i + k*Siz
       }
       link := func(u, v int) {
           adj[u] = append(adj[u], v)
       }
       linkit := func(u, v int) {
           link(u, v)
           var vn, un int
           if v > Siz {
               vn = v - Siz
           } else {
               vn = v + Siz
           }
           if u > Siz {
               un = u - Siz
           } else {
               un = u + Siz
           }
           link(vn, un)
       }

       // constraints
       for u := 1; u < n; u++ {
           for i := 2; i <= K; i++ {
               linkit(getind(u, i, 0), getind(u+1, i, 0))
           }
       }
       for u := 1; u <= n; u++ {
           for i := 3; i <= K; i++ {
               linkit(getind(u, i, 0), getind(u, i-1, 0))
           }
       }
       // operations
       for t := 0; t < m; t++ {
           opt := readInt()
           if opt == 1 {
               i := readInt()
               v := readInt()
               if v == K {
                   link(getind(i, v, 0), getind(i, v, 1))
               } else if v == 1 {
                   link(getind(i, 2, 1), getind(i, 2, 0))
               } else {
                   linkit(getind(i, v, 0), getind(i, v+1, 0))
               }
           } else if opt == 2 {
               i := readInt()
               j := readInt()
               v := readInt()
               if i > j {
                   i, j = j, i
               }
               if v <= K {
                   link(getind(j, v, 0), getind(j, v, 1))
               }
               for p := 2; p <= v && p <= K; p++ {
                   if v-p+1 <= K {
                       linkit(getind(i, p, 0), getind(j, max(p, v-p+1), 1))
                   }
               }
           } else {
               i := readInt()
               j := readInt()
               v := readInt()
               if i > j {
                   i, j = j, i
               }
               for p := 2; p <= (v+1>>1); p++ {
                   if v-p+1 <= K {
                       linkit(getind(i, p, 1), getind(j, v-p+1, 0))
                   } else {
                       link(getind(i, p, 1), getind(i, p, 0))
                   }
               }
           }
       }
       // tarjan
       dfn := make([]int, N+1)
       low := make([]int, N+1)
       bel := make([]int, N+1)
       vis := make([]bool, N+1)
       var S []int
       ind := 0
       cnt := 0
       var dfs func(u int)
       dfs = func(u int) {
           ind++
           dfn[u] = ind
           low[u] = ind
           vis[u] = true
           S = append(S, u)
           for _, v := range adj[u] {
               if dfn[v] == 0 {
                   dfs(v)
                   if low[v] < low[u] {
                       low[u] = low[v]
                   }
               } else if vis[v] && dfn[v] < low[u] {
                   low[u] = dfn[v]
               }
           }
           if low[u] == dfn[u] {
               cnt++
               for {
                   x := S[len(S)-1]
                   S = S[:len(S)-1]
                   bel[x] = cnt
                   vis[x] = false
                   if x == u {
                       break
                   }
               }
           }
       }
       for u := 1; u <= N; u++ {
           if dfn[u] == 0 {
               dfs(u)
           }
       }
       // check
       sat := true
       for u := 1; u <= Siz; u++ {
           if bel[u] == bel[u+Siz] {
               sat = false
               break
           }
       }
       if !sat {
           writer.WriteString("-1\n")
           continue
       }
       // output
       for x := 1; x <= n; x++ {
           lst := 1
           for i := 2; i <= K; i++ {
               if bel[getind(x, i, 0)] < bel[getind(x, i, 1)] {
                   lst++
               }
           }
           writer.WriteString(strconv.Itoa(lst))
           if x < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
