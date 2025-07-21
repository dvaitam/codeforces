package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   succ []int
   parent []int
   diff  []int
   visited []bool
   v     []int
   bval  []int
)

// find with path compression, returns root and distance (sum diff) to root
func find(x int) (int, int) {
   if parent[x] == x {
       return x, 0
   }
   p := parent[x]
   root, d := find(p)
   diff[x] += d
   parent[x] = root
   return parent[x], diff[x]
}

// union root x into root y: make parent[xRoot] = yRoot, set diff
func unite(xRoot, y int) {
   // y may not be root
   yRoot, _ := find(y)
   parent[xRoot] = yRoot
   // one step from xRoot to yRoot (via succ)
   diff[xRoot] = 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   succ = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &succ[i])
   }
   v = make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &v[i])
   }
   bval = make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &bval[i])
   }
   parent = make([]int, n+1)
   diff = make([]int, n+1)
   visited = make([]bool, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       diff[i] = 0
   }
   prevRes := 0
   for i := 0; i < m; i++ {
       ai := ( (v[i] + prevRes - 1) % n ) + 1
       rem := bval[i]
       cur := ai
       cnt := 0
       for rem > 0 {
           root, dist := find(cur)
           if visited[root] {
               break
           }
           if dist >= rem {
               break
           }
           // skip visited nodes (dist steps)
           rem -= dist
           // now at root (unvisited)
           cnt++
           rem--
           visited[root] = true
           // unite root into its successor chain
           unite(root, succ[root])
           cur = succ[root]
       }
       fmt.Fprintln(writer, cnt)
       prevRes = cnt
   }
}
