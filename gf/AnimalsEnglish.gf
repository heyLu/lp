concrete AnimalsEnglish of Animals = {
	lincat
		FullPhrase, Phrase, Individual, Species, Quality = {s: Str};

	lin
		ExcuseMeBut phrase = {s = "excuse me but" ++ phrase.s};

		Q individual quality = {s = "is" ++ individual.s ++ quality.s ++ "?"};

		Is individual quality = {s = individual.s ++ "is" ++ quality.s};

		This species = {s = "this" ++ species.s};
		That species = {s = "that" ++ species.s};

		QSpecies quality species = {s = quality.s ++ species.s};

		Cat = {s = "cat"};
		Dog = {s = "dog"};
		Crocodile = {s = "crocodile"};
		Camel = {s = "camel"};
		Rhinoceros = {s = "rhinoceros"};
		Pony = {s = "pony"};

		Very quality = {s = "very" ++ quality.s};

		Furry = {s = "furry"};
		Big = {s = "big"};
		Small = {s = "small"};
		Tiny = {s = "tiny"};
		Interesting = {s = "interesting"};
		Friendly = {s = "friendly"};
		Fluffy = {s = "fluffy"};
		Weird = {s = "weird"};
		Extraordinary = {s = "extraordinary"};
		Cute = {s = "cute"};
		Adorable = {s = "adorable"};
}
