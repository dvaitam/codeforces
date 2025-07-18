#include <cstdio>
#include <iostream>
#include <cmath>
#include <cstring>
#include <cctype>
using namespace std;
double ans=-1,x[6000],y[6000],z[6000],d[6000];
int n;
void qsort (int l, int r) {
	int i=l, j=r, m=(l+r)/2;
	double w=d[m];
	while (i<=j) {
		while (d[i]<w) i++;
		while (d[j]>w) j--;
		if (i<=j) {
			double tmp=d[i];
			d[i]=d[j];
			d[j]=tmp;
			tmp=x[i];
			x[i]=x[j];
			x[j]=tmp;
			tmp=y[i];
			y[i]=y[j];
			y[j]=tmp;
			tmp=z[i];
			z[i]=z[j];
			z[j]=tmp;
			i++;
			j--;
		}
	}
	if (i<r) qsort(i,r);
	if (j>l) qsort(l,j);
}
double sqr (double x) {
	return x*x;
}
int main () {
#ifndef ONLINE_JUDGE
	freopen ("input.txt","r",stdin);
	freopen ("output.txt","w",stdout);
#endif
	cin>>n;
	for (int i=1; i<=n; i++) cin>>x[i]>>y[i]>>z[i]; 
	for (int i=1; i<=n; i++) d[i]=sqrt(sqr(x[i]-x[1])+sqr(y[i]-y[1])+sqr(z[i]-z[1]));
	qsort(1,n);
	for (int i=2; i<=n; i++) 
		for (int j=i+1; j<=n; j++) {
			double when=0;
			double dist=sqrt(sqr(x[i]-x[j])+sqr(y[i]-y[j])+sqr(z[i]-z[j]));
			dist-=d[j]-d[i];
			when=d[j]+dist/2;
			if (ans<0 || when<ans) ans=when;
		}
	printf("%.19lf",ans);
}