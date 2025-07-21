package main

import (
   "bufio"
   "fmt"
   "os"
)

type card struct {
   idx  int
   c    int
   word string
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   // skills: a[i][j], b[i][j], i in [0,n), j in [0,2)
   a := make([][2]int, n)
   b := make([][2]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i][0], &b[i][0], &a[i][1], &b[i][1])
   }
   var m int
   fmt.Fscan(reader, &m)
   cards := make([]card, 0, m)
   for k := 0; k < m; k++ {
       var w string
       var c int
       fmt.Fscan(reader, &w)
       fmt.Fscan(reader, &c)
       cards = append(cards, card{idx: k, c: c, word: w})
   }
   // per team and card spent time
   d := make([][]int, n)
   for i := 0; i < n; i++ {
       d[i] = make([]int, m)
   }
   // results per team
   res := make([][]string, n)

   // deck queue
   deck := cards
   // current player pointer 0..2*n-1
   p := 0
   for len(deck) > 0 {
       timeLeft := t
       team := p % n
       player := 0
       if p >= n {
           player = 1
       }
       teammate := 1 - player
       // player's turn
       for timeLeft > 0 && len(deck) > 0 {
           cur := deck[0]
           deck = deck[1:]
           k := cur.idx
           need := cur.c - (a[team][player] + b[team][teammate]) - d[team][k]
           if need < 1 {
               need = 1
           }
           if need <= timeLeft {
               timeLeft -= need
               res[team] = append(res[team], cur.word)
               // card removed, do not re-add
           } else {
               // fail to guess, spend remaining time
               d[team][k] += timeLeft
               deck = append(deck, cur)
               timeLeft = 0
           }
       }
       p++
       if p == 2*n {
           p = 0
       }
   }
   // output
   for i := 0; i < n; i++ {
       fmt.Fprint(writer, len(res[i]))
       for _, w := range res[i] {
           fmt.Fprint(writer, " ", w)
       }
       fmt.Fprintln(writer)
   }
}
