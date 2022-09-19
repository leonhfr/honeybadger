package quiescence

import "context"

// None performs no quiescence search.
type None struct{}

// String implements the Interface interface.
func (None) String() string {
	return "None"
}

// Search implements the Interface interface.
func (None) Search(ctx context.Context, input Input) (*Output, error) {
	return &Output{
		Nodes: 1,
		Score: input.Evaluation.Evaluate(input.Position),
	}, nil
}
