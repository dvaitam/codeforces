package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && len(input) == 0 {
       return
   }
   s := strings.TrimSpace(input)
   negative := false
   if strings.HasPrefix(s, "-") {
       negative = true
       s = s[1:]
   }
   // split integer and fractional parts
   parts := strings.SplitN(s, ".", 2)
   intPart := parts[0]
   fracPart := ""
   if len(parts) == 2 {
       fracPart = parts[1]
   }
   // format integer part with commas
   var b strings.Builder
   n := len(intPart)
   for i, ch := range intPart {
       if i > 0 && (n-i)%3 == 0 {
           b.WriteByte(',')
       }
       b.WriteRune(ch)
   }
   formattedInt := b.String()
   // format fractional part to exactly two digits
   if len(fracPart) < 2 {
       fracPart = fracPart + strings.Repeat("0", 2-len(fracPart))
   } else if len(fracPart) > 2 {
       fracPart = fracPart[:2]
   }
   // build result
   result := fmt.Sprintf("$%s.%s", formattedInt, fracPart)
   if negative {
       result = fmt.Sprintf("(%s)", result)
   }
   fmt.Fprintln(os.Stdout, result)
}
