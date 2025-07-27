import random
random.seed(0)

# Problem A
with open('testcasesA.txt','w') as f:
    for _ in range(100):
        n = random.randint(1, 20)
        s = ''.join(random.choice('()') for _ in range(n))
        f.write(f"{n} {s}\n")

# Problem B
with open('testcasesB.txt','w') as f:
    for _ in range(100):
        n = random.randint(1,8)
        arr = [random.randint(0,20) for _ in range(n)]
        f.write(str(n)+' '+' '.join(map(str,arr))+'\n')

# Problem C
with open('testcasesC.txt','w') as f:
    for _ in range(100):
        n = random.randint(1,5)
        m = random.randint(0, n*n)
        weights = [random.randint(1,20) for _ in range(n)]
        edges = set()
        while len(edges) < m:
            u = random.randint(1,n)
            v = random.randint(1,n)
            edges.add((u,v))
        edges = list(edges)
        line = [str(n), str(len(edges))] + [str(x) for x in weights]
        for u,v in edges:
            line.append(str(u))
            line.append(str(v))
        f.write(' '.join(line)+'\n')

# Problem D
with open('testcasesD.txt','w') as f:
    for _ in range(100):
        n = random.randint(1,5)
        m = random.randint(1,5)
        l = [random.randint(1,m) for _ in range(n)]
        s = [random.randint(0,5) for _ in range(n)]
        c = [random.randint(-5,5) for _ in range(n+m)]
        line = [str(n), str(m)] + [str(x) for x in l] + [str(x) for x in s] + [str(x) for x in c]
        f.write(' '.join(line)+'\n')

# Problem E
with open('testcasesE.txt','w') as f:
    for _ in range(100):
        n = random.randint(1,6)
        arr = [random.randint(0,10) for _ in range(n)]
        f.write(str(n)+' '+' '.join(map(str,arr))+'\n')

# Problem F
with open('testcasesF.txt','w') as f:
    for _ in range(100):
        n = random.randint(2,5)
        m = random.randint(1,3)
        edges = []
        for i in range(2,n+1):
            parent = random.randint(1,i-1)
            edges.append((parent,i))
        lines = []
        for _ in range(m):
            a = random.randint(1,n)
            b = random.randint(1,n)
            while b==a:
                b = random.randint(1,n)
            lines.append((a,b))
        line = [str(n), str(m)]
        for u,v in edges:
            line.extend([str(u), str(v)])
        for a,b in lines:
            line.extend([str(a), str(b)])
        f.write(' '.join(line)+'\n')
