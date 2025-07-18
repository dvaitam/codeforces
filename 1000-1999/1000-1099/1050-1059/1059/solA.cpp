#include<math.h> 
#include<stdio.h>
#include<ctype.h>
#include<string.h> 
#define il inline
#define rg register
il int read ( ){
	rg int o=0,fl=1;char ch=getchar();
	while(!isdigit(ch))o^=ch=='-',ch=getchar();
	while(isdigit(ch))o=10*o+ch-48,ch=getchar();
	return fl?o:-o;
}
int n,L,a,l,r,ans;
int main ( ){
	n=read(),L=read(),a=read();
	for(int i=1;i<=n;i++){
		l=read();
		ans+=(l-r)/a;
		r=l+read();
	}
	ans+=(L-r)/a;
	printf("%d\n",ans);
	return 0;
}