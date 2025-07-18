#include <bits/stdc++.h>
using namespace std;
using namespace __gnu_cxx;


#define IOC std::ios::sync_with_stdio(false)
typedef long long LL;
typedef unsigned long long uLL;
typedef unsigned int uit;
typedef long double ld;
typedef pair<int,int> pii;
const int INF = 0x3f3f3f3f;
const long double enf = 1e-20;
const double pi = acos(-1.0);
inline LL read()
{
    LL X=0,w=0; char ch=0;
    while(!isdigit(ch)) {w|=ch=='-';ch=getchar();}
    while(isdigit(ch)) X=(X<<3)+(X<<1)+(ch^48),ch=getchar();
    return w?-X:X;
}
inline void write(LL x)
{
     if(x<0) putchar('-'),x=-x;
     if(x>9) write(x/10);
     putchar(x%10+'0');
}
//-------head------------


const int maxn = 2e6 + 10;

struct edge{
    int u,v;
}dict[maxn];
int in[maxn];

struct UFS{
    int fa[maxn];
    UFS(){
        for(int i = 0; i < maxn; ++i){
            fa[i] = i;
        }
    }
    int _find(int x){
        return x == fa[x] ? x : fa[x] = _find(fa[x]);
    }
    bool same(int x,int y){
        return _find(x) == _find(y);
    }
    void un(int x,int y){
        x = _find(x);
        y = _find(y);
        fa[x] = y;
    }
}vis;

bool check[maxn];

int main(){
    int n = read(),m = read();
    for(int i = 0; i < m; ++i){
        dict[i].u = read(),dict[i].v = read();
        ++in[dict[i].u],++in[dict[i].v];
    }
    int who = 1;
    for(int i = 1; i <= n; ++i){
        if(in[i] > in[who]){
            who = i;
        }
    }
    for(int i = 0; i < m; ++i){
        if(dict[i].u == who || dict[i].v == who){
            vis.un(dict[i].u,dict[i].v);
            check[i] = 1;
        }
    }
    for(int i = 0; i < m; ++i){
        if(check[i])    continue;
        if(!vis.same(dict[i].u,dict[i].v)){
            vis.un(dict[i].u,dict[i].v);
            check[i] = 1;
        }
    }
    for(int i = 0; i < m; ++i){
        if(check[i]){
            printf("%d %d\n",dict[i].u,dict[i].v);
        }
    }



    return 0;
}