package strategy

// Strategy represents an interface to the basic oprations of a model
type MainStrategy struct {

}

func (strat *MainStrategy) CalculateScore(score float32, status int, difficulty int) float32 {
	if status == 1 {
		score += strat.increaseScore(score, difficulty)
	} else {
		score -= strat.decreaseScore(score, difficulty)
	}

	return score
}

func increaseScore(score float32, difficulty int) float32 {
	var addedScore float32

	if(score < 40){
		addedScore = float32(difficulty)
	} else if score >= 40 && score < 50 {
		addedScore = float32(difficulty)*0.5
	} else {
		addedScore = float32(1)
	}

	return addedScore
}

func decreaseScore(score float32, difficulty int) float32 {
	var subScore float32

	if(score < 0){
		subScore = float32(difficulty)*2
	} else if score >= 0 && score < 10 {
		subScore = float32(difficulty)*1.5
	} else {
		subScore = float32(difficulty)
	}

	return subScore
}