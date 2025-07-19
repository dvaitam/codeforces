package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const (
   N   = 100005
   LOG = 18
)

// Link-Cut Tree arrays
var ch [N][2]int
var fa [N]int
var sz [N]int

func isRoot(x int) bool {
   p := fa[x]
   return ch[p][0] != x && ch[p][1] != x
}
func get(x int) int {
   if ch[fa[x]][1] == x {
       return 1
   }
   return 0
}
func pushUp(x int) {
   sz[x] = sz[ch[x][0]] + sz[ch[x][1]] + 1
}
func rotate(x int) {
   y := fa[x]
   z := fa[y]
   k := get(x)
   if !isRoot(y) {
       ch[z][get(y)] = x
   }
   // detach
   c := ch[x][1-k]
   if c != 0 {
       fa[c] = y
   }
   ch[y][k] = c
   ch[x][1-k] = y
   fa[y] = x
   fa[x] = z
   pushUp(y)
   pushUp(x)
}
func splay(x int) {
   for !isRoot(x) {
       if !isRoot(fa[x]) && get(x) == get(fa[x]) {
           rotate(fa[x])
       }
       rotate(x)
   }
}
func access(x int) {
   var p int
   for x != 0 {
       splay(x)
       ch[x][1] = p
       pushUp(x)
       p = x
       x = fa[x]
   }
}
func findRoot(x int) int {
   access(x)
   splay(x)
   for ch[x][0] != 0 {
       x = ch[x][0]
   }
   splay(x)
   return x
}
func link(x, y int) {
   splay(x)
   fa[x] = y
}
func cut(x int) {
   access(x)
   splay(x)
   left := ch[x][0]
   if left != 0 {
       fa[left] = 0
       ch[x][0] = 0
   }
   pushUp(x)
}
func ask(x int) int {
   access(x)
   splay(x)
   return sz[x]
}

// RMQ arrays
var a [N]int
var logN [N]int
var mx [LOG][N]int
var mn [LOG][N]int

func initST(n int) {
   logN[1] = 0
   for i := 2; i <= n; i++ {
       logN[i] = logN[i/2] + 1
   }
   for i := 1; i <= n; i++ {
       mx[0][i] = a[i]
       mn[0][i] = a[i]
   }
   for j := 1; j < LOG; j++ {
       for i := 1; i+(1<<j)-1 <= n; i++ {
           v1, v2 := mx[j-1][i], mx[j-1][i+(1<<(j-1))]
           if v1 > v2 {
               mx[j][i] = v1
           } else {
               mx[j][i] = v2
           }
           u1, u2 := mn[j-1][i], mn[j-1][i+(1<<(j-1))]
           if u1 < u2 {
               mn[j][i] = u1
           } else {
               mn[j][i] = u2
           }
       }
   }
}
func qry(l, r int) int {
   k := logN[r-l+1]
   x1, x2 := mx[k][l], mx[k][r-(1<<k)+1]
   y1, y2 := mn[k][l], mn[k][r-(1<<k)+1]
   maxv := x1
   if x2 > maxv {
       maxv = x2
   }
   minv := y1
   if y2 < minv {
       minv = y2
   }
   return maxv - minv
}
func fnd(pos, w, n int) int {
   l, r := pos, n+1
   for l+1 < r {
       mid := (l + r) >> 1
       if qry(pos, mid) > w {
           r = mid
       } else {
           l = mid
       }
   }
   return r
}

// Query and add structures
type Query struct {
   val, idx int
}
type Pair struct {
   x, y int
}

func lowerBound(c []Query, v int) int {
   return sort.Search(len(c), func(i int) bool { return c[i].val >= v })
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, w, q int
   fmt.Fscan(reader, &n, &w, &q)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
       sz[i] = 1
       fa[i] = 0
       ch[i][0], ch[i][1] = 0, 0
   }
   initST(n)
   c := make([]Query, q)
   for i := 0; i < q; i++ {
       var ki int
       fmt.Fscan(reader, &ki)
       c[i] = Query{val: w - ki, idx: i}
   }
   sort.Slice(c, func(i, j int) bool {
       if c[i].val != c[j].val {
           return c[i].val < c[j].val
       }
       return c[i].idx < c[j].idx
   })
   b := int(math.Min(100.0, math.Sqrt(float64(n))+0.5))
   add := make([][]Pair, q+1)
   for i := 1; i <= n; i++ {
       dif := 0
       nxt := fnd(i, dif, n)
       limit := i + b
       if limit > n {
           limit = n
       }
       for nxt <= limit {
           pos := lowerBound(c, dif)
           add[pos] = append(add[pos], Pair{x: i, y: nxt})
           dif = qry(i, nxt)
           nxt = fnd(i, dif, n)
       }
       pos := lowerBound(c, dif)
       add[pos] = append(add[pos], Pair{x: i, y: 0})
   }
   res := make([]int, q)
   // process queries
   for i := 0; i < q; i++ {
       for _, p := range add[i] {
           cut(p.x)
           if p.y != 0 {
               link(p.x, p.y)
           }
       }
       cur, ans := 1, 0
       for cur <= n {
           ans += ask(cur)
           cur = findRoot(cur)
           cur = fnd(cur, c[i].val, n)
       }
       res[c[i].idx] = ans - 1
   }
   for i := 0; i < q; i++ {
       fmt.Fprintln(writer, res[i])
   }
}
