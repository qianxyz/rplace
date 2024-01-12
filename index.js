const CELL_SIZE = 20; // px

const COLOR_GRAY = "#CCCCCC";
const COLOR_WHITE = "#FFFFFF";
const COLOR_BLACK = "#000000";

const BOARD_WIDTH = 16;
const BOARD_HEIGHT = 16;

const canvas = document.getElementById("place");
canvas.height = (CELL_SIZE + 1) * BOARD_HEIGHT + 1;
canvas.width = (CELL_SIZE + 1) * BOARD_WIDTH + 1;

const ctx = canvas.getContext('2d');

function drawGrid() {
  ctx.beginPath();
  ctx.strokeStyle = COLOR_GRAY;

  // Vertical lines.
  for (let i = 0; i <= BOARD_WIDTH; i++) {
    ctx.moveTo(i * (CELL_SIZE + 1) + 1, 0);
    ctx.lineTo(i * (CELL_SIZE + 1) + 1, (CELL_SIZE + 1) * BOARD_HEIGHT + 1);
  }

  // Horizontal lines.
  for (let j = 0; j <= BOARD_HEIGHT; j++) {
    ctx.moveTo(0, j * (CELL_SIZE + 1) + 1);
    ctx.lineTo((CELL_SIZE + 1) * BOARD_WIDTH + 1, j * (CELL_SIZE + 1) + 1);
  }

  ctx.stroke();
};

async function drawBoard() {
  const response = await fetch("http://" + document.location.host + "/board");
  const buffer = await response.arrayBuffer();
  const board = new Uint8Array(buffer);

  ctx.beginPath();

  for (let row = 0; row < BOARD_HEIGHT; row++) {
    for (let col = 0; col < BOARD_WIDTH; col++) {
      const idx = row * BOARD_WIDTH + col;
      const cell = board[Math.floor(idx / 8)] & (1 << (7 - idx % 8));

      ctx.fillStyle = cell == 0
        ? COLOR_BLACK
        : COLOR_WHITE;

      ctx.fillRect(
        col * (CELL_SIZE + 1) + 1,
        row * (CELL_SIZE + 1) + 1,
        CELL_SIZE,
        CELL_SIZE
      );
    }
  }

  ctx.stroke();
}

canvas.addEventListener("click", async event => {
  const boundingRect = canvas.getBoundingClientRect();

  const scaleX = canvas.width / boundingRect.width;
  const scaleY = canvas.height / boundingRect.height;

  const canvasLeft = (event.clientX - boundingRect.left) * scaleX;
  const canvasTop = (event.clientY - boundingRect.top) * scaleY;

  const row = Math.min(Math.floor(canvasTop / (CELL_SIZE + 1)), BOARD_HEIGHT - 1);
  const col = Math.min(Math.floor(canvasLeft / (CELL_SIZE + 1)), BOARD_WIDTH - 1);

  await fetch("http://" + document.location.host + "/toggle", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ "x": row, "y": col }),
  });
});

drawGrid();
drawBoard();

const socket = new WebSocket("ws://" + document.location.host + "/ws");

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);

  ctx.beginPath();

  ctx.fillStyle = data.color == 0
    ? COLOR_BLACK
    : COLOR_WHITE;

  ctx.fillRect(
    data.y * (CELL_SIZE + 1) + 1,
    data.x * (CELL_SIZE + 1) + 1,
    CELL_SIZE,
    CELL_SIZE
  );

  ctx.stroke();
}
