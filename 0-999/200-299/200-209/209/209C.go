package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU for int elements
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := 0; i <= n; i++ {
       p[i] = i
   }
   return &DSU{p: p}
}

func (d *DSU) Find(x int) int {
   // path compression
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx != ry {
       d.p[ry] = rx
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   dsu := NewDSU(n)
   deg := make([]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       deg[x]++
       deg[y]++
       dsu.Union(x, y)
   }
   // mark components with edges
   hasEdge := make([]bool, n+1)
   for v := 1; v <= n; v++ {
       if deg[v] > 0 {
           r := dsu.Find(v)
           hasEdge[r] = true
       }
   }
   // compute odd count per component
   oddMap := make([]int, n+1)
   totalOdd := 0
   for v := 1; v <= n; v++ {
       if deg[v]&1 == 1 {
           totalOdd++
           r := dsu.Find(v)
           oddMap[r]++
       }
   }
   // determine component containing node 1
   root1 := dsu.Find(1)
   // build list of other components with edges
   nonZero := []int{}
   zeroCnt := 0
   for r := 1; r <= n; r++ {
       if r == root1 || !hasEdge[r] {
           continue
       }
       if oddMap[r] > 0 {
           nonZero = append(nonZero, oddMap[r])
       } else {
           zeroCnt++
       }
   }
   // simulate connections
   odd1 := oddMap[root1]
   connections := 0
   // first connect components with odd nodes
   for _, oc := range nonZero {
       if odd1 > 0 {
           // connect odd-odd: use one odd from each
           odd1--
           totalOdd -= 2
       } else {
           // connect odd-even: odd in other, even in comp1
           odd1++
       }
       connections++
   }
   // then connect components without odd nodes
   for i := 0; i < zeroCnt; i++ {
       if odd1 > 0 {
           // connect odd-even: parity shifts within merged comp, totalOdd unchanged
       } else {
           // connect even-even: creates two odd
           odd1++
           totalOdd += 2
       }
       connections++
   }
   // finally, fix remaining odd vertices by pairing
   needed := totalOdd / 2
   ans := connections + needed
   fmt.Fprintln(writer, ans)
}
