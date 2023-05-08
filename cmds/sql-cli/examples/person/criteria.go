package person

import (
	"database/sql"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-batis/sqlmapper"
)

/*
 * Criteria
 */
type FilterBuilder struct {
	fb *sqlmapper.FilterBuilder
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{fb: &sqlmapper.FilterBuilder{}}
}

func (ub *FilterBuilder) OrderBy(ob string) *FilterBuilder {
	ub.fb.OrderBy(ob)
	return ub
}

func (ub *FilterBuilder) Or() *FilterBuilder {
	ub.fb.Or()
	return ub
}

func (ub *FilterBuilder) Build() sqlmapper.Filter {
	return ub.fb.Build()
}

func (ub *FilterBuilder) AndIdEqualTo(aId Max20Text) *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.SingleValue, Condition: "id = ", Value: aId})
	return ub
}

func (ub *FilterBuilder) AndLastnameEqualTo(aLastname Max20Text) *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.SingleValue, Condition: "lastname = ", Value: aLastname})
	return ub
}

func (ub *FilterBuilder) AndNicknameEqualTo(aNickname sql.NullString) *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.SingleValue, Condition: "nickname = ", Value: aNickname})
	return ub
}

func (ub *FilterBuilder) AndAgeEqualTo(aAge sql.NullInt32) *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.SingleValue, Condition: "age = ", Value: aAge})
	return ub
}

func (ub *FilterBuilder) AndAgeIsNull() *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.NoValue, Condition: "age is null "})
	return ub
}

func (ub *FilterBuilder) AndConsensusEqualTo(aConsensus sql.NullBool) *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.SingleValue, Condition: "consensus = ", Value: aConsensus})
	return ub
}

func (ub *FilterBuilder) AndConsensusIsNull() *FilterBuilder {
	ub.fb.And(sqlmapper.Criterion{Type: sqlmapper.NoValue, Condition: "consensus is null "})
	return ub
}
