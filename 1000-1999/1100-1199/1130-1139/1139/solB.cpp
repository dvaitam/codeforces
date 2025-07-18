#include<cstdio>
using namespace std;
int X,W;char ch;
inline int read()
{
	X=0,W=1;ch=getchar();
	while(ch<'0'||ch>'9'){if(ch=='-')W=-1;ch=getchar();}
	while(ch>='0'&&ch<='9'){X=(X<<1)+(X<<3)+ch-48;ch=getchar();}
	return X*W;
}
int a[200001];
long long ans;
int main()
{
	int n=read();
	for(int i=1;i<=n;i++)a[i]=read();
	int h=a[n];
	ans=h;
	for(int i=n-1;i>0&&h;i--){
		if(h>a[i])h=a[i],ans+=(long long)a[i];
		else h--,ans+=(long long)h;
	}
	printf("%lld",ans);
	return 0;
}