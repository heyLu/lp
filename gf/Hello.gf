abstract Hello = {
	flags startcat = Greeting;

	cat Greeting; Recipient; -- both are categories now

	fun
		Hello: Recipient -> Greeting;
		Goodbye: Recipient -> Greeting;
		World, Mum, Friends, Cat, Dog, RobotOverlords: Recipient;
}
