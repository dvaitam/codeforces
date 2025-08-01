#include<cstdio>
#include<algorithm>
using namespace std;
int i,j,x,y,N,k,h,Ans,p,len;
struct node
{
	int l,r,len;
}A[2005];
int P[2005],Q[2005],f[2005],g[2005],link[2005],Out[2005];
int res[2005][2005];
int cmp1(const int &i,const int &j)
{
	return A[i].len<A[j].len;
}
int cmp2(const int &i,const int &j)
{
	return A[i].l<A[j].l;
}
void Get(int x)
{
	int i;
	Out[++len]=x;
	for(i=1;i<res[x][0];++i)
		Get(res[x][i]);
}
int main()
{
	scanf("%d",&N);
	for(i=1;i<=N;++i)
	{
		scanf("%d%d",&x,&y);
		A[i].l=x-y+1,A[i].r=x+y;
	}
	for(i=1;i<=N;++i) A[i].len=A[i].r-A[i].l+1;
	for(i=1;i<=N;++i) P[i]=i;
	sort(P+1,P+N+1,cmp1);
	for(i=1;i<=N;++i) Q[i]=i;
	sort(Q+1,Q+N+1,cmp2);
	for(k=1;k<=N;++k)
	{
		i=P[k];
		for(j=1;j<=N;++j) f[j]=2000000005;
		for(j=1;j<=N;++j) link[j]=0;
		f[0]=A[i].l-1,p=0;
		for(h=1;h<=N;++h)
		{
			j=Q[h];
			if(A[j].r>A[i].r||A[j].l<A[i].l) continue;
			for(;f[p+1]<A[j].l;++p);
			if(A[j].r<f[p+g[j]])
			{
				f[p+g[j]]=A[j].r;
				link[p+g[j]]=j;
			}
		}
		for(j=N;f[j]==2000000005;--j);
		for(h=j;h;h-=g[link[h]]) res[i][++res[i][0]]=link[h];
		res[i][++res[i][0]]=i;
		g[i]=j+1;
	}
	for(j=1;j<=N;++j) f[j]=2000000005;
	for(j=1;j<=N;++j) link[j]=0;
	f[0]=A[1].l-1,p=0;
	for(h=1;h<=N;++h)
	{
		j=Q[h];
		for(;f[p+1]<A[j].l;++p);
		if(A[j].r<f[p+g[j]])
		{
			f[p+g[j]]=A[j].r;
			link[p+g[j]]=j;
		}
	}
	for(j=N;f[j]==2000000005;--j);
	printf("%d\n",j);
	for(h=j;h;h-=g[link[h]]) Get(link[h]);
	sort(Out+1,Out+j+1);
	for(i=1;i<j;++i) printf("%d ",Out[i]);
	printf("%d\n",Out[j]);
	return 0;
}