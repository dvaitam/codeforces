#include <bits/stdc++.h>

using namespace std;



int h,m;



int main()

{

	scanf("%d:%d",&h,&m);

	printf("%.9lf %.9lf\n",(h%12)*30+m/2.0,m*6.0);

}