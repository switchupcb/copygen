package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// RemoveUnpointedFields removes unpointed fields from a Generator.
func RemoveUnpointedFields(gen *models.Generator) {
	for _, function := range gen.Functions {
		for _, fromType := range function.From {
			fromType.Field.Fields = RelatedFields(fromType.Field.Fields, nil, nil)
		}

		for _, toType := range function.To {
			toType.Field.Fields = RelatedFields(toType.Field.Fields, nil, nil)
		}
	}
}

// RelatedFields returns solely related fields in a list of fields.
func RelatedFields(fields, related []*models.Field, cyclic map[*models.Field]bool) []*models.Field {
	if cyclic == nil {
		cyclic = make(map[*models.Field]bool)
	}

	for _, subfield := range fields {
		if !cyclic[subfield] {
			cyclic[subfield] = true
			if len(subfield.Fields) != 0 {
				related = RelatedFields(subfield.Fields, related, cyclic)
			}

			if subfield.To != nil || subfield.From != nil {
				related = append(related, subfield)
			}
		}
	}

	return related
}
