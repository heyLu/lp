<script>
  let isDrawing = false;
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
    ctx.fillRect(ev.offsetX, ev.offsetY, 3, 3);

    if (lastEv) {
      ctx.moveTo(lastEv.offsetX, lastEv.offsetY);
      ctx.lineTo(ev.offsetX, ev.offsetY);
      ctx.stroke();
    }

    lastEv = ev;
  }

  const clearDrawing = () => {
    const cv = document.querySelector("canvas");
    const ctx = cv.getContext("2d");
    ctx.clearRect(0, 0, 1000, 1000);

    isDrawing = false;
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

<button on:click={clearDrawing}>Clear</button>
