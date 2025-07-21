package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // skip the rest of the line after n
   _, _ = in.ReadString('\n')
   // read the text line
   text, _ := in.ReadString('\n')
   text = strings.TrimRight(text, "\r\n")
   // parse sentences separated by a punctuation (.,!,?) and a space
   var sentences []string
   start := 0
   for i := 0; i < len(text); i++ {
       c := text[i]
       if (c == '.' || c == '!' || c == '?') && i+1 < len(text) && text[i+1] == ' ' {
           sentences = append(sentences, text[start:i+1])
           start = i + 2
           i++
       }
   }
   if start < len(text) {
       sentences = append(sentences, text[start:])
   }
   // greedy pack sentences into messages of size <= n
   msgs := 0
   curr := 0
   for _, s := range sentences {
       L := len(s)
       if L > n {
           fmt.Println("Impossible")
           return
       }
       if curr == 0 {
           curr = L
       } else if curr+1+L <= n {
           curr += 1 + L
       } else {
           msgs++
           curr = L
       }
   }
   if curr > 0 {
       msgs++
   }
   fmt.Println(msgs)
}
