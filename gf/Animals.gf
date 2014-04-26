abstract Animals = {
	flags startcat = Phrase;

	cat
		FullPhrase; Phrase; Individual; Species; Quality;

	fun
		ExcuseMeBut: Phrase -> FullPhrase;
		Q: Individual -> Quality -> Phrase;
		Is: Individual -> Quality -> Phrase;
		This, That: Species -> Individual;
		QSpecies: Quality -> Species -> Species;
		Cat, Dog, Crocodile, Camel, Rhinoceros, Pony: Species;
		Very: Quality -> Quality;
		Furry, Big, Small, Tiny, Interesting, Friendly, Fluffy, Weird, Extraordinary, Cute, Adorable: Quality;
}
