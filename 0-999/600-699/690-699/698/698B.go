package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU with path compression and union by size
type DSU struct {
   parent []int
   size   []int
}

func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent, size}
}

func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

func (d *DSU) Union(x, y int) {
   px := d.Find(x)
   py := d.Find(y)
   if px == py {
       return
   }
   if d.size[px] > d.size[py] {
       d.parent[py] = px
       d.size[px] += d.size[py]
   } else {
       d.parent[px] = py
       d.size[py] += d.size[px]
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   dsu := NewDSU(n)
   var root int = -1
   selfloops := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if a[i] == i && root == -1 {
           root = i
       } else {
           u := dsu.Find(i)
           v := dsu.Find(a[i])
           if u == v {
               selfloops = append(selfloops, i)
           } else {
               dsu.Union(u, v)
           }
       }
   }
   changes := len(selfloops)
   if root == -1 {
       // choose one as root
       root = selfloops[len(selfloops)-1]
       selfloops = selfloops[:len(selfloops)-1]
   }
   // redirect root to itself
   a[root] = root
   // redirect other self loops to root
   for _, idx := range selfloops {
       a[idx] = root
   }
   // output
   fmt.Fprintln(writer, changes)
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, a[i])
   }
   fmt.Fprintln(writer)
}
