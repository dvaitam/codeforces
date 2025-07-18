#include<cstdio>
#include<algorithm>
#include<string>
#include<map>
using namespace std;
string P[110], AA, BB, U[1010];
int n, BS, m;
string R[2222];
int Num[3][222], Last[222];
long long K, C1[2222], C2[2222], TC[2222];
map<string, int>Map, Map2, MNum;
void Print(string a, int L) {
	if (a.size() < L)printf("%s\n", a.c_str());
	else {
		for (int i = 0; i < L; i++)printf("%c", a[i]);
		puts("");
	}
}
int main() {
	int i, j;
	scanf("%d%lld%d", &n, &K, &m);
	P[0] = "0";
	P[1] = "1";
	for (i = 0;; i++) {
		P[i + 2] = P[i] + P[i + 1];
		if (P[i + 1].size() >= 200) {
			BS = i + 1;
			break;
		}
	}
	AA = P[BS];
	BB = P[BS + 1];
	for (i = 0; i < 5; i++)P[BS + i + 2] = P[BS + i] + P[BS + i + 1];
	if (n <= BS + 3) {
		for (i = 0; i < P[n].size(); i++) {
			U[i] = P[n].substr(i);
		}
		sort(U, U + P[n].size());
		string r = U[K - 1];
		Print(r, m);
		return 0;
	}
	int sz = P[BS + 4].size();
	for (i = 0; i < sz; i++) {
		Map2[P[BS + 4].substr(i, min(sz - i, 200))] = 1;
	}
	int cnt = 0;
	for (auto &x : Map2) {
		MNum[x.first] = ++cnt;
		R[cnt] = x.first;
	}
	for (i = 0; i < sz; i++) {
		C2[MNum[P[BS + 4].substr(i, min(sz - i, 200))]]++;
	}
	int sz3 = P[BS + 3].size();
	for (i = 0; i < sz3; i++) {
		C1[MNum[P[BS + 3].substr(i, min(sz3 - i, 200))]]++;
	}
	for (i = 1; i <= 200; i++) {
		Last[i] = MNum[P[BS + 4].substr(sz - i)];
	}
	int sA = AA.size(), sB = BB.size();
	for (i = 1; i < 200; i++) {
		Num[0][i] = MNum[BB.substr(sB - i) + AA.substr(0, 200 - i)];
		Num[1][i] = MNum[BB.substr(sB - i) + BB.substr(0, 200 - i)];
	}
	long long INF = 2e18;
	for (i = BS + 5; i <= n; i++) {
		for (j = 1; j <= cnt; j++)TC[j] = min(C1[j] + C2[j], INF);
		if (i % 2 == (BS + 5) % 2) {
			for (j = 1; j < 200;j++)TC[Last[j]]--, TC[Num[0][j]]++;
		}
		else {
			for (j = 1; j < 200;j++)TC[Last[j]]--, TC[Num[1][j]]++;
		}
		for (j = 1; j <= cnt; j++) {
			C1[j] = C2[j];
			C2[j] = TC[j];
		}
	}
	long long s = 0;
	for (i = 1; i <= cnt; i++) {
		s += C2[i];
		if (s >= K) {
			Print(R[i], m);
			return 0;
		}
	}
	return 0;
}