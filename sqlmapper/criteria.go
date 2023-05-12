package sqlmapper

/*
 * Collection
 */
type CriterionType string

const (
	NoValue     CriterionType = "NoValue"
	SingleValue CriterionType = "SingleValue"
	ListValue   CriterionType = "ListValue"
)

type FilterBuilder struct {
	orListOfCriteria []Criteria
	orderBy          string
	limit            int
	offset           int
}

type Filter struct {
	OrListOfCriteria []Criteria
	OrderBy          string
	Limit            int
	Offset           int
}

type Criteria struct {
	AndListOfCriterion []Criterion
}

type Criterion struct {
	Type      CriterionType
	Condition string
	Value     interface{}
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{limit: -1, offset: -1}
}

func (f *FilterBuilder) OrderBy(orderBy string) *FilterBuilder {
	f.orderBy = orderBy
	return f
}

func (f *FilterBuilder) Limit(n int) *FilterBuilder {
	f.limit = n
	return f
}

func (f *FilterBuilder) Offset(n int) *FilterBuilder {
	f.offset = n
	return f
}

func (f *FilterBuilder) Or() *FilterBuilder {
	f.orListOfCriteria = append(f.orListOfCriteria, Criteria{})
	return f
}

func (f *FilterBuilder) And(c Criterion) *FilterBuilder {
	if len(f.orListOfCriteria) == 0 {
		f.orListOfCriteria = append(f.orListOfCriteria, Criteria{})
	}

	currentCriteria := len(f.orListOfCriteria) - 1
	f.orListOfCriteria[currentCriteria].AndListOfCriterion = append(f.orListOfCriteria[currentCriteria].AndListOfCriterion, c)
	return f
}

func (f *FilterBuilder) Build() Filter {

	filter := Filter{OrderBy: f.orderBy, Limit: f.limit, Offset: f.offset}
	for _, criteria := range f.orListOfCriteria {
		if len(criteria.AndListOfCriterion) > 0 {
			filter.OrListOfCriteria = append(filter.OrListOfCriteria, criteria)
		}
	}
	return filter
}
