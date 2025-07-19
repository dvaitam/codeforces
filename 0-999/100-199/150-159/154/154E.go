package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
   "sort"
)

// Point represents a 2D point.
type Point struct { x, y float64 }

func sub(a, b Point) Point { return Point{a.x - b.x, a.y - b.y} }
func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b Point) float64   { return a.x*b.x + a.y*b.y }
func absPt(a Point) float64     { return math.Hypot(a.x, a.y) }

// Item is an entry in the priority queue.
type Item struct {
   r   float64
   idx int
   ver int
}
// PriorityQueue implements a max-heap of Items by r.
type PriorityQueue []*Item
func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].r > pq[j].r }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq; n := len(old)
   item := old[n-1]
   *pq = old[0 : n-1]
   return item
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var R float64
   if _, err := fmt.Fscan(in, &n, &R); err != nil {
       return
   }
   pts := make([]Point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
   }
   if n == 1 {
       fmt.Printf("0.000000\n")
       return
   }
   // Sort by x, then y
   sort.Slice(pts, func(i, j int) bool {
       if math.Abs(pts[i].x-pts[j].x) > 1e-9 {
           return pts[i].x < pts[j].x
       }
       return pts[i].y < pts[j].y
   })
   // Build convex hull (monotone chain)
   q := make([]Point, 0, 2*n)
   for i := 0; i < n; i++ {
       for len(q) >= 2 && cross(sub(q[len(q)-1], q[len(q)-2]), sub(pts[i], q[len(q)-2])) <= 0 {
           q = q[:len(q)-1]
       }
       q = append(q, pts[i])
   }
   t := len(q)
   for i := n - 1; i >= 0; i-- {
       for len(q) > t && cross(sub(q[len(q)-1], q[len(q)-2]), sub(pts[i], q[len(q)-2])) <= 0 {
           q = q[:len(q)-1]
       }
       q = append(q, pts[i])
   }
   // Remove duplicate last point
   q = q[:len(q)-1]
   m := len(q)
   l := make([]int, m)
   r := make([]int, m)
   del := make([]bool, m)
   version := make([]int, m)
   for i := 0; i < m; i++ {
       l[i] = (i-1+m)%m
       r[i] = (i+1)%m
   }
   var clk int
   pqh := &PriorityQueue{}
   heap.Init(pqh)
   // delete index x from linked list
   delFunc := func(x int) {
       del[x] = true
       r[l[x]] = r[x]
       l[r[x]] = l[x]
   }
   // update candidate circle at x
   update := func(x int) {
       // skip if angle is acute
       if dot(sub(q[l[x]], q[x]), sub(q[r[x]], q[x])) > 0 {
           version[x] = -1
           return
       }
       a, b, c := q[x], q[l[x]], q[r[x]]
       ab := absPt(sub(a, b))
       ac := absPt(sub(a, c))
       bc := absPt(sub(b, c))
       area2 := math.Abs(cross(sub(a, b), sub(c, b)))
       radius := ab * ac * bc / (2 * area2)
       clk++
       version[x] = clk
       heap.Push(pqh, &Item{r: radius, idx: x, ver: clk})
   }
   // init queue
   for i := 0; i < m; i++ {
       if dot(sub(q[l[i]], q[i]), sub(q[r[i]], q[i])) < 0 {
           update(i)
       }
   }
   // remove points with radius > R
   for pqh.Len() > 0 {
       top := (*pqh)[0]
       if top.r <= R {
           break
       }
       item := heap.Pop(pqh).(*Item)
       if item.ver != version[item.idx] {
           continue
       }
       delFunc(item.idx)
       update(l[item.idx])
       update(r[item.idx])
   }
   // compute final area with circular arcs
   s := 0
   for del[s] {
       s++
   }
   ans := 0.0
   i := s
   for {
       ni := r[i]
       ans += cross(q[i], q[ni]) / 2
       d := absPt(sub(q[i], q[ni]))
       v := d / (2 * R)
       if v > 1 {
           v = 1
       } else if v < -1 {
           v = -1
       }
       angle := 2 * math.Asin(v)
       ans += R*R*(angle-math.Sin(angle)) / 2
       i = ni
       if i == s {
           break
       }
   }
   fmt.Printf("%.10f\n", ans)
}
