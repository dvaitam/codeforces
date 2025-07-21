package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   type pile struct {
       s       int
       front   int // sum when starting pick from front moves
       back    int // sum when starting pick from back moves
       diff    int // front - back
   }
   piles := make([]pile, n)
   totalCardsSum := 0
   for i := 0; i < n; i++ {
       var s int
       fmt.Fscan(reader, &s)
       cards := make([]int, s)
       for j := 0; j < s; j++ {
           fmt.Fscan(reader, &cards[j])
           totalCardsSum += cards[j]
       }
       // sum of front picks (odd local moves): sum of first (s+1)/2 cards
       // sum of back picks (even local moves): sum of last s/2 cards
       half := s / 2
       frontCnt := (s + 1) / 2
       sumFront, sumBack := 0, 0
       for j := 0; j < frontCnt; j++ {
           sumFront += cards[j]
       }
       for j := 0; j < half; j++ {
           sumBack += cards[s-1-j]
       }
       piles[i] = pile{s: s, front: sumFront, back: sumBack, diff: sumFront - sumBack}
   }
   // sort by absolute diff descending
   sort.Slice(piles, func(i, j int) bool {
       di := piles[i].diff
       dj := piles[j].diff
       if di < 0 {
           di = -di
       }
       if dj < 0 {
           dj = -dj
       }
       return di > dj
   })
   // simulate
   cielSum, jiroSum := 0, 0
   cielTurn := true
   for _, p := range piles {
       if cielTurn {
           cielSum += p.front
           jiroSum += p.back
       } else {
           cielSum += p.back
           jiroSum += p.front
       }
       // if pile size odd, turn flips after finishing pile
       if p.s%2 != 0 {
           cielTurn = !cielTurn
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprintf(writer, "%d %d", cielSum, jiroSum)
   writer.Flush()
}
