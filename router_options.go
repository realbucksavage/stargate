package stargate

// RouterOption represents a closure type that can be used to customize the behavior of Router created using NewRouter.
type RouterOption func(r *Router)

// WithMiddleware takes a middleware chain to be executed before all requests. The order of middleware passed to it is
// preserved.
func WithMiddleware(mw ...MiddlewareFunc) RouterOption {
	return func(r *Router) {
		r.middlewareFuncs = mw
	}
}

// WithLoadBalancer lets you set the LoadBalancerMaker of your choice.
func WithLoadBalancer(lb LoadBalancerMaker) RouterOption {
	return func(r *Router) {
		r.loadBalancerMaker = lb
	}
}
