package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU with rollback
type DSU struct {
   parent []int
   size   []int
   ops    []op
}
type op struct{ vRoot, uRoot, oldSize int }

func newDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent: parent, size: size, ops: make([]op, 0, 1024)}
}
func (d *DSU) find(x int) int {
   for d.parent[x] != x {
       x = d.parent[x]
   }
   return x
}
func (d *DSU) unite(u, v int) {
   ru := d.find(u)
   rv := d.find(v)
   if ru == rv {
       d.ops = append(d.ops, op{-1, 0, 0})
       return
   }
   if d.size[ru] < d.size[rv] {
       ru, rv = rv, ru
   }
   // attach rv under ru
   d.ops = append(d.ops, op{rv, ru, d.size[ru]})
   d.parent[rv] = ru
   d.size[ru] += d.size[rv]
}
func (d *DSU) rollback(to int) {
   for len(d.ops) > to {
       last := d.ops[len(d.ops)-1]
       d.ops = d.ops[:len(d.ops)-1]
       if last.vRoot == -1 {
           continue
       }
       // restore
       d.size[last.uRoot] = last.oldSize
       d.parent[last.vRoot] = last.vRoot
   }
}

// edge interval
type edge struct{ u, v int }

var (
   q         int
   t         []int
   x, y, z   []int
   dayNumber []int
   dayStart  []int
   answers   []int
   seg       [][]edge
   dsu       *DSU
)

func addEdge(pos, l, r, ql, qr, u, v int) {
   if ql >= r || qr <= l {
       return
   }
   if ql <= l && r <= qr {
       seg[pos] = append(seg[pos], edge{u, v})
       return
   }
   m := (l + r) >> 1
   addEdge(pos<<1, l, m, ql, qr, u, v)
   addEdge(pos<<1|1, m, r, ql, qr, u, v)
}

func dfs(pos, l, r int) {
   opsBefore := len(dsu.ops)
   for _, e := range seg[pos] {
       dsu.unite(e.u, e.v)
   }
   if l+1 == r {
       if l <= q && t[l] == 2 {
           root := dsu.find(z[l])
           answers[l] = dsu.size[root]
       }
   } else {
       m := (l + r) >> 1
       dfs(pos<<1, l, m)
       dfs(pos<<1|1, m, r)
   }
   dsu.rollback(opsBefore)
}

func readInt(r *bufio.Reader) int {
   c := make([]byte, 1)
   // skip non-digit
   for {
       if _, err := r.Read(c); err != nil {
           return 0
       }
       if c[0] >= '0' && c[0] <= '9' {
           break
       }
   }
   x := int(c[0] - '0')
   for {
       if _, err := r.Read(c); err != nil {
           break
       }
       if c[0] < '0' || c[0] > '9' {
           break
       }
       x = x*10 + int(c[0]-'0')
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   n := readInt(in)
   q = readInt(in)
   k := readInt(in)
   t = make([]int, q+2)
   x = make([]int, q+2)
   y = make([]int, q+2)
   z = make([]int, q+2)
   dayNumber = make([]int, q+2)
   dayStart = make([]int, q+3)
   answers = make([]int, q+2)
   // read queries and days
   day := 1
   dayStart[1] = 1
   for i := 1; i <= q; i++ {
       ti := readInt(in)
       t[i] = ti
       dayNumber[i] = day
       if ti == 1 {
           x[i] = readInt(in)
           y[i] = readInt(in)
       } else if ti == 2 {
           z[i] = readInt(in)
       } else if ti == 3 {
           // end of day
           day++
           if i+1 <= q {
               dayStart[day] = i + 1
           }
       }
   }
   maxDay := day
   // build segment tree
   size := 1
   for size < q+2 {
       size <<= 1
   }
   seg = make([][]edge, size<<1)
   // add edges intervals
   for i := 1; i <= q; i++ {
       if t[i] == 1 {
           d := dayNumber[i]
           dExpDay := d + k
           var r int
           if dExpDay+1 <= maxDay && dayStart[dExpDay+1] != 0 {
               r = dayStart[dExpDay+1]
           } else {
               r = q + 1
           }
           l := i
           if l < r {
               addEdge(1, 1, q+1, l, r, x[i], y[i])
           }
       }
   }
   // init DSU
   dsu = newDSU(n)
   // process
   dfs(1, 1, q+1)
   // output answers
   for i := 1; i <= q; i++ {
       if t[i] == 2 {
           fmt.Fprintln(out, answers[i])
       }
   }
}
