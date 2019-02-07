package n3grpc

import (
	"fmt"

	u "github.com/cdutwhu/go-util"
)

var (
	fPln = fmt.Println
	fPf  = fmt.Printf
	fSf  = fmt.Sprintf

	PE   = u.PanicOnError
	PE1  = u.PanicOnError1
	PH   = u.PanicHandle
	PC   = u.PanicOnCondition
	Must = u.Must
)
