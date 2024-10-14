package crypto

import (
	"fmt"

	"github.com/alecthomas/units"
)

type Params struct {
	T int `json:"t"`
	M int `json:"m"`
	P int `json:"p"`
	L int `json:"T"`
}

var DefaultParams = Params{
	T: 1,
	M: int(64 * units.KiB),
	P: 4,
	L: keySize,
}

func (p *Params) ToMap() map[string]int {
	return map[string]int{
		"t": p.T,
		"m": p.M,
		"p": p.P,
		"T": p.L,
	}
}

func (p *Params) String() string {
	return fmt.Sprintf("t=%d,m=%d,p=%d,T=%d", p.T, p.M, p.P, p.L)
}
