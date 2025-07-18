#include<bits/stdc++.h>
const int N=2e5+5;
#define re register 
#define ll long long
using namespace std;
template<class T>
inline void read(T &x)
{
	x=0; int f=1;
	static char ch=getchar();
	while(!isdigit(ch))
	{
		if(ch=='-')	f=-1;
		ch=getchar();
	}
	while(isdigit(ch)) x=x*10+ch-'0',ch=getchar();
	x*=f;
}
inline void write(int x)
{
	if(x<0)
	{
		putchar('-');
		x=-x;
	}
	if(x>9)	write(x/10);
	putchar(x%10+'0');
}
struct Node
{
	ll val; int pos;
	bool operator <(const Node &p) const
	{
		return this->val>p.val;
	}
}a[N],b[N];
int n,m,k,vis[N];
int main()
{
	read(n),read(m),read(k);
	for(re int i=1;i<=n;i++) read(a[i].val),a[i].pos=i,b[i]=a[i];
	int tar=m*k; ll sum=0;
	sort(b+1,b+n+1);
	for(int i=1;i<=tar;i++) vis[b[i].pos]=1,sum=sum+b[i].val;
	cout<<sum<<endl;
	for(int i=1,cnt=0,t=0;i<=n;i++)
	{
		if(t==k-1)	break;
		if(vis[i]) cnt++;
		if(cnt==m) {t++; write(i); putchar(' '); cnt=0;}
	}
	return 0;
}