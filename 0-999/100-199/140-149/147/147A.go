package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
   "unicode"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       // ignore error
   }
   line = strings.TrimRight(line, "\n")
   var tokens []string
   var curr strings.Builder
   for _, r := range line {
       if unicode.IsLetter(r) {
           curr.WriteRune(r)
       } else if r == '.' || r == ',' || r == '!' || r == '?' {
           if curr.Len() > 0 {
               tokens = append(tokens, curr.String())
               curr.Reset()
           }
           tokens = append(tokens, string(r))
       } else {
           if curr.Len() > 0 {
               tokens = append(tokens, curr.String())
               curr.Reset()
           }
           // skip spaces
       }
   }
   if curr.Len() > 0 {
       tokens = append(tokens, curr.String())
   }
   var out strings.Builder
   for i, tok := range tokens {
       // punctuation token
       if len(tok) == 1 && (tok == "." || tok == "," || tok == "!" || tok == "?") {
           out.WriteString(tok)
           out.WriteByte(' ')
       } else {
           out.WriteString(tok)
           // add space between words
           if i+1 < len(tokens) {
               next := tokens[i+1]
               if !(len(next) == 1 && (next == "." || next == "," || next == "!" || next == "?")) {
                   out.WriteByte(' ')
               }
           }
       }
   }
   fmt.Print(out.String())
}
