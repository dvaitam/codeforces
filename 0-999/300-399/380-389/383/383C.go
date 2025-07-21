package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

// Fenwick tree for range update and point query
type Fenwick struct {
   n    int
   bit  []int64
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, bit: make([]int64, n+2)}
}

// add v at index i
func (f *Fenwick) add(i int, v int64) {
   for ; i <= f.n; i += i & -i {
       f.bit[i] += v
   }
}

// range add v on [l..r]
func (f *Fenwick) rangeAdd(l, r int, v int64) {
   if l > r {
       return
   }
   f.add(l, v)
   f.add(r+1, -v)
}

// point query at i
func (f *Fenwick) query(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += f.bit[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   // read ints
   readInt := func() int {
       b, _ := in.ReadBytes(' ')
       // trim space
       s := string(b)
       s = s[:len(s)-1]
       x, _ := strconv.Atoi(s)
       return x
   }
   // better readNonSpace
   read := func() int {
       var x int
       fmt.Fscan(in, &x)
       return x
   }
   n := read()
   m := read()
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       a[i] = int64(read())
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := read()
       v := read()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   depth := make([]int, n+1)
   // Euler tour, iterative DFS
   time := 1
   type st struct{ node, parent, idx int }
   stack := make([]st, 0, n)
   stack = append(stack, st{1, 0, 0})
   depth[1] = 0
   for len(stack) > 0 {
       cur := &stack[len(stack)-1]
       u := cur.node
       if cur.idx == 0 {
           tin[u] = time
           time++
       }
       if cur.idx < len(adj[u]) {
           v := adj[u][cur.idx]
           cur.idx++
           if v == cur.parent {
               continue
           }
           depth[v] = depth[u] ^ 1
           stack = append(stack, st{v, u, 0})
       } else {
           tout[u] = time - 1
           stack = stack[:len(stack)-1]
       }
   }
   // Fenwicks
   fe := NewFenwick(n)
   fo := NewFenwick(n)
   // process queries
   for i := 0; i < m; i++ {
       t := read()
       if t == 1 {
           x := read()
           v := int64(read())
           l, r := tin[x], tout[x]
           if depth[x] == 0 {
               fe.rangeAdd(l, r, v)
               fo.rangeAdd(l, r, -v)
           } else {
               fo.rangeAdd(l, r, v)
               fe.rangeAdd(l, r, -v)
           }
       } else {
           x := read()
           var delta int64
           if depth[x] == 0 {
               delta = fe.query(tin[x])
           } else {
               delta = fo.query(tin[x])
           }
           res := a[x] + delta
           fmt.Fprintln(out, res)
       }
   }
}
