// Author : R.Raj

#include <bits/stdc++.h>
using namespace std;


#define fastio ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)
#define nl "\n"
#define inf 1e18
#define yes      cout<<"YES"<<nl
#define no       cout<<"NO"<<nl
#define couts(n) cout<<n<<" "
#define coutn(n) cout<<n<<nl
#define pb push_back
#define ff first
#define ss second
#define all(a) a.begin(),a.end()
#define rall(a) a.rbegin(),a.rend()
#define forl(i,a1,a2) for(i=a1;i<a2;i++)
#define forle(i,a1,a2) for(i=a1;i<=a2;i++)
#define forlr(i,a1,a2) for(i=a1;i>=a2;i--)
#define rep while
#define vsort(a) (sort(all(a)))
#define revsort(a) (sort(rall(a)))
#define vprint(a) for(auto it : a) couts(it); coutn("");

#define sum(a)     ( accumulate ((a).begin(), (a).end(), 0LL))
#define mine(a)    (*min_element((a).begin(), (a).end()))
#define maxe(a)    (*max_element((a).begin(), (a).end()))
#define mini(a)    ( min_element((a).begin(), (a).end()) - (a).begin())
#define maxi(a)    ( max_element((a).begin(), (a).end()) - (a).begin())
#define cnt(a,x)   ( count((a).begin(), (a).end(),(x) ))
#define lob(a, x)  (*lower_bound((a).begin(), (a).end(), (x)))
#define upb(a, x)  (*upper_bound((a).begin(), (a).end(), (x)))


#ifndef ONLINE_JUDGE
#define debug(x) cerr << #x <<" "; _print(x); cerr << endl;
#else
#define debug(x)
#endif


typedef long long ll;
typedef unsigned long long ull;
typedef long double lld;

void _print(int t) {cerr << t;}
void _print(ll t) {cerr << t;}
void _print(string t) {cerr << t;}
void _print(char t) {cerr << t;}
void _print(lld t) {cerr << t;}
void _print(double t) {cerr << t;}
void _print(ull t) {cerr << t;}

// template <class T, class V> void _print(pair <T, V> p);
// template <class T> void _print(vector <T> v);
// template <class T, class V> void _print(map <T, V> v);
// template <class T, class V> void _print(unordered_map <T, V> v);
// template <class T> void _print(set <T> v);
// template <class T> void _print(unordered_set <T> v);
// template <class T> void _print(multiset <T> v);
// template <class T> void _print(unordered_multiset <T> v);


template <class T, class V> void _print(pair <T, V> p) {cerr << "{"; _print(p.first); cerr << ","; _print(p.second); cerr << "}";}
template <class T> void _print(vector <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T, class V> void _print(map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T, class V> void _print(unordered_map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}

template <class T> void _print(set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T> void _print(unordered_set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T> void _print(multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T> void _print(unordered_multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}


#define int long long

typedef map<int,int> mapii;
typedef pair<int,int> pii;
typedef vector<int> vi;
typedef vector<pii> vpi;
typedef vector<vi> vvi;


template <class T> istream & operator>> (istream &in, vector<T> &v) {
    for (auto &vi : v) 
        in >> vi;
    return in;
}

template <class T> ostream & operator<< (ostream &out, vector<T> &v) {
    for (auto &vi : v) 
        out << vi << " ";
    return out;
}


void solve(){
    int n,m,k;
    int i,j;
    int x,y;
    
    cin >> n;
    
    vi arr(n);
    cin >> arr;

    set<int> s;
    forl(i , 0 , n)
        s.insert(arr[i]);
    if (s.size() == 1){
        coutn(0);
        return;
    }

    x = cnt(arr , 1);
    if (x > 0){
        coutn(-1);
        return;
    }

    vpi ans;
    int mn , ind; 
    rep(1){
        mn = mine(arr);
        ind = mini(arr);

        debug(mn);debug(ind);

        forl(i , 0 , n){
            // cerr << " ######### " << nl;
            rep(arr[i] > mn){
                ans.pb({i+1 , ind+1});
                debug(arr[i]);
                arr[i] = (arr[i]+mn-1)/mn;
            }
        }

        s.clear();
        forl(i , 0 , n)
            s.insert(arr[i]);
        if (s.size() == 1)
            break;
    }

    debug(arr);


    // forl(i , 0 , n){
    //     if (arr[i] == 2){
    //         ind = i;
    //         break;
    //     }
    // }

    // forl(i , 0 , n){
    //     rep(arr[i] > 2){
    //         ans.pb({i+1 , ind+1});
    //         // arr[i] = (arr[i]+1)/2;
    //         arr[i] /= 2;
    //     }
    // }
    
    coutn(ans.size());
    forl(i , 0 , ans.size())
        cout << ans[i].ff << " " << ans[i].ss << nl;

    // int di = 2;
    // vi temp = arr;
    // vsort(temp);

    // x = temp[0];
    // bool check = true;
    // forl(i , 1 , n){
    //     rep(x < temp[i]){
    //         if (x*temp[0] > temp[i]){
    //             check = false;
    //             break;
    //         }
    //         x*= temp[0];
    //     }
    //     if (x != temp[i])
    //         check = false;

    //     if (!check)
    //         break;
    // }
    // if (check)
    //     di = temp[0];


    // vpi ans , a;
    // x = cnt(arr , 2);
    // int ind = -1 , last;
    // if (x == 0 && !check){
    //     last = -1;
    //     for(auto it : s){
    //         if (last == -1)
    //             last = it;
    //         else {
    //             x = last;
    //             rep(x < it){
    //                 if (x*last > it){
    //                     break;
    //                 }
    //                 x *= last;
    //             }

    //         }
    //     }


    // }   
}

signed main() {
    fastio;

    int t = 1;
    cin >> t;
    while (t--)
        solve();
        // coutn(solve());
        // solve() ? yes : no;
}