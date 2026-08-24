package calc

type BudgetInput struct {
	Income  float64 `json:"income"`
	Needs   float64 `json:"needs"`
	Wants   float64 `json:"wants"`
	Savings float64 `json:"savings"`
}

type BudgetOutput struct {
	TargetNeeds   float64 `json:"target_needs"`
	TargetWants   float64 `json:"target_wants"`
	TargetSavings float64 `json:"target_savings"`
	DeltaNeeds    float64 `json:"delta_needs"`
	DeltaWants    float64 `json:"delta_wants"`
	DeltaSavings  float64 `json:"delta_savings"`
	Unallocated   float64 `json:"unallocated"`
	Overspent     bool    `json:"overspent"`
	SavingsRate   float64 `json:"savings_rate"`
}

func DefaultBudget() BudgetInput {
	return BudgetInput{
		Income:  120_000,
		Needs:   60_000,
		Wants:   25_000,
		Savings: 25_000,
	}
}

func Budget(in BudgetInput) BudgetOutput {
	out := BudgetOutput{
		TargetNeeds:   in.Income * 0.50,
		TargetWants:   in.Income * 0.30,
		TargetSavings: in.Income * 0.20,
	}
	out.DeltaNeeds = in.Needs - out.TargetNeeds
	out.DeltaWants = in.Wants - out.TargetWants
	out.DeltaSavings = in.Savings - out.TargetSavings
	out.Unallocated = in.Income - in.Needs - in.Wants - in.Savings
	out.Overspent = (in.Needs + in.Wants + in.Savings) > in.Income+0.005
	if in.Income > 0 {
		out.SavingsRate = (in.Savings / in.Income) * 100
	}
	return out
}
