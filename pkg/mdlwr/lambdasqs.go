package mdlwr

type LambdaSQSMiddlewareBuilder struct {
	middlewares []MiddlewareSQSFunc
}

func CreateLambdaSQSMiddlewareBuilder() *LambdaSQSMiddlewareBuilder {
	return &LambdaSQSMiddlewareBuilder{}
}

func (r *LambdaSQSMiddlewareBuilder) AddMiddleware(m MiddlewareSQSFunc) *LambdaSQSMiddlewareBuilder {
	r.middlewares = append(r.middlewares, m)
	return r
}

func (r *LambdaSQSMiddlewareBuilder) Build(f SQSHandlerFunction) SQSHandlerFunction {
	return r.buildMiddlewareChain(f)
}

func (r *LambdaSQSMiddlewareBuilder) buildMiddlewareChain(f SQSHandlerFunction) SQSHandlerFunction {
	h := f
	m := r.middlewares
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
