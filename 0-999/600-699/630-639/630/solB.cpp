#include <bits/stdc++.h>

using namespace std;

#include <cstdio>
#include <cmath>

int main()
{
    long long n, i, t;
    double r=1.0, x;
    scanf("%lld %lld", &n, &t);
    if (t==0) printf ("%f", (n+0.0));
    else {
        x=1.000000011;
        i=2;
        while (i<=t) {
            for (i=2; i<=t; i*=2) x*=x;
            r*=x;
            i/=2;
            t-=i;
            i=2;
            x=1.000000011;
        }
        for (i=0; i<t; i++) r*=1.000000011;
        r*=n;
        printf("%f", r);
    }
    return 0;
}