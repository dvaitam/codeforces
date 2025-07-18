#include <iostream>

#include <cstdlib>

#include <cstdio>

#include <iomanip>

#include <cmath>

#include <algorithm>

#include <queue>

#include <map>

#include <string>

#include <vector>

#include <cstring>

#include <set>



#define forn(i,n) for(int i = 0; i < (n); i++)

#define forsn(i,s,n) for(int i = (s); i < (n); i++)

#define all(v) ((v).begin, (v).end)

#define pb push_back

#define x first

#define y second

#define mp make_pair



using namespace std;



typedef pair<int,int> par;

typedef long long int tint;



double p[110];

double prod[110];

double revprod[110];



int main()

{

	int n; cin >> n;

	forn(i,n) cin >> p[i];

	sort(p,p+n);



	cout.setf(ios::fixed | ios::showpoint);

	cout.precision(15);



	if(fabs(p[n-1] - 1.0) < 1e-9) cout << 1.0 << endl;

	else

	{

		prod[0] = 1.0; forn(i,n) prod[i+1] = (1.0 - p[i]) * prod[i];

		revprod[n] = 1.0; for(int i = (n - 1); i >= 0; i--) revprod[i] = (1.0 - p[i]) * revprod[i+1];



		double best = 0.0;



		forn(i,n+1)

		{

			for(int j = n; j >= i; j--)

			{

				double pr = prod[i] * revprod[j];

				double cnd = 0.0;

				forn(k,i) cnd += p[k] * pr / (1.0 - p[k]);

				for(int k = n - 1; k >= j; k--) cnd += p[k] * pr / (1.0 - p[k]);



				//cout << i << " " << j << endl;

				//cout << cnd << endl;

				best = max(best, cnd);

			}

		}

		cout << best << endl;

	}



	return 0;

}