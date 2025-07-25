#include <iostream>
#include <vector>
#include <algorithm>
#include <map>
#include <string>

using namespace std;

map<string,int> f;
int K, N, a[10200];
string s, t;
vector <int> v[10200];

int check ( int n )
{
	for (int i = 0; i < N; i += 1)
		a[i] = 0;
	for (int i = 0; i < v[n].size(); i += 1)
		a[v[n][i]] = 1;
	a[n] = 1;
	int m = 0, r = 0;
	for (int i = 0; i < v[n].size(); i += 1)
		for (int j = 0; j < v[v[n][i]].size(); j += 1)
			if(a[v[v[n][i]][j]]%2 == 0)
				a[v[v[n][i]][j]] += 2, m = max(m, a[v[v[n][i]][j]]);
	for (int i = 0; i < N; i += 1)
		if(a[i] == m)
			r++;
	return r;
}

int main (int argc, char const* argv[])
{
	cin >> K;

	for (int i = 0; i < K; i += 1)
	{
		cin >> s >> t;
		if (f.find(s) == f.end())
			f[s] = N++;
		if (f.find(t) == f.end())
			f[t] = N++;

		v[f[s]].push_back(f[t]);
		v[f[t]].push_back(f[s]);
	}
	cout << N << '\n';
	for (map<string,int>::iterator it = f.begin(); it != f.end(); it++)
	{
		cout << it->first << ' ' << check(it->second) << '\n';
	}
	
	return 0;
}