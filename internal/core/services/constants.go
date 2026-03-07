package services

// Query limits for internal operations that need to fetch "all" records
// These are safety limits to prevent unbounded queries while being generous enough
// for typical use cases
const (
	// MaxEmployeesForMetrics is the maximum number of employees to fetch when calculating
	// department-level metrics. Assumes departments won't exceed 10,000 employees.
	MaxEmployeesForMetrics = 10000

	// MaxSubmissionsForMetrics is the maximum number of submissions to fetch when
	// calculating risk metrics for a department.
	MaxSubmissionsForMetrics = 10000

	// MaxAnswersForCalculation is the maximum number of answers to fetch when
	// calculating scores or validating submission completeness.
	MaxAnswersForCalculation = 10000

	// MaxQuestionsPerTemplate is the maximum number of questions to fetch from
	// a single assessment template. Assumes templates won't exceed 1,000 questions.
	MaxQuestionsPerTemplate = 1000

	// MaxAssignmentsPerQuestionnaire is the maximum number of department assignments
	// to fetch for a questionnaire.
	MaxAssignmentsPerQuestionnaire = 1000

	// DefaultPageSize is the default number of items to return in paginated API responses
	DefaultPageSize = 20

	// MaxPageSize is the maximum allowed page size for API responses
	MaxPageSize = 100
)
