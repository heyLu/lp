<script>
  let artsyMode = false;
  let isDrawing = false;
  let color = 'black';
  let lineWidth = 1;
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


  const draw = (/** @type MouseEvent */ ev) => {
    if (!isDrawing) {
      return;
    }

    const cv = document.querySelector("canvas");
    const ctx = cv.getContext("2d");
    // ctx.fillRect(ev.offsetX, ev.offsetY, 3, 3);

    if (lastEv) {
      // console.log(lastEv.offsetX, lastEv.offsetY, ev.offsetX, ev.offsetY);
      ctx.lineWidth = lineWidth;
      ctx.lineJoin = "round";
      ctx.lineCap = "round";
      ctx.strokeStyle = color;
      ctx.beginPath()
      ctx.moveTo(lastEv.offsetX, lastEv.offsetY);
      ctx.lineTo(ev.offsetX, ev.offsetY);
      ctx.stroke();
    }

    if (artsyMode) {
      lastEv = { offsetX: 0, offsetY: 0 };
    } else {
      lastEv = { offsetX: ev.offsetX, offsetY: ev.offsetY };
    }
  }

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
  on:mousemove={draw}>
</canvas>

<div>
  <label for="artsy">artsy mode:</label>
  <input type="checkbox" name="artsy" on:change={setMode} />
</div>

<label for="width">line width ({lineWidth}):</label>
<input type="range" name="width" value="1" min="1" max="10" on:change={setLineWidth} />

<input type="color" name="color" on:change={setColor} />

<button on:click={clearDrawing}>Clear</button>
