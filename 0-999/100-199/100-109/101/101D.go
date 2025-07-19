package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Edge represents a corridor to another hall with a traversal time
type Edge struct {
   to, v int
}
// Arr holds subtree size and total edge-sum for sorting
type Arr struct {
   size int
   sum   int64
}

var (
   n        int
   adj      [][]Edge
   sizeArr  []int
   sumArr   []int64
   timeArr  []int64
)

// dfs computes sizeArr, sumArr, and timeArr for subtree rooted at x
func dfs(x, fa int) {
   sizeArr[x] = 1
   // first pass: compute subtree sums and sizes
   for _, e := range adj[x] {
       y, v := e.to, e.v
       if y == fa {
           continue
       }
       sumArr[y] = int64(v) * 2
       dfs(y, x)
       sumArr[x] += sumArr[y]
       sizeArr[x] += sizeArr[y]
   }
   // second pass: accumulate time and prepare for weighted ordering
   var q []Arr
   for _, e := range adj[x] {
       y, v := e.to, e.v
       if y == fa {
           continue
       }
       timeArr[x] += timeArr[y] + int64(v)*int64(sizeArr[y])
       q = append(q, Arr{sizeArr[y], sumArr[y]})
   }
   // sort children by increasing sum/size ratio
   sort.Slice(q, func(i, j int) bool {
       return q[i].sum*int64(q[j].size) < q[j].sum*int64(q[i].size)
   })
   var s int64
   for _, a := range q {
       timeArr[x] += int64(a.size) * s
       s += a.sum
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   adj = make([][]Edge, n+1)
   for i := 1; i < n; i++ {
       var x, y, v int
       fmt.Fscan(reader, &x, &y, &v)
       adj[x] = append(adj[x], Edge{to: y, v: v})
       adj[y] = append(adj[y], Edge{to: x, v: v})
   }
   sizeArr = make([]int, n+1)
   sumArr = make([]int64, n+1)
   timeArr = make([]int64, n+1)
   dfs(1, 0)
   // output expected time over n-1 possible treasure locations
   result := float64(timeArr[1]) / float64(n-1)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%.10f\n", result)
}
