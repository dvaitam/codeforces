package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU implements disjoint set union with union by size and path compression
type DSU struct {
   parent []int
   size   []int
}

// NewDSU creates a new DSU for n elements (0..n-1)
func NewDSU(n int) *DSU {
   parent := make([]int, n)
   size := make([]int, n)
   for i := 0; i < n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent, size}
}

// Find returns the leader of x with path compression
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Merge unions the sets containing x and y. Returns true if merged.
func (d *DSU) Merge(x, y int) bool {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return false
   }
   // attach rx under ry
   d.parent[rx] = ry
   d.size[ry] += d.size[rx]
   return true
}

// ComponentSize returns the size of the component containing x
func (d *DSU) ComponentSize(x int) int {
   return d.size[d.Find(x)]
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       ans := 0
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
           if a[i] == 0 {
               a[i] = 1
               ans++
           }
       }

       // check connectivity via bits
       check := func() bool {
           dsu := NewDSU(n)
           for k := 0; k <= 30; k++ {
               first := -1
               mask := 1 << k
               for i := 0; i < n; i++ {
                   if a[i]&mask != 0 {
                       if first == -1 {
                           first = i
                       } else {
                           dsu.Merge(i, first)
                       }
                   }
               }
           }
           return dsu.ComponentSize(0) == n
       }

       // output result
       show := func() {
           fmt.Fprintln(writer, ans)
           for i := 0; i < n; i++ {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, a[i])
           }
           fmt.Fprintln(writer)
       }

       if check() {
           show()
           continue
       }
       // try increment each
       done := false
       for i := 0; i < n && !done; i++ {
           a[i]++
           ans++
           if check() {
               show()
               done = true
           }
           if !done {
               a[i]--
               ans--
           }
       }
       if done {
           continue
       }
       // try decrement each
       for i := 0; i < n && !done; i++ {
           a[i]--
           ans++
           if check() {
               show()
               done = true
           }
           if !done {
               a[i]++
               ans--
           }
       }
       if done {
           continue
       }
       // fallback: adjust two elements with maximal lowbit
       mx := 0
       for i := 0; i < n; i++ {
           low := a[i] & -a[i]
           if low > mx {
               mx = low
           }
       }
       // decrement first with lowbit mx
       for i := 0; i < n; i++ {
           if a[i]&-a[i] == mx {
               a[i]--
               ans++
               break
           }
       }
       // increment last with lowbit mx
       for i := n - 1; i >= 0; i-- {
           if a[i]&-a[i] == mx {
               a[i]++
               ans++
               break
           }
       }
       show()
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   solve(reader, writer)
}
