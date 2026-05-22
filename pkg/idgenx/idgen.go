package idgenx

import (
	"github.com/yitter/idgenerator-go/idgen"
)

func Init(workerID uint16) {
	options := idgen.NewIdGeneratorOptions(workerID)
	idgen.SetIdGenerator(options)
}
