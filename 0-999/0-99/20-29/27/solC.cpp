#include<stdio.h>
int main()
{
	int n, a[100005];
	scanf("%d", &n);
	for(int i = 1; i <= n; i++)
		scanf("%d", &a[i]);
	int i = 2;
	while(a[i] == a[1])
	    i++;
	if(a[i] > a[1])
		for(int j = i + 1; j <= n; j++)
			if(a[j] < a[j - 1])
			{
				printf("3\n");
				printf("1 %d %d\n", j - 1, j);
				return 0;
			}
	if(a[i] < a[1])
		for(int j = i + 1; j <= n; j++)
			if(a[j] > a[j-1])
			{
				printf("3\n");
				printf("1 %d %d\n", j - 1, j);
				return 0;
			}
	printf("0\n");
	return 0;
}