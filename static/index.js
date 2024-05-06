const BOARD_WIDTH = 16;
const BOARD_HEIGHT = 16;

// https://lospec.com/palette-list/r-place
const COLORS = [
  "#FFFFFF",
  "#E4E4E4",
  "#888888",
  "#222222",
  "#FFA7D1",
  "#E50000",
  "#E59500",
  "#A06A42",
  "#E5D900",
  "#94E044",
  "#02BE01",
  "#00D3DD",
  "#0083C7",
  "#0000EA",
  "#CF6EE4",
  "#820080",
];

var color = 0;

$(document).ready(async function () {
  // Color pickers
  // TODO: The check is hard to see on darker colors.
  // Is there a better way to highlight current selection?
  $(".color").each(function (index) {
    $(this).css("background-color", COLORS[index]);
    $(this).click(function () {
      $(".color").eq(color).html(`&nbsp;`);
      color = index;
      $(this).html(`<i class="fa-solid fa-check"></i>`);
    });
  });

  // Set the canvas size to fill up the container as much as possible
  const canvas = document.getElementById("place");
  const cellSize = Math.min(
    Math.floor($(".canvas-container").width() / BOARD_WIDTH),
    Math.floor($(".canvas-container").height() / BOARD_HEIGHT),
  );
  canvas.width = cellSize * BOARD_WIDTH;
  canvas.height = cellSize * BOARD_HEIGHT;

  // Fetch the whole board
  // $.ajax doesn't handle raw bytes well, use fetch API instead
  const response = await fetch("/board");
  const buffer = await response.arrayBuffer();
  const board = new Uint8Array(buffer);

  // Draw the initial board
  const ctx = canvas.getContext("2d");
  ctx.beginPath();
  for (let row = 0; row < BOARD_HEIGHT; row++) {
    for (let col = 0; col < BOARD_WIDTH; col++) {
      const idx = row * BOARD_WIDTH + col;
      const b = board[Math.floor(idx / 2)]; // two colors packed in one byte
      const color = idx % 2 ? b & 15 : b >> 4;

      ctx.fillStyle = COLORS[color];
      ctx.fillRect(col * cellSize, row * cellSize, cellSize, cellSize);
    }
  }
  ctx.stroke();

  // Click on canvas to draw
  canvas.addEventListener("click", async (event) => {
    // Get the real clicking position
    const boundingRect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / boundingRect.width;
    const scaleY = canvas.height / boundingRect.height;
    const canvasLeft = (event.clientX - boundingRect.left) * scaleX;
    const canvasTop = (event.clientY - boundingRect.top) * scaleY;

    const row = Math.floor(canvasTop / cellSize);
    const col = Math.floor(canvasLeft / cellSize);
    $.post("pixel", { row, col, color });
  });

  // Establish websocket connection
  // FIX: This will not work on deployment
  const socket = new WebSocket("ws://" + document.location.host + "/ws");

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    ctx.beginPath();
    ctx.fillStyle = COLORS[data.color];
    ctx.fillRect(data.col * cellSize, data.row * cellSize, cellSize, cellSize);
    ctx.stroke();
  };
});
