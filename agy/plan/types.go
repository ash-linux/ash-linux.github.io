package plan

type Status string

const (
	StatusDone  Status = "done"
	StatusWIP   Status = "wip"
	StatusPlanned Status = "planned"
)

type Phase struct {
	Name     string
	Deliverables []Deliverable
}

type Deliverable struct {
	Text   string
	Status Status
}

type PriorityItem struct {
	Priority string // P0, P1, P2
	Phase    string
	Feature  string
	Effort   string
	Impact   string
}

type Feature struct {
	Number int
	Title  string
	Theme  string
	Status Status
	Desc   string
	Phase  string
}

type Plan struct {
	Phases         []Phase
	PriorityItems  []PriorityItem
	Features       []Feature
	FeatureThemes  []string
}
