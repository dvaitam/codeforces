package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a [5]int
   var b [5]int
   for i := 1; i <= 4; i++ {
       if _, err := fmt.Fscan(reader, &a[i], &b[i]); err != nil {
           return
       }
   }
   // payoff matrix: payoff[c1][c2] for team1 choice c1 and team2 choice c2
   var payoff [2][2]int
   // team1 choices: c1=0: p1 attack, p2 defence; c1=1: p2 attack, p1 defence
   // team2 choices: c2=0: p3 attack, p4 defence; c2=1: p4 attack, p3 defence
   for c1 := 0; c1 < 2; c1++ {
       var t1atk, t1def int
       if c1 == 0 {
           t1atk = b[1]
           t1def = a[2]
       } else {
           t1atk = b[2]
           t1def = a[1]
       }
       for c2 := 0; c2 < 2; c2++ {
           var t2atk, t2def int
           if c2 == 0 {
               t2atk = b[3]
               t2def = a[4]
           } else {
               t2atk = b[4]
               t2def = a[3]
           }
           // determine result
           if t1def > t2atk && t1atk > t2def {
               payoff[c1][c2] = 1
           } else if t2def > t1atk && t2atk > t1def {
               payoff[c1][c2] = -1
           } else {
               payoff[c1][c2] = 0
           }
       }
   }
   // team2 response: for each c1, team2 picks c2 minimizing payoff
   r0 := payoff[0][0]
   if payoff[0][1] < r0 {
       r0 = payoff[0][1]
   }
   r1 := payoff[1][0]
   if payoff[1][1] < r1 {
       r1 = payoff[1][1]
   }
   // team1 picks c1 maximizing r
   result := r0
   if r1 > result {
       result = r1
   }
   switch result {
   case 1:
       fmt.Println("Team 1")
   case -1:
       fmt.Println("Team 2")
   default:
       fmt.Println("Draw")
   }
}
