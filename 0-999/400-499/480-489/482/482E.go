package main

import (
   "bufio"
   "fmt"
   "os"
)

// Link-Cut Tree based solution translated from C++ solE.cpp
type Node struct {
   c        [2]*Node
   f        *Node
   cnt, sum int64
   val, siz, lazy, treeSize int
}

var (
   null *Node
   nodes []*Node
   adj [][]int
   ans int64
   n, m int
   rdr = bufio.NewReader(os.Stdin)
   wrt = bufio.NewWriter(os.Stdout)
)

func (o *Node) newNode() {
   o.sum = 0; o.val = 0; o.lazy = 0
   o.c[0], o.c[1], o.f = null, null, null
}
func (o *Node) addSize(d int) {
   if o == null { return }
   o.treeSize += d
   o.cnt += 2 * int64(o.siz) * int64(d)
   o.lazy += d
}
func (o *Node) up() {
   o.sum = o.c[0].sum + int64(o.siz)*int64(o.val) + o.c[1].sum
}
func (o *Node) down() {
   if o.lazy != 0 {
       o.c[0].addSize(o.lazy)
       o.c[1].addSize(o.lazy)
       o.lazy = 0
   }
}
func (o *Node) isRoot() bool {
   return o.f == null || (o.f.c[0] != o && o.f.c[1] != o)
}
func (o *Node) signDown() {
   if !o.isRoot() {
       o.f.signDown()
   }
   o.down()
}
func (o *Node) setChild(x *Node, d int) {
   o.c[d] = x; x.f = o
}
// rotate performs a single rotation in link-cut tree
func (o *Node) rotate(dir int) {
   p := o.f; g := p.f
   // move o's child in direction dir^1 to p
   p.setChild(o.c[dir], dir^1)
   // link o to g
   if !p.isRoot() {
       if p == g.c[1] {
           g.setChild(o, 1)
       } else {
           g.setChild(o, 0)
       }
   } else {
       o.f = g
   }
   // put p as o's child in direction dir
   o.setChild(p, dir)
   p.up()
}
func (o *Node) splay() {
   o.signDown()
   for !o.isRoot() {
       if o.f.isRoot() {
           if o == o.f.c[0] {
               o.rotate(1)
           } else {
               o.rotate(0)
           }
       } else {
           p := o.f; g := p.f
           if p == g.c[0] {
               if o == p.c[0] {
                   p.rotate(1)
                   o.rotate(1)
               } else {
                   o.rotate(0)
                   o.rotate(1)
               }
           } else {
               if o == p.c[1] {
                   p.rotate(0)
                   o.rotate(0)
               } else {
                   o.rotate(1)
                   o.rotate(0)
               }
           }
       }
   }
   o.up()
}
// access makes o the root of the preferred path and updates siz
func (o *Node) access() {
   var x *Node = null
   for u := o; u != null; u = u.f {
       if x != null {
           // bring x's leftmost child to root
           y := x
           for y.c[0] != null {
               y = y.c[0]
           }
           y.splay()
           x = y
       }
       u.splay()
       u.siz = u.treeSize - x.treeSize
       u.setChild(x, 1)
       u.up()
       x = u
   }
   o.splay()
}

// DFS to initialize cnt, treeSize and siz
func dfs(u int) {
   node := nodes[u]
   node.cnt = 1; node.treeSize = 1
   for _, v := range adj[u] {
       dfs(v)
       node.cnt += 2 * int64(node.treeSize) * int64(nodes[v].treeSize)
       node.treeSize += nodes[v].treeSize
   }
   ans += node.cnt * int64(node.val)
   node.siz = node.treeSize
}

func readInt() int {
   var x int
   var c byte
   for {
       b, err := rdr.ReadByte()
       if err != nil { break }
       if b >= '0' && b <= '9' {
           c = b; x = int(b - '0'); break
       }
   }
   for {
       b, err := rdr.ReadByte()
       if err != nil || b < '0' || b > '9' { break }
       x = x*10 + int(b - '0')
   }
   return x
}
func readOp() byte {
   for {
       b, err := rdr.ReadByte()
       if err != nil { return 0 }
       if b >= 'A' && b <= 'Z' { return b }
   }
}

func main() {
   defer wrt.Flush()
   n = readInt()
   nodes = make([]*Node, n+1)
   adj = make([][]int, n+1)
   // init null
   null = &Node{}
   null.c[0], null.c[1], null.f = null, null, null
   nodes[0] = null
   for i := 1; i <= n; i++ {
       nodes[i] = &Node{}
       nodes[i].newNode()
   }
   // read tree
   for i := 2; i <= n; i++ {
       p := readInt()
       adj[p] = append(adj[p], i)
       nodes[i].f = nodes[p]
   }
   // read values
   for i := 1; i <= n; i++ {
       nodes[i].val = readInt()
   }
   ans = 0
   dfs(1)
   tot := float64(n) * float64(n)
   fmt.Fprintf(wrt, "%.10f\n", float64(ans)/tot)
   m = readInt()
   for i := 0; i < m; i++ {
       op := readOp()
       x := readInt(); y := readInt()
       if op == 'P' {
           nodes[y].access()
           nodes[x].splay()
           if nodes[x].f == null {
               x, y = y, x
           }
           // cut x
           nodes[x].access()
           o := nodes[x].c[0]
           ans -= 2 * o.sum * int64(nodes[x].treeSize)
           o.addSize(-nodes[x].treeSize)
           nodes[x].c[0].f = null
           nodes[x].c[0] = null
           nodes[x].up()
           o.up()
           // link to y
           nodes[y].access()
           ans += 2 * nodes[y].sum * int64(nodes[x].treeSize)
           nodes[y].addSize(nodes[x].treeSize)
           nodes[y].up()
           nodes[x].f = nodes[y]
       } else {
           nodes[x].access()
           ans += nodes[x].cnt * int64(y - nodes[x].val)
           nodes[x].val = y
           nodes[x].up()
       }
       fmt.Fprintf(wrt, "%.10f\n", float64(ans)/tot)
   }
}
