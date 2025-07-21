package main

import (
   "bufio"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var c byte
   var x int
   // skip non-numbers
   for {
       b, err := reader.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c >= '0' && c <= '9' {
           break
       }
   }
   x = int(c - '0')
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
   return x
}

func main() {
   defer writer.Flush()
   n := readInt()
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // HLD prep
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   size := make([]int, n+1)
   heavy := make([]int, n+1)
   // dfs1: compute parent, depth, order
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   parent[1] = 0
   depth[1] = 0
   stack = append(stack, 1)
   for len(stack) > 0 {
       v := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, v)
       for _, u := range adj[v] {
           if u != parent[v] {
               parent[u] = v
               depth[u] = depth[v] + 1
               stack = append(stack, u)
           }
       }
   }
   // compute sizes and heavy child
   for i := n - 1; i >= 0; i-- {
       v := order[i]
       size[v] = 1
       maxSize := 0
       for _, u := range adj[v] {
           if parent[u] == v {
               size[v] += size[u]
               if size[u] > maxSize {
                   maxSize = size[u]
                   heavy[v] = u
               }
           }
       }
   }
   // decompose
   head := make([]int, n+1)
   pos := make([]int, n+1)
   curPos := 0
   type task struct{ v, h int }
   tstack := []task{{1, 1}}
   for len(tstack) > 0 {
       t := tstack[len(tstack)-1]
       tstack = tstack[:len(tstack)-1]
       v, h := t.v, t.h
       // follow heavy path
       for c := v; c != 0; c = heavy[c] {
           head[c] = h
           pos[c] = curPos
           curPos++
           // push light children
           for _, u := range adj[c] {
               if parent[u] == c && u != heavy[c] {
                   tstack = append(tstack, task{u, u})
               }
           }
       }
   }
   // segment tree for fill and empty times
   st := newSegTree(n)
   q := readInt()
   time := 1
   for i := 0; i < q; i++ {
       op := readInt()
       v := readInt()
       switch op {
       case 1:
           // fill subtree
           l := pos[v]
           r := pos[v] + size[v] - 1
           st.updateFill(l, r, time)
       case 2:
           // empty path to root
           for x := v; x != 0; x = parent[head[x]] {
               l := pos[head[x]]
               r := pos[x]
               st.updateEmpty(l, r, time)
           }
       case 3:
           // query
           idx := pos[v]
           ft, et := st.query(idx)
           if ft > et {
               writer.WriteByte('1')
           } else {
               writer.WriteByte('0')
           }
           writer.WriteByte('\n')
       }
       time++
   }
}

// segment tree supporting range assign for fill and empty times
type segTree struct{
   n, size int
   fillT, fillL []int
   empT, empL   []int
}

func newSegTree(n int) *segTree {
   s := &segTree{n: n}
   sz := 1
   for sz < n {
       sz <<= 1
   }
   s.size = sz
   s.fillT = make([]int, 2*sz)
   s.fillL = make([]int, 2*sz)
   s.empT  = make([]int, 2*sz)
   s.empL  = make([]int, 2*sz)
   for i := range s.fillL {
       s.fillL[i] = -1
       s.empL[i] = -1
   }
   return s
}

func (s *segTree) push(x int) {
   if s.fillL[x] != -1 {
       v := s.fillL[x]
       lc, rc := 2*x, 2*x+1
       s.fillT[lc], s.fillL[lc] = v, v
       s.fillT[rc], s.fillL[rc] = v, v
       s.fillL[x] = -1
   }
   if s.empL[x] != -1 {
       v := s.empL[x]
       lc, rc := 2*x, 2*x+1
       s.empT[lc], s.empL[lc] = v, v
       s.empT[rc], s.empL[rc] = v, v
       s.empL[x] = -1
   }
}

func (s *segTree) updateFill(l, r, v int) {
   s.uf(1, 0, s.size-1, l, r, v)
}
func (s *segTree) uf(x, lx, rx, l, r, v int) {
   if l > rx || r < lx {
       return
   }
   if l <= lx && rx <= r {
       s.fillT[x], s.fillL[x] = v, v
       return
   }
   s.push(x)
   m := (lx + rx) >> 1
   s.uf(2*x, lx, m, l, r, v)
   s.uf(2*x+1, m+1, rx, l, r, v)
}

func (s *segTree) updateEmpty(l, r, v int) {
   s.ue(1, 0, s.size-1, l, r, v)
}
func (s *segTree) ue(x, lx, rx, l, r, v int) {
   if l > rx || r < lx {
       return
   }
   if l <= lx && rx <= r {
       s.empT[x], s.empL[x] = v, v
       return
   }
   s.push(x)
   m := (lx + rx) >> 1
   s.ue(2*x, lx, m, l, r, v)
   s.ue(2*x+1, m+1, rx, l, r, v)
}

func (s *segTree) query(i int) (int, int) {
   return s.q(1, 0, s.size-1, i)
}
func (s *segTree) q(x, lx, rx, i int) (int, int) {
   if lx == rx {
       return s.fillT[x], s.empT[x]
   }
   s.push(x)
   m := (lx + rx) >> 1
   if i <= m {
       return s.q(2*x, lx, m, i)
   }
   return s.q(2*x+1, m+1, rx, i)
}
