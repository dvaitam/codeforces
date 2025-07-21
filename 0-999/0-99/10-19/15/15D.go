package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
   "sort"
)

// 2D BIT for point updates and prefix sum queries
type BIT2D struct {
   n, m int
   t    [][]int
}

func NewBIT2D(n, m int) *BIT2D {
   t := make([][]int, n+1)
   for i := range t {
       t[i] = make([]int, m+1)
   }
   return &BIT2D{n: n, m: m, t: t}
}

// add v at (x,y), 0-based
func (b *BIT2D) Update(x, y, v int) {
   for i := x + 1; i <= b.n; i += i & -i {
       row := b.t[i]
       for j := y + 1; j <= b.m; j += j & -j {
           row[j] += v
       }
   }
}

// sum of rectangle [0..x) x [0..y)
func (b *BIT2D) Sum(x, y int) int {
   s := 0
   for i := x; i > 0; i -= i & -i {
       row := b.t[i]
       for j := y; j > 0; j -= j & -j {
           s += row[j]
       }
   }
   return s
}

// query sum of [x1..x2], [y1..y2]
func (b *BIT2D) RectSum(x1, y1, x2, y2 int) int {
   return b.Sum(x2+1, y2+1) - b.Sum(x1, y2+1) - b.Sum(x2+1, y1) + b.Sum(x1, y1)
}

type Rect struct {
   cost    int64
   i, j    int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, a, b int
   if _, err := fmt.Fscan(in, &n, &m, &a, &b); err != nil {
       return
   }
   h := make([][]int64, n)
   for i := 0; i < n; i++ {
       h[i] = make([]int64, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &h[i][j])
       }
   }
   // prefix sums
   ps := make([][]int64, n+1)
   for i := range ps {
       ps[i] = make([]int64, m+1)
   }
   for i := 0; i < n; i++ {
       rowSum := int64(0)
       for j := 0; j < m; j++ {
           rowSum += h[i][j]
           ps[i+1][j+1] = ps[i][j+1] + rowSum
       }
   }
   // row minima of width b
   wm := m - b + 1
   rowMin := make([][]int64, n)
   for i := 0; i < n; i++ {
       rowMin[i] = make([]int64, wm)
       dq := list.New()
       for j := 0; j < m; j++ {
           for dq.Len() > 0 {
               e := dq.Back()
               if h[i][e.Value.(int)] > h[i][j] {
                   dq.Remove(e)
               } else {
                   break
               }
           }
           dq.PushBack(j)
           if dq.Front().Value.(int) <= j-b {
               dq.Remove(dq.Front())
           }
           if j >= b-1 {
               rowMin[i][j-b+1] = h[i][dq.Front().Value.(int)]
           }
       }
   }
   // full minima of size a x b
   hm := n - a + 1
   mins := make([][]int64, hm)
   for i := 0; i < hm; i++ {
       mins[i] = make([]int64, wm)
   }
   for j := 0; j < wm; j++ {
       dq := list.New()
       for i := 0; i < n; i++ {
           for dq.Len() > 0 {
               e := dq.Back()
               if rowMin[e.Value.(int)][j] > rowMin[i][j] {
                   dq.Remove(e)
               } else {
                   break
               }
           }
           dq.PushBack(i)
           if dq.Front().Value.(int) <= i-a {
               dq.Remove(dq.Front())
           }
           if i >= a-1 {
               mins[i-a+1][j] = rowMin[dq.Front().Value.(int)][j]
           }
       }
   }
   // prepare rectangles
   area := int64(a * b)
   total := hm * wm
   rects := make([]Rect, 0, total)
   for i := 0; i < hm; i++ {
       for j := 0; j < wm; j++ {
           sum := ps[i+a][j+b] - ps[i][j+b] - ps[i+a][j] + ps[i][j]
           cost := sum - area*mins[i][j]
           rects = append(rects, Rect{cost: cost, i: i, j: j})
       }
   }
   sort.Slice(rects, func(i, j int) bool {
       if rects[i].cost != rects[j].cost {
           return rects[i].cost < rects[j].cost
       }
       if rects[i].i != rects[j].i {
           return rects[i].i < rects[j].i
       }
       return rects[i].j < rects[j].j
   })
   bit := NewBIT2D(n, m)
   results := make([][3]int64, 0)
   // select
   for _, r := range rects {
       // check overlap
       if bit.RectSum(r.i, r.j, r.i+a-1, r.j+b-1) > 0 {
           continue
       }
       // select
       results = append(results, [3]int64{int64(r.i + 1), int64(r.j + 1), r.cost})
       // mark occupied
       for x := r.i; x < r.i+a; x++ {
           for y := r.j; y < r.j+b; y++ {
               bit.Update(x, y, 1)
           }
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, len(results))
   for _, v := range results {
       fmt.Fprintf(out, "%d %d %d\n", v[0], v[1], v[2])
   }
}
