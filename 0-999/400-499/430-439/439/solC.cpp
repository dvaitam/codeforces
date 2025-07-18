/* Author haleyk10198 */

/* �@��:  haleyk10198 */

#include <iostream>

#include <fstream>

#include <sstream>

#include <cstdlib>

#include <cstdio>

#include <vector>

#include <map>

#include <queue>

#include <cmath>

#include <algorithm>

#include <cstring>

#include <iomanip>

#include <ctime>

#include <string>

#include <set>



#define MOD 1000000007

#define INF 2147483647

#define PI 3.1415926535897932384626433

#define ll long long

#define pii pair<int,int>

#define mp(x,y) make_pair((x),(y))



using namespace std;



vector<int> odd,even;



int main(){

	int n,k,p;

	scanf("%d%d%d",&n,&k,&p);

	for(int i=0;i<n;i++){

		int val;

		scanf("%d",&val);

		if(val&1)

			odd.push_back(val);

		else

			even.push_back(val);

	}

	if(odd.size()<(k-p)||(odd.size()-(k-p))&1||(odd.size()-(k-p))/2+even.size()<p)

		printf("NO\n");

	else{

		printf("YES\n");

		if(p>0){

			int i,j;

			for(i=0;i<k-p;i++)

				printf("1 %d\n",odd[i]);

			copy(odd.begin()+(k-p),odd.end(),back_inserter(even));

			for(i=j=0;j<p-1;i++,j++){

				if(even[i]&1){

					printf("2 %d %d\n",even[i],even[i+1]);

					i++;

				}

				else

					printf("1 %d\n",even[i]);

			}

			printf("%d ",even.size()-i);

			for(;i<even.size();i++)

				printf("%d%c",even[i],i==even.size()-1?'\n':' ');

		}

		else{

			copy(even.begin(),even.end(),back_inserter(odd));

			for(int i=0;i<k-p-1;i++)

				printf("1 %d\n",odd[i]);

			printf("%d ",odd.size()-(k-p-1));

			for(int i=k-p-1;i<odd.size();i++)

				printf("%d%c",odd[i],i==odd.size()-1?'\n':' ');

		}

	}

	return 0;

}