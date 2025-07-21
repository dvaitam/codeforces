package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

type block struct {
   tryLine   int
   catchLine int
   exType    string
   message   string
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   _, err := fmt.Fscanf(reader, "%d\n", &n)
   if err != nil {
       return
   }
   scanner := bufio.NewScanner(reader)
   tryStack := make([]int, 0, n)
   blocks := make([]block, 0, n/2)
   throwLine := -1
   throwType := ""
   lineNum := 0
   for scanner.Scan() {
       line := scanner.Text()
       lineNum++
       s := strings.TrimSpace(line)
       if len(s) == 0 {
           continue
       }
       if strings.HasPrefix(s, "try") && (s == "try" || strings.HasPrefix(s, "try")) {
           tryStack = append(tryStack, lineNum)
       } else if strings.HasPrefix(s, "catch") {
           // parse catch(type, "message")
           p1 := strings.Index(s, "(")
           p2 := strings.LastIndex(s, ")")
           inner := s[p1+1 : p2]
           // split by first comma
           ci := strings.Index(inner, ",")
           typ := strings.TrimSpace(inner[:ci])
           msgPart := strings.TrimSpace(inner[ci+1:])
           // remove surrounding quotes
           if len(msgPart) >= 2 && msgPart[0] == '"' && msgPart[len(msgPart)-1] == '"' {
               msgPart = msgPart[1 : len(msgPart)-1]
           }
           // pop last try
           t := tryStack[len(tryStack)-1]
           tryStack = tryStack[:len(tryStack)-1]
           blocks = append(blocks, block{tryLine: t, catchLine: lineNum, exType: typ, message: msgPart})
       } else if strings.HasPrefix(s, "throw") {
           // parse throw(type)
           p1 := strings.Index(s, "(")
           p2 := strings.LastIndex(s, ")")
           inner := s[p1+1 : p2]
           tm := strings.TrimSpace(inner)
           throwType = tm
           throwLine = lineNum
       }
   }
   // find matching block
   bestCatch := n + 1
   bestMsg := ""
   for _, b := range blocks {
       if b.tryLine < throwLine && throwLine < b.catchLine && b.exType == throwType {
           if b.catchLine < bestCatch {
               bestCatch = b.catchLine
               bestMsg = b.message
           }
       }
   }
   if bestMsg == "" {
       fmt.Println("Unhandled Exception")
   } else {
       fmt.Println(bestMsg)
   }
}
