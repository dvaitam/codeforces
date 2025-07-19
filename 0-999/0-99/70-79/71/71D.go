package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

type pair struct{ r, s int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   crank := "A23456789TJQK"
   suit := "CDHS"
   // deck availability
   inDeck := make([][]bool, 13)
   for i := range inDeck {
       inDeck[i] = make([]bool, 4)
       for j := 0; j < 4; j++ {
           inDeck[i][j] = true
       }
   }
   // read cards
   card := make([][]string, n)
   var jokers []int
   for i := 0; i < n; i++ {
       card[i] = make([]string, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &card[i][j])
           r, s := cardData(card[i][j], crank, suit)
           if r >= 0 {
               inDeck[r][s] = false
           } else {
               jokers = append(jokers, s)
           }
       }
   }
   // place function closure
   var placeJokers []int
   jreplace := [2][]pair{}
   place := func(y, x int) bool {
       used := make([]bool, 13)
       var localJokers []int
       for dy := 0; dy < 3; dy++ {
           for dx := 0; dx < 3; dx++ {
               cy, cx := y+dy, x+dx
               r, s := cardData(card[cy][cx], crank, suit)
               if r < 0 {
                   localJokers = append(localJokers, s)
               } else {
                   if used[r] {
                       return false
                   }
                   used[r] = true
               }
           }
       }
       // record jokers in this square
       for _, v := range localJokers {
           placeJokers = append(placeJokers, v)
       }
       // possible replacements
       for i := 0; i < 13; i++ {
           if !used[i] {
               for j := 0; j < 4; j++ {
                   if inDeck[i][j] && len(localJokers) > 0 {
                       t := localJokers[len(localJokers)-1]
                       jreplace[t] = append(jreplace[t], pair{i, j})
                       if len(localJokers) > 1 {
                           localJokers = localJokers[:len(localJokers)-1]
                           break
                       }
                   }
               }
           }
       }
       return true
   }
   // search two non-overlapping squares
   for y1 := 0; y1+3 <= n; y1++ {
       for x1 := 0; x1+3 <= m; x1++ {
           for y2 := 0; y2+3 <= n; y2++ {
               for x2 := 0; x2+3 <= m; x2++ {
                   // non-overlap
                   if y2+3 <= y1 || y1+3 <= y2 || x1+3 <= x2 || x2+3 <= x1 {
                       // reset
                       placeJokers = placeJokers[:0]
                       jreplace = [2][]pair{{}, {}}
                       if !place(y1, x1) {
                           continue
                       }
                       if !place(y2, x2) {
                           continue
                       }
                       // determine replacements
                       var rep [2]pair
                       // assign for jokers in squares
                       if len(placeJokers) == 1 {
                           t := placeJokers[0]
                           if len(jreplace[t]) == 0 {
                               continue
                           }
                           rep[t] = jreplace[t][0]
                       } else if len(placeJokers) == 2 {
                           if len(jreplace[0]) == 0 || len(jreplace[1]) == 0 {
                               continue
                           }
                           rep[0] = jreplace[0][0]
                           rep[1] = jreplace[1][0]
                           if rep[0] == rep[1] && len(jreplace[0]) > 1 {
                               rep[0] = jreplace[0][1]
                           }
                           if rep[0] == rep[1] && len(jreplace[1]) > 1 {
                               rep[1] = jreplace[1][1]
                           }
                       }
                       // assign remaining jokers if any
                       if len(jokers) > len(placeJokers) {
                           if len(placeJokers) == 1 {
                               other := 1 - placeJokers[0]
                               for i := 0; i < 13; i++ {
                                   for j := 0; j < 4; j++ {
                                       if inDeck[i][j] && (i != rep[placeJokers[0]].r || j != rep[placeJokers[0]].s) {
                                           rep[other] = pair{i, j}
                                       }
                                   }
                               }
                           } else {
                               cnt := 0
                               for i := 0; i < 13; i++ {
                                   for j := 0; j < 4; j++ {
                                       if inDeck[i][j] {
                                           rep[cnt] = pair{i, j}
                                           cnt = (cnt + 1) % 2
                                       }
                                   }
                               }
                           }
                       }
                       // output result
                       fmt.Fprintln(out, "Solution exists.")
                       if len(jokers) == 1 {
                           t := jokers[0]
                           p := rep[t]
                           fmt.Fprintf(out, "Replace J%d with %c%c.\n", t+1, crank[p.r], suit[p.s])
                       } else if len(jokers) == 2 {
                           fmt.Fprintf(out, "Replace J1 with %c%c and J2 with %c%c.\n", crank[rep[0].r], suit[rep[0].s], crank[rep[1].r], suit[rep[1].s])
                       } else {
                           fmt.Fprintln(out, "There are no jokers.")
                       }
                       fmt.Fprintf(out, "Put the first square to (%d, %d).\n", y1+1, x1+1)
                       fmt.Fprintf(out, "Put the second square to (%d, %d).\n", y2+1, x2+1)
                       return
                   }
               }
           }
       }
   }
   fmt.Fprintln(out, "No solution.")
}

func cardData(card string, crank, suit string) (int, int) {
   if card[0] == 'J' && card[1] >= '1' && card[1] <= '2' {
       return -1, int(card[1] - '1')
   }
   r := strings.IndexByte(crank, card[0])
   s := strings.IndexByte(suit, card[1])
   return r, s
}
