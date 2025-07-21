package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

var (
   s   string
   pos int
   cnt map[string]int
   ans int64
)

// parse reads an employee description starting at s[pos]
func parse() {
   // read name
   start := pos
   for pos < len(s) {
       c := s[pos]
       if c < 'A' || c > 'Z' {
           break
       }
       pos++
   }
   name := s[start:pos]
   // count uncomfortable situations for this name
   ans += int64(cnt[name])
   // add this name to ancestors
   cnt[name]++
   // parse subordinates if any
   if pos < len(s) && s[pos] == ':' {
       pos++ // skip ':'
       for {
           parse()
           if pos < len(s) && s[pos] == ',' {
               pos++ // skip ',' and continue
               continue
           }
           break
       }
   }
   // skip terminating '.'
   if pos < len(s) && s[pos] == '.' {
       pos++
   }
   // remove this name from ancestors
   cnt[name]--
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   s = strings.TrimSpace(line)
   cnt = make(map[string]int)
   ans = 0
   pos = 0
   parse()
   fmt.Println(ans)
}
