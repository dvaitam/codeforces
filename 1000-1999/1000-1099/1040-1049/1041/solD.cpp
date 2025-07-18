#include<cstdio>
#include<algorithm>
#define maxn 200005
using namespace std;
int n,h,ans,sum;
struct chj{
	int x,y;
}a[maxn];
inline char nc(){
	static char buf[100000],*L=buf,*R=buf;
	return L==R&&(R=(L=buf)+fread(buf,1,100000,stdin),L==R)?EOF:*L++;
}
inline int read(){
	int ret=0,f=1;char ch=nc();
	while (ch<'0'||ch>'9'){if (ch=='-') f=-1;ch=nc();}
	while (ch>='0'&&ch<='9') ret=ret*10+ch-'0',ch=nc();
	return ret*f;
}
int main(){
	n=read();h=read();
	for (int i=1;i<=n;i++) a[i]=(chj){read(),read()};
	int i=1,j=1;
	sum=a[1].y-a[1].x;
	while (h){
		if (j==n) break;
		int t=a[j+1].x-a[j].y;
		if (h-t>0) h-=t;else break;
		j++;
		sum+=t+a[j].y-a[j].x;
	}
	ans=sum+h;
	if (j==n){printf("%d\n",ans);return 0;}
	while (i<=n){
		i++;
		sum-=a[i].x-a[i-1].x;
		h+=a[i].x-a[i-1].y;
		while (j+1<=n){
			int t=a[j+1].x-a[j].y;
			if (h-t>0) h-=t;else break;
			j++;
			sum+=t+a[j].y-a[j].x;
		}
		if (j==n) break;
		ans=max(ans,sum+h);
	}
	ans=max(ans,sum+h);
	printf("%d\n",ans);
	return 0;
}