module scripts

go 1.16

require (
	github.com/google/uuid v1.2.0 // indirect
	internal/pkg/helper/formatrequest v1.0.0
)

replace internal/pkg/helper/formatrequest => ./internal/pkg/helper/formatrequest
replace enums/ScoreStatus => ./enums/ScoreStatus
