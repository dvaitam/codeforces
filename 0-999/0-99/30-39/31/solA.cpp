#include <stdio.h>
#include <string.h>
#include <iostream>
#include <algorithm>
using namespace std;

int a[105];

int main()
{	
	//freopen("A.in", "r", stdin);

	int n, i, j, k;
	scanf("%d", &n);
	for (i = 0; i < n; i++)
		scanf("%d", &a[i]);
	int flag = 1;
	for (i = 0; flag && i < n; i++)
		for (j = i + 1; flag && j < n; j++)
			for (k = j + 1; flag && k < n; k++)
			{
				if (a[i] + a[j] == a[k])
				{
					printf("%d %d %d\n", k + 1, j + 1, i + 1);
					flag = 0;
				}
				if (flag && a[i] + a[k] == a[j])
				{
					printf("%d %d %d\n", j + 1, k + 1, i + 1);
					flag = 0;
				}
				if (flag && a[k] + a[j] == a[i])
				{
					printf("%d %d %d\n", i + 1, j + 1, k + 1);
					flag = 0;
				}
			}
	if (flag) printf("-1\n");
	return 0;
}