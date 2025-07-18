#include<bits/stdc++.h>

#include<queue>

#include<stack>

#define INF 0x3f3f3f3f

typedef long long ll;

using namespace std;

bool tag[65],ans[2005];

int main(){

	int t;

	scanf("%d",&t);

	while(t--){

		ll k;

		memset(ans,0,sizeof(ans));

		scanf("%lld",&k);

		if(k&1ll){

			printf("-1\n");

			continue;

		}

		for(int i=0;i<62;i++){

			if(k>>i&1)

				tag[i]=true;

			else

				tag[i]=false;

		}

		int ptr=0;

		ans[ptr]=true;

		if(tag[1]){

			ans[++ptr]=true;

		}

		for(int i=2;i<62;i++){

			if(tag[i]){

				ans[++ptr]=true;

				ptr+=i-1;

				ans[ptr]=true;

			}

		}

		printf("%d\n",ptr);

		for(int i=ptr;i>=1;i--){

			if(ans[i])

				printf("1 ");

			else

				printf("0 ");

		}

		printf("\n");

	}

	return 0;//

}