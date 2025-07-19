package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

var (
   hp    [2]int
   dt    [2]int
   l     [2]int
   r     [2]int
   p     [2]int
   z     [2]int
   state [2][2][205]float64
   w     int
)

// goStep simulates a shot by the opponent of f and updates the state
func goStep(f int) {
   wp := 1 - w
   // initialize probabilities for non-piercing (misses)
   for j := 0; j <= hp[f]; j++ {
       state[wp][f][j] = float64(p[1-f]) * 0.01 * state[w][f][j]
   }
   // sliding window sum for piercing (hits)
   var s float64
   for j := hp[f]; j >= 0; j-- {
       if j+l[1-f] <= hp[f] {
           s += state[w][f][j+l[1-f]]
       }
       if j+r[1-f]+1 <= hp[f] {
           s -= state[w][f][j+r[1-f]+1]
       }
       state[wp][f][j] += 0.01 * float64(100-p[1-f]) / float64(z[1-f]) * s
   }
   // account for overkill into zero state
   for j := 0; j <= hp[f]; j++ {
       if j >= r[1-f] {
           continue
       }
       a := j - r[1-f]
       b := min(j-l[1-f], -1)
       pq := float64(b-a+1) / float64(z[1-f]) * float64(100-p[1-f]) * 0.01
       state[wp][f][0] += pq * state[w][f][j]
   }
   // copy unchanged states for the other player
   for j := 0; j <= hp[1-f]; j++ {
       state[wp][1-f][j] = state[w][1-f][j]
   }
   w = wp
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   for i := 0; i < 2; i++ {
       fmt.Fscan(reader, &hp[i], &dt[i], &l[i], &r[i], &p[i])
       z[i] = r[i] - l[i] + 1
   }
   if p[0] == 100 || p[1] == 100 {
       if p[0] == 100 {
           fmt.Printf("%.6f\n", 0.0)
       } else {
           fmt.Printf("%.6f\n", 1.0)
       }
       return
   }
   w = 0
   state[w][0][hp[0]] = 1.0
   state[w][1][hp[1]] = 1.0
   var result float64
outer:
   for t := 0; ; t++ {
       if t%dt[0] == 0 {
           result -= state[w][1][0] * (1.0 - state[w][0][0])
           goStep(1)
           result += state[w][1][0] * (1.0 - state[w][0][0])
           q := state[w][0][0] + state[w][1][0] - state[w][0][0]*state[w][1][0]
           if q+1e-7 > 1.0 {
               break outer
           }
       }
       if t%dt[1] == 0 {
           goStep(0)
       }
   }
   fmt.Printf("%.6f\n", result)
}
