package nef

type baseOp[F ExistsInFS] struct {
	fS   F
	calc PathCalc
	root string
}

func (m *baseOp[F]) peek(name string) (exists, isDir bool) {
	if m.fS.DirectoryExists(name) {
		return true, true
	}

	if m.fS.FileExists(name) {
		return true, false
	}

	return false, false
}
