package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
)

// Edge represents an undirected weighted edge
type Edge struct {
   u, v int
   w    int64
}

// DSU supports union-find with path compression
type DSU struct {
   parent []int
}

// NewDSU initializes a DSU for 1..n
func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &DSU{parent: p}
}

// Find returns the representative of x
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union merges the sets of x and y
func (d *DSU) Union(x, y int) {
   px := d.Find(x)
   py := d.Find(y)
   if px != py {
       d.parent[px] = py
   }
}

var rdr = bufio.NewReader(os.Stdin)
var wrt = bufio.NewWriter(os.Stdout)

// readInt reads next integer from stdin
func readInt() int {
   var x int
   var neg bool
   b, _ := rdr.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = rdr.ReadByte()
   }
   if b == '-' {
       neg = true
       b, _ = rdr.ReadByte()
   }
   for b >= '0' && b <= '9' {
       x = x*10 + int(b-'0')
       b, _ = rdr.ReadByte()
   }
   if neg {
       x = -x
   }
   return x
}

// readInt64 reads next 64-bit integer
func readInt64() int64 {
   var x int64
   var neg bool
   b, _ := rdr.ReadByte()
   for (b < '0' || b > '9') && b != '-' {
       b, _ = rdr.ReadByte()
   }
   if b == '-' {
       neg = true
       b, _ = rdr.ReadByte()
   }
   for b >= '0' && b <= '9' {
       x = x*10 + int64(b-'0')
       b, _ = rdr.ReadByte()
   }
   if neg {
       x = -x
   }
   return x
}

func main() {
   defer wrt.Flush()
   n := readInt()
   m := readInt()
   a := make([]int64, n+1)
   var mn int64 = 1<<62
   pos := 1
   for i := 1; i <= n; i++ {
       a[i] = readInt64()
       if a[i] < mn {
           mn = a[i]
           pos = i
       }
   }
   edges := make([]Edge, 0, m+n)
   for i := 0; i < m; i++ {
       u := readInt()
       v := readInt()
       w := readInt64()
       edges = append(edges, Edge{u: u, v: v, w: w})
   }
   // add edges to the minimal a-node
   for i := 1; i <= n; i++ {
       if i == pos {
           continue
       }
       edges = append(edges, Edge{u: i, v: pos, w: a[i] + mn})
   }
   // sort edges by weight
   sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
   dsu := NewDSU(n)
   var ans int64
   for _, e := range edges {
       if dsu.Find(e.u) != dsu.Find(e.v) {
           dsu.Union(e.u, e.v)
           ans += e.w
       }
   }
   wrt.WriteString(strconv.FormatInt(ans, 10))
}
