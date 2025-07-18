#include<cstdio>
int main(){
	int n;
	scanf("%d",&n);
	int c=n%10;
	n-=c;
	if(c>=5)n+=10;
	printf("%d\n",n);
	return 0;
}