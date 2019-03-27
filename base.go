package rentals

var DbModels = []interface{}{
	&User{},
	&UserSession{},
	&Apartment{},
}

type uid uint
