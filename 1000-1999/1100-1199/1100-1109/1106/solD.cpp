#include<iostream>
#include<cstdio>
#include<cmath>
#include<string>
#include<cstring>
#include<cstdlib>
#include<algorithm>
#include<vector>
#include<queue>
#include<stack>
#include<set>
#include<map>
using namespace std;
#define isNum(a) (a>='0'&&a<='9')
#define SP putchar (' ')
#define EL putchar ('\n')
#define inf 2147483647
#define N 100005
#define M 200005
#define File(a) freopen(a".in", "r", stdin), freopen(a".out", "w", stdout)
template<class T>inline void read(T&);
template<class T>inline void write(const T&);
void add(int, int);
int hed[N], to[M], nxt[M], id;
priority_queue<int, vector<int>, greater<int> >k;
bool vis[N];
int main () {
    int n, m;
    read(n), read(m);
    vis[1]=true;
    for (int i=1; i<=m; ++i) {
        int x, y;
        read(x), read(y);
        add(x, y);
        add(y, x);
        if (x==1) {
            if (!vis[y]) {
                k.push(y);
                vis[y]=true;
            }
        }
        if (y==1) {
            if (!vis[x]) {
                k.push(x);
                vis[x]=true;
            }
        }
    }
    write(1);SP;
    while (!k.empty()) {
        int u=k.top();
        k.pop();
        write(u), SP;
        for (int i=hed[u]; i; i=nxt[i]) {
            int v=to[i];
            if (!vis[v]) {
                k.push(v);
                vis[v]=true;
            }
        }
    }
    EL;
    return 0;
}
template<class T>void read(T &Re) {
    T k=0;
    char ch=getchar();
    int flag=1;
    while (!isNum(ch)) {
        if (ch=='-') {
            flag=-1;
        }
        ch=getchar();
    }
    while (isNum(ch)) {
        k=(k<<1)+(k<<3)+ch-'0';
        ch=getchar();
    }
    Re=flag*k;
}
template<class T>void write(const T& Wr) {
    if (Wr<0) {
        putchar('-');
        write(-Wr);
    }else {
        if (Wr<10) {
            putchar(Wr+'0');
        }else {
            write(Wr/10);
            putchar((Wr%10)+'0');
        }
    }
}
void add(int x, int y) {
    nxt[++id]=hed[x];
    hed[x]=id;
    to[id]=y;
}