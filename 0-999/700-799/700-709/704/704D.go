package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1000000000

var (
   n, m, r, b int
   fl bool
   x, y []int
   b1, c1 []int
   cx, cy, dx, dy []int
   S, T, SS, TT int
   head, nxt, to, capArr []int
   eCount int
   Sx, xT []int
   h, cur, q []int
   eid []int
)

func addEdge(u, v, c int) {
   to[eCount] = v
   capArr[eCount] = c
   nxt[eCount] = head[u]
   head[u] = eCount
   eCount++
}

// insertEdge adds edge u->v with capacity c and reverse edge v->u with 0 capacity.
// Returns the index of the reverse edge.
func insertEdge(u, v, c int) int {
   // forward edge id = eCount
   addEdge(u, v, c)
   // backward edge id = eCount
   addEdge(v, u, 0)
   return eCount - 1
}

// insertLower adds an edge with lower bound l and upper bound r.
// It accumulates demands in Sx and xT, and adds edge of capacity r-l.
func insertLower(u, v, l, r int) int {
   Sx[v] += l
   xT[u] += l
   return insertEdge(u, v, r-l)
}

func build() {
   SS = 1
   TT = T + 1
   // add edges from SS and to TT for lower bound demands
   for i := S; i <= T; i++ {
       if Sx[i] > 0 {
           insertEdge(SS, i, Sx[i])
       }
       if xT[i] > 0 {
           insertEdge(i, TT, xT[i])
       }
   }
}

func bfs() bool {
   for i := SS; i <= TT; i++ {
       h[i] = -1
   }
   headQ, tailQ := 0, 0
   q[tailQ] = SS
   h[SS] = 0
   tailQ++
   for headQ < tailQ {
       u := q[headQ]
       headQ++
       for i := head[u]; i != -1; i = nxt[i] {
           v := to[i]
           if capArr[i] > 0 && h[v] < 0 {
               h[v] = h[u] + 1
               q[tailQ] = v
               tailQ++
           }
       }
   }
   return h[TT] >= 0
}

func dfs(u, f int) int {
   if u == TT {
       return f
   }
   used := 0
   for i := cur[u]; i != -1; i = nxt[i] {
       v := to[i]
       if capArr[i] > 0 && h[v] == h[u] + 1 {
           can := f - used
           if capArr[i] < can {
               can = capArr[i]
           }
           pushed := dfs(v, can)
           if pushed > 0 {
               capArr[i] -= pushed
               capArr[i^1] += pushed
               used += pushed
               if used == f {
                   cur[u] = i
                   break
               }
           }
       }
       cur[u] = i
   }
   if used == 0 {
       h[u] = -1
   }
   return used
}

func dinic() {
   for bfs() {
       for i := SS; i <= TT; i++ {
           cur[i] = head[i]
       }
       dfs(SS, INF)
   }
}

func findPos(a []int, x int) int {
   l, r := 1, n+1
   for l+1 < r {
       mid := (l + r) >> 1
       if a[mid] <= x {
           l = mid
       } else {
           r = mid
       }
   }
   return l
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &r, &b)
   if r < b {
       r, b = b, r
       fl = true
   }
   x = make([]int, n+1)
   y = make([]int, n+1)
   b1 = make([]int, n+1)
   c1 = make([]int, n+1)
   cx = make([]int, n+1)
   cy = make([]int, n+1)
   dx = make([]int, n+1)
   dy = make([]int, n+1)
   eid = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
       b1[i] = x[i]
       c1[i] = y[i]
   }
   sort.Ints(b1[1:])
   sort.Ints(c1[1:])
   for i := 1; i <= n; i++ {
       xi := findPos(b1, x[i])
       yi := findPos(c1, y[i])
       x[i] = xi
       y[i] = yi
       cx[xi]++
       cy[yi]++
   }
   for i := 1; i <= n; i++ {
       dx[i] = cx[i]
       dy[i] = cy[i]
   }
   for i := 0; i < m; i++ {
       var t, l, d int
       fmt.Fscan(in, &t, &l, &d)
       tmp := l
       if t == 1 {
           pos := findPos(b1, l)
           if b1[pos] == tmp && dx[pos] > d {
               dx[pos] = d
           }
       } else {
           pos := findPos(c1, l)
           if c1[pos] == tmp && dy[pos] > d {
               dy[pos] = d
           }
       }
   }
   S = 2
   T = 2*n + 3
   // initialize graph arrays
   totNodes := T + 2
   head = make([]int, totNodes)
   // allocate enough space for edges: around 20 * totNodes
   capSize := totNodes * 20
   nxt = make([]int, capSize)
   to = make([]int, capSize)
   capArr = make([]int, capSize)
   Sx = make([]int, totNodes)
   xT = make([]int, totNodes)
   h = make([]int, totNodes)
   cur = make([]int, totNodes)
   q = make([]int, totNodes)
   for i := range head {
       head[i] = -1
   }
   // insert x->y edges
   for i := 1; i <= n; i++ {
       eid[i] = insertLower(x[i]+2, y[i]+n+2, 0, 1)
   }
   // x nodes: connect S
   for i := 1; i <= n; i++ {
       c := cx[i]
       dval := dx[i]
       lbound := (c - dval) / 2
       if (c - dval) % 2 != 0 {
           lbound++
       }
       ubound := (c + dval) / 2
       if lbound > ubound {
           fmt.Fprintln(out, -1)
           return
       }
       insertLower(S, i+2, lbound, ubound)
   }
   // y nodes: connect to T
   for i := 1; i <= n; i++ {
       c := cy[i]
       dval := dy[i]
       lbound := (c - dval) / 2
       if (c - dval) % 2 != 0 {
           lbound++
       }
       ubound := (c + dval) / 2
       if lbound > ubound {
           fmt.Fprintln(out, -1)
           return
       }
       insertLower(i+n+2, T, lbound, ubound)
   }
   build()
   dinic()
   // add edge T->S with infinite cap
   idTSBack := insertEdge(T, S, INF)
   dinic()
   // check feasibility: edges from SS
   for i := head[SS]; i != -1; i = nxt[i] {
       if capArr[i] > 0 {
           fmt.Fprintln(out, -1)
           return
       }
   }
   tFlow := capArr[idTSBack]
   res := int64(r)*int64(tFlow) + int64(b)*int64(n-tFlow)
   fmt.Fprintln(out, res)
   // print assignment
   for i := 1; i <= n; i++ {
       used := capArr[eid[i]] // flow on x->y original edge
       var c byte
       if (used == 1) != fl {
           c = 'r'
       } else {
           c = 'b'
       }
       out.WriteByte(c)
   }
   out.WriteByte('\n')
}
