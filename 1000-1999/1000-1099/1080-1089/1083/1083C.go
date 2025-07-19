package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXLOG = 19

var (
   r = bufio.NewReader(os.Stdin)
   w = bufio.NewWriter(os.Stdout)
   n int
   a []int
   where []int
   v [][]int
   tin, tout []int
   gl []int
   par [][]int
   curT int
   tree []Node
   ans Node
   mex int
)

// Node represents two endpoints a, b and their LCA g; a==-1 indicates bad
type Node struct {
   a, b, g int
}

func (nd *Node) bad() bool {
   return nd.a == -1
}

func (nd *Node) beBad() {
   nd.a, nd.b, nd.g = -1, -1, -1
}

// isAncestor checks if x is ancestor of y
func isAncestor(x, y int) bool {
   return tin[x] <= tin[y] && tin[y] <= tout[x]
}

// onTheWay checks if y is on path from x to z
func onTheWay(x, y, z int) bool {
   return isAncestor(x, y) && isAncestor(y, z)
}

func lca(x, y int) int {
   if isAncestor(x, y) {
       return x
   }
   if isAncestor(y, x) {
       return y
   }
   for i := MAXLOG - 1; i >= 0; i-- {
       p := par[x][i]
       if !isAncestor(p, y) {
           x = p
       }
   }
   return par[x][0]
}

// add extends node to include x
func (nd *Node) add(x int) {
   if nd.bad() {
       return
   }
   if nd.a == nd.b {
       nd.b = x
       nd.g = lca(nd.a, nd.b)
       if nd.a == nd.g {
           nd.a, nd.b = nd.b, nd.a
       }
   } else {
       if isAncestor(nd.a, x) {
           nd.a = x
           return
       }
       if nd.b == nd.g {
           if onTheWay(nd.b, x, nd.a) {
               return
           }
           if onTheWay(x, nd.b, nd.a) {
               nd.b = x
               nd.g = x
               return
           }
           ng := lca(nd.a, x)
           if isAncestor(ng, nd.b) {
               nd.b = x
               nd.g = ng
               return
           }
           nd.beBad()
           return
       }
       if isAncestor(nd.b, x) {
           nd.b = x
           return
       }
       if onTheWay(nd.g, x, nd.a) || onTheWay(nd.g, x, nd.b) {
           return
       }
       nd.beBad()
   }
}

func merge(s, t Node) Node {
   if s.bad() || t.bad() {
       return Node{-1, -1, -1}
   }
   res := s
   res.add(t.a)
   res.add(t.b)
   return res
}

// build segment tree over values [l..r]
func build(cur, l, r int) {
   if l == r {
       tree[cur] = Node{where[l], where[l], where[l]}
   } else {
       m := (l + r) >> 1
       lc := cur << 1
       build(lc, l, m)
       build(lc|1, m+1, r)
       tree[cur] = merge(tree[lc], tree[lc|1])
   }
}

// update position pos in segment tree
func update(cur, l, r, pos int) {
   if l == r {
       tree[cur] = Node{where[l], where[l], where[l]}
   } else {
       m := (l + r) >> 1
       lc := cur << 1
       if pos <= m {
           update(lc, l, m, pos)
       } else {
           update(lc|1, m+1, r, pos)
       }
       tree[cur] = merge(tree[lc], tree[lc|1])
   }
}

// getAns finds maximum r where merging nodes [0..r] is valid
func getAns(cur, l, r int) {
   res := merge(ans, tree[cur])
   if !res.bad() {
       mex = r
       ans = res
   } else if l < r {
       m := (l + r) >> 1
       lc := cur << 1
       getAns(lc, l, m)
       if mex == m {
           getAns(lc|1, m+1, r)
       }
   }
}

// dfs computes tin, tout, par, gl
func dfs(cur, p int) {
   curT++
   tin[cur] = curT
   par[cur][0] = p
   for i := 1; i < MAXLOG; i++ {
       par[cur][i] = par[par[cur][i-1]][i-1]
   }
   for _, to := range v[cur] {
       gl[to] = gl[cur] + 1
       dfs(to, cur)
   }
   tout[cur] = curT
}

func main() {
   defer w.Flush()
   fmt.Fscan(r, &n)
   a = make([]int, n)
   where = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(r, &a[i])
       where[a[i]] = i
   }
   v = make([][]int, n)
   for i := 0; i < n-1; i++ {
       var x int
       fmt.Fscan(r, &x)
       v[x-1] = append(v[x-1], i+1)
   }
   tin = make([]int, n)
   tout = make([]int, n)
   gl = make([]int, n)
   par = make([][]int, n)
   for i := range par {
       par[i] = make([]int, MAXLOG)
   }
   dfs(0, 0)
   tree = make([]Node, 4*n)
   build(1, 0, n-1)
   var q int
   fmt.Fscan(r, &q)
   for q > 0 {
       q--
       var t int
       fmt.Fscan(r, &t)
       if t == 1 {
           var x, y int
           fmt.Fscan(r, &x, &y)
           x--
           y--
           a[x], a[y] = a[y], a[x]
           where[a[x]] = x
           where[a[y]] = y
           update(1, 0, n-1, a[x])
           update(1, 0, n-1, a[y])
       } else {
           ans = Node{where[0], where[0], where[0]}
           mex = 0
           getAns(1, 0, n-1)
           fmt.Fprintln(w, mex+1)
       }
   }
}
