#include<iostream>
#include<iomanip>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<string>
#include<algorithm>
#include<cmath>
#include<queue>
#include<stack>
using namespace std;
inline void read(int &x)
{
	x=0;char s=getchar();
	bool flag=false;
	while (s<'0'||s>'9')
	{
		if (s=='-')
			flag=true;
		s=getchar();
	}
	while ('0'<=s&&s<='9')
	{
		x=x*10+(s-'0');
		s=getchar();
	}
	if (flag)
		x=-x;
}
int n;
int x,y;
int ans=0;
int main()
{
	//freopen("������.in","r",stdin);
	//freopen("������.out","w",stdout);
	read(n);
	for (int i=1;i<=n;i++)
	{
		read(x),read(y);
		ans=max(ans,x+y);
	}
	printf("%d",ans);
	return 0;
}