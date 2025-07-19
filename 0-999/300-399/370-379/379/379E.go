package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Line represents a linear segment over [x0, x1]: f(x) = a*x + b
type Line struct {
   a, b   int
   x0, x1 float64
}

func eval(a, b int, x float64) float64 {
   return float64(a)*x + float64(b)
}

func P(a, b int, x0, x1 float64) float64 {
   return (eval(a, b, x0) + eval(a, b, x1)) * (x1 - x0) / 2.0
}

// lowerBound finds first index in lines where lines[i].a >= a
func lowerBound(lines []Line, a int) int {
   return sort.Search(len(lines), func(i int) bool { return lines[i].a >= a })
}

func removeAt(lines []Line, idx int) []Line {
   return append(lines[:idx], lines[idx+1:]...)
}

func insertAt(lines []Line, idx int, ln Line) []Line {
   lines = append(lines, Line{})
   copy(lines[idx+1:], lines[idx:])
   lines[idx] = ln
   return lines
}

func min(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b float64) float64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // S holds for each strip a set of visible line segments sorted by slope a
   S := make([][]Line, k)
   for i := 0; i < n; i++ {
       var y0 int
       fmt.Fscan(reader, &y0)
       var area float64
       for j := 0; j < k; j++ {
           var y1 int
           fmt.Fscan(reader, &y1)
           a := y1 - y0
           b := y0
           segs := S[j]
           // initialize
           if len(segs) == 0 {
               // full [0,1]
               segs = []Line{{a: a, b: b, x0: 0.0, x1: 1.0}}
               area += P(a, b, 0.0, 1.0)
               S[j] = segs
               y0 = y1
               continue
           }
           left, right := 1e18, -1e18
           // forward pass: remove or cut segments to the right
           idx := lowerBound(segs, a)
           for idx < len(segs) {
               l := segs[idx]
               if eval(a, b, l.x0) >= eval(l.a, l.b, l.x0) && eval(a, b, l.x1) >= eval(l.a, l.b, l.x1) {
                   area += P(a, b, l.x0, l.x1) - P(l.a, l.b, l.x0, l.x1)
                   left = min(left, l.x0)
                   right = max(right, l.x1)
                   segs = removeAt(segs, idx)
                   continue
               }
               // partial override on left side
               if eval(a, b, l.x0) > eval(l.a, l.b, l.x0) {
                   xsr := float64(b-l.b) / float64(l.a-a)
                   area += P(a, b, l.x0, xsr) - P(l.a, l.b, l.x0, xsr)
                   left = min(left, l.x0)
                   right = max(right, xsr)
                   // keep right part of original
                   segs[idx] = Line{a: l.a, b: l.b, x0: xsr, x1: l.x1}
               }
               break
           }
           // backward pass: remove or cut segments to the left
           idx = lowerBound(segs, a)
           if idx > 0 {
               idx--
               for idx >= 0 {
                   l := segs[idx]
                   if eval(a, b, l.x0) >= eval(l.a, l.b, l.x0) && eval(a, b, l.x1) >= eval(l.a, l.b, l.x1) {
                       area += P(a, b, l.x0, l.x1) - P(l.a, l.b, l.x0, l.x1)
                       left = min(left, l.x0)
                       right = max(right, l.x1)
                       segs = removeAt(segs, idx)
                       idx--
                       continue
                   }
                   if eval(a, b, l.x1) > eval(l.a, l.b, l.x1) {
                       xsr := float64(b-l.b) / float64(l.a-a)
                       area += P(a, b, xsr, l.x1) - P(l.a, l.b, xsr, l.x1)
                       left = min(left, xsr)
                       right = max(right, l.x1)
                       // keep left part
                       segs[idx] = Line{a: l.a, b: l.b, x0: l.x0, x1: xsr}
                   }
                   break
               }
           }
           // insert new visible segment
           if left <= right {
               insertPos := lowerBound(segs, a)
               segs = insertAt(segs, insertPos, Line{a: a, b: b, x0: left, x1: right})
           }
           S[j] = segs
           y0 = y1
       }
       // output visible area for this piece
       fmt.Fprintf(writer, "%.10f\n", area)
   }
}
