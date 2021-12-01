package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// RelatedFields returns solely related fields in a list of fields.
func RelatedFields(fields, related []*models.Field) []*models.Field {
	for i := range fields {
		if (fields[i].To != nil || fields[i].From != nil) &&
			fields[i].Parent.From == nil && fields[i].Parent.To == nil {
			related = append(related, fields[i])
		}
	}
	return related
}
