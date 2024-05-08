package mdlwr

type LambdaRestMiddlewareBuilder struct {
	middlewares []MiddlewareFunc
}

func CreateLambdaRESTMiddlewareBuilder() *LambdaRestMiddlewareBuilder {
	return &LambdaRestMiddlewareBuilder{}
}

func (r *LambdaRestMiddlewareBuilder) AddMiddleware(m MiddlewareFunc) *LambdaRestMiddlewareBuilder {
	r.middlewares = append(r.middlewares, m)
	return r
}

func (r *LambdaRestMiddlewareBuilder) Build(f HandlerFunc) HandlerFunc {
	return r.buildMiddlewareChain(f)
}

func (r *LambdaRestMiddlewareBuilder) buildMiddlewareChain(f HandlerFunc) HandlerFunc {
	h := f
	m := r.middlewares
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
