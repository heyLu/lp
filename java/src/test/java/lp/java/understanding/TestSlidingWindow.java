package lp.java.understanding;

import static org.junit.Assert.*;
import org.junit.Before;
import org.junit.Test;

public class TestSlidingWindow {
    private SlidingWindow<Integer> window;

    @Before
    public void setUp() {
        window = new SlidingWindow<Integer>(10);
    }

    @Test
    public void initiallyTheSlidingWindowHoldsNoValues() {
        assertTrue(window.getCurrentValues().isEmpty());
    }

    @Test
    public void testIsInWindow() {
        window.addValue(0, 0);
        assertFalse("before window start", window.isInWindow(window.getWindowStart() - 1));
        assertTrue("at window start", window.isInWindow(window.getWindowStart()));
        assertTrue("after window start", window.isInWindow(window.getWindowStart() + 1));
        assertFalse("excluding window end", window.isInWindow(window.getWindowEnd()));
        assertFalse("after window end", window.isInWindow(window.getWindowEnd() + 10));
    }

    @Test
    public void addingAValueInTheCurrentWindowAddsItImmediately() {
        window.addValue(0, 0);
        window.addValue(1, window.getWindowStart());
        assertTrue(window.getCurrentValues().contains(1));
    }

    @Test
    public void addingAValueNotInTheCurrentWindowDoesntAddItImmediately() {
        window.addValue(0, 0);
        window.addValue(1, window.getWindowEnd());
        assertFalse(window.getCurrentValues().contains(1));
    }

    @Test
    public void theFirstAddedValueDeterminesTheWindowStart() {
        window.addValue(1, 42);
        assertEquals(42, window.getWindowStart());
    }

    @Test
    public void addingAValueOutsideOfTheCurrentWindowMakesItSlideable() {
        window.addValue(1, 0);
        assertFalse(window.canSlide());
        window.addValue(2, window.getWindowEnd());
        assertTrue(window.canSlide());
    }

    @Test
    public void slidingRemovesValuesNotInTheCurrentWindowAnymore() {
        addManyTimedValues(1, 0l, 2, 1l, 3, 5l, 4, 8l);
        assertWindowContains(1, 2, 3, 4);
        window.slide(2);
        assertWindowContains(3, 4);
    }

    private void addManyTimedValues(Object ...valuesWithTimestamps) {
        if (valuesWithTimestamps.length % 2 != 0) {
            throw new IllegalArgumentException("must pass even number of arguments");
        }
        for (int i = 0; i < valuesWithTimestamps.length; i += 2) {
            Integer value = (Integer) valuesWithTimestamps[i];
            Long timestamp = (Long) valuesWithTimestamps[i + 1];
            window.addValue(value, timestamp);
        }
    }

    @Test
    public void slidingAddsNewValuesFromTheNewWindow() {
        addManyTimedValues(1, 0l, 2, 1l, 3, 5l, 4, 8l, 5, 11l, 6, 20l);
        assertWindowContains(1, 2, 3, 4);
        window.slide(2);
        assertWindowContains(3, 4, 5);
        window.slide(10);
        assertWindowContains(6);
    }

    private void assertWindowContains(Integer ...values) {
        for (Integer value : values) {
            assertTrue("contains " + value, window.getCurrentValues().contains(value));
        }
    }
}
