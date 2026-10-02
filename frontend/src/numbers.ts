const numericInput = /^[+-]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][+-]?\d+)?$/

export type OperandError = 'invalid_format' | 'out_of_range'

type OperandParseResult =
  | { valid: true; value: number }
  | { valid: false; error: OperandError }

export function parseOperand(input: string): OperandParseResult {
  const value = input.trim()
  if (!numericInput.test(value)) return { valid: false, error: 'invalid_format' }

  const number = Number(value)
  if (!Number.isFinite(number)) return { valid: false, error: 'out_of_range' }
  return { valid: true, value: number }
}

export function formatResult(result: number): string {
  return Number(result.toPrecision(12)).toString()
}

export function needsApproximationNotice(result: number, operands: readonly number[]): boolean {
  const displayRounded = Number(formatResult(result)) !== result
  const integerPrecisionRisk = operands.some(
    (operand) => Number.isInteger(operand) && !Number.isSafeInteger(operand),
  )
  return displayRounded || integerPrecisionRisk
}
