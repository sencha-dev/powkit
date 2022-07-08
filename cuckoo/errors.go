package cuckoo

import (
	"fmt"
)

var (
	ErrPowTooBig      = fmt.Errorf("pow too big")
	ErrPowTooSmall    = fmt.Errorf("pow too small")
	ErrPowNotMatching = fmt.Errorf("pow not matching")
	ErrPowBranch      = fmt.Errorf("pow branch")
	ErrPowDeadEnd     = fmt.Errorf("pow dead end")
	ErrPowShortCycle  = fmt.Errorf("pow short cycle")
)
