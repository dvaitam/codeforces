#include<cmath>
#include<cstdio>
#include<cstring>
#include<algorithm>
#define gt getchar()
#define ll long long
#define File(s) freopen(s".in","r",stdin),freopen(s".out","w",stdout)
inline int in()
{
	int k=0;char ch=gt;
	while(ch<'-')ch=gt;
	while(ch>'-')k=k*10+ch-'0',ch=gt;
	return k;
}
int a[100005];
int main()
{
	int n=in(),m=in(),k=in(),ans=n;
	for(int i=1;i<=n;++i)a[i]=in();
	for(int i=n;i>=2;--i)ans+=(a[i]-=a[i-1]+1);
	std::sort(a+2,a+n+1);
	for(int i=1;i<k;++i)ans-=a[n-i+1];printf("%d\n",ans);
	return 0;
}