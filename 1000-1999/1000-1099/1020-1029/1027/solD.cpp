#include<iostream>
#include<cstdio>
#include<cstdlib>
#include<cstring>
#include<cmath>
#include<algorithm>
#define N 200005
using namespace std;
#define _(d) while(d((ch=getchar()-48)>=0))
int get(){
	char ch;_(!);int v=ch;_() v=(v<<1)+(v<<3)+ch;return v;
}
int n,c[N],a[N],top,s[N],ans;
bool f[N],used[N];
void work(int x){
	while(1){
		if(f[x]){
			if(!used[x]){
				while(top) used[s[top--]]=0;
				return;
			}
			int res=1e5;
			while(s[top+1]!=x) res=min(res,c[s[top]]),used[s[top]]=0,top--;
			ans+=res;
			while(top) used[s[top--]]=0;
			return;
		}
		f[x]=1;
		used[x]=1;
		s[++top]=x;
		x=a[x];
	}
}
int main(){
	n=get();
	for(int i=1;i<=n;i++) c[i]=get();
	for(int i=1;i<=n;i++) a[i]=get();
	for(int i=1;i<=n;i++)
		if(!f[i]) work(i);
	printf("%d\n",ans);
	return 0;
}