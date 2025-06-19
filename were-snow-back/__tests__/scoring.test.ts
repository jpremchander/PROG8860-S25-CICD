// Game scoring logic tests
describe("Scoring System", () => {
  function calculateScore(frameCount: number, baseScore = 0): number {
    return baseScore + Math.floor(frameCount / 60) * 10
  }

  function calculateSpeedMultiplier(timeElapsed: number): number {
    return 1 + Math.floor(timeElapsed / 2500) * 0.05
  }

  test("should calculate correct score based on frame count", () => {
    expect(calculateScore(0)).toBe(0)
    expect(calculateScore(60)).toBe(10)
    expect(calculateScore(120)).toBe(20)
    expect(calculateScore(300)).toBe(50)
  })

  test("should add base score to calculated score", () => {
    expect(calculateScore(60, 100)).toBe(110)
    expect(calculateScore(120, 50)).toBe(70)
  })

  test("should calculate speed multiplier correctly", () => {
    expect(calculateSpeedMultiplier(0)).toBe(1)
    expect(calculateSpeedMultiplier(2500)).toBe(1.05)
    expect(calculateSpeedMultiplier(5000)).toBe(1.1)
  })

  test("should handle edge cases in scoring", () => {
    expect(calculateScore(-10)).toBe(0)
    expect(calculateScore(59)).toBe(0)
    expect(calculateScore(61)).toBe(10)
  })
})
