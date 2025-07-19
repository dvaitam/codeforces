package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MAXN = 200005

type Edge struct{ to, idx int }
type Node struct {
   C, cnt int
   S, sum  int64
   tag     bool
}

func (nd *Node) flip() {
   nd.cnt = nd.C - nd.cnt
   nd.sum = nd.S - nd.sum
   nd.tag = !nd.tag
}

func merge(a, b Node) Node {
   return Node{
       C:   a.C + b.C,
       cnt: a.cnt + b.cnt,
       S:   a.S + b.S,
       sum: a.sum + b.sum,
   }
}

var (
   n       int
   E       [MAXN][]Edge
   siz, son [MAXN]int
   top, L, R [MAXN]int
   F, dep   [MAXN]int
   dw, rev  [MAXN]int
   dfncnt   int
   lf, yes  [MAXN]bool
   tree    []Node
   ans      []int
   Ncount   int
   ok       bool
   rdr      = bufio.NewReader(os.Stdin)
   wrt      = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   var c byte
   var neg bool
   for {
       b, err := rdr.ReadByte()
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
       b, _ := rdr.ReadByte()
       c = b
   }
   for ; c >= '0' && c <= '9'; b, _ := rdr.ReadByte() {
       x = x*10 + int(c-'0')
       c = b
   }
   if neg {
       x = -x
   }
   return x
}

func dfs1(u, fa int) {
   F[u] = fa
   dep[u] = dep[fa] + 1
   siz[u] = 1
   son[u] = 0
   for _, e := range E[u] {
       v := e.to
       if v == fa {
           continue
       }
       dw[v] = e.idx
       dfs1(v, u)
       siz[u] += siz[v]
       if siz[v] > siz[son[u]] {
           son[u] = v
       }
   }
}

func dfs2(u, tp int) {
   top[u] = tp
   dfncnt++
   L[u] = dfncnt
   rev[dfncnt] = u
   if son[u] == 0 {
       R[u] = dfncnt
       return
   }
   dfs2(son[u], tp)
   for _, e := range E[u] {
       v := e.to
       if v != F[u] && v != son[u] {
           dfs2(v, v)
       }
   }
   R[u] = dfncnt
}

func pushup(u int) {
   tree[u] = merge(tree[u<<1], tree[u<<1|1])
}

func pushdown(u int) {
   if tree[u].tag {
       tree[u<<1].flip()
       tree[u<<1|1].flip()
       tree[u].tag = false
   }
}

func upd(u, l, r, pos int) {
   if l == r {
       tree[u].C = 1
       tree[u].S = int64(dw[rev[l]])
       tree[u].cnt = 1
       tree[u].sum = tree[u].S
       return
   }
   pushdown(u)
   mid := (l + r) >> 1
   if pos <= mid {
       upd(u<<1, l, mid, pos)
   } else {
       upd(u<<1|1, mid+1, r, pos)
   }
   pushup(u)
}

func flipRange(u, l, r, ql, qr int) {
   if ql <= l && r <= qr {
       tree[u].flip()
       return
   }
   pushdown(u)
   mid := (l + r) >> 1
   if ql <= mid {
       flipRange(u<<1, l, mid, ql, qr)
   }
   if qr > mid {
       flipRange(u<<1|1, mid+1, r, ql, qr)
   }
   pushup(u)
}

func dfs3(u, fa int) int {
   if lf[u] {
       ans = append(ans, dw[u])
       return 1
   }
   ret := 0
   for _, e := range E[u] {
       v := e.to
       if yes[v] && v != fa {
           ret += dfs3(v, u)
       }
   }
   if ret == 0 {
       ans = append(ans, dw[u])
       return 1
   }
   return 0
}

func solve() {
   ans = ans[:0]
   dfs3(1, 0)
   sort.Ints(ans)
   fmt.Fprint(wrt, len(ans), ' ')
   for i, v := range ans {
       if i > 0 {
           fmt.Fprint(wrt, ' ')
       }
       fmt.Fprint(wrt, v)
   }
   fmt.Fprint(wrt, '\n')
}

func main() {
   defer wrt.Flush()
   n = readInt()
   for i := 1; i < n; i++ {
       u := readInt()
       v := readInt()
       E[u] = append(E[u], Edge{v, i})
       E[v] = append(E[v], Edge{u, i})
   }
   dfs1(1, 0)
   dfs2(1, 1)
   tree = make([]Node, 4*(n+1))
   lf[1] = true
   upd(1, 1, n, 1)
   Ncount = 1
   for {
       opt := readInt()
       if opt == 1 {
           x := readInt()
           upd(1, 1, n, L[x])
           Ncount++
           lf[F[x]] = false
           lf[x] = true
           yes[x] = true
           y := F[x]
           for top[y] != 1 {
               flipRange(1, 1, n, L[top[y]], L[y])
               y = F[top[y]]
           }
           flipRange(1, 1, n, 1, L[y])
           ok = (tree[1].cnt*2 == Ncount)
           if ok {
               fmt.Fprint(wrt, tree[1].sum, '\n')
           } else {
               fmt.Fprint(wrt, 0, '\n')
           }
       } else if opt == 2 {
           if !ok {
               fmt.Fprint(wrt, 0, '\n')
           } else {
               solve()
           }
       } else {
           break
       }
   }
}
