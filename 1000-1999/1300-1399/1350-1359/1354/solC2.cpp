/**

 *      Author:  Richw818

 *      Created: 02.06.2023 15:32:16

**/



#include <bits/stdc++.h>

using namespace std;



int main(){

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr);

    const double PI = acos(-1.0);

    int t; cin >> t;

    while(t--){

        int n; cin >> n;

        int consider = (n / 2) + (n / 2) % 2;

        double polyAngle = PI * (2 * n - 2) / (2 * n);

        double cRadius = 1 / (2 * cos(polyAngle / 2));

        double considerAngle = 2 * PI * consider / (2 * n);

        double part1 = cRadius * cos(considerAngle / 2);

        double part2 = cRadius * sin(considerAngle / 2);

        double ans = (part1 + part2) * sqrt(2);

        cout << fixed << setprecision(10) << ans << '\n';

    }

    return 0;

}