package storage

type CompileSource interface {
	string | QueryBuilder
}

type Compiler struct{}
type Lexer struct{}

func Compile[S CompileSource](src S) (*Program, error) {
	return nil, nil
}

func CompileStr(query string) (*Program, error) {
	return nil, nil
}

func CompileBuilder(query QueryBuilder) (*Program, error) {
	return nil, nil
}
