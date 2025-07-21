package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

// Widget represents a GUI widget, which may be simple or a container (HBox/VBox)
type Widget struct {
   kind     int      // 0=simple, 1=HBox, 2=VBox
   w, h     int      // for simple widget
   border   int      // for containers
   spacing  int      // for containers
   children []string // names of packed widgets
}

var widgets map[string]*Widget
var memoW, memoH map[string]int
var seen map[string]bool

// compute returns width and height of widget named name
func compute(name string) (int, int) {
   if w, ok := memoW[name]; ok {
       return w, memoH[name]
   }
   wgt := widgets[name]
   var w, h int
   switch wgt.kind {
   case 0:
       w, h = wgt.w, wgt.h
   case 1, 2:
       n := len(wgt.children)
       if n == 0 {
           w, h = 0, 0
       } else {
           // compute child sizes
           var sumW, sumH, maxW, maxH int
           for _, ch := range wgt.children {
               cw, chh := compute(ch)
               sumW += cw
               sumH += chh
               if cw > maxW {
                   maxW = cw
               }
               if chh > maxH {
                   maxH = chh
               }
           }
           // apply spacing and border
           if wgt.kind == 1 {
               // HBox: horizontal
               w = sumW + wgt.spacing*(n-1) + 2*wgt.border
               h = maxH + 2*wgt.border
           } else {
               // VBox: vertical
               w = maxW + 2*wgt.border
               h = sumH + wgt.spacing*(n-1) + 2*wgt.border
           }
       }
   }
   memoW[name] = w
   memoH[name] = h
   return w, h
}

func main() {
   in := bufio.NewScanner(os.Stdin)
   in.Split(bufio.ScanLines)
   // read n
   in.Scan()
   n, _ := strconv.Atoi(strings.TrimSpace(in.Text()))
   widgets = make(map[string]*Widget)
   // parse instructions
   for i := 0; i < n; i++ {
       in.Scan()
       line := in.Text()
       switch {
       case strings.HasPrefix(line, "Widget "):
           // Widget name(x,y)
           rest := line[len("Widget "):]
           idx := strings.Index(rest, "(")
           name := rest[:idx]
           nums := rest[idx+1 : len(rest)-1]
           parts := strings.Split(nums, ",")
           x, _ := strconv.Atoi(parts[0])
           y, _ := strconv.Atoi(parts[1])
           widgets[name] = &Widget{kind: 0, w: x, h: y}
       case strings.HasPrefix(line, "HBox "):
           name := line[len("HBox "):]
           widgets[name] = &Widget{kind: 1}
       case strings.HasPrefix(line, "VBox "):
           name := line[len("VBox "):]
           widgets[name] = &Widget{kind: 2}
       case strings.Contains(line, ".pack("):
           idx := strings.Index(line, ".pack(")
           parent := line[:idx]
           child := line[idx+6 : len(line)-1]
           wgt := widgets[parent]
           wgt.children = append(wgt.children, child)
       case strings.Contains(line, ".set_border("):
           idx := strings.Index(line, ".set_border(")
           name := line[:idx]
           val, _ := strconv.Atoi(line[idx+12 : len(line)-1])
           widgets[name].border = val
       case strings.Contains(line, ".set_spacing("):
           idx := strings.Index(line, ".set_spacing(")
           name := line[:idx]
           val, _ := strconv.Atoi(line[idx+13 : len(line)-1])
           widgets[name].spacing = val
       }
   }
   // compute sizes
   memoW = make(map[string]int)
   memoH = make(map[string]int)
   // output sorted by name
   names := make([]string, 0, len(widgets))
   for name := range widgets {
       names = append(names, name)
   }
   sort.Strings(names)
   for _, name := range names {
       w, h := compute(name)
       fmt.Printf("%s %d %d\n", name, w, h)
   }
}
