package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// RemoveUnpointedFields removes unpointed fields from a Generator.
func RemoveUnpointedFields(gen *models.Generator) {
	for _, function := range gen.Functions {
		for _, fromType := range function.From {
			fromType.Field.Fields = RelatedFields(fromType.Field.Fields, nil)
		}

		for _, toType := range function.To {
			toType.Field.Fields = RelatedFields(toType.Field.Fields, nil)
		}
	}
}

// RelatedFields returns solely related fields in a list of fields.
func RelatedFields(fields, related []*models.Field) []*models.Field {
	for i := len(fields) - 1; i > -1; i-- {
		if len(fields[i].Fields) != 0 {
			related = append(related, RelatedFields(fields[i].Fields, related)...)
		}

		if fields[i].To != nil || fields[i].From != nil {
			related = append(related, fields[i])
		}
	}

	return related
}
