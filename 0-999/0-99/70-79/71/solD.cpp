#include<iostream>
#include<algorithm>
#include<vector>
#include<string>
#define f first
#define s second
#define mp make_pair
using namespace std;

int n, m;
string crank = "A23456789TJQK";
string suit = "CDHS";
vector<vector<bool> > inDeck(13, vector<bool>(4, true));
vector<int> crankCount(13, 4);
vector<vector<string> > card;


pair<int, int> cardData(string card) {
	if(card[0] == 'J' && card[1] >= '1' && card[1] <= '2')
		return mp(-1, card[1] - '1');
	return mp(crank.find(card[0]), suit.find(card[1]));
}



//------------------------------------------------------
//placement function

vector<vector<pair<int, int> > > jreplace(2);
vector<int> placeJokers;

bool place(int y, int x) {
	vector<bool> used(13, false);
	vector<int> jokers;
	
	for(int dy = 0; dy < 3; dy++)
		for(int dx = 0; dx < 3; dx++) {
			int cy = y + dy, cx = x + dx;
			
			auto data = cardData(card[cy][cx]);
			if(data.f == -1)
				jokers.push_back(data.s);
			else {
				if(used[data.f])
					return false;
				used[data.f] = true;
			}
			
		}
	
	placeJokers.insert(placeJokers.end(), jokers.begin(), jokers.end() );
	
	for(int i=0;i<13;i++)
		if(!used[i]) {
			for(int j=0;j<4;j++)
				if(inDeck[i][j] && jokers.size() > 0) {
					jreplace[jokers.back()].push_back(mp(i, j) );
					
					if(jokers.size() > 1) {
						jokers.pop_back();
						break;
					}
				}
		}
	return true;
} 





int main() {
	ios_base::sync_with_stdio(false);
	cin.tie(0);
	
	cin >> n >> m;
	card.resize(n, vector<string>(m));
	vector<int> jokers;
	
	for(int i=0;i<n;i++)
		for(int j=0;j<m;j++) {
			cin >> card[i][j];
			
			auto data = cardData(card[i][j]);
			if(data.f != -1) {
				inDeck[data.f][data.s] = false;
				crankCount[data.f]--;
			}
			else
				jokers.push_back(data.s);
		}
	
	
	for(int y1 = 0; y1+3 <= n; y1++)
		for(int x1 = 0; x1+3 <= m; x1++)
			for(int y2 = 0; y2+3 <= n; y2++)
				for(int x2 = 0; x2+3 <= m; x2++)
					if(y2+3 <= y1 || y1+3 <= y2 || x1+3 <= x2 || x2+3 <= x1) {
						placeJokers.clear();
						jreplace.clear();
						jreplace.resize(2);
						if(!place(y1, x1))
							continue;
						if(!place(y2, x2))
							continue;
						
						pair<int, int> card[2];
						
						if(placeJokers.size() == 1) {
							if(jreplace[placeJokers[0]].size() == 0)
								continue;
							card[placeJokers[0]] = jreplace[placeJokers[0]][0];
						}
						if(placeJokers.size() == 2) {
							if(jreplace[0].size() == 0 || jreplace[1].size() == 0)
								continue;
							if(jreplace[0].size() == 1 && jreplace[1].size() == 1 && jreplace[0][0] == jreplace[1][0])
								continue;
							card[0] = jreplace[0][0];
							card[1] = jreplace[1][0];
							if(card[0] == card[1] && jreplace[0].size() > 1)
								card[0] = jreplace[0][1];
							if(card[0] == card[1] && jreplace[1].size() > 1)
								card[1] = jreplace[1][1];
						}
						
						
						if(jokers.size() > placeJokers.size()) {
							if(placeJokers.size() == 1) {
								for(int i=0;i<13;i++)
									for(int j=0;j<4;j++)
										if(inDeck[i][j] && mp(i, j) != card[placeJokers[0]])
											card[1-placeJokers[0]] = mp(i, j);
							}
							else {
								int cnt = 0;
								
								for(int i=0;i<13;i++)
									for(int j=0;j<4;j++)
										if(inDeck[i][j]) {
											card[cnt] = mp(i, j);
											cnt = (cnt+1)%2;
										}
							}
						}
						
						
						if(jokers.size() == 1) {
							cout << "Solution exists.\n";
							cout << "Replace J" << jokers[0] + 1 << " with " << crank[card[jokers[0]].f] << suit[card[jokers[0]].s] << ".\n";
						}
						
							
						else if(jokers.size() == 2) {	
							cout << "Solution exists.\n";
							cout << "Replace J1 with " << crank[card[0].f] << suit[card[0].s];
							cout << " and J2 with " << crank[card[1].f] << suit[card[1].s] << ".\n";
						}
						else
							cout << "Solution exists.\nThere are no jokers.\n";
						
						cout << "Put the first square to (" << y1+1 << ", " << x1+1 << ").\n";
						cout << "Put the second square to (" << y2+1 << ", " << x2+1 << ").\n";
						
						return 0;
					}
	
	
	cout << "No solution.\n";
	return 0;
}