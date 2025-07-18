//A. Bmail Computer Network 
#include<cstdio>
void in(int &x){
	int flag=1;char ch=getchar();x=0;
	while(ch<'0'||ch>'9'){if(ch=='-')flag=-flag;ch=getchar();}
	while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
}
void write(int x){
	if(x<0){putchar('-');x=-x;}
	if(x>9)write(x/10);
	putchar(x%10+'0');
}
int n,l,now,a[200005],ans[200005];
int main(){
	in(n);
	for(int i=2;i<=n;i++)
		in(a[i]);
	now=n;
	while(now!=1){
		ans[++l]=now;
		now=a[now];
	}
	write(now);
	for(int i=l;i>=1;i--){
		putchar(' ');
		write(ans[i]);
	}
	return 0;
}