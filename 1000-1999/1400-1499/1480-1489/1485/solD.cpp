#include<bits/stdc++.h>
using namespace std;
#define N 505
#define For(i,x,y)for(i=x;i<=(y);i++)
int a[N][N];
int gcd(int _,int __)
{
	return(!_?__:gcd(__%_,_));
}
inline int lcm(int _,int __)
{
	return _/gcd(_,__)*__;
}
int read()
{
	int A;
	bool K;
	char C;
	C=A=K=0;
	while(C<'0'||C>'9')K|=C=='-',C=getchar();
	while(C>'/'&&C<':')A=(A<<3)+(A<<1)+(C^48),C=getchar();
	return(K?-A:A);
}
void write(int X)
{
	if(X<0)putchar('-'),X=-X;
	if(X>9)write(X/10);
	putchar(X%10|48);
}
int main()
{
	int n,m,i,j,s=1;
	n=read(),m=read();
	For(i,1,n)
	For(j,1,m)a[i][j]=read(),s=lcm(s,a[i][j]);
	For(i,1,n)
	{
		For(j,1,m)write((i&1^j&1?a[i][j]*a[i][j]*a[i][j]*a[i][j]:0)+s),putchar(' ');
		putchar('\n');
	}
	return 0;
}