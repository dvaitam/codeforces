#include<algorithm>
#include<iostream>
#include<cstring>
#include<vector>
#include<cstdio>
#include<cmath>
#include<queue>
#include<map>
#include<set>
using namespace std;
#define MN 3000
#define K 12
inline int read()
{
    int x=0,f=1;char ch=getchar();
    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
    return x*f;
}
pair<int,int> a[MN+5],f[K][MN+5];
int n,lg[MN+5],ans[MN+5],d1=-1e9,p1,d2=-1e9,p2,d3=-1e9,p3;
pair<int,int> rmq(int l,int r)
{
	int x=lg[r-l+1];
	return max(f[x][l],f[x][r-(1<<x)+1]);
}
int main()
{
    n=read();
    for(int i=1;i<=n;++i)a[i]=make_pair(-read(),i);
    sort(a+1,a+n+1);
    for(int i=1;i<=n;++i)f[0][i]=make_pair(a[i+1].first-a[i].first,i);
    for(int i=1;i<K;++i)for(int j=1;j<=n;++j)f[i][j]=max(f[i-1][j],f[i-1][j+(1<<i-1)]);
    for(int i=3;i<=n;++i)lg[i]=lg[i+1>>1]+1;
    for(int i=1;i<=n;++i)for(int j=i;++j<=n;)if(i<=2*(j-i)&&j-i<=2*i)
    {
    	int l=max(1,max(i,j-i)+1>>1),r=min(n-j,min(i,j-i)*2);
    	if(l<=r&&j+l<=n)
    	{
    		pair<int,int> x=rmq(j+l,j+r);
    		if(a[i+1].first-a[i].first>d1||
			  (a[i+1].first-a[i].first==d1&&a[j+1].first-a[j].first>d2)||
			  (a[i+1].first-a[i].first==d1&&a[j+1].first-a[j].first==d2&&x.first>d3))
    			d1=a[i+1].first-a[i].first,p1=i,d2=a[j+1].first-a[j].first,p2=j,d3=x.first,p3=x.second;
		}
	}
	for(int i=1;i<=p1;++i)ans[a[i].second]=1;
	for(int i=p1;++i<=p2;)ans[a[i].second]=2;
	for(int i=p2;++i<=p3;)ans[a[i].second]=3;
	for(int i=p3;++i<=n;)ans[a[i].second]=-1;
	for(int i=1;i<=n;++i)printf("%d ",ans[i]); 
    return 0;
}