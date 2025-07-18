#include <stdio.h>
#include <algorithm>
#define MAXN 60

using namespace std;

int m1[MAXN][MAXN], m2[MAXN][MAXN], nlin1, ncol1, nlin2, ncol2, maxsum=0, maxx=0, maxy=0;

void testxy (int x, int y)
{
        int sum = 0;
        for (int i = max(0,-x); i < nlin1 && i+x < nlin2; i++)
                for (int j = max(0,-y); j < ncol1 && j+y < ncol2; j++)
                        sum += m1[i][j] * m2[i+x][j+y];
        if (sum > maxsum)
                maxsum = sum, maxx = x, maxy = y;
}


int main()
{
	scanf(" %d %d", &nlin1, &ncol1);
	for(int i=0; i<nlin1; i++)
		for(int j=0; j<ncol1; j++)
		{
			char bf;
			scanf(" %c", &bf);
			m1[i][j] = (int)bf-(int)'0';
		}
	scanf(" %d %d", &nlin2, &ncol2);
	for(int i=0; i<nlin2; i++)
		for(int j=0; j<ncol2; j++)
		{
			char bf;
			scanf(" %c", &bf);
			m2[i][j] = (int)bf-(int)'0';		
		}
	
	for (int x = min(-nlin2, -nlin1); x < max(nlin2,nlin1); x++)
                for (int y = min(-ncol2, -ncol1); y < max(ncol2,ncol1); y++)
                        testxy (x, y);

	printf("%d %d", maxx, maxy);

	return 0;
}