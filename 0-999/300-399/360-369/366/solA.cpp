#include <cstdio>
#include <iostream>
#include <algorithm>

using namespace std;

int main(){
	int n, g[5][5], pos, i, a, b;
	scanf("%d\n", &n);
	for( i=1 ; i<=4 ; i++ )
		scanf("%d %d %d %d\n", &g[i][1], &g[i][2], &g[i][3], &g[i][4]);

	bool flag = false;
	for( i=1 ; i<=4 ; i++ ){	
		if( min( g[i][1], g[i][2] ) + min( g[i][3], g[i][4] ) <= n ){
			flag = true;
			pos = i;
			a = min( g[i][1], g[i][2] );
			b = n - min( g[i][1], g[i][2] );
			break;
		}
	}
	if( flag )
		printf("%d %d %d\n", i, a, b );
	else printf("-1\n");
	return 0;
}