#include <bits/stdc++.h>

using namespace std;

#define int long long

double rast(double x1, double y1, double x2, double y2) {
	return sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1));
}

signed main() {
	ios_base::sync_with_stdio(false);
	cin.tie(0);

	double a, b, c;
	cin >> a >> b >> c;
	double x1, y1, x2, y2;
	cin >> x1 >> y1 >> x2 >> y2;
	double ans1 = fabs(x1 - x2) + abs(y1 - y2);
	if (a != 0 && b != 0) {
		double tmp1y = (-c - a * x1) / b;
		double tmp1x = (-c - b * y1) / a;
		double tmp2y = (-c - a * x2) / b;
		double tmp2x = (-c - b * y2) / a;
		double r1x = rast(x1, y1, tmp1x, y1);
		double r1y = rast(x1, y1, x1, tmp1y);
		double r2x = rast(x2, y2, tmp2x, y2);
		double r2y = rast(x2, y2, x2, tmp2y);
		//cout << r1x << ' ' << r1y << ' ' << r2x << ' ' << r2y << ' ' << tmp1x << ' ' << tmp1y << ' ' << tmp2x << ' ' << tmp2y << '\n';
		double ans2 = min(r1x + r2x + rast(tmp1x, y1, tmp2x, y2), min(r1x + r2y + rast(tmp1x, y1, x2, tmp2y), min(r1y + r2x + rast(x1, tmp1y, tmp2x, y2), r1y + r2y + rast(x1, tmp1y, x2, tmp2y))));
		ans1 = min(ans1, ans2);
		cout << fixed << setprecision(12) << ans1;
		return 0;
	}
	else if (a == 0) {
		double tmp1y = (-c - a * x1) / b;
		double tmp2y = (-c - a * x2) / b;
		double r1y = rast(x1, y1, x1, tmp1y);
		double r2y = rast(x2, y2, x2, tmp2y);
		//cout << r1x << ' ' << r1y << ' ' << r2x << ' ' << r2y << ' ' << tmp1x << ' ' << tmp1y << ' ' << tmp2x << ' ' << tmp2y << '\n';
		double ans2 = r1y + r2y + rast(x1, tmp1y, x2, tmp2y);
		ans1 = min(ans1, ans2);
		cout << fixed << setprecision(12) << ans1;
		return 0;
	}
	else if (b == 0) {
		double tmp1x = (-c - b * y1) / a;
		double tmp2x = (-c - b * y2) / a;
		double r1x = rast(x1, y1, tmp1x, y1);
		double r2x = rast(x2, y2, tmp2x, y2);
		//cout << r1x << ' ' << r1y << ' ' << r2x << ' ' << r2y << ' ' << tmp1x << ' ' << tmp1y << ' ' << tmp2x << ' ' << tmp2y << '\n';
		double ans2 = r1x + r2x + rast(tmp1x, y1, tmp2x, y2);
		ans1 = min(ans1, ans2);
		cout << fixed << setprecision(12) << ans1;
		return 0;
	}
	cout << fixed << setprecision(12) << ans1;
	return 0;
}