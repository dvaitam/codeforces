package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var k1 int
   fmt.Fscan(reader, &k1)
   p1 := make([]int, k1)
   for i := 0; i < k1; i++ {
       fmt.Fscan(reader, &p1[i])
   }
   var k2 int
   fmt.Fscan(reader, &k2)
   p2 := make([]int, k2)
   for i := 0; i < k2; i++ {
       fmt.Fscan(reader, &p2[i])
   }

   seen := make(map[string]struct{})
   rounds := 0
   for {
       // check for infinite loop
       state := encode(p1, p2)
       if _, ok := seen[state]; ok {
           fmt.Println(-1)
           return
       }
       seen[state] = struct{}{}

       if len(p1) == 0 {
           fmt.Printf("%d 2\n", rounds)
           return
       }
       if len(p2) == 0 {
           fmt.Printf("%d 1\n", rounds)
           return
       }

       // play one fight
       rounds++
       c1 := p1[0]
       c2 := p2[0]
       p1 = p1[1:]
       p2 = p2[1:]
       if c1 > c2 {
           // player1 wins: take opponent's card then own card
           p1 = append(p1, c2, c1)
       } else {
           p2 = append(p2, c1, c2)
       }
   }
}

func encode(a, b []int) string {
   // encode state as string
   var sb strings.Builder
   for i, v := range a {
       if i > 0 {
           sb.WriteByte(',')
       }
       sb.WriteString(fmt.Sprint(v))
   }
   sb.WriteByte('|')
   for i, v := range b {
       if i > 0 {
           sb.WriteByte(',')
       }
       sb.WriteString(fmt.Sprint(v))
   }
   return sb.String()
}
