#include<iostream>

#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cstdlib>

#include<vector>

#include<map>

#include<set>

#include<queue>

#include<stack>

#include<cmath>

#define Clr(a,x) memset(a,x,sizeof(a));

#define For(i,x,y) for (int i=x;i<=y;i++)

#define For_Edge(k,u) for (int k=head[u];k;k=e[k].next)

#define Dor(i,y,x) for (int i=y;i>=x;i--)

#define eps 1e-8

#define inf 2147483647

#define MAXN 300005

#define pa pair<long long,int>

using namespace std;



/*---------------------------------------------------------------*/

int read()

{

    int x=0,f=1;char ch=getchar();

    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}

    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}

    return x*f;

}

/*---------------------------------------------------------------*/

int n,m;

struct edge{int t,next,id,v;}e[MAXN*2];

int head[MAXN],ne=0,vis[MAXN];

long long dis[MAXN];

int st[MAXN],tp=0;

void Insert(int u,int v,int w,int id)

{e[++ne].t=v;e[ne].next=head[u];head[u]=ne;e[ne].v=w;e[ne].id=id;}

priority_queue<pa,vector<pa > ,greater<pa> > pq;

void Dij(int s)

{

    For(i,1,n) dis[i]=(1LL<<60);Clr(vis,0);

    dis[s]=0;pq.push(make_pair(0,s));

    while (!pq.empty())

    {

        int u=pq.top().second;pq.pop();

        if (vis[u]) continue;vis[u]=1;

        for (int k=head[u];k;k=e[k].next)

        {

            int v=e[k].t;

            if (dis[v]>dis[u]+(long long)e[k].v)

            {

                dis[v]=dis[u]+(long long)e[k].v;

                pq.push(make_pair(dis[v],v));

            }

        }

    }

}

int main()

{

    n=read();m=read();

    For(i,1,m)

    {

        int u=read(),v=read(),w=read();

        Insert(u,v,w,i);Insert(v,u,w,i);

    }

    int fir=read();Dij(fir);

    long long ans=0;

    For(u,1,n)

    { 

        if (u==fir) continue;

        int mn=2147483647,add=0;

        for (int k=head[u];k;k=e[k].next)

        {

            int v=e[k].t;

            if (dis[u]==dis[v]+e[k].v&&mn>e[k].v)

                mn=e[k].v,add=e[k].id;

        }

        ans+=(long long)mn;st[++tp]=add;

    }

    cout<<ans<<endl;

    For(i,1,tp) printf("%d ",st[i]);puts("");

    return 0;

}