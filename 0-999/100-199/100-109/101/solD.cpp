#include<cstdio>

#include<algorithm>

#define ll long long

using namespace std;

const int N=1e5+5,M=2*N;

int v,et,he[N];

struct edge{int l,to,v;}e[M];

int n,i;

int read(){

	char c=getchar();int k=0;for (;c<48||c>57;c=getchar());

	for (;c>47&&c<58;c=getchar()) k=(k<<3)+(k<<1)+c-48;return k;

}

void line(int x,int y){

	e[++et].l=he[x];he[x]=et;

	e[et].to=y;e[et].v=v;

}

int size[N],sum[N];ll time[N];

struct arr{int size,sum;}q[N];

bool operator < (arr A,arr B){

	return (ll)A.sum*B.size<(ll)B.sum*A.size;

}

void dfs(int x,int fa){

	size[x]=1;int m=0;

	for (int i=he[x];i;i=e[i].l){

		int y=e[i].to;if (y==fa) continue;

		sum[y]=e[i].v<<1;dfs(y,x);

		sum[x]+=sum[y];size[x]+=size[y];

	}

	for (int i=he[x];i;i=e[i].l){

		int y=e[i].to;if (y==fa) continue;

		time[x]+=time[y]+e[i].v*size[y];

		q[++m]=(arr){size[y],sum[y]};

	}

	sort(q+1,q+m+1);ll s=0;

	for (int i=1;i<=m;s+=q[i].sum,i++)

		time[x]+=q[i].size*s; 

}

int main(){

	for (n=read(),i=1;i<n;i++){

		int x=read(),y=read();v=read();

		line(x,y);line(y,x);

	}

	dfs(1,0);

	printf("%.10lf",time[1]/(double)(n-1));

}