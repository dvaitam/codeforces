#include <bits/stdc++.h>

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp> 



using namespace std;

using namespace __gnu_pbds; 



#define ff              first

#define ss              second

#define ll              int64_t

#define ld              long double

#define nl              cout<<"\n"

#define all(v)          v.begin(),v.end()

#define mset(a,v)       memset((a),(v),sizeof(a))

#define forn(i,a,b)     for(int64_t i=int64_t(a);i<int64_t(b);++i)

#define forb(i,a,b)     for(int64_t i=int64_t(a);i>=int64_t(b);--i)

#define fastio()        ios::sync_with_stdio(false); cin.tie(0); cout.tie(0);



#define mod         1'000'000'007

#define mod2        998'244'353 

#define inf         1'000'000'000'000'007

#define pi          3.14159265358979323846



template<class key,class cmp=std::less<key>>

using ordered_set=tree<key,null_type,cmp,rb_tree_tag,tree_order_statistics_node_update>;



template<class L,class R> ostream& operator<<(ostream& out,pair<L,R> &p)        {return out<<"("<<p.ff<<", "<<p.ss<<")";}

template<class T> ostream& operator<<(ostream& out,vector<T> &v)                {out<<"[";for(auto it=v.begin();it!=v.end();++it){if(it!=v.begin())out<<", ";out<<*it;}return out<<"]";}

template<class T> ostream& operator<<(ostream& out,deque<T> &v)                 {out<<"[";for(auto it=v.begin();it!=v.end();++it){if(it!=v.begin())out<<", ";out<<*it;}return out<<"]";}

template<class T> ostream& operator<<(ostream& out,set<T> &s)                   {out<<"{";for(auto it=s.begin();it!=s.end();++it){if(it!=s.begin())out<<", ";out<<*it;}return out<<"}";}

template<class L,class R> ostream& operator<<(ostream& out,map<L,R> &m)         {out<<"{";for(auto it=m.begin();it!=m.end();++it){if(it!=m.begin())out<<", ";out<<*it;}return out<<"}";}



void dbg_out() {cerr<<"]\n";}

template<typename Head,typename... Tail> 

void dbg_out(Head H,Tail... T) {cerr<<H;if(sizeof...(Tail))cerr<<", ";dbg_out(T...);}

#ifdef LOCAL

	#define dbg(...) cerr<<"["<<#__VA_ARGS__<<"] = [",dbg_out(__VA_ARGS__)

#else

	#define dbg(...)

#endif



//---------------------------------mars4---------------------------------



struct st

{

	ll p;

	vector<string> s;



	bool operator < (const st &o) const 

	{

		return p<o.p;

	}

};



vector<set<ll>> v(26);

vector<ll> color(26);

vector<ll> topsort;

bool flag=false;



void dfs(ll cur)

{

	color[cur]=1;

	for(ll i:v[cur])

	{

		if(color[i]==1)

		{

			flag=true;

		}

		else if(color[i]==0)

		{

			dfs(i);

		}

	}

	color[cur]=2;

	topsort.push_back(cur);

}



int main()

{

	fastio();

	ll z,n,m,t,k,i,j,l,d,h,r;

	cin>>n>>k;

	vector<st> a(n);

	vector<bool> mark(26);

	forn(i,0,n)

	{

		cin>>a[i].p;

		a[i].s=vector<string>(k);

		for(auto &j:a[i].s)

		{

			cin>>j;

			for(char ch:j)

			{

				mark[ch-'a']=true;

			}

		}

	}

	sort(all(a));

	string last=a[0].s[0];

	forn(i,0,n)

	{

		forn(j,0,k)

		{

			bool chk=false;

			forn(r,0,(ll)min(last.size(),a[i].s[j].size()))

			{

				if(last[r]!=a[i].s[j][r])

				{

					v[last[r]-'a'].insert(a[i].s[j][r]-'a');

					chk=true;

					break;

				}

			}

			if(!chk and last.size()>a[i].s[j].size())

			{

				flag=true;

			}

			last=a[i].s[j];

		}

	}

	forn(i,0,26)

	{

		if(mark[i] and !color[i])

		{

			dfs(i);

		}

	}

	if(flag)

	{

		cout<<"IMPOSSIBLE\n";

	}

	else

	{

		reverse(all(topsort));

		for(ll i:topsort)

		{

			cout<<char('a'+i);

		}

		nl;

	}

	cerr<<"\nTime elapsed: "<<1000*clock()/CLOCKS_PER_SEC<<"ms\n";

	return 0;

}