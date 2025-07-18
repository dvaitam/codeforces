#include<cstdio>
#include<iostream>
#include<algorithm>
#include<cstring>
#include<string>
#include<cmath>
#include<queue>
#include<map>
#include<bitset>
#include<deque>
#include<cstdlib>
#include<set>
#define N 2000003
#include<ctime>
#define ll long long
#define mp make_pair
using namespace std;
ll read()
{
	ll x=0,f=1;
	char c=getchar();
	while(c>'9'||c<'0')
	{
		if(c=='-') f=-1;
		c=getchar();
	}
	while(c>='0'&&c<='9')
	{
		x=x*10+c-'0';
		c=getchar();
	}
	return f*x;
}
int n,a[N],wei[N][2];
ll ans;
int main()
{
	n=read();
	wei[0][0]=wei[0][1]=1;
	for(int i=1;i<=n*2;++i) 
	{
		int x=read();
		if(wei[x][0]) wei[x][1]=i;
		else wei[x][0]=i;
	}
	for(int i=1;i<=n;++i) ans+=min(abs(wei[i][1]-wei[i-1][1])+abs(wei[i][0]-wei[i-1][0]),abs(wei[i][0]-wei[i-1][1])+abs(wei[i][1]-wei[i-1][0]));
	cout<<ans;
	return 0;
 }