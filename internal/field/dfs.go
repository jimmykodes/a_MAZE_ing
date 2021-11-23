package field

import "math/rand"

// dfs is a depth-first-search maze generation method
func (f *Field) dfs() {
	var (
		count     = 0
		deferFunc func()
		animate   func()
	)
	defer func() {
		if deferFunc != nil {
			deferFunc()
		}
	}()
	for {
		f.current.Visited = true
		if f.Animate {
			if animate == nil {
				animate, deferFunc = f.animator()
			}
			animate()
		}

		if f.current.IsEnd {
			f.current = f.current.Parent
			continue
		}
		// reset count per loop
		count = f.updateAvailable()

		if count == 0 {
			if p := f.current.Parent; p != nil {
				f.current = p
				continue
			} else {
				break
			}
		}

		next := f.available[rand.Intn(count)]
		next.Parent = f.current
		f.current = next
	}
}
