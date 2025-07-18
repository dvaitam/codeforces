#include <cstdio>
#include <math.h>
using namespace std;
int n, m;
int main()
{
	scanf ("%d", &n);
	int k;
	m=1;
	k=sqrt(n*m+1);
	while (sqrt((n*m)+1)!=k) 
	{
		m++;
		k=sqrt((n*m)+1);
	}
	printf ("%d", m);
}