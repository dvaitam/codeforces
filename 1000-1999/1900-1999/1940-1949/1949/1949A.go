package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Point represents a grid point
type Point struct {
   x, y int
}

// Solver finds the maximum clique in the adjacency graph
type Solver struct {
   n    int
   exi  [][]bool
   ans  []int
   tans []int
   f    []int
   id   []Point
}

// dfs attempts to build cliques starting from node x
func (s *Solver) dfs(x int) bool {
   for i := x + 1; i < s.n; i++ {
       if s.f[i]+len(s.tans) <= len(s.ans) {
           return false
       }
       if !s.exi[x][i] {
           continue
       }
       ok := true
       for _, j := range s.tans {
           if !s.exi[i][j] {
               ok = false
               break
           }
       }
       if ok {
           s.tans = append(s.tans, i)
           if s.dfs(i) {
               return true
           }
           s.tans = s.tans[:len(s.tans)-1]
       }
   }
   if len(s.tans) > len(s.ans) {
       s.ans = make([]int, len(s.tans))
       copy(s.ans, s.tans)
       return true
   }
   return false
}

// Run executes the solver and outputs the result
func (s *Solver) Run() {
   s.f = make([]int, s.n)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := s.n - 1; i >= 0; i-- {
       s.tans = s.tans[:0]
       s.tans = append(s.tans, i)
       s.dfs(i)
       s.f[i] = len(s.ans)
   }
   fmt.Fprintln(writer, len(s.ans))
   for _, idx := range s.ans {
       p := s.id[idx]
       fmt.Fprintln(writer, p.x, p.y)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var r float64
   if _, err := fmt.Fscan(reader, &n, &r); err != nil {
       return
   }
   mr := int(math.Ceil(r))
   // collect valid points
   id := make([]Point, 0)
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           if i >= mr && j >= mr && n-i >= mr && n-j >= mr {
               id = append(id, Point{i, j})
           }
       }
   }
   m := len(id)
   // build adjacency
   exi := make([][]bool, m)
   for i := range exi {
       exi[i] = make([]bool, m)
   }
   for i := 0; i < m; i++ {
       for j := 0; j < m; j++ {
           if i != j {
               dx := float64(id[i].x - id[j].x)
               dy := float64(id[i].y - id[j].y)
               if math.Hypot(dx, dy)-2*r > -1e-9 {
                   exi[i][j] = true
               }
           }
       }
   }
   solver := Solver{n: m, exi: exi, id: id}
   solver.Run()
}
