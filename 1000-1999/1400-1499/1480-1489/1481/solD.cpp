#include<bits/stdc++.h>



using namespace std;



// #include <ext/pb_ds/assoc_container.hpp>

// #include<ext/pb_ds/tree_policy.hpp>

// using namespace __gnu_pbds;

// using namespace __gnu_cxx;

// template <typename T> using oset = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>;



#define IO ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)

#define   int                      long long

typedef   double                   dd;

typedef   vector<int>              vi;

typedef   vector<char>             vc;

typedef   set<int>                 si;

typedef   map<int,int>             mp;

typedef   pair<int,int>            pr;

typedef   tuple<int,int, int>      tp;

typedef   vector<vector<int>>      vvi;

typedef   vector<pair<int,int>>    vpr;

#define   intmax            INT_MAX

#define   lintmax           LLONG_MAX

#define   endl              "\n"

#define   pb                push_back

#define   F                 first

#define   S                 second

#define   pi                acos(-1.0)

#define   size(a)           (int)a.size()

#define   all(v)            v.begin(), v.end()

#define   rall(v)           v.rbegin(), v.rend()

#define   bp(x)             __builtin_popcountll(x)



#define dbg(args...) { string _s = #args; replace(_s.begin(), _s.end(), ',', ' '); stringstream _ss(_s); istream_iterator<string> _it(_ss); err(_it, args); cerr << '\n';}

void err(istream_iterator<string> it) {}

template<typename T, typename... Args>

void err(istream_iterator<string> it, T a, Args... args) {

   cerr << *it << "= " << a << " ! ";

   err(++it, args...);

}



//int ax[] = {0,  0, -1, 1, -1,  1, -1, 1};

//int ay[] = {1, -1,  0, 0, -1, -1,  1, 1};



const double eps = 1e-8;

const int N = 1e6 + 5;



int ii = 1;



void solve() {

   int n, k; cin >> n >> k;

   char ch[n + 4][n + 4];

   for (int i = 1; i <= n; i++)

      for (int j = 1; j <= n; j++) cin >> ch[i][j];



   bool f = false;

   int n1 = 1, n2 = 2;

   for (int i = 1; i <= n; i++) {

      for (int j = 1; j <= n; j++) {

         if(i == j) continue;

         if(ch[i][j] == ch[j][i]) {f = true; n1 = i; n2 = j;}

      }

   }





   if((k&1) || f) {

      cout << "YES" << endl;

      f = true;

      cout << n2 << ' ';

      while(k--) {

         cout << (f ? n1 : n2) << ' ';

         f ^= 1;

      }

      cout << endl;

      return;

   }



   vi apc;

   for (int i = 1; i <= n; i++) {

      apc.clear();

      bool f = false;

      int a = -1, b = -1;

      for (int j = 1; j <= n; j++) {

         if(i == j) continue;

         if(ch[i][j] == 'a') a = j;

         if(ch[i][j] == 'b') b = j;

      }

      for (int j = 1; j <= n; j++) {

         if(i == j) continue;

         if(ch[j][i] == 'a' && a != -1) {

            apc.pb(j);

            apc.pb(i);

            apc.pb(a);

            f = true;

            break;

         }

         if(ch[j][i] == 'b' && b != -1) {

            apc.pb(j);

            apc.pb(i);

            apc.pb(b);

            f = true;

            break;  

         }

      }

      if(f) break;

   }



   if(size(apc) > 0) {

      cout << "YES" << endl;

      int m = k / 2;

      // n /= 2;

      if(m&1) {

         bool f = true;

         cout << apc[0] << ' ';

         while(m--) {

            cout << (f ? apc[1] : apc[0]) << ' ';

            f ^= 1;

         }

         f = false;

         m = k/2;

         // cout << apc[2] << ' ';

         while(m--) {

            cout << (f ? apc[1] : apc[2]) << ' ';

            f ^= 1;

         }

      } else {

         bool f = true;

         cout << apc[1] << ' ';

         while(m--) {

            cout << (f ? apc[0] : apc[1]) << ' ';

            f ^= 1;

         }

         f = false;

         m = k/2;

         // cout << apc[2] << ' ';

         while(m--) {

            cout << (f ? apc[1] : apc[2]) << ' ';

            f ^= 1;

         }

      }



      cout << endl;

   } else {

      cout << "NO" << endl;

   }







}



int32_t main() {

   IO;

   int TestCase = 1; 



   // #ifndef ONLINE_JUDGE

   // freopen("D:\\Atom\\sublime text 3\\input.in", "r", stdin);

   // freopen("D:\\Atom\\sublime text 3\\output.out", "w", stdout);

   // #endif



   cin >> TestCase;

   // scanf("%d", &TestCase);

   while(TestCase--)

      solve();



   //cerr << endl << "time : "<< (float)clock()/CLOCKS_PER_SEC << " s" << endl;



   return 0;

}







 

/**

   stuff you should look for

   * obvious case

   * check MIN/MAX value

   * int overflow, array bounds

   * special cases (n=1?)

   * do smth instead of nothing and stay organized

   * WRITE STUFF DOWN

   * DON'T GET STUCK ON ONE APPROACH

*/