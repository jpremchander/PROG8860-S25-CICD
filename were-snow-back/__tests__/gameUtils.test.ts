import { GAME_CONSTANTS, COLORS, IMAGES } from "../constants"

describe("Game Constants", () => {
  test("should have correct canvas dimensions", () => {
    expect(GAME_CONSTANTS.CANVAS_WIDTH).toBe(800)
    expect(GAME_CONSTANTS.CANVAS_HEIGHT).toBe(400)
  })

  test("should have valid movement speed", () => {
    expect(GAME_CONSTANTS.MOVEMENT_SPEED).toBeGreaterThan(0)
    expect(GAME_CONSTANTS.MOVEMENT_SPEED).toBeLessThanOrEqual(5)
  })

  test("should have valid player dimensions", () => {
    expect(GAME_CONSTANTS.PLAYER_WIDTH).toBeGreaterThan(0)
    expect(GAME_CONSTANTS.PLAYER_HEIGHT).toBeGreaterThan(0)
  })

  test("should have valid gravity constant", () => {
    expect(GAME_CONSTANTS.GRAVITY).toBeGreaterThan(0)
    expect(GAME_CONSTANTS.GRAVITY).toBeLessThan(1)
  })
})

describe("Game Colors", () => {
  test("should have valid color values", () => {
    expect(COLORS.sky).toMatch(/^#[0-9A-Fa-f]{6}$/)
    expect(COLORS.snow).toMatch(/^#[0-9A-Fa-f]{6}$/)
    expect(COLORS.skiTrail).toMatch(/^#[0-9A-Fa-f]{6}$/)
  })
})

describe("Game Images", () => {
  test("should have player image URL", () => {
    expect(IMAGES.PLAYER).toContain("https://")
    expect(IMAGES.PLAYER).toContain("snowboarder")
  })

  test("should have tree images array", () => {
    expect(Array.isArray(IMAGES.TREES)).toBe(true)
    expect(IMAGES.TREES.length).toBeGreaterThan(0)
    IMAGES.TREES.forEach((url) => {
      expect(url).toContain("https://")
    })
  })

  test("should have snowman images array", () => {
    expect(Array.isArray(IMAGES.SNOWMEN)).toBe(true)
    expect(IMAGES.SNOWMEN.length).toBeGreaterThan(0)
    IMAGES.SNOWMEN.forEach((url) => {
      expect(url).toContain("https://")
    })
  })
})
