<script>
  let artsyMode = false;
  let isDrawing = false;
  let haveDrawn = false;
  let color = localStorage.getItem('color') || 'black';
  let lineWidth = parseInt(localStorage.getItem('lineWidth')) || 1;
  let lastEv = null;

  let undoing = false;
  let undoHistory = [];
  window.undoHistory = undoHistory;

  const startDrawing = (/** @type MouseEvent */ ev) => {
    if (ev.button != 0) { // is not left-click/main button
      return;
    }
    isDrawing = true;
    haveDrawn = false;
    undoing = false;

    ev.preventDefault()
  }

  const undo = () => {
    if (!undoing) {
      undoHistory.pop(); // discard last step, which was the last scribble, restore the one before
      undoing = true;
    }

    const prev = undoHistory.pop();

    const cv = document.querySelector("canvas");
    const ctx = cv.getContext("2d");

    if (prev) {
      ctx.putImageData(prev, 0, 0);
    } else { // initial state, no history
      ctx.clearRect(0, 0, cv.width, cv.height);
    }
  }

  const stopDrawing = (/** @type MouseEvent */ ev) => {
    if (ev.button != 0) { // is not left-click/main button
      return;
    }

    if (isDrawing && haveDrawn) {
      console.log("save", ev);
      const cv = document.querySelector("canvas");
      const ctx = cv.getContext("2d");   
      undoHistory.push(ctx.getImageData(0, 0, cv.width, cv.height));
    }

    isDrawing = false;
    lastEv = null;
    haveDrawn = false;

    ev.preventDefault()
  }

  const startTouchDrawing = (/** @type TouchEvent */ ev) => {
    isDrawing = true;
    lastEv = null;
    ev.preventDefault();
  };

  const stopTouchDrawing = (/** @type TouchEvent */ ev) => {
    isDrawing = true;
    lastEv = null;
    ev.preventDefault()
  };

  const drawTouch = (/** @type TouchEvent */ ev) => {
    const x = ev.touches[0].clientX - ev.target.offsetLeft;
    const y = ev.touches[0].clientY - ev.target.offsetTop;

    draw(x, y);

    if (artsyMode) {
      lastEv = { offsetX: 0, offsetY: 0 };
    } else {
      lastEv = { offsetX: x, offsetY: y };
    }

    ev.preventDefault();
  };

  const drawMouse = (/** @type MouseEvent */ ev) => {
    if (!isDrawing) {
      return;
    }
    if (ev.button != 0) {
      return;
    }

    const x = (ev.clientX - ev.target.offsetLeft) * (ev.target.width / ev.target.offsetWidth);
    const y = (ev.clientY - ev.target.offsetTop) * (ev.target.height / ev.target.offsetHeight);
    draw(x, y);

    if (artsyMode) {
      lastEv = { offsetX: 0, offsetY: 0 };
    } else {
      lastEv = { offsetX: x, offsetY: y };
    }
  };

  const draw = (x, y) => {
    const cv = document.querySelector("canvas");
    const ctx = cv.getContext("2d");
    // ctx.fillRect(ev.offsetX, ev.offsetY, 3, 3);

    if (lastEv) {
      ctx.lineWidth = lineWidth;
      ctx.lineJoin = "round";
      ctx.lineCap = "round";
      ctx.strokeStyle = color;
      ctx.beginPath()
      ctx.moveTo(lastEv.offsetX, lastEv.offsetY);
      ctx.lineTo(x, y);
      ctx.stroke();

      haveDrawn = true;
    }
  };

  const clearDrawing = () => {
    const cv = document.querySelector("canvas");
    const ctx = cv.getContext("2d");
    ctx.clearRect(0, 0, 1000, 1000);

    isDrawing = false;
  }

  const setMode = (/** @type Event */ ev) => {
    artsyMode = ev.target.checked;
    lastEv = null;
  }

  const setLineWidth = (ev) => {
    lineWidth = ev.target.value;
  }

  const setColor = (ev) => {
    color = ev.target.value;
    localStorage.setItem('color', color);
  }
</script>

<style>
  canvas {
    border: 1px solid #ddd;
    max-width: 100%;
    max-height: 100%;
  }
</style>

<canvas
  width="300" height="300"

  on:mousedown={startDrawing}
  on:mouseup={stopDrawing} on:mouseleave={stopDrawing}
  on:mousemove={drawMouse}
  
  on:touchstart={startTouchDrawing}
  on:touchend={stopTouchDrawing} on:touchcancel={stopTouchDrawing}
  on:touchmove={drawTouch}
  >
</canvas>

<div>
  <label for="artsy">artsy mode:</label>
  <input type="checkbox" name="artsy" on:change={setMode} />
</div>

<label for="width">line width ({lineWidth}):</label>
<input type="range" name="width" value={lineWidth} min="1" max="10" on:change={setLineWidth} />

<input type="color" name="color" value={color} on:change={setColor} />

<button on:click={undo}>Undo ↩</button>

<button on:click={clearDrawing}>Clear</button>
