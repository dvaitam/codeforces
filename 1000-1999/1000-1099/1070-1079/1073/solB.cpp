#include<math.h>
#include<ctype.h>
#include<stdio.h>
#include<string.h>
#include<iostream>
#include<algorithm>
#define il inline
#define rg register
using namespace std;
const int O=200010;
il int read ( ){
	rg int o=0,fl=1;
	char ch=getchar();
	while(!isdigit(ch))fl^=ch=='-',ch=getchar();
	while(isdigit(ch))o=10*o+ch-48,ch=getchar();
	return fl?o:-o;
}
il void write(int x){
	if(x<0)putchar('-'),x=-x;
	if(x>9)write(x/10);
	putchar(x%10+'0');
}
int n,a[O],b,now;
int main ( ){
	n=read();
	for(int x,i=1;i<=n;i++)x=read(),a[x]=i;
	for(int i=1;i<=n;i++){
		b=read();
		write(max(0,a[b]-now));
		putchar(' ');
		if(now<a[b])now=a[b];
	}
	return 0;
}