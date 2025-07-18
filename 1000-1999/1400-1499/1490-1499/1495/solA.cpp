#include<bits/stdc++.h> 
#define inf 2139062143
#define MAXN 100100
#define MOD 1000000007
#define ll long long
#define ull unsigned long long
#define rep(i,s,t) for(register int i=(s),i##end=(t);i<=i##end;++i)
#define dwn(i,s,t) for(register int i=(s),i##end=(t);i>=i##end;--i)
#define ren for(register int i=fst[x];i;i=nxt[i])
#define Fill(a,b) memset(a,b,sizeof(a))
#define pls(a,b) ((a+b)%MOD+MOD)%MOD
#define mns(a,b) (((a)-(b))%MOD+MOD)%MOD
#define mul(a,b) (1LL*(a)*(b)%MOD)
#define pii pair<int,int>
#define mp(a,b) make_pair(a,b)
#define fi first
#define se second
using namespace std;
inline int read()
{
	int x=0;bool f=1;char ch=getchar();
	while(!isdigit(ch)) {if(ch=='-') f=0;ch=getchar();}
	while(isdigit(ch)) x=x*10+(ch^48),ch=getchar();
	return f?x:-x;
}
struct Dinic
{
    int fst[MAXN],nxt[MAXN<<1],to[MAXN<<1],cnt,val[MAXN<<1];
    int vis[MAXN],q[MAXN],l,r,S,T,tot,dis[MAXN],cur[MAXN];
    Dinic(){memset(fst,0,sizeof(fst));cnt=1,tot=0;}
    void add(int u,int v,int w) {nxt[++cnt]=fst[u],fst[u]=cnt,to[cnt]=v,val[cnt]=w;}
    void ins(int u,int v,int w) {add(u,v,w);add(v,u,0);}
    int bfs()
    {
        vis[T]=++tot,dis[T]=0,q[l=r=1]=T;int x;
        while(l<=r)
        {
            x=q[l++],cur[x]=fst[x];ren if(val[i^1]&&vis[to[i]]!=tot)
                dis[to[i]]=dis[x]+1,vis[to[i]]=tot,q[++r]=to[i];
        }
        return vis[S]==tot;
    }
    int dfs(int x,int a)
    {
        if(x==T||!a) return a;int flw=0,f;
        for(int& i=cur[x];i&&a;i=nxt[i])
            if(val[i]&&dis[to[i]]==dis[x]-1&&(f=dfs(to[i],min(a,val[i]))))
                val[i]-=f,val[i^1]+=f,a-=f,flw+=f;
        return flw;
    }
    int solve(int ss,int tt,int res=0)
        {S=ss,T=tt;while(bfs()) res+=dfs(S,inf);return res;}
}D;
int n,m,nn,a[MAXN],b[MAXN];
double ans;
int main()
{
	int x,y;rep(T,1,read())
	{
		n=read();nn=m=0;rep(i,1,n<<1) {x=read(),y=read();if(!x) a[++nn]=abs(y);else b[++m]=abs(x);}
		sort(a+1,a+n+1);sort(b+1,b+m+1);ans=0;
		rep(i,1,n) ans+=sqrt(1.0*a[i]*a[i]+1.0*b[i]*b[i]);
		printf("%.9lf\n",ans);
	}
}