package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Key struct { X, Y int }
type Pair struct { X, I int }
type Point struct { X, Y, I int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       solve(reader, writer)
   }
}

func solve(r *bufio.Reader, w *bufio.Writer) {
   var n, m int
   fmt.Fscan(r, &n, &m)
   m /= 2
   points := make([]Point, n)
   pointMap := make(map[Key]int, n)
   mmp1 := make(map[int][]Pair)
   mmp2 := make(map[int][]Pair)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(r, &x, &y)
       idx := i + 1
       points[i] = Point{x, y, idx}
       pointMap[Key{x, y}] = idx
       sum := x + y
       mmp1[sum] = append(mmp1[sum], Pair{X: x, I: idx})
       diff := x - y
       mmp2[diff] = append(mmp2[diff], Pair{X: x, I: idx})
   }
   for k, v := range mmp1 {
       sort.Slice(v, func(i, j int) bool { return v[i].X < v[j].X })
       mmp1[k] = v
   }
   for k, v := range mmp2 {
       sort.Slice(v, func(i, j int) bool { return v[i].X < v[j].X })
       mmp2[k] = v
   }
   done := false
   for _, p := range points {
       if done {
           break
       }
       x, y, idx := p.X, p.Y, p.I
       // Case 1: point at (x+m, y+m)
       if id2, ok := pointMap[Key{x + m, y + m}]; ok {
           key := (x + m) - (y - m)
           tmp := mmp2[key]
           i := sort.Search(len(tmp), func(i int) bool { return tmp[i].X >= x + m })
           if i < len(tmp) && tmp[i].X <= x + m + m {
               fmt.Fprintf(w, "%d %d %d\n", idx, id2, tmp[i].I)
               done = true
               break
           }
       }
       // Case 2: point at (x+m, y+m)
       if id2, ok := pointMap[Key{x + m, y + m}]; ok {
           key := (x - m) - (y + m)
           tmp := mmp2[key]
           i := sort.Search(len(tmp), func(i int) bool { return tmp[i].X >= x - m })
           if i < len(tmp) && tmp[i].X <= x {
               fmt.Fprintf(w, "%d %d %d\n", idx, id2, tmp[i].I)
               done = true
               break
           }
       }
       // Case 3: point at (x+m, y-m)
       if id2, ok := pointMap[Key{x + m, y - m}]; ok {
           key := (x + m) + (y + m)
           tmp := mmp1[key]
           i := sort.Search(len(tmp), func(i int) bool { return tmp[i].X >= x + m })
           if i < len(tmp) && tmp[i].X <= x + m + m {
               fmt.Fprintf(w, "%d %d %d\n", idx, id2, tmp[i].I)
               done = true
               break
           }
       }
       // Case 4: point at (x+m, y-m)
       if id2, ok := pointMap[Key{x + m, y - m}]; ok {
           key := (x - m) + (y - m)
           tmp := mmp1[key]
           i := sort.Search(len(tmp), func(i int) bool { return tmp[i].X >= x - m })
           if i < len(tmp) && tmp[i].X <= x {
               fmt.Fprintf(w, "%d %d %d\n", idx, id2, tmp[i].I)
               done = true
               break
           }
       }
   }
   if !done {
       fmt.Fprintf(w, "0 0 0\n")
   }
}
