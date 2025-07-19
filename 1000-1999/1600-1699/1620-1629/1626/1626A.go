package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var a string
       fmt.Fscan(reader, &a)
       var book [256]int
       var nex []byte
       for i := 0; i < len(a); i++ {
           b := a[i]
           book[b]++
           if book[b] == 2 {
               nex = append(nex, b)
           }
       }
       // build answer: repeat nex twice, then leftover unique letters
       var ans []byte
       ans = append(ans, nex...)
       ans = append(ans, nex...)
       for ch := byte('a'); ch <= byte('z'); ch++ {
           if book[ch] == 1 {
               ans = append(ans, ch)
           }
       }
       fmt.Println(string(ans))
   }
}
