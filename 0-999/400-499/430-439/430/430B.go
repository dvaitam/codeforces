package main

import (
   "bufio"
   "fmt"
   "os"
)

// Ball represents a ball with its color and whether it's original
type Ball struct {
   color int
   orig  bool
}

// simulate insertion of a ball of color x at position pos and return number of original balls destroyed
func simulate(c []int, x int, pos int) int {
   // build initial list with inserted ball
   n := len(c)
   balls := make([]Ball, 0, n+1)
   for i := 0; i < pos; i++ {
       balls = append(balls, Ball{color: c[i], orig: true})
   }
   balls = append(balls, Ball{color: x, orig: false})
   for i := pos; i < n; i++ {
       balls = append(balls, Ball{color: c[i], orig: true})
   }
   // simulate removals
   for {
       m := len(balls)
       removed := false
       newBalls := make([]Ball, 0, m)
       i := 0
       for i < m {
           j := i + 1
           for j < m && balls[j].color == balls[i].color {
               j++
           }
           if j-i >= 3 {
               // skip these balls
               removed = true
           } else {
               // keep these balls
               for k := i; k < j; k++ {
                   newBalls = append(newBalls, balls[k])
               }
           }
           i = j
       }
       if !removed {
           break
       }
       balls = newBalls
   }
   // count remaining original balls
   rem := 0
   for _, b := range balls {
       if b.orig {
           rem++
       }
   }
   // destroyed = initial originals - remaining originals
   return len(c) - rem
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k, x int
   if _, err := fmt.Fscan(reader, &n, &k, &x); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   ans := 0
   // try all insertion positions
   for pos := 0; pos <= n; pos++ {
       destroyed := simulate(c, x, pos)
       if destroyed > ans {
           ans = destroyed
       }
   }
   fmt.Println(ans)
}
