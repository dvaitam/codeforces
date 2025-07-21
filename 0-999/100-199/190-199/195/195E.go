package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod int64 = 1000000007

// DSU with weights: weight[u] is weight of edge from u to parent[u]
type DSU struct {
   parent []int
   weight []int64
}

// NewDSU initializes DSU for n elements (1-based)
func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   weight := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       weight[i] = 0
   }
   return &DSU{parent, weight}
}

// find returns root of u and sum of weights from u to root, with path compression
func (d *DSU) find(u int) (int, int64) {
   // collect path
   var nodes []int
   var wts []int64
   cur := u
   for d.parent[cur] != cur {
       nodes = append(nodes, cur)
       wts = append(wts, d.weight[cur])
       cur = d.parent[cur]
   }
   root := cur
   // compute total depth
   total := int64(0)
   for _, w := range wts {
       total += w
   }
   // compress path
   sum := int64(0)
   for i, node := range nodes {
       w := wts[i]
       d.parent[node] = root
       d.weight[node] = total - sum
       sum += w
   }
   return root, total
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   dsu := NewDSU(n)
   var ans int64
   for i := 1; i <= n; i++ {
       var k int
       fmt.Fscan(in, &k)
       for j := 0; j < k; j++ {
           var v int
           var x int64
           fmt.Fscan(in, &v, &x)
           root, depth := dsu.find(v)
           w := depth + x
           // add edge from root to i with weight w
           dsu.parent[root] = i
           dsu.weight[root] = w
           // accumulate answer modulo mod
           wmod := w % mod
           if wmod < 0 {
               wmod += mod
           }
           ans = (ans + wmod) % mod
       }
   }
   fmt.Println(ans)
