#include<bits/stdc++.h>
#define rep(i,l,r) for(int i=l;i<=r;i++)
#define per(i,r,l) for(int i=r;i>=l;i--)
#define for4(i,x) for(int i=head[x],y=e[i].go;i;i=e[i].next,y=e[i].go)
#define maxn (400000+5)
#define mod (1000000007)
#define ll long long
#define inf 1000000009
using namespace std;
inline int read()
{
    int x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=10*x+ch-'0';ch=getchar();}
    return x*f;
}
inline ll readll()
{
	ll x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=10*x+ch-'0';ch=getchar();}
    return x*f;
}
long long gcd(long long x,long long y){return y?gcd(y,x%y):x;}
long long power(long x,long y)
{
	long long t=1;
	for(;y;y>>=1,x=x*x%mod)
		if(y&1)t=t*x%mod;
	return t;
}
/*
struct edge
{
	int go,next;
}e[2*maxn];
void insert(int x,int y)
{
	e[++tot]=(edge){y,head[x]};head[x]=tot;
	e[++tot]=(edge){x,head[y]};head[y]=tot;
}
*/
ll n,L,S,w[maxn],dep[maxn],sum[maxn];
bool v[maxn];
queue<int>q;
int fa[maxn],son[maxn],id[maxn];
int find(int x){return id[x]==x?x:id[x]=find(id[x]);}
int main()
{
	n=read();L=readll();S=readll();
	rep(i,1,n)w[i]=read(),id[i]=i;
	rep(i,1,n)if(w[i]>S)
	{
		cout<<-1<<endl;
		return 0;
	}
	rep(i,2,n)fa[i]=read(),son[fa[i]]++;
	dep[1]=1;sum[1]=w[1];
	rep(i,2,n)dep[i]=dep[fa[i]]+1,sum[i]=sum[fa[i]]+w[i];
	rep(i,1,n)if(!son[i])q.push(i);
	int ans=0;
	while(!q.empty())
	{
		int x=q.front();q.pop();
		if(v[x])continue;
		ans++;
		ll tmp=0,dist=0;
		while(x&&tmp<=S)
		{
			int xx=find(x);
			dist+=dep[x]-dep[xx];
			tmp+=sum[x]-sum[xx];
			if(dist+1>L||tmp+w[x]>S)break;
			id[xx]=find(fa[xx]);
			dist++;
			tmp+=w[xx];
			v[xx]=1;
			son[fa[xx]]--;
			if(fa[xx]&&son[fa[xx]]==0&&!v[fa[xx]])q.push(fa[xx]);
			x=fa[xx];
		}
	}
	cout<<ans<<endl;
	return 0;
}