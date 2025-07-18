#include<bits/stdc++.h> 

using namespace std; 



using ld = long double;

using ll = long long;

using ull = unsigned long long;

using pii = pair<int, int>;

using pll = pair<ll, ll>;

using str = string;



#define vc vector

#define sz(a) (int)a.size()

#define all(A) A.begin(), A.end()

#define pb push_back

#define mp make_pair

#define fri(start , end , step)  for(int i = (int)start; i < (int)end ; i = i + step)

#define frj(start , end , step)  for(int j = (int)start; j < (int)end ; j = j + step)

#define ff first

#define ss second

#define forn(i, n) for (int i = 0; i < (int)n; ++i)

#define fast ios_base::sync_with_stdio(0); cin.tie(0)







#ifdef HOME

    #define debug(x) cerr<<#x<<": "<<x<<endl;

#else

    #define debug(x) ;

#endif 











mt19937 rnd(chrono::steady_clock::now().time_since_epoch().count());



const int MOD = 1e9 + 7;

const ld e = 1/1e10;

#define int ll



template<typename T>

ostream& operator<< (ostream &other, const vector<T> &v) {

    for (const T &x: v) other << x << ' ';

    other << '\n';

    return other;

}





inline int pw(int a, int b){

    if(!b) return 1;

    if(b & 1) return (1ll * a * pw(a ,  b - 1)) % MOD;

    int x = pw(a , b / 2);

    return (1ll * x * x) % MOD;

}



template<typename T> 

inline bool cmin(T& a , const T& b){

    if(a > b){

        a = b;

        return 1;

    }

    return 0;

}



template<typename T> 

inline bool cmax(T& a , const T& b){

    if(a < b){

        a = b;

        return 1;

    }

    return 0;

}







const int maxn = 1e6;

const ll inf = 1e12;





inline void run(){

    int n;

    cin >> n;

    vc<int> A1(n) , A2(n);

    forn(i , n) cin >> A1[i];

    forn(i , n) cin >> A2[i];

    sort(all(A1));

    sort(all(A2));

    bool f = 1;

    forn(i , n){

        if(A1[i] + A2[n - 1 - i] != A1[0] + A2[n - 1]) f = 0;

    }

    if(f){

        cout << "YES\n";

        

        for(int x : A1) cout << x << " ";

        cout << "\n";

        cout << 0 << " " <<  A1[0] + A2[n - 1] << "\n";

        return;

    }

    if(A1 == A2){

        cout << "YES\n";

       

        for(int x : A1) cout << x << " ";

        cout << "\n";

        cout << 0 << " " << 0 << "\n";

        return;

    }

    

    auto check = [&](int d , vc<int> &ans){

        debug(d);

        multiset<int> B1 , B2;

        forn(i , n){

            B1.insert(A1[i]);

            B2.insert(A2[i]);   

        }

        int h = 0;

        while(sz(B1)){

            int x = *B1.rbegin();

            int y = *B2.rbegin(); 

            if(max(x , y) <= d) break;

            if(x > y){

                auto it = B2.find(x - d);

                if(it == B2.end()){

                    return 0;

                }

                B1.erase(prev(B1.end()));

                B2.erase(it);

                ans[h++] = x; 



            }

            else{

                auto it = B1.find(y - d);

                if(it == B1.end()){

                    return 0;

                }

                B2.erase(prev(B2.end()));

                B1.erase(it);

                ans[h++] = -(y - d); 

            }

        }

        

        if(sz(B1) == 0) return 1;

        debug(sz(B1));

        debug(sz(B2));

        auto it1 = B1.begin();

        auto it2 = prev(B2.end());

        

        while(it1 != B1.end()){

            if(*it1 + *it2 != d) return 0;

            ans[h++] = *it1;

            it1++;

            if(it1 != B1.end()) it2--;

        }

        return 1;

    };

    int kek = 1e9;

    forn(i , n){

        vc<int> ans(n);

        if(A2[n - 1] > A1[n - 1]) {

            if(check(A2[n - 1] - A1[i] , ans)){

                cout << "YES\n";

                

                for(int x : ans) cout << x + kek << " ";

                cout << "\n";

                cout << kek << " " << A2[n - 1] - A1[i] + kek << "\n";

                return;

            }

        }

        else {

            if(check(A1[n - 1] - A2[i] , ans)){

                cout << "YES\n";

               

                for(int x : ans) cout << x + kek<< " ";

                cout << "\n";

                cout << kek << " " << A1[n - 1] - A2[i] + kek << "\n";

                return;

            }

            

        }

    }

    cout << "NO\n";

}







signed main(){

    fast;

    int tt;

    tt = 1;

    cin >> tt;

    

    while(tt--){

        run();

    }

    #ifdef HOME

        cerr << "time :" << fixed << setprecision(3) << ((ld) clock() / (double) CLOCKS_PER_SEC) << endl;

    #endif

    return 0;

}