#include<cmath>
#include<cstdio>
using namespace std;
const int N=1e5+5;
int n;
double px[N],pn[N];
double a[N],b[N],f[N],g[N];
void Init(){
	scanf("%d",&n);
	for(int i=1;i<=n;i++) scanf("%lf",px+i);
	for(int i=1;i<=n;i++) scanf("%lf",pn+i);
}
double sqr(double x){return x*x;}
int main(){
	Init();
	for(int i=1;i<=n;i++){
		double t=sqr(px[i]+pn[i]-f[i-1]+g[i-1])-4*(px[i]-f[i-1]*(px[i]+pn[i]));
		if(t<0) t=0;
		a[i]=(px[i]+pn[i]-f[i-1]+g[i-1]-sqrt(t))/2;
		b[i]=px[i]+pn[i]-a[i];
		f[i]=f[i-1]+a[i];
		g[i]=g[i-1]+b[i];
	}
	for(int i=1;i<=n;i++) printf("%.7f%c",a[i],i==n?'\n':' ');
	for(int i=1;i<=n;i++) printf("%.7f%c",b[i],i==n?'\n':' ');
	return 0;
}