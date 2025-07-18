#include<cstdio>
#include<cctype>
#include<algorithm>
using namespace std;

#define re register int
#define ll long long

const int INF=0x3f3f3f3f;

int read()
{
	int x=0,f=0;
	char ch=getchar();
	while(!isdigit(ch))
	{
		f|=ch=='-';
		ch=getchar();
	}
	while(isdigit(ch))
	{
		x=(x<<1)+(x<<3)+(ch^48);
		ch=getchar();
	}
	return f?-x:x;
}
inline void print(int x)
{
	if(x==0)
	{
		putchar('0');
		return;
	}
	else if(x<0)
		putchar('-'),x=-x;
	int stk[111],top=0;
	while(x)
	{
		stk[++top]=x%10;
		x/=10;
	}
	for(re i=top;i;i--)
		putchar(stk[i]|48);
}

int n,k;
struct node
{
	int x,w;
	bool operator<(const node &a) const
	{
		return w>a.w;
	}
}digit[200100];

bool check(int mid)
{
	int ans=0;
	for(re i=1;i<=200000&&digit[i].w;i++)
		ans+=digit[i].w/mid;
	return ans>=k;
}

int main()
{
	n=read(),k=read();
	for(re i=1;i<=200000;i++)
		digit[i].x=i;
	for(re i=1;i<=n;i++)
		digit[read()].w++;
	sort(digit+1,digit+200001);
	int l=1,r=200000,mid=(l+r)>>1,cnt=0;
	while(l<=r)
	{
		if(check(mid)) l=mid+1;
		else r=mid-1;
		mid=(l+r)>>1;
	}
	for(re i=1;i<=200000&&digit[i].w>=r;i++)
		for(re j=1;j<=digit[i].w/r&&cnt<k;cnt++,j++)
			print(digit[i].x),putchar(' ');
	putchar('\n');
	return 0;
}