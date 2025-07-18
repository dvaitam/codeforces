#include<bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

using namespace std;

using namespace __gnu_pbds;

#define TEST int tt; cin>>tt; while(tt--)

#define FOR(v,s,e) for(int v=int(s);v<int(e);v++)

#define FF first

#define SS second

#define newline cout<<"\n"; 

#define vt vector

#define SZ(v) (int)v.size() 

#define all(a) (a).begin(),(a).end()

#define rall(a) (a).rbegin(),(a).rend()

#define LOWB(a,i,j,x) ( lower_bound((a).begin()+(i), (a).begin()+(j), (x)) - (a).begin()) // a[ptr]>=key

#define UPPB(a,i,j,x) ( upper_bound((a).begin()+(i), (a).begin()+(j), (x)) - (a).begin()) // just greater than key

#define  SUM(a,i,j)   ( accumulate ((a).begin()+(i), (a).begin()+(j), 0ll))

#define UNIQ(a) (a).erase(unique((a).begin(),(a).end()),(a).end())

#define pb push_back

#define pi M_PI

// // #pragma GCC optimize ("trapv")

// -----------------------------------debugging and IO --------------------------------------

bool __DEBUGG=false;

template<typename T>

void debug(const T&single){

    if(!__DEBUGG)cout<<"{";

    cout<<single<<"}"; flush(cout);

    __DEBUGG=false;

}

template<typename T,typename...Paras>

void debug(const T& single,const Paras&...vals){

    if(!__DEBUGG)cout<<"{",__DEBUGG=true; 

    cout<<single<<", ";

    debug(vals...);

}

// ----------------------------------- cuiaoxiang ---------------------------------------------

template<class T,class U>

ostream& operator<<(ostream& outpr,pair<T,U>&prhr){outpr<<prhr.FF<<" "<<prhr.SS; return outpr;}

template<class T,class U>

istream& operator>>(istream& inpr,pair<T,U>&prhr){inpr>>prhr.FF>>prhr.SS; return inpr;}



template<class T> 

ostream& operator<<(ostream& outv,vector<T>&vec){for(auto x:vec){outv<<x<<" ";}return outv;}

template<class T>

istream& operator>>(istream& inpv,vector<T>&vec){for(auto &x:vec){inpv>>x;}return inpv;}

template<class T> 

ostream& operator<<(ostream& outv,set<T>&vec){for(auto x:vec){outv<<x<<" ";}return outv;}

//cautiously use operator overloading they are very slow, stream overloading is okay asofnow



template <class T> bool ckmin(T& a, const T& b) { return b < a ? a = b, 1 : 0; }

template <class T> bool ckmax(T& a, const T& b) { return a < b ? a = b, 1 : 0; }



template <class T> auto vect(const T& v, int n) { return vector<T>(n, v); } // base case of 1D vector

template <class T, class... D> auto vect(const T& v, int n, D... m) {       // v remains same as init value

  return vector<decltype(vect(v, m...))>(n, vect(v, m...));

}

struct NFS {

  NFS() {

    cin.tie(nullptr);   ios::sync_with_stdio(false);    cout << fixed << setprecision(12);

  };

} nfs;

template <class T> static constexpr T inf = numeric_limits<T>::max() / 2;

template<class T> 

using PQ = priority_queue<T, vector<T>, greater<T>>;

using CN = complex<double>;

template<class T>

using ordered_set = tree <T, null_type,less<T>, rb_tree_tag,tree_order_statistics_node_update>;

typedef long long int int64;

typedef std::pair<int,int> pii;

typedef std::pair<int64,int64> pii64;

const int MOD=1e9+7,MOD0=998244353, MOD32 = 2147483647;// (1<<31)-1 

const int N=1e5+1,LOG9=31,LOG6=21;

const double EPS = 1e-9;

// ------------------------------------------------------------------------------------------





int main(){

    TEST{

        int W,H;    

        cin>>H>>W;

        vt<string> grid;

        FOR(i,0,H){

            string x;   cin>>x;

            grid.pb(x);

        }

        // solve ----------------

        vt<vt<int>> ans;

        auto handleit = [&](const array<pii,2> &col1,const array<pii,2>& col2)->void{

            // empty col1 and arrange col2

            pii tl = col1[0],bl=col1[1]; // {tl,tr}      {x,y}

            pii tr = col2[0],br=col2[1]; // {bl,br}

            int cnt = (grid[tl.SS][tl.FF]=='1' ? 1 : 0) + (grid[bl.SS][bl.FF]=='1' ? 1 : 0)

                        + (grid[tr.SS][tr.FF]=='1' ? 1 : 0)

                        + (grid[br.SS][br.FF]=='1' ? 1 : 0);



            // debug(cnt,grid[tl.SS][tl.FF],grid[bl.SS][bl.FF],grid[tr.SS][tr.FF],grid[br.SS][br.FF]);newline;

            if(cnt==1){

                if(grid[tl.SS][tl.FF]=='0' && grid[bl.SS][bl.FF]=='0')return; // already Ok!

                vt<int> ansnow={tr.FF,tr.SS,br.FF,br.SS};

                if(grid[tl.SS][tl.FF]=='1')ansnow.pb(tl.FF),ansnow.pb(tl.SS),grid[tl.SS][tl.FF]='0';

                if(grid[bl.SS][bl.FF]=='1')ansnow.pb(bl.FF),ansnow.pb(bl.SS),grid[bl.SS][bl.FF]='0';



                grid[tr.SS][tr.FF]='1';

                grid[br.SS][br.FF]='1';



                ans.pb(ansnow);

            }else if(cnt==2){

                vt<int> ansnow;

                if(grid[tr.SS][tr.FF]=='1' && grid[br.SS][br.FF]=='1')return; // already ok!

                // got 2 '1's

                if(grid[tl.SS][tl.FF]=='1')ansnow.pb(tl.FF),ansnow.pb(tl.SS);//,grid[tl.SS][tl.FF]='0';

                if(grid[bl.SS][bl.FF]=='1')ansnow.pb(bl.FF),ansnow.pb(bl.SS);//,grid[bl.SS][bl.FF]='0';

                if(grid[tr.SS][tr.FF]=='1')ansnow.pb(tr.FF),ansnow.pb(tr.SS);//,grid[tr.SS][tr.FF]='0';

                if(grid[br.SS][br.FF]=='1')ansnow.pb(br.FF),ansnow.pb(br.SS);//,grid[br.SS][br.FF]='0';

                // need 1 '0'

                bool found=false;

                if(grid[tr.SS][tr.FF]=='0' && !found)ansnow.pb(tr.FF),ansnow.pb(tr.SS),found=true;//,grid[tr.SS][tr.FF]='0';

                if(grid[br.SS][br.FF]=='0' && !found)ansnow.pb(br.FF),ansnow.pb(br.SS),found=true;//,grid[br.SS][br.FF]='0';

                if(grid[tl.SS][tl.FF]=='0' && !found)ansnow.pb(tl.FF),ansnow.pb(tl.SS),found=true;//,grid[tl.SS][tl.FF]='1';

                if(grid[bl.SS][bl.FF]=='0' && !found)ansnow.pb(bl.FF),ansnow.pb(bl.SS),found=true;//,grid[bl.SS][bl.FF]='1';

                

                for(int i=0;i<SZ(ansnow);i+=2){

                    int y = ansnow[i+1],x = ansnow[i];

                    grid[y][x]= grid[y][x]=='1' ? '0' : '1';

                }

                ans.pb(ansnow);



            }else if(cnt==3){

                vt<int> ansnow;

                if(grid[tl.SS][tl.FF]=='1')ansnow.pb(tl.FF),ansnow.pb(tl.SS),grid[tl.SS][tl.FF]='0';

                if(grid[bl.SS][bl.FF]=='1')ansnow.pb(bl.FF),ansnow.pb(bl.SS),grid[bl.SS][bl.FF]='0';

                if(grid[tr.SS][tr.FF]=='1')ansnow.pb(tr.FF),ansnow.pb(tr.SS),grid[tr.SS][tr.FF]='0';

                if(grid[br.SS][br.FF]=='1')ansnow.pb(br.FF),ansnow.pb(br.SS),grid[br.SS][br.FF]='0';



                ans.pb(ansnow);

            }else if(cnt==4){

                vt<int> ansnow = {tl.FF,tl.SS,bl.FF,bl.SS,tr.FF,tr.SS};



                grid[tr.SS][tr.FF]='0';

                grid[bl.SS][bl.FF]='0';

                grid[tl.SS][tl.FF]='0';

                ans.pb(ansnow);

            }

                    

        };

        function<void(array<pii,2>,array<pii,2>)> handlefull = [&](const array<pii,2> &col1,const array<pii,2>& col2)->void{

            // empty col1 and arrange col2

            pii tl = col1[0],bl=col1[1]; // {tl,tr}      {x,y}

            pii tr = col2[0],br=col2[1]; // {bl,br}

            int cnt = (grid[tl.SS][tl.FF]=='1' ? 1 : 0) + (grid[bl.SS][bl.FF]=='1' ? 1 : 0)

                        + (grid[tr.SS][tr.FF]=='1' ? 1 : 0)

                        + (grid[br.SS][br.FF]=='1' ? 1 : 0);

            // debug("full",cnt);newline;

            if(cnt==1){

                vt<int> ansnow;

                // get 1 '1' 

                if(grid[tl.SS][tl.FF]=='1')ansnow.pb(tl.FF),ansnow.pb(tl.SS);//,grid[tl.SS][tl.FF]='0';

                if(grid[bl.SS][bl.FF]=='1')ansnow.pb(bl.FF),ansnow.pb(bl.SS);//,grid[bl.SS][bl.FF]='0';

                if(grid[br.SS][br.FF]=='1')ansnow.pb(br.FF),ansnow.pb(br.SS);//,grid[bl.SS][bl.FF]='0';

                if(grid[tr.SS][tr.FF]=='1')ansnow.pb(tr.FF),ansnow.pb(tr.SS);//,grid[bl.SS][bl.FF]='0';

                // get 2 '0's

                int cnt0=0;

                if(grid[tl.SS][tl.FF]=='0' && cnt0+1<=2)ansnow.pb(tl.FF),ansnow.pb(tl.SS),cnt0++;

                if(grid[bl.SS][bl.FF]=='0' && cnt0+1<=2)ansnow.pb(bl.FF),ansnow.pb(bl.SS),cnt0++;

                if(grid[tr.SS][tr.FF]=='0' && cnt0+1<=2)ansnow.pb(tr.FF),ansnow.pb(tr.SS),cnt0++;

                if(grid[br.SS][br.FF]=='0' && cnt0+1<=2)ansnow.pb(br.FF),ansnow.pb(br.SS),cnt0++;

                

                ans.pb(ansnow);

                for(int i=0;i<SZ(ansnow);i+=2){

                    int y = ansnow[i+1],x = ansnow[i];

                    grid[y][x]= grid[y][x]=='1' ? '0' : '1';

                }

                

                handlefull(col1,col2);

            }else if(cnt==2){

                vt<int> ansnow;

                // get 1 '1's 

                int cnt1=0;

                if(grid[tl.SS][tl.FF]=='1' && cnt1+1<=1)ansnow.pb(tl.FF),ansnow.pb(tl.SS),cnt1++;

                if(grid[bl.SS][bl.FF]=='1' && cnt1+1<=1)ansnow.pb(bl.FF),ansnow.pb(bl.SS),cnt1++;

                if(grid[br.SS][br.FF]=='1' && cnt1+1<=1)ansnow.pb(br.FF),ansnow.pb(br.SS),cnt1++;

                if(grid[tr.SS][tr.FF]=='1' && cnt1+1<=1)ansnow.pb(tr.FF),ansnow.pb(tr.SS),cnt1++;

                // get 2 '0's

                int cnt0=0;

                if(grid[tl.SS][tl.FF]=='0' && cnt0+1<=2)ansnow.pb(tl.FF),ansnow.pb(tl.SS),cnt0++;

                if(grid[bl.SS][bl.FF]=='0' && cnt0+1<=2)ansnow.pb(bl.FF),ansnow.pb(bl.SS),cnt0++;

                if(grid[tr.SS][tr.FF]=='0' && cnt0+1<=2)ansnow.pb(tr.FF),ansnow.pb(tr.SS),cnt0++;

                if(grid[br.SS][br.FF]=='0' && cnt0+1<=2)ansnow.pb(br.FF),ansnow.pb(br.SS),cnt0++;



                ans.pb(ansnow);

                

                for(int i=0;i<SZ(ansnow);i+=2){

                    int y = ansnow[i+1],x = ansnow[i];

                    grid[y][x]= grid[y][x]=='1' ? '0' : '1';

                }



                handlefull(col1,col2);



            }else if(cnt==3){

                vt<int> ansnow;

                if(grid[tl.SS][tl.FF]=='1')ansnow.pb(tl.FF),ansnow.pb(tl.SS),grid[tl.SS][tl.FF]='0';

                if(grid[bl.SS][bl.FF]=='1')ansnow.pb(bl.FF),ansnow.pb(bl.SS),grid[bl.SS][bl.FF]='0';

                if(grid[tr.SS][tr.FF]=='1')ansnow.pb(tr.FF),ansnow.pb(tr.SS),grid[tr.SS][tr.FF]='0';

                if(grid[br.SS][br.FF]=='1')ansnow.pb(br.FF),ansnow.pb(br.SS),grid[br.SS][br.FF]='0';

                ans.pb(ansnow);

                return ;

            }else if(cnt==4){

                vt<int> ansnow = {tl.FF,tl.SS,bl.FF,bl.SS,tr.FF,tr.SS};



                grid[tr.SS][tr.FF]='0';

                grid[bl.SS][bl.FF]='0';

                grid[tl.SS][tl.FF]='0';

                ans.pb(ansnow);

                handlefull(col1,col2);

            }

            

        };

        // every 2xW matrix ---- till W-1

        if(H>2 || W>2){

            FOR(j,0,H-1){

                FOR(i,0,W-2){

                    array<pii,2> col1 = {make_pair(i,j),make_pair(i,j+1)};

                    array<pii,2> col2 = {make_pair(i+1,j),make_pair(i+1,j+1)};

                    handleit(col1,col2);

                }

                j++;

            }

            if(H%2){

                int j=H-2;

                FOR(i,0,W-2){

                    array<pii,2> col1 = {make_pair(i,j),make_pair(i,j+1)};

                    array<pii,2> col2 = {make_pair(i+1,j),make_pair(i+1,j+1)};

                    handleit(col1,col2);

                }

            }

            // top to bottom ------- till H-1

            // int i=W-1;

            if(H>2){

                FOR(j,0,H-2){

                    int i=W-1;

                    array<pii,2> col1 = {make_pair(i,j),make_pair(i-1,j)};

                    array<pii,2> col2 = {make_pair(i,j+1),make_pair(i-1,j+1)};

                    handleit(col1,col2);

                }

            }

        }

        // handle last 2x2 -----------------------

        // debug("a");

        int i=W-2;

        int j=H-2;

        array<pii,2> col1 = {make_pair(i,j),make_pair(i,j+1)};

        array<pii,2> col2 = {make_pair(i+1,j),make_pair(i+1,j+1)};

        handlefull(col1,col2);

        // -------------------------------------------------------------------

        // debug("b");

        cout<<SZ(ans);newline;

        for(vt<int> &g : ans){

            for(int i=0;i<SZ(g);i+=2)cout<<g[i+1]+1<<" "<<g[i]+1<<" ";

            newline;

        }        

    }

    



    return 0;

}