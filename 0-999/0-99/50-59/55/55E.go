package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Vec struct {
   x, y int64
   half int
}

type Vecs []Vec

func (v Vecs) Len() int { return len(v) }
func (v Vecs) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v Vecs) Less(i, j int) bool {
   if v[i].half != v[j].half {
       return v[i].half < v[j].half
   }
   // cross(v[i], v[j]) > 0
   return v[i].x*v[j].y - v[i].y*v[j].x > 0
}

func countTriangles(n int, px, py int64, xs, ys []int64) int64 {
   // build and sort vectors by angle
   vecs := make(Vecs, n)
   for i := 0; i < n; i++ {
       dx := xs[i] - px
       dy := ys[i] - py
       half := 1
       if dy > 0 || (dy == 0 && dx > 0) {
           half = 0
       }
       vecs[i] = Vec{x: dx, y: dy, half: half}
   }
   sort.Sort(vecs)
   // duplicate for circular wrap
   vecs = append(vecs, vecs...)
   var bad int64
   j := 0
   // two pointers
   for i := 0; i < n; i++ {
       if j < i+1 {
           j = i + 1
       }
       for j < i+n && (vecs[i].x*vecs[j].y - vecs[i].y*vecs[j].x) > 0 {
           j++
       }
       k := int64(j - i - 1)
       if k >= 2 {
           bad += k * (k - 1) / 2
       }
   }
   total := int64(n)
   total = total * (total - 1) * (total - 2) / 6
   return total - bad
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       var px, py int64
       fmt.Fscan(in, &px, &py)
       ans := countTriangles(n, px, py, xs, ys)
       fmt.Fprintln(out, ans)
   }
}
