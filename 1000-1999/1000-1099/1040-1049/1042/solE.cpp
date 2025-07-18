#include<cstdio>
#include<algorithm>
using namespace std;
int gi(){
	int x=0,w=1;char ch=getchar();
	while ((ch<'0'||ch>'9')&&ch!='-') ch=getchar();
	if (ch=='-') w=0,ch=getchar();
	while (ch>='0'&&ch<='9') x=(x<<3)+(x<<1)+ch-'0',ch=getchar();
	return w?x:-x;
}
const int N = 1005;
const int mod = 998244353;
struct node{
	int v,x,y;
	bool operator < (const node &b) const
		{return v<b.v;}
}a[N*N];
int n,m,nm,s,X,X2,Y,Y2,sum,inv[N*N],ans[N][N];
int main(){
	n=gi();m=gi();nm=n*m;
	for (int i=1;i<=n;++i)
		for (int j=1;j<=m;++j)
			a[(i-1)*m+j]=(node){gi(),i,j};
	sort(a+1,a+nm+1);
	inv[0]=inv[1]=1;
	for (int i=2;i<=nm;++i) inv[i]=1ll*inv[mod%i]*(mod-mod/i)%mod;
	for (int i=1,j=1;i<=nm;i=j=j+1){
		while (j<nm&&a[j+1].v==a[i].v) ++j;
		for (int k=i;k<=j;++k){
			ans[a[k].x][a[k].y]=(1ll*s*a[k].x*a[k].x%mod+X2+1ll*X*a[k].x%mod+1ll*s*a[k].y*a[k].y%mod+Y2+1ll*Y*a[k].y%mod+sum)%mod*inv[s]%mod;
		}
		for (int k=i;k<=j;++k){
			++s;
			X=(X-2*a[k].x+mod)%mod;X2=(X2+a[k].x*a[k].x)%mod;
			Y=(Y-2*a[k].y+mod)%mod;Y2=(Y2+a[k].y*a[k].y)%mod;
			sum=(sum+ans[a[k].x][a[k].y])%mod;
		}
	}
	int r=gi(),c=gi();
	printf("%d\n",ans[r][c]);
	return 0;
}