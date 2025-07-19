package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       cnt := map[rune]int{'R': 0, 'P': 0, 'S': 0}
       for _, ch := range s {
           cnt[ch]++
       }
       var ans strings.Builder
       // All moves have equal count
       if cnt['R'] == cnt['P'] && cnt['P'] == cnt['S'] {
           for _, ch := range s {
               switch ch {
               case 'R':
                   ans.WriteRune('P')
               case 'P':
                   ans.WriteRune('S')
               case 'S':
                   ans.WriteRune('R')
               }
           }
       } else {
           // choose the most frequent move and play winning move
           cmx := 'R'
           maxc := cnt['R']
           if cnt['P'] > maxc {
               cmx = 'P'
               maxc = cnt['P']
           }
           if cnt['S'] > maxc {
               cmx = 'S'
               maxc = cnt['S']
           }
           var play rune
           switch cmx {
           case 'R':
               play = 'P'
           case 'P':
               play = 'S'
           case 'S':
               play = 'R'
           }
           for i := 0; i < len(s); i++ {
               ans.WriteRune(play)
           }
       }
       fmt.Fprintln(writer, ans.String())
   }
}
