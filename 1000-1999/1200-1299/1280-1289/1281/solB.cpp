#include <bits/stdc++.h>
 
#define ll long long
#define pb push_back
#define pf push_front
#define pof pop_front
#define pob pop_back
#define s second
#define f first
#define mp make_pair
#define lb lower_bound
#define in insert
#define IOS ios_base::sync_with_stdio(0); cin.tie(0);cout.tie(0)
#define pi 3.1415926535
 
using namespace std;
 
ll a, b[500], n, m, k, c, d[200200], mx;
ll ans;
 
vector <int> v;
vector <int> v2;
pair <char, int> p[30000];
//map <int, bool> ma;
set <int> st;
//char cc, cc2;
//int was[10000000];

bool ch, ch2, ch3, ch4;
pair <int, int> qq[200200];
const int N = (int)2e5 + 123, inf = 1e9, mod = 1e9 + 7;
const ll INF = 1e18;
bool was[(int)2e5+5];
bool was2[(int)2e5+50];
int main(){  
	IOS;
	cin >> n;
	for(int i = 1; i <= n; i++){
		string s, s1;
	 	cin >> s >> s1;
	 	for(int j = 0; j < s.size(); j++){
	 		b[s[j] - 'A']++; 	
	 	}
	 	int cc = 0, id = -1;
	 	for(int l = 0; l < s.size(); l++){
	 	 	b[s[l] - 'A']--;
	 	 	for(int j = 0; j < 26; j++){
	 	 		//cout << b[j] << ' ' << s[l] - 'A' << ' ' << j << '\n';
	 	 	 	if(b[j] > 0 && s[l] - 'A' > j){
	 	 	 	 	id = l; //cout << "!!!\n";
	 	 	 	 	cc = j;
	 	 	 	 	break;
	 	 	 	}
	 		}
	 	 	if(id >= 0){
	 	 	 	break;
	 	 	}
	 	}
	 	if(s < s1){		
	 	 	cout << s << '\n';
	 	 	for(int j = 0; j <= 27; j++){
	 	 	 	b[j] = 0;
	 	 	}
	 	 	continue;
	 	}
	 	if(id == -1){ //cout << 1;		
	 	 	cout << "---" << '\n';  
	 	}
	 	else{
	 	 	for(int j = s.size() - 1; j > id; j--){
	 	 	 	if(s[j] == cc + 'A'){
	 	 	 	 	s[j] = s[id];
	 	 	 	 	break;
	 	 		}
	 	 	}
	 	 	s[id] = cc + 'A';
	 	 	if(s < s1){		
	 	 		cout << s << '\n';
	 	 	}
	 	 	else{		             //cout << 2;
	 	 	 	cout << "---" << '\n';
	 	 	}
	 	 	for(int j = 0; j <= 27; j++){
	 	 	 	b[j] = 0;
	 	 	}
	 	}
  	}
    return 0;
}