package httpsuite

// ProblemBuilder builds ProblemDetails declaratively.
type ProblemBuilder struct {
	problem *ProblemDetails
}

// Problem starts a declarative ProblemDetails builder.
func Problem(status int) *ProblemBuilder {
	return &ProblemBuilder{
		problem: NewProblemDetails(status, "", "", ""),
	}
}

// Type sets the problem type URL.
func (b *ProblemBuilder) Type(problemType string) *ProblemBuilder {
	b.problem.Type = problemType
	return b
}

// Title sets the problem title.
func (b *ProblemBuilder) Title(title string) *ProblemBuilder {
	b.problem.Title = title
	return b
}

// Detail sets the problem detail.
func (b *ProblemBuilder) Detail(detail string) *ProblemBuilder {
	b.problem.Detail = detail
	return b
}

// Instance sets the problem instance.
func (b *ProblemBuilder) Instance(instance string) *ProblemBuilder {
	b.problem.Instance = instance
	return b
}

// Extension sets a single problem extension.
func (b *ProblemBuilder) Extension(key string, value any) *ProblemBuilder {
	if b.problem.Extensions == nil {
		b.problem.Extensions = make(map[string]interface{})
	}
	b.problem.Extensions[key] = value
	return b
}

// Extensions merges multiple problem extensions.
func (b *ProblemBuilder) Extensions(values map[string]any) *ProblemBuilder {
	for key, value := range values {
		b.Extension(key, value)
	}
	return b
}

// Build returns the configured ProblemDetails.
func (b *ProblemBuilder) Build() *ProblemDetails {
	clone := *b.problem
	if b.problem.Extensions != nil {
		clone.Extensions = make(map[string]interface{}, len(b.problem.Extensions))
		for key, value := range b.problem.Extensions {
			clone.Extensions[key] = value
		}
	}
	return &clone
}
