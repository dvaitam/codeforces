package main

import (
   "fmt"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   var HPY, ATKY, DEFY int
   var HPM, ATKM, DEFM int
   var costH, costA, costD int
   if _, err := fmt.Scan(&HPY, &ATKY, &DEFY); err != nil {
       return
   }
   fmt.Scan(&HPM, &ATKM, &DEFM)
   fmt.Scan(&costH, &costA, &costD)

   const INF = 1 << 60
   ans := INF

   // minimal attack boost to deal damage
   minAtkBoost := max(0, DEFM+1-ATKY)
   // reasonable upper bound for attack boost
   maxAtkBoost := HPM + DEFM + 1
   // maximal defense boost to potentially reduce incoming damage to zero
   maxDefBoost := max(0, ATKM-DEFY)

   for da := minAtkBoost; da <= maxAtkBoost; da++ {
       yATK := ATKY + da
       dmgM := yATK - DEFM
       if dmgM <= 0 {
           continue
       }
       // turns to kill monster
       turns := (HPM + dmgM - 1) / dmgM

       for dd := 0; dd <= maxDefBoost; dd++ {
           yDEF := DEFY + dd
           dmgY := ATKM - yDEF
           if dmgY < 0 {
               dmgY = 0
           }
           // compute needed HP
           var needHP int
           if dmgY == 0 {
               needHP = 1
           } else {
               needHP = turns*dmgY + 1
           }
           dh := max(0, needHP-HPY)

           cost := dh*costH + da*costA + dd*costD
           if cost < ans {
               ans = cost
           }
       }
   }
   fmt.Println(ans)
}
