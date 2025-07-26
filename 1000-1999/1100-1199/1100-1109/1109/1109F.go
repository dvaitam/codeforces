package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxN = 200005

type Point struct{ x, y int }

var (
   r, c, n int
   f [][]int
   pos [maxN]Point
   eg  [maxN][]int
   nodes [maxN]Node
   st []*Node
   segMin [4 * maxN]int
   segCnt [4 * maxN]int
   segFg  [4 * maxN]int
)

// Link-Cut Tree Node
type Node struct {
   s   [2]*Node
   f   *Node
   rev bool
}

func (x *Node) isRoot() bool {
   return x.f == nil || (x.f.s[0] != x && x.f.s[1] != x)
}
func (x *Node) dir() int {
   if x.f != nil && x.f.s[1] == x {
       return 1
   }
   return 0
}
func (x *Node) setChild(c *Node, d int) {
   x.s[d] = c
   if c != nil {
       c.f = x
   }
}
func (x *Node) push() {
   if x.rev {
       x.s[0], x.s[1] = x.s[1], x.s[0]
       if x.s[0] != nil {
           x.s[0].rev = !x.s[0].rev
       }
       if x.s[1] != nil {
           x.s[1].rev = !x.s[1].rev
       }
       x.rev = false
   }
}
func (x *Node) update() {}

func rotate(x *Node) {
   p := x.f
   d := x.dir()
   if !p.isRoot() {
       p.f.setChild(x, p.dir())
   } else {
       x.f = p.f
   }
   // toggle child based on direction
   p.setChild(x.s[1-d], d)
   x.setChild(p, 1-d)
   p.update()
}

func splay(x *Node) {
   st = st[:0]
   y := x
   for {
       st = append(st, y)
       if y.isRoot() {
           break
       }
       y = y.f
   }
   for i := len(st) - 1; i >= 0; i-- {
       st[i].push()
   }
   for !x.isRoot() {
       if x.f.isRoot() {
           rotate(x)
       } else if x.dir() == x.f.dir() {
           rotate(x.f)
           rotate(x)
       } else {
           rotate(x)
           rotate(x)
       }
   }
   x.update()
}

func expose(x *Node) *Node {
   var last *Node
   for y := x; y != nil; y = y.f {
       splay(y)
       y.s[1] = last
       y.update()
       last = y
   }
   return last
}

func findRoot(x *Node) *Node {
   expose(x)
   splay(x)
   for x.s[0] != nil {
       x.push()
       x = x.s[0]
   }
   splay(x)
   return x
}

func makeRoot(x *Node) {
   expose(x)
   splay(x)
   x.rev = !x.rev
   x.push()
}

func link(x, y *Node) {
   makeRoot(x)
   // makeRoot(y)
   x.f = y
}

func cut(x, y *Node) {
   makeRoot(x)
   expose(y)
   splay(y)
   if y.s[0] == x {
       y.s[0].f = nil
       y.s[0] = nil
   }
}

// Segment Tree funcs
func segUpdate(p int) {
   if segMin[p*2] < segMin[p*2+1] {
       segMin[p] = segMin[p*2]
       segCnt[p] = segCnt[p*2]
   } else if segMin[p*2] > segMin[p*2+1] {
       segMin[p] = segMin[p*2+1]
       segCnt[p] = segCnt[p*2+1]
   } else {
       segMin[p] = segMin[p*2]
       segCnt[p] = segCnt[p*2] + segCnt[p*2+1]
   }
}

func segApply(p, v int) {
   segMin[p] += v
   segFg[p] += v
}

func segPush(p int) {
   if segFg[p] != 0 {
       segApply(p*2, segFg[p])
       segApply(p*2+1, segFg[p])
       segFg[p] = 0
   }
}

func segBuild(p, l, r int) {
   segFg[p] = 0
   if l == r {
       segMin[p] = 0
       segCnt[p] = 1
   } else {
       m := (l + r) >> 1
       segBuild(p*2, l, m)
       segBuild(p*2+1, m+1, r)
       segUpdate(p)
   }
}

func segModify(p, l, r, tl, tr, v int) {
   if tl > tr {
       return
   }
   if tl <= l && r <= tr {
       segApply(p, v)
       return
   }
   segPush(p)
   m := (l + r) >> 1
   if tr <= m {
       segModify(p*2, l, m, tl, tr, v)
   } else if tl > m {
       segModify(p*2+1, m+1, r, tl, tr, v)
   } else {
       segModify(p*2, l, m, tl, m, v)
       segModify(p*2+1, m+1, r, m+1, tr, v)
   }
   segUpdate(p)
}

func segQuery(p, l, r, tl, tr int) (int, int) {
   if tl <= l && r <= tr {
       return segMin[p], segCnt[p]
   }
   segPush(p)
   m := (l + r) >> 1
   if tr <= m {
       return segQuery(p*2, l, m, tl, tr)
   } else if tl > m {
       return segQuery(p*2+1, m+1, r, tl, tr)
   }
   lmin, lcnt := segQuery(p*2, l, m, tl, m)
   rmin, rcnt := segQuery(p*2+1, m+1, r, m+1, tr)
   if lmin < rmin {
       return lmin, lcnt
   } else if rmin < lmin {
       return rmin, rcnt
   }
   return lmin, lcnt + rcnt
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &r, &c)
   n = r * c
   f = make([][]int, r+2)
   for i := range f {
       f[i] = make([]int, c+2)
   }
   for i := 1; i <= r; i++ {
       for j := 1; j <= c; j++ {
           fmt.Fscan(in, &f[i][j])
           pos[f[i][j]] = Point{i, j}
       }
   }
   st = make([]*Node, 0, n+5)
   segBuild(1, 1, n)
   pr := 0
   var ans int64
   dx := [4]int{0, 1, 0, -1}
   dy := [4]int{1, 0, -1, 0}
   for i := 1; i <= n; i++ {
       if pr < i {
           pr = i
           segModify(1, 1, n, i, n, 1)
       }
       for pr < n {
           z := pr + 1
           p := pos[z]
           ta := make([]int, 0, 4)
           cc := make([]*Node, 0, 4)
           for k := 0; k < 4; k++ {
               nx, ny := p.x+dx[k], p.y+dy[k]
               if nx >= 1 && nx <= r && ny >= 1 && ny <= c {
                   v := f[nx][ny]
                   if v >= i && v < z {
                       ta = append(ta, v)
                       cc = append(cc, &nodes[v])
                   }
               }
           }
           ok := true
           mcc := len(cc)
           for a := 0; a < mcc; a++ {
               for b := a + 1; b < mcc; b++ {
                   if findRoot(cc[a]) == findRoot(cc[b]) {
                       ok = false
                       break
                   }
               }
               if !ok {
                   break
               }
           }
           if !ok {
               break
           }
           pr++
           segModify(1, 1, n, pr, n, 1)
           for _, v := range ta {
               link(&nodes[pr], &nodes[v])
               eg[v] = append(eg[v], pr)
               segModify(1, 1, n, pr, n, -1)
           }
       }
       mn, cnt := segQuery(1, 1, n, i, pr)
       if mn == 1 {
           ans += int64(cnt)
       }
       segModify(1, 1, n, i, n, -1)
       for _, v := range eg[i] {
           segModify(1, 1, n, v, n, 1)
           cut(&nodes[i], &nodes[v])
       }
   }
   fmt.Fprintln(out, ans)
}
