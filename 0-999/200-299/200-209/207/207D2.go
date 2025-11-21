package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	matchPrefix = iota
	matchExact
)

type keyword struct {
	topic   int
	pattern string
	weight  int
	mode    int
}

type phraseBoost struct {
	phrase string
	topic  int
	weight int
}

var keywords = func() []keyword {
	add := func(topic, weight, mode int, patterns ...string) []keyword {
		items := make([]keyword, 0, len(patterns))
		for _, p := range patterns {
			items = append(items, keyword{topic: topic, pattern: p, weight: weight, mode: mode})
		}
		return items
	}

	var list []keyword

	list = append(list, add(1, 6, matchPrefix,
		"cultur", "museum", "literat", "theatr", "concert", "festival", "opera", "ballet",
		"gallery", "sculpt", "orchestra", "folklor",
	)...)
	list = append(list, add(1, 5, matchPrefix,
		"music", "film", "cinema", "screen", "author", "novel", "writer", "poet", "drama",
		"director", "artist", "artwork", "design", "fashion", "exhibit",
	)...)
	list = append(list, add(1, 4, matchPrefix,
		"album", "song", "lyric", "choreograph", "choir", "broadway", "museum", "historic",
	)...)
	list = append(list, add(1, 6, matchPrefix,
		"культур", "искусств", "музе", "худож", "музык", "пев", "фильм", "кино", "карти",
		"роман", "литератур", "театр", "спектак", "балет", "концерт", "ансамбл", "альбом",
		"пьес", "режиссер", "режиссёр", "премь", "галер", "композитор", "артист", "исполн",
		"сценар", "книг", "писател", "пластинк", "экспозиц", "живопис", "скульпт",
	)...)

	list = append(list, add(2, 6, matchPrefix,
		"polit", "government", "ministr", "president", "parliament", "congress", "senate",
		"diplomat", "governor", "cabinet", "chancellor", "security", "coalition", "campaign",
		"election", "referend", "policy", "federal", "council", "legisl", "state", "govern",
	)...)
	list = append(list, add(2, 5, matchPrefix,
		"duma", "kremlin", "embassy", "mayor", "municip", "opposit", "minister", "supreme",
		"constitution", "dissident", "military", "army", "conflict", "border", "alliance",
	)...)
	list = append(list, add(2, 6, matchPrefix,
		"правитель", "министр", "президент", "премьер", "парламент", "депутат", "дум",
		"кремл", "госдум", "выбор", "голос", "коалиц", "оппозиц", "закон", "конституц",
		"совет", "администрац", "губернат", "силов", "арм", "военн", "безопасн", "правоохр",
		"мэр", "суд", "прокурат", "посол", "дипломат", "комитет", "реформ",
	)...)

	list = append(list, add(3, 6, matchPrefix,
		"econom", "market", "trade", "business", "company", "corporate", "industr", "product",
		"invest", "capital", "financ", "budget", "profit", "revenue", "credit", "bank",
		"currency", "export", "import", "stock", "share", "bond", "retail", "supply",
		"demand", "inflation", "logistic", "freight", "harvest", "agricultur", "contract",
	)...)
	list = append(list, add(3, 5, matchPrefix,
		"oil", "gas", "energy", "metal", "steel", "coal", "mining", "shipping", "insur",
		"tariff", "custom", "venture", "startup", "pricing", "valuation", "merger",
	)...)
	list = append(list, add(3, 6, matchPrefix,
		"эконом", "рынк", "компан", "фирм", "бизн", "промыш", "завод", "производ", "товар",
		"прибыл", "убыт", "банк", "кредит", "инвест", "капитал", "финанс", "бюджет", "деньг",
		"доллар", "рубл", "евро", "валют", "цена", "стоим", "продаж", "экспорт", "импорт",
		"нефт", "газ", "энерг", "металл", "сталь", "угл", "сельхоз", "урож", "налог", "доход",
		"выруч", "дивиден", "контракт", "сделк", "предприним", "бирж", "акци", "облига",
		"инфляц", "ипотек", "поставк", "портфель",
	)...)

	list = append(list, add(1, 3, matchExact,
		"art", "arts", "song", "songs", "poem", "poems", "aria", "folk",
	)...)
	list = append(list, add(2, 3, matchExact,
		"law", "laws", "army", "war", "nato", "un", "uno", "osce",
	)...)
	list = append(list, add(3, 3, matchExact,
		"oil", "gas", "loan", "loans", "gdp", "vat", "ipo", "opec",
	)...)

	return list
}()

var phraseBoosts = []phraseBoost{
	{"state duma", 2, 7},
	{"white house", 2, 5},
	{"prime minister", 2, 6},
	{"presidential election", 2, 7},
	{"foreign ministry", 2, 6},
	{"security council", 2, 4},
	{"supreme court", 2, 4},
	{"united nations", 2, 4},
	{"human rights", 2, 3},
	{"state department", 2, 4},
	{"министерство культуры", 1, 4},
	{"министерство иностранных", 2, 5},
	{"государственная дума", 2, 7},
	{"верховный суд", 2, 4},
	{"центральный банк", 3, 6},
	{"центробанк", 3, 5},
	{"мировой банк", 3, 5},
	{"фондовый рынок", 3, 6},
	{"фондовая биржа", 3, 6},
	{"станция метро", 1, 2},
	{"фильм фестив", 1, 6},
	{"film festival", 1, 6},
	{"art exhibition", 1, 5},
	{"book fair", 1, 4},
	{"literary award", 1, 5},
	{"music award", 1, 4},
	{"stock exchange", 3, 6},
	{"central bank", 3, 6},
	{"trade balance", 3, 4},
	{"gross domestic product", 3, 6},
	{"oil company", 3, 5},
	{"gas field", 3, 4},
	{"oil price", 3, 4},
}

var yoReplacer = strings.NewReplacer("ё", "е")

func main() {
	reader := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	raw := string(data)
	normalized := normalize(raw)
	tokens := tokenize(normalized)

	var scores [4]int
	for _, token := range tokens {
		if token == "" {
			continue
		}
		applyKeywordScore(token, &scores)
		applyExactTokenHeuristics(token, &scores)
		if isNumericToken(token) {
			if len(token) >= 4 {
				scores[3] += 3
			} else {
				scores[3] += 1
			}
		}
	}

	applyPhraseHeuristics(normalized, &scores)
	applyCharacterHeuristics(normalized, len(tokens), &scores)

	// generic fallbacks based on document structure
	if scores[1] == 0 && strings.Contains(normalized, "art") {
		scores[1]++
	}
	if scores[2] == 0 && strings.Contains(normalized, "minister") {
		scores[2]++
	}
	if scores[3] == 0 && strings.Contains(normalized, "market") {
		scores[3]++
	}

	result := 1
	for topic := 2; topic <= 3; topic++ {
		if scores[topic] > scores[result] {
			result = topic
		} else if scores[topic] == scores[result] && topic > result {
			result = topic
		}
	}

	fmt.Println(result)
}

func normalize(s string) string {
	lower := strings.ToLower(s)
	return yoReplacer.Replace(lower)
}

func tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
		return true
	})
}

func applyKeywordScore(token string, scores *[4]int) {
	for _, kw := range keywords {
		if kw.mode == matchExact {
			if token == kw.pattern {
				scores[kw.topic] += kw.weight
			}
			continue
		}
		if strings.HasPrefix(token, kw.pattern) {
			scores[kw.topic] += kw.weight
		}
	}
}

func applyExactTokenHeuristics(token string, scores *[4]int) {
	switch token {
	case "art", "arts", "artist", "artists":
		scores[1] += 3
	case "song", "songs", "album", "albums":
		scores[1] += 3
	case "novel", "novels", "book", "books":
		scores[1] += 2
	case "poet", "poetry", "poems":
		scores[1] += 3
	case "opera", "ballet":
		scores[1] += 4
	case "award", "awards", "oscar", "grammy", "emmy", "bafta":
		scores[1] += 2
	case "war", "wars", "army", "brigade":
		scores[2] += 3
	case "law", "laws", "court", "courts":
		scores[2] += 2
	case "vote", "votes", "voter", "voters":
		scores[2] += 3
	case "referendum":
		scores[2] += 4
	case "nato", "un", "uno", "osce", "eu", "european":
		scores[2] += 3
	case "russia", "kremlin", "moscow", "duma":
		scores[2] += 3
	case "usa", "america", "washington":
		scores[2] += 2
	case "president", "premier", "minister", "governor":
		scores[2] += 2
	case "oil", "gas", "diesel":
		scores[3] += 3
	case "bank", "banks":
		scores[3] += 2
	case "loan", "loans", "credit", "credits":
		scores[3] += 2
	case "percent", "percentage", "percentages":
		scores[3] += 2
	case "million", "billion", "trillion":
		scores[3] += 3
	case "gdp", "ipo", "opec", "wto", "imf":
		scores[3] += 4
	case "eur", "usd", "rur", "rub", "uah", "yen", "yuan":
		scores[3] += 3
	case "stock", "stocks", "share", "shares":
		scores[3] += 2
	case "company", "companies", "business":
		scores[3] += 2
	case "budget":
		scores[3] += 2
	case "налог", "налоги":
		scores[3] += 3
	case "выборы", "выбор":
		scores[2] += 4
	case "министр", "министер":
		scores[2] += 3
	case "губернатор":
		scores[2] += 3
	case "госдума", "дума":
		scores[2] += 4
	case "суд":
		scores[2] += 2
	case "армия", "армии":
		scores[2] += 3
	case "нефть", "газ":
		scores[3] += 4
	case "рубль", "рублей", "руб":
		scores[3] += 3
	case "доллар", "долларов":
		scores[3] += 3
	case "евро":
		scores[3] += 3
	case "биржа", "биржи":
		scores[3] += 4
	case "банк", "банка", "банки":
		scores[3] += 3
	case "культура", "искусство":
		scores[1] += 4
	case "театр", "театра":
		scores[1] += 3
	case "музыка", "песня", "песни":
		scores[1] += 3
	case "роман":
		scores[1] += 3
	case "поэт", "поэты":
		scores[1] += 3
	}
}

func isNumericToken(token string) bool {
	if token == "" {
		return false
	}
	hasDigit := false
	for _, r := range token {
		if r >= '0' && r <= '9' {
			hasDigit = true
			continue
		}
		if r == ',' || r == '.' {
			continue
		}
		return false
	}
	return hasDigit
}

func applyPhraseHeuristics(text string, scores *[4]int) {
	for _, ph := range phraseBoosts {
		if strings.Contains(text, ph.phrase) {
			scores[ph.topic] += ph.weight
		}
	}
}

func applyCharacterHeuristics(text string, tokenCount int, scores *[4]int) {
	digits := 0
	for _, ch := range text {
		if ch >= '0' && ch <= '9' {
			digits++
		}
	}
	if digits > tokenCount {
		scores[3] += 4
	} else if digits*2 > tokenCount {
		scores[3] += 2
	}

	percentCount := strings.Count(text, "%")
	if percentCount > 0 {
		scores[3] += percentCount * 2
	}

	currencySymbols := strings.Count(text, "$") + strings.Count(text, "€") + strings.Count(text, "£") + strings.Count(text, "¥")
	if currencySymbols > 0 {
		scores[3] += currencySymbols * 3
	}

	if strings.Contains(text, "state of emergency") || strings.Contains(text, "martial law") {
		scores[2] += 4
	}
	if strings.Contains(text, "cultural heritage") {
		scores[1] += 4
	}
}
