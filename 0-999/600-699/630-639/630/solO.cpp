#include <iostream>
#include <cstdio>
#include <algorithm>
#include <vector>
#include <map>
#include <queue>
#include <cmath>

using namespace std;

#define y1 roman_kaban
#define rank oryshych_konb
//const int mod = int(1e9) + 7;

const int MX = 200500;

int main()
{
    //freopen("in.txt","r", stdin);
    ios_base::sync_with_stdio(false);
    double px, py, vx, vy;
    cin >> px >> py >> vx >> vy;
    int a, b, c, d;
    cin >> a >> b >> c >>d;
    double t = sqrt(vx * vx + vy * vy);
    vx /= t;
    vy /= t;
    double vx1 = vy;
    double vy1 = -vx;
    printf("%.10f %.10f\n", px + vx * b, py + vy * b);
    printf("%.10f %.10f\n", px - vx1 * a / 2, py - vy1 * a / 2);
    printf("%.10f %.10f\n", px - vx1 * c / 2, py - vy1 * c / 2);
    printf("%.10f %.10f\n", px - vx1 * c / 2 - vx * d, py - vy1 * c / 2 - vy * d);
    printf("%.10f %.10f\n", px + vx1 * c / 2 - vx * d, py + vy1 * c / 2 - vy * d);
    printf("%.10f %.10f\n", px + vx1 * c / 2, py + vy1 * c / 2);
    printf("%.10f %.10f\n", px + vx1 * a / 2, py + vy1 * a / 2);

    return 0;
}