const CELL_SIZE = 25; // px
const SEP_WIDTH = 2;

const COLOR_WHITE = "#FFFFFF";
const COLOR_BLACK = "#000000";

const BOARD_WIDTH = 16;
const BOARD_HEIGHT = 16;

const canvas = document.getElementById("place");
canvas.height = (CELL_SIZE + SEP_WIDTH) * BOARD_HEIGHT + SEP_WIDTH;
canvas.width = (CELL_SIZE + SEP_WIDTH) * BOARD_WIDTH + SEP_WIDTH;

const ctx = canvas.getContext('2d');

// Fetch and draw the whole board.
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
        col * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
        row * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
        CELL_SIZE,
        CELL_SIZE
      );
    }
  }

  ctx.stroke();
}

// Trigger request on click.
canvas.addEventListener("click", async event => {
  const boundingRect = canvas.getBoundingClientRect();

  const scaleX = canvas.width / boundingRect.width;
  const scaleY = canvas.height / boundingRect.height;

  const canvasLeft = (event.clientX - boundingRect.left) * scaleX;
  const canvasTop = (event.clientY - boundingRect.top) * scaleY;

  // do nothing if click on grid, not cell
  const clickOnGrid =
    canvasLeft % (CELL_SIZE + SEP_WIDTH) < SEP_WIDTH ||
    canvasTop % (CELL_SIZE + SEP_WIDTH) < SEP_WIDTH;
  if (clickOnGrid) { return }

  const row = Math.floor(canvasTop / (CELL_SIZE + SEP_WIDTH));
  const col = Math.floor(canvasLeft / (CELL_SIZE + SEP_WIDTH));

  await fetch("http://" + document.location.host + "/toggle", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ "x": row, "y": col }),
  });
});

// Establish websocket connection.
const socket = new WebSocket("ws://" + document.location.host + "/ws");

socket.onmessage = (event) => {
  const data = JSON.parse(event.data);

  ctx.beginPath();

  ctx.fillStyle = data.color == 0
    ? COLOR_BLACK
    : COLOR_WHITE;

  ctx.fillRect(
    data.y * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
    data.x * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
    CELL_SIZE,
    CELL_SIZE
  );

  ctx.stroke();
}

// Draw the initial state of board once.
drawBoard();
