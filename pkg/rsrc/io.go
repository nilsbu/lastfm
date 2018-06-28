package rsrc

// Reader is an interface for reading resources.
type Reader interface {
	Read(loc Locator) (data []byte, err error)
}

// Writer is an interface for writing resources.
type Writer interface {
	Write(data []byte, loc Locator) error
}

// Remover is an interface for removing a resources.
type Remover interface {
	Remove(loc Locator) error
}

type Updater interface {
	Update(loc Locator) (data []byte, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

type IO interface {
	ReadWriter
	Remover
}
