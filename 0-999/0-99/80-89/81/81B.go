package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       return
   }
   s = strings.TrimRight(s, "\r\n")

   // tokenize: numbers, commas, three dots
   // token types
   const (
       TOK_NUM = iota
       TOK_COMMA
       TOK_DOTS
   )
   var tokTypes []int
   var tokVals []string
   n := len(s)
   for i := 0; i < n; {
       switch {
       case s[i] >= '0' && s[i] <= '9':
           j := i
           for j < n && s[j] >= '0' && s[j] <= '9' {
               j++
           }
           tokTypes = append(tokTypes, TOK_NUM)
           tokVals = append(tokVals, s[i:j])
           i = j
       case i+2 < n && s[i] == '.' && s[i+1] == '.' && s[i+2] == '.':
           tokTypes = append(tokTypes, TOK_DOTS)
           tokVals = append(tokVals, "...")
           i += 3
       case s[i] == ',':
           tokTypes = append(tokTypes, TOK_COMMA)
           tokVals = append(tokVals, ",")
           i++
       case s[i] == ' ':
           // skip spaces
           i++
       default:
           // unexpected char, skip
           i++
       }
   }

   var b strings.Builder
   m := len(tokTypes)
   for i, t := range tokTypes {
       switch t {
       case TOK_NUM:
           if i > 0 && tokTypes[i-1] == TOK_NUM {
               b.WriteByte(' ')
           }
           b.WriteString(tokVals[i])
       case TOK_COMMA:
           b.WriteString(",")
           if i < m-1 {
               b.WriteByte(' ')
           }
       case TOK_DOTS:
           if b.Len() > 0 {
               // ensure exactly one space before dots
               if last := b.String(); last[len(last)-1] != ' ' {
                   b.WriteByte(' ')
               }
           }
           b.WriteString("...")
       }
   }
   fmt.Print(b.String())
}
