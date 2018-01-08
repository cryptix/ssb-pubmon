package hexagen

import (
	"log"
)

func next(current, prev FaceAddr, input bool) (out FaceAddr) {
	type caseStruct struct {
		a bool
		b bool
		Orientation
	}

	if !adjacent(current, prev) {
		log.Println("error: current and prev not adjacent")
	}

	defer func() {
		if !adjacent(current, out) {
			log.Println("error: next face not adjacent to current one!", current, out)
		}
	}()

	switch (caseStruct{
		current.X == prev.X,
		current.Y == prev.Y,
		current.Orientation,
	}) {
	case caseStruct{true, true, Up}: // both X and Y equal
		if input {
			return FaceAddr{X: current.X, Y: current.Y - 1, Orientation: Down}
		} else {
			return FaceAddr{X: current.X - 1, Y: current.Y, Orientation: Down}
		}

	case caseStruct{true, false, Up}: // X equal, Y not
		if input {
			return FaceAddr{X: current.X - 1, Y: current.Y, Orientation: Down}
		} else {
			return FaceAddr{X: current.X, Y: current.Y, Orientation: Down}
		}

	case caseStruct{false, true, Up}: // X equal, Y not
		if input {
			return FaceAddr{X: current.X, Y: current.Y, Orientation: Down}
		} else {
			return FaceAddr{X: current.X, Y: current.Y - 1, Orientation: Down}
		}

	case caseStruct{true, true, Down}: // both X and Y equal
		if input {
			return FaceAddr{X: current.X, Y: current.Y + 1, Orientation: Up}
		} else {
			return FaceAddr{X: current.X + 1, Y: current.Y, Orientation: Up}
		}

	case caseStruct{true, false, Down}: // X equal, Y not
		if input {
			return FaceAddr{X: current.X + 1, Y: current.Y, Orientation: Up}
		} else {
			return FaceAddr{X: current.X, Y: current.Y, Orientation: Up}
		}

	case caseStruct{false, true, Down}: // X equal, Y not
		if input {
			return FaceAddr{X: current.X, Y: current.Y, Orientation: Up}
		} else {
			return FaceAddr{X: current.X, Y: current.Y + 1, Orientation: Up}
		}
	}
	log.Println(current, prev, input)
	panic("error: no correct branch!")
}

// wrap at hex borders. not clean. FIXME
func wrap(addr FaceAddr) FaceAddr {
	if inhexagon(addr) {
		return addr
	}

	// top wrap
	if (addr == FaceAddr{0, 4, Up}) {
		return FaceAddr{2, 0, Up}
	}

	if (addr == FaceAddr{1, 4, Up}) {
		return FaceAddr{3, 0, Up}
	}

	// bottom wrap
	if (addr == FaceAddr{2, -1, Down}) {
		return FaceAddr{0, 3, Down}
	}

	if (addr == FaceAddr{3, -1, Down}) {
		return FaceAddr{1, 3, Down}
	}

	//bottom-left wrap
	if (addr == FaceAddr{0, 1, Up}) {
		return FaceAddr{3, 2, Up}
	}

	if (addr == FaceAddr{1, 0, Up}) {
		return FaceAddr{2, 3, Up}
	}

	// top-right wrap
	if (addr == FaceAddr{2, 3, Down}) {
		return FaceAddr{0, 1, Down}
	}

	if (addr == FaceAddr{3, 2, Down}) {
		return FaceAddr{1, 0, Down}
	}

	// top-left wrap
	if (addr == FaceAddr{-1, 3, Down}) {
		return FaceAddr{3, 1, Down}
	}

	if (addr == FaceAddr{-1, 2, Down}) {
		return FaceAddr{3, 1, Down}
	}

	// bottom-right wrap
	if (addr == FaceAddr{4, 0, Up}) {
		return FaceAddr{0, 2, Up}
	}

	if (addr == FaceAddr{4, 1, Up}) {
		return FaceAddr{0, 3, Up}
	}

	log.Println(addr)
	panic("error: wrapping not properly implmented")
}
