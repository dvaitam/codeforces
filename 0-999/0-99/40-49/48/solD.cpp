#include<stdio.h>
int a[100001]={0},b[100001]={0};
int main()
{
	int n,i;
	bool t;
	scanf("%d",&n);
	for (i=1;i<=n;i++)
	{
	 	scanf("%d",&a[i]);
		b[a[i]]++;
    }
    t=true;
	for (i=2;i<=100000;i++) if (b[i]>b[i-1]) {t=false; break;}
	if (t)
	{
		printf("%d\n",b[1]);
		for (i=1;i<=n;i++)
		{
			printf("%d ",b[a[i]]);
			b[a[i]]--;
		}
		printf("\n");
	}
	else
	{
		printf("-1\n");
	}
}