package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   total := 2 * n
   p := make([]int, total+1)
   for i := 1; i <= total; i++ {
       fmt.Fscan(reader, &p[i])
   }
   partner := make([]int, total+1)
   pairs := make([][2]int, 0, m)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       partner[a] = b
       partner[b] = a
       pairs = append(pairs, [2]int{a, b})
   }
   var t int
   fmt.Fscan(reader, &t)
   juryPrio := make([]int, total)
   for i := 0; i < total; i++ {
       fmt.Fscan(reader, &juryPrio[i])
   }

   used := make([]bool, total+1)
   forcedForMe := 0
   forcedForOpp := 0
   myTurn := (t == 1)

   // simulate 2n moves
   for moves := 0; moves < total; moves++ {
       if myTurn {
           var pick int
           // forced pick
           if forcedForMe != 0 && !used[forcedForMe] {
               pick = forcedForMe
               forcedForMe = 0
           } else {
               // choose best remaining pair diff
               bestDiff := -1
               bestPick := 0
               for _, pr := range pairs {
                   a, b := pr[0], pr[1]
                   if !used[a] && !used[b] {
                       d := abs(p[a] - p[b])
                       if d > bestDiff {
                           bestDiff = d
                           if p[a] > p[b] {
                               bestPick = a
                           } else {
                               bestPick = b
                           }
                       }
                   }
               }
               if bestDiff >= 0 {
                   pick = bestPick
               } else {
                   // pick highest power among unused
                   maxP := -1
                   for i := 1; i <= total; i++ {
                       if !used[i] {
                           if p[i] > maxP {
                               maxP = p[i]
                               pick = i
                           }
                       }
                   }
               }
           }
           // make pick
           used[pick] = true
           fmt.Fprintln(writer, pick)
           // set forced for opponent if pair
           if partner[pick] != 0 && !used[partner[pick]] {
               forcedForOpp = partner[pick]
           }
       } else {
           var pick int
           if forcedForOpp != 0 && !used[forcedForOpp] {
               pick = forcedForOpp
               forcedForOpp = 0
           } else {
               for _, h := range juryPrio {
                   if !used[h] {
                       pick = h
                       break
                   }
               }
           }
           used[pick] = true
           // set forced for me if pair
           if partner[pick] != 0 && !used[partner[pick]] {
               forcedForMe = partner[pick]
           }
       }
       myTurn = !myTurn
   }
}
