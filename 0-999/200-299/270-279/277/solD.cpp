#include <iostream>
#include <fstream>
#include <iomanip>
#include <algorithm>
using namespace std;

struct problem
{
	long long scoresmall, scorelarge;
	long long timesmall, timelarge;
	long long probpass, probfail;

	bool operator < (const problem& a) const{
		return probfail*a.probpass*timelarge < a.probfail*probpass*a.timelarge;
	}
};

struct results
{
	long double score, time;
	results(): score(0), time(0) {}
};

void getmax(results& a, results b)
{
	if (a.score < b.score) a = b;
	else if (a.score == b.score && a.time > b.time) a.time = b.time;
}

int main() {
	ios_base::sync_with_stdio(false);
	cin.tie(0);
	int n, t;
	cin >> n >> t;

	double d;
	problem p[n];
	for (int i = 0; i < n; ++i) {
		cin >> p[i].scoresmall >> p[i].scorelarge >> p[i].timesmall >> p[i].timelarge >> d;
		p[i].probfail = d*1e6+0.5;
		p[i].probpass = 1e6-p[i].probfail;
	}
	sort(p, p+n);

	int k;
	results r[t+1], tmp;
	for (int i = 0; i < n; ++i) {
		for (int j = t-p[i].timesmall; j >= 0; --j) {
			if (j != 0 && r[j].score == 0) continue;
			k = j+p[i].timesmall;
			tmp.score = r[j].score + p[i].scoresmall*1e6;
			tmp.time = r[j].time + p[i].timesmall*1e6;
			getmax(r[k], tmp);

			k += p[i].timelarge;
			if (k > t) continue;
			tmp.score += p[i].probpass * p[i].scorelarge;
			tmp.time += p[i].probpass * (k-tmp.time/1e6);
			getmax(r[k], tmp);
		}
	}
	for (int i = 0; i <= t; ++i)
		getmax(tmp, r[i]);
	cout << fixed << setprecision(9) << tmp.score/1e6 << ' ' << tmp.time/1e6;
	return 0;
}