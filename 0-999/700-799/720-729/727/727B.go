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
   line, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       return
   }
   line = strings.TrimSpace(line)
   n := len(line)
   i := 0
   var totalCents int64
   for i < n {
       // skip name (letters)
       for i < n && line[i] >= 'a' && line[i] <= 'z' {
           i++
       }
       // parse price (digits and dots)
       start := i
       for i < n && ((line[i] >= '0' && line[i] <= '9') || line[i] == '.') {
           i++
       }
       price := line[start:i]
       // split dollars and cents
       idx := strings.LastIndex(price, ".")
       var dollarsStr, centsStr string
       if idx == -1 {
           dollarsStr = strings.ReplaceAll(price, ".", "")
       } else if len(price)-idx-1 == 2 {
           dollarsStr = strings.ReplaceAll(price[:idx], ".", "")
           centsStr = price[idx+1:]
       } else {
           dollarsStr = strings.ReplaceAll(price, ".", "")
       }
       // parse to numbers
       var dollars int64
       for _, ch := range dollarsStr {
           dollars = dollars*10 + int64(ch-'0')
       }
       var cents int64
       if len(centsStr) == 2 {
           cents = int64(centsStr[0]-'0')*10 + int64(centsStr[1]-'0')
       }
       totalCents += dollars*100 + cents
   }
   // format total
   totalDollars := totalCents / 100
   remCents := totalCents % 100
   ds := fmt.Sprintf("%d", totalDollars)
   // group every three digits
   var sb strings.Builder
   L := len(ds)
   first := L % 3
   if first == 0 {
       first = 3
   }
   sb.WriteString(ds[:first])
   for j := first; j < L; j += 3 {
       sb.WriteByte('.')
       sb.WriteString(ds[j : j+3])
   }
   if remCents > 0 {
       sb.WriteByte('.')
       sb.WriteString(fmt.Sprintf("%02d", remCents))
   }
   fmt.Println(sb.String())
}
