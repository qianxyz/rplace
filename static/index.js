const BOARD_WIDTH = 16;
const BOARD_HEIGHT = 16;

$(document).ready(function () {
  const canvas = document.getElementById("place");
  const ctx = canvas.getContext("2d");

  // Set the canvas size to fill up the container as much as possible
  const cellSize = Math.min(
    Math.floor($(".canvas-container").width() / BOARD_WIDTH),
    Math.floor($(".canvas-container").height() / BOARD_HEIGHT),
  );
  canvas.width = cellSize * BOARD_WIDTH;
  canvas.height = cellSize * BOARD_HEIGHT;
});

// const API = "http://" + document.location.host;
// const WS = "ws://" + document.location.host + "/ws";
//
// // Fetch and draw the whole board.
// async function drawBoard() {
//   const response = await fetch(API + "/board");
//   const buffer = await response.arrayBuffer();
//   const board = new Uint8Array(buffer);
//
//   ctx.beginPath();
//
//   for (let row = 0; row < BOARD_HEIGHT; row++) {
//     for (let col = 0; col < BOARD_WIDTH; col++) {
//       const idx = row * BOARD_WIDTH + col;
//       const cell = board[Math.floor(idx / 8)] & (1 << (7 - (idx % 8)));
//
//       ctx.fillStyle = cell == 0 ? COLOR_BLACK : COLOR_WHITE;
//
//       ctx.fillRect(
//         col * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
//         row * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
//         CELL_SIZE,
//         CELL_SIZE,
//       );
//     }
//   }
//
//   ctx.stroke();
// }
//
// // Trigger request on click.
// canvas.addEventListener("click", async (event) => {
//   const boundingRect = canvas.getBoundingClientRect();
//
//   const scaleX = canvas.width / boundingRect.width;
//   const scaleY = canvas.height / boundingRect.height;
//
//   const canvasLeft = (event.clientX - boundingRect.left) * scaleX;
//   const canvasTop = (event.clientY - boundingRect.top) * scaleY;
//
//   // do nothing if click on grid, not cell
//   const clickOnGrid =
//     canvasLeft % (CELL_SIZE + SEP_WIDTH) < SEP_WIDTH ||
//     canvasTop % (CELL_SIZE + SEP_WIDTH) < SEP_WIDTH;
//   if (clickOnGrid) {
//     return;
//   }
//
//   const row = Math.floor(canvasTop / (CELL_SIZE + SEP_WIDTH));
//   const col = Math.floor(canvasLeft / (CELL_SIZE + SEP_WIDTH));
//
//   await fetch(API + "/toggle", {
//     method: "POST",
//     headers: { "Content-Type": "application/json" },
//     body: JSON.stringify({ row: row, col: col }),
//   });
// });
//
// // Establish websocket connection.
// const socket = new WebSocket(WS);
//
// socket.onmessage = (event) => {
//   const data = JSON.parse(event.data);
//
//   ctx.beginPath();
//
//   ctx.fillStyle = data.color == 0 ? COLOR_BLACK : COLOR_WHITE;
//
//   ctx.fillRect(
//     data.col * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
//     data.row * (CELL_SIZE + SEP_WIDTH) + SEP_WIDTH,
//     CELL_SIZE,
//     CELL_SIZE,
//   );
//
//   ctx.stroke();
// };
//
// // Draw the initial state of board once.
// drawBoard();
