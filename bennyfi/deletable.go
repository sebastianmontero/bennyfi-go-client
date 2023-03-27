package bennyfi

import (
	"fmt"
	"time"

	"github.com/sebastianmontero/eos-go-toolbox/util"
)

type Deletable struct {
	DeletedDate string `json:"deleted_date"`
}

func (m *Deletable) GetDeletedDate() time.Time {
	deletedDate, err := util.ToTime(m.DeletedDate)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse deletedDate: %v to asset", m.DeletedDate))
	}
	return deletedDate
}

func (m *Deletable) IsDeleted() bool {
	deletedDate := m.GetDeletedDate()
	return deletedDate.Unix() > 0
}
