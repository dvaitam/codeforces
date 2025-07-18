#include<cstdio>
#include<cstring>
#include<algorithm>
#include<queue>
#include<cmath>
#define For(i,a,b) for(int i=a,_=b;i<=_;++i) 
#define Fdn(i,b,a) for(int i=b;i>=a;--i)
#define sqr(x) ((x)*(x))
#define y1 y11
#define bin(x) (1<<((x)-1))
#define ful(x) ((1<<(x))-1)
using namespace std;
typedef long long LL;
typedef unsigned long long ULL;
const int MAXN=2005, MAXM=200005, mo=1e9+7, INF=0x3f3f3f3f;
const LL LINF=0x3f3f3f3f3f3f3f3fLL;
typedef int ArrN[MAXN];
const double eps=1e-12, pi=acos(-1);

int m,n,X,a[MAXN],b[MAXN],ans;
int f[MAXN],g[MAXN];

int main(){
	scanf("%d%d",&n,&m);
	For(i,1,n) scanf("%d",a+i), a[i]+=a[i-1];
	For(i,1,m) scanf("%d",b+i), b[i]+=b[i-1];
	scanf("%d",&X);
	
	For(i,1,n){
		f[i]=2e9+1;
		For(j,i,n){
			f[i]=min(f[i],a[j]-a[j-i]);
		}
	}
	For(i,1,m){
		g[i]=2e9+1;
		For(j,i,m){
			g[i]=min(g[i],b[j]-b[j-i]);
		}
	}		
	For(i,1,n) For(j,1,m) if((LL)f[i]*g[j]<=X) ans=max(ans,i*j);
	printf("%d\n",ans);
	
    return 0; 
}