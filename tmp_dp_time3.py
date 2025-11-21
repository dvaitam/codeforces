from math import gcd
N=2000
max_steps=2*N-2
reachable=[set() for _ in range(max_steps+1)]
reachable[0].add(1)
maxsize=0
for t in range(max_steps):
    total=t+2
    for x in list(reachable[t]):
        y=total-x
        nx=x+1; ny=total+1-nx
        if nx<=N and ny<=N and gcd(nx,ny)==1:
            reachable[t+1].add(nx)
        nx=x; ny=total+1-nx
        if nx<=N and ny<=N and gcd(nx,ny)==1:
            reachable[t+1].add(nx)
    if len(reachable[t+1])>maxsize:
        maxsize=len(reachable[t+1])
print('max states per t',maxsize)
print('total states',sum(len(s) for s in reachable))
