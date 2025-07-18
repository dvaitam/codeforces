/*{{{*/



//#pragma GCC optimize("O3")



#include <bits/stdc++.h>



#ifdef __linux__

#define getchar getchar_unlocked

#define putchar putchar_unlocked

#endif



typedef std::pair<int, int> pii;



std::string program_name = __FILE__;

std::string iput = program_name.substr(0, program_name.length() - 4) + ".in";

std::string oput = program_name.substr(0, program_name.length() - 4) + ".out";



template <class T> inline bool chkmin(T &x, T y) { return x > y ? x = y, true : false; }

template <class T> inline bool chkmax(T &x, T y) { return x < y ? x = y, true : false; }



template <class T> inline T &read(T &x)

{

    static int f;

    static char c; 

    for (f = 1; !isdigit(c = getchar()); ) {

        if (c == '-')

            f = -1;

    }

    for (x = 0; isdigit(c); c = getchar()) {

        x = x * 10 + c - 48;

    }

    return x *= f;

}

 

template <class T> void write(T x, const char p = '\n')

{

    static int top;

    static int s[30];

    if (x < 0) {

        x = -x;

        putchar('-');

    }

    do s[++ top] = x % 10 + 48;

    while (x /= 10);

    while (top) 

        putchar(s[top --]);

    putchar(p);

}

/*}}}*/

 

const int mo = 1000000007;

const int oo = 0x3f3f3f3f;

const int maxn = 2e5 + 5;

const int maxm = 2e5 + 5;

const int Log = 19;



struct Edge { int u, v, w, c, id; } E[maxm];

inline bool cmpE(const Edge &A, const Edge &B) { return A.w < B.w; }



int fa[maxn];

inline int find(int x) { return x == fa[x] ? x : fa[x] = find(fa[x]); }



int st[maxn], cnte;



struct edge { int to, next, w; } e[maxm << 1];



inline void addedge(int u, int v, int w)

{

    e[++ cnte] = (edge) {v, st[u], w};

    st[u] = cnte;

}

 

int n;

int m;

int s;

int a[maxn];

bool vis[maxm];



int dep[maxn];

int Fa[maxn][Log];

int max[maxn][Log];



void preDFS(int x, int fa)

{

    dep[x] = dep[fa] + 1;

    for (int i = st[x]; i; i = e[i].next) {

        int v = e[i].to;

        if (v == fa) continue;

        Fa[v][0] = x;

        max[v][0] = e[i].w;

        for (int j = 1; j < Log; ++ j) {

            Fa[v][j] = Fa[Fa[v][j - 1]][j - 1];

            max[v][j] = std::max(max[v][j - 1], max[Fa[v][j - 1]][j - 1]);

        }

        preDFS(v, x);

    }

}



int getmaxw(int u, int v)

{

    if (dep[u] < dep[v]) std::swap(u, v);



    int ret = 0;



    for (int i = 0, k = dep[u] - dep[v]; i < Log; ++ i)

        if (k >> i & 1)

            chkmax(ret, max[u][i]), u = Fa[u][i];



    if (u == v) return ret;



    for (int i = Log - 1; ~i; -- i) {

        if (Fa[u][i] != Fa[v][i]) {

            chkmax(ret, max[u][i]);

            chkmax(ret, max[v][i]);

            u = Fa[u][i];

            v = Fa[v][i];

        }

    }



    chkmax(ret, max[u][0]);

    chkmax(ret, max[v][0]);



    return ret;

}



void exec() 

{

    static int to[maxn];



    read(n);

    read(m);

    for (int i = 1; i <= m; ++ i) read(E[i].w);

    for (int i = 1; i <= m; ++ i) read(E[i].c);

    for (int i = 1; i <= m; ++ i) read(E[i].u), read(E[i].v), E[i].id = i;

    read(s);



    int minc = oo, pc;

    long long totw = 0;



    std::sort(E + 1, E + m + 1, cmpE);

    for (int i = 1; i <= m; ++ i) to[E[i].id] = i;

    for (int i = 1; i <= n; ++ i) fa[i] = i;

    for (int i = 1; i <= m; ++ i) {

        int A = find(E[i].u);

        int B = find(E[i].v);

        if (A != B) {

            fa[A] = B;

            totw += E[i].w;

            if (chkmin(minc, E[i].c)) pc = i;

            addedge(E[i].u, E[i].v, E[i].w);

            addedge(E[i].v, E[i].u, E[i].w);

            vis[i] = true;

        }

    }

    preDFS(1, 0);



    int pos = 0;

    long long ans = totw - s / minc;

    for (int i = 1; i <= m; ++ i) {

        if (!vis[i] && E[i].c < minc) {

            if (chkmin(ans, totw - getmaxw(E[i].u, E[i].v) + E[i].w - s / E[i].c))

                pos = i;

        }

    }

    write(ans);

    if (pos == 0) {

        for (int i = 1; i <= m; ++ i) {

            if (vis[to[i]]) {

                if (to[i] == pc)

                    write(i, ' '), write(E[to[i]].w - s / minc);

                else

                    write(i, ' '), write(E[to[i]].w);

            }

        }

    }

    else {

        int u = E[pos].u, v = E[pos].v;

        int U, V, W = 0;

        while (u != v) {

            if (dep[u] < dep[v]) std::swap(u, v);

            if (chkmax(W, max[u][0]))

                U = u, V = Fa[u][0];

            u = Fa[u][0];

        }

        if (U < V) std::swap(U, V);

        for (int i = 1; i <= m; ++ i) {

            if (vis[to[i]] && (std::max(E[to[i]].u, E[to[i]].v) != U || std::min(E[to[i]].u, E[to[i]].v) != V))

                write(i, ' '), write(E[to[i]].w);

            else

                if (to[i] == pos)

                    write(i, ' '), write(E[to[i]].w - s / E[to[i]].c);

        }

    }

}



/*{{{*/

int main() 

{

    if (fopen(iput.c_str(), "r") != NULL) {

        freopen(iput.c_str(), "r", stdin);

        freopen(oput.c_str(), "w", stdout);

    }



    do exec();

    while (0);



    fclose(stdin);

    fclose(stdout);



    return 0;

}

/*}}}*/