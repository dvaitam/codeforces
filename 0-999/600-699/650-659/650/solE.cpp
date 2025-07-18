#include<cstdio>

using namespace std;

const int N=500500;

int i,j,k,n,En1,En2,ch,x,y,ans,t,En;

int fa[N],f[N],h1[N],h2[N],fg[N],out[N],Fa[N],d[N],h[N],tail[N];

struct edge {int s,n;} E1[N<<1],E2[N<<1];

struct Edge {int s,n,x,y;} E[N<<1];

struct cc {int x,f;} A[N];

void R(int &x)

{

	x=0;ch=getchar();

	while (ch<'0' || '9'<ch) ch=getchar();

	while ('0'<=ch && ch<='9') x=x*10+ch-'0',ch=getchar();

}

void W(int x)

{

	if (x>=10) W(x/10);

	putchar(x%10+'0');

}

int getf(int x)

{

	if (f[x]==x) return x;

	return f[x]=getf(f[x]);

}

void E_add1(int x,int y)

{

	E1[++En1].s=y;E1[En1].n=h1[x];h1[x]=En1;

	E1[++En1].s=x;E1[En1].n=h1[y];h1[y]=En1;

}

void E_add2(int x,int y)

{

	E2[++En2].s=y;E2[En2].n=h2[x];h2[x]=En2;

}

void E_add(int x,int y,int tx,int ty)

{

	if (!h[x]) tail[x]=En+1;

	E[++En].s=y;E[En].n=h[x];h[x]=En;

	E[En].x=tx;E[En].y=ty;

	if (!h[y]) tail[y]=En+1;

	E[++En].s=x;E[En].n=h[y];h[y]=En;

	E[En].x=tx;E[En].y=ty;

}

void dfs1(int x,int F)

{

	fa[x]=F;

	for (int k=h1[x];k;k=E1[k].n) if (E1[k].s!=F) dfs1(E1[k].s,x);

}

int main()

{

	R(n);ans=n-1;

	for (i=1;i<=n;i++) f[i]=i;

	for (i=1;i<n;i++)

	{

		R(x);R(y);

		E_add1(x,y);

	}

	dfs1(1,0);

	for (i=1;i<n;i++)

	{

		R(x);R(y);

		if (fa[x]==y && getf(x)!=getf(y))

		{

			f[f[x]]=f[y];ans--;

			continue;

		}

		if (fa[y]==x && getf(x)!=getf(y))

		{

			f[f[y]]=f[x];ans--;

			continue;

		}

		E_add2(x,y);

	}

	for (i=1;i<=n;i++) if (getf(i)!=getf(fa[i]))

	{

		fg[i]=1;

		if (!fa[i]) continue;

		A[f[i]].x=i;

		A[f[i]].f=fa[i];

		out[f[fa[i]]]++;

		Fa[f[i]]=f[fa[i]];

	}

	for (i=1;i<=n;i++)

		for (k=h2[i];k;k=E2[k].n) if (f[i]!=f[E2[k].s]) E_add(f[i],f[E2[k].s],i,E2[k].s);

	W(ans);puts("");

	j=0;

	for (i=1;i<=n;i++) if (fg[i] && !out[i]) d[++j]=i;

	i=0;

	while (i<j)

	{

		i++;

		if (d[i]==1) break;

		W(A[d[i]].x);putchar(' ');W(A[d[i]].f);putchar(' ');

		for (int &k=h[d[i]];k;k=E[k].n) if (getf(d[i])!=getf(E[k].s))

		{

			W(E[k].x);putchar(' ');W(E[k].y);puts("");

			E[tail[f[E[k].s]]].n=h[f[d[i]]];

			tail[f[E[k].s]]=tail[f[d[i]]];

			f[f[d[i]]]=f[E[k].s];

			k=E[k].n;

			break;

		}

		if (--out[Fa[d[i]]]==0) d[++j]=Fa[d[i]];

	}

}