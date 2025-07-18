#include <bits/stdc++.h>

     

    #define ull unsigned long long

     

    #define int long long

     

    #define ld long double

     

    #define mod 1000000007

     

    #define INF 1000000000

     

    #define rep(i,x,y) for(int i = x; i<y; ++i)

    

    #define repq(i,x,y) for(int i = x; i<=y; ++i)

    

    #define rev(i,x,y) for(int i = x ; i >= y; --i)

     

    #define pb push_back 

    

    #define pii pair<int, int> 

    

    #define fi first

    

    #define se second

    

    using namespace std;

     

     

    /////////////////////////////////////////Fast (I/O)////////////////////////////////////

     

    inline int readInt(){

    	int sign = 1, x = 0, c = getc(stdin);

    	while (c <= 32){

    		c = getc(stdin);

    	}

    	if (c == '-'){

    		sign = -1;

    		c = getc(stdin);

    	}

    	while ('0' <= c && c <= '9' && c != EOF){

    		x = (x<<1) + (x<<3) + c - '0'; // x*10 == (x << 1) + (x << 3)

    		c = getc(stdin);

    	}

    	return x*sign;

    }

     

    inline void writeLong(long long x){

    	if (x < 0){

    		putc('-', stdout);

    		x = -x;

    	}

    	char s[40];

    	int n = 0;

    	while (x || !n){ // if x == 0 it won't be printed! (added !n) :)

    		s[n++] = '0' + x % 10;

    		x /= 10;

    	}

    	while (n--){

    		putc(s[n], stdout);

    	}

    }

    ///////////////////////////////Auxiliary functions////////////////////////////

     

    inline bool des (int i , int j) { return i > j; }

    inline bool asc (int i , int j) { return i < j; }

    bool pairComp_s (const pii &x , const pii &y) { return asc(x.se, y.se); }

    

    int fact(int x) { int ans = 1; while(x) ans*= (x--); return ans; }

    int gcd(int x, int y) { return y ? gcd(y, x % y) : x ;}

    int lcm(int x, int y) { return x * 1ll * y / gcd(x, y); }

    int ncr(int n, int r) { return fact(n) / (fact(n-r) * fact(r)); }

    

    ////////////////////////////////////Functions////////////////////////////////////

    

    int bin_pow(int a, int b){

        

        int res = 1ull;

        int tmp = a;

        

        while(b){

            if(b & 1) res = res * tmp % mod;

            tmp = tmp * tmp % mod;

            b>>=1;

        }

        

        return res;

        

    }

    

 

    

    ////////////////////////////////////Variables/////////////////////////////////

     

    int hexagonal_spiral[][2] = {{1,2} , {-1, 2} , {-2, 0} , {-1, -2} , {1, -2} , {2, 0}};

    int rep[] = {1, 0, 1, 1, 1, 1}; 



    const int M = 2e6 + 1;

    int n, m, k;

    int a[500005];

    //int a[M], b[M], c[300];

    //set<int> myset;

    //map<int, int> cnt;

    //map<int, int>:: iterator mit;

    //map<int, int>:: reverse_iterator rit;

    //vector<int> v;

    //vector<pii> v;

    //priority_queue<int> q;

    int sum, MAX, MIN;

    int st[3];

    int x;

    

    

    bool test(int x){

        

        for(int i = 0, b = 0; i < 3; i++){

            b = upper_bound(a, a+n, a[st[i] = b] + x) - a;

            if(b == n) return false;

        }

        

        return true;

    }

    /////////////////////////////////Main Code//////////////////////////////////////

    

    int32_t main(){

        

       n = readInt();

       

       rep(i, 0, n) a[i] = readInt();

       

       sort(a, a+n);

       

       int l = 0, r = INF, mid;

       while(l < r){

           

           mid = l + (r - l) / 2;

           

           if(test(mid)) l = mid + 1;

           

           else r = mid;

       }

       

       test(r);

       //cout << l << " " << mid << " " << r <<endl;

       

       cout << setprecision(6) << fixed;

       cout << l/2.0 << endl;

       

       rep(i, 0, 3) cout << a[st[i]] + l/2.0 << " ";

       

       cout << endl;



       

        

    return 0;

     

    }