#include <cstdio>
#include <cstring>
#include <iostream>
#include <algorithm>
using namespace std;
#define ll long long
#define N 200005
int n;
char A[N], B[N];
int main()
{
	scanf("%d", &n);
	scanf("%s%s", A + 1, B + 1);
	int s1 = 0, s2 = 0;
	for (int i = 1; i <= n; i++)
	{
		if (A[i] == 'a' && B[i] == 'b')
			s1++;
		if (A[i] == 'b' && B[i] == 'a')
			s2++;
	}
	int ans = s1 / 2 + s2 / 2;
	s1 %= 2, s2 %= 2;
	if (s1 ^ s2)
		puts("-1");
	else
	{
		printf("%d\n", ans + s1 + s2);
		int l1 = 0, l2 = 0;
		for (int i = 1; i <= n; i++)
		{
			if (A[i] == 'a' && B[i] == 'b')
			{
				if (l1)
					printf("%d %d\n", l1, i), l1 = 0;
				else
					l1 = i;
			}
			if (A[i] == 'b' && B[i] == 'a')
			{
				if (l2)
					printf("%d %d\n", l2, i), l2 = 0;
				else
					l2 = i;
			}
		}
		if (s1)
		{
			printf("%d %d\n", l1, l1);
			printf("%d %d\n", l1, l2);
		}
	}
}