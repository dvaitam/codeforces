package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Point struct {
   X, Y int64
}

type Ship struct {
   pos Point
   idx int
   typ int
}

var ans []int

func cross(a, b Point) int64 {
   return a.X*b.Y - a.Y*b.X
}

func dot(a, b Point) int64 {
   return a.X*b.X + a.Y*b.Y
}

// CompareY: by Y then X
func compareY(a, b Point) bool {
   if a.Y != b.Y {
       return a.Y < b.Y
   }
   return a.X < b.X
}

// internal as in C++ code
func internal(a, b Point) bool {
   c := cross(a, b)
   if c != 0 {
       return c > 0
   }
   // !CompareY(b, a)
   return !(compareY(b, a))
}

// extract items in circular slice [head, tail)
func extract(items []Ship, head, tail int) []Ship {
   var out []Ship
   if head < tail {
       out = append(out, items[head:tail]...)
   } else {
       out = append(out, items[head:]...)
       out = append(out, items[:tail]...)
   }
   return out
}

func solve(oitems []Ship) {
   N := len(oitems)
   if N == 2 {
       if oitems[0].typ == 1 {
           ans[oitems[0].idx] = oitems[1].idx
       } else {
           ans[oitems[1].idx] = oitems[0].idx
       }
       return
   }
   // make a copy
   items := make([]Ship, N)
   copy(items, oitems)
   // compute mid sum
   var midX, midY int64
   for i := 0; i < N; i++ {
       midX += items[i].pos.X
       midY += items[i].pos.Y
   }
   // shift and scale
   mul := int64(N)
   for i := 0; i < N; i++ {
       items[i].pos.X = items[i].pos.X*mul - midX
       items[i].pos.Y = items[i].pos.Y*mul - midY
   }
   // sort by angle around origin
   sort.Slice(items, func(i, j int) bool {
       a := items[i].pos
       b := items[j].pos
       fa := (a.Y > 0) || (a.Y == 0 && a.X >= 0)
       fb := (b.Y > 0) || (b.Y == 0 && b.X >= 0)
       if fa != fb {
           return fa
       }
       c := cross(a, b)
       if c != 0 {
           return c > 0
       }
       // tie by distance
       return dot(a, a) < dot(b, b)
   })
   // two pointers
   tail := 0
   sum := 0
   for head := 0; head < N; head++ {
       for internal(items[head].pos, items[tail].pos) {
           sum += items[tail].typ
           tail++
           if tail == N {
               tail = 0
           }
       }
       if sum == 0 {
           // restore original positions
           for i := 0; i < N; i++ {
               items[i].pos.X = (items[i].pos.X + midX) / mul
               items[i].pos.Y = (items[i].pos.Y + midY) / mul
           }
           // split and recurse
           part1 := extract(items, head, tail)
           part2 := extract(items, tail, head)
           solve(part1)
           solve(part2)
           return
       }
       sum -= items[head].typ
   }
   panic("unreachable")
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var N int
   fmt.Fscan(reader, &N)
   items := make([]Ship, 2*N)
   for i := 0; i < N; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       items[i] = Ship{pos: Point{X: x, Y: y}, idx: i, typ: 1}
   }
   for i := 0; i < N; i++ {
       var x, y int64
       fmt.Fscan(reader, &x, &y)
       items[i+N] = Ship{pos: Point{X: x, Y: y}, idx: i, typ: -1}
   }
   ans = make([]int, N)
   solve(items)
   for i := 0; i < N; i++ {
       // output 1-based
       fmt.Fprintln(writer, ans[i]+1)
   }
}
