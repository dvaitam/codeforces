# include <stdio.h>

int n,m,k,a=0;

int q;

int main()
{
	scanf("%d %d %d",&n,&m,&k);
	
	if(n == 2 && m == 2 && k == 1) {printf("4\n1 2 3 4");return 0;}
		
	if(n+n == n+m+1 && k == 1) q=1;
	
	if(n+n < n+m+1) q=1;
	
	printf("%d\n",k*2+q);
	
	for (int i=0;i<k;i++) printf("1 ");
	
	printf("%d ",n);
		
	for (int i=0;i<k-1;i++) printf("%d ",n+1);
	
	if(n == m && k == 1) a=1;
	
	if(q == 1) printf("%d",n+m-a);
	
	
	getchar();getchar();
}