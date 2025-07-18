/*input
1 2 3
*/
#include<bits/stdc++.h>
typedef std::vector<int> vi;
typedef std::vector<vi> vii;
typedef std::vector<bool> vb;
typedef long long ll;
typedef std::pair<int, int> ii;
typedef std::vector<vb> vbb;
typedef std::vector<ii> vp;
int main(){
	int N, d, k;
	std::cin >> N >> d >> k;
	if(d >= N || (k == 1 && N > 2)) return !printf("NO\n");
	if(N == 1 && d == 1) return !printf("YES\n");
	vi count(N+1, 0);
	vi didatstumas(N+1);
	vp ats;
	for (int i = 1; i <= d; ++i)
	{
		ats.push_back({i, i + 1});
		count[i]++;
		count[i+1]++;
		didatstumas[i] = std::min(i - 1, d - i + 1);
		didatstumas[i+1] = std::min(i, d - i); 
	}
	int i = d+2;
	int j = 2;
	while(i <= N){
		while(j < i && (count[j] == k || didatstumas[j] == 0)){
			j++;
		}
		//printf("i = %d, j = %d\n", i, j);
		if(i == j) return !printf("NO\n");
		ats.push_back({i, j});
		count[i]++;
		count[j]++;
		didatstumas[i] = didatstumas[j] - 1;
		i++;
	}
	if(ats.size() != N-1) return !printf("NO\n");
	printf("YES\n");
	for(auto u: ats){
		printf("%d %d\n", u.first, u.second);
	}
	return 0;
}