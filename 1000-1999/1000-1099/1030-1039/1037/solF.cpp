#include <bits/stdc++.h>
using namespace std;
typedef long long LL;
const int N=1000005,mod=1e9+7;
LL read(){
	LL x=0,f=1;
	char ch=getchar();
	while (!isdigit(ch)&&ch!='-')
		ch=getchar();
	if (ch=='-')
		f=-1,ch=getchar();
	while (isdigit(ch))
		x=(x<<1)+(x<<3)+ch-48,ch=getchar();
	return x*f;
}
int n,k,a[N],pre[N],nxt[N];
int st[N],top;
void Get_Pre(){
	st[top=1]=0;
	for (int i=1;i<=n;i++){
		while (a[i]>=a[st[top]])
			top--;
		pre[i]=st[top];
		st[++top]=i;
	}
}
void Get_Nxt(){
	st[top=1]=n+1;
	for (int i=n;i>=1;i--){
		while (a[i]>a[st[top]])
			top--;
		nxt[i]=st[top];
		st[++top]=i;
	}
}
int calc(int L,int R){
//	printf("%d %d\n",L,R);
	int a=L/k,b=R/k;
	int x=L%k,y=R%k;
	int ans=a+b;
	ans=(1LL*a*k%mod*b+ans)%mod;
	ans=(1LL*x*b+ans)%mod;
	ans=(1LL*y*a+ans)%mod;
	if (x==0||y==0)
		return ans;
	ans=(ans+max(0,x+y-k+1))%mod;
	return ans;
}
int main(){
	n=read(),k=read()-1;
	for (int i=1;i<=n;i++)
		a[i]=read();
	a[0]=a[n+1]=1e9+1;
	Get_Pre();
	Get_Nxt();
//	for (int i=1;i<=n;i++)
//		printf("%d %d\n",pre[i],nxt[i]);
	int ans=0;
	for (int i=1;i<=n;i++)
		ans=(1LL*calc(i-pre[i]-1,nxt[i]-i-1)*a[i]+ans)%mod;
	printf("%d",ans);
	return 0;
}