#include <iostream>

#include <cstdio>

#include <cmath>



using namespace std;





int T;

double R, r, _r, K;



int main()

{

    scanf("%d", &T);



    while( T-- )

    {

        scanf("%lf %lf %lf", &R, &r, &K);



        _r = R-r;



        double x1 = (R*2+r*2)/2;

        double y1 = (R*2-r*2)*K;

        double d1 = sqrt(x1*x1+y1*y1);



        printf("%.9f\n", R*2*r/(d1-_r)-R*2*r/(d1+_r) );

    }

}