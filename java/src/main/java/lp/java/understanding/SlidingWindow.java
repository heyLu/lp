package lp.java.understanding;

import java.util.ArrayList;
import java.util.List;

public class SlidingWindow<T> {
    private static final long NEVER = -1;

    private List<Long> currentTimestamps;
    private List<T> currentValues;
    private List<Long> newTimestamps;
    private List<T> newValues;
    private long windowStart;
    private long windowSize;

    public SlidingWindow(long windowSize) {
        currentTimestamps = new ArrayList<Long>();
        currentValues = new ArrayList<T>();
        newTimestamps = new ArrayList<Long>();
        newValues = new ArrayList<T>();
        windowStart = NEVER;
        this.windowSize = windowSize;
    }

    public List<T> getCurrentValues() {
        return currentValues;
    }

    public void addValue(T value, long timestamp) {
        if (windowStart == NEVER) {
            windowStart = timestamp;
        }
        if (isInWindow(timestamp)) {
            currentTimestamps.add(timestamp);
            currentValues.add(value);
        } else {
            newTimestamps.add(timestamp);
            newValues.add(value);
        }
    }

    public long getWindowStart() {
        return windowStart;
    }

    public boolean isInWindow(long timestamp) {
        return timestamp >= windowStart && timestamp < windowStart + windowSize;
    }

    public long getWindowSize() {
        return windowSize;
    }

    public long getWindowEnd() {
        return windowStart + windowSize;
    }

    public boolean canSlide() {
        return !newValues.isEmpty();
    }

    public void slide(long slideBy) {
        windowStart += slideBy;
        for (int i = 0; i < currentValues.size(); i++) {
            long timestamp = currentTimestamps.get(i);
            if (!isInWindow(timestamp)) {
                currentTimestamps.remove(i);
                currentValues.remove(i);
                i -= 1;
            }
        }

        for (int i = 0; i < newValues.size(); i++) {
            long timestamp = newTimestamps.get(i);
            if (isInWindow(timestamp)) {
                T value = newValues.get(i);
                currentTimestamps.add(timestamp);
                currentValues.add(value);
                newTimestamps.remove(i);
                newValues.remove(i);
                i -= 1;
            }
        }
    }
}
