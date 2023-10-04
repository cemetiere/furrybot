package balance

type BalanceService struct {
	userBalance map[int64]int64
}

func CreateNewBalanceService() *BalanceService {
	return &BalanceService{make(map[int64]int64)}
}

func (balance *BalanceService) createUserBalance(userId int64) {
	balance.userBalance[userId] = 0
}

// Returns a user's balance, creates the users balance with 0 money
// if the user is not found.
func (balance *BalanceService) GetBalance(userId int64) int64 {
	val, ok := balance.userBalance[userId]

	if !ok {
		balance.createUserBalance(userId)
	}

	return val
}

func (balance *BalanceService) IncreaseBalance(userId, amount int64) {
	val := balance.GetBalance(userId)

	balance.userBalance[userId] = val + amount
}

func (balance *BalanceService) DecreaseBalance(userId, amount int64) bool {
	val := balance.GetBalance(userId)

	if val < amount {
		return false
	}

	balance.userBalance[userId] = val - amount

	return true
}
