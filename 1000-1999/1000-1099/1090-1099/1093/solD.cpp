//#pragma GCC optimize ("O3")
#pragma GCC target ("sse4")
#include<bits/stdc++.h>
#define up(x,y) (x<y?x=y:0)
#define dn(x,y) (x>y?x=y:0)
#define ad(x,y) (x=(x+(y))%M)
#define mu(x,y) (x=x*(y)%M)
#define fi first
#define se second
#define pq priority_queue
#define br putchar('\n')
#define pb push_back
#define isinf(x) (x>=INF?-1:x)

using namespace std;

typedef long long ll;
typedef pair<int,int> pii;

inline ll read()
{
    char c;ll out=0,f=1;
    for (c=getchar();(c<'0'||c>'9')&&c!='-';c=getchar());if (c=='-'){f=-1;c=getchar();}
    for (;c>='0'&&c<='9';c=getchar()){out=(out<<3)+(out<<1)+c-'0';}return out*f;
}
void write(ll x)
{
	if (x<0){putchar('-');write(-x);return;}
    if (x>9){write(x/10);}putchar(x%10+'0');
}

const int N=300010;
const long long M=998244353;

void add(int u,int v);
void dfs(int u);

int head[N],nxt[N<<1],to[N<<1],cnt;
long long n,m,t,c[N],b,w,ans,two[N];
bool flag;

int main()
{
	int i,u,v;
	
	t=read();
	
	two[0]=1;
	for (i=1;i<N;++i)
	{
		two[i]=two[i-1]*2%M;
	}
	
	while (t--)
	{
		n=read();
		m=read();
		memset(head,0,sizeof(int)*(n+1));
		memset(c,0,sizeof(long long)*(n+1));
		cnt=0;
		ans=1;
		flag=false;
		for (i=0;i<m;++i)
		{
			u=read();
			v=read();
			add(u,v);
			add(v,u);
		}
		for (i=1;!flag&&i<=n;++i)
		{
			if (c[i]==0)
			{
				b=w=0;
				c[i]=1;
				dfs(i);
				mu(ans,two[w]+two[b]);
			}
		}
		if (flag)
		{
			ans=0;
		}
		write(ans);
		br;
	}
	
	return 0;
}

void add(int u,int v)
{
	nxt[++cnt]=head[u];
	head[u]=cnt;
	to[cnt]=v;
}

void dfs(int u)
{
	if (c[u]==1)
	{
		++b;
	}
	else
	{
		++w;
	}
	int i,v;
	for (i=head[u];i;i=nxt[i])
	{
		v=to[i];
		if (c[v]==0)
		{
			c[v]=-c[u];
			dfs(v);
			if (flag)
			{
				return;
			}
		}
		else if (c[v]==c[u])
		{
			flag=true;
			return;
		}
	}
}