#include<deque>
#include<queue>
#include<vector>
#include<algorithm>
#include<iostream>
#include<set>
#include<cmath>
#include<tuple>
#include<string>
#include<chrono>
#include<functional>
#include<iterator>
#include<random>
#include<unordered_set>
#include<array>
#include<map>
#include<iomanip>
#include<assert.h>
#include<bitset>
#include<stack>
#include<memory>
using namespace std;
typedef long long int llint;
typedef long double lldo;
#define mp make_pair
#define mt make_tuple
#define pub push_back
#define puf push_front
#define pob pop_back
#define pof pop_front
#define fir first
#define sec second
#define res resize
#define ins insert
#define era erase
/*
cout<<fixed<<setprecision(20);
cin.tie(0);
ios::sync_with_stdio(false);
*/
const llint mod=1000000007;
const llint big=2.19e15+1;
const long double pai=3.141592653589793238462643383279502884197;
const long double eps=1e-15;
template <class T,class U>bool mineq(T& a,U b){if(a>b){a=b;return true;}return false;}
template <class T,class U>bool maxeq(T& a,U b){if(a<b){a=b;return true;}return false;}
llint gcd(llint a,llint b){if(a%b==0){return b;}else return gcd(b,a%b);}
llint lcm(llint a,llint b){if(a==0){return b;}return a/gcd(a,b)*b;}
template<class T> void SO(T& ve){sort(ve.begin(),ve.end());}
template<class T> void REV(T& ve){reverse(ve.begin(),ve.end());}
template<class T>llint LBI(vector<T>&ar,T in){return lower_bound(ar.begin(),ar.end(),in)-ar.begin();}
template<class T>llint UBI(vector<T>&ar,T in){return upper_bound(ar.begin(),ar.end(),in)-ar.begin();}
int main(void){
	int n,i,j;cin>>n;
	vector<int>go(n);
	for(i=0;i<n;i++){
		int a;cin>>a;a--;
		go[a]=i;
	}
	//上->下に転置
	bool dame=0;
	for(i=0;i<n;i++){
		if(go[i]!=i){dame=1;break;}
	}
	string mto(n,'.');
	if(!dame){
		cout<<n<<endl;
		for(i=0;i<n;i++){cout<<mto<<endl;}
		return 0;
	}
	go[0]=-1;
	cout<<n-1<<endl;
	for(int h=0;h<n;h++){
		int dL=-1,dR=-1;
		for(i=0;i<n;i++){if(go[i]!=i){dL=i;break;}}
		for(i=n-1;i>=0;i--){if(go[i]!=i){dR=i;break;}}
		string gen=mto;
		
		if(dL==dR){}
		else if(go[dL]==-1){
			int t=go[dR];
			gen[dL]='/';
			gen[dR]='/';
			gen[t]='/';
			go[dL]=go[t];
			go[t]=go[dR];
			go[dR]=-1;
		}else{
			int t=go[dL];
			gen[dR]='\\';
			gen[dL]='\\';
			gen[t]='\\';
			go[dR]=go[t];
			go[t]=go[dL];
			go[dL]=-1;
		}
		cout<<gen<<endl;
	}
	return 0;
}