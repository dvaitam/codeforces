package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() (int, error) {
   var x int
   _, err := fmt.Fscan(reader, &x)
   return x, err
}

// BIT supports range add, point query
type BIT struct {
   n    int
   tree []int64
}

func newBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int64, n+2)}
}

func (b *BIT) add(idx int, v int64) {
   for i := idx; i <= b.n; i += i & -i {
       b.tree[i] += v
   }
}

// range add [l,r] by v
func (b *BIT) rangeAdd(l, r int, v int64) {
   if l <= r {
       b.add(l, v)
       if r+1 <= b.n {
           b.add(r+1, -v)
       }
   }
}

// point query at idx
func (b *BIT) pointQuery(idx int) int64 {
   var res int64
   for i := idx; i > 0; i -= i & -i {
       res += b.tree[i]
   }
   return res
}

func main() {
   defer writer.Flush()
   n, _ := readInt()
   // build tree
   adj := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       u, _ := readInt()
       v, _ := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // init arrays
   dep := make([]int, n+1)
   siz := make([]int, n+1)
   id := make([]int, n+1)
   parent := make([]int, n+1)
   vec := make([][]int, n+2)
   tot := 0
   maxDep := 0
   // dfs iterative
   type item struct{ u, p, d, state int }
   stack := make([]item, 0, n*2)
   stack = append(stack, item{1, 0, 1, 0})
   for len(stack) > 0 {
       it := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, p, d, st := it.u, it.p, it.d, it.state
       if st == 0 {
           dep[u] = d
           if d > maxDep {
               maxDep = d
           }
           tot++
           id[u] = tot
           parent[u] = p
           vec[d] = append(vec[d], u)
           // post state
           stack = append(stack, item{u, p, d, 1})
           // push children
           for _, v := range adj[u] {
               if v == p {
                   continue
               }
               stack = append(stack, item{v, u, d + 1, 0})
           }
       } else {
           // post-order: compute size
           s := 1
           for _, v := range adj[u] {
               if v == p {
                   continue
               }
               s += siz[v]
           }
           siz[u] = s
       }
   }
   m, _ := readInt()
   type Op struct{ v, temp int; x int64 }
   ops := make([]Op, m)
   for i := 0; i < m; i++ {
       v, _ := readInt()
       d, _ := readInt()
       x, _ := readInt()
       // limit d so temp <= maxDep
       if d > maxDep-dep[v] {
           d = maxDep - dep[v]
       }
       temp := dep[v] + d
       ops[i] = Op{v: v, temp: temp, x: int64(x)}
   }
   sort.Slice(ops, func(i, j int) bool {
       return ops[i].temp > ops[j].temp
   })
   bit := newBIT(n)
   ans := make([]int64, n+1)
   now := maxDep
   for _, op := range ops {
       for now > op.temp {
           for _, u := range vec[now] {
               ans[u] = bit.pointQuery(id[u])
           }
           now--
       }
       // apply
       l := id[op.v]
       r := id[op.v] + siz[op.v] - 1
       bit.rangeAdd(l, r, op.x)
   }
   for now > 0 {
       for _, u := range vec[now] {
           ans[u] = bit.pointQuery(id[u])
       }
       now--
   }
   // output
   for i := 1; i <= n; i++ {
       fmt.Fprintf(writer, "%d", ans[i])
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
