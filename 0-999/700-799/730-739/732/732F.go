package main

import (
   "bufio"
   "os"
   "strconv"
)

var (
   br   = bufio.NewReader(os.Stdin)
   bw   = bufio.NewWriter(os.Stdout)
   n, m int
   ecnt int
   to   []int
   nxt  []int
   g    []int
   used []bool
   pos  []int
   low  []int
   dfn  int
   vis  []bool
   inq  []bool
   ans  [][2]int
   q    []int
   sz   int
)

func F() int {
   c, _ := br.ReadByte()
   for (c < '0' || c > '9') && c != '-' {
       c, _ = br.ReadByte()
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = br.ReadByte()
   }
   x := 0
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, _ = br.ReadByte()
   }
   return x * sign
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func ins(u, v int) {
   ecnt++
   to[ecnt] = v
   nxt[ecnt] = g[u]
   g[u] = ecnt
   ecnt++
   to[ecnt] = u
   nxt[ecnt] = g[v]
   g[v] = ecnt
}

func dfs1(x, fa int) {
   dfn++
   pos[x] = dfn
   low[x] = dfn
   for i := g[x]; i != 0; i = nxt[i] {
       v := to[i]
       if v == fa {
           continue
       }
       if pos[v] == 0 {
           dfs1(v, x)
           if low[v] > pos[x] {
               used[i] = true
               used[i^1] = true
           }
           low[x] = min(low[x], low[v])
       } else {
           low[x] = min(low[x], pos[v])
       }
   }
}

func dfs2(x int) {
   vis[x] = true
   sz++
   for i := g[x]; i != 0; i = nxt[i] {
       if used[i] {
           continue
       }
       eid := i >> 1
       ans[eid][0] = x
       ans[eid][1] = to[i]
       v := to[i]
       if !vis[v] {
           dfs2(v)
       }
   }
}

func bfs(x int) {
   h, t := 0, 0
   q[t] = x
   inq[x] = true
   for h <= t {
       u := q[h]
       h++
       for i := g[u]; i != 0; i = nxt[i] {
           v := to[i]
           if inq[v] {
               continue
           }
           if used[i] {
               eid := i >> 1
               ans[eid][0] = v
               ans[eid][1] = u
           }
           inq[v] = true
           t++
           q[t] = v
       }
   }
}

func main() {
   defer bw.Flush()
   n = F()
   m = F()
   // allocate
   size := 2*m + 2
   to = make([]int, size)
   nxt = make([]int, size)
   g = make([]int, n+1)
   used = make([]bool, size)
   pos = make([]int, n+1)
   low = make([]int, n+1)
   vis = make([]bool, n+1)
   inq = make([]bool, n+1)
   ans = make([][2]int, m+1)
   q = make([]int, n+1)
   // initialize edge counter
   ecnt = 1
   // read edges
   for i := 1; i <= m; i++ {
       u := F()
       v := F()
       ins(u, v)
   }
   dfs1(1, 0)
   mx := 0
   mxl := 1
   for i := 1; i <= n; i++ {
       if !vis[i] {
           sz = 0
           dfs2(i)
           if sz > mx {
               mx = sz
               mxl = i
           }
       }
   }
   bfs(mxl)
   // output
   bw.WriteString(strconv.Itoa(mx))
   bw.WriteByte('\n')
   for i := 1; i <= m; i++ {
       bw.WriteString(strconv.Itoa(ans[i][0]))
       bw.WriteByte(' ')
       bw.WriteString(strconv.Itoa(ans[i][1]))
       bw.WriteByte('\n')
   }
}
