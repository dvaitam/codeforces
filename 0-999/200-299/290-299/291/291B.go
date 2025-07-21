package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read the entire command line string, may contain spaces and quotes
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   // Trim trailing newline characters
   s = strings.TrimRight(s, "\r\n")
   var lexemes []string
   n := len(s)
   for i := 0; i < n; {
       // Skip spaces
       if s[i] == ' ' {
           i++
           continue
       }
       // Quoted lexeme
       if s[i] == '"' {
           i++
           start := i
           // Find closing quote
           for i < n && s[i] != '"' {
               i++
           }
           // Extract content inside quotes
           lexemes = append(lexemes, s[start:i])
           // Skip closing quote
           if i < n && s[i] == '"' {
               i++
           }
       } else {
           // Unquoted lexeme: read until next space
           start := i
           for i < n && s[i] != ' ' {
               i++
           }
           lexemes = append(lexemes, s[start:i])
       }
   }
   // Print lexemes with angle brackets
   for _, lex := range lexemes {
       fmt.Printf("<%s>\n", lex)
   }
}
