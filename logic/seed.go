package logic

type Logic uint8

const (
	_ Logic = iota
	LogicGlitchless
	LogicGlitched
	LogicNone
)

type Seed struct{}

func (s Seed) WinCondition()        {}
func (s Seed) Goals()               {}
func (s Seed) TokenPool()           {}
func (s Seed) ShuffableTokens()     {}
func (s Seed) LocationPool()        {}
func (s Seed) ShuffableLocations()  {}
func (s Seed) PreplacedLocations()  {}
func (s Seed) EntrancePool()        {}
func (s Seed) PreplacedEntrances()  {}
func (s Seed) ShufflableEntrances() {}
func (s Seed) EnabledTricks()       {}
func (s Seed) HintDistribution()    {}
func (s Seed) Logic() Logic         { var l Logic; return l }
