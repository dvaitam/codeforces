package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   k := len(s)
   if k%2 != 0 {
       fmt.Println("-1")
       return
   }
   n := k / 2
   // seven-segment bit masks: bit0=a,1=b,2=c,3=d,4=e,5=f,6=g
   masks := [10]int{
       0b0111111, // 0: a b c d e f
       0b0000110, // 1: b c
       0b1011011, // 2: a b d e g
       0b1001111, // 3: a b c d g
       0b1100110, // 4: f g b c
       0b1101101, // 5: a f g c d
       0b1111101, // 6: a f e d c g
       0b0000111, // 7: a b c
       0b1111111, // 8: a b c d e f g
       0b1101111, // 9: a b c d f g
   }
   // overlap count and best opportunities
   var overlap [10][10]int
   bestOpp := [10]int{}
   M := 0
   for i := 0; i < 10; i++ {
       for j := 0; j < 10; j++ {
           cnt := bitsCount(masks[i] & masks[j])
           overlap[i][j] = cnt
           if cnt > M {
               M = cnt
           }
           if cnt > bestOpp[i] {
               bestOpp[i] = cnt
           }
       }
   }
   // original degree
   sdig := make([]int, k)
   for i := 0; i < k; i++ {
       sdig[i] = int(s[i] - '0')
   }
   D0 := 0
   for i := 0; i < n; i++ {
       D0 += overlap[sdig[i]][sdig[i+n]]
   }
   // state per pair
   state := make([]uint8, n) // bit0:first fixed, bit1:second fixed
   firstD := make([]int, n)
   secondD := make([]int, n)
   sumBoth := 0
   sumSingle := 0
   freePairs := n
   // result digits
   T := make([]byte, k)
   // helper to test candidate at pos
   test := func(pos, d int) bool {
       i := pos
       side := 0
       if pos >= n {
           i = pos - n
           side = 1
       }
       st := state[i]
       // simulate deltas
       deltaBoth, deltaSingle, deltaFree := 0, 0, 0
       if st == 0 {
           // from free to single
           deltaSingle = bestOpp[d]
           deltaFree = -1
       } else if st == 1 && side == 1 {
           // first fixed, now second
           od := firstD[i]
           deltaSingle = -bestOpp[od]
           deltaBoth = overlap[od][d]
       } else if st == 2 && side == 0 {
           // second fixed, now first
           od := secondD[i]
           deltaSingle = -bestOpp[od]
           deltaBoth = overlap[d][od]
       }
       maxPossible := sumBoth + deltaBoth + sumSingle + deltaSingle + (freePairs+deltaFree)*M
       return maxPossible > D0
   }
   // helper to apply candidate permanently
   apply := func(pos, d int) {
       i := pos
       side := 0
       if pos >= n {
           i = pos - n
           side = 1
       }
       st := state[i]
       if st == 0 {
           sumSingle += bestOpp[d]
           freePairs--
           if side == 0 {
               state[i] = 1
               firstD[i] = d
           } else {
               state[i] = 2
               secondD[i] = d
           }
       } else if st == 1 && side == 1 {
           od := firstD[i]
           sumSingle -= bestOpp[od]
           sumBoth += overlap[od][d]
           state[i] = 3
           secondD[i] = d
       } else if st == 2 && side == 0 {
           od := secondD[i]
           sumSingle -= bestOpp[od]
           sumBoth += overlap[d][od]
           state[i] = 3
           firstD[i] = d
       }
       // set in result
       T[pos] = byte('0' + d)
   }
   // greedy prefix
   pivot := -1
   for idx := 0; idx < k; idx++ {
       sd := sdig[idx]
       // try greater digits
       for d := sd + 1; d < 10; d++ {
           if test(idx, d) {
               apply(idx, d)
               pivot = idx
               break
           }
       }
       if pivot == idx {
           break
       }
       // else must match original to continue
       if !test(idx, sd) {
           // cannot even match, no solution
           fmt.Println("-1")
           return
       }
       apply(idx, sd)
   }
   if pivot < 0 {
       fmt.Println("-1")
       return
   }
   // fill suffix
   for idx := pivot + 1; idx < k; idx++ {
       for d := 0; d < 10; d++ {
           if test(idx, d) {
               apply(idx, d)
               break
           }
       }
   }
   fmt.Println(string(T))
}

func bitsCount(x int) int {
   // builtin popcount
   return int((uint(x) * 0x200040008001 & 0x111111111111111) % 0xf)
}
