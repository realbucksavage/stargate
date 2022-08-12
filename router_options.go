package stargate

type RouterOption func(r *Router)

func WithMiddleware(mw ...MiddlewareFunc) RouterOption {
	return func(r *Router) {
		r.middlewareFuncs = mw
	}
}

func WithLoadBalancer(lb LoadBalancerMaker) RouterOption {
	return func(r *Router) {
		r.loadBalancerMaker = lb
	}
}
