<script>
  let artsyMode = false;
  let isDrawing = false;
  let color = localStorage.getItem('color') || 'black';
  let lineWidth = parseInt(localStorage.getItem('lineWidth')) || 1;
  let lastEv = null;

  const startDrawing = (/** @type MouseEvent */ ev) => {
    if (ev.button != 0) { // is not left-click/main button
      ev.preventDefault();
      return;
    }
    isDrawing = true;
  }

  const stopDrawing = (/** @type MouseEvent */ ev) => {
    if (ev.button != 0) { // is not left-click/main button
      ev.preventDefault();
      return;
    }
    isDrawing = false;
    lastEv = null;
  }

  const startTouchDrawing = (/** @type TouchEvent */ ev) => {
    isDrawing = true;
    lastEv = null;
  };

  const stopTouchDrawing = (/** @type TouchEvent */ ev) => {
    isDrawing = true;
    lastEv = null;
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

    draw(ev.offsetX, ev.offsetY);

    if (artsyMode) {
      lastEv = { offsetX: 0, offsetY: 0 };
    } else {
      lastEv = { offsetX: ev.offsetX, offsetY: ev.offsetY };
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
  }
</style>

<canvas
  width="300" height="300"

  on:mousedown={startDrawing}
  on:mouseup={stopDrawing} on:mouseleave={stopDrawing}
  on:mousemove={drawMouse}
  
  on:touchstart={startTouchDrawing}
  on:touchend={stopTouchDrawing}
  on:touchmove={drawTouch}
  >
</canvas>

<div>
  <label for="artsy">artsy mode:</label>
  <input type="checkbox" name="artsy" on:change={setMode} />
</div>

<label for="width">line width ({lineWidth}):</label>
<input type="range" name="width" value="1" min="1" max="10" on:change={setLineWidth} />

<input type="color" name="color" on:change={setColor} />

<button on:click={clearDrawing}>Clear</button>
