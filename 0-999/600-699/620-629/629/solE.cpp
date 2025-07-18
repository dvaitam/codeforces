#include <cstdio>

#define maxm 400010

struct edge{int to,len,next;}E[maxm];

int cnt,last[maxm],fa[maxm],top[maxm],deep[maxm],son[maxm],val[maxm];

long long siz[maxm],sum[maxm];

void addedge(int a,int b,int len=0){

	E[++cnt]=(edge){b,len,last[a]},last[a]=cnt;

}

void dfs1(int x){

	deep[x]=deep[fa[x]]+1;siz[x]=1;

	for(int i=last[x];i;i=E[i].next){

		int to=E[i].to;

		if(fa[x]!=to&&!fa[to]){

			val[to]=E[i].len;

			fa[to]=x;

			dfs1(to);

			siz[x]+=siz[to];

			sum[x]+=sum[to]+siz[to];

			if(siz[son[x]]<siz[to])son[x]=to;

		}

	}

}

void dfs2(int x){

	if(x==son[fa[x]])top[x]=top[fa[x]];

	else top[x]=x;

	for(int i=last[x];i;i=E[i].next)if(fa[E[i].to]==x)dfs2(E[i].to);

}

void init(int root){dfs1(root),dfs2(root);}

int query(int x,int y){

	for(;top[x]!=top[y];deep[top[x]]>deep[top[y]]?x=fa[top[x]]:y=fa[top[y]]);

	return deep[x]<deep[y]?x:y;

}

long long src[maxm];

int n,m;

void dfs3(int x)

{

    src[x]+=sum[x];

    for(int i=last[x];i;i=E[i].next)

    {

        if(fa[E[i].to]==x)

        {

            src[E[i].to]+=src[x]-sum[E[i].to]-siz[E[i].to]+n-siz[E[i].to];

            dfs3(E[i].to);

        }

    }

}

int solve(int u,int f)

{

	for(;;)

	{

		if(deep[fa[top[u]]]>deep[f])

			u=fa[top[u]];

		else

		if(fa[u]==f)

			return u;

		else 

		if(deep[top[u]]>deep[f])

			u=top[u];

		else

			return son[f]; 

	}

}

double ask(int u,int v)

{

	int lca=query(u,v);

	if(u!=lca&&v!=lca)

	{

		return(double)sum[u]/siz[u]+(double)sum[v]/siz[v]+deep[u]+deep[v]+1-2*deep[lca];

	}

	else

	{

		if(u==lca)cnt=u,u=v,v=cnt;int nw=solve(u,v);

		register double x,y;

		x=sum[u];x/=siz[u];

		y=src[v]-sum[nw]-siz[nw];y/=(n-siz[nw]);

		return x+y+deep[u]-deep[v]+1;

	}

}

int main(){

	scanf("%d%d",&n,&m);

	int x,y;

	for(int i=1;i<n;i++)

	{

		scanf("%d%d",&x,&y);addedge(x,y);addedge(y,x);

	}

	dfs1(1),dfs2(1),dfs3(1);

	for(int i=1;i<=m;i++)

	{

		double ans;

		scanf("%d%d",&x,&y);

		ans=ask(x,y);if(ans<0)return puts("nan"),0;

		printf("%.7lf\n",ans);

	}

	return 0 ;

}