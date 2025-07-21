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
   // read tokens (each "int" or "pair")
   tokens := make([]string, 0, n*2)
   for i := 0; ; i++ {
       var tok string
       if _, err := fmt.Fscan(reader, &tok); err != nil {
           break
       }
       tokens = append(tokens, tok)
   }
   m := len(tokens)
   // frame for pending pair, state: 0 = left child pending, 1 = right child pending
   type frame struct{ state int }
   frames := make([]frame, 0, n)
   var out strings.Builder
   for i, tok := range tokens {
       if tok == "int" {
           out.WriteString("int")
           // finish subtrees: close completed pair frames
           for len(frames) > 0 && frames[len(frames)-1].state == 1 {
               out.WriteByte('>')
               frames = frames[:len(frames)-1]
           }
           if len(frames) > 0 && frames[len(frames)-1].state == 0 {
               // after left child, write comma and expect right child
               out.WriteByte(',')
               frames[len(frames)-1].state = 1
           } else if len(frames) == 0 && i != m-1 {
               // complete type before consuming all tokens
               fmt.Println("Error occurred")
               return
           }
       } else if tok == "pair" {
           out.WriteString("pair<")
           frames = append(frames, frame{state: 0})
       } else {
           // invalid token
           fmt.Println("Error occurred")
           return
       }
   }
   // if any pending frames remain, error
   if len(frames) != 0 {
       fmt.Println("Error occurred")
       return
   }
   fmt.Println(out.String())
}
