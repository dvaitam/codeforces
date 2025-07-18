#include<bits/stdc++.h>



using namespace std;



//Constants

typedef long long ll;

typedef pair<long long, long long> pii;

const ll mod7 = 1000000007;



#define fori(x) for (ll i = 0; i < x; i++)

#define forj(x) for (ll j = 0; j < x; j++)

#define fork(x) for (ll k = 0; k < x; k++)



#define rep(u, s, x) for (ll u = s; u < x; u++)

#define rem(u, x, s) for (ll u = s - 1; u >= x; u--)



#define DBG(x) cout << (x) << endl;

#define DBGT(x, t) cout << (t) << "" << (x) << endl;

#define DBGARRAY(x, s, t) cout << t << ""; for (unsigned int z = 0; z < s; z++) cout << x[z] << " "; cout << endl;

#define DBGVECTOR(x, t) cout << t << ""; for (unsigned int z = 0; z < x.size(); z++) cout << x[z] << " "; cout << endl;

#define DBGVECTORPAR(x, t, p) cout << t << ""; for (unsigned int z = 0; z < x.size(); z++) cout << x[z].p << " "; cout << endl;



ll GCD(ll a, ll b) { return b == 0 ? a : GCD(b, a % b); }

ll tonumber(string stringNum) { ll a; stringstream convert; convert << stringNum; convert >> a; return a; }

string tostring(ll num) { string s; stringstream convert; convert << num; convert >> s; return s; }



//Structs



//Global Variables

ll A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z;

ll temp;

char tempChar;

string tempString;



//Function Declarations



int main()

{

    cin >> N >> L;



    vector<ll> nums;

    fori(N)

    {

        cin >> A;

        nums.push_back(A);

    }



    sort(nums.begin(), nums.end());

    //DBGVECTOR(nums, "nums: ");



    double biggest = max(nums.front(), L - nums.back());

    fori(N - 1)

    {

        biggest = max(biggest, (nums[i + 1] - nums[i]) / 2.0);

    }



    cout << setprecision(12) << fixed << biggest;

}