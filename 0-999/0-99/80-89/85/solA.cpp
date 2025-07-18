#include <string>
#include <vector>
#include <cstdio>
#include <cmath>
#include <cstdlib>
#include <cctype>
#include <algorithm>
#include <map>
#include <iostream>
#include <sstream>
#include <queue>
#include <cstring>
#include <ctime>

using namespace std;

#define SZ 105

char dat[4][SZ];

int main()
{
	int n, i;
	scanf("%d", &n);
	if(n % 2 == 0)
	{
		int j = 0;
		for(i = 0;i < n - 1; i += 2)
		{
			dat[0][i] = dat[0][i + 1] = 'a' + j;
			dat[3][i] = dat[3][i + 1] = 'g' + j;
			j = (j + 1) % 2;
		}
		j = 0;
		for(i = 1; i < n - 2; i += 2)
		{
			dat[1][i] = dat[1][i + 1] = 'c' + j;
			dat[2][i] = dat[2][i + 1] = 'e' + j;
			j = (j + 1) % 2;
		}
		dat[1][0] = dat[2][0] = 'y';
		dat[1][n - 1] = dat[2][n - 1] = 'z';
	}
	else
	{
		int j = 0;
		for(i = 0;i < n - 1; i += 2)
		{
			dat[0][i] = dat[0][i + 1] = 'a' + j;
			dat[1][i] = dat[1][i + 1] = 'g' + j;
			j = (j + 1) % 2;
		}
		j = 0;
		for(i = 1; i <= n - 2; i += 2)
		{
			dat[2][i] = dat[2][i + 1] = 'c' + j;
			dat[3][i] = dat[3][i + 1] = 'e' + j;
			j = (j + 1) % 2;
		}
		dat[2][0] = dat[3][0] = 'y';
		dat[0][n - 1] = dat[1][n - 1] = 'z';
	}
	dat[0][n] = dat[1][n] = dat[2][n] = dat[3][n] = '\0';
	for(i = 0; i < 4; i++)
	{
		printf("%s\n", dat[i]);
	}
	return 0;
}