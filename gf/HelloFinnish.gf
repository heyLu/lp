concrete HelloFinnish of Hello = {
	lincat Greeting, Recipient = {s : Str} ;

	lin
		Hello recip = {s = "terve" ++ recip.s} ;
		World = {s = "maailma"} ;
		Mum = {s = "äiti"} ;
		Friends = {s = "ystävät"} ;
}
