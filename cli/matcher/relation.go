package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// RelatedFields returns solely related fields in a list of fields.
func RelatedFields(fields, related []*models.Field) []*models.Field {
	for i := range fields {
		if len(fields[i].Fields) != 0 {
			related = RelatedFields(fields[i].Fields, related)
		}

		if fields[i].To != nil || fields[i].From != nil {
			related = append(related, fields[i])
		}
	}

	return related
}
