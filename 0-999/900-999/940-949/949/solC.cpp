#include <cstdio>
#include <cstring>
#define N 100005
using namespace std;
namespace fastIN{
    #define frdBUF_SIZE 21000000
    char frdbuf[frdBUF_SIZE],*frdp1,*frdpend;
    #define rd(x) {\
        while (*++frdp1>>4!=3); x=0;\
        do x=(x<<1)+(x<<3)+((*frdp1)&15);\
        while (*++frdp1>>4==3);\
    }
 inline char gc(){return *+frdp1;}
    inline void initfrd(){frdp1=frdbuf-1,frdpend=frdbuf+fread(frdbuf,1,frdBUF_SIZE,stdin);}
};
using namespace fastIN;
namespace fastOUT{
    #define fotOUT_SIZE 4000000
    char fotout[fotOUT_SIZE],*fotop1=fotout-1;
    const char fotDigit[200] = {
        '0','0','0','1','0','2','0','3','0','4','0','5','0','6','0','7','0','8','0','9','1','0','1','1','1','2','1','3','1','4','1','5','1','6','1','7','1','8','1','9',
        '2','0','2','1','2','2','2','3','2','4','2','5','2','6','2','7','2','8','2','9','3','0','3','1','3','2','3','3','3','4','3','5','3','6','3','7','3','8','3','9',
        '4','0','4','1','4','2','4','3','4','4','4','5','4','6','4','7','4','8','4','9','5','0','5','1','5','2','5','3','5','4','5','5','5','6','5','7','5','8','5','9',
        '6','0','6','1','6','2','6','3','6','4','6','5','6','6','6','7','6','8','6','9','7','0','7','1','7','2','7','3','7','4','7','5','7','6','7','7','7','8','7','9',
        '8','0','8','1','8','2','8','3','8','4','8','5','8','6','8','7','8','8','8','9','9','0','9','1','9','2','9','3','9','4','9','5','9','6','9','7','9','8','9','9'
    };
    inline void out11(const int x){*++fotop1=fotDigit[x<<1],*++fotop1=fotDigit[x<<1|1];}
    inline void out12(const int x){if (x>=10)*++fotop1=fotDigit[x<<1];*++fotop1=fotDigit[x<<1|1];}
    inline void out21(const int x){out11(x/100);out11(x%100);}
    inline void out22(const int x){out12(x/100);out11(x%100);}
    inline void out31(const int x){out21(x/10000);out21(x%10000);}
    inline void out32(const int x){out22(x/10000);out21(x%10000);}
    inline void out42(const long long x){out32(x/100000000);out31(x%100000000);}
    inline void prt(int x){
        if (x<0)*++fotop1='-',x=~x+1;
        if (x<100)out12(x);else
        if (x<10000)out22(x);else
        if (x<1000000)out22(x/100),out11(x%100);else
        if (x<100000000)out32(x);else
        out32(x/100),out11(x%100);
    }
    inline void prt(long long x){
        if (x<0)*++fotop1='-',x=~x+1;
        if (x<100)out12(x);else
        if (x<10000)out22(x);else
        if (x<1000000)out22(x/100),out11(x%100);else
        if (x<100000000)out32(x);else
        if (x<10000000000)out32(x/100),out11(x%100);else
        if (x<1000000000000)out32(x/10000),out21(x%10000);else
        if (x<100000000000000)out32(x/1000000),out21(x/100%10000),out11(x%100);else
        if (x<10000000000000000)out42(x);else
        if (x<1000000000000000000)out42(x/100),out11(x%100);
        else out42(x/10000),out21(x%10000);
    }
    inline void prt(const char x){*++fotop1=x;}
    inline void flushfot(){fwrite(fotout,1,fotop1-fotout+1,stdout);}
}
using namespace fastOUT;
int head[N], node[N * 2], nex[N * 2], tote,low[N], dfn[N], sz[N], scc[N], S[N], dout[N], ccnt, times;
void dfs(int u)
{
    low[u] = dfn[u] = ++times;
    S[++S[0]] = u;
    for (int tmp = head[u]; ~tmp; tmp = nex[tmp])
    {
        if (!dfn[node[tmp]])
        {
            dfs(node[tmp]);
            if (low[node[tmp]] < low[u])
                low[u] = low[node[tmp]];
        }
        else if (!scc[node[tmp]] && dfn[node[tmp]] < low[u])
            low[u] = dfn[node[tmp]];
    }
    if (dfn[u] == low[u])
    {
        ccnt++;
        for(;;)
        {
            int v = S[S[0]--];
            scc[v] = ccnt;
            sz[ccnt]++;
            if (v == u)
                break;
        }
    }
}
void addedge(int u, int v)
{
    node[tote] = v;
    nex[tote] = head[u];
    head[u] = tote++;
}
int w[N];
int n, m, h;
int main()
{
    initfrd();
    rd(n); rd(m); rd(h);
    memset(head, -1, sizeof(head));
    for (int i = 1; i <= n; i++)
        rd(w[i]);
    int u, v;
    for (int i = 1; i <= m; i++)
    {
        rd(u); rd(v);
        if ((w[u] + 1) % h == w[v])
            addedge(u, v);
        if ((w[v] + 1) % h == w[u])
            addedge(v, u);
    }
    S[0] = 0;
    ccnt = 0;
    times = 0;
    for (int i = 1; i <= n; i++)
        if (!dfn[i])
            dfs(i);
    for (int i = 1; i <= n; i++)
        for (int j = head[i]; j != -1; j = nex[j])
        {
            int v = node[j];
            if (scc[i] != scc[v])
                dout[scc[i]]++;
        }
    int Ans = 1e9, t = 0;
    for (int i = 1; i <= ccnt; i++)
        if (!dout[i] && sz[i] < Ans)
            Ans = sz[i], t = i;
    prt( Ans); prt('\n');
    for (int i = 1; i <= n; i++)
        if (scc[i] == t)
            prt( i),prt('\n');
    flushfot();
}