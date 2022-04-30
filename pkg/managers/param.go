package managers

import "github.com/google/uuid"

func stringToUUID(s string) *uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &u
}
