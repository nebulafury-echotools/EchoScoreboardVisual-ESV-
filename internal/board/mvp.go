package board

func echoArenaMVPScore(stats map[string]float64) (float64, bool) {
	points, hasPoints := stats["points"]
	assists, hasAssists := stats["assists"]
	saves, hasSaves := stats["saves"]
	score := points*5 + assists*3 + saves*2
	return score, hasPoints && hasAssists && hasSaves
}
