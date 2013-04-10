package understanding;

import org.junit.Test;

import static org.junit.Assert.*;

/**
 * Documenting the behaviour of output parameters in Java.
 */
public class TestOutputParameters {
	// FIXME: The name already suggests that this should go in either a
	//        Strings class or in a separate ImmutableObjects class.
	//        However, because immutable objects probably behave
	//        differently than mutable objects especially in the context
	//        of output parameters, at least a note regarding their
	//        interactions should be retained.
	@Test
	public void stringsAreImmutable() {
		String input = "the thing";
		this.tryAssigningAStringArgument(input, "will not be assigned");
		assertEquals(input, "the thing");
	}
	
	public void tryAssigningAStringArgument(String s, String assignment) {
		s = assignment;
	}
	
	@Test
	public void aMethodCantChangeItsParameters() {
		int output = 0;
		this.increment(output);
		assertEquals(output, 0);
	}
	
	public void increment(int i) {
		i += 1;
	}
	
	@Test
	public void thisActsLikeAnOutputArgument() {
		SelfChangingObject selfChangingObject = new SelfChangingObject("hello");
		selfChangingObject.change();
		assertEquals(selfChangingObject.value, "bye");
	}
	
	private class SelfChangingObject {
		String value;

		SelfChangingObject(String value) {
			this.value = value;
		}

		void change() {
			this.value = "bye";
		}
	}
}
