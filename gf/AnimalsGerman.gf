concrete AnimalsGerman of Animals = {
	lincat
		Phrase, Individual, Species, Quality = {s: Str};

	lin
		Is individual quality = {s = individual.s ++ "ist" ++ quality.s};

		This species = {s = "dieses" ++ species.s};
		That species = {s = "das" ++ species.s};

		QSpecies quality species = {s = quality.s ++ species.s};

		Cat = {s = "Katze"};
		Dog = {s = "Hund"};
		Crocodile = {s = "Krokodil"};
		Camel = {s = "Kamel"};
		Rhinoceros = {s = "Rhinozeros"};
		Pony = {s = "Pony"};

		Very quality = {s = "sehr" ++ quality.s};

		Furry = {s = "haarig"};
		Big = {s = "big"};
		Small = {s = "klein"};
		Tiny = {s = "tiny"};
		Interesting = {s = "interessant"};
		Friendly = {s = "freundlich"};
		Fluffy = {s = "flauschig"};
		Weird = {s = "seltsam"};
		Extraordinary = {s = "außerordentlich"};
		Cute = {s = "süß"};
}
