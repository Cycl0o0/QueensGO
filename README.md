# QueensGO

A puzzle game adapted from the famous **Queens** game on LinkedIn. Place one queen per row, column, and color region on the board — no two queens may touch, even diagonally.

Built entirely in Go with [Ebitengine](https://ebitengine.org/) for the GUI.

## Rules

1. Place exactly one queen in each row, column, and colored region
2. No two queens may be adjacent (horizontally, vertically, or diagonally)
3. Click a cell to cycle through: empty → X (elimination mark) → Queen → empty

## Difficulty Levels

| Difficulty | Grid Size |
|------------|-----------|
| Easy       | 5×5       |
| Normal     | 6×6       |
| Hard       | 7×7       |
| Expert     | 8×8       |
| Extreme    | 9×9       |

All puzzles are generated algorithmically with a guaranteed unique solution.

## Build & Run

```bash
go build -o queensgo .
./queensgo
```

Or simply:

```bash
go run .
```

### Requirements

- Go 1.21+
- On Linux: `sudo apt install libgl1-mesa-dev xorg-dev` (for Ebitengine)

## Licence

This project is licensed under the [GNU Affero General Public License v3.0](LICENCE).

---

Made with <3 by [Cycl0o0](https://github.com/Cycl0o0)
