import random, string

random.seed(0)

def rand_digits(length):
    return ''.join(random.choice('0123456789') for _ in range(length))

# Problem A
def gen_A(n_cases=100):
    with open('testcasesA.txt', 'w') as f:
        for _ in range(n_cases):
            n = random.randint(2, 5)
            length = random.randint(1, 10)
            first = rand_digits(length)
            numbers = [first]
            for _ in range(n-1):
                if random.random() < 0.5:
                    # share some prefix
                    p = random.randint(0, length)
                    s = first[:p] + rand_digits(length-p)
                else:
                    s = rand_digits(length)
                numbers.append(s)
            f.write(str(n) + ' ' + ' '.join(numbers) + '\n')

# Problem B
def gen_B(n_cases=100):
    with open('testcasesB.txt', 'w') as f:
        for _ in range(n_cases):
            m = random.randint(1, 1000)
            a = random.randint(0, 1000)
            b = random.randint(0, 1000)
            r0 = random.randint(0, m-1)
            f.write(f"{a} {b} {m} {r0}\n")

# Problem C
def gen_C(n_cases=100):
    with open('testcasesC.txt', 'w') as f:
        for _ in range(n_cases):
            n = random.randint(1, 4)
            m = random.randint(1, 3)
            t = 0
            vals = []
            for _ in range(n):
                t += random.randint(1, 3)
                x = random.randint(1, 5)
                vals.append((t, x))
            line = [str(n), str(m)]
            for ti, xi in vals:
                line.append(str(ti))
                line.append(str(xi))
            f.write(' '.join(line) + '\n')

# Problem D
def gen_D(n_cases=100):
    with open('testcasesD.txt', 'w') as f:
        for _ in range(n_cases):
            a = random.randint(1, 50)
            n = random.randint(1, 50)
            f.write(f"{a} {n}\n")

# Generate random HTML tree for problem E
def gen_tree(depth=0):
    tag = random.choice(['a','b','c','d'])
    if depth > 2 or random.random() < 0.3:
        if random.random() < 0.5:
            return f"<{tag}/>", [ [tag] ]
        else:
            return f"<{tag}></{tag}>", [ [tag] ]
    num_children = random.randint(0,2)
    child_strings = []
    paths = []
    for _ in range(num_children):
        cs, ps = gen_tree(depth+1)
        child_strings.append(cs)
        for p in ps:
            paths.append([tag]+p)
    inner = ''.join(child_strings)
    return f"<{tag}>{inner}</{tag}>", ([[tag]]+paths)


def gen_E(n_cases=100):
    with open('testcasesE.txt', 'w') as f:
        for _ in range(n_cases):
            doc, paths = gen_tree()
            # flatten paths
            all_paths = paths
            m = random.randint(1, 3)
            queries = []
            for i in range(m):
                if random.random()<0.5 and all_paths:
                    q = random.choice(all_paths)
                else:
                    l = random.randint(1,2)
                    q = [random.choice(['a','b','c','d']) for _ in range(l)]
                queries.append(' '.join(q))
            line = doc + '|' + str(m)
            for q in queries:
                line += '|' + q
            f.write(line + '\n')

if __name__ == '__main__':
    gen_A()
    gen_B()
    gen_C()
    gen_D()
    gen_E()
