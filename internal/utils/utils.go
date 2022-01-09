package utils

import (
	"math/rand"
	"time"
)

// GenerateID generates an ID between 0 and 500
func GenerateID() (rndID int64) {
	// we have to generate an ID, as in contrast to spring/quarkus, gorm does not provide the generate-id on ORM side
	// (as we do not want to migrate the DB to keep compatibility between the different implementations)
	// we keep this logic simple and just generate a random int64. we do not expect too much strips, so the chance of a collision should be low
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	rndID = r1.Int63n(500)
	return
}
