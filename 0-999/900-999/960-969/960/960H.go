package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   MAXN = 50005
   MAXM = 15000005
)

var (
   n, m, q, C int
   deep [MAXN]int
   top  [MAXN]int
   fa   [MAXN]int
   sizeArr [MAXN]int
   son  [MAXN]int
   st   [MAXN]int
   ed   [MAXN]int
   colArr [MAXN]int
   costArr [MAXN]int
   depthArr [MAXN]int64
   typArr   [MAXN]int64
   v   [][]int
   root []int
   tree []Node
   dfsTime  int
   totNodes int
)

// Node for segment tree
type Node struct {
   l, r int
   sum, add int64
}

func dfs(u int) {
   sizeArr[u] = 0
   for _, vtx := range v[u] {
       deep[vtx] = deep[u] + 1
       fa[vtx] = u
       dfs(vtx)
       sizeArr[u] += sizeArr[vtx]
       if sizeArr[vtx] > sizeArr[son[u]] {
           son[u] = vtx
       }
   }
   sizeArr[u]++
}

func dfs2(u, tp int) {
   dfsTime++
   st[u] = dfsTime
   top[u] = tp
   if son[u] != 0 {
       dfs2(son[u], tp)
   }
   for _, vtx := range v[u] {
       if vtx == son[u] {
           continue
       }
       dfs2(vtx, vtx)
   }
   ed[u] = dfsTime
}

func down(t, length int) {
   addv := tree[t].add
   if addv == 0 {
       return
   }
   if tree[t].l == 0 {
       totNodes++
       tree[t].l = totNodes
   }
   if tree[t].r == 0 {
       totNodes++
       tree[t].r = totNodes
   }
   lc := tree[t].l
   rc := tree[t].r
   tree[lc].sum += addv * int64(length-(length>>1))
   tree[rc].sum += addv * int64(length>>1)
   tree[lc].add += addv
   tree[rc].add += addv
   tree[t].add = 0
}

func query(ll, rr, l, r, t int) int64 {
   if t == 0 {
       return 0
   }
   if ll <= l && r <= rr {
       return tree[t].sum
   }
   down(t, r-l+1)
   mid := (l + r) >> 1
   var res int64
   if ll <= mid {
       res += query(ll, rr, l, mid, tree[t].l)
   }
   if rr > mid {
       res += query(ll, rr, mid+1, r, tree[t].r)
   }
   return res
}

func modify(ll, rr, c, l, r, t int) {
   if ll <= l && r <= rr {
       tree[t].sum += int64(c) * int64(r-l+1)
       tree[t].add += int64(c)
       return
   }
   down(t, r-l+1)
   mid := (l + r) >> 1
   if ll <= mid {
       if tree[t].l == 0 {
           totNodes++
           tree[t].l = totNodes
       }
       modify(ll, rr, c, l, mid, tree[t].l)
   }
   if rr > mid {
       if tree[t].r == 0 {
           totNodes++
           tree[t].r = totNodes
       }
       modify(ll, rr, c, mid+1, r, tree[t].r)
   }
   tree[t].sum = tree[tree[t].l].sum + tree[tree[t].r].sum
}

func ask(rt, u int) int64 {
   var s int64
   for u != 0 {
       s += query(st[top[u]], st[u], 1, n, rt)
       u = fa[top[u]]
   }
   return s
}

func addPath(rt, u, delta int) {
   for u != 0 {
       modify(st[top[u]], st[u], delta, 1, n, rt)
       u = fa[top[u]]
   }
}

func ins(color, u int) {
   depthArr[color] += int64(deep[u])
   typArr[color] += 2 * ask(root[color], u)
   typArr[color] += int64(deep[u])
   addPath(root[color], u, 1)
}

func del(color, u int) {
   depthArr[color] -= int64(deep[u])
   addPath(root[color], u, -1)
   typArr[color] -= 2 * ask(root[color], u)
   typArr[color] -= int64(deep[u])
}

func work(u int) float64 {
   a := float64(typArr[u]) * float64(costArr[u]) * float64(costArr[u])
   b := 2.0 * float64(depthArr[u]) * float64(costArr[u]) * float64(C)
   c := float64(n) * float64(C) * float64(C)
   return (a - b + c) / float64(n)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m, &q, &C)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &colArr[i])
   }
   v = make([][]int, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &fa[i])
       v[fa[i]] = append(v[fa[i]], i)
   }
   for i := 1; i <= m; i++ {
       fmt.Fscan(reader, &costArr[i])
   }
   deep[1] = 1
   dfs(1)
   dfsTime = 0
   dfs2(1, 1)
   // initialize segment tree roots
   tree = make([]Node, MAXM)
   root = make([]int, m+1)
   totNodes = 0
   for i := 1; i <= m; i++ {
       totNodes++
       root[i] = totNodes
   }
   for u := 1; u <= n; u++ {
       ins(colArr[u], u)
   }
   for ; q > 0; q-- {
       var ty int
       fmt.Fscan(reader, &ty)
       if ty == 1 {
           var u, newc int
           fmt.Fscan(reader, &u, &newc)
           del(colArr[u], u)
           colArr[u] = newc
           ins(colArr[u], u)
       } else {
           var u int
           fmt.Fscan(reader, &u)
           res := work(u)
           fmt.Fprintf(writer, "%.10f\n", res)
       }
   }
}
