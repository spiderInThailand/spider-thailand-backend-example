package usecase

import (
	"fmt"
	"reflect"
	"spider-go/model"

	"github.com/golang/mock/gomock"
)

// =======================================================
// GROUP EqSpiderStatistic
// =======================================================
func EqSpiderStatistic(spiderStatistics model.SpiderStatistics) gomock.Matcher {
	return eqSpiderStatistic{expectValue: spiderStatistics}
}

type eqSpiderStatistic struct {
	expectValue model.SpiderStatistics
}

func (eq eqSpiderStatistic) String() string {
	return fmt.Sprintf("matches with expectedValue %v", eq.expectValue)
}

func (eq eqSpiderStatistic) Matches(x interface{}) bool {
	actualValue, ok := x.(model.SpiderStatistics)

	if !ok {
		return false
	}

	eq.expectValue.CreatedAt = actualValue.CreatedAt
	eq.expectValue.UpdatedAt = actualValue.UpdatedAt

	isEqual := reflect.DeepEqual(eq.expectValue, actualValue)
	return isEqual

}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

// =======================================================
// GROUP EqSpiderInfo
// =======================================================
func EqSpiderInfo(spiderInfo model.SpiderInfo) gomock.Matcher {
	return eqSpiderInfo{expectValue: spiderInfo}
}

type eqSpiderInfo struct {
	expectValue model.SpiderInfo
}

func (eq eqSpiderInfo) String() string {
	return fmt.Sprintf("matches with expectedValue %v", eq.expectValue)
}

func (eq eqSpiderInfo) Matches(x interface{}) bool {
	actualValue, ok := x.(model.SpiderInfo)

	if !ok {
		return false
	}

	eq.expectValue.SpiderUUID = actualValue.SpiderUUID
	eq.expectValue.CreatedAt = actualValue.CreatedAt
	eq.expectValue.UpdatedAt = actualValue.UpdatedAt

	isEqual := reflect.DeepEqual(eq.expectValue, actualValue)

	return isEqual
}

// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
