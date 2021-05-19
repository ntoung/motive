package scorestatus

type ScoreStatus string

const (
	Processing ScoreStatus = "PROCESSING"
	Enqueued   ScoreStatus = "ENQUEUED"
	Done       ScoreStatus = "DONE"
)
