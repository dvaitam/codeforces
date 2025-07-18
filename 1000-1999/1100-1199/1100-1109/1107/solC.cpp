#include<iostream> 
#include<cstdio>
#include<cmath>
#include<cstdlib>
#include<cstring>
#include<algorithm>
using namespace std;
#define ll long long
#define N 200010
char getc(){char c=getchar();while ((c<'A'||c>'Z')&&(c<'a'||c>'z')&&(c<'0'||c>'9')) c=getchar();return c;}
int gcd(int n,int m){return m==0?n:gcd(m,n%m);}
int read()
{
	int x=0,f=1;char c=getchar();
	while (c<'0'||c>'9') {if (c=='-') f=-1;c=getchar();}
	while (c>='0'&&c<='9') x=(x<<1)+(x<<3)+(c^48),c=getchar();
	return x*f;
}
int n,m,a[N];
ll ans;
char s[N];
signed main()
{
#ifndef ONLINE_JUDGE
	freopen("a.in","r",stdin);
	freopen("a.out","w",stdout);
	const char LL[]="%I64d\n";
#endif
	n=read(),m=read();
	for (int i=1;i<=n;i++) a[i]=read();
	scanf("%s",s+1);
	for (int i=1;i<=n;i++)
	{
		int t=i;
		while (t<n&&s[t+1]==s[i]) t++;
		sort(a+i,a+t+1);reverse(a+i,a+t+1);
		for (int j=i;j<=t&&j<i+m;j++) ans+=a[j];
		i=t;
	}
	cout<<ans;
	return 0;
	//NOTICE LONG LONG!!!!!
}
//ѡһ��ȫ1���У�ʹ�����м�ֵ��+�����ּ�ֵ�����
//g[i][j]ǰiλ��ѡ���г���Ϊjʱ������ֵ 
//ת�ƿ���ö����һ��1��λ�� 
//���߽����һλ��1ѡ��