package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
2 5
5-b 3-e
4-e 1-e 1-d 3-e 2-b
4 5
5-d 4-b 2-b 5-d
1-a 2-e 1-c 1-c 4-e
4 4
4-e 4-b 3-a 1-b
4-b 3-d 3-d 5-d
5 3
5-e 4-e 2-c 1-c 5-b
3-e 5-e 1-b
5 3
3-a 1-d 4-a 3-a 4-b
1-c 4-d 1-a
5 5
1-d 5-c 5-c 5-b 1-c
1-a 1-e 5-a 2-d 3-e
3 2
1-c 3-c 2-d
4-d 5-d
5 5
1-e 5-c 4-b 3-d 3-e
3-e 3-a 4-e 3-a 4-e
5 2
1-c 4-c 3-e 3-d 1-e
1-a 3-c
4 3
5-e 3-b 3-b 3-c
5-c 3-d 1-a
5 2
3-e 2-c 2-c 2-d 1-a
5-c 3-b
4 2
1-c 2-e 4-c 2-a
1-e 2-c
5 2
3-c 1-e 3-e 2-d 3-e
3-d 3-d
1 5
1-d
2-e 1-c 5-d 5-d 2-a
5 2
5-a 3-c 2-d 1-d 2-a
1-c 2-a
4 5
3-a 3-c 3-e 2-e
3-a 1-e 5-c 5-a 3-e
1 1
4-b
4-c
5 5
4-c 5-d 5-c 3-b 1-b
4-b 4-e 2-e 5-e 3-a
5 5
4-d 3-d 4-c 2-b 3-c
4-d 5-e 1-a 3-b 1-e
1 5
5-d
1-e 4-c 2-c 5-b 3-a
2 1
5-a 4-c
3-d
4 5
4-b 2-d 1-b 5-e
3-e 3-a 4-b 2-c 1-a
3 2
2-b 4-b 5-e
1-d 1-d
2 5
4-d 1-c
5-a 5-a 2-d 3-c 5-c
5 2
1-c 3-c 5-e 5-b 5-d
2-a 4-d
2 2
4-b 3-e
1-b 3-b
1 4
2-a
2-d 1-e 5-d 2-d
2 5
1-e 5-e
2-c 2-e 4-b 2-c 1-d
1 3
4-c
3-d 5-e 4-e
3 3
5-b 3-e 4-a
1-d 1-e 4-d
4 2
5-c 3-b 1-b 4-c
4-a 2-c
1 4
3-b
1-a 4-a 1-e 1-d
4 5
4-d 2-c 5-d 5-b
3-e 5-d 2-a 3-c 3-b
1 3
4-e
4-d 1-b 2-a
5 2
1-e 2-d 2-e 5-b 2-b
5-a 3-e
5 5
4-e 5-c 3-e 1-c 1-b
2-a 5-c 2-d 4-b 4-d
3 4
1-d 2-c 2-b
3-e 1-d 4-a 2-e
1 5
4-b
1-c 2-c 5-b 5-e 3-c
4 4
1-c 1-a 5-b 4-b
3-e 1-d 1-c 1-d
1 5
5-c
2-c 2-e 2-c 1-a 1-d
1 5
5-d
2-a 4-d 2-e 1-d 1-c
2 1
2-d 5-a
5-c
2 1
3-c 4-c
1-e
5 5
5-d 5-d 5-b 4-d 2-b
3-a 4-e 1-e 5-c 4-d
1 2
5-e
2-b 5-e
2 5
5-a 4-a
5-c 2-a 3-e 1-b 5-d
1 1
3-e
5-c
3 2
4-a 2-d 2-a
1-c 1-a
5 1
4-c 1-e 4-a 2-e 3-a
4-d
2 5
1-b 5-c
4-b 1-e 2-c 5-d 4-b
2 4
3-e 4-a
5-d 2-d 2-d 4-c
1 3
2-e
1-c 3-c 5-a
3 5
2-b 4-d 5-d
5-b 2-e 1-e 2-d 3-b
5 2
5-d 2-a 2-b 4-b 1-b
3-d 3-e
2 5
5-b 2-d
3-e 1-c 3-a 5-d 2-b
5 3
4-a 5-c 3-b 5-b 1-a
2-d 4-b 1-a
4 3
1-c 4-a 1-d 4-b
3-e 5-d 2-d
2 5
4-d 3-b
3-d 1-b 5-e 2-d 5-d
4 4
2-a 3-d 4-b 5-a
5-b 3-a 3-b 1-c
2 1
4-b 1-d
4-b
4 2
2-b 4-c 4-c 2-a
2-e 2-e
5 4
2-b 5-a 5-a 3-c 3-d
1-c 2-c 3-a 3-a
2 4
3-e 1-c
2-d 2-a 5-d 1-b
3 1
5-c 4-b 5-e
2-d
5 3
1-c 4-e 3-e 5-a 5-d
5-b 4-d 3-b
4 1
2-d 1-d 1-a 3-d
3-d
3 5
5-e 3-e 4-c
4-e 5-c 1-d 4-a 5-a
2 5
1-c 2-a
4-c 1-e 1-b 2-e 3-e
2 4
3-d 1-b
5-a 1-e 4-d 5-a
2 5
4-b 3-d
5-c 1-d 4-e 4-a 3-c
1 2
3-a
5-e 2-e
2 3
3-b 4-e
1-a 2-a 4-d
2 3
4-c 4-e
4-c 4-b 3-c
5 1
4-b 5-c 2-e 5-a 4-d
4-a
5 3
1-d 4-e 4-d 1-c 3-c
2-a 1-b 1-b
3 4
2-b 4-b 1-b
5-e 5-a 3-a 3-c
4 4
1-a 3-b 2-a 4-a
4-c 2-a 4-c 3-c
1 2
4-b
1-e 3-e
3 1
3-b 1-b 4-b
2-b
4 1
5-d 5-e 5-a 3-b
4-d
5 2
3-e 2-e 5-c 5-a 1-d
1-e 2-c
2 1
1-d 4-e
5-b
4 5
5-b 2-b 2-c 2-a
2-c 4-b 5-c 3-d 4-b
5 2
5-a 2-d 5-a 2-d 3-a
1-e 1-b
2 3
5-c 3-c
4-a 2-c 4-d
5 4
1-e 5-e 3-b 3-a 5-e
1-e 2-c 2-d 3-e
2 2
1-d 3-d
5-e 4-c
2 5
5-d 4-b
4-e 1-b 3-e 5-a 1-b
3 1
5-e 2-d 5-c
2-e
2 5
4-d 3-b
4-c 1-a 2-a 5-e 2-d
1 1
4-e
4-d
5 5
5-b 4-a 1-a 2-a 3-c
1-b 2-e 5-b 5-b 1-c
5 4
4-b 2-c 5-b 2-c 3-e
3-d 5-a 5-d 4-b
3 1
4-e 1-d 3-e
4-b
2 5
2-c 5-d
2-e 1-c 5-e 2-b 3-a
4 4
5-a 5-d 3-a 2-b
2-c 2-d 5-c 4-e
5 1
1-a 3-d 2-c 4-a 3-d
5-a
1 2
3-a
5-d 3-b
5 1
1-a 2-b 4-d 2-a 1-b
3-d
5 1
3-b 2-d 5-d 4-e 1-c
2-d
3 5
4-d 1-b 5-a
1-c 3-b 3-a 1-d 4-c
2 3
1-c 1-c
3-c 5-e 4-c
1 5
4-a
1-e 3-d 4-b 2-c 2-c
4 3
2-a 1-a 2-c 2-c
5-c 4-e 5-d
2 4
5-e 1-b
1-b 5-a 1-a 2-a
4 4
3-e 3-d 2-a 3-a
1-e 4-e 2-c 3-a
4 4
4-e 5-e 5-d 2-c
3-e 4-e 4-e 1-d
5 4
1-d 4-e 1-b 4-b 5-d
1-b 3-e 4-d 4-b
5 4
5-c 2-b 1-b 1-b 2-a
4-b 1-e 5-e 5-b
1 1
4-a
4-b
2 1
3-d 1-d
1-e
4 5
3-b 3-b 5-d 3-a
1-e 5-a 2-c 4-b 2-b
5 5
2-a 1-c 2-d 1-c 1-a
3-a 3-c 3-d 4-b 1-c
3 3
4-a 1-e 4-c
5-e 2-b 4-d
1 1
3-e
1-e
3 3
5-b 2-b 5-d
5-e 3-e 1-a
4 1
1-a 4-e 1-a 3-e
2-a
5 3
5-b 4-a 1-b 5-c 5-c
5-e 5-b 2-d
2 3
1-d 5-e
1-a 1-e 2-e
4 4
2-d 5-a 5-b 4-d
4-e 2-b 1-d 3-c
3 1
2-b 2-a 4-e
2-e
3 4
4-e 5-d 2-a
4-c 4-e 3-b 5-a
1 5
5-d
4-a 3-d 5-d 1-a 1-c
3 3
5-a 5-d 1-c
3-a 1-c 1-c
5 5
3-b 1-e 3-c 2-e 3-b
5-e 2-b 1-a 4-e 4-b
3 5
4-b 1-d 4-b
3-b 4-e 3-e 4-d 5-b
3 3
2-b 5-e 1-b
1-e 3-e 3-b
1 5
3-d
5-c 3-c 5-b 3-a 5-c
2 2
3-e 2-c
4-b 3-b
2 2
2-e 5-c
3-d 1-d
3 5
5-e 2-c 4-e
4-a 1-b 1-d 1-a 5-b
4 4
3-d 4-e 2-e 2-d
2-d 1-d 3-b 2-c
5 1
3-e 5-b 2-d 5-a 5-d
2-b
5 1
2-c 5-e 4-e 3-a 2-a
2-a
5 1
2-b 5-e 1-a 4-d 3-b
2-b
5 2
1-d 4-a 4-d 5-c 3-a
3-e 2-d
3 4
3-a 3-c 3-b
5-e 1-c 4-c 2-e
4 1
2-a 1-d 2-d 2-d
1-c
4 1
5-a 2-d 4-c 2-e
2-e
1 1
2-b
5-d
2 5
4-a 4-a
5-e 4-a 3-c 3-a 1-c
5 5
1-c 1-e 3-b 2-a 5-d
5-c 5-c 1-d 5-d 4-a
3 5
1-d 3-a 3-e
1-e 4-e 4-a 5-c 1-d
2 4
1-c 4-d
2-a 4-e 1-a 5-e
4 1
1-a 4-b 3-e 4-b
5-c
1 2
3-d
5-d 1-a
4 5
1-c 1-b 3-e 3-b
3-c 5-e 3-d 2-b 5-a
4 2
2-e 5-e 1-d 5-e
1-e 5-b
2 2
2-c 5-c
3-e 4-e
3 1
2-a 5-e 3-c
5-b
2 1
1-a 5-d
1-c
3 4
4-d 3-e 5-c
3-d 3-d 3-a 5-b
3 1
4-e 2-d 2-a
3-b
1 2
1-b
2-c 3-c
3 3
1-a 1-e 5-e
5-a 5-b 3-b
4 3
5-a 3-c 1-b 1-c
1-e 3-c 2-b
3 3
3-e 1-a 5-d
4-b 1-d 1-b
4 3
3-e 3-a 3-d 5-e
3-a 1-b 3-c
1 3
4-b
1-d 1-b 1-c
2 1
1-e 3-d
5-a
1 4
5-e
3-b 1-e 4-a 1-d
4 4
2-d 1-b 3-c 3-e
2-b 5-e 4-b 1-c
4 3
5-a 5-a 3-a 4-e
2-a 3-c 2-e
2 2
1-e 3-b
1-d 3-d
5 4
4-b 3-c 5-d 2-a 5-d
1-e 4-a 2-c 3-c
3 2
2-c 2-b 5-b
4-b 2-b
1 4
5-c
1-a 1-e 4-d 2-b
5 4
4-a 3-a 1-d 5-d 1-d
1-d 4-d 4-a 4-b
2 4
3-e 5-c
5-c 4-d 3-a 1-d
2 3
1-c 3-b
2-d 4-c 3-b
3 2
1-a 5-e 3-e
5-c 5-a
5 1
4-b 4-b 3-a 4-e 4-d
5-c
1 1
1-c
4-e
5 4
1-b 2-b 1-d 1-c 5-a
4-e 1-a 3-d 3-b
2 2
5-e 2-a
2-d 1-b
4 2
5-d 1-a 4-c 2-e
3-e 5-b
5 1
5-b 2-b 1-d 1-d 1-e
5-e
2 4
4-a 5-c
5-e 4-e 2-b 4-c
1 1
4-c
3-e
2 2
1-c 5-c
3-c 5-b
2 2
3-c 1-c
3-b 5-e
3 3
5-e 1-d 4-d
1-d 4-b 3-c
5 3
3-e 1-c 4-b 5-b 4-b
2-b 1-e 1-b
3 5
4-e 5-b 2-a
3-a 1-a 1-b 3-c 5-c
1 3
2-e
5-a 2-b 3-b
4 5
5-d 5-b 4-e 3-e
2-c 1-b 2-d 4-d 4-b
1 1
1-a
2-e
2 1
5-a 4-e
4-a
5 1
1-a 4-a 1-d 4-e 5-c
4-b
2 5
1-a 3-c
4-a 4-d 5-c 2-b 5-b
3 1
2-d 1-a 2-a
1-b
1 5
4-e
1-c 5-a 3-c 3-e 2-a
5 2
4-e 4-b 4-e 2-c 1-e
4-c 1-d
1 5
5-e
5-a 1-b 3-e 4-a 3-d
4 5
1-e 2-c 2-c 1-c
3-b 3-d 2-e 4-c 5-c
2 2
5-a 3-d
2-d 1-b
1 1
4-d
3-a
1 5
4-b
5-a 2-c 5-c 1-e 3-d
3 4
1-d 3-c 3-e
1-e 5-a 1-b 3-c
3 2
5-d 5-c 3-a
4-e 2-b
5 3
2-e 3-e 3-d 5-a 5-c
2-e 4-c 1-c
3 5
1-d 5-a 2-b
5-a 3-d 4-d 4-e 1-a
4 5
5-d 5-b 1-c 2-b
3-a 5-e 3-e 2-c 1-d
4 3
3-a 3-d 1-b 2-b
5-e 2-d 1-a
2 2
1-e 4-b
2-d 2-e
5 5
5-a 3-a 2-b 3-d 3-a
5-d 4-c 5-a 4-a 1-b
5 1
2-a 4-b 4-e 2-a 5-d
3-d
3 3
2-b 2-e 3-b
2-e 3-c 5-d
3 4
5-e 1-d 5-a
1-e 4-a 3-a 4-c
5 5
5-e 3-e 5-a 1-a 4-a
3-b 2-e 5-e 1-a 1-d
5 3
1-c 5-b 3-e 5-c 1-b
2-e 5-e 3-d
5 2
1-e 4-e 5-b 2-a 2-a
2-b 1-d
2 5
2-b 1-e
3-b 5-d 1-d 4-a 5-d
5 3
3-e 2-e 3-b 4-e 3-c
1-b 1-a 1-b
2 5
2-e 4-a
2-d 2-c 3-c 5-d 2-c
5 4
3-a 2-b 2-b 1-a 1-a
2-a 3-e 1-e 2-d
4 1
5-d 4-b 1-a 4-c
2-d
5 5
3-d 1-e 5-d 5-d 2-e
4-c 2-b 2-b 3-c 5-b
5 1
4-c 4-c 2-e 5-b 3-a
4-b
2 5
3-b 2-e
1-b 3-e 3-b 1-e 5-b
3 5
1-b 1-a 4-e
4-e 2-e 4-a 5-b 4-a
2 5
4-e 2-c
1-d 3-e 5-c 3-d 3-d
4 2
1-e 1-e 2-a 3-d
1-d 5-e
4 2
5-c 1-a 4-a 1-a
2-d 3-d
4 2
1-a 4-a 1-e 1-c
5-a 1-d
5 5
4-d 1-a 4-b 4-c 1-c
4-d 3-a 3-a 2-a 3-e
5 2
5-e 4-b 1-a 1-a 5-c
3-c 3-a
1 1
1-a
1-a
3 2
1-e 5-e 3-a
5-c 4-c
3 2
2-a 4-e 5-e
3-e 4-d
4 2
2-b 3-b 1-d 2-d
3-e 5-a
3 1
3-e 5-c 3-d
1-c
5 5
1-e 4-a 1-e 5-c 5-d
1-e 1-d 4-e 3-e 2-e
2 5
3-a 1-c
1-c 2-e 4-c 4-a 5-a
4 1
2-d 2-d 2-c 1-b
2-d
5 3
2-b 5-a 1-c 4-a 4-b
3-d 3-c 4-a
1 1
1-d
2-a
1 1
1-b
4-c
2 4
3-d 5-d
1-c 5-d 5-d 3-a
5 1
4-b 3-a 2-e 2-e 3-e
1-d
4 2
4-e 3-e 1-e 5-c
3-b 2-d
3 1
2-c 5-a 4-c
4-a
5 4
3-b 1-e 1-c 4-a 5-c
5-e 4-b 2-d 3-d
3 2
1-c 2-c 3-d
5-a 1-a
4 2
3-b 2-b 2-d 2-d
4-c 5-b
5 5
5-d 2-b 3-c 1-b 3-d
2-d 2-b 2-a 4-a 2-d
4 4
2-c 5-e 3-c 1-c
2-b 5-d 1-c 5-e
5 1
3-e 2-e 5-c 3-c 1-e
3-d
1 2
3-b
2-d 3-c
4 5
2-c 4-b 5-d 2-c
4-b 2-d 2-c 4-d 4-a
1 5
5-a
5-a 3-c 5-e 1-b 1-c
2 3
5-d 1-a
1-b 1-e 5-c
2 2
5-d 4-d
1-c 5-e
5 5
5-b 2-c 4-a 5-c 3-d
2-b 3-b 5-a 3-c 4-d
4 1
2-a 5-a 4-a 2-d
1-a
2 1
1-e 1-c
1-e
4 3
4-a 3-d 1-a 2-e
1-c 4-d 2-b
1 4
1-e
1-d 1-d 1-b 5-d
3 4
5-a 1-e 1-c
5-c 4-a 2-a 2-b
4 3
1-a 5-d 1-b 4-a
1-c 3-d 1-a
5 4
2-e 4-e 1-a 5-b 3-d
5-c 2-e 4-b 2-a
3 3
3-c 1-e 2-e
1-d 2-e 5-e
3 5
2-d 4-b 5-d
2-a 1-a 5-a 5-c 4-b
2 2
5-d 5-a
5-d 1-e
1 2
2-d
4-d 4-b
5 4
3-a 5-a 2-e 5-e 2-e
1-c 2-d 4-c 3-e
5 4
3-b 3-e 1-e 5-a 2-c
2-d 5-c 4-d 5-b
2 5
4-a 5-a
1-c 1-b 1-e 5-b 5-d
5 1
4-e 1-b 5-b 2-c 1-d
1-a
5 1
1-d 5-b 1-e 5-b 3-c
4-e
5 2
5-c 3-e 5-a 2-c 2-a
5-c 2-a
3 2
2-b 5-b 4-a
4-a 1-b
5 3
2-a 5-e 1-c 1-c 2-a
1-c 4-a 5-a
2 3
5-d 4-b
5-b 2-a 5-b
1 1
1-d
3-e
2 3
2-e 4-e
1-e 1-c 2-e
2 3
1-d 5-c
1-e 1-e 5-b
2 2
1-d 1-d
5-c 4-d
1 2
5-b
5-a 2-c
4 5
1-c 4-c 5-b 5-a
4-d 2-a 3-a 1-b 5-c
5 3
1-c 5-e 5-a 2-a 3-c
2-b 5-b 3-b
5 5
1-c 3-c 4-b 1-e 3-a
1-d 2-a 1-e 5-e 3-e
5 1
3-d 1-a 1-c 5-c 1-b
2-a
5 5
3-c 5-e 4-c 5-a 2-a
2-e 1-a 4-d 1-a 1-e
4 5
1-b 3-d 2-a 5-c
5-b 4-e 3-d 5-c 5-e
4 2
4-a 1-a 2-b 1-d
4-a 1-e
5 2
3-e 5-d 5-e 4-a 3-b
2-e 3-c
5 1
5-c 2-e 4-b 1-a 4-a
5-b
5 2
5-c 5-e 5-e 2-b 2-c
2-a 4-c
2 2
3-e 5-d
1-a 1-d
2 5
4-a 1-b
1-c 5-c 4-e 4-e 4-e
5 4
4-a 1-e 3-d 4-c 4-a
1-c 4-d 4-d 4-d
4 4
1-c 5-c 3-a 1-c
5-a 5-c 2-c 2-c
1 2
1-e
5-a 4-b
1 4
1-c
5-d 5-e 2-b 5-e
3 2
2-a 4-d 4-b
5-d 4-b
2 5
2-b 1-a
1-a 3-e 5-c 2-e 1-d
5 3
2-e 4-d 3-a 2-e 2-c
5-b 1-c 5-a
3 1
5-b 1-a 2-b
4-b
1 3
5-e
5-d 1-a 3-b
2 1
1-d 1-d
1-a
3 3
5-d 4-b 2-d
1-a 1-c 1-c
5 5
5-b 1-d 2-e 3-a 2-d
3-e 4-d 3-c 5-e 4-d
4 4
2-a 3-c 4-d 1-c
2-e 5-c 3-a 3-d
5 3
2-e 5-a 4-a 5-c 5-c
1-e 2-c 1-c
1 5
1-e
4-d 1-e 5-e 4-a 2-b
1 4
2-c
4-c 1-a 2-a 4-d
4 2
4-b 5-e 1-c 1-b
1-e 2-a
5 5
5-a 2-c 1-d 1-b 3-e
2-e 1-a 5-d 5-d 1-c
3 1
3-a 5-a 2-e
4-b
4 4
3-e 5-e 5-d 1-d
5-a 2-e 2-b 4-d
1 2
3-c
1-e 5-e
3 3
3-b 1-e 5-c
5-e 4-a 5-a
5 4
4-c 5-c 5-c 5-d 2-a
3-e 2-d 5-e 4-d
4 1
4-a 1-b 5-e 3-a
3-e
5 2
1-e 3-e 4-c 3-c 2-d
1-a 1-a
5 5
4-d 5-b 5-d 4-d 5-d
4-b 2-c 1-c 5-a 4-d
4 1
3-e 1-c 2-e 4-a
2-d
1 1
1-a
4-c
1 3
2-e
3-a 3-e 2-c
2 2
2-e 2-d
2-c 1-a
3 2
5-d 1-d 4-b
5-c 2-e
5 3
5-e 1-c 4-b 2-c 3-e
2-d 1-a 5-d
2 1
5-e 1-e
4-c
3 2
4-e 1-c 4-c
5-c 4-a
5 4
5-b 3-d 5-c 1-a 1-a
1-d 5-b 3-d 4-b
3 4
4-e 2-a 4-d
1-c 3-c 1-e 2-e
3 2
3-a 3-b 2-e
3-e 1-a
3 5
2-a 4-c 5-b
4-e 3-e 2-b 5-c 2-b
4 3
3-e 4-d 2-e 1-c
1-a 4-e 3-c
5 2
5-e 2-d 5-c 1-d 2-b
1-e 5-a
4 5
4-a 2-e 2-d 1-c
2-a 1-a 2-a 5-d 5-c
3 5
1-a 4-b 2-e
1-c 4-a 2-e 3-b 1-b
4 3
3-a 3-b 1-d 1-c
1-b 1-c 1-a
5 2
2-e 1-d 2-a 4-c 3-d
4-c 2-b
4 2
5-c 5-e 1-d 1-a
1-a 4-b
2 1
2-e 1-c
3-a
4 2
3-a 5-c 4-d 2-e
1-e 2-e
3 1
1-e 1-e 5-c
4-e
5 4
2-c 3-c 5-e 5-d 5-e
4-b 5-d 1-e 5-d
2 3
5-c 5-c
5-d 1-b 5-d
3 2
2-d 5-b 4-e
2-c 3-a
5 4
1-e 2-a 5-e 2-e 5-c
5-e 5-d 3-c 1-c
5 5
4-d 4-d 3-d 2-b 4-a
3-a 3-d 1-d 4-e 5-e
2 1
5-a 1-a
1-a
4 4
5-d 5-d 1-a 2-d
3-e 5-c 4-b 2-d
4 3
1-e 5-e 4-e 1-a
3-b 1-e 1-a
2 5
4-c 3-e
2-c 5-a 4-a 2-c 3-e
5 3
4-c 4-c 5-b 1-e 1-c
2-d 5-d 5-d
2 2
3-e 5-e
2-c 3-c
2 4
5-b 1-b
5-b 4-b 5-e 3-c
4 4
5-e 1-d 5-e 4-d
2-d 1-e 5-b 5-c
5 4
1-e 2-a 5-a 1-d 2-e
1-a 4-e 2-e 1-d
3 4
5-e 3-b 2-b
5-d 5-d 5-e 2-a
4 4
2-e 1-e 3-a 4-b
1-e 3-b 5-b 1-c
4 3
5-c 2-c 1-e 2-e
5-b 1-d 2-c
2 2
2-a 4-e
3-c 3-c
2 5
3-c 1-a
4-b 4-e 4-c 2-c 2-a
2 1
2-e 2-c
4-a
5 5
3-c 3-a 3-a 1-e 5-c
3-a 2-a 3-c 2-d 1-e
5 3
5-b 1-e 2-e 2-a 3-b
5-e 2-e 4-d
5 3
3-d 1-e 3-e 1-e 3-d
5-e 3-b 2-d
5 2
1-a 5-e 5-a 5-c 3-c
5-b 1-d
4 5
3-c 5-e 2-a 2-d
5-e 1-a 4-b 5-b 2-b
5 5
2-b 2-e 2-d 3-c 5-c
5-e 4-b 4-d 4-e 4-c
5 4
1-a 2-b 3-c 2-c 3-a
5-d 5-e 2-e 1-e
4 5
2-e 3-b 4-e 3-b
1-e 4-e 2-e 3-d 2-a
5 4
2-c 5-d 4-e 5-a 1-b
1-c 5-d 4-e 2-b
3 1
2-d 3-a 5-e
5-c
3 3
3-b 1-a 4-e
5-a 3-a 1-e
1 4
3-e
4-b 5-e 3-c 2-e
1 5
1-e
3-b 5-d 1-d 1-d 3-e
2 4
1-d 5-a
1-b 3-a 4-e 2-a
5 5
2-d 5-a 1-d 5-a 1-a
2-d 3-c 1-e 2-e 2-a
1 2
1-a
2-c 2-b
4 3
2-b 3-c 2-c 2-d
5-e 2-d 1-a
4 5
1-c 1-a 5-c 5-c
1-d 4-d 1-d 2-e 4-c
1 5
3-d
3-d 3-b 5-b 5-d 5-d
1 1
2-b
3-a
3 1
2-d 3-d 2-c
2-e
2 3
5-d 3-e
3-e 1-b 4-c
5 4
4-a 4-e 3-d 1-c 2-b
2-e 1-c 5-b 1-e
3 3
5-c 3-a 1-d
4-a 1-d 4-b
5 4
4-a 5-e 4-d 5-e 3-e
1-e 5-e 1-c 1-e
5 3
1-d 3-a 3-d 4-d 4-b
5-e 3-c 2-c
1 5
4-b
3-c 3-a 5-c 1-c 5-d
5 4
3-b 3-d 2-e 2-b 1-a
1-d 4-a 1-c 4-b
5 4
1-a 5-e 2-d 3-c 1-a
5-e 2-c 4-c 1-e
4 4
5-e 1-e 2-d 3-d
5-d 2-d 5-d 4-d
4 5
5-a 4-e 4-b 2-a
3-c 4-b 3-a 3-c 1-b
5 1
3-e 5-e 2-a 2-e 1-d
2-a
3 2
4-a 2-a 5-b
5-c 2-b
1 5
3-d
2-c 3-c 3-e 1-e 2-b
4 1
5-d 4-b 1-c 2-c
5-c
4 1
2-a 5-a 1-e 4-e
3-a
4 3
3-c 3-e 5-a 4-b
3-c 3-c 4-e
5 2
5-e 2-d 4-a 3-e 2-c
5-c 4-b
2 2
1-c 5-e
2-b 4-b
2 3
5-b 1-b
3-e 4-e 2-a
5 4
4-b 5-b 3-a 1-a 5-a
5-b 2-c 4-b 5-b
3 3
3-e 4-c 1-a
5-e 2-e 3-a
2 5
3-b 3-a
5-c 3-a 5-d 5-e 5-b
2 1
5-e 3-e
1-b
4 3
3-c 4-b 3-d 5-d
4-d 1-c 1-a
5 5
5-a 4-e 3-b 5-b 5-c
1-d 5-d 2-e 3-b 5-a
5 1
3-a 3-e 1-b 2-d 1-a
1-a
2 3
1-d 5-a
1-a 3-c 3-b
3 3
4-e 5-b 1-c
4-c 2-e 1-b
3 5
3-d 1-b 1-d
5-d 1-a 4-d 3-c 1-e
5 4
2-d 5-e 5-a 2-e 2-c
5-b 3-e 2-d 5-e
3 3
1-d 1-b 4-d
2-d 5-e 3-e
3 3
3-d 3-d 5-c
5-b 5-a 3-a
3 5
1-a 1-a 4-a
4-a 4-e 2-a 3-c 3-c
3 1
4-e 2-a 2-c
1-a
4 1
1-d 1-e 2-b 4-d
5-c
5 5
5-b 5-a 4-c 4-b 1-c
4-e 3-b 5-b 4-e 2-b
3 5
5-b 1-b 1-b
5-a 2-e 5-c 1-e 3-d
2 1
4-b 1-e
2-e
4 4
1-a 5-a 5-d 2-d
2-b 2-c 5-e 1-c
1 1
2-e
1-b
4 1
1-c 5-b 1-a 5-b
5-d
4 5
3-c 1-d 4-c 5-b
3-a 1-c 3-d 4-d 4-b
2 1
5-e 4-e
3-e
3 4
3-a 4-a 5-b
2-b 4-a 2-d 2-b
5 1
4-e 2-c 5-b 2-e 2-d
3-d
4 3
5-c 4-d 1-e 5-d
3-e 3-d 4-d
5 3
5-c 2-d 3-d 3-b 1-a
5-e 5-c 5-a
5 2
4-b 2-b 3-e 4-c 1-a
2-e 5-a
5 3
3-b 1-d 5-c 2-b 1-a
5-d 1-a 4-a
4 3
2-b 1-e 2-c 4-d
2-a 4-d 4-d
5 4
4-e 5-d 5-e 1-e 4-b
3-b 3-b 5-e 5-e
2 5
1-b 5-c
2-b 5-c 2-a 3-d 2-d
4 2
1-c 5-d 3-d 5-b
5-d 1-d
4 4
3-e 1-a 4-c 5-a
5-b 1-d 5-a 3-a
4 1
2-d 1-e 2-c 4-b
5-b
1 2
3-a
1-d 1-e
2 4
2-a 1-c
1-a 3-a 4-e 2-e
2 5
3-c 4-e
1-d 5-a 2-d 3-d 5-c`

type segment struct {
	len  int64
	char byte
}

func parseSegments(tokens []string) ([]segment, error) {
	segs := make([]segment, 0, len(tokens))
	for _, tok := range tokens {
		parts := strings.Split(tok, "-")
		if len(parts) != 2 || len(parts[1]) == 0 {
			return nil, fmt.Errorf("invalid token %q", tok)
		}
		l, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		c := parts[1][0]
		if len(segs) > 0 && segs[len(segs)-1].char == c {
			segs[len(segs)-1].len += l
		} else {
			segs = append(segs, segment{len: l, char: c})
		}
	}
	return segs, nil
}

func equal(a, b segment) bool {
	return a.len == b.len && a.char == b.char
}

func prefixFunction(p []segment) []int {
	pi := make([]int, len(p))
	for i := 1; i < len(p); i++ {
		j := pi[i-1]
		for j > 0 && !equal(p[i], p[j]) {
			j = pi[j-1]
		}
		if equal(p[i], p[j]) {
			j++
		}
		pi[i] = j
	}
	return pi
}

func solveCase(tTokens, sTokens []string) (string, error) {
	tSegs, err := parseSegments(tTokens)
	if err != nil {
		return "", err
	}
	sSegs, err := parseSegments(sTokens)
	if err != nil {
		return "", err
	}
	n := len(tSegs)
	m := len(sSegs)

	if m == 0 {
		return "0", nil
	}

	if m == 1 {
		var ans int64
		pat := sSegs[0]
		for _, seg := range tSegs {
			if seg.char == pat.char && seg.len >= pat.len {
				ans += seg.len - pat.len + 1
			}
		}
		return fmt.Sprint(ans), nil
	}

	if m == 2 {
		var ans int64
		for i := 0; i < n-1; i++ {
			if tSegs[i].char == sSegs[0].char && tSegs[i+1].char == sSegs[1].char &&
				tSegs[i].len >= sSegs[0].len && tSegs[i+1].len >= sSegs[1].len {
				ans++
			}
		}
		return fmt.Sprint(ans), nil
	}

	midPattern := sSegs[1 : m-1]
	pi := prefixFunction(midPattern)
	ans := 0
	for i, j := 0, 0; i < n; i++ {
		for j > 0 && !equal(tSegs[i], midPattern[j]) {
			j = pi[j-1]
		}
		if equal(tSegs[i], midPattern[j]) {
			j++
			if j == len(midPattern) {
				start := i - len(midPattern) + 1
				if start-1 >= 0 && start+len(midPattern) < n {
					left := tSegs[start-1]
					right := tSegs[start+len(midPattern)]
					if left.char == sSegs[0].char && left.len >= sSegs[0].len &&
						right.char == sSegs[m-1].char && right.len >= sSegs[m-1].len {
						ans++
					}
				}
				j = pi[j-1]
			}
		}
	}
	return fmt.Sprint(ans), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	fields := strings.Fields(testcasesD)
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "empty testcases")
		os.Exit(1)
	}
	ptr := 0
	t, err := strconv.Atoi(fields[ptr])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid test count: %v\n", err)
		os.Exit(1)
	}
	ptr++

	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if ptr+2 > len(fields) {
			fmt.Fprintf(os.Stderr, "unexpected end of tests at case %d\n", caseIdx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[ptr])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad n: %v\n", caseIdx, err)
			os.Exit(1)
		}
		ptr++
		m, err := strconv.Atoi(fields[ptr])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad m: %v\n", caseIdx, err)
			os.Exit(1)
		}
		ptr++
		if ptr+n+m > len(fields) {
			fmt.Fprintf(os.Stderr, "case %d incomplete data\n", caseIdx)
			os.Exit(1)
		}
		tTokens := fields[ptr : ptr+n]
		ptr += n
		sTokens := fields[ptr : ptr+m]
		ptr += m

		expected, err := solveCase(tTokens, sTokens)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", caseIdx, err)
			os.Exit(1)
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, tok := range tTokens {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(tok)
		}
		input.WriteByte('\n')
		for i, tok := range sTokens {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(tok)
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseIdx, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
