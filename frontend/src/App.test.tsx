import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { afterEach, beforeEach, expect, test, vi } from 'vitest'
import App from './App'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  fetchMock.mockReset()
  vi.stubGlobal('fetch', fetchMock)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.useRealTimers()
})

function enterOperands(left: string, right: string) {
  fireEvent.change(screen.getByLabelText('First number'), { target: { value: left } })
  fireEvent.change(screen.getByLabelText('Second number'), { target: { value: right } })
}

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

test('renders the calculator application in the main landmark', () => {
  render(<App />)

  const main = screen.getByRole('main')
  expect(
    within(main).getByRole('heading', { name: 'Calculator', level: 1 }),
  ).toBeInTheDocument()
})

test.each([
  ['', '2', 'First number has an invalid format. Use a number such as -2.5 or 1e3.'],
  ['2', '  ', 'Second number has an invalid format. Use a number such as -2.5 or 1e3.'],
  ['1,5', '2', 'First number has an invalid format. Use a number such as -2.5 or 1e3.'],
  ['2', '5abc', 'Second number has an invalid format. Use a number such as -2.5 or 1e3.'],
  ['9'.repeat(400), '2', 'First number is outside the supported numeric range. Enter a smaller absolute value.'],
  ['2', '-1e309', 'Second number is outside the supported numeric range. Enter a smaller absolute value.'],
])('identifies the invalid field and reason (case %#) without calling the API', (left, right, message) => {
  render(<App />)
  enterOperands(left, right)
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(screen.getByRole('alert')).toHaveTextContent(message)
  expect(fetchMock).not.toHaveBeenCalled()
  const field = message.startsWith('First') ? 'First number' : 'Second number'
  const invalidInput = screen.getByLabelText(field)
  const validInput = screen.getByLabelText(field === 'First number' ? 'Second number' : 'First number')
  expect(invalidInput).toHaveFocus()
  expect(invalidInput).toHaveAttribute('aria-invalid', 'true')
  expect(invalidInput).toHaveAccessibleDescription(expect.stringContaining(message))
  expect(validInput).toHaveAttribute('aria-invalid', 'false')
  expect(validInput).toHaveAttribute('aria-describedby', 'number-format')
  fireEvent.change(invalidInput, { target: { value: '1.5' } })
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  expect(invalidInput).toHaveAttribute('aria-invalid', 'false')
  expect(invalidInput).toHaveAttribute('aria-describedby', 'number-format')
})

test('sends the entered operands and shows the API result', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 10 }))
  render(<App />)
  enterOperands(' 8 ', '2')
  expect(screen.getByRole('combobox', { name: 'Operation' })).toHaveValue('add')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^10$/)
  expect(screen.queryByText(/^Approximate result\./)).not.toBeInTheDocument()
  await waitFor(() => expect(screen.getByRole('button', { name: 'Calculate' })).toHaveFocus())
  expect(fetchMock).toHaveBeenCalledWith('/api/add', expect.objectContaining({
    body: '{"left":8,"right":2}',
  }))
  fireEvent.change(screen.getByLabelText('Second number'), { target: { value: '3' } })
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test('clears output when selecting subtraction and sends operands in order', async () => {
  fetchMock.mockResolvedValueOnce(jsonResponse({ result: 10 }))
    .mockResolvedValueOnce(jsonResponse({ result: -6 }))
  render(<App />)
  enterOperands('2', '8')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  await screen.findByRole('status', { name: 'Result' })

  const operation = screen.getByRole('combobox', { name: 'Operation' })
  fireEvent.change(operation, { target: { value: 'subtract' } })
  expect(operation).toHaveValue('subtract')
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^-6$/)
  expect(fetchMock).toHaveBeenLastCalledWith('/api/subtract', expect.objectContaining({
    body: '{"left":2,"right":8}',
  }))

  fireEvent.change(operation, { target: { value: 'add' } })
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test('clears a subtraction error when the operation changes', async () => {
  fetchMock.mockResolvedValue(jsonResponse({
    error: { code: 'non_finite_result', message: 'The calculation result must be finite.' },
  }, 400))
  render(<App />)
  enterOperands('1e308', '-1e308')
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'subtract' } })
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('alert')).toHaveTextContent('The calculation result must be finite.')
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  await waitFor(() => expect(screen.getByRole('button', { name: 'Calculate' })).toHaveFocus())
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'add' } })
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
})

test('selects multiplication, sends signed operands, and clears its result on operation change', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: -16 }))
  render(<App />)
  enterOperands('-8', '2')
  fireEvent.change(screen.getByRole('combobox', { name: 'Operation' }), {
    target: { value: 'multiply' },
  })
  expect(screen.getByRole('combobox')).toHaveValue('multiply')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^-16$/)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/multiply', expect.objectContaining({
    body: '{"left":-8,"right":2}',
  }))
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'add' } })
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test('selects division and preserves operand order', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 0.25 }))
  render(<App />)
  enterOperands('2', '8')
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'divide' } })
  expect(screen.getByRole('combobox')).toHaveValue('divide')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^0.25$/)
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/divide', expect.objectContaining({
    body: '{"left":2,"right":8}',
  }))
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'multiply' } })
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test.each(['0', '-0', '1e-999', '-1e-999'])('shows the backend division-by-zero error for divisor %s and allows retry', async (divisor) => {
  fetchMock.mockResolvedValueOnce(jsonResponse({
    error: { code: 'division_by_zero', message: 'Cannot divide by zero.' },
  }, 400)).mockResolvedValueOnce(jsonResponse({ result: 4 }))
  render(<App />)
  enterOperands('8', divisor)
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'divide' } })
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('alert')).toHaveTextContent('Cannot divide by zero.')
  expect(fetchMock).toHaveBeenCalledExactlyOnceWith('/api/divide', expect.objectContaining({
    body: '{"left":8,"right":0}',
  }))
  expect(screen.getByLabelText('First number')).toBeEnabled()
  expect(screen.getByLabelText('Second number')).toBeEnabled()
  expect(screen.getByRole('combobox')).toBeEnabled()
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  await waitFor(() => expect(screen.getByRole('button', { name: 'Calculate' })).toHaveFocus())
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
  fireEvent.change(screen.getByLabelText('Second number'), { target: { value: '2' } })
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^4$/)
  expect(fetchMock).toHaveBeenCalledTimes(2)
})

test.each([
  ['0.1', '0.2', 0.30000000000000004, '≈ 0.3'],
  ['-8', '8', 0, '0'],
] as const)('displays %s + %s: API %s, UI %s', async (left, right, result, display) => {
  fetchMock.mockResolvedValue(jsonResponse({ result }))
  render(<App />)
  enterOperands(left, right)
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveProperty('textContent', display)
})

test('explains presentation rounding and clears the notice when an operand changes', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 0.3333333333333333 }))
  render(<App />)
  enterOperands('1', '3')
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'divide' } })
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  const output = await screen.findByRole('status', { name: 'Result' })
  expect(output).toHaveTextContent(/^≈ 0\.333333333333$/)
  expect(output).toHaveAccessibleDescription('Approximate result. Large numbers can lose precision; results are shown with up to 12 significant digits.')
  fireEvent.change(screen.getByLabelText('Second number'), { target: { value: '4' } })
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
  expect(screen.queryByText(/^Approximate result\./)).not.toBeInTheDocument()
})

test('warns for a large input even when the API result does not need display rounding', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: 1e20 }))
  render(<App />)
  enterOperands('100000000000000000001', '0')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^≈ 100000000000000000000$/)
  expect(screen.getByText(/^Approximate result\./)).toBeInTheDocument()
  expect(fetchMock).toHaveBeenCalledWith('/api/add', expect.objectContaining({
    body: '{"left":100000000000000000000,"right":0}',
  }))
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'multiply' } })
  expect(screen.queryByText(/^Approximate result\./)).not.toBeInTheDocument()
})

test('clears approximation feedback on a new request and keeps it cleared after failure', async () => {
  let rejectRequest!: (reason: Error) => void
  fetchMock.mockResolvedValueOnce(jsonResponse({ result: 0.30000000000000004 }))
    .mockImplementationOnce(() => new Promise((_, reject) => { rejectRequest = reject }))
  render(<App />)
  enterOperands('0.1', '0.2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  await screen.findByRole('status', { name: 'Result' })

  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
  expect(screen.queryByText(/^Approximate result\./)).not.toBeInTheDocument()
  await act(async () => { rejectRequest(new TypeError('Failed to fetch')) })
  expect(await screen.findByRole('alert')).toHaveTextContent('The calculator service is unavailable.')
  expect(screen.queryByText(/^Approximate result\./)).not.toBeInTheDocument()
})

test('disables controls while pending and prevents a duplicate submission', async () => {
  let resolveResponse!: (response: Response) => void
  fetchMock.mockReturnValue(new Promise((resolve) => { resolveResponse = resolve }))
  render(<App />)
  enterOperands('8', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(screen.getByLabelText('First number')).toBeDisabled()
  expect(screen.getByLabelText('Second number')).toBeDisabled()
  expect(screen.getByRole('combobox')).toBeDisabled()
  const button = screen.getByRole('button', { name: 'Calculating…' })
  expect(button).toBeDisabled()
  fireEvent.submit(button.closest('form')!)
  expect(fetchMock).toHaveBeenCalledTimes(1)

  await act(async () => { resolveResponse(jsonResponse({ result: 10 })) })
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^10$/)
  expect(screen.getByLabelText('First number')).toBeEnabled()
  expect(screen.getByLabelText('Second number')).toBeEnabled()
  expect(screen.getByRole('combobox')).toBeEnabled()
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
})

test('shows an API error and supports a manual retry', async () => {
  fetchMock.mockResolvedValueOnce(jsonResponse({
    error: { code: 'non_finite_result', message: 'The calculation result must be finite.' },
  }, 400)).mockResolvedValueOnce(jsonResponse({ result: 10 }))
  render(<App />)
  enterOperands('1e308', '1e308')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('alert')).toHaveTextContent('The calculation result must be finite.')
  expect(fetchMock).toHaveBeenCalledTimes(1)
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  enterOperands('8', '2')
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent(/^10$/)
  expect(fetchMock).toHaveBeenCalledTimes(2)
})

test('shows a network failure and restores the controls', async () => {
  fetchMock.mockRejectedValue(new TypeError('Failed to fetch'))
  render(<App />)
  enterOperands('8', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))

  expect(await screen.findByRole('alert')).toHaveTextContent('The calculator service is unavailable.')
  expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'false')
  expect(screen.getByLabelText('Second number')).toHaveAttribute('aria-invalid', 'false')
  expect(screen.getByLabelText('First number')).toBeEnabled()
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test('clears field validation feedback when the operation changes', () => {
  render(<App />)
  enterOperands('1,5', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'true')
  fireEvent.change(screen.getByRole('combobox'), { target: { value: 'multiply' } })
  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  expect(screen.getByLabelText('First number')).toHaveAttribute('aria-invalid', 'false')
  expect(screen.getByLabelText('First number')).toHaveAttribute('aria-describedby', 'number-format')
  expect(fetchMock).not.toHaveBeenCalled()
})

test('does not display a malformed API result', async () => {
  fetchMock.mockResolvedValue(jsonResponse({ result: '10' }))
  render(<App />)
  enterOperands('8', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  expect(await screen.findByRole('alert')).toHaveTextContent('The server returned an invalid response.')
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})

test('shows the timeout and restores the form after 10 seconds', async () => {
  vi.useFakeTimers()
  fetchMock.mockImplementation((_, options) => new Promise((_, reject) => {
    options?.signal?.addEventListener('abort', () => reject(new DOMException('Aborted', 'AbortError')))
  }))
  render(<App />)
  enterOperands('8', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  await act(async () => { await vi.advanceTimersByTimeAsync(10_000) })

  expect(screen.getByRole('alert')).toHaveTextContent('The request timed out.')
  expect(screen.getByRole('button', { name: 'Calculate' })).toBeEnabled()
  expect(screen.getByRole('button', { name: 'Calculate' })).toHaveFocus()
  expect(screen.getByLabelText('First number')).toBeEnabled()
  expect(fetchMock).toHaveBeenCalledTimes(1)
})

test('clears the previous result when a later request fails', async () => {
  fetchMock.mockResolvedValueOnce(jsonResponse({ result: 10 }))
    .mockRejectedValueOnce(new TypeError('Failed to fetch'))
  render(<App />)
  enterOperands('8', '2')
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  await screen.findByRole('status', { name: 'Result' })
  fireEvent.click(screen.getByRole('button', { name: 'Calculate' }))
  await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('The calculator service is unavailable.'))
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
})
