#include<cstdio>
#define rep(i,n) for (int i=1;i<=n;++i)
const int N=2005;
double f[N][N]; bool a[N],b[N]; int n,m,A,B,x,y;
int main()
{
	scanf("%d%d",&n,&m);
	rep(i,m) scanf("%d%d",&x,&y),a[x]=b[y]=1;
	rep(i,n) A+=a[i],B+=b[i];
	for (int i=n;i>=A;--i) for (int j=n;j>=B;--j) if (i<n || j<n)
		f[i][j]=((n-i)*j*f[i+1][j]+i*(n-j)*f[i][j+1]+(n-i)*(n-j)*f[i+1][j+1]+n*n)/(n*n-i*j);
	printf("%.9lf\n",f[A][B]); return 0;
}