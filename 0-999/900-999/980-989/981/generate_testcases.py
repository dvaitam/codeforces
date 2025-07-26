import random, string
random.seed(0)

def gen_A(n=100):
    with open('testcasesA.txt','w') as f:
        for _ in range(n):
            l=random.randint(1,50)
            s=''.join(random.choice(string.ascii_lowercase) for _ in range(l))
            f.write(s+'\n')

def gen_B(n=100):
    with open('testcasesB.txt','w') as f:
        for _ in range(n):
            n1=random.randint(1,5)
            f.write(str(n1)+'\n')
            indices=random.sample(range(1,30), n1)
            for a in indices:
                x=random.randint(1,100)
                f.write(f"{a} {x}\n")
            m=random.randint(1,5)
            f.write(str(m)+'\n')
            indices2=random.sample(range(1,30), m)
            for b in indices2:
                y=random.randint(1,100)
                f.write(f"{b} {y}\n")

def gen_C(n=100):
    with open('testcasesC.txt','w') as f:
        for _ in range(n):
            nodes=random.randint(2,8)
            f.write(str(nodes)+'\n')
            for i in range(2,nodes+1):
                p=random.randint(1,i-1)
                f.write(f"{p} {i}\n")

def gen_D(n=100):
    with open('testcasesD.txt','w') as f:
        for _ in range(n):
            n1=random.randint(1,8)
            k=random.randint(1,n1)
            f.write(f"{n1} {k}\n")
            arr=[random.randint(0,1000) for _ in range(n1)]
            f.write(' '.join(map(str,arr))+'\n')

def gen_E(n=100):
    with open('testcasesE.txt','w') as f:
        for _ in range(n):
            n1=random.randint(1,8)
            q=random.randint(1,5)
            f.write(f"{n1} {q}\n")
            for _ in range(q):
                l=random.randint(1,n1)
                r=random.randint(l,n1)
                x=random.randint(1,n1)
                f.write(f"{l} {r} {x}\n")

def gen_F(n=100):
    with open('testcasesF.txt','w') as f:
        for _ in range(n):
            n1=random.randint(1,5)
            L=random.randint(10,100)
            f.write(f"{n1} {L}\n")
            a=[random.randint(0,L) for _ in range(n1)]
            b=[random.randint(0,L) for _ in range(n1)]
            f.write(' '.join(map(str,a))+'\n')
            f.write(' '.join(map(str,b))+'\n')

def gen_G(n=100):
    with open('testcasesG.txt','w') as f:
        for _ in range(n):
            n1=random.randint(1,4)
            q=random.randint(1,5)
            f.write(f"{n1} {q}\n")
            for _ in range(q):
                t=random.randint(1,2)
                if t==1:
                    l=random.randint(1,n1)
                    r=random.randint(l,n1)
                    x=random.randint(1,n1)
                    f.write(f"1 {l} {r} {x}\n")
                else:
                    l=random.randint(1,n1)
                    r=random.randint(l,n1)
                    f.write(f"2 {l} {r}\n")

def gen_H(n=100):
    with open('testcasesH.txt','w') as f:
        for _ in range(n):
            nodes=random.randint(2,5)
            k=random.randint(1,min(3,nodes))
            f.write(f"{nodes} {k}\n")
            for i in range(2,nodes+1):
                p=random.randint(1,i-1)
                f.write(f"{p} {i}\n")

def main():
    gen_A()
    gen_B()
    gen_C()
    gen_D()
    gen_E()
    gen_F()
    gen_G()
    gen_H()

if __name__=='__main__':
    main()
