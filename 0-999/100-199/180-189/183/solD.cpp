#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cctype>

#include<ctime>

#include<cstdlib>

#include<string>

#include<queue>

#include<cmath>

#include<set>

#include<map>

#include<bitset>

#define Rep(x,a,b) for (int x=a;x<=b;x++)

#define Drp(x,a,b) for (int x=a;x>=b;x--)

#define Cross(x,a) for (int x=head[a];~x;x=next[x])

#define ll long long

#define oo (1<<29)

#define mk(a,b) make_pair(a,b)

#define fr first

#define sc second

using namespace std;

inline int IN(){

	int x=0,ch=getchar(),f=1;

	while (!isdigit(ch)&&(ch!='-')&&(ch!=EOF)) ch=getchar();

	if (ch=='-'){f=-1;ch=getchar();}

	while (isdigit(ch)){x=(x<<1)+(x<<3)+ch-'0';ch=getchar();}

	return x*f;

}

inline void OUT(ll x){

	if (x<0) putchar('-'),x=-x;

	if (x>=10) OUT(x/10),putchar(x%10+'0');

		else putchar(x+'0');

}

const int N=3005,M=305;

int n,m,p;

double a[N][M],g[M][N][2],Ans;

bool G[M];

int main(){

	n=IN(),m=IN();

	Rep(i,1,n) Rep(j,1,m) a[i][j]=IN()/1000.0;

	Rep(i,1,m) Rep(j,1,n) g[i][j][G[i]]=(g[i][j-1][G[i]^1]+1)*a[j][i]+g[i][j-1][G[i]]*(1-a[j][i]);

	Rep(i,1,n){

		double Mx=0;

		Rep(j,1,m) if (g[j][n][G[j]]-g[j][n][G[j]^1]>Mx) Mx=g[j][n][G[j]]-g[j][n][G[j]^1],p=j;

		Ans+=Mx,G[p]^=1;

		Rep(j,1,n) g[p][j][G[p]]=(g[p][j-1][G[p]^1]+1)*a[j][p]+g[p][j-1][G[p]]*(1-a[j][p]);

	}

	printf("%.15lf\n",Ans);

	return 0;

}