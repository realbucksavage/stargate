package stargate

func newHealthCounter(maxPings int) *healthCounter {
	return &healthCounter{maxPings: maxPings, alive: true}
}

type healthCounter struct {
	unhealthyPings int
	maxPings       int
	alive          bool
}

func (h *healthCounter) countHealthy() {
	if h.unhealthyPings > 0 {
		h.unhealthyPings--
	}

	if !h.alive {
		h.alive = h.unhealthyPings == 0
	}
}

func (h *healthCounter) countUnhealthy() {
	if h.unhealthyPings < h.maxPings {
		h.unhealthyPings++
	}

	if h.alive {
		h.alive = h.unhealthyPings != h.maxPings
	}
}

func (h *healthCounter) ok() bool {
	return h.alive
}
