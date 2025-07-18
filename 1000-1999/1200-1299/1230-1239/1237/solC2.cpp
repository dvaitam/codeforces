#include<iostream>
#include<cstdio>
#include<algorithm>
#include<cstring>
#include<cmath>
#include<set>
#include<map>
using namespace std;
typedef long long ll;
inline int read()
{
	int x=0,f=1; char ch=getchar();
	while(ch<'0'||ch>'9') { if(ch=='-') f=-1; ch=getchar(); }
	while(ch>='0'&&ch<='9') { x=(x<<1)+(x<<3)+(ch^48); ch=getchar(); }
	return x*f;
}
const int N=5e4+7;
int n;
struct dat {
	int x,y,z,id;
	inline bool operator < (const dat &tmp) const {
		if(x!=tmp.x) return x<tmp.x;
		return y!=tmp.y ? y<tmp.y : z<tmp.z;
	}
}A[N];
bool vis[N];
int main()
{
	n=read();
	for(int i=1;i<=n;i++)
		A[i].x=read(),A[i].y=read(),A[i].z=read(),A[i].id=i;
	sort(A+1,A+n+1);
	int l=1;
	for(int i=2;i<=n;i++)
		if(l && A[i].x==A[l].x && A[i].y==A[l].y)
		{
			vis[i]=vis[l]=1;
			printf("%d %d\n",A[i].id,A[l].id);
			l=0;
		}
		else l=i;
	l=1; while(vis[l]) l++;
	for(int i=l+1;i<=n;i++)
	{
		if(vis[i]) continue;
		if(l && A[i].x==A[l].x)
		{
			vis[i]=vis[l]=1;
			printf("%d %d\n",A[i].id,A[l].id);
			l=0;
		}
		else l=i;
	}
	l=1; while(vis[l]) l++;
	for(int i=l+1;i<=n;i++)
	{
		if(vis[i]) continue;
		if(l) printf("%d %d\n",A[i].id,A[l].id),l=0;
		else l=i;
	}
	return 0;
}