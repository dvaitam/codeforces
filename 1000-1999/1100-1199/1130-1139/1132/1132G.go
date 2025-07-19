package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU maintains a differential array with union-find to count greedy subsequence lengths
type DSU struct {
   p  []int
   ld []int
   lb []int
}

// init initializes DSU for n elements
func (d *DSU) init(n int) {
   d.p = make([]int, n)
   d.ld = make([]int, n)
   d.lb = make([]int, n)
   for i := 0; i < n; i++ {
       d.p[i] = i
       d.lb[i] = i
   }
}

// find returns root of x with path compression
func (d *DSU) find(x int) int {
   for d.p[x] != x {
       d.p[x] = d.p[d.p[x]]
       x = d.p[x]
   }
   return x
}

// join merges sets containing x and y
func (d *DSU) join(x, y int) {
   x = d.find(x)
   y = d.find(y)
   if x == y {
       return
   }
   d.p[x] = y
   d.lb[y] = d.lb[x]
   d.ld[y] += d.ld[x]
}

// upd increments range [i, j) and merges segments with non-negative differential
func (d *DSU) upd(i, j int) {
   n := len(d.ld)
   i = d.find(i)
   d.ld[i]++
   if j < n {
       d.ld[j]--
   }
   for {
       root := d.find(i)
       if d.ld[root] < 0 || d.lb[root] == 0 {
           break
       }
       left := d.lb[root] - 1
       d.join(left, root)
   }
}

// upd2 removes the element sliding out of the window of size k
func (d *DSU) upd2(i, k int) {
   if i < k {
       return
   }
   j := d.find(i - k + 1)
   for d.lb[j] > i-k {
       left := d.lb[j] - 1
       d.join(left, j)
   }
}

// qry returns the current count at position 0
func (d *DSU) qry() int {
   return d.ld[d.find(0)]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   var d DSU
   d.init(n)
   l := make([]int, n)
   for i := 0; i < n; i++ {
       d.upd2(i, k)
       if i == 0 {
           l[i] = -1
       } else {
           l[i] = i - 1
           for l[i] >= 0 && a[l[i]] < a[i] {
               l[i] = l[l[i]]
           }
       }
       d.upd(l[i]+1, i+1)
       if i >= k-1 {
           fmt.Fprint(writer, d.qry(), " ")
       }
   }
}
