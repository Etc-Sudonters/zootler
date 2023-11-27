package infra

import "sudonters/zootler/storage"

// translates zootr logic syntax to equivalent zootler dsl
type LogicRewriter interface {
	Rewrite(ZootrLogicFile) storage.Program // sure :shrug:
}
