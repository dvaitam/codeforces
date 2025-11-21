from math import gcd
N=200
max_steps=2*N-2
reachable=[set() for _ in range(max_steps+1)]
reachable[0].add(1)
for t in range(max_steps):
    total=t+2
    for x in reachable[t]:
        y=total-x
        # L wins
        nx=x+1; ny=total+1-nx
        if nx<=N and ny<=N and gcd(nx,ny)==1:
            reachable[t+1].add(nx)
        # F wins
        nx=x; ny=total+1-nx
        if nx<=N and ny<=N and gcd(nx,ny)==1:
            reachable[t+1].add(nx)
print('max states per t',max(len(s) for s in reachable))
print('total states',sum(len(s) for s in reachable))
