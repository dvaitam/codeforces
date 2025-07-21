package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var A, B int
   if _, err := fmt.Fscan(in, &A, &B); err != nil {
       return
   }
   grid := make([][]byte, A)
   for i := 0; i < A; i++ {
       var s string
       fmt.Fscan(in, &s)
       // assume s length equals B
       grid[i] = []byte(s)
   }
   // collect divisors
   var divA, divB []int
   for x := 1; x <= A; x++ {
       if A % x == 0 {
           divA = append(divA, x)
       }
   }
   for y := 1; y <= B; y++ {
       if B % y == 0 {
           divB = append(divB, y)
       }
   }
   type pair struct{ x, y int }
   var good []pair
   for _, X := range divA {
       for _, Y := range divB {
           seen := make(map[string]struct{})
           ok := true
           rows := A / X
           cols := B / Y
           for i := 0; i < rows && ok; i++ {
               for j := 0; j < cols; j++ {
                   // extract piece at (i,j)
                   piece := make([][]byte, X)
                   for xi := 0; xi < X; xi++ {
                       piece[xi] = grid[i*X+xi][j*Y : j*Y+Y]
                   }
                   // get canonical rep
                   reps := []string{repr(piece)}
                   // 180 rotation
                   reps = append(reps, repr(rot180(piece)))
                   // if square, include 90 and 270
                   if X == Y {
                       reps = append(reps, repr(rot90(piece)))
                       reps = append(reps, repr(rot270(piece)))
                   }
                   // pick minimal
                   min := reps[0]
                   for _, s := range reps[1:] {
                       if s < min {
                           min = s
                       }
                   }
                   if _, found := seen[min]; found {
                       ok = false
                       break
                   }
                   seen[min] = struct{}{}
               }
           }
           if ok {
               good = append(good, pair{X, Y})
           }
       }
   }
   // output count and minimal by area then X
   fmt.Println(len(good))
   // find minimal
   best := good[0]
   bestArea := best.x * best.y
   for _, p := range good[1:] {
       area := p.x * p.y
       if area < bestArea || (area == bestArea && p.x < best.x) {
           best = p
           bestArea = area
       }
   }
   fmt.Printf("%d %d\n", best.x, best.y)
}

// repr returns a string representation of piece rows separated by '|'
func repr(p [][]byte) string {
   // flatten with row separator
   var s []byte
   for i, row := range p {
       s = append(s, row...)
       if i < len(p)-1 {
           s = append(s, '|')
       }
   }
   return string(s)
}

func rot180(p [][]byte) [][]byte {
   X := len(p)
   Y := len(p[0])
   r := make([][]byte, X)
   for i := 0; i < X; i++ {
       r[i] = make([]byte, Y)
       for j := 0; j < Y; j++ {
           r[i][j] = p[X-1-i][Y-1-j]
       }
   }
   return r
}

func rot90(p [][]byte) [][]byte {
   X := len(p)
   Y := len(p[0])
   // result Y x X
   r := make([][]byte, X)
   // since X==Y for usage, we can size same
   for i := 0; i < X; i++ {
       r[i] = make([]byte, X)
   }
   for i := 0; i < X; i++ {
       for j := 0; j < Y; j++ {
           r[j][X-1-i] = p[i][j]
       }
   }
   return r
}

func rot270(p [][]byte) [][]byte {
   // rot270 = rot180(rot90)
   return rot180(rot90(p))
}
