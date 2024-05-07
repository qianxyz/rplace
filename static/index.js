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

const COOLDOWN = 5000; // 5s

var color = 0;
var isCoolingDown = false;

$(document).ready(async function () {
  // Color pickers
  $(".color").each(function (index) {
    const hex = COLORS[index];
    $(this).css("background-color", hex);

    // Calculate the color of the tick symbol
    const r = parseInt(hex.slice(1, 3), 16),
      g = parseInt(hex.slice(3, 5), 16),
      b = parseInt(hex.slice(5, 7), 16);
    // https://stackoverflow.com/a/3943023/112731
    const tickColor =
      r * 0.299 + g * 0.587 + b * 0.114 > 186 ? "#000000" : "#FFFFFF";

    $(this).click(function () {
      $(".color").eq(color).html(`&nbsp;`);
      color = index;
      $(this).html(
        `<i class="fa-solid fa-check" style="color: ${tickColor}"></i>`,
      );
    });
  });
  // Init click on the default color to show the tick
  $(".color").eq(color).trigger("click");

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
  const response = await fetch("board");
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
    if (isCoolingDown) {
      return;
    }
    isCoolingDown = true;
    setTimeout(() => (isCoolingDown = false), COOLDOWN);

    // Cooldown animation
    $("#cooldownBar").width("0%");
    $("#cooldownBar").animate({ width: "100%" }, COOLDOWN);

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
  const url = new URL("ws", document.baseURI);
  url.protocol = url.protocol.replace("http", "ws");
  const socket = new WebSocket(url);

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    ctx.beginPath();
    ctx.fillStyle = COLORS[data.color];
    ctx.fillRect(data.col * cellSize, data.row * cellSize, cellSize, cellSize);
    ctx.stroke();
  };
});
