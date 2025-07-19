package main

import (
   "fmt"
)

// dp[r][s][p] stores probability for state (r, s, p)
var dp [101][101][101]float64

// prob returns probability that 'rock' wins when counts are (r, s, p)
func prob(r, s, p int) float64 {
   if dp[r][s][p] != 0 {
       return dp[r][s][p]
   }
   // Base cases
   if r == 0 {
       return 0
   }
   if p == 0 {
       return 1
   }
   if s == 0 {
       return 0
   }
   // Total possible interactions
   rs := float64(r * s)
   sp := float64(s * p)
   pr := float64(p * r)
   total := rs + sp + pr
   // Recurrence: consider which pair interacts next
   // rock-paper: r*p/total leads to rock count unchanged? but rock loses one: r-1
   // rock-scissors: r*s/total leads to scissors decrease: s-1
   // scissors-paper: s*p/total leads to paper decrease: p-1
   d := (pr/total)*prob(r-1, s, p)
   d += (rs/total)*prob(r, s-1, p)
   d += (sp/total)*prob(r, s, p-1)
   dp[r][s][p] = d
   return d
}

func main() {
   var r, s, p int
   if _, err := fmt.Scan(&r, &s, &p); err != nil {
       return
   }
   // Compute probabilities for rock, scissors, paper wins
   // prob(r, s, p) gives rock wins
   // For scissors wins: rotate (s, p, r)
   // For paper wins: rotate (p, r, s)
   fR := prob(r, s, p)
   fS := prob(s, p, r)
   fP := prob(p, r, s)
   fmt.Printf("%.9f %.9f %.9f\n", fR, fS, fP)
}
