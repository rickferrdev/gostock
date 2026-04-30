package helpers

type Builder[T any] struct {
	res T
}

func New[T any]() *Builder[T] {
	return &Builder[T]{}
}

func (builder *Builder[T]) Set(value T) *Builder[T] {
	builder.res = value
	return builder
}

func (builder *Builder[T]) Build() T {
	return builder.res
}
