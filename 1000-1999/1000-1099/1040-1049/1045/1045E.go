package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   n        int
   x, y, c  []int
   pts      []int
   chull    []int
   ansL, ansR []int
)

func addEdge(u, v int) {
   ansL = append(ansL, u)
   ansR = append(ansR, v)
}

func printAns() {
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(ansL))
   for i := range ansL {
       fmt.Fprintln(w, ansL[i], ansR[i])
   }
}

func ccw(i, j, k int) bool {
   return x[i]*y[j] + x[j]*y[k] + x[k]*y[i] - x[j]*y[i] - x[k]*y[j] - x[i]*y[k] > 0
}

// recursive construction
func runRec(i, j, k int, inner []int, sep bool) {
   if len(inner) == 0 {
       if !sep {
           addEdge(j, k)
       }
       return
   }
   pivot := -1
   for _, p := range inner {
       if pivot == -1 || ccw(j, p, pivot) {
           pivot = p
       }
   }
   var jInner, kInner []int
   for _, p := range inner {
       if p == pivot {
           continue
       }
       if ccw(pivot, i, p) {
           jInner = append(jInner, p)
       } else {
           kInner = append(kInner, p)
       }
   }
   if c[pivot] == c[i] {
       runRec(j, pivot, i, jInner, false)
       runRec(k, i, pivot, kInner, true)
       if !sep {
           addEdge(j, k)
       }
   } else {
       runRec(i, j, pivot, jInner, false)
       runRec(i, pivot, k, kInner, sep)
   }
}

// initial run to collect inner points
func run(i, j, k int, sep bool) {
   var inner []int
   for l := 0; l < n; l++ {
       if l == i || l == j || l == k {
           continue
       }
       if ccw(i, j, l) && ccw(j, k, l) && ccw(k, i, l) {
           inner = append(inner, l)
       }
   }
   runRec(i, j, k, inner, sep)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n)
   x = make([]int, n)
   y = make([]int, n)
   c = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &x[i], &y[i], &c[i])
   }
   if n == 1 {
       fmt.Println(0)
       return
   }
   if n == 2 {
       if c[0] == c[1] {
           fmt.Println(1)
           fmt.Println(0, 1)
       } else {
           fmt.Println(0)
       }
       return
   }
   // pivot for sorting
   pivot := 0
   for i := 1; i < n; i++ {
       if x[i] < x[pivot] || (x[i] == x[pivot] && y[i] < y[pivot]) {
           pivot = i
       }
   }
   // prepare pts
   pts = make([]int, 0, n-1)
   for i := 0; i < n; i++ {
       if i != pivot {
           pts = append(pts, i)
       }
   }
   sort.Slice(pts, func(a, b int) bool {
       return ccw(pivot, pts[a], pts[b])
   })
   // graham scan
   chull = make([]int, 0, n)
   chull = append(chull, pivot)
   for _, cur := range pts {
       for len(chull) >= 2 && !ccw(chull[len(chull)-2], chull[len(chull)-1], cur) {
           chull = chull[:len(chull)-1]
       }
       chull = append(chull, cur)
   }
   chSize := len(chull)
   // count color changes
   numChange := 0
   for i := 0; i < chSize; i++ {
       j := (i + 1) % chSize
       if c[chull[i]] != c[chull[j]] {
           numChange++
       }
   }
   // handle cases
   if numChange == 0 {
       cCh := c[chull[0]]
       center := -1
       for i := 0; i < n; i++ {
           if c[i] != cCh {
               center = i
               break
           }
       }
       if center == -1 {
           // all same
           fmt.Println(n - 1)
           for _, p := range pts {
               fmt.Println(pivot, p)
           }
           return
       }
       // recursive around hull
       for i := 0; i < chSize-1; i++ {
           run(center, chull[i], chull[i+1], false)
       }
       run(center, chull[chSize-1], chull[0], true)
       printAns()
   } else if numChange == 2 {
       // find first change
       p1 := -1
       for i := 1; i < chSize; i++ {
           if c[chull[i-1]] != c[chull[i]] {
               p1 = i
               break
           }
       }
       // rotate
       chull = append(chull[p1:], chull[:p1]...)
       // find second change
       p2 := -1
       for i := 1; i < chSize; i++ {
           if c[chull[i-1]] != c[chull[i]] {
               p2 = i
               break
           }
       }
       // split
       for i := 0; i < p2-1; i++ {
           run(chull[p2], chull[i], chull[i+1], false)
       }
       for i := p2; i < chSize-1; i++ {
           run(chull[0], chull[i], chull[i+1], false)
       }
       printAns()
   } else {
       fmt.Println("Impossible")
   }
}
