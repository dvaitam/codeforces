/* In The Name Of God */

#include <bits/stdc++.h>



# define xx first

# define yy second

# define pb push_back

# define pp pop_back

# define eps 1e-9



using namespace std;

typedef long long ll;

typedef pair<int,int> pii;

typedef vector<int> vint;

const int maxn = 200*200;

vint v[maxn]; 

string ans;

int n,t,num[maxn],number;

void dfs(int root){

	while(!v[root].empty()){

		int tmp = v[root].back();

		v[root].pp();

		number++;

		dfs(tmp);

	}

	ans+=char(root%200);

}

int main(){

	ios_base::sync_with_stdio (0);cin.tie(0);

	cin>>n;

	for(int i=1 ; i<=n ; i++){

		string s;cin>>s;

		int u1 = s[0]*200 + s[1];

		int u2 = s[1]*200 + s[2];

		v[u1].pb(u2);

		num[u1]--;

		num[u2]++;

	}

	int root,ted;root=ted=0;

	for(int i=0 ; i<200*200 ; i++)if(num[i]<=-1){

		root = i;

		ted++;

	}

	if(num[root]<-1 || ted>=2){

		cout<<"NO\n";

		return 0;

	}

	if(num[root]==0){

		for(int i=0 ; i<200*200 ; i++)if(!v[i].empty())

			root = i;

	}

	dfs(root);

	reverse(ans.begin(),ans.end());

	ans = char(root/200) + ans;

	if(ans.size() != n+2)

		cout<<"NO\n";

	else

		cout<<"YES\n"<<ans<<endl;

	return 0;

}