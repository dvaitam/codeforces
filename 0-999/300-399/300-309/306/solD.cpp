#include <cstdio>
#include <cstring>
#include <cmath>
#include <algorithm>
#include <iostream>
#include <vector>
#include <complex>
using namespace std;
#define X real()
#define Y imag()
typedef long long ll;
typedef complex<double> cd;
const double PI = acos(-1.0);
const int MAX = 1e3+5;
const int base = 690;
cd point[MAX];
cd u, v;
cd v1, v2;
int n;
double a1, b1, c1;
double a2, b2, c2;
double det, XX1, YY1;
double l[MAX];
double uu[MAX];

int main()
{
    scanf("%d", &n);
    if ( n < 5 )
        printf("No solution\n");
    else {
        u = cd( cos(2.0*PI/n), sin(2.0*PI/n) );
        v = cd( base, 0.0);
        point[1] = cd(0.0, 0.0);

        for(int i=2; i<=n-1; ++i) {
            point[i] = point[i-1] + v;
            v = v*u;
            v = v*(abs(v) - 0.01)/abs(v);
        }

        v1 = v;
        v2 = v*u;

        a1 = v1.X;
        b1 = v2.X;
        c1 = point[1].X - point[n-1].X;
        a2 = v1.Y;
        b2 = v2.Y;
        c2 = point[1].Y - point[n-1].Y;
        det = a1*b2 - a2*b1;
        XX1 = (c1*b2 - c2*b1) / det;
        YY1 = (a1*c2 - a2*c1) / det;
        point[n] = point[n-1] + XX1*v1;

        point[n+1] = point[n];
        for(int i=1; i<=n; ++i) printf("%.9f %.9f\n", point[i].X, point[i].Y);
    }
    return 0;
}