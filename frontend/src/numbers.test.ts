import { expect, test } from 'vitest'
import { formatResult, needsApproximationNotice, parseOperand } from './numbers'

test.each([
  ['0', 0], ['-8', -8], ['1.25', 1.25], ['1e3', 1000],
  ['-2.5e-2', -0.025], ['  +.5  ', 0.5], ['1.', 1],
  ['1e308', 1e308], ['5e-324', Number.MIN_VALUE],
  ['1.7976931348623157e308', Number.MAX_VALUE],
  ['-1.7976931348623157e308', -Number.MAX_VALUE],
  ['1e-999', 0],
  ['-0', -0], ['-1e-999', -0], ['+.5E+2', 50],
  ['0008.00', 8], ['\t\n-2.5\r ', -2.5],
] as const)('parses valid input %s', (input, expected) => {
  expect(parseOperand(input)).toEqual({ valid: true, value: expected })
})

test.each([
  ['ordinary integers', 10, [8, 2], false],
  ['exact binary fraction', 0.25, [2, 8], false],
  ['presentation rounding', 0.30000000000000004, [0.1, 0.2], true],
  ['a repeating fraction', 0.3333333333333333, [1, 3], true],
  ['the safe integer boundary', 0, [9007199254740991, -9007199254740991], false],
  ['an unsafe left integer without display rounding', 1e20, [1e20, 0], true],
  ['an unsafe right integer', -1, [9007199254740991, -9007199254740992], true],
  ['a large negative operand', -1e20, [-1e20, 0], true],
] as const)('identifies approximation feedback for %s', (_, result, operands, expected) => {
  expect(needsApproximationNotice(result, operands)).toBe(expected)
})

test.each([
  '', '  ', '1,5', '1 000', '0x10', '0b11', '5abc',
  'NaN', 'Infinity', '1e', '--1',
  '.', '+', '-', '-.', '1e+', '1e-', '1e2e3', '1.2.3',
  '1_000', '0o10', '2+2', '1/0', '1\n2', '1\u00002',
  '−2', '１２', '١٢', '1\u200b', '<script>alert(1)</script>',
])('identifies invalid numeric format %j', (input) => {
  expect(parseOperand(input)).toEqual({ valid: false, error: 'invalid_format' })
})

test.each([
  ['positive exponent overflow', '1e309'],
  ['negative exponent overflow', '-1e309'],
  ['beyond the maximum finite boundary', '1.7976931348623159e308'],
  ['a long pasted digit sequence', '9'.repeat(400)],
])('identifies a number outside the supported range: %s', (_, input) => {
  expect(parseOperand(input)).toEqual({ valid: false, error: 'out_of_range' })
})

test.each([
  [0, '0'], [-2.5, '-2.5'], [0.30000000000000004, '0.3'],
  [1.23456789012345, '1.23456789012'], [1e30, '1e+30'], [1e-20, '1e-20'],
  [-0, '0'], [Number.MIN_VALUE, '5e-324'],
  [Number.MAX_VALUE, '1.79769313486e+308'],
  [-Number.MAX_VALUE, '-1.79769313486e+308'],
] as const)('formats result %s as %s', (result, expected) => {
  expect(formatResult(result)).toBe(expected)
})
