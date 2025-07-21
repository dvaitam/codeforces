package main

import (
   "bufio"
   "fmt"
   "os"
)

// Candy represents a candy with height requirement and mass benefit
type Candy struct {
   h int
   m int
}

// simulate runs the eating process starting with a given type (0 or 1)
// and returns the number of candies eaten
func simulate(startType int, candies0, candies1 []Candy, initialX int) int {
   x := initialX
   next := startType
   eaten := 0
   used0 := make([]bool, len(candies0))
   used1 := make([]bool, len(candies1))
   for {
       bestIdx, bestMass := -1, -1
       if next == 0 {
           for i, c := range candies0 {
               if !used0[i] && c.h <= x && c.m > bestMass {
                   bestMass = c.m
                   bestIdx = i
               }
           }
           if bestIdx == -1 {
               break
           }
           used0[bestIdx] = true
           x += candies0[bestIdx].m
       } else {
           for i, c := range candies1 {
               if !used1[i] && c.h <= x && c.m > bestMass {
                   bestMass = c.m
                   bestIdx = i
               }
           }
           if bestIdx == -1 {
               break
           }
           used1[bestIdx] = true
           x += candies1[bestIdx].m
       }
       eaten++
       // alternate type
       next ^= 1
   }
   return eaten
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   fmt.Fscan(reader, &n, &x)
   candies0 := make([]Candy, 0, n)
   candies1 := make([]Candy, 0, n)
   for i := 0; i < n; i++ {
       var t, h, m int
       fmt.Fscan(reader, &t, &h, &m)
       if t == 0 {
           candies0 = append(candies0, Candy{h: h, m: m})
       } else {
           candies1 = append(candies1, Candy{h: h, m: m})
       }
   }
   // try both starting types
   ans0 := simulate(0, candies0, candies1, x)
   ans1 := simulate(1, candies0, candies1, x)
   if ans1 > ans0 {
       fmt.Println(ans1)
   } else {
       fmt.Println(ans0)
   }
