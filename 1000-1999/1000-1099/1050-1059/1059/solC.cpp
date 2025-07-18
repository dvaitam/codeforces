#include<cstdio>
#define ll long long
#define dd double
using namespace std;
void write(int x){
	if(x<0) x=-x,putchar('-');
	if(x>9) write(x/10);
	putchar(x%10^48);
}
int main(){
	int n,x,s=1,num=0;
	scanf("%d",&n);
	if(n==1){printf("1");return 0;}
	if(n==3){printf("1 1 3");return 0;}
	x=n;
	while(x){
		x-=num,num=x+1>>1;
		if(x==1)
			s>>=1,write(n-n%s);
		else for(int i=1;i<=num;i++) write(s),putchar(' ');
		s<<=1;
	}
	return 0;
}