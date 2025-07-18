#include<iostream>
#include<cstdio>
#include<ctype.h>
#include<algorithm>
using namespace std;
inline int read(){
    int x=0,f=0;char ch=getchar();
    while(!isdigit(ch))f|=ch=='-',ch=getchar();
    while(isdigit(ch))x=x*10+(ch^48),ch=getchar();
    return f?-x:x;
}
long long a[100007];
double sum=0.0,ans=0.0;
int main(){
	long long n=read(),k=read(),m=read();
	for(int i=1;i<=n;++i)a[i]=read(),ans+=a[i];ans/=n;
	sort(a+1,a+1+n);
	for(int i=n;i>=1;--i){
		sum+=a[i];
		long long o=m-i+1;
		if(o<0)continue;
		long long add=min((n-i+1)*k,o);
		ans=max(ans,(sum+add)/(n-i+1));
	}
	printf("%.10lf",ans);
    return 0;
}