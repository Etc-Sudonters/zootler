package infra

type ZootrLogicFile struct {
	content []byte
}

type LogicFileEntry struct{}

type LogicFileIterator struct {
	file ZootrLogicFile
	pos  int
}

func (l *LogicFileIterator) Advance() bool {
	return false
}

func (l *LogicFileIterator) Current() (*LogicFileEntry, error) {
	return nil, nil
}
