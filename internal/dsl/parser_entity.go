package dsl

func handleEntity(p *Parser) (node, error) {
	var entity EntityNode
	p.advanceToken()

	return entity, nil
}
