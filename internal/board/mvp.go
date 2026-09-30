package board

type awardMultiplier struct {
	stat       string
	multiplier float64
}

var echoArenaAwardMultipliers = []awardMultiplier{
	{stat: "stuns", multiplier: 5},
	{stat: "threepointgoals", multiplier: 50},
	{stat: "bouncegoals", multiplier: 30},
	{stat: "steals", multiplier: 20},
	{stat: "assists", multiplier: 30},
	{stat: "saves", multiplier: 30},
	{stat: "joustswon", multiplier: 40},
	{stat: "shotsongoal", multiplier: 10},
	{stat: "goals", multiplier: 20},
	{stat: "longestpossession", multiplier: 4},
	{stat: "passes", multiplier: 20},
	{stat: "catches", multiplier: 30},
	{stat: "clears", multiplier: 20},
	{stat: "interceptions", multiplier: 20},
	{stat: "onepointgoals", multiplier: 20},
	{stat: "hattricks", multiplier: 100},
}

func echoArenaMVPScore(stats map[string]float64) (float64, bool) {
	var score float64
	complete := true
	for _, award := range echoArenaAwardMultipliers {
		value, exists := stats[award.stat]
		if !exists {
			complete = false
			continue
		}
		score += value * award.multiplier
	}
	return score, complete
}
