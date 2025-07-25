#include<cstdio>
#include<cstring>
#include<cmath>
#include<vector>
#include<utility>
#define X first
#define Y second
#define mp make_pair
#define pb push_back
#define rep(i,n) for(int i=0;i<n;i++)
#define inf 1e9
#define eps 1e-6
using namespace std;
int x[444],y[444];
double slk[888],lx[888],ly[888];
bool visx[888],visy[888];
double g[888][888];
int match[888];
int n,m,flw;
inline double sqr(int x){
	return x*x*1.0;
}
inline double dist(int i,int j){
	return -sqrt(sqr(x[i]-x[j])+sqr(y[i]-y[j]));
}
inline void add_edge(int i,int j,double d){
	g[i][j]=g[i][j+n]=d;
	(d>lx[i]) && (lx[i]=d);
	return;
}
inline bool hun(int u){
	visx[u]=1;
	rep(v,m)
		if(g[u][v]!=-inf){
			double d=lx[u]+ly[v]-g[u][v];
			if(d>eps){
				(d<slk[v]) && (slk[v]=d);
				continue;
			}
			if(visy[v])
				continue;
			visy[v]=1;
			if(match[v]<0 || hun(match[v])){
				match[v]=u;
				return 1;
			}
		}
	return 0;
}
void km(){
	memset(match,-1,sizeof(match));
	rep(i,n){
		while(1){
			rep(j,m)
				slk[j]=inf;
			memset(visx,0,sizeof(visx));
			memset(visy,0,sizeof(visy));
			if(hun(i))
				break;
			double d=inf;
			rep(j,m)
				if(!visy[j])
					(slk[j]<d) && (d=slk[j]);
			if(d>=inf){
				flw++;
				break;
			}
			rep(j,m){
				if(visx[j])
					lx[j]-=d;
				if(visy[j])
					ly[j]+=d;
			}
		}
		if(flw>1)
			break;
	}
	if(flw>1){
		puts("-1");
		return;
	}
	double ans=0.0;
	rep(i,m)
		if(match[i]>-1)
			ans-=g[match[i]][i];
	printf("%.6lf\n",ans);
	return;
}
int main(){
	scanf("%d",&n);
	m=n<<1;
	rep(i,n)
		rep(j,m)
			g[i][j]=-inf;
	rep(i,n)
		scanf("%d%d",x+i,y+i);
	rep(i,n){
		lx[i]=-inf;
		rep(j,n)
			if(y[i]<y[j]){
				double d=dist(i,j);
				add_edge(i,j,d);
			}
	}
	km();
	return 0;
}