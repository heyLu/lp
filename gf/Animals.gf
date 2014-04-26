abstract Animals = {
	flags startcat = Phrase;

	cat
		Phrase; Individual; Species; Quality;

	fun
		Is: Individual -> Quality -> Phrase;
		This, That: Species -> Individual;
		QSpecies: Quality -> Species -> Species;
		Cat, Dog, Crocodile, Camel, Rhinoceros: Species;
		Very: Quality -> Quality;
		Furry, Big, Small, Tiny, Interesting, Friendly, Fluffy, Weird, Extraordinary: Quality;
}
