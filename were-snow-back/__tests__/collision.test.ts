// Game logic tests for collision detection
describe("Collision Detection Logic", () => {
  const mockPlayer = {
    x: 100,
    y: 200,
    width: 32,
    height: 32,
  }

  const mockObstacle = {
    x: 100,
    y: 200,
    width: 32,
    height: 48,
  }

  function checkCollision(player: typeof mockPlayer, obstacle: typeof mockObstacle): boolean {
    const dx = Math.abs(player.x - obstacle.x)
    const dy = Math.abs(player.y - obstacle.y)
    return dx < player.width / 2 && dy < player.height / 2
  }

  test("should detect collision when player and obstacle overlap", () => {
    const result = checkCollision(mockPlayer, mockObstacle)
    expect(result).toBe(true)
  })

  test("should not detect collision when player and obstacle are far apart", () => {
    const distantObstacle = { ...mockObstacle, x: 300, y: 300 }
    const result = checkCollision(mockPlayer, distantObstacle)
    expect(result).toBe(false)
  })

  test("should not detect collision when player is above obstacle", () => {
    const aboveObstacle = { ...mockObstacle, y: mockPlayer.y - 50 }
    const result = checkCollision(mockPlayer, aboveObstacle)
    expect(result).toBe(false)
  })

  test("should detect collision at boundary conditions", () => {
    const boundaryObstacle = {
      ...mockObstacle,
      x: mockPlayer.x + mockPlayer.width / 2 - 1,
      y: mockPlayer.y,
    }
    const result = checkCollision(mockPlayer, boundaryObstacle)
    expect(result).toBe(true)
  })
})
